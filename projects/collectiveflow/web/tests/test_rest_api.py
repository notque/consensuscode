"""
CollectiveFlow Web Application - REST API Endpoint Tests

Tests for the new REST API endpoints added to support programmatic interaction:
- POST /api/proposals          — create proposals via API
- POST /api/proposals/<id>/consultation — add consultation input
- PUT  /api/proposals/<id>/status       — advance proposal status
- GET  /api/collective/stats            — collective statistics
- GET  /api/proposals?status=...        — filtered listing
- GET  /api/proposals/<id>              — canonical plural route

Also validates:
- Consistent error response format {"error": "...", "code": "..."}
- Input validation for all write endpoints
- Status transition enforcement
"""

import pytest
import json


# ---------------------------------------------------------------------------
# POST /api/proposals — create proposals
# ---------------------------------------------------------------------------

class TestCreateProposalAPI:
    """Tests for creating proposals through the JSON API."""

    @pytest.mark.api
    def test_create_proposal_success(self, client, temp_data_dir):
        """A valid POST creates a proposal and returns 201."""
        response = client.post('/api/proposals',
            data=json.dumps({
                'title': 'API-created proposal',
                'description': 'Created through the REST API',
                'proposer': 'api-test-agent',
                'urgency': 'medium',
                'affected_areas': ['api', 'testing'],
            }),
            content_type='application/json')

        assert response.status_code == 201
        data = json.loads(response.data)

        assert data['title'] == 'API-created proposal'
        assert data['description'] == 'Created through the REST API'
        assert data['proposer'] == 'api-test-agent'
        assert data['urgency'] == 'medium'
        assert data['status'] == 'proposed'
        assert 'id' in data
        assert 'date' in data
        assert data['consultations'] == []

    @pytest.mark.api
    def test_create_proposal_returns_location_header(self, client, temp_data_dir):
        """201 response includes a Location header pointing to the new resource."""
        response = client.post('/api/proposals',
            data=json.dumps({
                'title': 'Location header test',
                'description': 'Check the Location header',
            }),
            content_type='application/json')

        assert response.status_code == 201
        assert 'Location' in response.headers
        assert '/api/proposals/' in response.headers['Location']

    @pytest.mark.api
    def test_create_proposal_defaults(self, client, temp_data_dir):
        """Omitted optional fields get sensible defaults."""
        response = client.post('/api/proposals',
            data=json.dumps({
                'title': 'Defaults test',
                'description': 'Only required fields',
            }),
            content_type='application/json')

        assert response.status_code == 201
        data = json.loads(response.data)
        assert data['proposer'] == 'api-user'
        assert data['urgency'] == 'medium'
        assert data['affected_areas'] == []

    @pytest.mark.api
    def test_create_proposal_missing_title(self, client, temp_data_dir):
        """Missing title returns 422 with VALIDATION_ERROR."""
        response = client.post('/api/proposals',
            data=json.dumps({'description': 'No title'}),
            content_type='application/json')

        assert response.status_code == 422
        data = json.loads(response.data)
        assert data['code'] == 'VALIDATION_ERROR'
        assert 'title' in data['error'].lower()

    @pytest.mark.api
    def test_create_proposal_missing_description(self, client, temp_data_dir):
        """Missing description returns 422 with VALIDATION_ERROR."""
        response = client.post('/api/proposals',
            data=json.dumps({'title': 'No description'}),
            content_type='application/json')

        assert response.status_code == 422
        data = json.loads(response.data)
        assert data['code'] == 'VALIDATION_ERROR'
        assert 'description' in data['error'].lower()

    @pytest.mark.api
    def test_create_proposal_invalid_urgency(self, client, temp_data_dir):
        """Invalid urgency value returns 422."""
        response = client.post('/api/proposals',
            data=json.dumps({
                'title': 'Bad urgency',
                'description': 'Testing invalid urgency',
                'urgency': 'critical',
            }),
            content_type='application/json')

        assert response.status_code == 422
        data = json.loads(response.data)
        assert data['code'] == 'VALIDATION_ERROR'
        assert 'urgency' in data['error'].lower()

    @pytest.mark.api
    def test_create_proposal_affected_areas_not_list(self, client, temp_data_dir):
        """Non-list affected_areas returns 422."""
        response = client.post('/api/proposals',
            data=json.dumps({
                'title': 'Bad areas',
                'description': 'Testing non-list areas',
                'affected_areas': 'not-a-list',
            }),
            content_type='application/json')

        assert response.status_code == 422
        data = json.loads(response.data)
        assert data['code'] == 'VALIDATION_ERROR'

    @pytest.mark.api
    def test_create_proposal_wrong_content_type(self, client, temp_data_dir):
        """Non-JSON content type returns 415."""
        response = client.post('/api/proposals',
            data='title=oops',
            content_type='application/x-www-form-urlencoded')

        assert response.status_code == 415
        data = json.loads(response.data)
        assert data['code'] == 'INVALID_CONTENT_TYPE'

    @pytest.mark.api
    def test_create_proposal_is_retrievable(self, client, temp_data_dir):
        """A created proposal can be fetched back via GET."""
        create_resp = client.post('/api/proposals',
            data=json.dumps({
                'title': 'Roundtrip test',
                'description': 'Create then GET',
            }),
            content_type='application/json')

        proposal_id = json.loads(create_resp.data)['id']

        get_resp = client.get(f'/api/proposals/{proposal_id}')
        assert get_resp.status_code == 200
        data = json.loads(get_resp.data)
        assert data['title'] == 'Roundtrip test'


