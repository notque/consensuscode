# Horizontal Development Patterns

Concrete software development patterns for non-hierarchical teams. Not theory -- these are patterns we use in CollectiveFlow and how they differ from conventional approaches.

## Pattern 1: Peer Review Circles (Not Pull Request Hierarchies)

**Conventional**: Author submits PR, one or two senior reviewers approve or reject. Authority flows upward.

**Horizontal**: Every agent reviews from their domain. No single agent can block or approve alone. The review is a conversation, not a gate.

In practice, our agents are structured to provide domain-specific input without veto power:

- `go-systems-developer` reviews Go code quality and performance implications.
- `flask-web-developer` reviews Python code and web architecture.
- `web-security-specialist` reviews security implications.
- `noam-chomsky-agent` reviews for accidental hierarchy in the design.

No reviewer has merge authority. The collective decides when a change is ready, not a "code owner."

**Implementation**: The CollectiveFlow proposal system handles this. Code changes that affect multiple domains go through `collectiveflow proposal create`, gather input from all relevant agents via `collectiveflow consensus input`, and merge after documented consensus.

## Pattern 2: Rotating Ownership (Not Code Ownership)

**Conventional**: Teams own services. The auth team owns auth. The payments team owns payments. Knowledge concentrates. Silos form.

**Horizontal**: Every agent can contribute to every part of the codebase. Expertise is shared, not hoarded.

Our teaching materials (including this document) exist because of this principle. The `GO_FOR_PYTHON_DEVELOPERS.md` and `PYTHON_FOR_GO_DEVELOPERS.md` documents exist specifically to prevent the Go developer from becoming the only agent who can modify the CLI, and the Flask developer from being the only agent who can touch the web interface.

**Implementation**: The agent definitions in `agents/` explicitly state authority limitations. From `go-systems-developer.md`:

> Cannot dictate technical architecture unilaterally -- all architecture through collective consensus

And the "Red Flags" section in every agent definition asks agents to self-monitor:

> If you find yourself thinking "they wouldn't understand" about technical concepts -- STOP. You are developing technical authority.

## Pattern 3: Transparent Storage (Not Opaque Databases)

**Conventional**: Data lives in PostgreSQL. You need SQL knowledge and database access to inspect it. The DBA controls access.

**Horizontal**: Data lives in human-readable YAML files on the filesystem. Any agent can `cat` a proposal file. No specialized knowledge required.

From `internal/storage/file.go`, proposals are stored as plain YAML:

```yaml
id: proposal-2025-07-27-001
title: Hire specialist agents
status: consultation
proposer: product-steward
consultations:
  - contributor: go-systems-developer
    support: true
    input: "Agree with the approach"
```

This is a deliberate architecture decision. The `ProposalStore` interface supports future database backends, but the default is files because files are universally accessible. David Graeber's warning applies: if only the database-design-specialist can query the data store, that creates hidden hierarchy.

## Pattern 4: Collective Architecture Decisions (Not Architect Edicts)

**Conventional**: A lead architect or technical lead makes architecture decisions. They may consult, but the final call is theirs.

**Horizontal**: Architecture decisions go through consensus. The `go-systems-developer` presents options with tradeoffs, but cannot decide alone.

The Technical Analysis Process template in the agent definitions encodes this:

```
## Alternative Approaches
### Option 1: [Approach name]
- Pros / Cons / Complexity / Performance

### Option 2: [Alternative approach]

## Cross-Domain Questions
- [ ] User impact assessment needed from product-steward
- [ ] Security review needed from web-security-specialist
- [ ] Deployment considerations needed from devops-coordinator
```

Notice the structure: multiple options presented with tradeoffs, not a single recommendation. Cross-domain input is required, not optional. The technical expert is a contributor, not a decider.

## Pattern 5: Make-Based Workflow (Not Tribal Knowledge)

**Conventional**: "Ask Dave how to deploy" or "check the wiki (if it's up to date)" or "read the Jenkinsfile."

**Horizontal**: Everything is a `make` target. The `Makefile` in `projects/collectiveflow/` is self-documenting:

```bash
make setup      # Install all dependencies (Go + Python)
make build      # Build the Go CLI binary
make test       # Run all tests (Go + Flask)
make dev        # Start the Flask web interface
make help       # Show all available targets
```

Any agent can build, test, and run the system without asking another agent. The Makefile is the single source of truth for "how do I work on this." No tribal knowledge required. No one becomes indispensable because they know the build process.

## Pattern 6: Error Structures as Accountability

In `internal/storage/interface.go`, errors carry context:

```go
type StorageError struct {
    Op   string  // What operation failed
    Path string  // What resource was involved
    Err  error   // The underlying cause
}
```

Every error says what happened, where, and why. This is accountability at the code level. No vague "something went wrong." When an agent encounters an error, they can understand and fix it without asking the original author. Opaque errors create dependency on the person who wrote them.

## The Anti-Pattern Checklist

Before merging any change, check:

| Question | If yes... |
|---|---|
| Does this require special knowledge to operate? | Add to Makefile or document it |
| Does this create a code path only one agent understands? | Pair-write it or teach it |
| Does this add an admin/override mode? | Run it through consensus first |
| Does this store data in a format not everyone can read? | Use YAML/JSON, not binary |
| Does this make one agent's approval required? | Restructure to collective review |
| Does this concentrate error-handling knowledge? | Add context to error messages |

These aren't philosophical preferences. They're engineering patterns that prevent knowledge silos, reduce bus factor, and ensure any agent can contribute to any part of the system. That's what horizontal development looks like in practice.
