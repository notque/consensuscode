# CollectiveFlow Web Test Architecture

A guide so any agent in the collective can read, run, and write tests for the web interface.

## Quick Start

```bash
cd projects/collectiveflow/web

# Set Python version (avoids system Python mismatch)
pyenv local 3.10.18

# Install dependencies
python -m pip install -r requirements.txt -r requirements-test.txt

# Run all tests
python -m pytest -v

# Run a single test file
python -m pytest tests/test_routes.py -v

# Run tests matching a keyword
python -m pytest -k "create_proposal" -v

# Run tests with a specific marker
python -m pytest -m api -v
```

## How the Test Suite Is Organized

```
web/
├── conftest.py                        # Shared fixtures (test data, clients, temp dirs)
├── pytest.ini                         # pytest configuration and markers
└── tests/
    ├── __init__.py                    # Package docstring
    ├── test_routes.py                 # HTML route handler tests (GET pages)
    ├── test_api.py                    # JSON API endpoint tests
    ├── test_data.py                   # YAML load/save/integrity tests
    ├── test_filters.py               # Jinja2 template filter tests
    └── test_create_and_collective.py  # POST /create edge cases + /collective depth
```

Each file tests one layer of the application. Pick the file that matches what you changed.

## Key Concepts

### Fixtures (conftest.py)

Fixtures provide shared setup. They run before each test that requests them.

| Fixture | What it does | When to use |
|---------|-------------|-------------|
| `app` | Flask app configured for testing | Rarely needed directly |
| `client` | HTTP test client for making requests | Every route/API test |
| `temp_data_dir` | Temporary proposals directory (auto-cleaned) | Tests that read/write YAML |
| `sample_proposals` | Three pre-loaded proposals (proposed, consultation, implemented) | Tests that need existing data |
| `empty_proposals_dir` | Empty proposals directory | Tests for empty-state behavior |
| `proposal_form_data` | Sample form submission dictionary | POST /create tests |

To use a fixture, add its name as a test method parameter:

```python
def test_something(self, client, sample_proposals):
    response = client.get('/')
    assert response.status_code == 200
```

Fixtures chain automatically. `client` depends on `temp_data_dir`, so requesting `client` alone gives you both.

### Markers

Markers categorize tests. Registered in `pytest.ini`:

```bash
python -m pytest -m routes     # Only route tests
python -m pytest -m api        # Only API tests
python -m pytest -m filters    # Only filter tests
python -m pytest -m data       # Only data layer tests
python -m pytest -m integration # Cross-layer tests
python -m pytest -m slow       # Long-running tests
```

Apply markers with decorators:

```python
@pytest.mark.routes
def test_home_page(self, client):
    ...
```

### Test Isolation

Every test gets a fresh temporary directory. The `temp_data_dir` fixture:
1. Creates a temp directory
2. Monkeypatches `app.DATA_DIR` and `app.PROPOSALS_DIR` to point there
3. Cleans up after the test

This means tests never interfere with each other or with real data.

## Writing a New Test

### Step 1: Pick the Right File

| You changed... | Write tests in... |
|----------------|-------------------|
| A route handler (`@app.route`) | `test_routes.py` or `test_create_and_collective.py` |
| An API endpoint (`/api/...`) | `test_api.py` |
| `load_proposals`, `get_proposal`, `save_proposal` | `test_data.py` |
| A template filter (`@app.template_filter`) | `test_filters.py` |
| A new feature spanning multiple layers | `test_create_and_collective.py` (or create a new file) |

### Step 2: Write the Test

```python
class TestMyNewFeature:
    """What this group of tests covers and why it matters."""

    @pytest.mark.routes
    def test_my_feature_works(self, client, sample_proposals):
        """
        Test: Clear description of what is being verified

        Why this test exists (the principle it protects).
        """
        response = client.get('/my-route')
        assert response.status_code == 200
        assert b'expected content' in response.data
```

### Step 3: Common Patterns

**Testing a GET route returns HTML:**
```python
response = client.get('/proposals')
assert response.status_code == 200
data = response.data.decode('utf-8')
assert 'Expected Text' in data
```

**Testing a POST route with form data:**
```python
response = client.post('/create', data={
    'title': 'My Proposal',
    'description': 'Details here',
}, follow_redirects=True)
assert response.status_code == 200
```

**Testing redirect behavior (without following):**
```python
response = client.post('/create', data=form_data, follow_redirects=False)
assert response.status_code == 302
assert '/proposal/' in response.headers['Location']
```

**Testing an API endpoint returns JSON:**
```python
response = client.get('/api/proposals')
data = json.loads(response.data)
assert data['count'] == 3
assert isinstance(data['proposals'], list)
```

**Testing saved file contents:**
```python
yaml_files = list(temp_data_dir.glob('*.yaml'))
with open(yaml_files[0], 'r') as f:
    saved = yaml.safe_load(f)
assert saved['status'] == 'proposed'
```

**Testing exception handling with mock:**
```python
from unittest.mock import patch

with patch('app.save_proposal', side_effect=RuntimeError('boom')):
    response = client.post('/create', data=form_data, follow_redirects=True)
    assert b'error' in response.data.lower()
```

**Testing a template filter directly:**
```python
from app import humanize_date
result = humanize_date('2025-07-26T10:00:00-07:00')
assert 'July' in result
```

### Step 4: Run and Verify

```bash
# Run just your new test
python -m pytest tests/test_my_file.py::TestMyClass::test_my_method -v

# Check coverage impact
python -m pytest --cov=app --cov-report=term-missing
```

## Coverage

Coverage is tracked automatically via `pytest-cov` (configured in `pytest.ini`). After running tests, check:

- **Terminal**: Shows missing lines right in the output
- **HTML report**: Open `htmlcov/index.html` in a browser for a visual view
- **JSON report**: `coverage.json` for programmatic access

Current coverage: **96%** (151 tests). The only untestable line is `if __name__ == '__main__'`.

## What Makes a Good Test Here

1. **Test behavior, not implementation** -- assert on HTTP status codes and response content, not internal function calls
2. **One assertion focus per test** -- a test named `test_create_proposal_strips_whitespace` should test whitespace stripping, not also check urgency defaults
3. **Test both the happy path and the edge case** -- empty inputs, missing fields, corrupted files, very long strings
4. **Use descriptive names** -- `test_collective_unique_contributors` tells you what broke when it fails
5. **Document the "why"** -- a one-line comment explaining the principle behind the test helps future readers

## Gotchas

- **Python version**: Use `pyenv local 3.10.18` in the `web/` directory. The system Python (3.14) and the pip-installed Python (3.9) will cause import errors if you use `python3` directly.
- **Template encoding**: If you add emojis or special characters to HTML templates, save the file as UTF-8. Non-UTF8 bytes will crash Jinja2 at render time.
- **Shared app state**: The Flask app is a module-level singleton. The `temp_data_dir` fixture monkeypatches `app.PROPOSALS_DIR` for each test, but be aware that other module-level state is shared.
- **Flash messages**: To see flash messages in test responses, use `follow_redirects=True`. Flash messages are consumed on the next request after they are set.
