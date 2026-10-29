# Installation Guide

This repository currently contains a Rust packet-capture core and a Go traffic-analysis component. To use the code as it exists today, you need both toolchains and the appropriate system privileges for packet capture.

## Prerequisites

- Rust toolchain (`cargo`, `rustc`)
- Go toolchain
- Linux: `libpcap` development package
- Windows: Npcap with WinPcap compatibility mode enabled
- Administrator/root privileges for packet capture

## Linux

1. Install Rust:

   ```bash
   curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
   ```

2. Install Go and `libpcap` development files:

   ```bash
   sudo apt-get update
   sudo apt-get install -y golang-go libpcap-dev
   ```

3. Build the Rust capture core:

   ```bash
   cd src/Rust-core
   cargo build --release
   ```

4. Build the Go analyzer module:

   ```bash
   cd src/Go-analyzer
   go build ./...
   ```

5. Run the capture side with elevated privileges:

   ```bash
   cd src/Rust-core
   cargo run --release
   ```

The current Rust code connects to `/tmp/apex_sniffer.sock` on Linux.

## Windows

1. Install Rust via `rustup`.
2. Install Go from the official Go installer.
3. Install Npcap and enable the WinPcap compatibility mode.
4. Build the Rust side:

   ```powershell
   cd src\Rust-core
   cargo build --release
   ```

5. Build the Go module:

   ```powershell
   cd src\Go-analyzer
   go build ./...
   ```

6. Run the capture code in an elevated terminal:

   ```powershell
   cd src\Rust-core
   cargo run --release
   ```

The current Windows implementation attempts to open `\\.\pipe\apex_sniffer`.

## Important note

The repository is still a prototype. The source code does not currently expose a complete CLI with `--interface`, `--config`, or a standalone analyzer binary. The safe and accurate way to build today is to use the per-project commands above.
