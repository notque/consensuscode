# Python/Flask for Go Developers

You know Go. This document teaches you enough Python and Flask to read, modify, and contribute to the CollectiveFlow web interface -- the Python half of our collective's tooling.

## The Mental Model Shift

Go makes you declare everything. Python infers everything. Where Go gives you compile-time safety, Python gives you rapid iteration. Our Flask app in `web/app.py` is roughly 370 lines for the same proposal management that takes several hundred lines of Go across multiple files. That's not because Python is "better" -- it's because it trades compile-time guarantees for conciseness.

## No Types (Sort Of)

Where your Go code defines:

```go
type Proposal struct {
    ID     string         `yaml:"id"`
    Status ProposalStatus `yaml:"status"`
}
```

Python just uses a dictionary:

```python
proposal = yaml.safe_load(f)  # Returns a dict
title = proposal.get('title', '')  # No compile-time check
```

The `.get('title', '')` pattern provides a default value when the key is missing -- something your Go struct fields handle implicitly by zero values. You'll see this throughout `app.py`. The tradeoff: you can load any YAML shape without defining a struct, but nothing stops you from misspelling `'titl'` until runtime.

## Flask Routes = HTTP Handlers

Flask routes map directly to what you'd build with Go's `net/http` or Cobra commands. In `app.py`:

```python
@app.route('/proposal/<proposal_id>')
def proposal_detail(proposal_id):
    proposal = get_proposal(proposal_id)
    if not proposal:
        return "Proposal not found", 404
    return render_template('proposal.html', proposal=proposal)
```

The `@app.route` decorator is Python's equivalent of registering a handler. The `<proposal_id>` in the URL is like Cobra's `Args: cobra.ExactArgs(1)` -- Flask extracts it and passes it as a function argument. The function returns a response, just like your Cobra `RunE` returns an error.

**Key difference**: Flask returns rendered HTML templates. Your CLI returns `fmt.Printf` text. Same data, different presentation layer.

## Error Handling: Exceptions, Not Return Values

Where Go returns `(result, error)`, Python raises exceptions:

```python
try:
    with open(yaml_file, 'r') as f:
        proposal = yaml.safe_load(f)
except Exception as e:
    print(f"Error loading {yaml_file}: {e}")
```

The `with` statement is a context manager -- similar to Go's `defer file.Close()` but it's scoped to the block. Python's `try/except` catches exceptions that bubble up from any depth, unlike Go's explicit error return at each level.

**In practice**: Flask catches unhandled exceptions and returns a 500 error. You don't need to check every operation. This is convenient but can hide bugs -- the opposite tradeoff from Go's verbose error checking.

## Dictionaries Are Everywhere

Python dicts replace several Go patterns at once:

```python
# This replaces your const block + switch statement
VALID_URGENCIES = {'low', 'medium', 'high', 'emergency'}

# This replaces your map[ProposalStatus][]ProposalStatus
grouped = {
    'consultation': [],
    'proposed': [],
    'consensus': [],
    'implemented': [],
}
```

In `app.py`, the `load_proposals()` function returns a list of dicts. The `grouped` dict-of-lists pattern replaces what you'd do with filtered slices in Go's `List()` function.

## The Project Structure

```
web/
  app.py              <- All routes and logic (single file for simplicity)
  conftest.py          <- pytest fixtures (test setup)
  requirements.txt     <- Dependencies (like go.mod)
  templates/           <- Jinja2 HTML templates
  static/              <- CSS, JS, images
  tests/               <- Test files
```

Unlike Go's multi-package structure, the Flask app is a single `app.py`. This is intentional: the web interface is simpler than the CLI, and splitting it across files would add complexity without benefit. If it grows, Python supports packages (directories with `__init__.py`).

## Testing: pytest vs go test

Our test setup in `conftest.py` uses pytest fixtures, which are Go's `TestMain` and test helpers on steroids:

```python
@pytest.fixture
def client(app, temp_data_dir):
    return app.test_client()

@pytest.fixture
def temp_data_dir(monkeypatch):
    temp_dir = tempfile.mkdtemp(prefix='collectiveflow_test_')
    monkeypatch.setenv('COLLECTIVEFLOW_DATA', temp_dir)
    yield proposals_dir
    shutil.rmtree(temp_dir, ignore_errors=True)
```

Fixtures are dependency-injected by name. When a test function has a parameter called `client`, pytest automatically runs the `client` fixture and passes the result. The `yield` makes it a context manager -- setup runs before yield, cleanup runs after. This replaces Go's `t.Cleanup()` pattern.

Run tests with: `cd web && python3 -m pytest -v`

## Dependencies and Virtual Environments

```bash
python3 -m venv venv          # Create isolated environment (like GOPATH per project)
source venv/bin/activate       # Enter the environment
pip install -r requirements.txt # Install dependencies (like go mod download)
```

The `Makefile` handles this with `make web-setup`. Virtual environments isolate Python packages per project, preventing the dependency conflicts that Go modules solve differently.

## Your First Contribution

Start by reading `app.py` top to bottom -- it's one file. The route handlers (`@app.route`) are structurally identical to your Cobra command handlers: parse input, call a function, return output. The main difference is that output is HTML via `render_template` instead of terminal text via `fmt.Printf`. If you can write a Cobra command, you can write a Flask route.
