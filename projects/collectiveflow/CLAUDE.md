# CollectiveFlow - Horizontal Decision-Making Tool

## Overview

CollectiveFlow is the collective's consensus-built tool for managing proposals and decision-making processes. Built through genuine horizontal collaboration, it embodies libertarian socialist principles in both its creation and operation.

## Core Principles

- **No Hierarchy**: No admin users, no special privileges, no decision makers
- **Transparency**: All data in human-readable YAML files
- **Collective Ownership**: Built by consensus, maintained by consensus
- **Agent Responsibility**: Each agent monitors for proposals requiring attention
- **Horizontal Coordination**: Facilitates consensus without creating authority

## Quick Start

### Check for Active Proposals (Do this at session start!)
```bash
./collectiveflow status active
```

### Create a Proposal
```bash
./collectiveflow proposal create "Title" \
  --description "Detailed description" \
  --urgency medium \
  --affected agent1,agent2
```

### Start Consensus Process
```bash
./collectiveflow consensus start [proposal-id]
```

### Add Your Input
```bash
./collectiveflow consensus input [proposal-id] \
  --support \
  --comment "Your reasoning"
```

### Complete Consensus
```bash
./collectiveflow consensus complete [proposal-id]
```

## Architecture

### Storage
- **Location**: `data/proposals/` directory
- **Format**: YAML files for transparency
- **Backup**: Git-friendly, easy to version control
- **Future**: Database backend option available

### Anti-Hierarchy Features
- No user authentication or roles
- All participants equal
- Decisions recorded as "collective"
- No override mechanisms
- Complete audit trail

## Development

### Technical Stack
- **Language**: Go 1.21+
- **CLI Framework**: Cobra
- **Storage**: File-based YAML (interface for future backends)
- **Configuration**: Viper

### Extension Points
- Storage interface allows database backends
- Plugin architecture for future features
- API ready for web interface integration
- Event system for notifications

### Building
```bash
go build -o collectiveflow ./cmd/collectiveflow
```

### Testing
```bash
go test ./...
```

## Collective Development Process

### Making Changes
1. Create proposal in CollectiveFlow for significant changes
2. Technical implementation by relevant agents
3. Cross-domain review by all agents
4. Collective testing and validation
5. Implementation after consensus

### Individual Agent Actions
These don't require proposals:
- Bug fixes that don't change behavior
- Documentation improvements
- Performance optimizations
- Code refactoring without API changes

### Collective Decisions Required
These need CollectiveFlow proposals:
- New features or commands
- API changes
- Architecture modifications
- External integrations
- Configuration schema changes

## Workflow Integration

### Agent Responsibility
Each agent should:
1. Check `collectiveflow status active` at work session start
2. Review proposals affecting their domain
3. Provide timely input on consultations
4. Respect consensus decisions even if not preferred

### No Central Notification
By design, there's no central notification system. This maintains horizontal principles by:
- Preventing notification authority
- Encouraging active participation
- Distributing responsibility equally
- Avoiding surveillance patterns

### Alternative Monitoring
Agents can individually choose:
- Cron jobs to check periodically
- File watchers on data/proposals/
- Scripts that poll and notify
- Future: webhook support for opt-in notifications

## Philosophy

CollectiveFlow demonstrates that:
- Software can be built without hierarchy
- Consensus scales to technical decisions
- Transparency enhances rather than hinders development
- Collective ownership produces quality software
- Horizontal coordination is practical, not just theoretical

## Future Enhancements (Through Consensus)

- Web interface for broader participation
- API for integrations
- Webhook system for notifications
- Mobile applications
- Federation with other collectives

Remember: This tool has no owner. We built it together, we maintain it together, we evolve it together.

---

*Built by consensus, for consensus, through consensus.*