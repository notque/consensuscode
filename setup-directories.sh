#!/bin/bash

# Create the core directory structure for Consensus Code

# Main Claude agents directory
mkdir -p .claude/agents

# Collective coordination directories
mkdir -p collective/decisions
mkdir -p collective/proposals/pending
mkdir -p collective/proposals/implemented
mkdir -p collective/consultations
mkdir -p collective/mediation/active
mkdir -p collective/mediation/resolved
mkdir -p collective/resources/shared-tools
mkdir -p collective/resources/documentation
mkdir -p collective/resources/standards
mkdir -p collective/status
mkdir -p collective/tracking

# Documentation directories
mkdir -p docs/go
mkdir -p docs/python
mkdir -p docs/consensus-process

# Project structure
mkdir -p src
mkdir -p tests

echo "Directory structure created for Consensus Code collective!"
echo ""
echo "Next steps:"
echo "1. Place consensus-base.md in .claude/agents/"
echo "2. Add other specialized agents"
echo "3. Initialize collective coordination files"
echo "4. Configure git repository for collaborative development"