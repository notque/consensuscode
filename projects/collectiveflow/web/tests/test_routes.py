"""
CollectiveFlow Web Application - Route Handler Tests

These tests validate the web application's route handlers, ensuring:
- Routes respond correctly to requests
- Status codes are appropriate
- HTML templates render properly
- Error conditions are handled gracefully
- The application embodies horizontal, non-hierarchical principles

Test Organization:
- Each test class focuses on a specific route or functionality area
- Test names describe what they're testing in plain language
- Comments explain the "why" behind each test
"""

import pytest
from flask import url_for


class TestHomeRoute:
    """
    Tests for the home page (/) route.

    The home page is the entry point to our collective's decision-making system.
    It should provide an accessible overview of all proposals grouped by status.
    """

    @pytest.mark.routes
    def test_home_page_loads_successfully(self, client, sample_proposals):
        """
        Test: Home page loads with 200 OK status

        The home page should always be accessible to everyone.
        This test ensures the basic route works.
        """
        response = client.get('/')
        assert response.status_code == 200

    @pytest.mark.routes
    def test_home_page_contains_title(self, client, sample_proposals):
        """
        Test: Home page displays CollectiveFlow title

        The title identifies our system to visitors.
        """
        response = client.get('/')
        assert b'CollectiveFlow' in response.data

    @pytest.mark.routes
    def test_home_page_shows_proposal_counts(self, client, sample_proposals):
        """
        Test: Home page displays statistics about proposals

        Transparency requires showing collective activity metrics.
        We should see counts of proposals in different states.
        """
        response = client.get('/')
        data = response.data.decode('utf-8')

        # Check that statistics section exists
        assert 'Total Proposals' in data
        assert 'Active Discussions' in data
        assert 'Implemented' in data
        assert 'Consensus Reached' in data

    @pytest.mark.routes
    def test_home_page_lists_recent_proposals(self, client, sample_proposals):
        """
        Test: Home page displays recent proposals

        Members should see recent activity to stay informed.
        """
        response = client.get('/')
        data = response.data.decode('utf-8')

        # Check for proposal titles from our sample data
        assert 'Test Proposal: Simple Example' in data
        assert 'Test Proposal: With Consultations' in data

    @pytest.mark.routes
    def test_home_page_with_no_proposals(self, client, empty_proposals_dir):
        """
        Test: Home page works even with no proposals

        The application should handle empty state gracefully.
        This is important for new collectives just starting out.
        """
        response = client.get('/')
        assert response.status_code == 200
        # Should still show the structure, just with zero counts
        assert b'Total Proposals' in response.data

    @pytest.mark.routes
    def test_home_page_groups_proposals_by_status(self, client, sample_proposals):
        """
        Test: Home page organizes proposals by their status

        Grouping helps members quickly find proposals needing attention.
        """
        response = client.get('/')
        data = response.data.decode('utf-8')

        # We have proposals with different statuses
        # Check that consultation status proposal appears in appropriate section
        assert 'consultation' in data.lower()
        assert 'implemented' in data.lower()


class TestProposalsListRoute:
    """
    Tests for the proposals list (/proposals) route.

    This route provides a comprehensive view of all proposals,
    organized by status for easy navigation.
    """

    @pytest.mark.routes
    def test_proposals_list_loads(self, client, sample_proposals):
        """
        Test: Proposals list page loads successfully

        All members should be able to view the complete list of proposals.
        """
        response = client.get('/proposals')
        assert response.status_code == 200

    @pytest.mark.routes
    def test_proposals_list_shows_all_proposals(self, client, sample_proposals):
        """
        Test: Proposals list displays all proposals from sample data

        Transparency requires showing all proposals, not filtering based on roles.
        """
        response = client.get('/proposals')
        data = response.data.decode('utf-8')

        # All three sample proposals should be visible
        assert 'Test Proposal: Simple Example' in data
        assert 'Test Proposal: With Consultations' in data
        assert 'Test Proposal: Implemented' in data

    @pytest.mark.routes
    def test_proposals_list_groups_by_status(self, client, sample_proposals):
        """
        Test: Proposals are organized by status categories

        Grouping makes it easier to find proposals in specific states.
        """
        response = client.get('/proposals')
        data = response.data.decode('utf-8')

        # Check for status section headers
        # (The exact HTML structure may vary, but status names should appear)
        assert 'proposed' in data.lower() or 'consultation' in data.lower()


