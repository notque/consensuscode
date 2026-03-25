"""
CollectiveFlow Web Application - Test Configuration and Fixtures

This module provides shared pytest fixtures for testing the CollectiveFlow web application.
Following our collective principles, these fixtures are:
- Well-documented for knowledge sharing
- Reusable across all test modules
- Transparent in their behavior
- Accessible to all test authors without special knowledge

Fixtures are pytest's way of providing test data and setup - they run before tests
and provide the resources tests need. Think of them as shared resources for the collective.
"""

import os
import tempfile
import shutil
from pathlib import Path
from datetime import datetime
import yaml
import pytest
from app import app as flask_app
from storage import YAMLStorage


# Test Data Fixtures - Sample Proposals
# These represent the types of data our application handles

SAMPLE_PROPOSAL_SIMPLE = {
    'id': 'test-proposal-001',
    'title': 'Test Proposal: Simple Example',
    'description': 'A simple test proposal for basic functionality testing',
    'proposer': 'test-agent',
    'date': '2025-07-26T10:00:00-07:00',
    'status': 'proposed',
    'urgency': 'medium',
    'affected_areas': ['testing'],
    'consensus_status': 'New proposal submitted',
    'consensus_history': [{
        'timestamp': '2025-07-26T10:00:00-07:00',
        'event': 'proposal_created',
        'actor': 'test-agent',
        'details': 'Created with urgency: medium'
    }],
    'consultations': []
}

SAMPLE_PROPOSAL_WITH_CONSULTATIONS = {
    'id': 'test-proposal-002',
    'title': 'Test Proposal: With Consultations',
    'description': 'A proposal with multiple consultations for testing consensus features',
    'proposer': 'test-agent-1',
    'date': '2025-07-26T11:00:00-07:00',
    'status': 'consultation',
    'urgency': 'high',
    'affected_areas': ['testing', 'development'],
    'consensus_status': 'Under active consultation',
    'consensus_history': [
        {
            'timestamp': '2025-07-26T11:00:00-07:00',
            'event': 'proposal_created',
            'actor': 'test-agent-1',
            'details': 'Created with urgency: high'
        },
        {
            'timestamp': '2025-07-26T11:05:00-07:00',
            'event': 'status_changed',
            'actor': 'test-agent-1',
            'details': 'Status changed from proposed to consultation'
        }
    ],
    'consultations': [
        {
            'contributor': 'test-agent-2',
            'timestamp': '2025-07-26T11:10:00-07:00',
            'input': 'I support this proposal with some suggestions.',
            'support': True
        },
        {
            'contributor': 'test-agent-3',
            'timestamp': '2025-07-26T11:15:00-07:00',
            'input': 'I have concerns about the timeline.',
            'concerns': ['Timeline seems aggressive'],
            'support': False
        }
    ]
}

SAMPLE_PROPOSAL_IMPLEMENTED = {
    'id': 'test-proposal-003',
    'title': 'Test Proposal: Implemented',
    'description': 'A completed proposal for testing implemented status',
    'proposer': 'test-agent-1',
    'date': '2025-07-25T10:00:00-07:00',
    'status': 'implemented',
    'urgency': 'low',
    'affected_areas': ['documentation'],
    'consensus_status': 'Implemented successfully',
    'decision': {
        'result': 'approved',
        'timestamp': '2025-07-26T10:00:00-07:00',
        'rationale': 'All agents reached consensus'
    },
    'consensus_history': [
        {
            'timestamp': '2025-07-25T10:00:00-07:00',
            'event': 'proposal_created',
            'actor': 'test-agent-1',
            'details': 'Created with urgency: low'
        },
        {
            'timestamp': '2025-07-26T10:00:00-07:00',
            'event': 'decision_recorded',
            'actor': 'collective',
            'details': 'Decision: approved'
        }
    ],
    'consultations': [
        {
            'contributor': 'test-agent-2',
            'timestamp': '2025-07-25T12:00:00-07:00',
            'input': 'Strong support for this documentation improvement.',
            'support': True
        }
    ]
}


