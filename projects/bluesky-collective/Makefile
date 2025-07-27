# Bluesky Collective Makefile
# Build system for consensus-based Bluesky client

.PHONY: all build test clean install lint fmt vet deps help

# Build configuration
BINARY_NAME=bluesky-collective
BUILD_DIR=build
PKG=github.com/consensuscode/bluesky-collective
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_TIME?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"
BUILD_FLAGS=-trimpath $(LDFLAGS)

# Default target
all: build

# Build the binary
build: deps
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/bluesky-collective

# Run tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=$(BUILD_DIR)/coverage.out ./...

# Test with coverage report
coverage: test
	@echo "Generating coverage report..."
	go tool cover -html=$(BUILD_DIR)/coverage.out -o $(BUILD_DIR)/coverage.html
	@echo "Coverage report: $(BUILD_DIR)/coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# Install binary to system
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

# Lint code
lint:
	@echo "Running linters..."
	golangci-lint run

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Development build with debugging symbols
dev: deps
	@echo "Building development version..."
	@mkdir -p $(BUILD_DIR)
	go build -race -o $(BUILD_DIR)/$(BINARY_NAME)-dev ./cmd/bluesky-collective

# Quick check (fmt + vet + test)
check: fmt vet test

# Full quality check
full-check: check lint coverage

# Initialize development environment
init:
	@echo "Initializing development environment..."
	go mod tidy
	@if [ ! -f .bluesky-collective.yaml ]; then \
		echo "Creating default config..."; \
		$(BUILD_DIR)/$(BINARY_NAME) config init || true; \
	fi

# Run the tool with example usage
demo: build
	@echo "Running demo..."
	$(BUILD_DIR)/$(BINARY_NAME) config init
	@echo ""
	@echo "Example commands:"
	@echo "  $(BUILD_DIR)/$(BINARY_NAME) propose --text 'Hello from the collective!' --reasoning 'Introduction post'"
	@echo "  $(BUILD_DIR)/$(BINARY_NAME) status"
	@echo "  $(BUILD_DIR)/$(BINARY_NAME) vote --proposal proposal-123 --position support --reasoning 'Good first post'"

# Show help
help:
	@echo "Available targets:"
	@echo "  build      Build the binary"
	@echo "  test       Run tests"
	@echo "  coverage   Run tests with coverage report"
	@echo "  clean      Clean build artifacts"
	@echo "  install    Install binary to system"
	@echo "  lint       Run linters"
	@echo "  fmt        Format code"
	@echo "  vet        Vet code"
	@echo "  deps       Download dependencies"
	@echo "  tidy       Tidy dependencies"
	@echo "  dev        Build development version"
	@echo "  check      Quick quality check (fmt + vet + test)"
	@echo "  full-check Full quality check (includes lint + coverage)"
	@echo "  init       Initialize development environment"
	@echo "  demo       Build and show example usage"
	@echo "  help       Show this help"