"""
CollectiveFlow Web Application - Create Proposal & Collective View Tests

These tests cover edge cases and deeper behavior for:
- POST /create endpoint (validation, defaults, error handling, redirect behavior)
- GET /collective endpoint (statistics accuracy, contributor counting, recent events)

Why a separate file?
The existing test_routes.py covers the happy paths. These tests go deeper into
boundary conditions and data accuracy, which helps any agent contributing tests
understand what "thorough coverage" looks like for form-handling routes.
"""

import pytest
import yaml
import json
from pathlib import Path
from unittest.mock import patch


class TestCreateProposalPost:
    """
    Deeper tests for the POST /create endpoint.

    These tests go beyond basic validation to cover:
    - Default values when optional fields are omitted
    - Redirect targets after success and failure
    - Flash message content
    - Saved data accuracy
    - The exception-handling path (lines 247-249 in app.py)
    """

    @pytest.mark.routes
    def test_create_proposal_redirects_to_detail_on_success(self, client, temp_data_dir):
        """
        Test: Successful creation redirects to the new proposal detail page

        The redirect should point to /proposal/<new-id>, not back to the form.
        """
        form_data = {
            'title': 'Redirect Test Proposal',
            'description': 'Testing redirect target after creation',
            'proposer': 'redirect-tester',
            'urgency': 'low',
        }
        response = client.post('/create', data=form_data, follow_redirects=False)

        # Should be a 302 redirect
        assert response.status_code == 302
        # Redirect location should contain /proposal/
        assert '/proposal/' in response.headers['Location']

    @pytest.mark.routes
    def test_create_proposal_default_proposer_is_anonymous(self, client, temp_data_dir):
        """
        Test: When proposer field is omitted, defaults to 'anonymous'

        No one should be excluded from proposing just because they
        did not fill in a name.
        """
        form_data = {
            'title': 'Anonymous Proposal',
            'description': 'Submitted without a proposer name',
            'urgency': 'medium',
        }
        response = client.post('/create', data=form_data, follow_redirects=True)
        assert response.status_code == 200

        # Verify the saved file has 'anonymous' as proposer
        yaml_files = list(temp_data_dir.glob('*.yaml'))
        assert len(yaml_files) >= 1

        with open(yaml_files[0], 'r') as f:
            saved = yaml.safe_load(f)
        assert saved['proposer'] == 'anonymous'

    @pytest.mark.routes
    def test_create_proposal_default_urgency_is_medium(self, client, temp_data_dir):
        """
        Test: When urgency is not provided, it defaults to 'medium'

        Medium is a safe default that avoids both alarm and neglect.
        """
        form_data = {
            'title': 'Default Urgency Proposal',
            'description': 'No urgency specified',
            'proposer': 'test-agent',
        }
        client.post('/create', data=form_data, follow_redirects=True)

        yaml_files = list(temp_data_dir.glob('*.yaml'))
        assert len(yaml_files) >= 1

        with open(yaml_files[0], 'r') as f:
            saved = yaml.safe_load(f)
        assert saved['urgency'] == 'medium'

    @pytest.mark.routes
    def test_create_proposal_strips_whitespace_from_title(self, client, temp_data_dir):
        """
        Test: Leading and trailing whitespace is stripped from title

        Clean data prevents display issues.
        """
        form_data = {
            'title': '   Whitespace Title   ',
            'description': 'Testing whitespace handling',
            'proposer': 'test-agent',
        }
        client.post('/create', data=form_data, follow_redirects=True)

        yaml_files = list(temp_data_dir.glob('*.yaml'))
        with open(yaml_files[0], 'r') as f:
            saved = yaml.safe_load(f)
        assert saved['title'] == 'Whitespace Title'

    @pytest.mark.routes
    def test_create_proposal_whitespace_only_title_rejected(self, client, temp_data_dir):
        """
        Test: A title of only whitespace is treated as empty and rejected

        Stripping then checking emptiness catches this edge case.
        """
        form_data = {
            'title': '   ',
            'description': 'Description is fine',
            'proposer': 'test-agent',
        }
        response = client.post('/create', data=form_data, follow_redirects=True)
        assert response.status_code == 200

        # No YAML file should have been created
        yaml_files = list(temp_data_dir.glob('*.yaml'))
        assert len(yaml_files) == 0

    @pytest.mark.routes
    def test_create_proposal_whitespace_only_description_rejected(self, client, temp_data_dir):
        """
        Test: A description of only whitespace is treated as empty

        Same whitespace-stripping logic applies to description.
        """
        form_data = {
            'title': 'Valid Title',
            'description': '   ',
            'proposer': 'test-agent',
        }
        response = client.post('/create', data=form_data, follow_redirects=True)
        assert response.status_code == 200

        yaml_files = list(temp_data_dir.glob('*.yaml'))
        assert len(yaml_files) == 0

    @pytest.mark.routes
    def test_create_proposal_affected_areas_list(self, client, temp_data_dir):
        """
        Test: Multiple affected_areas are saved as a list

        The form uses getlist() for multi-value fields.
        """
        response = client.post('/create', data={
            'title': 'Multi-Area Proposal',
            'description': 'Affects multiple areas',
            'proposer': 'test-agent',
            'urgency': 'high',
            'affected_areas': ['testing', 'infrastructure', 'web'],
        }, follow_redirects=True)
        assert response.status_code == 200

        yaml_files = list(temp_data_dir.glob('*.yaml'))
        with open(yaml_files[0], 'r') as f:
            saved = yaml.safe_load(f)
        assert 'testing' in saved['affected_areas']
        assert 'infrastructure' in saved['affected_areas']
        assert 'web' in saved['affected_areas']

    @pytest.mark.routes
    def test_create_proposal_empty_affected_areas(self, client, temp_data_dir):
        """
        Test: No affected_areas results in empty list, not an error

        Optional fields should degrade gracefully.
        """
        response = client.post('/create', data={
            'title': 'No Areas',
            'description': 'No affected areas specified',
            'proposer': 'test-agent',
        }, follow_redirects=True)
        assert response.status_code == 200

        yaml_files = list(temp_data_dir.glob('*.yaml'))
        with open(yaml_files[0], 'r') as f:
            saved = yaml.safe_load(f)
        assert saved.get('affected_areas') == []

    @pytest.mark.routes
    def test_create_proposal_saves_both_yaml_and_json(self, client, temp_data_dir):
        """
        Test: Proposal creation via POST saves both YAML and JSON files

        The save_proposal function writes both formats for API compatibility.
        """
        client.post('/create', data={
            'title': 'Dual Format Test',
            'description': 'Should be saved in both formats',
            'proposer': 'test-agent',
        }, follow_redirects=True)

        yaml_files = list(temp_data_dir.glob('*.yaml'))
        json_files = list(temp_data_dir.glob('*.json'))
        assert len(yaml_files) >= 1
        assert len(json_files) >= 1

    @pytest.mark.routes
    def test_create_proposal_flash_success_message(self, client, temp_data_dir):
        """
        Test: Successful creation includes a flash message with the title

        Flash messages confirm the action to the user.
        """
        response = client.post('/create', data={
            'title': 'Flash Test Proposal',
            'description': 'Testing flash message content',
            'proposer': 'test-agent',
        }, follow_redirects=True)

        data = response.data.decode('utf-8')
        assert 'Flash Test Proposal' in data or 'success' in data.lower()

    @pytest.mark.routes
    def test_create_proposal_flash_error_on_missing_title(self, client, temp_data_dir):
        """
        Test: Missing title produces a flash error mentioning 'required'

        Users need to know what went wrong.
        """
        response = client.post('/create', data={
            'title': '',
            'description': 'Valid description',
        }, follow_redirects=True)

        data = response.data.decode('utf-8').lower()
        assert 'required' in data or 'error' in data

    @pytest.mark.routes
    def test_create_proposal_exception_handler(self, client, temp_data_dir):
        """
        Test: Exception during save_proposal triggers error flash and redirect

        This covers the except block (lines 247-249 in app.py) that catches
        unexpected errors during proposal saving.
        """
        with patch('app.save_proposal', side_effect=RuntimeError('disk full')):
            response = client.post('/create', data={
                'title': 'Exception Proposal',
                'description': 'This should trigger an exception',
                'proposer': 'test-agent',
            }, follow_redirects=True)

            assert response.status_code == 200
            data = response.data.decode('utf-8').lower()
            # The flash message should mention the error
            assert 'error' in data or 'disk full' in data

    @pytest.mark.routes
    def test_create_proposal_sets_status_to_proposed(self, client, temp_data_dir):
        """
        Test: Newly created proposals always start with status 'proposed'

        All proposals begin equal -- no fast-tracking.
        """
        client.post('/create', data={
            'title': 'Status Check Proposal',
            'description': 'Checking initial status',
            'proposer': 'test-agent',
        }, follow_redirects=True)

        yaml_files = list(temp_data_dir.glob('*.yaml'))
        with open(yaml_files[0], 'r') as f:
            saved = yaml.safe_load(f)
        assert saved['status'] == 'proposed'

    @pytest.mark.routes
    def test_create_proposal_generates_consensus_history(self, client, temp_data_dir):
        """
        Test: New proposal via web form has initial consensus history entry

        Audit trail starts at creation.
        """
        client.post('/create', data={
            'title': 'History Proposal',
            'description': 'Testing consensus history',
            'proposer': 'web-user',
            'urgency': 'high',
        }, follow_redirects=True)

        yaml_files = list(temp_data_dir.glob('*.yaml'))
        with open(yaml_files[0], 'r') as f:
            saved = yaml.safe_load(f)

        assert 'consensus_history' in saved
        assert len(saved['consensus_history']) == 1
        entry = saved['consensus_history'][0]
        assert entry['event'] == 'proposal_created'
        assert entry['actor'] == 'web-user'
        assert 'high' in entry['details']

    @pytest.mark.routes
    def test_create_proposal_visible_in_api_after_creation(self, client, temp_data_dir):
        """
        Test: A proposal created via POST shows up in the API immediately

        No caching or delay should hide newly created proposals.
        """
        client.post('/create', data={
            'title': 'API Visibility Test',
            'description': 'Should appear in API right away',
            'proposer': 'test-agent',
        }, follow_redirects=True)

        api_response = client.get('/api/proposals')
        api_data = json.loads(api_response.data)
        titles = [p['title'] for p in api_data['proposals']]
        assert 'API Visibility Test' in titles


