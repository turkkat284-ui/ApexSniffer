#!/bin/bash

# ApexSniffer Setup Script
# This script installs dependencies for the ApexSniffer project.

# --- Colors for output ---
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# --- Header ---
echo -e "${BLUE}======================================================================${NC}"
echo -e "${BLUE}                ApexSniffer Project Setup Script                      ${NC}"
echo -e "${BLUE}======================================================================${NC}"

# --- Usage ---
usage() {
    echo "Usage: $0"
    echo "This script will detect your OS and install the following dependencies:"
    echo "  - Rust Toolchain (rustup, cargo)"
    echo "  - Go Toolchain (1.21+)"
    echo "  - Packet Capture Libraries (libpcap)"
    echo "  - Build Tools (make, gcc/clang)"
    echo ""
    echo "Note: Windows users should use WSL or refer to the manual installation guide."
    exit 1
}

# Check for help flag
if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    usage
fi

# --- OS Detection ---
OS_TYPE="$(uname -s)"
case "${OS_TYPE}" in
    Linux*)
        echo -e "${GREEN}[*] Detected OS: Linux${NC}"
        PKG_MANAGER="apt-get" # Defaulting to Debian/Ubuntu.
        # Check if apt-get exists
        if ! command -v apt-get &> /dev/null; then
            echo -e "${RED}[!] This script currently only supports Debian/Ubuntu-based Linux distributions.${NC}"
            exit 1
        fi
        ;;
    Darwin*)
        echo -e "${GREEN}[*] Detected OS: MacOS${NC}"
        PKG_MANAGER="brew"
        if ! command -v brew &> /dev/null; then
            echo -e "${YELLOW}[!] Homebrew not found. Installing Homebrew...${NC}"
            /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        fi
        ;;
    *)
        echo -e "${RED}[!] Unsupported OS: ${OS_TYPE}${NC}"
        echo -e "${YELLOW}Windows users: Please use WSL (Windows Subsystem for Linux) to run this script.${NC}"
        exit 1
        ;;
esac

# --- Dependency Installation Functions ---

install_rust() {
    if command -v cargo &> /dev/null; then
        echo -e "${GREEN}[+] Rust is already installed.${NC}"
    else
        echo -e "${YELLOW}[*] Installing Rust toolchain...${NC}"
        curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
        source "$HOME/.cargo/env"
    fi
}

install_go() {
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        # Compare versions (requires 1.21+)
        if [[ "$(printf '%s\n' "1.21" "$GO_VERSION" | sort -V | head -n1)" == "1.21" ]]; then
            echo -e "${GREEN}[+] Go is already installed (Version: $GO_VERSION).${NC}"
        else
            echo -e "${YELLOW}[*] Go version $GO_VERSION is too old. Updating to latest stable...${NC}"
            install_go_latest
        fi
    else
        echo -e "${YELLOW}[*] Installing Go toolchain...${NC}"
        install_go_latest
    fi
}

install_go_latest() {
    if [[ "$PKG_MANAGER" == "apt-get" ]]; then
        sudo apt-get update
        sudo apt-get install -y golang
    elif [[ "$PKG_MANAGER" == "brew" ]]; then
        brew install go
    fi
}

install_build_tools() {
    echo -e "${YELLOW}[*] Installing build tools and pcap libraries...${NC}"
    if [[ "$PKG_MANAGER" == "apt-get" ]]; then
        sudo apt-get update
        sudo apt-get install -y build-essential make libpcap-dev gcc
    elif [[ "$PKG_MANAGER" == "brew" ]]; then
        brew install make libpcap gcc
    fi
}

# --- Execution ---

# 1. Install Build Tools & Pcap
install_build_tools

# 2. Install Rust
install_rust

# 3. Install Go
install_go

# --- Verification ---
echo -e "\n${BLUE}--- Verifying Installations ---${NC}"

VERIFIED=true

# Verify Rust
if command -v cargo &> /dev/null; then
    echo -e "${GREEN}✓ Rust/Cargo: Installed${NC}"
else
    echo -e "${RED}✗ Rust/Cargo: Missing${NC}"
    VERIFIED=false
fi

# Verify Go
if command -v go &> /dev/null; then
    echo -e "${GREEN}✓ Go: Installed${NC}"
else
    echo -e "${RED}✗ Go: Missing${NC}"
    VERIFIED=false
fi

# Verify Make
if command -v make &> /dev/null; then
    echo -e "${GREEN}✓ Make: Installed${NC}"
else
    echo -e "${RED}✗ Make: Missing${NC}"
    VERIFIED=false
fi

# Verify C Compiler
if command -v gcc &> /dev/null || command -v clang &> /dev/null; then
    echo -e "${GREEN}✓ C Compiler: Installed${NC}"
else
    echo -e "${RED}✗ C Compiler: Missing${NC}"
    VERIFIED=false
fi

if [ "$VERIFIED" = true ]; then
    echo -e "\n${GREEN}All dependencies installed successfully!${NC}"

    # Run make all to build the project
    echo -e "${BLUE}[*] Running 'make all' to build the project...${NC}"
    if make all; then
        echo -e "${GREEN}Build successful! ApexSniffer is ready to use.${NC}"
    else
        echo -e "${RED}[!] Build failed. Please check the logs above for errors.${NC}"
        exit 1
    fi
else
    echo -e "\n${RED}[!] Some dependencies failed to install. Please check the errors above.${NC}"
    exit 1
fi

echo -e "${BLUE}======================================================================${NC}"
echo -e "${GREEN} Setup Complete! ${NC}"
echo -e "${BLUE}======================================================================${NC}"
