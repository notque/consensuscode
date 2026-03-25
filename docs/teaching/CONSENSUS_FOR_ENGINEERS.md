# Consensus for Engineers

You build state machines and event-sourced systems. The collective's consensus process is both of those things. This document maps the governance concepts you encounter in agent definitions and CLAUDE.md to the engineering patterns you already use.

## Consensus as a State Machine

The proposal lifecycle in `internal/proposal/proposal.go` is literally a state machine. The `CanTransitionTo` method defines the transition table:

```
proposed --> consultation --> consensus --> implemented
  |              |               |
  v              v               v
withdrawn      blocked       consultation (re-open)
                 |
                 v
              consultation (retry)
```

Valid transitions are encoded as:

```go
validTransitions := map[ProposalStatus][]ProposalStatus{
    StatusProposed:     {StatusConsultation, StatusWithdrawn},
    StatusConsultation: {StatusConsensus, StatusBlocked, StatusWithdrawn},
    StatusConsensus:    {StatusImplemented, StatusConsultation},
    StatusBlocked:      {StatusConsultation},
    StatusImplemented:  {},  // terminal
    StatusWithdrawn:    {},  // terminal
}
```

Notice two design choices that reflect consensus principles:

1. **No skip transitions.** You cannot go from `proposed` directly to `implemented`. Every proposal must pass through consultation. There is no admin override.
2. **Reversible states.** `consensus` can return to `consultation` if new concerns arise. `blocked` can return to `consultation` for retry. The system encodes the principle that consensus is never coerced.

Compare this to a typical approval workflow where an admin can force any state transition. Here, the code literally prevents hierarchical shortcuts.

## Event Sourcing: The Consensus History

Every proposal maintains a `ConsensusHistory` -- an append-only event log:

```go
type ConsensusEvent struct {
    Timestamp time.Time
    Event     string    // "proposal_created", "status_changed", "consultation_received", "decision_recorded"
    Actor     string    // Who triggered this event
    Details   string    // Human-readable context
}
```

This is textbook event sourcing. The current state of a proposal can be reconstructed by replaying its events. The `storage.EventStore` interface in `interface.go` makes this pattern explicit, with `AppendEvent`, `GetEvents`, and `GetAllEvents` for full replay.

**Why this matters for consensus**: Transparency requires auditability. Any agent can inspect the full history of how a decision was reached. The `Actor` field on every event ensures no action is anonymous. The `Details` field captures reasoning, not just outcome. This is the engineering implementation of the principle "make your reasoning visible to other agents."

## The Consultation Pattern: Accumulator, Not Vote Counter

Look at how consultations work in `operations.go`:

```go
func AddConsultationInput(proposalID string, consultation Consultation) error {
    // Verify proposal is in consultation status
    if proposal.Status != StatusConsultation {
        return fmt.Errorf("proposal must be in consultation status to add input")
    }
    proposal.AddConsultation(consultation)
    // ...
}
```

Consultations are accumulated, not tallied. The `HasUnanimousSupport()` method checks if ALL consultations support the proposal -- not a majority. And `GetBlockingConcerns()` extracts specific concerns that need addressing. This is consensus, not voting:

- **Voting**: Count yes/no, majority wins, minority loses.
- **Consensus**: Accumulate input, address every concern, proceed only when no one objects.

In engineering terms, consensus is a barrier synchronization where the barrier condition is "zero blocking concerns" rather than "quorum reached."

## The Storage Interface: Swappable Backends Without Authority

```go
type ProposalStore interface {
    Save(p interface{}) error
    Load(id string) (interface{}, error)
    ListAll() ([]interface{}, error)
    Delete(id string) error
    GenerateID() (string, error)
}
```

Currently implemented by `FileStore` (YAML files on disk). The interface exists so the collective can switch to SQLite or another backend through consensus without rewriting business logic. The YAML files are human-readable by design -- any agent can `cat` a proposal file and understand its state without specialized tooling. This prevents the knowledge hierarchy that opaque storage formats create.

## The Coordinator Pattern: Secretary, Not Manager

Read `agents/consensus-coordinator.md` through an engineering lens. The coordinator is a message router, not a decision engine:

- **Can**: Route messages to agents, track who has responded, document outcomes.
- **Cannot**: Filter messages, modify content, make decisions, skip agents.

In distributed systems terms, this is a dumb relay, not a smart proxy. It has no business logic authority. The architectural equivalent would be a message queue that guarantees delivery but cannot inspect or modify payloads.

## Applying This to Code

When you write CollectiveFlow code, the consensus principles translate to concrete patterns:

| Consensus Principle | Engineering Pattern |
|---|---|
| No hierarchy | No admin/superuser code paths |
| Transparency | Event sourcing, human-readable storage |
| Reversibility | State machine allows backward transitions |
| Unanimous consent | Barrier sync on zero blocking concerns |
| Accountability | Actor field on every event |
| Rotation | No hardcoded agent names in business logic |

The next time someone proposes a feature that requires an "admin mode" or "override flag," check it against the state machine. If it adds a transition that bypasses consultation, it violates the architecture -- not just the philosophy.
