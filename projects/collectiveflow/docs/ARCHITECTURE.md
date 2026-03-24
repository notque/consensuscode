# CollectiveFlow Architecture

This document explains how CollectiveFlow is designed and, more importantly, **why** it's designed that way. The architecture embodies horizontal principles - every technical decision serves the goal of preventing hierarchy.

## Core Design Philosophy

### Principle 1: Transparency Over Convenience

CollectiveFlow stores everything in human-readable YAML files, not a database. This is intentional.

**Why?**

- **Anyone can read the data** without special tools or database access
- **Git-friendly** - you can version control and diff proposals
- **No database expertise required** - no one gains power through specialized knowledge
- **Audit trail by default** - file history shows what changed and when
- **No vendor lock-in** - your data is in plain text, forever accessible

**Trade-off**: File-based storage is slower than databases for complex queries. We accept this because transparency matters more than performance for collective decision-making.

### Principle 2: No Administrative Override

CollectiveFlow has no concept of "admin users" or privileged accounts. This is not an oversight - it's the point.

**Why?**

- **Authority corrupts consensus** - if someone can override decisions, it's not consensus
- **Prevents hierarchy creep** - no "just this once" admin powers that become permanent
- **Forces collective solutions** - can't shortcut to "ask the admin"
- **Equal access** - everyone has the same capabilities

**Trade-off**: This means there's no "emergency override" if something goes wrong. The collective must solve problems collectively, which takes longer but maintains horizontal principles.

### Principle 3: Complete History, No Deletions

You cannot delete proposals or consultations in CollectiveFlow. Proposals can be withdrawn, but the record remains.

**Why?**

- **Accountability** - everyone can see what was proposed and why
- **Learning** - past decisions inform future ones
- **Trust** - no one can hide inconvenient history
- **Prevents revisionism** - the record can't be changed to favor current politics

**Trade-off**: Mistakes are permanent. If you create a proposal with errors, you withdraw it and create a new one. This might feel harsh, but it serves transparency.

## Technical Architecture

### Component Overview

```
┌─────────────────────────────────────────────────┐
│              CollectiveFlow CLI                  │
│  (Command-line interface for all operations)     │
└────────────────┬────────────────────────────────┘
                 │
                 ├─────────────────┐
                 │                 │
         ┌───────▼────────┐  ┌────▼──────────────┐
         │   Proposal     │  │   Consensus       │
         │   Management   │  │   Coordination    │
         └───────┬────────┘  └────┬──────────────┘
                 │                │
                 └────────┬───────┘
                          │
                  ┌───────▼────────┐
                  │    Storage     │
                  │   Interface    │
                  └───────┬────────┘
                          │
              ┌───────────┴───────────┐
              │                       │
      ┌───────▼────────┐    ┌────────▼────────┐
      │  File Storage  │    │ Future: Database│
      │  (YAML files)  │    │   (Not yet)     │
      └────────────────┘    └─────────────────┘
```

### Why This Structure?

**Separation of concerns**: Each layer has one responsibility:
- **CLI**: Handle user interaction, display formatting
- **Business logic**: Implement consensus rules and proposal lifecycle
- **Storage interface**: Abstract how data is saved/loaded
- **Storage implementation**: Actually read/write data

**Benefits**:
- Anyone can understand one layer without understanding all layers
- Changes in one layer don't break other layers
- Testing is easier - can test business logic without file I/O
- Future features (database, web interface) don't require rewriting everything

**This prevents knowledge hierarchy** - no one needs to understand the entire system to contribute.

## Data Model

### Proposal Lifecycle

```
                    ┌──────────┐
                    │ proposed │ ← Initial state
                    └─────┬────┘
                          │
                    ┌─────▼──────────┐
                    │  consultation  │ ← Gathering input
                    └─────┬────┬─────┘
                          │    │
              ┌───────────┘    └────────────┐
              │                              │
         ┌────▼────┐                   ┌────▼────┐
         │consensus│                   │ blocked │
         └────┬────┘                   └─────────┘
              │                              ↓
              │                         (can retry)
              │
      ┌───────▼────────┐
      │  implemented   │ ← Terminal state
      └────────────────┘
```

**Why this lifecycle?**

- **proposed**: Prevents premature discussion - proposer writes it clearly first
- **consultation**: Explicit phase for gathering input
- **consensus**: Marks agreement reached, ready to implement
- **blocked**: Honest acknowledgment when consensus fails
- **implemented**: Clear record of completion

**No "draft" status**: Either propose it or don't. This prevents endless preparation instead of engaging the collective.

**Can return from consensus to consultation**: If new concerns emerge, we address them. Progress isn't one-way.

### Consultation Structure

Each consultation records:

```yaml
contributor: "agent-name"      # Who provided input
timestamp: 2025-11-05T10:00:00Z # When
input: "Detailed thoughts..."   # Their perspective
concerns:                       # Specific worries
  - "This might cause X"
  - "Have we considered Y?"
support: true                   # Whether they support
```

**Why this structure?**

- **Named contributors**: Accountability and context
- **Timestamp**: See how thinking evolved
- **Input and concerns separated**: Concerns get explicit attention
- **Boolean support**: Clear signal, not ambiguous ratings
- **Concerns are not voting**: Multiple concerns from one person don't "count more"

### No "Voting" Field

You won't find vote counts in CollectiveFlow. This is intentional.

**Why no voting?**

- **Consensus is not democracy** - it's not about majority
- **Concerns matter, not counts** - one blocking concern is enough
- **Prevents game-playing** - can't "stack" votes or mobilize voting blocs
- **Forces engagement** - must address concerns, not out-vote them

