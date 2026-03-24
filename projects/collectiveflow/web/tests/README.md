# CollectiveFlow Web Application - Test Suite

This test suite validates the CollectiveFlow Flask web application, embodying our collective principles of transparency, accessibility, and horizontal organization.

## Philosophy

These tests are designed following our collective values:

- **Knowledge Democratization**: Tests are clearly documented so anyone can understand them
- **Transparency**: Test coverage shows exactly what is and isn't validated
- **Accessibility**: Tests use clear naming and extensive comments
- **No Hierarchy**: All tests are equal - no "important" vs "minor" tests
- **Teaching Tools**: Tests serve as documentation and learning resources

## Quick Start

### Install Test Dependencies

```bash
pip install -r requirements-test.txt
```

This installs pytest and all testing utilities.

### Run All Tests

```bash
# Run all tests with coverage report
pytest

# Run with verbose output
pytest -v

# Run specific test file
pytest tests/test_routes.py

# Run specific test class
pytest tests/test_routes.py::TestHomeRoute

# Run specific test
pytest tests/test_routes.py::TestHomeRoute::test_home_page_loads_successfully
```

### Run Tests by Marker

Tests are organized with markers for easy filtering:

```bash
# Run only route tests
pytest -m routes

# Run only API tests
pytest -m api

# Run only filter tests
pytest -m filters

# Run only data/YAML tests
pytest -m data

# Run integration tests
pytest -m integration

# Skip slow tests
pytest -m "not slow"
```

## Test Organization

### Test Files

- **`test_routes.py`**: Route handlers and HTML responses
- **`test_filters.py`**: Jinja2 custom filters
- **`test_data.py`**: YAML loading, saving, and data persistence
- **`test_api.py`**: JSON API endpoints

### Test Structure

Each test file follows this pattern:

```python
class TestFeatureArea:
    """
    Tests for a specific feature or component.
    Class docstring explains what's being tested.
    """

    @pytest.mark.category
    def test_specific_behavior(self, fixtures):
        """
        Test: Clear description of what this test validates

        Additional explanation of why this test matters
        and what it's checking.
        """
        # Test implementation with clear steps
        pass
```

### Fixtures

Shared test data and setup is in `conftest.py`:

- **`app`**: Flask application configured for testing
- **`client`**: Flask test client for making requests
- **`temp_data_dir`**: Temporary directory for test data (auto-cleaned)
- **`sample_proposals`**: Pre-created test proposals
- **`empty_proposals_dir`**: Empty data directory for testing edge cases
- **`proposal_form_data`**: Sample form submission data

## Coverage

### Current Coverage Targets

The test suite aims for comprehensive coverage:

- **Routes**: All HTTP endpoints
- **Templates**: Template rendering and variable usage
- **Filters**: All custom Jinja2 filters
- **Data Operations**: YAML loading, saving, and integrity
- **API**: All JSON endpoints and error cases
- **Error Handling**: 404s, malformed data, edge cases

### Viewing Coverage Reports

```bash
# Run tests with coverage
pytest --cov=app --cov-report=html

# Open HTML coverage report
open htmlcov/index.html
```

Coverage reports show:
- Which lines of code are tested
- Which branches are covered
- Which functions are called by tests

## Test Categories

### Unit Tests

Test individual functions and components in isolation:

```bash
pytest -m "not integration"
```

Examples:
- Testing a single filter function
- Testing data loading from one YAML file
- Testing one API endpoint

### Integration Tests

Test multiple components working together:

```bash
pytest -m integration
```

Examples:
- Template rendering with filters and data loading
- API returning data that was saved through the web interface
- End-to-end proposal creation flow

### Slow Tests

Tests that take longer to run (e.g., testing with many proposals):

```bash
# Run fast tests only
pytest -m "not slow"

# Run slow tests
pytest -m slow
```

## Writing New Tests

### Test Naming Convention

Test names should clearly describe what they test:

```python
# Good names
def test_home_page_loads_successfully()
def test_api_returns_json_content_type()
def test_humanize_date_handles_invalid_input()

# Less clear names
def test_home()
def test_api()
def test_filter()
```

### Test Documentation

Every test should have:

1. **Descriptive name**: What behavior is being tested
2. **Docstring**: Why this test matters and what it validates
3. **Clear steps**: Comments explaining non-obvious parts
4. **Assertions with messages**: Help debugging when tests fail

Example:

```python
@pytest.mark.routes
def test_home_page_shows_recent_proposals(self, client, sample_proposals):
    """
    Test: Home page displays recent proposals

    Members need to see recent activity to stay informed about
    collective decision-making. This test ensures the home page
    includes proposal listings.
    """
    response = client.get('/')
    data = response.data.decode('utf-8')

    # Check that we see proposal titles from test data
    assert 'Test Proposal: Simple Example' in data, \
        "Home page should display proposal titles"
```

### Using Fixtures

Leverage existing fixtures for common test needs:

