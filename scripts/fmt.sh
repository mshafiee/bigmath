#!/bin/bash
# Format all Go files and check if any files were modified

set -e

echo "Running go fmt..."

# Run go fmt and capture output
if ! output=$(go fmt ./... 2>&1); then
    echo "Error running go fmt:"
    echo "$output"
    exit 1
fi

# Check if any files were modified
if [ -n "$output" ]; then
    echo "The following files were modified by go fmt:"
    echo "$output"
    echo ""
    echo "Please run 'go fmt ./...' to format your code."
    exit 1
fi

# Alternative check: use go fmt -d to see what would change
if diff_output=$(go fmt -d ./... 2>&1); then
    if [ -n "$diff_output" ]; then
        echo "The following files need formatting:"
        echo "$diff_output"
        echo ""
        echo "Please run 'go fmt ./...' to format your code."
        exit 1
    fi
fi

echo "All files are properly formatted!"

