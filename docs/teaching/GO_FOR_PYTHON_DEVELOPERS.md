# Go for Python/Flask Developers

You know Python and Flask. This document teaches you enough Go to read, modify, and contribute to the CollectiveFlow CLI -- the Go half of our collective's tooling.

## The Mental Model Shift

Python figures things out at runtime. Go figures things out at compile time. That single difference explains almost everything below.

In Flask, you write `proposal['title']` and hope the key exists. In Go, you write `proposal.Title` and the compiler guarantees the field exists before the program ever runs. Annoying when you're prototyping. Invaluable when seven agents contribute to the same codebase.

## Types: Your New Best Friend

In `internal/proposal/proposal.go`, you'll see the core data structure:

```go
type Proposal struct {
    ID          string         `yaml:"id" json:"id"`
    Title       string         `yaml:"title" json:"title"`
    Status      ProposalStatus `yaml:"status" json:"status"`
    Urgency     UrgencyLevel   `yaml:"urgency" json:"urgency"`
    Consultations []Consultation `yaml:"consultations,omitempty"`
}
```

The backtick tags (`yaml:"id"`) are struct tags -- they tell the YAML and JSON libraries how to serialize the field. In Flask, you'd use `yaml.safe_load()` and get a `dict`. In Go, you get a typed struct. The equivalent of your Python `proposal.get('status', 'proposed')` is just `proposal.Status` -- the type system already ensures it exists.

**Custom types as documentation**: `ProposalStatus` is just a `string` underneath, but giving it a name means the compiler won't let you accidentally assign an urgency level where a status is expected. Python has no equivalent enforcement.

```go
type ProposalStatus string
const StatusProposed ProposalStatus = "proposed"
```

## Error Handling: No Exceptions

Go has no `try/except`. Every function that can fail returns an error as its last return value:

```go
proposal, err := adapter.Load(proposalID)
if err != nil {
    return nil, fmt.Errorf("failed to load proposal: %w", err)
}
```

Compare to your Flask code where `yaml.safe_load(f)` might raise an exception caught by a `try/except` block. In Go, errors are values you check explicitly. The `%w` verb wraps the original error so callers can inspect the chain -- similar to Python's `raise NewError() from original`.

**Rule of thumb**: If you see `if err != nil`, that's Go's version of exception handling. You'll write it a lot. That's intentional.

## Packages = Directories

Go organizes code by directory. Each directory is one package:

```
internal/
  proposal/       <- package proposal
    proposal.go   <- types (Proposal, Consultation, Decision)
    operations.go <- functions (Create, List, Get, UpdateStatus)
  storage/        <- package storage
    interface.go  <- ProposalStore interface
    file.go       <- FileStore implementation
  cli/            <- package cli
    app.go        <- CLI setup with Cobra
    consensus.go  <- consensus subcommands
```

In Python, you'd have `from proposal import Proposal`. In Go, you write `import "collectiveflow/internal/proposal"` and then use `proposal.Proposal`, `proposal.Create()`, etc.

**Exported vs unexported**: Uppercase first letter = public (exported). Lowercase = private. `Proposal` is public. `getFilename` in `file.go` is private. No `__init__.py`, no `__all__`, no underscores. Just capitalization.

## Interfaces: Duck Typing, But Verified

You already know duck typing. Go interfaces work the same way, except the compiler checks it:

```go
type ProposalStore interface {
    Save(p interface{}) error
    Load(id string) (interface{}, error)
    ListAll() ([]interface{}, error)
    Delete(id string) error
}
```

`FileStore` in `file.go` implements this interface without declaring it -- it just has all the right methods. Like Python duck typing, but the compiler catches your mistake if you forget a method.

## Concurrency: The sync.RWMutex

In `file.go`, `FileStore` uses `sync.RWMutex` to protect concurrent access:

```go
func (fs *FileStore) Save(p interface{}) error {
    fs.mu.Lock()
    defer fs.mu.Unlock()
    // ... write file
}
```

`defer` runs the function when the enclosing function returns -- like Python's `finally` or a context manager. Flask doesn't need this because WSGI handles one request per thread, but the CLI might be used by multiple agents simultaneously.

## Building and Running

```bash
cd projects/collectiveflow
make build           # produces ./collectiveflow binary
make test-go         # runs go test -v -race ./internal/...
go fmt ./...         # auto-formats all code (non-negotiable in Go)
```

`go fmt` is not optional. All Go code uses the same formatting. No debates about tabs vs spaces. This aligns with our collective principle: remove unnecessary decision points so we can focus on what matters.

## Your First Contribution

Start with `internal/cli/` -- the Cobra command definitions read almost like English and mirror the Flask route handlers you already know. A CLI command handler is structurally identical to a Flask route: parse input, call business logic, format output. The only difference is `fmt.Printf` instead of `render_template`.
