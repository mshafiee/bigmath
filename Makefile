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

# Test build on multiple architectures
test-build:
	@echo "Testing build on multiple architectures..."
	@echo "Testing amd64..."
	@GOOS=linux GOARCH=amd64 go build ./... || (echo "Build failed on amd64" && exit 1)
	@echo "Testing arm64..."
	@GOOS=linux GOARCH=arm64 go build ./... || (echo "Build failed on arm64" && exit 1)
	@echo "Testing s390x (generic platform)..."
	@GOOS=linux GOARCH=s390x go build ./... || (echo "Build failed on s390x - check for missing fallback implementations!" && exit 1)
	@echo "Testing ppc64le (generic platform)..."
	@GOOS=linux GOARCH=ppc64le go build ./... || (echo "Build failed on ppc64le - check for missing fallback implementations!" && exit 1)
	@echo "All architecture builds passed!"

# Run both fmt and lint
check: fmt lint test-build
	@echo "All checks passed!"

# Check if golangci-lint is installed
check-golangci-lint:
	@which golangci-lint > /dev/null || (echo "Error: golangci-lint is not installed. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)

