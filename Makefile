# Variables
RUST_SRC = src/Rust-core
GO_SRC = src/Go-analyzer
BIN_DIR = bin

# Build targets
all: build-core build-analyzer

build-core:
	@echo "Building Rust Core..."
	cd $(RUST_SRC) && cargo build --release
	mkdir -p $(BIN_DIR)
	cp $(RUST_SRC)/target/release/rust-core $(BIN_DIR)/apex-core

build-analyzer:
	@echo "Building Go Analyzer..."
	cd $(GO_SRC) && go build -o $(BIN_DIR)/apex-analyzer ./cmd/analyzer

run-core:
	@echo "Running Rust Core..."
	./$(BIN_DIR)/apex-core

run-analyzer:
	@echo "Running Go Analyzer..."
	./$(BIN_DIR)/apex-analyzer

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
	cd $(RUST_SRC) && cargo clean

test:
	cd $(RUST_SRC) && cargo test
	cd $(GO_SRC) && go test ./...

.PHONY: all build-core build-analyzer run-core run-analyzer clean test