@pytest.fixture
def app():
    """
    Fixture: Flask Application Instance

    Provides a Flask application configured for testing.
    This fixture ensures each test gets a fresh app instance with test configuration.

    Yields:
        Flask app: Application instance ready for testing
    """
    flask_app.config.update({
        'TESTING': True,
        'DEBUG': False,
        'SECRET_KEY': 'test-secret-key-for-testing-only',
    })

    yield flask_app


@pytest.fixture
def client(app, temp_data_dir):
    """
    Fixture: Flask Test Client

    Provides a test client for making HTTP requests to the application.
    This is how we simulate browser requests in our tests.

    The test client automatically uses the temp_data_dir fixture,
    ensuring each test has isolated data storage.

    Args:
        app: Flask app fixture
        temp_data_dir: Temporary data directory fixture

    Yields:
        FlaskClient: Test client for making requests
    """
    return app.test_client()


@pytest.fixture
def temp_data_dir(monkeypatch):
    """
    Fixture: Temporary Data Directory

    Creates a temporary directory for test data, ensuring tests don't interfere
    with real data or each other. This embodies our principle of isolation -
    each test runs in its own clean environment.

    After creating the directory, this fixture swaps the app's storage backend
    to a fresh YAMLStorage instance pointing at the temp directory. This
    ensures that all routes (which call storage.load_proposals(), etc.)
    operate on the test data, regardless of which backend is configured
    globally.

    The directory is automatically cleaned up after the test completes.

    Args:
        monkeypatch: Pytest's monkeypatch fixture for environment modification

    Yields:
        Path: Path to temporary proposals directory (inside the temp data dir)
    """
    # Create temporary directory
    temp_dir = tempfile.mkdtemp(prefix='collectiveflow_test_')
    proposals_dir = Path(temp_dir) / 'proposals'
    proposals_dir.mkdir(parents=True, exist_ok=True)

    # Configure app to use temp directory via storage backend
    monkeypatch.setenv('COLLECTIVEFLOW_DATA', temp_dir)

    import app as app_module
    from storage import YAMLStorage
    monkeypatch.setattr(app_module, 'storage', YAMLStorage(temp_dir))

    # Point the storage abstraction at the temp directory so that
    # load_proposals / get_proposal / save_proposal all use it.
    monkeypatch.setattr(app_module, 'storage', YAMLStorage(temp_dir))

    # Also patch PROPOSALS_DIR so that update_proposal() (which writes
    # directly to PROPOSALS_DIR rather than through the storage abstraction)
    # targets the temp directory too.
    monkeypatch.setattr(app_module, 'PROPOSALS_DIR', proposals_dir)

    yield proposals_dir

    # Cleanup after test
    shutil.rmtree(temp_dir, ignore_errors=True)


@pytest.fixture
def csrf_token(client):
    """
    Fixture: CSRF Token for POST/PUT/DELETE Requests

    The app's csrf_protect() before_request hook rejects any non-API
    POST/PUT/DELETE that lacks a valid _csrf_token matching the session.
    This fixture performs a GET to establish a session, then extracts
    the CSRF token so tests can include it in form data.

    Usage in tests:
        def test_something(self, client, csrf_token, temp_data_dir):
            form_data = {'title': 'X', 'description': 'Y', '_csrf_token': csrf_token}
            response = client.post('/create', data=form_data)

    Args:
        client: Flask test client fixture (with session support)

    Returns:
        str: A valid CSRF token matching the current session
    """
    # Make a GET request to establish a session and generate a CSRF token
    with client.session_transaction() as sess:
        import secrets
        token = secrets.token_hex(32)
        sess['_csrf_token'] = token
    return token


