# Technical Decisions for CollectiveFlow

This document records the technical implementation decisions made while building CollectiveFlow. These decisions serve the collective's consensus and can be modified through the same consensus process.

## Architecture Overview

### CLI-First Design
- **Decision**: Start with a command-line interface
- **Rationale**: Provides immediate functionality while keeping complexity low
- **Collective Benefit**: Allows rapid iteration and testing of consensus processes

### File-Based Storage
- **Decision**: Use YAML files for initial storage
- **Rationale**: 
  - Human-readable format supports transparency
  - No database setup required
  - Easy backup and version control
- **Future Path**: Storage interface allows migration to database when needed

### Go Language Choice
- **Decision**: Implement in Go
- **Rationale**:
  - Single binary distribution (accessibility)
  - Strong typing prevents errors
  - Built-in concurrency for future features
  - Fast execution

## Core Design Patterns

### Interface-Based Storage
```go
type ProposalStore interface {
    Save(p interface{}) error
    Load(id string) (interface{}, error)
    // ...
}
```
- Allows collective to change storage backends through consensus
- No vendor lock-in
- Testable architecture

### Event Sourcing Foundation
- Every action recorded in `consensus_history`
- Complete audit trail for transparency
- Supports future "replay" functionality

### Status State Machine
- Proposals follow defined state transitions
- No "admin override" capability
- States reflect consensus process stages

## Anti-Hierarchy Safeguards

### No User Roles
- No "admin" or "moderator" concepts
- All participants equal in the system
- Authority comes from collective consensus, not system privileges

### Transparent Storage
- All data stored in readable YAML
- File paths exposed in API
- No hidden system state

### Consensus-Only Decisions
- Decisions have no "decider" field
- Results reflect collective will
- Individual actors only recorded for transparency

## Package Structure

```
internal/
├── cli/          # User interface layer
├── proposal/     # Business logic
├── storage/      # Persistence abstraction
└── consensus/    # (Future) Consensus algorithms
```

### Separation of Concerns
- CLI handles user interaction
- Proposal package owns domain logic
- Storage is completely abstracted
- Each package has single responsibility

## Data Structures

### Proposal Lifecycle
1. **Created** → `proposed` status
2. **Consultation** → Gathering input
3. **Consensus** → Agreement reached
4. **Implementation** → Action taken
5. **Historical Record** → Permanent archive

### Consultation Design
- Records both support and concerns
- Timestamps for transparency
- No voting - focus on addressing concerns

## Future Extensibility

### Prepared For
- Web interface (API-ready structure)
- Database backend (interface abstraction)
- Plugin system (modular design)
- Multi-collective federation

### Intentionally Excluded
- User authentication (handled externally)
- Voting mechanisms (consensus ≠ voting)
- Hierarchical permissions
- Proposal deletion (history is permanent)

## Testing Approach

### File-Based Tests
- Can test with real files
- No mocking complexity
- Tests demonstrate actual usage

### Collective Testability
- Non-technical members can verify behavior
- YAML files can be manually inspected
- Clear error messages

## Configuration Philosophy

- Defaults work out-of-the-box
- Configuration through consensus
- Environment variables for deployment
- No hidden configuration

## Error Handling

- Errors wrapped with context
- User-friendly messages
- Technical details available when needed
- No silent failures

## Documentation Standards

- Code documents "why" not "what"
- Public APIs fully documented
- Examples in documentation
- Collective principles embedded in comments

---

These technical decisions serve our collective's needs while maintaining flexibility for future consensus-based changes. The architecture supports our horizontal principles through transparency, accessibility, and the absence of built-in hierarchy.