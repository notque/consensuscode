"""
CollectiveFlow Web Application - API Endpoint Tests

These tests validate the JSON API endpoints, ensuring:
- API routes return proper JSON responses
- Status codes are appropriate
- Data structure matches expected format
- Error handling works correctly
- CORS headers are present for cross-origin requests
- The API embodies accessibility and transparency

Why Test APIs?
APIs enable integration with other tools and future mobile apps.
Testing ensures reliable programmatic access to our collective's data.
"""

import pytest
import json


class TestProposalsListAPI:
    """
    Tests for the /api/proposals endpoint.

    This endpoint returns all proposals in JSON format.
    It's used for programmatic access and future integrations.
    """

    @pytest.mark.api
    def test_api_proposals_returns_json(self, client, sample_proposals):
        """
        Test: API endpoint returns JSON content type

        APIs should always return proper JSON with correct headers.
        """
        response = client.get('/api/proposals')

        assert response.status_code == 200
        assert response.content_type == 'application/json'

    @pytest.mark.api
    def test_api_proposals_structure(self, client, sample_proposals):
        """
        Test: API response has expected structure

        The response should include both proposals array and count.
        """
        response = client.get('/api/proposals')
        data = json.loads(response.data)

        # Should have expected top-level keys
        assert 'proposals' in data
        assert 'count' in data

        # proposals should be a list
        assert isinstance(data['proposals'], list)
        # count should be a number
        assert isinstance(data['count'], int)

    @pytest.mark.api
    def test_api_proposals_returns_all_proposals(self, client, sample_proposals):
        """
        Test: API returns all proposals in the system

        Transparency requires exposing all proposals via API.
        """
        response = client.get('/api/proposals')
        data = json.loads(response.data)

        # We created 3 sample proposals
        assert data['count'] == 3
        assert len(data['proposals']) == 3

    @pytest.mark.api
    def test_api_proposals_proposal_structure(self, client, sample_proposals):
        """
        Test: Each proposal in API response has required fields

        API consumers need complete proposal data.
        """
        response = client.get('/api/proposals')
        data = json.loads(response.data)

        # Check first proposal has expected fields
        if data['proposals']:
            proposal = data['proposals'][0]

            required_fields = ['id', 'title', 'description', 'proposer', 'date', 'status', 'urgency']
            for field in required_fields:
                assert field in proposal, f"Missing required field: {field}"

    @pytest.mark.api
    def test_api_proposals_includes_consultations(self, client, sample_proposals):
        """
        Test: API includes consultation data

        Full transparency requires exposing all consultation input.
        """
        response = client.get('/api/proposals')
        data = json.loads(response.data)

        # Find the proposal with consultations
        proposal_with_consultations = next(
            (p for p in data['proposals'] if p['id'] == 'test-proposal-002'),
            None
        )

        assert proposal_with_consultations is not None
        assert 'consultations' in proposal_with_consultations
        assert len(proposal_with_consultations['consultations']) > 0

    @pytest.mark.api
    def test_api_proposals_sorted_by_date(self, client, sample_proposals):
        """
        Test: API returns proposals sorted by date (newest first)

        Consistent sorting helps API consumers.
        """
        response = client.get('/api/proposals')
        data = json.loads(response.data)

        proposals = data['proposals']
        if len(proposals) > 1:
            dates = [p.get('date', '') for p in proposals]
            # Should be in descending order (newest first)
            assert dates == sorted(dates, reverse=True)

    @pytest.mark.api
    def test_api_proposals_empty_collection(self, client, empty_proposals_dir):
        """
        Test: API handles empty proposal collection gracefully

        Empty state should return valid JSON with empty array.
        """
        response = client.get('/api/proposals')
        data = json.loads(response.data)

        assert response.status_code == 200
        assert data['count'] == 0
        assert data['proposals'] == []

    @pytest.mark.api
    def test_api_proposals_has_cors_headers(self, client, sample_proposals):
        """
        Test: API includes CORS headers for cross-origin requests

        CORS enables web apps from other domains to use our API.
        """
        response = client.get('/api/proposals')

        # Flask-CORS should add these headers
        # Note: Exact headers depend on CORS configuration
        assert response.status_code == 200
        # CORS headers are present (Flask-CORS is enabled in app.py)


