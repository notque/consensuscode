# CollectiveFlow Architecture

## System Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                     CollectiveFlow System                           │
│                  (Horizontal Decision-Making Tool)                  │
└─────────────────────────────────────────────────────────────────────┘
                                 │
                                 ▼
        ┌────────────────────────────────────────────┐
        │        User Interaction Layer              │
        ├────────────────────────────────────────────┤
        │                                            │
        │  ┌──────────────┐      ┌───────────────┐  │
        │  │   Go CLI     │      │  Web Browser  │  │
        │  │ collectiveflow│      │  Interface    │  │
        │  └──────────────┘      └───────────────┘  │
        │         │                      │           │
        │         │                      │           │
        └─────────┼──────────────────────┼───────────┘
                  │                      │
                  ▼                      ▼
        ┌─────────────────┐    ┌──────────────────┐
        │  Go Business    │    │  Python Flask    │
        │  Logic Layer    │    │  Web Server      │
        │  (internal/)    │    │  (web/app.py)    │
        └─────────────────┘    └──────────────────┘
                  │                      │
                  │        OR            │
                  │                      │
                  │            ┌──────────────────┐
                  │            │  Go Web Server   │
                  │            │  (internal/web/) │
                  │            └──────────────────┘
                  │                      │
                  └──────────┬───────────┘
                             │
                             ▼
                ┌─────────────────────────┐
                │  Storage Interface      │
                │  (language-neutral)     │
                └─────────────────────────┘
                             │
                             ▼
                ┌─────────────────────────┐
                │   YAML File Storage     │
                │  data/proposals/*.yaml  │
                │                         │
                │  • Language-neutral     │
                │  • Human-readable       │
                │  • Git-friendly         │
                │  • Transparent          │
                └─────────────────────────┘
```

## Component Interaction Patterns

### Pattern 1: CLI Workflow (Go-only path)

```
Agent → collectiveflow CLI → Go Business Logic → YAML Files
          (fast, typed)        (validated)        (persistent)

Example:
$ ./collectiveflow proposal create "Add testing" --urgency medium
   ↓
  Proposal struct created (Go types)
   ↓
  Validated (type-safe)
   ↓
  Saved to data/proposals/proposal-2025-11-05-001.yaml
   ↓
  Success message displayed
```

### Pattern 2: Web Browsing (Python path)

```
User → Web Browser → Python Flask → YAML Files → Jinja2 Templates → HTML
        (rich UI)     (rapid dev)    (readable)    (flexible)

Example:
User visits http://localhost:5000/proposals
   ↓
Flask reads data/proposals/*.yaml
   ↓
Python parses YAML (yaml.safe_load)
   ↓
Jinja2 renders proposals.html
   ↓
User sees community bulletin board interface
```

### Pattern 3: Web Browsing (Go path)

```
User → Web Browser → Go Web Server → YAML Files → Go Templates → HTML
        (rich UI)     (embedded)      (readable)    (compiled)

Example:
User visits http://localhost:8080/proposals
   ↓
Go loads data/proposals/*.yaml
   ↓
Go parses YAML (yaml.Unmarshal)
   ↓
html/template renders proposals.html
   ↓
User sees embedded web interface
```

### Pattern 4: Polyglot Integration (Language-neutral)

```
Any Tool → YAML Files ← Any Other Tool
  (shell)   (neutral)     (python/go/etc)

Examples:
# Shell script monitoring
$ watch -n 60 "cat data/proposals/*.yaml | grep 'status: consultation'"

# Python script for notifications
import yaml
for file in Path('data/proposals').glob('*.yaml'):
    proposal = yaml.safe_load(file.read_text())
    if proposal['status'] == 'consultation':
        send_notification(proposal)

# Go tool for analysis
files, _ := filepath.Glob("data/proposals/*.yaml")
for _, file := range files {
    var p Proposal
    yaml.Unmarshal(readFile(file), &p)
    // Analyze...
}
```

## Layer Responsibilities

### User Interaction Layer

**Go CLI** (`cmd/collectiveflow/main.go`)
- **Responsibility**: Command-line interface for agents
- **Strengths**: Fast, type-safe, single binary
- **Use Cases**: Agent automation, scripting, local workflow

**Web Browsers**
- **Responsibility**: Visual interface for humans
- **Strengths**: Rich UI, familiar interaction
- **Use Cases**: Browsing proposals, reading discussions

### Application Layer

**Go Business Logic** (`internal/`)
```
internal/
├── proposal/     - Domain types and validation
│   ├── proposal.go        (Proposal struct, validation)
│   ├── operations.go      (Create, update, consensus logic)
│   └── storage_adapter.go (Storage interface adapter)
│
├── storage/      - Persistence abstraction
│   ├── interface.go       (ProposalStore interface)
│   └── file.go            (YAML file implementation)
│
├── cli/          - Command-line interface
│   ├── app.go             (Cobra app setup)
│   ├── proposal.go        (Proposal commands)
│   ├── consensus.go       (Consensus commands)
│   ├── status.go          (Status commands)
│   └── web.go             (Web server command)
│
└── web/          - Embedded web server
    ├── server.go          (HTTP handlers)
    ├── templates/         (Go templates)
    └── static/            (CSS, JS)
```

**Responsibilities**:
- ✅ **Type safety**: Go's static typing prevents errors
- ✅ **Validation**: Enforce consensus rules (no admin overrides!)
- ✅ **State machine**: Manage proposal lifecycle
- ✅ **Anti-hierarchy**: Code structure prevents privilege escalation

**Python Flask Web** (`web/app.py`)
```
web/
├── app.py              - Flask application
├── templates/          - Jinja2 templates
│   ├── base.html
│   ├── proposals.html
│   ├── proposal.html
│   └── collective.html
├── static/
│   └── css/style.css   - Community bulletin board styling
└── requirements.txt    - Python dependencies
```

**Responsibilities**:
- ✅ **Rapid development**: Quick template iterations
- ✅ **Rich UI**: Complex HTML/CSS easier in Jinja2
- ✅ **Accessibility**: Python widely understood
- ✅ **Flexibility**: Easy to customize for user feedback

### Storage Layer

**YAML File Storage** (`data/proposals/`)

**Format**:
```yaml
id: proposal-2025-11-05-001
title: "Add comprehensive test coverage"
description: "Implement table-driven tests for domain logic"
proposer: "golang-general-engineer"
date: 2025-11-05T14:30:00Z
status: proposed
urgency: medium
affected_areas:
  - testing
  - quality-assurance
consensus_history:
  - timestamp: 2025-11-05T14:30:00Z
    event: "proposal_created"
    actor: "golang-general-engineer"
consultations: []
```

**Responsibilities**:
- ✅ **Language-neutral**: Any tool can read/write
- ✅ **Human-readable**: Transparency principle
- ✅ **Git-friendly**: Version control ready
- ✅ **Simple**: No database complexity
- ✅ **Accessible**: Standard YAML parsers everywhere

## Horizontal Design Principles in Architecture

### 1. No Privileged Components

```
┌──────────────────────────────────────────────────┐
│  ALL paths have equal authority:                 │
│                                                   │
│  Go CLI       ←→  YAML Files  ←→  Python Web     │
│  Python Script ←→  YAML Files  ←→  Shell Script  │
│  Go Web       ←→  YAML Files  ←→  Custom Tool    │
│                                                   │
│  Storage is the source of truth, not any app     │
└──────────────────────────────────────────────────┘
```

**Anti-pattern avoided**: No "master" application that others must go through

### 2. Knowledge Distribution Through Language Diversity

```
┌────────────────────────────────────────────┐
│  Agent Skills Respected:                   │
│                                             │
│  Go-comfortable agents   → Use Go CLI      │
│  Python-comfortable     → Use Python web   │
│  Shell-comfortable      → Use grep/sed     │
│  Web-comfortable        → Use browser      │
│                                             │
│  All paths equally valid and supported     │
└────────────────────────────────────────────┘
```

**Benefit**: Prevents "rotation illusion" (Graeber) - no false meritocracy

### 3. Transparency Through Storage Format

```
YAML Storage = Democratic Technology
├── Any editor can modify
├── Git tracks all changes
├── Diff shows exact modifications
├── No binary black boxes
├── No database passwords/permissions
└── Collective can inspect everything
```

**Political alignment**: Technology enables transparency, not hierarchy

### 4. Horizontal State Machine

```go
// From internal/proposal/proposal.go
func (p *Proposal) CanTransitionTo(newStatus ProposalStatus) bool {
    // Valid transitions defined by consensus process
    // NO administrative override paths
    // NO "emergency" bypass mechanisms
    // NO "superuser" state changes
}
```

**Code enforces horizontal principles**:
- No `if user.IsAdmin() { allowAnything() }`
- No backdoor state transitions
- Collective consensus = code structure

## Integration Patterns

### Pattern A: Go CLI + Python Web (Current)

```
Development:
  Agent writes proposal → Go CLI → YAML
  User browses → Python Flask → Reads YAML → Beautiful HTML

Deployment:
  $ ./collectiveflow proposal create ...    (Go binary)
  $ flask run                               (Python web)
```

**Strengths**:
- Fast CLI for agents
- Rapid web development iteration
- Language strengths leveraged

### Pattern B: Go CLI + Go Web (Alternative)

```
Development:
  Agent writes proposal → Go CLI → YAML
  User browses → Go Web Server → Reads YAML → Embedded HTML

Deployment:
  $ ./collectiveflow proposal create ...    (Go binary)
  $ ./collectiveflow web serve              (Same binary!)
```

**Strengths**:
- Single binary deployment
- Embedded templates (no external files)
- Faster template rendering

### Pattern C: Polyglot Ecosystem (Future)

```
┌─────────────────────────────────────────────┐
│  Any tool can participate:                  │
│                                              │
│  Go CLI          →  YAML  ←  Python Web     │
│  Python Script   →  YAML  ←  Go Web         │
│  Shell Monitor   →  YAML  ←  Ruby Tool      │
│  Node.js Bot     →  YAML  ←  Rust Analyzer  │
│                                              │
│  Storage format enables unlimited extension  │
└─────────────────────────────────────────────┘
```

**Vision**: Technology-agnostic collective tooling

## Deployment Scenarios

### Scenario 1: Local Development (Current)

```bash
# Terminal 1: Python web server
$ cd projects/collectiveflow/web
$ source venv/bin/activate
$ flask run
 * Running on http://127.0.0.1:5000

# Terminal 2: Agent workflow
$ ./collectiveflow status active
$ ./collectiveflow proposal create "New feature"
```

**Deployment complexity**: Low
**Language requirements**: Python + Go
**Best for**: Active development

### Scenario 2: Single Binary Deployment (Go-only)

```bash
# Single process, embedded web
$ ./collectiveflow web serve --addr :8080
🌐 CollectiveFlow Web Interface
================================
Server starting at http://localhost:8080

# Separate terminal for CLI
$ ./collectiveflow proposal create "New feature"
```

**Deployment complexity**: Very low
**Language requirements**: Go only
**Best for**: Production, simple deployments

### Scenario 3: Systemd Service (Production)

```ini
# /etc/systemd/system/collectiveflow-web.service
[Unit]
Description=CollectiveFlow Web Interface
After=network.target

[Service]
Type=simple
User=collectiveflow
WorkingDirectory=/opt/collectiveflow
ExecStart=/opt/collectiveflow/collectiveflow web serve --addr :8080
Restart=always

[Install]
WantedBy=multi-user.target
```

**Deployment complexity**: Medium
**Language requirements**: Go binary
**Best for**: 24/7 availability

## Anti-Hierarchy Safeguards in Code

### Safeguard 1: No User Model

```go
// ❌ What we DON'T have:
type User struct {
    Username string
    Role     string  // admin, moderator, etc.
    IsAdmin  bool
}

// ✅ What we DO have:
type Consultation struct {
    Contributor string    // Identity for transparency
    Support     bool      // Position, not authority
    Concerns    []string  // Voice, not veto power
}
```

### Safeguard 2: Collective Decisions

```go
// ❌ What we DON'T have:
type Decision struct {
    Approver  string  // Someone who decides
    Decider   string  // Authority figure
}

// ✅ What we DO have:
type Decision struct {
    Result    DecisionResult  // Collective outcome
    Rationale string          // Collective reasoning
    // No "decider" field - decisions are collective
}
```

### Safeguard 3: No Override Mechanisms

```go
// ❌ What we DON'T have:
func (p *Proposal) ForceApprove(admin string) {
    if admin.IsAdmin() {
        p.Status = StatusImplemented
    }
}

// ✅ What we DO have:
func (p *Proposal) CanTransitionTo(newStatus ProposalStatus) bool {
    // Validates against consensus-defined state machine
    // No bypass paths, no emergency overrides
}
```

### Safeguard 4: Transparent Audit Trail

```go
type ConsensusEvent struct {
    Timestamp time.Time  // When it happened
    Event     string     // What happened
    Actor     string     // Who participated (transparency)
    Details   string     // Why it happened
}

// Every action recorded, no deletions, no hiding
p.ConsensusHistory = append(p.ConsensusHistory, event)
```

## Future Architecture Possibilities

### Possibility 1: Database Backend (Through Consensus)

```go
// Storage interface already supports this
type PostgresStore struct {
    db *sql.DB
}

func (s *PostgresStore) Save(p interface{}) error {
    // SQL implementation
}

// Same interface, different backend
// Collective decides when/if to migrate
```

### Possibility 2: API Layer for Federation

```go
// internal/api/server.go (future)
type API struct {
    store storage.ProposalStore
}

func (a *API) ListProposals(w http.ResponseWriter, r *http.Request) {
    proposals, _ := a.store.ListAll()
    json.NewEncoder(w).Encode(proposals)
}

// Enables collective-to-collective communication
```

### Possibility 3: WebSocket for Real-time Updates

```go
// internal/web/websocket.go (future)
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
    // Notify browsers of new proposals
    // Still horizontal - no privileged connections
}
```

### Possibility 4: Plugin System

```go
// internal/plugin/interface.go (future)
type Plugin interface {
    Name() string
    OnProposalCreated(p *Proposal) error
    OnConsensusReached(p *Proposal) error
}

// Collective can add notification methods
// Email, Slack, Bluesky, etc.
// Without modifying core code
```

## Conclusion

CollectiveFlow's architecture demonstrates:

1. **Horizontal principles in code**: Type system prevents hierarchy
2. **Language diversity**: Multiple valid paths, no gatekeeping
3. **Transparent storage**: YAML enables universal participation
4. **Clean abstraction**: Storage interface supports evolution
5. **Accessibility**: Simple tools, no complex infrastructure

The Go implementation **complements** the Python web interface, creating a **polyglot system** that respects agent diversity and prevents knowledge hierarchies.

**This is architectural libertarian socialism in practice.**
