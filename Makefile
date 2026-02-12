# Makefile for quantds

.PHONY: all build test lint fmt gen-docs help clean

# Default target
all: test build

# Build the project
build:
	go build ./...

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Generate documentation (Supported Data Sources table)
gen-docs:
	go run cmd/gendoc/main.go

# Clean build artifacts
clean:
	go clean
	rm -rf dist/

# Show help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all       Test and build"
	@echo "  build     Build the project"
	@echo "  test      Run tests"
	@echo "  fmt       Format code"
	@echo "  lint      Run linter"
	@echo "  gen-docs  Generate supported data sources table in README.md"
	@echo "  clean     Clean build artifacts"
	@echo "  help      Show this help message"
