# CollectiveFlow Web Testing Guide

This document provides a quick-start guide to testing the CollectiveFlow Flask web application.

## Setup

### 1. Install Test Dependencies

```bash
cd projects/collectiveflow/web
pip install -r requirements-test.txt
```

This installs:
- `pytest` - Testing framework
- `pytest-cov` - Coverage reporting
- `pytest-flask` - Flask-specific testing utilities
- `pytest-mock` - Mocking support
- Additional testing utilities

### 2. Verify Installation

```bash
pytest --version
```

You should see pytest version information.

## Running Tests

### Basic Commands

```bash
# Run all tests
pytest

# Run with verbose output
pytest -v

# Run with coverage report
pytest --cov=app --cov-report=term-missing

# Run specific test file
pytest tests/test_routes.py

# Run specific test class
pytest tests/test_routes.py::TestHomeRoute

# Run specific test
pytest tests/test_routes.py::TestHomeRoute::test_home_page_loads_successfully
```

### Run by Category

```bash
# Route handler tests
pytest -m routes

# API endpoint tests
pytest -m api

# Jinja2 filter tests
pytest -m filters

# Data/YAML tests
pytest -m data

# Integration tests
pytest -m integration

# Skip slow tests
pytest -m "not slow"
```

### Coverage Reports

```bash
# Generate HTML coverage report
pytest --cov=app --cov-report=html

# View report
open htmlcov/index.html
```

## Test Structure

```
web/
├── pytest.ini                 # Pytest configuration
├── conftest.py               # Shared fixtures and test setup
├── requirements-test.txt     # Test dependencies
├── tests/
│   ├── __init__.py          # Test package marker
│   ├── README.md            # Comprehensive test documentation
│   ├── fixtures/            # Test data files
│   │   └── README.md        # Fixture documentation
│   ├── test_routes.py       # Route handler tests (529 lines)
│   ├── test_filters.py      # Jinja2 filter tests (490 lines)
│   ├── test_data.py         # YAML data tests (632 lines)
│   └── test_api.py          # API endpoint tests (583 lines)
```

**Total**: ~2,680 lines of test code with extensive documentation

## Test Coverage

The test suite covers:

### Route Handlers (`test_routes.py`)
- ✅ Home page (`/`)
- ✅ Proposals list (`/proposals`)
- ✅ Proposal detail (`/proposal/<id>`)
- ✅ About page (`/about`)
- ✅ Collective page (`/collective`)
- ✅ Proposal creation form (`/create`)
- ✅ Navigation and links
- ✅ Error handling (404s)
- ✅ Responsive design elements

### Jinja2 Filters (`test_filters.py`)
- ✅ `humanize_date` - Date formatting
- ✅ `status_emoji` - Status indicators
- ✅ `urgency_color` - Urgency styling
- ✅ Edge cases and error handling
- ✅ Filter integration in templates

### Data Operations (`test_data.py`)
- ✅ `load_proposals()` - Loading all proposals
- ✅ `get_proposal()` - Loading specific proposal
- ✅ `save_proposal()` - Creating new proposals
- ✅ YAML file operations
- ✅ Data integrity across save/load cycles
- ✅ Error handling (corrupted files, missing data)
- ✅ Special character handling

### API Endpoints (`test_api.py`)
- ✅ `/api/proposals` - List all proposals
- ✅ `/api/proposal/<id>` - Get specific proposal
- ✅ JSON response format
- ✅ CORS headers
- ✅ Error responses (404s)
- ✅ Data consistency with web interface
- ✅ Performance with large datasets

## Quick Test Examples

### Verify Basic Functionality

```bash
# Test that home page loads
pytest tests/test_routes.py::TestHomeRoute::test_home_page_loads_successfully -v

# Test that proposals load from YAML
pytest tests/test_data.py::TestLoadProposals::test_load_proposals_with_data -v

# Test that API returns JSON
pytest tests/test_api.py::TestProposalsListAPI::test_api_proposals_returns_json -v
```

### Test Coverage Report

```bash
# Run all tests with coverage
pytest --cov=app --cov-report=html

# Check coverage percentage
pytest --cov=app --cov-report=term
```

Expected coverage:
- Routes: High coverage (most routes tested)
- Filters: Complete coverage (all filters tested)
- Data operations: Complete coverage (all functions tested)
- Template rendering: Good coverage (via route tests)

## Common Issues

### Import Errors

If you see `ModuleNotFoundError: No module named 'app'`:

```bash
# Make sure you're in the web directory
cd projects/collectiveflow/web

# Run pytest from there
pytest
```

### Fixture Errors

If tests fail due to missing fixtures:

```bash
# Clean up any old test data
rm -rf /tmp/collectiveflow_test_*

# Re-run tests
pytest
```

### Permission Errors

If you see permission errors on test data:

```bash
# Clean up temp directories
rm -rf /tmp/collectiveflow_test_*

# Ensure you have write permissions
pytest tests/test_data.py -v
```

## Test Principles

Following our collective values, these tests are:

- **Transparent**: Extensive documentation and clear naming
- **Accessible**: Anyone can understand and modify tests
- **Non-hierarchical**: All tests are equal in importance
- **Educational**: Tests serve as documentation and teaching tools

Every test includes:
- Descriptive name explaining what it tests
- Docstring explaining why it matters
- Clear assertions with helpful messages
- Comments for non-obvious logic

## Next Steps

1. **Read the test documentation**: `tests/README.md` has comprehensive information
2. **Explore the tests**: Read through test files to understand coverage
3. **Run the tests**: Execute the test suite and review results
4. **Check coverage**: Generate HTML coverage report to see what's tested
5. **Add tests**: When adding features, add corresponding tests

## Documentation

- **`tests/README.md`**: Comprehensive test suite documentation
- **`tests/fixtures/README.md`**: Test fixture documentation
- **`conftest.py`**: Fixture definitions with detailed docstrings
- **Test files**: Each test has explanatory docstrings

## CI/CD Integration

The test suite is designed for continuous integration:

```bash
# Typical CI command
pytest --cov=app --cov-report=json --cov-report=term -v
```

Features:
- Fast execution (uses temporary directories)
- Deterministic (no flaky tests)
- Clear failure messages
- Coverage reporting in multiple formats

---

**Remember**: Tests embody our collective principles. They should be as accessible and transparent as the code they validate.

For detailed information, see **`tests/README.md`**.