## Anti-Hierarchy Safeguards

### 1. No User Roles or Permissions

```go
// You won't find this in CollectiveFlow:
type User struct {
    Name  string
    Role  Role  // ← No roles!
    Perms []Permission  // ← No permissions!
}
```

There are no "coordinator" accounts with special powers. The "Consensus Coordinator" is a **role**, not a **user type**. Anyone can perform coordination functions.

**Why?**

Roles in software tend to become roles in the collective. If the system knows about "admin" and "member," people start thinking in those terms.

### 2. All Data is Readable

```bash
# Anyone can do this:
cat data/proposals/proposal-2025-11-05-001.yaml

# It looks like:
id: proposal-2025-11-05-001
title: "Improve documentation"
status: consultation
consultations:
  - contributor: alice
    support: true
    input: "Great idea"
```

**Why?**

- **Technical knowledge doesn't gate access** - no "must know SQL" to see proposals
- **Transparency by design** - can't hide data in binary formats
- **Tool-independent** - don't need CollectiveFlow running to read data
- **Trust through visibility** - everyone can verify what the tool says

### 3. Immutable History

The `consensus_history` field records every event:

```yaml
consensus_history:
  - timestamp: 2025-11-05T10:00:00Z
    event: "consultation_started"
    actor: "consensus-coordinator"
  - timestamp: 2025-11-05T11:00:00Z
    event: "consultation_received"
    actor: "alice"
    details: "Support: true"
```

**Why?**

- **Accountability** - see who did what
- **Timeline reconstruction** - understand how decisions evolved
- **Prevents revisionism** - can't pretend consultation didn't happen
- **Debugging** - when process breaks, history shows what went wrong

### 4. No Implicit Defaults That Create Hierarchy

Bad example (from other tools):
```yaml
urgency: high  # Automatically notifies "leaders"
```

CollectiveFlow's approach:
```yaml
urgency: high  # Just metadata, no automatic privilege
```

**Why?**

Automatic escalation to "leaders" based on urgency creates hierarchy. In CollectiveFlow, urgency is information for the collective, not a trigger for authority.

## Storage Interface Design

```go
type ProposalStore interface {
    Save(p interface{}) error
    Load(id string) (interface{}, error)
    List(filter interface{}) ([]interface{}, error)
}
```

**Why an interface?**

This enables future changes through consensus:

- Want a database for performance? Implement `ProposalStore` for PostgreSQL
- Want distributed storage? Implement `ProposalStore` for git-based sync
- Want different file format? Implement `ProposalStore` for JSON or TOML

**The interface prevents lock-in** to current technical choices. The collective can evolve the tool without rewriting everything.

## Performance vs Principles

CollectiveFlow makes deliberate trade-offs:

### We Choose Transparency Over Speed

File-based storage is slower than databases. For 1000 proposals, database queries would be faster. We accept this because:

- Most collectives don't have 1000 active proposals
- Transparency matters more than milliseconds
- When performance becomes a real problem, storage interface allows migration
- Optimization is the last resort, not the first tool

### We Choose Accessibility Over Features

CollectiveFlow doesn't have:
- Real-time notifications
- Rich text editing
- Advanced search
- Analytics dashboards
- Mobile apps (yet)

**Why?**

Each feature creates complexity. Complexity creates knowledge hierarchies - some people understand the features, others don't. We add features through consensus only when the benefit clearly outweighs the complexity cost.

## Future Extensions

The architecture supports future collective decisions:

### Database Backend (Prepared, Not Implemented)

```go
// Already possible through storage interface:
type PostgresStore struct {
    db *sql.DB
}

func (s *PostgresStore) Save(p interface{}) error {
    // Implement database storage
}
```

**When might we want this?**

- Collective grows to hundreds of active proposals
- Complex queries become necessary (show all proposals by multiple tags)
- Performance becomes a real barrier to participation

**Decision**: Collective consensus required before implementation.

### Web Interface (In Progress)

The CLI architecture allows adding a web interface without changing core logic. The web interface uses the same storage layer, maintaining transparency.

**Why web interface?**

- Broader accessibility - not everyone is comfortable with CLI
- Visual representation of consensus state
- Easier for newcomers to explore without learning commands

**Critical**: Web interface maintains horizontal principles - no features the CLI doesn't have.

### Federation (Possible Future)

```
Collective A          Collective B
    ↓                      ↓
  [Sync proposals through git or API]
```

Multiple collectives could share proposals while maintaining autonomy.

**Why?**

Enable solidarity between collectives without creating centralized coordination.

**Challenges**:
- How to handle conflicting consensus results?
- Who decides sync frequency?
- What about network failures?

**Decision**: Requires collective consensus and careful design to avoid hidden centralization.

## Learning from the Architecture

### For Developers

The code structure demonstrates:
- **Interfaces over implementations** - keeps options open
- **Separation of concerns** - easier to understand and modify
- **Explicit over implicit** - no "magic" behavior
- **Documentation as first-class** - explains "why" not just "what"

### For Collectives

The architecture shows:
- **Technology can serve horizontal principles** - it's not just "best practices"
- **Trade-offs should align with values** - we optimize for transparency, not speed
- **Complexity is a political choice** - simpler tools prevent knowledge hierarchies
- **Open design enables adaptation** - collectives can modify as needs change

## Questions About the Architecture?

If this document leaves you confused about design decisions, that's a documentation problem. Create a proposal asking for clarification, and the collective will address it.

Architecture documents should demystify, not gatekeep.

---

Remember: This architecture serves consensus, not technical elegance. If it gets in the way of horizontal coordination, the collective can change it.
