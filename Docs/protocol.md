# Binary Wire Protocol Specification

The ApexSniffer binary protocol is designed for minimal serialization overhead and maximum parsing speed.

## Frame Format

Every packet sent through the IPC pipeline follows this structure:

| Offset | Field | Length | Type | Description |
| :--- | :--- | :--- | :--- | :--- |
| 0 | Magic | 2 Bytes | uint16 | Constant `0xAPEX` to identify the start of a frame. |
| 2 | Timestamp | 8 Bytes | int64 | Nanoseconds since Unix Epoch. |
| 10 | Length | 4 Bytes | uint32 | Length of the payload in bytes. |
| 14 | Payload | Variable | bytes | The raw packet data or parsed metadata. |

## Protocol Rules

1. **Endianness**: All multi-byte integers are stored in **Big-Endian** (Network Byte Order).
2. **Alignment**: Frames are not padded; the payload begins immediately after the length field.
3. **Integrity**: The Magic bytes are used for synchronization. If the parser encounters bytes other than `0xAPEX` at the start of a frame, it scans forward until the next magic sequence is found.
