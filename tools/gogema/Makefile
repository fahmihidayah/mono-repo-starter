.PHONY: build install run test clean help dev generate example fmt lint tidy release

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=gogema
BINARY_PATH=../../dist/apps/cli/$(BINARY_NAME)
INSTALL_PATH=$(HOME)/go/bin/$(BINARY_NAME)
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p ../../dist/apps/cli
	go build $(LDFLAGS) -o $(BINARY_PATH) cmd/main.go
	@echo "✓ Binary created at: $(BINARY_PATH)"
	@echo "  Size: $$(du -h $(BINARY_PATH) | cut -f1)"

# Build with optimizations for release
build-release:
	@echo "Building $(BINARY_NAME) for release..."
	@mkdir -p ../../dist/apps/cli
	go build $(LDFLAGS) -ldflags="-s -w" -o $(BINARY_PATH) cmd/main.go
	@echo "✓ Release binary created at: $(BINARY_PATH)"
	@echo "  Size: $$(du -h $(BINARY_PATH) | cut -f1)"

# Install the binary globally
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@cp $(BINARY_PATH) $(INSTALL_PATH)
	@chmod +x $(INSTALL_PATH)
	@echo "✓ $(BINARY_NAME) installed successfully"
	@echo "  Run 'gogema --help' to get started"

# Uninstall the binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_PATH)
	@echo "✓ $(BINARY_NAME) uninstalled"

# Run the application directly (without building)
run:
	@echo "Running $(BINARY_NAME)..."
	go run cmd/main.go

# Run with example project
example:
	@echo "Running $(BINARY_NAME) with example project..."
	go run cmd/main.go --path example --verbose

# Run with example project in dry-run mode
example-dry:
	@echo "Dry run with example project..."
	go run cmd/main.go --path example --verbose --dry-run

# Generate code for a specific path
generate:
	@if [ -z "$(PATH_ARG)" ]; then \
		echo "Usage: make generate PATH_ARG=/path/to/project"; \
		exit 1; \
	fi
	@echo "Generating code from $(PATH_ARG)..."
	go run cmd/main.go --path $(PATH_ARG) --verbose

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✓ Dependencies installed"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Code formatted"

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run: make install-tools"; \
		exit 1; \
	fi

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy
	@echo "✓ Dependencies tidied"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf ../../dist/apps/cli
	@rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "✓ Development tools installed"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p ../../dist/apps/cli
	@echo "  Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ../../dist/apps/cli/$(BINARY_NAME)-linux-amd64 cmd/main.go
	@echo "  Building for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o ../../dist/apps/cli/$(BINARY_NAME)-linux-arm64 cmd/main.go
	@echo "  Building for macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o ../../dist/apps/cli/$(BINARY_NAME)-darwin-amd64 cmd/main.go
	@echo "  Building for macOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o ../../dist/apps/cli/$(BINARY_NAME)-darwin-arm64 cmd/main.go
	@echo "  Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ../../dist/apps/cli/$(BINARY_NAME)-windows-amd64.exe cmd/main.go
	@echo "✓ All platforms built successfully"
	@ls -lh ../../dist/apps/cli/

# Create a release package
release: clean build-release
	@echo "Creating release package..."
	@mkdir -p ../../dist/apps/cli/release
	@cp $(BINARY_PATH) ../../dist/apps/cli/release/
	@cp README.md ../../dist/apps/cli/release/
	@cp -r templates ../../dist/apps/cli/release/
	@cp -r example ../../dist/apps/cli/release/
	@cd ../../dist/apps/cli && tar -czf gogema-$(VERSION).tar.gz release/
	@rm -rf ../../dist/apps/cli/release
	@echo "✓ Release package created: dist/apps/cli/gogema-$(VERSION).tar.gz"

# Development mode - watch for changes and rebuild
dev:
	@echo "Starting development mode..."
	@echo "Watching for file changes..."
	@while true; do \
		inotifywait -e modify -r cmd/ internal/ 2>/dev/null || \
		fswatch -1 -r cmd/ internal/ 2>/dev/null || \
		(echo "Install 'fswatch' or 'inotifywait' for auto-reload support" && sleep 5); \
		clear; \
		echo "Rebuilding..."; \
		make build; \
		echo "\n✓ Rebuild complete. Waiting for changes..."; \
	done

# Verify the built binary
verify: build
	@echo "Verifying binary..."
	@$(BINARY_PATH) --version
	@$(BINARY_PATH) --help
	@echo "✓ Binary verification complete"

# Quick test - build and run with example
quick-test: build example-dry

# Show version information
version:
	@echo "GoGema CLI"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"

# Show project statistics
stats:
	@echo "Project Statistics:"
	@echo "==================="
	@echo "Go files: $$(find . -name '*.go' | wc -l)"
	@echo "Total lines of code: $$(find . -name '*.go' -exec cat {} \; | wc -l)"
	@echo "Templates: $$(find templates -name '*.tmpl' | wc -l)"
	@echo "Dependencies: $$(go list -m all | wc -l)"
	@echo ""
	@echo "Directory sizes:"
	@du -sh cmd internal templates

# Help command
help:
	@echo "GoGema CLI - Makefile Commands"
	@echo "=============================="
	@echo ""
	@echo "Build Commands:"
	@echo "  make build           - Build the application"
	@echo "  make build-release   - Build optimized release binary"
	@echo "  make build-all       - Build for all platforms (Linux, macOS, Windows)"
	@echo "  make install         - Build and install binary to ~/go/bin"
	@echo "  make uninstall       - Remove installed binary"
	@echo ""
	@echo "Development Commands:"
	@echo "  make run             - Run without building"
	@echo "  make example         - Run with example project"
	@echo "  make example-dry     - Dry run with example project"
	@echo "  make generate PATH_ARG=/path - Generate code from specific path"
	@echo "  make dev             - Watch mode (requires fswatch/inotifywait)"
	@echo ""
	@echo "Testing Commands:"
	@echo "  make test            - Run tests"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo "  make verify          - Build and verify binary"
	@echo "  make quick-test      - Build and run example (dry run)"
	@echo ""
	@echo "Code Quality Commands:"
	@echo "  make fmt             - Format Go code"
	@echo "  make lint            - Lint code with golangci-lint"
	@echo "  make tidy            - Tidy dependencies"
	@echo ""
	@echo "Maintenance Commands:"
	@echo "  make deps            - Install/update dependencies"
	@echo "  make install-tools   - Install development tools"
	@echo "  make clean           - Remove build artifacts"
	@echo ""
	@echo "Release Commands:"
	@echo "  make release         - Create release package"
	@echo "  make version         - Show version information"
	@echo ""
	@echo "Information Commands:"
	@echo "  make stats           - Show project statistics"
	@echo "  make help            - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build && ./$(BINARY_PATH) --help"
	@echo "  make install && gogema --path /path/to/project"
	@echo "  make generate PATH_ARG=../api/project-plan"
