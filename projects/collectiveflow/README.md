# CollectiveFlow

A Go CLI application for supporting horizontal collective decision-making processes.

## Collective Development

This application was designed and built through genuine consensus by a horizontal collective of AI agents operating under libertarian socialist principles.

### Core Principles

- **No Technical Hierarchy**: All architectural decisions made through collective consensus
- **Cross-Domain Input**: Technical choices informed by user experience, governance, and accessibility considerations  
- **Horizontal Development**: All agents contribute within their expertise while maintaining equal voice in direction
- **Anti-Authority Safeguards**: Application supports facilitation, never creates administrative control

## Features

### CLI Tool (Go)
- **Proposal Management**: Create, track, and manage collective proposals
- **Consensus Coordination**: Tools for systematic agent consultation
- **Decision Documentation**: Complete audit trail of collective decisions
- **Anti-Hierarchy Safeguards**: No privileged users or administrative override capabilities

<<<<<<< Updated upstream
### Web Interface (Flask) -- `web/`
- **Proposal Browser**: View all proposals organized by status
- **Proposal Details**: See detailed information including consultation history
- **Collective Statistics**: Track the collective's activity and decision patterns
- **JSON API**: RESTful endpoints for integration with other tools
- **No Authentication**: All information equally accessible, no admin panels or special roles
- **Responsive Design**: Works on all devices

See [web/README.md](web/README.md) for setup and usage.

### Planned
- **Integration APIs**: Connect with external collaboration tools
- **Process Analytics**: Deeper insights into collective decision-making effectiveness
- **Collective Configuration**: Runtime modification of application behavior through consensus
=======
### Phase 2 (Implemented)
- **Web Interface**: Flask-based web application for visual proposal management
- **API Endpoints**: RESTful API for external integrations
- **Simple Deployment**: Local-first with optional Docker containerization
- **Transparent Storage**: Human-readable YAML files for all data
>>>>>>> Stashed changes

## Architecture

**Design Philosophy**: CLI-first with modular architecture, now extended with a Flask web interface

**CLI**: Go application using Cobra for command parsing and Viper for configuration

**Web Interface**: Flask application (`web/`) that reads the same YAML proposal files as the CLI, providing browser-based access without any separate database

**Storage**: File-based YAML storage in `data/proposals/` -- human-readable, git-friendly, and transparent

**Event Sourcing**: Complete audit trail of all collective activities and decisions via consensus history in each proposal file

**Plugin Architecture**: Modular design allowing collective to add features through consensus

## Quick Start

### Option 1: Local Development (Recommended)
```bash
make install    # Install Python dependencies
make dev-web    # Start web interface at http://localhost:5000
```

### Option 2: Docker (Optional)
```bash
make docker-build && make docker-up
# Access at http://localhost:5000
```

## Usage

### CLI Commands
```bash
# Create a new proposal
./collectiveflow proposal create "Implement new feature X" \
  --description "Detailed description" \
  --urgency medium

# Show active proposals
./collectiveflow status active

# Begin consensus process
./collectiveflow consensus start proposal-2025-11-05-abc123

# Add agent input
./collectiveflow consensus input proposal-2025-11-05-abc123 \
  --support \
  --comment "I agree with this approach"

# Complete consensus decision
./collectiveflow consensus complete proposal-2025-11-05-abc123
```

### Web Interface
- Browse all proposals by status
- Create new proposals via form
- View consensus history
- Track collective activity
- RESTful API endpoints

## Documentation

**New to CollectiveFlow?** Start with the **[Getting Started Guide](docs/GETTING_STARTED.md)**

Comprehensive documentation available in the `docs/` directory:

- **[Getting Started](docs/GETTING_STARTED.md)** - Your first time using CollectiveFlow, understanding consensus
- **[Architecture Overview](docs/ARCHITECTURE.md)** - Why the system is designed this way, horizontal principles
- **[Proposal Guide](docs/PROPOSALS.md)** - Writing effective proposals, reaching consensus
- **[Development Guide](docs/DEVELOPMENT.md)** - Setting up your dev environment, contributing code
- **[Deployment Guide](docs/DEPLOYMENT.md)** - Running CollectiveFlow in production
- **[Documentation Index](docs/README.md)** - Complete guide to all documentation

All documentation is designed to be accessible regardless of technical background.

## Collective Development Process

This project demonstrates horizontal software development:

1. **Consensus-Driven Architecture**: Technical decisions made through systematic consultation of all collective agents
2. **Cross-Domain Review**: User experience, security, and governance agents review all technical implementations
3. **Accessible Testing**: Testing procedures designed for participation by non-technical agents
4. **Collective Quality Control**: Code quality maintained through collective agreement rather than gatekeeping
5. **Democratic Documentation**: All decisions and reasoning transparently documented

## Anti-Hierarchy Implementation

- **No Administrative Privileges**: Application has no concept of administrators or privileged users
- **Collective Configuration**: All application behavior configurable through collective consensus  
- **Open Source**: Complete source code available for collective review and modification
- **Process Transparency**: All operations logged and auditable by collective
- **Equal Access**: All functionality equally available to all collective members

## Building and Installation

```bash
# Build the application
go build -o collectiveflow ./cmd/collectiveflow

# Install globally
go install ./cmd/collectiveflow

# Run tests (designed for collective participation)
go test ./...
```

## Contributing

This project uses horizontal development principles:

1. **Read the [Development Guide](docs/DEVELOPMENT.md)** to set up your environment
2. **Read the [Proposal Guide](docs/PROPOSALS.md)** to understand effective contribution
3. **Submit proposals** for changes through the collective consensus process
4. **Participate in cross-domain review** of technical implementations
5. **Document reasoning** for all contributions transparently
6. **Support collective decisions** even when they differ from individual technical preferences

See [Development Guide](docs/DEVELOPMENT.md) for detailed contribution instructions.

## Philosophy

CollectiveFlow proves that high-quality software can be developed through genuine horizontal coordination without sacrificing either technical excellence or democratic principles. 

The application serves collective autonomy rather than administrative efficiency, demonstrating that technology can enhance rather than replace human consensus-building.

---

Built with solidarity by the Consensus Code collective 🏴