# ---------------------------------------------------------------------------
# POST /api/proposals/<id>/consultation — add consultation input
# ---------------------------------------------------------------------------

class TestConsultationAPI:
    """Tests for adding consultation input to proposals."""

    @pytest.mark.api
    def test_add_consultation_success(self, client, sample_proposals):
        """Valid consultation input returns 201."""
        response = client.post('/api/proposals/test-proposal-001/consultation',
            data=json.dumps({
                'contributor': 'api-agent',
                'input': 'I support this proposal fully.',
                'support': True,
            }),
            content_type='application/json')

        assert response.status_code == 201
        data = json.loads(response.data)
        assert data['message'] == 'Consultation added'
        assert data['consultation']['contributor'] == 'api-agent'
        assert data['consultation']['support'] is True
        assert data['proposal_id'] == 'test-proposal-001'

    @pytest.mark.api
    def test_add_consultation_with_concerns(self, client, sample_proposals):
        """Consultation with concerns list is stored properly."""
        response = client.post('/api/proposals/test-proposal-001/consultation',
            data=json.dumps({
                'contributor': 'concerned-agent',
                'input': 'I have reservations about the timeline.',
                'support': False,
                'concerns': ['Timeline too aggressive', 'Needs more testing'],
            }),
            content_type='application/json')

        assert response.status_code == 201
        data = json.loads(response.data)
        assert data['consultation']['support'] is False
        assert len(data['consultation']['concerns']) == 2

    @pytest.mark.api
    def test_add_consultation_persists(self, client, sample_proposals):
        """Added consultation appears when proposal is fetched again."""
        client.post('/api/proposals/test-proposal-001/consultation',
            data=json.dumps({
                'contributor': 'persistence-agent',
                'input': 'Testing persistence',
            }),
            content_type='application/json')

        get_resp = client.get('/api/proposals/test-proposal-001')
        data = json.loads(get_resp.data)

        contributors = [c['contributor'] for c in data['consultations']]
        assert 'persistence-agent' in contributors

    @pytest.mark.api
    def test_add_consultation_updates_history(self, client, sample_proposals):
        """Adding consultation appends to consensus_history."""
        client.post('/api/proposals/test-proposal-001/consultation',
            data=json.dumps({
                'contributor': 'history-agent',
                'input': 'Testing history update',
            }),
            content_type='application/json')

        get_resp = client.get('/api/proposals/test-proposal-001')
        data = json.loads(get_resp.data)

        events = [h['event'] for h in data.get('consensus_history', [])]
        assert 'consultation_added' in events

    @pytest.mark.api
    def test_add_consultation_missing_contributor(self, client, sample_proposals):
        """Missing contributor returns 422."""
        response = client.post('/api/proposals/test-proposal-001/consultation',
            data=json.dumps({'input': 'No contributor'}),
            content_type='application/json')

        assert response.status_code == 422
        data = json.loads(response.data)
        assert data['code'] == 'VALIDATION_ERROR'
        assert 'contributor' in data['error'].lower()

    @pytest.mark.api
    def test_add_consultation_missing_input(self, client, sample_proposals):
        """Missing input text returns 422."""
        response = client.post('/api/proposals/test-proposal-001/consultation',
            data=json.dumps({'contributor': 'agent'}),
            content_type='application/json')

        assert response.status_code == 422
        data = json.loads(response.data)
        assert data['code'] == 'VALIDATION_ERROR'
        assert 'input' in data['error'].lower()

    @pytest.mark.api
    def test_add_consultation_to_nonexistent_proposal(self, client, sample_proposals):
        """Consultation on nonexistent proposal returns 404."""
        response = client.post('/api/proposals/nonexistent/consultation',
            data=json.dumps({
                'contributor': 'agent',
                'input': 'hello',
            }),
            content_type='application/json')

        assert response.status_code == 404
        data = json.loads(response.data)
        assert data['code'] == 'NOT_FOUND'

    @pytest.mark.api
    def test_add_consultation_to_implemented_proposal(self, client, sample_proposals):
        """Consultation on an implemented proposal returns 409."""
        response = client.post('/api/proposals/test-proposal-003/consultation',
            data=json.dumps({
                'contributor': 'late-agent',
                'input': 'Too late',
            }),
            content_type='application/json')

        assert response.status_code == 409
        data = json.loads(response.data)
        assert data['code'] == 'INVALID_STATE'

    @pytest.mark.api
    def test_add_consultation_wrong_content_type(self, client, sample_proposals):
        """Non-JSON content type returns 415."""
        response = client.post('/api/proposals/test-proposal-001/consultation',
            data='contributor=x&input=y',
            content_type='application/x-www-form-urlencoded')

        assert response.status_code == 415


