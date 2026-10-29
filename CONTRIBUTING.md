# Contributing to ApexSniffer

First off, thank you for considering contributing to ApexSniffer! It is a high-performance system where efficiency and correctness are paramount.

## 🤝 Contribution Workflow

1. **Fork the Repository:** Create your own fork of the project.
2. **Create a Feature Branch:** Use a descriptive name (e.g., `feature/packet-decoding-optimization`).
3. **Implement Changes:** Ensure your code follows the project's style guidelines.
4. **Run Tests:** Verify that your changes don't introduce regressions.
5. **Submit a Pull Request:** Provide a clear description of the changes and why they are necessary.

## 🛠 Coding Standards

### Rust Core
- Follow standard Rust naming conventions (`snake_case` for functions/variables).
- Prefer `zero-copy` abstractions where possible to maintain performance.
- Use `cargo fmt` before committing.

### Go Analyzer
- Follow standard Go idioms and `gofmt`.
- Utilize Go channels for concurrent processing.
- Avoid heavy locks in the hot path of the analyzer.

## 🧪 Testing

- **Unit Tests:** Write unit tests for new logic in both Rust and Go.
- **Integration Tests:** Use the `scripts/` directory to run end-to-end capture tests with sample PCAP files.
- **Performance Benchmarks:** If you are optimizing a critical path, provide benchmark results showing the improvement.

## 📖 Documentation

When adding a new feature, please update the corresponding documentation in `ARCHITECTURE.md` or the `/Docs` folder.

## 💬 Communication

For major changes, please open an **Issue** first to discuss the proposed design. This ensures we are aligned on the architecture before significant effort is spent.
