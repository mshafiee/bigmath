#!/bin/bash
# Run all linters (go vet and golangci-lint)

set -e

echo "Running linters..."

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo "Error: golangci-lint is not installed."
    echo "Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    exit 1
fi

# Run go vet (assembly file warnings are expected and non-fatal)
echo "Running go vet..."
vet_output=$(go vet ./... 2>&1) || true
# Filter out assembly file (.s) errors and package headers
non_asm_errors=$(echo "$vet_output" | grep -v "\.s:" | grep -v "^#" | grep -v "^$" || true)
if [ -n "$non_asm_errors" ]; then
    echo "go vet found errors in Go files:"
    echo "$non_asm_errors"
    exit 1
fi

# Run golangci-lint
echo "Running golangci-lint..."
if ! golangci-lint run; then
    echo "golangci-lint failed!"
    exit 1
fi

echo "All linters passed!"