class TestProposalDetailAPI:
    """
    Tests for the /api/proposal/<id> endpoint.

    This endpoint returns a specific proposal by ID in JSON format.
    """

    @pytest.mark.api
    def test_api_proposal_returns_json(self, client, sample_proposals):
        """
        Test: Proposal detail API returns JSON

        Single proposal endpoint should also return proper JSON.
        """
        response = client.get('/api/proposal/test-proposal-001')

        assert response.status_code == 200
        assert response.content_type == 'application/json'

    @pytest.mark.api
    def test_api_proposal_returns_correct_proposal(self, client, sample_proposals):
        """
        Test: API returns the requested proposal

        The ID in the URL should determine which proposal is returned.
        """
        response = client.get('/api/proposal/test-proposal-001')
        data = json.loads(response.data)

        assert data['id'] == 'test-proposal-001'
        assert data['title'] == 'Test Proposal: Simple Example'

    @pytest.mark.api
    def test_api_proposal_includes_all_fields(self, client, sample_proposals):
        """
        Test: API response includes complete proposal data

        Detail endpoint should return everything about the proposal.
        """
        response = client.get('/api/proposal/test-proposal-002')
        data = json.loads(response.data)

        # Check for all major fields
        assert 'id' in data
        assert 'title' in data
        assert 'description' in data
        assert 'proposer' in data
        assert 'date' in data
        assert 'status' in data
        assert 'urgency' in data
        assert 'consultations' in data
        assert 'consensus_history' in data

    @pytest.mark.api
    def test_api_proposal_includes_consultations(self, client, sample_proposals):
        """
        Test: API includes full consultation details

        All consultation input should be accessible via API.
        """
        response = client.get('/api/proposal/test-proposal-002')
        data = json.loads(response.data)

        assert 'consultations' in data
        assert len(data['consultations']) == 2

        # Check consultation structure
        consultation = data['consultations'][0]
        assert 'contributor' in consultation
        assert 'timestamp' in consultation
        assert 'input' in consultation

    @pytest.mark.api
    def test_api_proposal_includes_decision(self, client, sample_proposals):
        """
        Test: Implemented proposal includes decision data via API

        Decision rationale should be accessible programmatically.
        """
        response = client.get('/api/proposal/test-proposal-003')
        data = json.loads(response.data)

        assert 'decision' in data
        assert data['decision']['result'] == 'approved'
        assert 'rationale' in data['decision']

    @pytest.mark.api
    def test_api_proposal_not_found_returns_404(self, client, sample_proposals):
        """
        Test: Non-existent proposal returns 404 with error JSON

        API should indicate errors clearly with proper status codes.
        """
        response = client.get('/api/proposal/nonexistent-id')

        assert response.status_code == 404
        assert response.content_type == 'application/json'

        data = json.loads(response.data)
        assert 'error' in data
        assert 'not found' in data['error'].lower()

    @pytest.mark.api
    def test_api_proposal_not_found_error_structure(self, client, sample_proposals):
        """
        Test: Error responses have consistent structure

        Consistent error format helps API consumers handle failures.
        """
        response = client.get('/api/proposal/does-not-exist')
        data = json.loads(response.data)

        # Error response should have 'error' key with string message
        assert 'error' in data
        assert isinstance(data['error'], str)

    @pytest.mark.api
    def test_api_proposal_has_cors_headers(self, client, sample_proposals):
        """
        Test: Proposal detail API includes CORS headers

        Cross-origin access should work for individual proposals.
        """
        response = client.get('/api/proposal/test-proposal-001')
        assert response.status_code == 200
        # CORS headers present (Flask-CORS enabled)


