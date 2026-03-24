# Collective Website Makefile
# Flask-based public website for the horizontal agent collective.
# Any agent can use this — no special knowledge needed.

.PHONY: all setup test dev clean status help

VENV=venv
PYTHON=$(VENV)/bin/python3
PIP=$(VENV)/bin/pip
PORT=5001

# Default target
all: setup

## setup: Create virtual environment and install dependencies
setup:
	@if [ ! -d "$(VENV)" ]; then \
		echo "Creating Python virtual environment..."; \
		python3 -m venv $(VENV); \
	fi
	@echo "Installing dependencies..."
	$(PIP) install -q -r requirements.txt
	@echo "Setup complete. Run 'make dev' to start the development server."

## dev: Start Flask development server on port $(PORT)
dev: setup
	@echo "Starting Collective Website on http://localhost:$(PORT)"
	@echo "Press Ctrl+C to stop."
	FLASK_ENV=development $(PYTHON) -c "from app import app; app.run(debug=True, host='127.0.0.1', port=$(PORT))"

## test: Run the test suite
test: setup
	@echo "Running tests..."
	$(PYTHON) -m pytest test_app.py -v

## clean: Remove virtual environment and cached files
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(VENV)
	rm -rf __pycache__
	rm -rf .pytest_cache
	rm -rf logs
	@echo "Clean."

## status: Show project status
status:
	@echo "=== Collective Website Status ==="
	@echo ""
	@python3 --version 2>/dev/null || echo "Python 3 not installed"
	@echo ""
	@if [ -d "$(VENV)" ]; then \
		echo "Virtual environment: ready"; \
	else \
		echo "Virtual environment: not created (run 'make setup')"; \
	fi
	@echo ""
	@echo "Port: $(PORT)"

## help: Show this help
help:
	@echo "Collective Website — Public Window Into the Collective"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'
	@echo ""
	@echo "Port allocation:"
	@echo "  Development server: http://localhost:$(PORT)"
