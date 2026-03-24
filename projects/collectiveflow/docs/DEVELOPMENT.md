# Development Setup Guide

This guide helps you set up a development environment for CollectiveFlow, whether you're fixing a bug, adding a feature, or just exploring how it works.

**Philosophy**: Development environments should be accessible to everyone. If this guide assumes knowledge you don't have, that's a documentation problem - create a proposal to fix it.

## Prerequisites

What you need before starting:

### Required
- **Go 1.21 or newer** - CollectiveFlow is written in Go
  - Check if you have it: `go version`
  - Install: https://golang.org/dl/
  - Why 1.21+: We use modern Go features like `min()`, `max()`, and the `slices` package

- **Git** - For version control
  - Check if you have it: `git --version`
  - Install: https://git-scm.com/downloads

### Optional but Helpful
- **Python 3.8+** - For the web interface
  - Check: `python3 --version`
  - Install: https://www.python.org/downloads/

- **Make** - For build automation
  - Check: `make --version`
  - Usually pre-installed on Linux/macOS
  - Windows: Install via chocolatey or mingw

- **Text editor or IDE** - Whatever you prefer
  - Suggestions: VS Code, GoLand, vim, emacs, nano
  - No "correct" choice - use what works for you

## Initial Setup

### 1. Get the Code

```bash
# Clone the repository
cd ~/projects  # or wherever you keep code
git clone https://github.com/your-org/consensuscode.git
cd consensuscode/projects/collectiveflow
```

**What just happened?**

You now have a local copy of CollectiveFlow. The directory structure looks like:

```
collectiveflow/
├── cmd/              ← Entrypoint code (main.go)
├── internal/         ← Core logic (not importable by other projects)
│   ├── cli/          ← Command-line interface
│   ├── proposal/     ← Business logic
│   ├── storage/      ← Data persistence
│   └── web/          ← Web server
├── web/              ← Web interface (Python/Flask)
├── data/             ← Where proposals are stored
├── docs/             ← Documentation (you're reading it!)
└── go.mod            ← Go dependencies
```

### 2. Verify Go Setup

```bash
# Check Go version
go version

# Should show: go version go1.21.x or newer
```

If you see an older version, update Go before continuing.

### 3. Download Dependencies

```bash
# From the collectiveflow directory:
go mod download

# This fetches libraries CollectiveFlow depends on
```

**What are these dependencies?**

Check `go.mod` to see exactly what we depend on:
- `github.com/spf13/cobra` - CLI framework
- `gopkg.in/yaml.v3` - YAML parsing
- A few others for configuration and HTTP handling