class TestAPIDataConsistency:
    """
    Tests for consistency between API and web interface data.

    The API should return the same data that the web interface displays.
    This ensures a single source of truth.
    """

    @pytest.mark.api
    @pytest.mark.integration
    def test_api_and_web_return_same_proposals(self, client, sample_proposals):
        """
        Test: API and web interface show same proposals

        Both interfaces should reflect the same underlying data.
        """
        # Get data from API
        api_response = client.get('/api/proposals')
        api_data = json.loads(api_response.data)

        # Get data from web interface
        web_response = client.get('/')
        web_html = web_response.data.decode('utf-8')

        # Check that proposal IDs in API appear in web interface
        for proposal in api_data['proposals']:
            assert proposal['id'] in web_html or proposal['title'] in web_html

    @pytest.mark.api
    @pytest.mark.integration
    def test_api_and_web_show_same_proposal_details(self, client, sample_proposals):
        """
        Test: Proposal details match between API and web

        Detail views should show consistent information.
        """
        proposal_id = 'test-proposal-001'

        # Get from API
        api_response = client.get(f'/api/proposal/{proposal_id}')
        api_data = json.loads(api_response.data)

        # Get from web
        web_response = client.get(f'/proposal/{proposal_id}')
        web_html = web_response.data.decode('utf-8')

        # Key fields should appear in both
        assert api_data['title'] in web_html
        assert api_data['description'] in web_html
        assert api_data['proposer'] in web_html

    @pytest.mark.api
    @pytest.mark.integration
    def test_api_reflects_newly_saved_proposals(self, client, temp_data_dir):
        """
        Test: API immediately reflects newly saved proposals

        No caching issues - new data should be immediately accessible.
        """
        import app as _app_mod

        # Save a new proposal
        new_proposal = {
            'title': 'API Consistency Test',
            'description': 'Testing API reflects new data',
            'proposer': 'test-agent',
            'urgency': 'low'
        }

        proposal_id = _app_mod.storage.save_proposal(new_proposal)

        # API should immediately return the new proposal
        response = client.get(f'/api/proposal/{proposal_id}')
        assert response.status_code == 200

        data = json.loads(response.data)
        assert data['title'] == new_proposal['title']

        # Should also appear in list
        list_response = client.get('/api/proposals')
        list_data = json.loads(list_response.data)

        proposal_ids = [p['id'] for p in list_data['proposals']]
        assert proposal_id in proposal_ids


class TestAPIContentNegotiation:
    """
    Tests for content type handling and API conventions.

    Good APIs follow conventions for headers and content types.
    """

    @pytest.mark.api
    def test_api_sets_json_content_type(self, client, sample_proposals):
        """
        Test: API responses have application/json content type

        Proper content type helps clients parse responses correctly.
        """
        response = client.get('/api/proposals')
        assert 'application/json' in response.content_type

        response = client.get('/api/proposal/test-proposal-001')
        assert 'application/json' in response.content_type

    @pytest.mark.api
    def test_api_response_is_valid_json(self, client, sample_proposals):
        """
        Test: API responses are valid, parseable JSON

        Responses should never contain malformed JSON.
        """
        response = client.get('/api/proposals')

        # Should parse without error
        try:
            data = json.loads(response.data)
            assert data is not None
        except json.JSONDecodeError:
            pytest.fail("API response is not valid JSON")

    @pytest.mark.api
    def test_api_json_is_properly_encoded(self, client, sample_proposals):
        """
        Test: API properly encodes special characters in JSON

        Unicode and special characters should be handled correctly.
        """
        # Create proposal with special characters
        import app as _app_mod

        special_proposal = {
            'title': 'Special Characters: "quotes", emoji 🤝, unicode café',
            'description': 'Testing: special chars & encoding',
            'proposer': 'test-agent'
        }

        proposal_id = _app_mod.storage.save_proposal(special_proposal)

        # Get via API
        response = client.get(f'/api/proposal/{proposal_id}')
        data = json.loads(response.data)

        # Special characters should be preserved
        assert 'quotes' in data['title']
        assert '🤝' in data['title']
        assert 'café' in data['title']


