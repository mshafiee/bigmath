.PHONY: fmt lint lint-fix check

# Format all Go files
fmt:
	@echo "Running go fmt..."
	@go fmt ./...

# Run all linters
lint: check-golangci-lint
	@echo "Running go vet..."
	@go vet ./... 2>&1 | grep -v "\.s:" || true
	@echo "Running golangci-lint..."
	@golangci-lint run

# Run linters with auto-fix
lint-fix: check-golangci-lint
	@echo "Running golangci-lint with --fix..."
	@golangci-lint run --fix

# Run both fmt and lint
check: fmt lint
	@echo "All checks passed!"

# Check if golangci-lint is installed
check-golangci-lint:
	@which golangci-lint > /dev/null || (echo "Error: golangci-lint is not installed. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)

