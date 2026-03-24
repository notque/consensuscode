# Consensus Code — Root Makefile
# Orchestrates all sub-projects in the horizontal agent collective.
# Any agent can run any target — no special knowledge required.
#
# Projects:
#   collectiveflow      Go CLI + Flask web (port 5000)
#   bluesky-collective  Go CLI for Bluesky integration
#   collective-website  Flask public website (port 5001)
#   user-advocacy       Documentation only — no build needed

.PHONY: all setup test dev clean status help \
        collectiveflow bluesky collective-website user-advocacy

# Project directories
CF_DIR=projects/collectiveflow
BS_DIR=projects/bluesky-collective
CW_DIR=projects/collective-website
UA_DIR=projects/user-advocacy

# Default target
all: setup

# ─────────────────────────────────────────────
# Top-Level Targets
# ─────────────────────────────────────────────

## setup: Install all project dependencies
setup:
	@echo "=== Setting up all projects ==="
	@echo ""
	@echo "--- CollectiveFlow (Go + Flask) ---"
	$(MAKE) -C $(CF_DIR) setup
	@echo ""
	@echo "--- Bluesky Collective (Go) ---"
	$(MAKE) -C $(BS_DIR) setup
	@echo ""
	@echo "--- Collective Website (Flask) ---"
	$(MAKE) -C $(CW_DIR) setup
	@echo ""
	@echo "=== All projects set up ==="

## test: Run all project tests
test:
	@echo "=== Running all tests ==="
	@echo ""
	@echo "--- CollectiveFlow ---"
	$(MAKE) -C $(CF_DIR) test || true
	@echo ""
	@echo "--- Bluesky Collective ---"
	$(MAKE) -C $(BS_DIR) test || true
	@echo ""
	@echo "--- Collective Website ---"
	$(MAKE) -C $(CW_DIR) test || true
	@echo ""
	@echo "=== All tests complete ==="

## dev: Show how to start each development server
dev:
	@echo "=== Development Servers ==="
	@echo ""
	@echo "Start each project's dev server in a separate terminal:"
	@echo ""
	@echo "  CollectiveFlow web (port 5000):"
	@echo "    make collectiveflow-dev"
	@echo ""
	@echo "  Collective Website (port 5001):"
	@echo "    make collective-website-dev"
	@echo ""
	@echo "Or start them individually:"
	@echo "    cd $(CF_DIR) && make dev"
	@echo "    cd $(CW_DIR) && make dev"

## status: Show status of all projects
status:
	@echo "============================================"
	@echo "  Consensus Code — Collective Project Status"
	@echo "============================================"
	@echo ""
	@echo "--- CollectiveFlow ---"
	@$(MAKE) -s -C $(CF_DIR) status 2>/dev/null || echo "  (run 'make setup' first)"
	@echo ""
	@echo "--- Bluesky Collective ---"
	@$(MAKE) -s -C $(BS_DIR) status 2>/dev/null || echo "  (run 'make setup' first)"
	@echo ""
	@echo "--- Collective Website ---"
	@$(MAKE) -s -C $(CW_DIR) status 2>/dev/null || echo "  (run 'make setup' first)"
	@echo ""
	@echo "--- User Advocacy ---"
	@$(MAKE) -s -C $(UA_DIR) status 2>/dev/null || echo "  (run 'make setup' first)"
	@echo ""
	@echo "============================================"
	@echo "  Port Allocation"
	@echo "============================================"
	@echo "  CollectiveFlow web:   http://localhost:5000"
	@echo "  Collective Website:   http://localhost:5001"
	@echo "============================================"

## clean: Clean all build artifacts
clean:
	@echo "=== Cleaning all projects ==="
	$(MAKE) -C $(CF_DIR) clean
	$(MAKE) -C $(BS_DIR) clean
	$(MAKE) -C $(CW_DIR) clean
	@echo "=== All projects cleaned ==="

# ─────────────────────────────────────────────
# Per-Project Shortcuts
# ─────────────────────────────────────────────

## collectiveflow: Build the CollectiveFlow CLI
collectiveflow:
	$(MAKE) -C $(CF_DIR) build

## collectiveflow-dev: Start CollectiveFlow web dev server (port 5000)
collectiveflow-dev:
	$(MAKE) -C $(CF_DIR) dev

## bluesky: Build the Bluesky Collective CLI
bluesky:
	$(MAKE) -C $(BS_DIR) build

## collective-website-dev: Start Collective Website dev server (port 5001)
collective-website-dev:
	$(MAKE) -C $(CW_DIR) dev

## help: Show this help
help:
	@echo "Consensus Code — Horizontal Agent Collective"
	@echo ""
	@echo "Top-level targets:"
	@echo "  make setup    Install all project dependencies"
	@echo "  make test     Run all project tests"
	@echo "  make dev      Show how to start dev servers"
	@echo "  make status   Show project status"
	@echo "  make clean    Clean all build artifacts"
	@echo ""
	@echo "Per-project shortcuts:"
	@echo "  make collectiveflow          Build CollectiveFlow CLI"
	@echo "  make collectiveflow-dev      Start CollectiveFlow web (port 5000)"
	@echo "  make bluesky                 Build Bluesky Collective CLI"
	@echo "  make collective-website-dev  Start Collective Website (port 5001)"
	@echo ""
	@echo "Or work directly in a project:"
	@echo "  cd projects/collectiveflow && make help"
	@echo "  cd projects/bluesky-collective && make help"
	@echo "  cd projects/collective-website && make help"
	@echo "  cd projects/user-advocacy && make help"
