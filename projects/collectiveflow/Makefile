# CollectiveFlow Makefile
# Build system for the Go CLI and Flask web interface
# Any agent in the collective can use this — no special knowledge needed.

.PHONY: all build test clean dev setup help web-setup web-dev web-test

# Build configuration
BINARY_NAME=collectiveflow
GO_PKG=./internal/...
WEB_DIR=web
VENV=$(WEB_DIR)/venv

# Default target
all: build

## setup: Install all dependencies (Go + Python)
setup: deps web-setup

## build: Build the CollectiveFlow Go CLI binary
build: deps
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) ./...
	@echo "Built: ./$(BINARY_NAME)"

## test: Run all tests (Go + Flask web)
test: test-go web-test

## test-go: Run Go unit tests
test-go:
	@echo "Running Go tests..."
	go test -v -race $(GO_PKG)

## dev: Start the Flask web interface for local development (port 5000)
dev: web-setup
	@echo "Starting CollectiveFlow web on http://localhost:5000"
	cd $(WEB_DIR) && . venv/bin/activate && python3 app.py

## clean: Remove build artifacts and virtual environments
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -rf $(VENV)
	rm -rf $(WEB_DIR)/__pycache__
	@echo "Clean."

## deps: Download Go module dependencies
deps:
	@echo "Downloading Go dependencies..."
	go mod download

## fmt: Format Go source code
fmt:
	go fmt ./...

## vet: Run Go vet checks
vet:
	go vet ./...

## check: Run fmt + vet + test
check: fmt vet test

## web-setup: Create Python venv and install Flask dependencies
web-setup:
	@if [ ! -d "$(VENV)" ]; then \
		echo "Creating Python virtual environment..."; \
		python3 -m venv $(VENV); \
	fi
	@echo "Installing Python dependencies..."
	$(VENV)/bin/pip install -q -r $(WEB_DIR)/requirements.txt

## web-dev: Start only the Flask web interface
web-dev: web-setup
	cd $(WEB_DIR) && . venv/bin/activate && python3 app.py

## web-test: Run Flask web tests (if test files exist)
web-test: web-setup
	@if [ -f "$(WEB_DIR)/conftest.py" ] || ls $(WEB_DIR)/tests/*.py >/dev/null 2>&1; then \
		echo "Running Flask web tests..."; \
		cd $(WEB_DIR) && . venv/bin/activate && python3 -m pytest -v; \
	else \
		echo "No web tests found — skipping."; \
	fi

## run: Build and run the CLI with arguments (usage: make run ARGS="status active")
run: build
	./$(BINARY_NAME) $(ARGS)

## status: Show project status
status:
	@echo "=== CollectiveFlow Status ==="
	@echo ""
	@echo "Go module:"
	@go version 2>/dev/null || echo "  Go not installed"
	@echo ""
	@echo "Go dependencies:"
	@go mod verify 2>/dev/null && echo "  Verified." || echo "  Not verified (run 'make deps')."
	@echo ""
	@echo "Web interface:"
	@if [ -d "$(VENV)" ]; then \
		echo "  Python venv: ready"; \
	else \
		echo "  Python venv: not created (run 'make web-setup')"; \
	fi
	@echo ""
	@echo "Binary:"
	@if [ -f "$(BINARY_NAME)" ]; then \
		echo "  $(BINARY_NAME): built"; \
	else \
		echo "  $(BINARY_NAME): not built (run 'make build')"; \
	fi

## help: Show this help
help:
	@echo "CollectiveFlow — Horizontal Decision-Making Tool"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'
	@echo ""
	@echo "Port allocation:"
	@echo "  Flask web interface: http://localhost:5000"