```python
def test_something(client, sample_proposals, temp_data_dir):
    # client: Make HTTP requests
    response = client.get('/proposals')

    # sample_proposals: Pre-created test data
    # (automatically created in temp_data_dir)

    # temp_data_dir: Path to temporary data directory
    # Use if you need to inspect files directly
    assert (temp_data_dir / 'test-proposal-001.yaml').exists()
```

### Marking Tests

Add appropriate markers:

```python
@pytest.mark.routes      # Route handler test
@pytest.mark.api        # API endpoint test
@pytest.mark.filters    # Filter test
@pytest.mark.data       # Data/YAML test
@pytest.mark.integration  # Integration test
@pytest.mark.slow       # Long-running test
```

## Testing Principles

### Test Isolation

Each test should be independent:

- ✅ Tests can run in any order
- ✅ Each test gets fresh data via fixtures
- ✅ Temporary directories are cleaned up
- ❌ Don't rely on state from previous tests
- ❌ Don't modify shared fixtures

### Clear Assertions

Assertions should be obvious:

```python
# Good: Clear what's being checked
assert response.status_code == 200
assert 'CollectiveFlow' in response.data.decode('utf-8')

# Less clear: What does True mean?
assert check_response(response)
```

### Test One Thing

Each test should validate one specific behavior:

```python
# Good: Tests one thing
def test_home_page_returns_200():
    response = client.get('/')
    assert response.status_code == 200

def test_home_page_contains_title():
    response = client.get('/')
    assert b'CollectiveFlow' in response.data

# Less good: Tests multiple unrelated things
def test_home_page():
    response = client.get('/')
    assert response.status_code == 200
    assert b'CollectiveFlow' in response.data
    assert len(response.data) > 0
    # ... many more assertions
```

### Handle Edge Cases

Good tests cover edge cases:

- Empty data (no proposals)
- Invalid input (malformed YAML)
- Missing data (nonexistent proposal IDs)
- Special characters (unicode, emojis)
- Boundary conditions (very long strings)

## Continuous Integration

These tests are designed to run in CI/CD pipelines:

```bash
# CI command
pytest --cov=app --cov-report=json --cov-report=term
```

The test suite:
- Runs fast (fixtures use temp directories)
- Is deterministic (no time-based flakiness)
- Provides clear failure messages
- Generates coverage reports

## Troubleshooting

### Tests Fail Locally

```bash
# Clean up any leftover test data
rm -rf /tmp/collectiveflow_test_*

# Reinstall dependencies
pip install -r requirements-test.txt

# Run with verbose output to see details
pytest -v
```

### Import Errors

Ensure you're in the correct directory:

```bash
cd projects/collectiveflow/web
pytest
```

### Coverage Not Generated

```bash
# Explicitly specify coverage options
pytest --cov=app --cov-report=html --cov-report=term-missing
```

### Specific Test Failing

Run just that test with maximum verbosity:

```bash
pytest tests/test_routes.py::TestHomeRoute::test_home_page_loads_successfully -vv
```

## Contributing Tests

When adding features to CollectiveFlow, add tests:

1. **Write tests first** (TDD approach) or alongside code
2. **Document your tests** with clear docstrings
3. **Use descriptive names** that explain what's being tested
4. **Add markers** so tests can be filtered
5. **Test edge cases** not just happy paths
6. **Keep tests simple** - complex tests are hard to maintain

### Test Review Checklist

Before submitting test code:

- [ ] Tests have clear, descriptive names
- [ ] Each test has a docstring explaining what and why
- [ ] Tests use appropriate fixtures
- [ ] Tests are properly marked (@pytest.mark.*)
- [ ] Edge cases are covered
- [ ] Tests pass individually and as a suite
- [ ] Code coverage hasn't decreased

## Learning Resources

### Understanding Pytest

- **Fixtures**: Shared test setup and data
- **Markers**: Tags for organizing and filtering tests
- **Parametrize**: Run same test with different inputs
- **Mocking**: Isolate code by replacing dependencies

### Understanding Flask Testing

```python
# Test client makes requests
client.get('/route')      # GET request
client.post('/route', data={...})  # POST with form data

# Response object
response.status_code      # HTTP status
response.data            # Raw bytes
response.content_type    # MIME type
```

### Understanding the Codebase

Read the application code alongside tests:

1. Start with `app.py` to understand routes
2. Look at templates to see what's rendered
3. Check `conftest.py` to understand available fixtures
4. Read tests to see how features are validated

## Questions?

Tests should be clear enough to answer most questions, but if you're unsure:

1. **Check test docstrings** - they explain the "why"
2. **Look at similar tests** - patterns are consistent
3. **Review `conftest.py`** - fixture documentation
4. **Run tests with `-v`** - see detailed output

Remember: These tests are teaching tools. If something isn't clear, that's a documentation issue we should fix!

---

**Built by the collective, for the collective, through consensus.**

Testing embodies our principles:
- Transparency (coverage shows what's tested)
- Accessibility (clear documentation)
- Horizontal organization (all tests equal)
- Knowledge sharing (tests teach how the system works)
