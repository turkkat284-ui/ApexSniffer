# ApexSniffer

ApexSniffer is a professional-grade, high-performance network sniffing and logging system. It leverages a hybrid architecture, combining the raw speed and memory safety of **Rust** for packet capture with the powerful analysis and concurrency capabilities of **Go** for data processing and logging.

## 🚀 Architecture

ApexSniffer utilizes a producer-consumer model to ensure zero packet loss even under high traffic loads.

- **Rust Core (Producer):** Interfaces directly with the NIC using `libpcap`/`npcap`/`af_packet`. It handles kernel-level BPF filtering and streams raw packets through a high-throughput IPC channel.
- **Go Analyzer (Consumer):** Consumes the raw byte stream via IPC, decodes packets using `gopacket`, applies analysis logic, and persists the data to structured logs.

For a detailed technical breakdown, see [ARCHITECTURE.md](ARCHITECTURE.md).

## 🛠 Installation

### Prerequisites

- **Rust:** Latest stable version (via rustup)
- **Go:** 1.21+
- **Npcap (Windows):** Required for packet capture on Windows.
- **Libpcap (Linux):** Required for packet capture on Linux.

### Building from Source

```bash
# Clone the repository
git clone https://github.com/your-username/ApexSniffer.git
cd ApexSniffer

# Build using the provided Makefile
make all
```

## 💻 Usage

### Running the System

To start the sniffer, you typically need administrative/root privileges to access the network interface.

```bash
# Start the Rust Core (Producer)
sudo ./bin/apex-core --interface eth0

# Start the Go Analyzer (Consumer)
./bin/apex-analyzer --config config/config.yaml
```

### Configuration

Configuration is handled via the `config/` directory. You can specify:
- Target interfaces
- BPF filters
- Log output formats (JSON, Parquet)
- Analysis rules

## ⚡ Performance Features

- **Zero-Copy Pipeline:** Minimized data movement between the NIC and the analyzer.
- **Kernel-Level Filtering:** BPF filters ensure only relevant traffic reaches user space.
- **Asynchronous IPC:** High-performance named pipes/UDS prevent the analyzer from bottlenecking the capture core.
- **Parallel Decoding:** Go's concurrency model allows for simultaneous decoding of multiple packet streams.

## 📄 License

Refer to [DISCLAIMER.md](DISCLAIMER.md) for legal usage terms.
