# User Advocacy Makefile
# This project is documentation and templates — no build step required.
# Makefile provided for consistency with other collective projects.

.PHONY: status help

## status: List all docs, guides, templates, and tools in this project
status:
	@echo "=== User Advocacy Framework Status ==="
	@echo ""
	@echo "Documentation:"
	@ls -1 docs/*.md 2>/dev/null | sed 's/^/  /' || echo "  (none)"
	@echo ""
	@echo "Guides:"
	@ls -1 guides/*.md 2>/dev/null | sed 's/^/  /' || echo "  (none)"
	@echo ""
	@echo "Templates:"
	@ls -1 templates/*.md 2>/dev/null | sed 's/^/  /' || echo "  (none)"
	@echo ""
	@echo "Tools:"
	@ls -1 tools/*.md 2>/dev/null | sed 's/^/  /' || echo "  (none)"

## help: Show this help
help:
	@echo "User Advocacy — Horizontal User Needs Framework"
	@echo ""
	@echo "This project contains documentation and templates only."
	@echo "No build, test, or dev targets needed."
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'