**Why so few dependencies?** Fewer dependencies means:
- Less code to audit for security
- Easier to understand what the software does
- Fewer opportunities for supply chain attacks
- Prevents knowledge hierarchy (don't need to learn many frameworks)

### 4. Build the CLI

```bash
# Build CollectiveFlow
go build -o collectiveflow ./cmd/collectiveflow

# Test it works
./collectiveflow --help
```

**What is `go build` doing?**

It compiles all the Go code into a single executable binary named `collectiveflow`. This binary has no external dependencies - you can copy it to another machine and it'll just work.

**If build fails**:
- Check Go version: `go version` (need 1.21+)
- Check you're in the right directory: `pwd` (should end in /collectiveflow)
- Read the error message - Go errors are usually helpful
- Create a proposal if the error is confusing (documentation problem!)

### 5. Run Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./internal/proposal
```

**What are tests?**

Tests are code that verifies other code works correctly. They:
- Catch bugs before they reach users
- Document expected behavior
- Give confidence that changes don't break things

**Test output**:
```
ok      collectiveflow/internal/proposal    0.123s
ok      collectiveflow/internal/storage     0.089s
ok      collectiveflow/internal/cli         0.156s
```

This means all tests passed. If you see `FAIL`, something's broken.

**If tests fail**:
- Don't panic - maybe they're expected to fail (check with the collective)
- Read the failure output - it shows what went wrong
- Ask in a proposal if you can't figure it out

## Development Workflow

### Making Changes

#### 1. Understand What You're Changing

Before writing code:

- **For bugs**: Reproduce the problem locally
- **For features**: Understand the collective's consensus (should be a proposal)
- **For refactoring**: Understand current behavior (tests document this)

#### 2. Make Your Changes

Edit files in your favorite editor. The code is organized by responsibility:

**CLI code** (`internal/cli/`):
- Handles user interaction
- Parses command-line arguments
- Formats output for display

**Business logic** (`internal/proposal/`):
- Core consensus rules
- Proposal lifecycle
- Validation logic

**Storage** (`internal/storage/`):
- Reading/writing YAML files
- Data persistence
- Storage interface

**Web interface** (`web/`):
- Flask app (Python)
- HTML templates
- REST API

#### 3. Test Your Changes

```bash
# Run existing tests to ensure nothing broke
go test ./...

# If you added new functionality, write tests for it
# Tests go in files named *_test.go
```

**Example test** (`internal/proposal/proposal_test.go`):

```go
func TestProposalValidation(t *testing.T) {
    p := &Proposal{
        Title: "Test proposal",
        Proposer: "test-agent",
        Urgency: UrgencyMedium,
        Status: StatusProposed,
    }

    err := p.Validate()
    if err != nil {
        t.Errorf("Valid proposal failed validation: %v", err)
    }
}
```

**Why write tests?**

- Documents how your code should work
- Catches regressions (when future changes break your code)
- Shows examples of using your code
- Makes code reviews easier

#### 4. Build and Test Manually

```bash
# Rebuild
go build -o collectiveflow ./cmd/collectiveflow

# Test the actual commands
./collectiveflow proposal create "Test proposal"
./collectiveflow proposal list
./collectiveflow status
```

**Manual testing is important** even with automated tests. Interact with the software as a user would.

### Code Style

CollectiveFlow follows standard Go conventions:

#### Formatting

```bash
# Format all code (Go does this automatically)
go fmt ./...

# Check code quality
go vet ./...
```

**`go fmt`** standardizes formatting. No debates about tabs vs spaces - Go decides.

**`go vet`** catches common mistakes like unreachable code or suspicious constructs.

#### Naming Conventions

**Good names**:
```go
proposalStorage  // Clear what it stores
isValidStatus    // Clear what it checks
GetProposalByID  // Clear what it does
```

**Bad names**:
```go
ps               // What's ps?
validate         // Validate what?
Get              // Get what?
```

**Why naming matters**: Code is read more than written. Clear names are documentation.

#### Comments

Write comments that explain **why**, not **what**:

**Bad comment** (explains what code does - obvious from reading it):
```go
// Add consultation to proposal
proposal.AddConsultation(consultation)
```

**Good comment** (explains why):
```go
// Add consultation and record in history for consensus transparency
proposal.AddConsultation(consultation)
```

**When to comment**:
- Explaining non-obvious decisions
- Warning about edge cases
- Documenting collective consensus that led to this code
- Referencing related proposals

#### Package Documentation

Every package should have a doc comment:

```go
// Package proposal provides core data structures and operations for
// collective proposals. This package embodies horizontal decision-making
// principles - no proposal has inherent authority, and all proposals
// require collective consensus.
package proposal
```

This appears in Go documentation and helps newcomers understand purpose.

## Understanding the Codebase

### Where to Start Reading

If you're new to the codebase, read in this order:

1. **`internal/proposal/proposal.go`** - Core data structures
   - Defines what a proposal is
   - Shows the consensus process model
   - Documents state transitions

2. **`internal/storage/interface.go`** - Storage abstraction
   - Interface that storage implementations must satisfy
   - Shows what operations are needed

3. **`internal/storage/file.go`** - File-based storage
   - Current implementation
   - Shows how YAML files are read/written

4. **`internal/cli/proposal.go`** - CLI commands
   - How user commands map to operations
   - Input validation and error handling

5. **`internal/cli/consensus.go`** - Consensus commands
   - Consultation tracking
   - Status checking

### Common Patterns

#### Error Handling

CollectiveFlow uses Go's standard error handling:

```go
proposal, err := storage.Load(proposalID)
if err != nil {
    return fmt.Errorf("failed to load proposal: %w", err)
}
```

**`%w`** wraps the error, preserving context. Error messages should help users fix problems:

**Bad error**: "Invalid input"
**Good error**: "Proposal ID must be in format proposal-YYYY-MM-DD-NNN, got: abc123"

#### Validation

Validate early, at the boundaries:

```go
func CreateProposal(newProposal proposal.New) error {
    // Validate input before creating
    if newProposal.Title == "" {
        return fmt.Errorf("proposal title cannot be empty")
    }

    // Create proposal
    p := &proposal.Proposal{
        Title: newProposal.Title,
        // ...
    }

    // Validate the complete structure
    if err := p.Validate(); err != nil {
        return fmt.Errorf("proposal validation failed: %w", err)
    }

    return storage.Save(p)
}
```

#### Interfaces for Flexibility

CollectiveFlow uses interfaces to allow future changes:

```go
type ProposalStore interface {
    Save(p interface{}) error
    Load(id string) (interface{}, error)
    List(filter interface{}) ([]interface{}, error)
}
```

This means:
- Current file-based storage can be swapped for database
- Tests can use mock storage
- Multiple implementations can coexist

## Working with the Web Interface

The web interface is a separate Flask application that reads the same data.

### Setup

```bash
cd web

# Create virtual environment
python3 -m venv venv

# Activate it
source venv/bin/activate     # On Linux/macOS
# or
venv\Scripts\activate        # On Windows

# Install dependencies
pip install -r requirements.txt
```

**What's a virtual environment?**

A virtual environment isolates Python dependencies for this project. Without it, installing Flask would affect your entire system. With it, each project has its own dependencies.

### Running the Web Interface

```bash
# Make sure you're in the web/ directory with venv activated
python app.py

# Or specify a different port if 5000 is in use
python -c "from app import app; app.run(debug=True, port=5001)"
```

Open http://localhost:5000 in your browser.

### Web Interface Structure

```
web/
├── app.py              ← Flask application
├── templates/          ← HTML templates
│   ├── base.html       ← Base template (header, footer)
│   ├── index.html      ← Home page
│   ├── proposals.html  ← Proposal list
│   └── detail.html     ← Proposal detail
└── static/             ← CSS, JavaScript, images
    └── styles.css      ← Custom styles (Tailwind CDN used)
```

### Making Changes to Web Interface

1. Edit templates or `app.py`
2. Flask auto-reloads in debug mode (see changes immediately)
3. Test in browser
4. No need to rebuild (Python is interpreted)

## Debugging

### CLI Debugging

**Print debugging** (simple and effective):

```go
fmt.Printf("DEBUG: proposalID = %s\n", proposalID)
fmt.Printf("DEBUG: proposal = %+v\n", proposal)
```

**`%+v`** prints structs with field names, very useful for debugging.

**Delve debugger** (more powerful):

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the CLI
dlv debug ./cmd/collectiveflow -- proposal list
```

Set breakpoints, step through code, inspect variables interactively.

### Web Interface Debugging

Flask runs in debug mode by default during development:

```python
# In app.py
if __name__ == '__main__':
    app.run(debug=True)  # Debug mode enabled
```

This gives you:
- Detailed error pages with stack traces
- Auto-reload when you change code
- Interactive debugger in the browser (!)

**Never enable debug mode in production** - it's a security risk.

## Testing

### Writing Tests

Tests go in files named `*_test.go`:

```go
package proposal

import "testing"

func TestProposalCreation(t *testing.T) {
    p := &Proposal{
        Title: "Test",
        Proposer: "test-agent",
        Urgency: UrgencyLow,
        Status: StatusProposed,
    }

    if err := p.Validate(); err != nil {
        t.Errorf("Validation failed: %v", err)
    }

    if p.Title != "Test" {
        t.Errorf("Expected title 'Test', got '%s'", p.Title)
    }
}
```

### Running Specific Tests

```bash
# Run one test function
go test -run TestProposalCreation ./internal/proposal

# Run tests matching a pattern
go test -run "TestProposal.*" ./...

# Run tests with coverage
go test -cover ./...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Coverage shows which code is tested. Aim for high coverage, but don't obsess over 100%.

### Table-Driven Tests

Go pattern for testing multiple cases:

```go
func TestProposalValidation(t *testing.T) {
    tests := []struct {
        name      string
        proposal  Proposal
        wantError bool
    }{
        {
            name: "valid proposal",
            proposal: Proposal{
                Title: "Test",
                Proposer: "agent",
                Urgency: UrgencyMedium,
                Status: StatusProposed,
            },
            wantError: false,
        },
        {
            name: "missing title",
            proposal: Proposal{
                Proposer: "agent",
                Urgency: UrgencyMedium,
                Status: StatusProposed,
            },
            wantError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.proposal.Validate()
            if (err != nil) != tt.wantError {
                t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
            }
        })
    }
}
```

This pattern makes it easy to test many cases without duplicating test code.

## Common Development Tasks

### Adding a New CLI Command

1. **Create the command function** in `internal/cli/`:

```go
// internal/cli/example.go
func newExampleCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "example",
        Short: "Example command",
        Long:  "Detailed description of what this does",
        RunE: func(cmd *cobra.Command, args []string) error {
            // Implementation here
            return nil
        },
    }

    // Add flags if needed
    cmd.Flags().String("option", "", "Description of option")

    return cmd
}
```

2. **Register the command** in `internal/cli/app.go`:

```go
rootCmd.AddCommand(newExampleCmd())
```

3. **Test it**:

```bash
go build -o collectiveflow ./cmd/collectiveflow
./collectiveflow example --help
```

### Adding Business Logic

1. **Add the functionality** to appropriate package
2. **Write tests** for it
3. **Update CLI** to expose it
4. **Document** in code comments and user docs

### Changing Data Structures

**Be careful!** Changing data structures affects stored proposals.

1. **Add fields** carefully:
   ```go
   // Safe - new optional field
   NewField string `yaml:"new_field,omitempty"`
   ```

2. **Don't remove fields** without migration plan
3. **Don't change field types** without migration
4. **Test** with existing proposal files

If you must break compatibility:
- Create a migration tool
- Document the change
- Provide upgrade path
- Requires collective consensus

## Development Anti-Patterns

### Anti-Pattern 1: Secretly Adding Hierarchy

```go
// Bad - introduces authority concept
type User struct {
    Name     string
    IsAdmin  bool  // ← NO! Violates horizontal principles
}
```

If you're tempted to add "admin" or "moderator" concepts, stop. Create a proposal to discuss why you think it's needed - there's usually a horizontal alternative.

### Anti-Pattern 2: Hiding Information

```go
// Bad - data not transparent
func (s *storage) SaveProposal(p Proposal) error {
    // Save to binary format that humans can't read
    return saveBinary(p)
}
```

All data should be human-readable. If you need binary format for performance, it should be an option, not the only choice.

### Anti-Pattern 3: Shortcuts Around Consensus

```go
// Bad - bypasses consensus rules
func (p *Proposal) ForceApprove() {
    p.Status = StatusConsensus  // ← NO! Skips consultation
}
```

Code should enforce consensus, never bypass it.

### Anti-Pattern 4: Complexity for Complexity's Sake

```go
// Bad - overengineered
type ProposalFactoryStrategyFactoryBuilder struct {
    // 500 lines of abstraction
}
```

Simpler code is more accessible. Avoid enterprise patterns unless they solve real problems.

## Getting Help

### Resources

1. **Go documentation**: https://golang.org/doc/
2. **Effective Go**: https://golang.org/doc/effective_go (style guide)
3. **Go by Example**: https://gobyexample.com/ (practical examples)
4. **Cobra CLI docs**: https://cobra.dev/

### Asking Questions

If you're stuck:

1. **Check documentation** - might already be answered
2. **Look at similar code** - how do existing parts work?
3. **Create a proposal** - "Help understanding X" is valid
4. **Share debugging steps** - what have you tried?

**Don't feel bad about asking** - if the code or docs are confusing, that's a collective problem to fix.

## Making Your First Contribution

Start small:

**Good first contributions**:
- Fix a typo in documentation
- Add a test case for existing functionality
- Improve error messages
- Add examples to documentation

**Not good first contributions**:
- Rewrite major components
- Add complex features
- Change core architecture

**Process**:
1. Make your change
2. Test it
3. Create a proposal describing the change
4. Get collective consensus
5. Merge when approved

## Summary

Development environment setup:
```bash
# Clone code
git clone [repo]
cd consensuscode/projects/collectiveflow

# Build CLI
go build -o collectiveflow ./cmd/collectiveflow

# Run tests
go test ./...

# Set up web interface (optional)
cd web
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python app.py
```

Development workflow:
1. Understand what you're changing
2. Make changes
3. Test (automated and manual)
4. Get collective consensus
5. Merge

Remember: If this guide is confusing, that's a documentation problem. Create a proposal to improve it.

---

Happy developing! Your contributions help build software that serves horizontal principles.