# ---------------------------------------------------------------------------
# PUT /api/proposals/<id>/status — status transitions
# ---------------------------------------------------------------------------

class TestStatusTransitionAPI:
    """Tests for updating proposal status with transition enforcement."""

    @pytest.mark.api
    def test_valid_transition_proposed_to_consultation(self, client, sample_proposals):
        """proposed -> consultation is a valid transition."""
        response = client.put('/api/proposals/test-proposal-001/status',
            data=json.dumps({
                'status': 'consultation',
                'actor': 'coordinator-agent',
                'reason': 'Ready for collective input',
            }),
            content_type='application/json')

        assert response.status_code == 200
        data = json.loads(response.data)
        assert data['previous_status'] == 'proposed'
        assert data['new_status'] == 'consultation'

    @pytest.mark.api
    def test_status_change_persists(self, client, sample_proposals):
        """Status change is reflected on subsequent GET."""
        client.put('/api/proposals/test-proposal-001/status',
            data=json.dumps({
                'status': 'consultation',
                'actor': 'agent',
            }),
            content_type='application/json')

        get_resp = client.get('/api/proposals/test-proposal-001')
        data = json.loads(get_resp.data)
        assert data['status'] == 'consultation'

    @pytest.mark.api
    def test_status_change_updates_history(self, client, sample_proposals):
        """Status change appends to consensus_history."""
        client.put('/api/proposals/test-proposal-001/status',
            data=json.dumps({
                'status': 'consultation',
                'actor': 'history-checker',
                'reason': 'Testing history',
            }),
            content_type='application/json')

        get_resp = client.get('/api/proposals/test-proposal-001')
        data = json.loads(get_resp.data)

        last_event = data['consensus_history'][-1]
        assert last_event['event'] == 'status_changed'
        assert last_event['actor'] == 'history-checker'
        assert 'Testing history' in last_event['details']

    @pytest.mark.api
    def test_invalid_transition_proposed_to_implemented(self, client, sample_proposals):
        """proposed -> implemented is not allowed (must go through consultation/consensus)."""
        response = client.put('/api/proposals/test-proposal-001/status',
            data=json.dumps({
                'status': 'implemented',
                'actor': 'impatient-agent',
            }),
            content_type='application/json')

        assert response.status_code == 409
        data = json.loads(response.data)
        assert data['code'] == 'INVALID_TRANSITION'

    @pytest.mark.api
    def test_invalid_transition_from_terminal_status(self, client, sample_proposals):
        """implemented is a terminal status — no transitions allowed."""
        response = client.put('/api/proposals/test-proposal-003/status',
            data=json.dumps({
                'status': 'proposed',
                'actor': 'agent',
            }),
            content_type='application/json')

        assert response.status_code == 409
        data = json.loads(response.data)
        assert data['code'] == 'INVALID_TRANSITION'
        assert 'terminal' in data['error'].lower() or 'none' in data['error'].lower()

    @pytest.mark.api
    def test_invalid_status_value(self, client, sample_proposals):
        """Unknown status value returns 422."""
        response = client.put('/api/proposals/test-proposal-001/status',
            data=json.dumps({
                'status': 'approved',
                'actor': 'agent',
            }),
            content_type='application/json')

        assert response.status_code == 422
        data = json.loads(response.data)
        assert data['code'] == 'VALIDATION_ERROR'

    @pytest.mark.api
    def test_missing_status_field(self, client, sample_proposals):
        """Missing status field returns 422."""
        response = client.put('/api/proposals/test-proposal-001/status',
            data=json.dumps({'actor': 'agent'}),
            content_type='application/json')

        assert response.status_code == 422
        data = json.loads(response.data)
        assert data['code'] == 'VALIDATION_ERROR'

    @pytest.mark.api
    def test_missing_actor_field(self, client, sample_proposals):
        """Missing actor field returns 422."""
        response = client.put('/api/proposals/test-proposal-001/status',
            data=json.dumps({'status': 'consultation'}),
            content_type='application/json')

        assert response.status_code == 422
        data = json.loads(response.data)
        assert data['code'] == 'VALIDATION_ERROR'
        assert 'actor' in data['error'].lower()

    @pytest.mark.api
    def test_status_nonexistent_proposal(self, client, sample_proposals):
        """Status change on nonexistent proposal returns 404."""
        response = client.put('/api/proposals/nonexistent/status',
            data=json.dumps({
                'status': 'consultation',
                'actor': 'agent',
            }),
            content_type='application/json')

        assert response.status_code == 404
        data = json.loads(response.data)
        assert data['code'] == 'NOT_FOUND'

    @pytest.mark.api
    def test_status_wrong_content_type(self, client, sample_proposals):
        """Non-JSON content type returns 415."""
        response = client.put('/api/proposals/test-proposal-001/status',
            data='status=consultation&actor=agent',
            content_type='application/x-www-form-urlencoded')

        assert response.status_code == 415

    @pytest.mark.api
    def test_valid_transition_proposed_to_withdrawn(self, client, sample_proposals):
        """proposed -> withdrawn is valid (proposer can withdraw)."""
        response = client.put('/api/proposals/test-proposal-001/status',
            data=json.dumps({
                'status': 'withdrawn',
                'actor': 'test-agent',
                'reason': 'No longer needed',
            }),
            content_type='application/json')

        assert response.status_code == 200
        data = json.loads(response.data)
        assert data['new_status'] == 'withdrawn'

    @pytest.mark.api
    def test_valid_transition_consultation_to_consensus(self, client, sample_proposals):
        """consultation -> consensus is valid."""
        response = client.put('/api/proposals/test-proposal-002/status',
            data=json.dumps({
                'status': 'consensus',
                'actor': 'coordinator',
            }),
            content_type='application/json')

        assert response.status_code == 200
        data = json.loads(response.data)
        assert data['new_status'] == 'consensus'