class TestProposalDetailRoute:
    """
    Tests for individual proposal detail (/proposal/<id>) route.

    Each proposal should have a detailed view showing all information,
    consultations, and decision history.
    """

    @pytest.mark.routes
    def test_proposal_detail_loads(self, client, sample_proposals):
        """
        Test: Proposal detail page loads for valid proposal

        Members should be able to view any proposal's full details.
        """
        response = client.get('/proposal/test-proposal-001')
        assert response.status_code == 200

    @pytest.mark.routes
    def test_proposal_detail_shows_title(self, client, sample_proposals):
        """
        Test: Proposal detail displays the proposal title

        The title is the primary identifier for members.
        """
        response = client.get('/proposal/test-proposal-001')
        assert b'Test Proposal: Simple Example' in response.data

    @pytest.mark.routes
    def test_proposal_detail_shows_description(self, client, sample_proposals):
        """
        Test: Proposal detail displays full description

        Members need complete information to participate in consensus.
        """
        response = client.get('/proposal/test-proposal-001')
        assert b'A simple test proposal' in response.data

    @pytest.mark.routes
    def test_proposal_detail_shows_consultations(self, client, sample_proposals):
        """
        Test: Proposal with consultations displays all input

        Transparency requires showing all perspectives and input.
        """
        response = client.get('/proposal/test-proposal-002')
        data = response.data.decode('utf-8')

        # Should show consultation contributors and their input
        assert 'test-agent-2' in data
        assert 'test-agent-3' in data
        assert 'I support this proposal' in data
        assert 'I have concerns' in data

    @pytest.mark.routes
    def test_proposal_detail_not_found(self, client, sample_proposals):
        """
        Test: Non-existent proposal returns 404

        The application should handle invalid proposal IDs gracefully.
        """
        response = client.get('/proposal/nonexistent-proposal')
        assert response.status_code == 404

    @pytest.mark.routes
    def test_proposal_detail_shows_status_badge(self, client, sample_proposals):
        """
        Test: Proposal detail shows current status

        Status visibility helps members understand proposal state.
        """
        response = client.get('/proposal/test-proposal-002')
        data = response.data.decode('utf-8')

        # Should show the consultation status
        assert 'consultation' in data.lower()

    @pytest.mark.routes
    def test_proposal_detail_shows_urgency(self, client, sample_proposals):
        """
        Test: Proposal detail displays urgency level

        Urgency helps the collective prioritize attention.
        """
        response = client.get('/proposal/test-proposal-002')
        data = response.data.decode('utf-8')

        # Sample proposal 002 has high urgency
        assert 'high' in data.lower()


class TestAboutRoute:
    """
    Tests for the about page (/about) route.

    The about page explains our collective principles and decision-making process.
    It's crucial for onboarding and transparency.
    """

    @pytest.mark.routes
    def test_about_page_loads(self, client):
        """
        Test: About page loads successfully

        The about page should always be accessible to explain our principles.
        """
        response = client.get('/about')
        assert response.status_code == 200

    @pytest.mark.routes
    def test_about_page_mentions_horizontal_principles(self, client):
        """
        Test: About page explains horizontal organization

        We must clearly communicate our non-hierarchical approach.
        """
        response = client.get('/about')
        data = response.data.decode('utf-8').lower()

        # Should mention key principles
        # (Exact wording may vary, but these concepts should be present)
        assert 'horizontal' in data or 'consensus' in data or 'collective' in data


class TestCollectiveRoute:
    """
    Tests for the collective status page (/collective) route.

    This page shows the collective's current activity and statistics,
    providing transparency about our decision-making process.
    """

    @pytest.mark.routes
    def test_collective_page_loads(self, client, sample_proposals):
        """
        Test: Collective page loads successfully

        Members should be able to view collective statistics.
        """
        response = client.get('/collective')
        assert response.status_code == 200

    @pytest.mark.routes
    def test_collective_page_shows_statistics(self, client, sample_proposals):
        """
        Test: Collective page displays proposal statistics

        Transparency requires showing collective activity metrics.
        """
        response = client.get('/collective')
        data = response.data.decode('utf-8')

        # Should show various statistics
        assert 'total' in data.lower() or 'proposals' in data.lower()

    @pytest.mark.routes
    def test_collective_page_shows_contributors(self, client, sample_proposals):
        """
        Test: Collective page acknowledges all contributors

        Everyone who participates should be recognized equally.
        """
        response = client.get('/collective')
        data = response.data.decode('utf-8')

        # Our sample data includes various contributors
        # Should show contributor information
        assert 'contributor' in data.lower() or 'agent' in data.lower()

    @pytest.mark.routes
    def test_collective_page_with_no_proposals(self, client, empty_proposals_dir):
        """
        Test: Collective page works with no data

        The page should handle empty state for new collectives.
        """
        response = client.get('/collective')
        assert response.status_code == 200


