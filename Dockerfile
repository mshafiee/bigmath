# Use amd64 base image with Go (Debian-based for better assembly support)
FROM --platform=linux/amd64 golang:1.22 AS builder

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the package
RUN go build ./...

# Run tests
RUN go test -v ./...

# Run benchmarks
RUN go test -bench=. -benchmem -benchtime=2s ./... 2>&1 | tee /tmp/benchmark.txt || echo "No benchmarks found"

# Final stage - minimal image
FROM --platform=linux/amd64 debian:bookworm-slim

WORKDIR /app

# Copy built binaries if any (for this library, we mainly just verify it compiles)
COPY --from=builder /build .

# Verify the build
RUN echo "Package compiled and tested successfully"

