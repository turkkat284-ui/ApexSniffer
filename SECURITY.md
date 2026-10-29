# Security Policy

## 🛡 Security Considerations

ApexSniffer operates at a low level of the networking stack and requires elevated privileges (root/Administrator) to capture raw packets. This introduces several security considerations:

1. **Privileged Execution:** The Rust Core must run as root to access the NIC. We recommend using capabilities (e.g., `setcap` on Linux) to grant only `CAP_NET_RAW` and `CAP_NET_ADMIN` instead of full root access.
2. **Input Validation:** Raw network data is untrusted. The Go Analyzer performs rigorous decoding and validation via `gopacket` to prevent buffer overflows or injection attacks.
3. **IPC Security:** Communication between the Rust Core and Go Analyzer happens via local IPC (Named Pipes/UDS). Ensure permissions on these pipes are restricted to the system user running the process.
4. **Resource Exhaustion:** A high-volume traffic spike can potentially exhaust system memory if the analyzer cannot keep up. The system uses a bounded ring buffer to prevent OOM (Out of Memory) crashes.

## 🚩 Reporting a Vulnerability

We take the security of ApexSniffer seriously. If you discover a security vulnerability, please do **not** open a public issue. Instead, follow these steps:

1. **Email the Maintainers:** Send a detailed report to `security@apexsniffer.example` (placeholder).
2. **Include Details:**
   - Description of the vulnerability.
   - Steps to reproduce (Proof of Concept).
   - Potential impact.
   - Suggested fix.
3. **Coordinated Disclosure:** We will acknowledge receipt of your report within 48 hours and work with you to fix the issue before public disclosure.

## 🔒 Security Auditing

We encourage security researchers to audit our IPC protocol and the Go-based decoding logic, as these are the primary attack surfaces for malformed packet injection.
