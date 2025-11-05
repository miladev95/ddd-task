.PHONY: help build run test clean lint fmt docs

help:
	@echo "Task Management System - DDD Architecture"
	@echo "=========================================="
	@echo ""
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run all tests"
	@echo "  make test-unit   - Run unit tests only"
	@echo "  make test-int    - Run integration tests only"
	@echo "  make lint        - Run linter"
	@echo "  make fmt         - Format code"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make docs        - Open architecture documentation"
	@echo ""

build:
	@echo "Building application..."
	go build -o bin/task-management main.go
	@echo "Build complete. Binary: bin/task-management"

run:
	@echo "Running application..."
	go run main.go

test:
	@echo "Running all tests..."
	go test -v ./tests/unit ./tests/integration

test-unit:
	@echo "Running unit tests..."
	go test -v ./tests/unit

test-int:
	@echo "Running integration tests..."
	go test -v ./tests/integration

lint:
	@echo "Running linter..."
	go vet ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

example:
	@echo "Running example..."
	go run examples/usage_example.go

.PHONY: all
all: fmt lint test build