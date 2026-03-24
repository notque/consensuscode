# Test Fixtures

This directory contains sample data files used for testing the CollectiveFlow web application.

## Purpose

Test fixtures provide:
- **Consistent test data**: All tests use the same baseline data
- **Realistic scenarios**: Fixtures represent actual proposal structures
- **Isolation**: Tests don't depend on production data
- **Documentation**: Fixtures show expected data formats

## Using Fixtures in Tests

Fixtures are automatically created in temporary directories by the `conftest.py` configuration. You don't typically load files from this directory directly.

Instead, use pytest fixtures:

```python
def test_something(sample_proposals):
    # sample_proposals fixture creates test data automatically
    # in a temporary directory that's cleaned up after the test
    pass
```

## Manual Fixture Files

If you need to add permanent fixture files for specific test scenarios:

1. Create YAML files in this directory
2. Add a pytest fixture in `conftest.py` to load them
3. Document the fixture's purpose

## Sample Data Structure

Fixtures follow the standard proposal structure:

```yaml
id: proposal-2025-07-26-001
title: Example Proposal
description: Detailed description
proposer: agent-name
date: 2025-07-26T10:00:00-07:00
status: proposed
urgency: medium
affected_areas:
  - area1
  - area2
consensus_status: New proposal submitted
consensus_history:
  - timestamp: 2025-07-26T10:00:00-07:00
    event: proposal_created
    actor: agent-name
    details: Created with urgency: medium
consultations: []
```

## Fixture Principles

Following our collective values:
- **Transparent**: Fixtures are in human-readable YAML
- **Accessible**: Anyone can understand and modify fixtures
- **Documented**: Each fixture explains its purpose
- **Non-hierarchical**: No fixture is more important than others