# ---------------------------------------------------------------------------
# GET /api/collective/stats — collective statistics
# ---------------------------------------------------------------------------

class TestCollectiveStatsAPI:
    """Tests for the collective statistics endpoint."""

    @pytest.mark.api
    def test_stats_returns_200(self, client, sample_proposals):
        """Stats endpoint returns 200 with JSON."""
        response = client.get('/api/collective/stats')

        assert response.status_code == 200
        assert 'application/json' in response.content_type

    @pytest.mark.api
    def test_stats_structure(self, client, sample_proposals):
        """Stats response has expected shape."""
        response = client.get('/api/collective/stats')
        data = json.loads(response.data)

        assert 'total_proposals' in data
        assert 'status_counts' in data
        assert 'total_consultations' in data
        assert 'contributor_count' in data
        assert 'contributors' in data

    @pytest.mark.api
    def test_stats_counts_match(self, client, sample_proposals):
        """Stats counts are consistent with sample data."""
        response = client.get('/api/collective/stats')
        data = json.loads(response.data)

        assert data['total_proposals'] == 3
        # test-proposal-001 (proposed), test-proposal-002 (consultation), test-proposal-003 (implemented)
        assert data['status_counts']['proposed'] == 1
        assert data['status_counts']['consultation'] == 1
        assert data['status_counts']['implemented'] == 1

    @pytest.mark.api
    def test_stats_consultations_count(self, client, sample_proposals):
        """Total consultations are correctly summed."""
        response = client.get('/api/collective/stats')
        data = json.loads(response.data)

        # test-proposal-002 has 2 consultations, test-proposal-003 has 1
        assert data['total_consultations'] == 3

    @pytest.mark.api
    def test_stats_contributors(self, client, sample_proposals):
        """Contributors list includes proposers and consultation contributors."""
        response = client.get('/api/collective/stats')
        data = json.loads(response.data)

        # Proposers: test-agent, test-agent-1
        # Consultation contributors: test-agent-2, test-agent-3
        assert data['contributor_count'] >= 3
        assert isinstance(data['contributors'], list)
        # Should be sorted
        assert data['contributors'] == sorted(data['contributors'])

    @pytest.mark.api
    def test_stats_empty_collection(self, client, empty_proposals_dir):
        """Stats work correctly with zero proposals."""
        response = client.get('/api/collective/stats')
        data = json.loads(response.data)

        assert data['total_proposals'] == 0
        assert data['total_consultations'] == 0
        assert data['contributor_count'] == 0
        assert data['contributors'] == []

    @pytest.mark.api
    def test_stats_all_statuses_present(self, client, sample_proposals):
        """Status counts dict includes all valid statuses, even those at zero."""
        response = client.get('/api/collective/stats')
        data = json.loads(response.data)

        for status in ['proposed', 'consultation', 'consensus', 'implemented', 'blocked', 'withdrawn']:
            assert status in data['status_counts']