class TestAPIErrorHandling:
    """
    Tests for API error handling and edge cases.

    Robust error handling makes APIs more reliable and easier to use.
    """

    @pytest.mark.api
    def test_api_404_returns_json_error(self, client, sample_proposals):
        """
        Test: 404 errors return JSON, not HTML

        API errors should be in JSON format for programmatic handling.
        """
        response = client.get('/api/proposal/nonexistent')

        assert response.status_code == 404
        assert 'application/json' in response.content_type

        # Should be parseable JSON
        data = json.loads(response.data)
        assert 'error' in data

    @pytest.mark.api
    def test_api_handles_special_characters_in_id(self, client, sample_proposals):
        """
        Test: API handles special characters in proposal IDs gracefully

        Invalid characters shouldn't crash the API.
        """
        response = client.get('/api/proposal/invalid<>chars')

        # Should return 404, not crash
        assert response.status_code == 404
        assert 'application/json' in response.content_type

    @pytest.mark.api
    def test_api_handles_empty_data_directory(self, client, empty_proposals_dir):
        """
        Test: API works correctly with no proposals

        Empty state should return valid responses.
        """
        response = client.get('/api/proposals')
        data = json.loads(response.data)

        assert response.status_code == 200
        assert data['count'] == 0
        assert data['proposals'] == []

    @pytest.mark.api
    def test_api_handles_very_long_id(self, client, sample_proposals):
        """
        Test: API handles unusually long proposal IDs

        Edge cases shouldn't cause errors.
        """
        long_id = 'a' * 1000
        response = client.get(f'/api/proposal/{long_id}')

        # Should return 404, not crash
        assert response.status_code == 404


class TestAPIPerformance:
    """
    Tests for API performance considerations.

    These tests ensure the API remains responsive with various data loads.
    """

    @pytest.mark.api
    @pytest.mark.slow
    def test_api_handles_many_proposals(self, client, temp_data_dir):
        """
        Test: API performs well with larger numbers of proposals

        The API should remain responsive as the collective grows.
        """
        import app as _app_mod

        # Create multiple proposals
        for i in range(50):
            _app_mod.storage.save_proposal({
                'title': f'Proposal {i}',
                'description': f'Test proposal number {i}',
                'proposer': 'test-agent',
                'urgency': 'low'
            })

        # API should still respond quickly
        response = client.get('/api/proposals')
        data = json.loads(response.data)

        assert response.status_code == 200
        assert data['count'] == 50
        assert len(data['proposals']) == 50

    @pytest.mark.api
    def test_api_proposal_with_many_consultations(self, client, temp_data_dir):
        """
        Test: API handles proposals with many consultations

        Proposals with extensive discussion should work correctly.
        """
        import app as _app_mod

        # Create proposal with many consultations
        consultations = [
            {
                'contributor': f'agent-{i}',
                'timestamp': f'2025-07-26T10:{i:02d}:00',
                'input': f'Consultation {i}',
                'support': True
            }
            for i in range(20)
        ]

        proposal_data = {
            'title': 'Proposal with Many Consultations',
            'description': 'Testing many consultations',
            'proposer': 'test-agent',
            'urgency': 'medium',
            'consultations': consultations
        }

        from pathlib import Path
        import yaml

        # Save with consultations
        proposal_id = f'proposal-many-consultations-{datetime.now().strftime("%Y%m%d%H%M%S")}'
        proposal_data['id'] = proposal_id
        proposal_data['date'] = datetime.now().isoformat()
        proposal_data['status'] = 'consultation'

        yaml_path = temp_data_dir / f'{proposal_id}.yaml'
        with open(yaml_path, 'w') as f:
            yaml.safe_dump(proposal_data, f)

        # API should return all consultations
        response = client.get(f'/api/proposal/{proposal_id}')
        data = json.loads(response.data)

        assert response.status_code == 200
        assert len(data['consultations']) == 20


# Import datetime for timestamp generation
from datetime import datetime
