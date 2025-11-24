# Contributing to bigmath

Thank you for your interest in contributing to bigmath! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue with:
- A clear description of the bug
- Steps to reproduce
- Expected vs actual behavior
- Go version and platform information
- Minimal code example if possible

### Suggesting Features

Feature suggestions are welcome! Please open an issue describing:
- The use case for the feature
- How it would benefit users
- Any implementation ideas

### Submitting Pull Requests

1. **Fork the repository** and create a branch from `main`
2. **Make your changes** following the coding standards below
3. **Add tests** for new functionality
4. **Update documentation** if needed
5. **Ensure tests pass**: `go test ./...`
6. **Ensure code compiles**: `go build ./...`
7. **Submit a pull request** with a clear description

## Coding Standards

### Go Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` to format code
- Run `golint` and fix warnings
- Keep functions focused and small
- Add comments for exported functions and types

### Assembly Code

- Follow Go assembly conventions
- Add comments explaining optimizations
- Test on both AMD64 and ARM64 when applicable
- Include fallback implementations

### Testing

- Write tests for all new functionality
- Aim for high test coverage
- Include edge cases and error conditions
- Use table-driven tests where appropriate

### Documentation

- Update README.md for user-facing changes
- Update DOCS.md for API changes
- Add examples for new features
- Keep comments clear and concise

## Development Setup

```bash
# Clone your fork
git clone https://github.com/your-username/bigmath.git
cd bigmath

# Run tests
go test ./...

# Run benchmarks
go test -bench=. ./...

# Check code quality
go vet ./...
golint ./...
```

## Commit Messages

- Use clear, descriptive commit messages
- Start with a verb in imperative mood (e.g., "Add", "Fix", "Update")
- Reference issue numbers when applicable
- Keep the first line under 72 characters

Example:
```
Add BigLog10 function for base-10 logarithm

Implements log10(x) = ln(x) / ln(10) with proper
precision handling and rounding modes.

Fixes #123
```

## Questions?

Feel free to open an issue for any questions about contributing!