@pytest.fixture
def sample_proposals(temp_data_dir):
    """
    Fixture: Sample Proposal Data

    Creates a set of sample proposals in the test data directory.
    These proposals cover different statuses and scenarios, providing
    realistic test data for our application.

    This fixture demonstrates transparency - test data is clearly defined
    and available for all tests to use.

    Args:
        temp_data_dir: Temporary data directory fixture

    Returns:
        list: List of proposal dictionaries that were created
    """
    proposals = [
        SAMPLE_PROPOSAL_SIMPLE,
        SAMPLE_PROPOSAL_WITH_CONSULTATIONS,
        SAMPLE_PROPOSAL_IMPLEMENTED
    ]

    # Write each proposal to a YAML file
    for proposal in proposals:
        yaml_path = temp_data_dir / f"{proposal['id']}.yaml"
        with open(yaml_path, 'w') as f:
            yaml.safe_dump(proposal, f, default_flow_style=False, sort_keys=False)

    return proposals


@pytest.fixture
def empty_proposals_dir(temp_data_dir):
    """
    Fixture: Empty Proposals Directory

    Provides a clean proposals directory with no data.
    Useful for testing application behavior with no proposals.

    Args:
        temp_data_dir: Temporary data directory fixture

    Returns:
        Path: Path to empty proposals directory
    """
    # temp_data_dir already creates an empty directory
    return temp_data_dir


@pytest.fixture
def proposal_form_data():
    """
    Fixture: Proposal Form Data

    Provides sample form data for testing proposal creation.
    This represents what a user would submit through the web form.

    Returns:
        dict: Form data dictionary
    """
    return {
        'title': 'Test Proposal from Form',
        'description': 'This is a test proposal submitted via the web form.',
        'proposer': 'web-user',
        'urgency': 'medium',
        'affected_areas': ['testing', 'web-interface']
    }


@pytest.fixture
def mock_datetime(monkeypatch):
    """
    Fixture: Mocked Datetime

    Provides a fixed datetime for testing time-dependent functionality.
    This ensures tests are deterministic and don't fail due to timing issues.

    Args:
        monkeypatch: Pytest's monkeypatch fixture

    Returns:
        datetime: Fixed datetime object
    """
    class MockDatetime:
        @staticmethod
        def now():
            return datetime(2025, 7, 26, 12, 0, 0)

        @staticmethod
        def isoformat():
            return '2025-07-26T12:00:00'

    return MockDatetime


# Helper function for tests (not a fixture, but useful for all tests)
def assert_proposal_structure(proposal_data):
    """
    Helper: Validate Proposal Structure

    Checks that a proposal has all required fields.
    This helper embodies our commitment to clear expectations -
    we document what a valid proposal looks like.

    Args:
        proposal_data (dict): Proposal data to validate

    Raises:
        AssertionError: If proposal structure is invalid
    """
    required_fields = [
        'id', 'title', 'description', 'proposer',
        'date', 'status', 'urgency'
    ]

    for field in required_fields:
        assert field in proposal_data, f"Missing required field: {field}"

    # Validate status is one of known values
    valid_statuses = ['proposed', 'consultation', 'consensus', 'implemented', 'blocked', 'withdrawn']
    assert proposal_data['status'] in valid_statuses, f"Invalid status: {proposal_data['status']}"

    # Validate urgency is one of known values
    valid_urgencies = ['low', 'medium', 'high', 'emergency']
    assert proposal_data['urgency'] in valid_urgencies, f"Invalid urgency: {proposal_data['urgency']}"


# Export sample data for direct use in tests
__all__ = [
    'app',
    'client',
    'temp_data_dir',
    'sample_proposals',
    'empty_proposals_dir',
    'proposal_form_data',
    'mock_datetime',
    'assert_proposal_structure',
    'SAMPLE_PROPOSAL_SIMPLE',
    'SAMPLE_PROPOSAL_WITH_CONSULTATIONS',
    'SAMPLE_PROPOSAL_IMPLEMENTED',
]
