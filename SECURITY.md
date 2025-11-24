# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Which versions are eligible for receiving such patches depends on the CVSS v3.0 Rating:

| Version | Supported          |
| ------- | ------------------ |
| Latest  | :white_check_mark: |
| < Latest| :x:                |

## Reporting a Vulnerability

Please report (suspected) security vulnerabilities to **[security@example.com](mailto:security@example.com)**. You will receive a response within 48 hours. If the issue is confirmed, we will release a patch as soon as possible depending on complexity but historically within a few days.

## Security Considerations

This library performs arbitrary-precision arithmetic operations. While we strive for correctness:

- **Precision**: Results are computed with the specified precision, but intermediate calculations may use higher precision
- **Rounding**: Rounding modes follow IEEE 754-2008 standards
- **Overflow**: Large numbers may cause performance issues but should not cause crashes
- **Assembly Code**: Assembly implementations are tested for correctness but may have platform-specific behavior

If you discover a security vulnerability, please do **NOT** open a public issue. Instead, email us directly.