# ---------------------------------------------------------------------------
# GET /api/proposals — filtering
# ---------------------------------------------------------------------------

class TestProposalFiltering:
    """Tests for query-string filtering on the proposals list."""

    @pytest.mark.api
    def test_filter_by_status(self, client, sample_proposals):
        """?status=proposed returns only proposed proposals."""
        response = client.get('/api/proposals?status=proposed')
        data = json.loads(response.data)

        assert data['count'] == 1
        assert all(p['status'] == 'proposed' for p in data['proposals'])

    @pytest.mark.api
    def test_filter_by_urgency(self, client, sample_proposals):
        """?urgency=high returns only high-urgency proposals."""
        response = client.get('/api/proposals?urgency=high')
        data = json.loads(response.data)

        assert data['count'] == 1
        assert all(p['urgency'] == 'high' for p in data['proposals'])

    @pytest.mark.api
    def test_filter_no_match(self, client, sample_proposals):
        """Filtering with no matches returns empty list, not 404."""
        response = client.get('/api/proposals?status=blocked')
        data = json.loads(response.data)

        assert response.status_code == 200
        assert data['count'] == 0
        assert data['proposals'] == []

    @pytest.mark.api
    def test_filter_combined(self, client, sample_proposals):
        """Multiple filters are ANDed together."""
        response = client.get('/api/proposals?status=consultation&urgency=high')
        data = json.loads(response.data)

        assert data['count'] == 1
        assert data['proposals'][0]['id'] == 'test-proposal-002'

    @pytest.mark.api
    def test_no_filter_returns_all(self, client, sample_proposals):
        """No query parameters returns all proposals."""
        response = client.get('/api/proposals')
        data = json.loads(response.data)

        assert data['count'] == 3


