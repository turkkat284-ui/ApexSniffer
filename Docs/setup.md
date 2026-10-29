# Setup Guide

This guide provides step-by-step instructions for installing ApexSniffer on Windows and Linux.

## Windows Installation

ApexSniffer requires **Npcap** for packet capture on Windows.

1. **Install Npcap**:
   - Download the Npcap installer from [npcap.com](https://npcap.com/).
   - During installation, ensure you check the option **"Install Npcap in WinPcap API-compatible Mode"**.
2. **Build ApexSniffer**:
   - Ensure you have the Go compiler installed.
   - Run `make build` (if using a Makefile) or `go build -o apexsniffer .`.
3. **Run with Privileges**:
   - Since packet capture requires administrative access, run the terminal as **Administrator**.
   - Execute: `.\apexsniffer.exe`

## Linux Installation

ApexSniffer requires **libpcap** for packet capture on Linux.

1. **Install libpcap**:
   - On Ubuntu/Debian: `sudo apt-get install libpcap-dev`
   - On CentOS/RHEL: `sudo yum install libpcap-devel`
2. **Build ApexSniffer**:
   - Ensure you have the Go compiler installed.
   - Run `make build` or `go build -o apexsniffer .`.
3. **Run with Privileges**:
   - Packet capture requires `CAP_NET_RAW` capabilities or root access.
   - Execute: `sudo ./apexsniffer`
   - *Alternative (non-root)*: `sudo setcap cap_net_raw,cap_net_admin=eip ./apexsniffer`
