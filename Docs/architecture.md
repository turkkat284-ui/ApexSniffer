# Architecture

ApexSniffer is designed for high-performance packet capture and analysis with a focus on minimal overhead and zero-copy data paths.

## Binary IPC
To ensure maximum throughput between the capture engine and the analysis pipeline, ApexSniffer utilizes a Binary Inter-Process Communication (IPC) mechanism. Instead of using heavy formats like JSON or XML, data is serialized into a compact binary format that mirrors the memory layout of the captured packets.

## BPF Filtering
The system leverages Berkeley Packet Filter (BPF) programs to filter traffic at the kernel level. By pushing filters into the kernel, ApexSniffer avoids the overhead of copying packets to user-space that would eventually be discarded. 

## Async Pipeline
The analysis pipeline is implemented as an asynchronous series of stages:
1. **Ingestion Stage**: Captures raw packets from the network interface.
2. **Filtering Stage**: Applies BPF and user-defined rules.
3. **Parsing Stage**: Decodes binary data into structured formats.
4. **Export/Analysis Stage**: Sends data to the final consumer (e.g., a dashboard or log file).

This pipeline uses a non-blocking architecture to ensure that the capture engine never stalls, even during heavy analysis loads.