# ---------------------------------------------------------------------------
# GET /api/proposals/<id> — canonical plural route
# ---------------------------------------------------------------------------

class TestCanonicalProposalRoute:
    """Tests for the canonical /api/proposals/<id> route."""

    @pytest.mark.api
    def test_plural_route_works(self, client, sample_proposals):
        """GET /api/proposals/<id> returns the proposal."""
        response = client.get('/api/proposals/test-proposal-001')
        assert response.status_code == 200

        data = json.loads(response.data)
        assert data['id'] == 'test-proposal-001'

    @pytest.mark.api
    def test_plural_route_404(self, client, sample_proposals):
        """GET /api/proposals/<id> returns 404 with error format."""
        response = client.get('/api/proposals/nonexistent')
        assert response.status_code == 404

        data = json.loads(response.data)
        assert 'error' in data
        assert 'code' in data
        assert data['code'] == 'NOT_FOUND'

    @pytest.mark.api
    def test_singular_route_still_works(self, client, sample_proposals):
        """Legacy /api/proposal/<id> route still returns data."""
        response = client.get('/api/proposal/test-proposal-001')
        assert response.status_code == 200

        data = json.loads(response.data)
        assert data['id'] == 'test-proposal-001'


# ---------------------------------------------------------------------------
# Error format consistency
# ---------------------------------------------------------------------------

class TestErrorFormat:
    """All API errors must return {"error": "message", "code": "ERROR_CODE"}."""

    @pytest.mark.api
    def test_404_has_error_and_code(self, client, sample_proposals):
        """404 responses include both error and code fields."""
        response = client.get('/api/proposals/nope')
        data = json.loads(response.data)

        assert 'error' in data
        assert 'code' in data
        assert isinstance(data['error'], str)
        assert isinstance(data['code'], str)

    @pytest.mark.api
    def test_422_has_error_and_code(self, client, temp_data_dir):
        """422 validation errors include both fields."""
        response = client.post('/api/proposals',
            data=json.dumps({}),
            content_type='application/json')

        assert response.status_code == 422
        data = json.loads(response.data)
        assert 'error' in data
        assert 'code' in data

    @pytest.mark.api
    def test_415_has_error_and_code(self, client, temp_data_dir):
        """415 content-type errors include both fields."""
        response = client.post('/api/proposals',
            data='not json',
            content_type='text/plain')

        assert response.status_code == 415
        data = json.loads(response.data)
        assert 'error' in data
        assert 'code' in data

    @pytest.mark.api
    def test_409_has_error_and_code(self, client, sample_proposals):
        """409 conflict errors include both fields."""
        response = client.put('/api/proposals/test-proposal-001/status',
            data=json.dumps({
                'status': 'implemented',
                'actor': 'agent',
            }),
            content_type='application/json')

        assert response.status_code == 409
        data = json.loads(response.data)
        assert 'error' in data
        assert 'code' in data
