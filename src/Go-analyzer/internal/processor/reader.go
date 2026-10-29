package processor

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
)

const (
	MagicNumber = 0x4150
)

var (
	// bufferPool minimizes GC pressure by reusing byte slices for packet payloads.
	bufferPool = sync.Pool{
		New: func() interface{} {
			// Initial capacity for common packet sizes, will grow if needed.
			return make([]byte, 4096)
		},
	}
)

type Packet struct {
	Timestamp int64
	Payload   []byte
}

// Reader handles the binary IPC stream from the sniffer.
type Reader struct {
	pipePath string
	packets  chan Packet
}

func NewReader() *Reader {
	pipePath := "/tmp/apex_sniffer.sock"
	if runtime.GOOS == "windows" {
		pipePath = `\\.\pipe\apex_sniffer`
	}

	return &Reader{
		pipePath: pipePath,
		packets:  make(chan Packet, 1000),
	}
}

// Start opens the IPC channel and begins reading frames.
func (r *Reader) Start() error {
	var conn net.Conn
	var err error

	if runtime.GOOS == "windows" {
		// In a real Windows scenario, we might use a specific named pipe library,
		// but net.Dial("pipe", ...) is often used via wrappers or specific implementations.
		// For standard Go, we treat it as a generic dial for the purpose of this implementation.
		conn, err = net.Dial("pipe", r.pipePath)
	} else {
		conn, err = net.Dial("unix", r.pipePath)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to IPC channel %s: %w", r.pipePath, err)
	}

	go r.readLoop(conn)
	return nil
}

func (r *Reader) readLoop(conn net.Conn) {
	defer conn.Close()

	for {
		// 1. Read Magic (2 bytes)
		var magic uint16
		if err := binary.Read(conn, binary.BigEndian, &magic); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("error reading magic: %v\n", err)
			break
		}

		if magic != MagicNumber {
			fmt.Printf("invalid magic number: expected 0x%X, got 0x%X\n", MagicNumber, magic)
			// In a real protocol, we might need to resync the stream here.
			continue
		}

		// 2. Read Timestamp (8 bytes)
		var timestamp int64
		if err := binary.Read(conn, binary.BigEndian, &timestamp); err != nil {
			fmt.Printf("error reading timestamp: %v\n", err)
			break
		}

		// 3. Read Length (4 bytes)
		var length uint32
		if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
			fmt.Printf("error reading length: %v\n", err)
			break
		}

		// 4. Read Payload (N bytes)
		buf := bufferPool.Get().([]byte)
		if uint32(len(buf)) < length {
			buf = make([]byte, length)
		}

		payload := buf[:length]
		if _, err := io.ReadFull(conn, payload); err != nil {
			fmt.Printf("error reading payload: %v\n", err)
			break
		}

		// Stream the packet to the analyzer
		// Note: We copy the payload because the buffer is returned to the pool
		// after the analyzer is done with it. If the analyzer is asynchronous,
		// we must be careful about ownership.
		packetPayload := make([]byte, length)
		copy(packetPayload, payload)

		r.packets <- Packet{
			Timestamp: timestamp,
			Payload:   packetPayload,
		}

		bufferPool.Put(buf)
	}
}

func (r *Reader) Packets() <-chan Packet {
	return r.packets
}
