# Performance Optimizations

To prevent packet loss during high-traffic bursts, ApexSniffer employs several memory and concurrency optimizations.

## sync.Pool for Buffer Reuse

Allocating memory for every captured packet would trigger frequent Garbage Collection (GC) cycles, leading to "stop-the-world" pauses and dropped packets. 

ApexSniffer uses `sync.Pool` to maintain a pool of pre-allocated byte slices. When a packet is captured:
1. A buffer is retrieved from the pool.
2. The packet data is copied into this buffer.
3. Once the analysis pipeline is finished with the packet, the buffer is returned to the pool.

This reduces heap allocations to nearly zero during steady-state operation.

## MPSC Optimization

The pipeline utilizes a **Multi-Producer Single-Consumer (MPSC)** pattern for packet processing.

- **Multi-Producer**: Multiple capture threads (one per NIC or CPU core) can push packets into the ingestion queue simultaneously.
- **Single-Consumer**: A dedicated worker process handles the sequencing and analysis of packets.

By using a lock-free ring buffer for the MPSC queue, ApexSniffer avoids mutex contention, allowing the capture engine to operate at wire speed without waiting for the analysis logic to complete.
