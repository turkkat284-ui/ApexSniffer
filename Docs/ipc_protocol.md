# IPC Protocol Specification

ApexSniffer uses a custom binary framing protocol for low-latency communication between the Rust Core and Go Analyzer.

## Frame Structure

| Offset | Field | Size | Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| 0 | `Magic` | 2 Bytes | `u16` | Always `0x4150` ('AP') |
| 2 | `Timestamp` | 8 Bytes | `i64` | Nanoseconds since epoch |
| 10 | `Length` | 4 Bytes | `u32` | Length of the payload (N) |
| 14 | `Payload` | N Bytes | `bytes` | Raw Ethernet frame |

## Error Handling
If the Go Analyzer encounters a frame without the `0x4150` magic header, it will drop the current buffer and seek the next magic header to resynchronize the stream.
