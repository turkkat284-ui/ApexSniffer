use pcap::{Capture, Device};
use anyhow::{Context, Result};
use std::time::{SystemTime, UNIX_EPOCH, Duration};
use std::io::{Write, BufWriter};
use tokio::sync::mpsc;
use bytes::Bytes;

#[cfg(windows)]
use std::os::windows::io::AsRawHandle;
#[cfg(windows)]
use std::fs::OpenOptions;

#[cfg(unix)]
use std::os::unix::net::UnixStream;

/// Represents a captured packet with its timestamp and raw payload.
#[derive(Debug)]
pub struct PacketData {
    pub timestamp: SystemTime,
    pub payload: bytes::Bytes,
}

/// Trait for the IPC writer to abstract over Windows Named Pipes and Unix Domain Sockets.
pub trait IpcWriter: Write + Send + Sync {
    fn connect() -> Result<Self> where Self: Sized;
}

/// Windows implementation using Named Pipes.
#[cfg(windows)]
pub struct WindowsNamedPipeWriter {
    file: std::fs::File,
}

#[cfg(windows)]
impl IpcWriter for WindowsNamedPipeWriter {
    fn connect() -> Result<Self> {
        // Note: In a real scenario, the Go analyzer would create the pipe.
        // Here we attempt to open a pipe named \\.\pipe\apex_sniffer
        let pipe_path = r"\\.\pipe\apex_sniffer";
        let file = OpenOptions::new()
            .write(true)
            .open(pipe_path)
            .with_context(|| format!("Failed to open named pipe: {}. Make sure the Go analyzer is running.", pipe_path))?;

        Ok(WindowsNamedPipeWriter { file })
    }
}

#[cfg(windows)]
impl Write for WindowsNamedPipeWriter {
    fn write(&mut self, buf: &[u8]) -> std::io::Result<usize> {
        self.file.write(buf)
    }
    fn flush(&mut self) -> std::io::Result<()> {
        self.file.flush()
    }
}

/// Linux/Unix implementation using Unix Domain Sockets.
#[cfg(unix)]
pub struct UnixSocketWriter {
    stream: UnixStream,
}

#[cfg(unix)]
impl IpcWriter for UnixSocketWriter {
    fn connect() -> Result<Self> {
        let socket_path = "/tmp/apex_sniffer.sock";
        let stream = UnixStream::connect(socket_path)
            .with_context(|| format!("Failed to connect to Unix socket: {}. Make sure the Go analyzer is running.", socket_path))?;

        Ok(UnixSocketWriter { stream })
    }
}

#[cfg(unix)]
impl Write for UnixSocketWriter {
    fn write(&mut self, buf: &[u8]) -> std::io::Result<usize> {
        self.stream.write(buf)
    }
    fn flush(&mut self) -> std::io::Result<()> {
        self.stream.flush()
    }
}

/// Handles the framing and writing of packets to the IPC channel.
pub struct PacketStreamer<W: Write> {
    writer: BufWriter<W>,
}

impl<W: Write> PacketStreamer<W> {
    pub fn new(writer: W) -> Self {
        Self {
            writer: BufWriter::new(writer),
        }
    }

    pub fn write_packet(&mut self, packet: &PacketData) -> Result<()> {
        // Magic: 0x4150 (2 bytes)
        let magic: u16 = 0x4150;

        // Timestamp: i64 (8 bytes) - nanoseconds since epoch
        let timestamp = packet.timestamp
            .duration_since(UNIX_EPOCH)
            .context("Time went backwards")?
            .as_nanos() as i64;

        // Length: u32 (4 bytes)
        let length = packet.payload.len() as u32;

        // Build frame
        // Use big-endian (network byte order) for binary protocol
        self.writer.write_all(&magic.to_be_bytes())?;
        self.writer.write_all(&timestamp.to_be_bytes())?;
        self.writer.write_all(&length.to_be_bytes())?;
        self.writer.write_all(&packet.payload)?;

        Ok(())
    }

    pub fn flush(&mut self) -> Result<()> {
        self.writer.flush().context("Failed to flush IPC writer")
    }
}

/// The CaptureEngine handles the interaction with the network interface.
pub struct CaptureEngine {
    capture: Capture,
}

impl CaptureEngine {
    pub fn new(device_name: &str, filter: &str) -> Result<Self> {
        let device = Device::list()
            .context("Failed to list devices")?
            .into_iter()
            .find(|d| d.name == device_name)
            .context(format!("Device {} not found", device_name))?;

        let mut capture = Capture::from_device(device)
            .context("Failed to create capture from device")?
            .promisc(true)
            .snaplen(65535)
            .timeout(1000)
            .open()
            .context("Failed to open capture device")?;

        capture.filter(filter, true)
            .context(format!("Failed to apply BPF filter: {}", filter))?;

        Ok(CaptureEngine { capture })
    }

    pub fn run<F>(&mut self, mut handler: F) -> Result<()>
    where
        F: FnMut(PacketData),
    {
        println!("Starting packet capture loop...");
        loop {
            match self.capture.next_packet() {
                Ok(packet) => {
                    let data = PacketData {
                        timestamp: SystemTime::now(),
                        payload: Bytes::copy_from_slice(packet.data),
                    };
                    handler(data);
                }
                Err(pcap::Error::TimeoutExpired) => {
                    continue;
                }
                Err(e) => {
                    return Err(anyhow::anyhow!("Packet capture error: {:?}", e));
                }
            }
        }
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    env_logger::init();

    // 1. List devices
    let devices = Device::list().context("Failed to list devices")?;
    println!("Available devices:");
    for (i, device) in devices.iter().enumerate() {
        println!("{}: {}", i, device.name);
    }

    if devices.is_empty() {
        return Err(anyhow::anyhow!("No network devices found"));
    }

    let selected_device = &devices[0].name;
    println!("Selected device: {}", selected_device);

    let filter = "udp";
    println!("Applying BPF filter: {}", filter);

    // 2. Create a channel for decoupling capture and writing
    let (tx, mut rx) = mpsc::channel::<PacketData>(1024);

    // 3. Spawn a separate task to consume the channel and write to IpcWriter
    tokio::spawn(async move {
        println!("Connecting to IPC channel...");
        let writer: Box<dyn IpcWriter> = {
            #[cfg(windows)]
            {
                Box::new(WindowsNamedPipeWriter::connect().expect("Failed to connect to Windows Named Pipe"))
            }
            #[cfg(unix)]
            {
                Box::new(UnixSocketWriter::connect().expect("Failed to connect to Unix Socket"))
            }
        };

        let mut streamer = PacketStreamer::new(writer);
        println!("Connected to IPC channel. Streaming binary frames...");

        let mut last_flush = tokio::time::Instant::now();
        let flush_interval = Duration::from_millis(10);

        while let Some(packet) = rx.recv().await {
            if let Err(e) = streamer.write_packet(&packet) {
                eprintln!("Failed to write packet to IPC: {:?}", e);
            }

            if last_flush.elapsed() >= flush_interval {
                if let Err(e) = streamer.flush() {
                    eprintln!("Failed to flush IPC writer: {:?}", e);
                }
                last_flush = tokio::time::Instant::now();
            }
        }
    });

    // 4. Initialize CaptureEngine
    let mut engine = CaptureEngine::new(selected_device, filter)?;

    // 5. Run the loop and send to channel
    engine.run(|packet| {
        if let Err(e) = tx.try_send(packet) {
            eprintln!("Failed to send packet to channel: {:?}", e);
        }
    })?;

    Ok(())
}