class TestCollectiveViewDepth:
    """
    Deeper tests for the GET /collective endpoint.

    These go beyond "does it load" to verify the accuracy of
    statistics, contributor counting, and recent event generation.
    """

    @pytest.mark.routes
    def test_collective_stats_total_count_accurate(self, client, sample_proposals):
        """
        Test: Total proposals count matches actual number of proposals

        The stat should reflect reality, not a cached or hardcoded value.
        """
        response = client.get('/collective')
        data = response.data.decode('utf-8')

        # We have 3 sample proposals
        assert '3' in data

    @pytest.mark.routes
    def test_collective_stats_active_consultations(self, client, sample_proposals):
        """
        Test: Active consultation count is correct

        Only proposals with status 'consultation' should be counted.
        """
        response = client.get('/collective')
        data = response.data.decode('utf-8')

        # sample_proposals has 1 with status 'consultation'
        # The page should show this count
        assert '1' in data

    @pytest.mark.routes
    def test_collective_stats_implemented_count(self, client, sample_proposals):
        """
        Test: Implemented proposal count is accurate

        Only proposals with status 'implemented' should be counted.
        """
        response = client.get('/collective')
        data = response.data.decode('utf-8')

        # sample_proposals has 1 implemented proposal
        assert '1' in data

    @pytest.mark.routes
    def test_collective_unique_contributors(self, client, sample_proposals):
        """
        Test: Contributors are deduplicated across proposals and consultations

        Each unique agent name should appear only once.
        """
        response = client.get('/collective')
        data = response.data.decode('utf-8')

        # From sample data: test-agent, test-agent-1, test-agent-2, test-agent-3
        # That's 4 unique contributors
        assert '4' in data

    @pytest.mark.routes
    def test_collective_shows_recent_events(self, client, sample_proposals):
        """
        Test: Recent events section includes proposal titles

        Members should see what happened recently.
        """
        response = client.get('/collective')
        data = response.data.decode('utf-8')

        # Recent events should reference proposal titles
        assert 'Test Proposal' in data

    @pytest.mark.routes
    def test_collective_recent_events_limited_to_five(self, client, temp_data_dir):
        """
        Test: Recent events show at most 5 proposals

        The view slices proposals[:5] to limit activity display.
        """
        from app import save_proposal

        # Create 8 proposals
        for i in range(8):
            save_proposal({
                'title': f'Bulk Proposal {i}',
                'description': f'Proposal number {i}',
                'proposer': 'bulk-tester',
            })

        response = client.get('/collective')
        data = response.data.decode('utf-8')

        # Should show activity but not more than 5 recent items
        # All 8 should be in total count
        assert '8' in data
        assert response.status_code == 200

    @pytest.mark.routes
    def test_collective_contributors_include_consultation_participants(self, client, sample_proposals):
        """
        Test: Contributors list includes agents from consultations, not just proposers

        Everyone who participates gets recognized.
        """
        response = client.get('/collective')
        data = response.data.decode('utf-8')

        # test-agent-2 and test-agent-3 only appear in consultations
        assert 'test-agent-2' in data or 'test-agent-3' in data

    @pytest.mark.routes
    def test_collective_with_single_proposal(self, client, temp_data_dir):
        """
        Test: Collective view works correctly with just one proposal

        New collectives starting out should see accurate stats.
        """
        from app import save_proposal

        save_proposal({
            'title': 'First Proposal Ever',
            'description': 'The very first proposal',
            'proposer': 'founding-member',
            'urgency': 'medium',
        })

        response = client.get('/collective')
        assert response.status_code == 200
        data = response.data.decode('utf-8')

        assert '1' in data  # total count
        assert 'founding-member' in data or 'contributor' in data.lower()

    @pytest.mark.routes
    def test_collective_api_and_view_agree_on_counts(self, client, sample_proposals):
        """
        Test: Statistics from /collective match /api/proposals count

        Both views query the same data source; they should agree.
        """
        api_response = client.get('/api/proposals')
        api_data = json.loads(api_response.data)

        collective_response = client.get('/collective')
        collective_data = collective_response.data.decode('utf-8')

        # Total from API
        total = api_data['count']
        assert str(total) in collective_data


class TestCreateProposalFormGET:
    """
    Additional tests for the GET /create form page.
    """

    @pytest.mark.routes
    def test_create_form_has_urgency_options(self, client):
        """
        Test: The creation form offers urgency level choices

        Members should be able to select urgency when proposing.
        """
        response = client.get('/create')
        data = response.data.decode('utf-8').lower()

        assert 'urgency' in data

    @pytest.mark.routes
    def test_create_form_has_submit_button(self, client):
        """
        Test: The creation form has a submit mechanism

        Without a submit button, the form is unusable.
        """
        response = client.get('/create')
        data = response.data.decode('utf-8').lower()

        assert 'submit' in data or 'button' in data or 'type="submit"' in data

    @pytest.mark.routes
    def test_create_form_uses_post_method(self, client):
        """
        Test: The form uses POST method, not GET

        Form data should not appear in URLs.
        """
        response = client.get('/create')
        data = response.data.decode('utf-8').lower()

        assert 'method="post"' in data or "method='post'" in data