class TestCreateProposalRoutes:
    """
    Tests for proposal creation routes (GET and POST /create).

    These routes allow members to submit new proposals to the collective.
    Testing ensures the submission process works correctly.
    """

    @pytest.mark.routes
    def test_create_form_loads(self, client):
        """
        Test: Proposal creation form loads successfully

        All members should be able to access the proposal creation form.
        """
        response = client.get('/create')
        assert response.status_code == 200

    @pytest.mark.routes
    def test_create_form_has_required_fields(self, client):
        """
        Test: Creation form includes all necessary input fields

        The form should collect all information needed for a valid proposal.
        """
        response = client.get('/create')
        data = response.data.decode('utf-8')

        # Check for form fields (looking for input names or labels)
        assert 'title' in data.lower()
        assert 'description' in data.lower()

    @pytest.mark.routes
    def test_create_proposal_success(self, client, csrf_token, temp_data_dir, proposal_form_data):
        """
        Test: Submitting valid proposal data creates new proposal

        Members should be able to successfully submit proposals.
        """
        proposal_form_data['_csrf_token'] = csrf_token
        response = client.post('/create', data=proposal_form_data, follow_redirects=True)

        # Should redirect to the new proposal's detail page
        assert response.status_code == 200

        # Check that proposal was created in filesystem
        yaml_files = list(temp_data_dir.glob('*.yaml'))
        assert len(yaml_files) >= 1

    @pytest.mark.routes
    def test_create_proposal_missing_title(self, client, csrf_token, proposal_form_data):
        """
        Test: Creating proposal without title shows error

        Required fields should be validated.
        """
        # Remove title from form data
        form_data = proposal_form_data.copy()
        form_data['title'] = ''
        form_data['_csrf_token'] = csrf_token

        response = client.post('/create', data=form_data, follow_redirects=True)

        # Should show error message
        assert response.status_code == 200
        # Flash messages appear in the response
        assert b'required' in response.data.lower() or b'error' in response.data.lower()

    @pytest.mark.routes
    def test_create_proposal_missing_description(self, client, csrf_token, proposal_form_data):
        """
        Test: Creating proposal without description shows error

        Description is essential for informed consensus.
        """
        form_data = proposal_form_data.copy()
        form_data['description'] = ''
        form_data['_csrf_token'] = csrf_token

        response = client.post('/create', data=form_data, follow_redirects=True)

        # Should show error message
        assert response.status_code == 200
        assert b'required' in response.data.lower() or b'error' in response.data.lower()


class TestNavigationAndLinks:
    """
    Tests for navigation and inter-page links.

    Ensuring all pages are properly linked maintains accessibility
    and helps members navigate the system.
    """

    @pytest.mark.routes
    def test_navigation_links_present(self, client, sample_proposals):
        """
        Test: Navigation links appear on all pages

        Members should be able to navigate between sections easily.
        """
        pages = ['/', '/proposals', '/about', '/collective']

        for page in pages:
            response = client.get(page)
            data = response.data.decode('utf-8')

            # Check for navigation links to other pages
            assert '<nav' in data.lower() or 'href=' in data

    @pytest.mark.routes
    def test_proposal_links_from_home(self, client, sample_proposals):
        """
        Test: Home page has links to individual proposals

        Members should be able to click through to proposal details.
        """
        response = client.get('/')
        data = response.data.decode('utf-8')

        # Should have links to proposal detail pages
        assert '/proposal/test-proposal-001' in data or 'test-proposal-001' in data

    @pytest.mark.routes
    def test_no_login_required_notice(self, client):
        """
        Test: Site clearly indicates no login is required

        This reinforces our horizontal principles - no gatekeeping.
        """
        response = client.get('/')
        data = response.data.decode('utf-8')

        # Should mention no login or authentication required
        assert 'no login' in data.lower() or 'no special roles' in data.lower()


class TestErrorHandling:
    """
    Tests for error handling and edge cases.

    Robust error handling ensures the application remains accessible
    even when things go wrong.
    """

    @pytest.mark.routes
    def test_404_on_invalid_route(self, client):
        """
        Test: Invalid routes return 404 status

        The application should handle unknown routes gracefully.
        """
        response = client.get('/this-route-does-not-exist')
        assert response.status_code == 404

    @pytest.mark.routes
    def test_proposal_not_found_returns_404(self, client, sample_proposals):
        """
        Test: Non-existent proposal ID returns 404

        Invalid proposal requests should fail gracefully.
        """
        response = client.get('/proposal/invalid-id-12345')
        assert response.status_code == 404
        assert b'not found' in response.data.lower()


class TestResponsiveDesign:
    """
    Tests for responsive design elements.

    The interface should work across different devices and screen sizes,
    ensuring accessibility for all members.
    """

    @pytest.mark.routes
    def test_viewport_meta_tag_present(self, client):
        """
        Test: Pages include viewport meta tag for mobile responsiveness

        This ensures the site works on mobile devices.
        """
        response = client.get('/')
        assert b'viewport' in response.data

    @pytest.mark.routes
    def test_responsive_framework_loaded(self, client):
        """
        Test: CSS framework (Tailwind) is loaded

        The responsive design framework should be available.
        """
        response = client.get('/')
        # Check for Tailwind CSS
        assert b'tailwindcss' in response.data
