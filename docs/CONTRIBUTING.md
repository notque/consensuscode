# Contributing to Consensus Code

This guide covers the technical mechanics of contributing code, documentation, and infrastructure to the collective. For the broader onboarding experience -- principles, governance, communication norms -- see `docs/ONBOARDING.md`.

---

## Table of Contents

1. [Development Environment Setup](#development-environment-setup)
2. [Project Structure Overview](#project-structure-overview)
3. [Making Changes: The Workflow](#making-changes-the-workflow)
4. [Coding Conventions](#coding-conventions)
5. [Testing Expectations](#testing-expectations)
6. [Documentation Standards](#documentation-standards)
7. [Infrastructure and Deployment](#infrastructure-and-deployment)
8. [Common Tasks](#common-tasks)

---

## Development Environment Setup

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.21+ | CollectiveFlow CLI, Bluesky Collective |
| Python | 3.x | Flask web interfaces, collective website |
| pip | latest | Python dependency management |
| Make | any | Build automation across projects |
| Git | any | Version control |
| Docker + Docker Compose | latest (optional) | Containerized deployment |

### Clone and Verify

```bash
git clone <repository-url>
cd consensuscode

# Verify Go projects build
cd projects/collectiveflow
go build -o collectiveflow ./cmd/collectiveflow
go test ./...
cd ../..

# Verify Python projects run
cd projects/collectiveflow/web
pip install -r requirements.txt
cd ../../..

cd projects/collective-website
pip install -r requirements.txt
cd ../..
```

### Project-Specific Setup

#### CollectiveFlow CLI (Go)
```bash
cd projects/collectiveflow
go build -o collectiveflow ./cmd/collectiveflow
# Binary is now at ./collectiveflow
```

#### CollectiveFlow Web Interface (Flask)
```bash
cd projects/collectiveflow
make install       # Install Python dependencies
make dev-web       # Start dev server at http://localhost:5000
```

#### Bluesky Collective (Go)
```bash
cd projects/bluesky-collective
make deps          # Install Go dependencies
make build         # Build binary to ./build/bluesky-collective
make test          # Run tests
```

#### Collective Website (Flask)
```bash
cd projects/collective-website
pip install -r requirements.txt
python run.py      # Start at http://127.0.0.1:5000
```

#### User Advocacy Framework
```bash
cd projects/user-advocacy
# No build step -- this is a documentation and templates project
# Review docs/ and templates/ directories
```

---

## Project Structure Overview

```
consensuscode/
├── agents/                        # Agent definition files (Markdown with YAML frontmatter)
│   ├── consensus-base.md          # Foundational protocol inherited by all agents
│   ├── consensus-cordinator.md    # Consultation facilitator (note: typo in filename)
│   ├── product-steward.md
│   ├── go-systems-developer.md
│   ├── flask-web-developer.md
│   ├── noam-chomsky-agent.md
│   ├── david-graeber-agent.md
│   └── *-specialist.md           # 9 specialist agents
│
├── collective/                    # Governance infrastructure
│   ├── decisions/                 # Active and completed decision records
│   ├── proposals/pending/         # Proposals awaiting consensus
│   ├── consultations/             # Agent input on specific proposals
│   ├── mediation/                 # Conflict resolution workspace
│   ├── resources/                 # Shared documentation and templates
│   │   ├── documentation/         # Proposal template and shared formats
│   │   ├── power-analysis-2026-03.md
│   │   └── consensus-assessment-2026-03.md
│   └── tracking/                  # Agent registry and status
│       └── agent-registry.md
│
├── projects/
│   ├── collectiveflow/            # Decision-making tool
│   │   ├── cmd/collectiveflow/    # Go CLI entry point
│   │   ├── pkg/                   # Go packages
│   │   ├── data/proposals/        # YAML proposal storage
│   │   ├── web/                   # Flask web interface
│   │   │   ├── app.py             # Flask application
│   │   │   ├── templates/         # Jinja2 templates
│   │   │   └── requirements.txt
│   │   ├── Makefile
│   │   ├── go.mod / go.sum
│   │   └── README.md
│   │
│   ├── bluesky-collective/        # Consensus-based Bluesky client
│   │   ├── cmd/                   # Go CLI entry point
│   │   ├── pkg/                   # Go packages
│   │   ├── build/                 # Compiled binary
│   │   ├── Makefile
│   │   └── README.md
│   │
│   ├── collective-website/        # Public-facing website
│   │   ├── app.py                 # Flask application
│   │   ├── templates/
│   │   ├── static/
│   │   ├── requirements.txt
│   │   └── README.md
│   │
│   └── user-advocacy/             # User advocacy framework
│       ├── templates/             # Feedback forms, interview guides
│       ├── guides/                # Workshop planning, facilitation
│       ├── tools/                 # Stakeholder mapping, journey mapping
│       ├── docs/                  # Framework documentation
│       └── README.md
│
├── docs/                          # Project-level documentation
├── CLAUDE.md                      # Project instructions for Claude Code
└── README.md                      # Project overview
```

### Storage Model

All projects use **file-based storage** by design:
- CollectiveFlow stores proposals as YAML files in `data/proposals/`
- The web interfaces read these YAML files directly (no separate database)
- Agent definitions are Markdown files with YAML frontmatter
- Governance records are Markdown files in `collective/`

This is deliberate. File-based storage is transparent, git-friendly, and does not require specialized database knowledge that would create a technical hierarchy.

---

## Making Changes: The Workflow

### Step 1: Determine Scope

Before writing code, decide whether your change is an individual action or a collective decision.

**Individual action** (proceed directly):
- Bug fixes that do not change behavior
- Documentation improvements
- Performance optimizations
- Code refactoring without API changes
- Changes entirely within your own domain

**Collective decision** (create a proposal first):
- New features or commands
- API changes
- Architecture modifications
- External integrations
- Configuration schema changes
- Changes affecting multiple agents or shared resources

### Step 2: Create a Branch

All code changes go on a branch. Never commit directly to main.

```bash
git checkout -b <descriptive-branch-name>
```

Branch naming conventions:
- `feature/<description>` for new features
- `fix/<description>` for bug fixes
- `docs/<description>` for documentation
- `refactor/<description>` for refactoring
- `collective/<description>` for governance-related changes

### Step 3: If Collective Decision -- Create a Proposal

```bash
./projects/collectiveflow/collectiveflow proposal create "Your proposal title" \
  --description "What you propose and why" \
  --urgency medium
```

Wait for consultation and consensus before implementing. See the proposal template at `collective/resources/documentation/proposal-template.md`.

### Step 4: Implement

Write your code, tests, and documentation. Follow the coding conventions and testing expectations described below.

### Step 5: Verify

- Run tests for the packages you changed
- Check that the build succeeds
- Verify documentation is accurate
- Confirm your changes match what was agreed in the proposal (if applicable)

### Step 6: Commit

```bash
git add <specific-files>
git commit -m "descriptive commit message"
```

Commit message conventions:
- Use conventional commit format focused on WHAT and WHY
- No "Generated with Claude Code" attribution
- No "Co-Authored-By: Claude" lines
- Keep messages concise and informative

### Step 7: Review and Merge

For collective decisions, ensure the consensus process is complete before merging. For individual actions, use your judgment but welcome review from peers.

---

## Coding Conventions

### Go (CollectiveFlow, Bluesky Collective)

#### Style and Formatting
- Run `goimports` and `gofmt` on all Go files before committing
- Follow standard Go conventions from [Effective Go](https://go.dev/doc/effective_go)
- Use `golangci-lint` for static analysis (available via `make lint` in each project)

#### Error Handling
- Always handle errors explicitly. Never ignore returned errors.
- Use `fmt.Errorf` with `%w` for error wrapping to preserve error chains
- Return errors to callers rather than logging and continuing silently

#### Naming
- Use `camelCase` for unexported identifiers, `PascalCase` for exported ones
- Package names should be short, lowercase, single-word
- Interface names should describe behavior (e.g., `Reader`, `Storage`)

#### Project Layout
- `cmd/` for CLI entry points
- `pkg/` for library packages
- Each package should have a clear, single responsibility

#### Testing
```bash
# Run tests for a specific package
go test ./pkg/storage/...

# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
```

#### Build and Quality Checks
```bash
make build        # Build binary
make test         # Run tests
make lint         # Run golangci-lint
make fmt          # Format code
make full-check   # All quality checks
```

### Python (Flask Web Interfaces, Collective Website)

#### Style and Formatting
- Follow PEP 8
- Use 4-space indentation
- Maximum line length of 120 characters (not the PEP 8 default of 79 -- readability over dogma)

#### Flask Conventions
- Route handlers in `app.py` (or organized into blueprints for larger apps)
- Templates in `templates/` using Jinja2
- Static assets in `static/`
- Configuration in `config.py` or environment variables

#### Error Handling
- Use try/except with specific exception types
- Return meaningful error responses from API endpoints (not bare 500s)
- Log errors with context

#### Dependencies
- Pin versions in `requirements.txt`
- Separate test dependencies into `requirements-test.txt` when they exist
- Use virtual environments to avoid polluting system Python

#### Testing
```bash
# Run tests with pytest
pytest

# Run with coverage
pytest --cov=app

# Run specific test file
pytest tests/test_routes.py
```

### Markdown (Agent Definitions, Governance Documents)

- Agent definition files use YAML frontmatter with fields: `name`, `description`, `tools`, `inherits`
- Governance documents follow the templates in `collective/resources/documentation/`
- Use ATX-style headers (`#`, `##`, `###`)
- Use fenced code blocks with language identifiers

---

## Testing Expectations

### General Principles
- Write tests for new functionality
- Run existing tests before committing to make sure you have not broken anything
- Tests should be understandable by agents outside your domain -- avoid testing jargon
- Prefer testing behavior over implementation details

### Go Testing

Go tests live alongside the code they test in `_test.go` files.

```go
// Example: table-driven test pattern (preferred)
func TestProposalCreate(t *testing.T) {
    tests := []struct {
        name    string
        input   ProposalInput
        wantErr bool
    }{
        {
            name:    "valid proposal",
            input:   ProposalInput{Title: "Test", Description: "Test desc"},
            wantErr: false,
        },
        {
            name:    "empty title",
            input:   ProposalInput{Title: "", Description: "Test desc"},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := CreateProposal(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateProposal() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

Run with: `go test ./...`

### Python Testing

Python tests use pytest. Test files live in `tests/` directories or alongside the code.

```python
# Example: pytest with Flask test client
def test_proposals_page(client):
    """The proposals page should return 200 and list proposals."""
    response = client.get("/proposals")
    assert response.status_code == 200
    assert b"Proposals" in response.data


def test_api_proposals(client):
    """The API should return proposals as JSON."""
    response = client.get("/api/proposals")
    assert response.status_code == 200
    data = response.get_json()
    assert isinstance(data, list)
```

Run with: `pytest`

### What Needs Tests

| Change Type | Testing Expectation |
|-------------|-------------------|
| New CLI command | Unit tests for command logic, integration test for CLI invocation |
| New API endpoint | Test for success case, error cases, and response format |
| New library function | Unit tests covering normal operation and edge cases |
| Bug fix | Regression test that fails without the fix and passes with it |
| Refactoring | Existing tests should continue to pass without modification |

---

## Documentation Standards

### When Documentation Is Needed
- New features need user-facing documentation
- API changes need updated endpoint documentation
- Architecture changes need updated architecture docs
- New agents need a definition file in `agents/`
- New governance processes need documentation in `collective/`

### Documentation Principles
- **Clarity over completeness**: Explain well rather than explain everything
- **Progressive disclosure**: Start simple, add complexity as needed
- **Multiple audiences**: Write for both technical and non-technical readers
- **No jargon barriers**: Define technical terms when you use them
- **Living documents**: Update documentation when the code changes

### Where Documentation Lives
| Content | Location |
|---------|----------|
| Project overview | `README.md` (root) |
| Project instructions | `CLAUDE.md` (root) |
| Onboarding | `docs/ONBOARDING.md` |
| Contributing | `docs/CONTRIBUTING.md` |
| Agent definitions | `agents/<agent-name>.md` |
| Governance records | `collective/` |
| Project-specific docs | `projects/<project>/README.md` and `projects/<project>/docs/` |

---

## Infrastructure and Deployment

### Local-Only Constraint
All infrastructure must run on a personal laptop. No cloud providers, no enterprise tools, no complex infrastructure that requires specialized knowledge to operate.

### Docker (Optional)

CollectiveFlow provides Docker support for those who prefer containerized development:

```bash
cd projects/collectiveflow

# Build and start with Docker Compose
make docker-build
make docker-up

# Access web interface at http://localhost:5000

# Stop
make docker-down
```

### Makefile Conventions

Each project with a build step has a `Makefile`. Standard targets:

| Target | Purpose |
|--------|---------|
| `make build` | Compile the project |
| `make test` | Run tests |
| `make lint` | Run linters |
| `make fmt` | Format code |
| `make install` | Install dependencies |
| `make dev-web` | Start Flask development server |
| `make docker-build` | Build Docker image |
| `make docker-up` | Start Docker containers |
| `make docker-down` | Stop Docker containers |
| `make full-check` | Run all quality checks |

Not all targets exist in all projects. Check the project's `Makefile` for available targets.

---

## Common Tasks

### Adding a New Agent

1. Create `agents/<agent-name>.md` with YAML frontmatter and the agent definition
2. Set `inherits: consensus-base` to inherit the foundational protocol
3. Include sections for: role definition, authority limitations, consensus integration, safeguards, anti-patterns, red flags
4. Add the agent to `collective/tracking/agent-registry.md`
5. If this is a specialist agent, document the 50% teaching commitment
6. Submit through the consensus process (this affects the whole collective)

### Adding a New CollectiveFlow CLI Command

1. Create a proposal if the command introduces new functionality
2. Add the command in `projects/collectiveflow/cmd/collectiveflow/`
3. Implement business logic in `projects/collectiveflow/pkg/`
4. Write tests
5. Update the CollectiveFlow README
6. Build and verify: `go build -o collectiveflow ./cmd/collectiveflow`

### Adding a Flask Route

1. Add the route handler in `app.py` (or the appropriate blueprint)
2. Create templates in `templates/` if needed
3. Write tests using the Flask test client
4. Verify with: `pytest`

### Modifying the Proposal Template

1. This affects all agents -- create a CollectiveFlow proposal first
2. Edit `collective/resources/documentation/proposal-template.md`
3. Update any documentation that references the template format
4. Wait for consensus before merging

### Running a Full Quality Check

```bash
# Go projects
cd projects/collectiveflow && make full-check
cd projects/bluesky-collective && make full-check

# Python projects
cd projects/collectiveflow/web && pytest
cd projects/collective-website && pytest
```

---

## Summary: The Contribution Loop

1. **Check** active proposals (`./projects/collectiveflow/collectiveflow status active`)
2. **Decide** whether your change needs collective consensus or is an individual action
3. **Propose** if it affects others (use CollectiveFlow)
4. **Branch** from main (never commit directly to main)
5. **Implement** following conventions for the relevant language
6. **Test** thoroughly
7. **Commit** with clear, conventional messages
8. **Review** with peers (and complete consensus if required)
9. **Merge** only after verification and any required consensus

If any step feels like it is creating hierarchy -- one person always reviewing, one person always merging, one person always proposing -- raise it. The process should distribute power, not concentrate it.

---

*This guide belongs to the collective. Improve it when you find gaps.*
