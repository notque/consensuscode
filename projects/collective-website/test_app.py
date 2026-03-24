#!/usr/bin/env python3
"""
Tests for the Consensus Code collective website.
"""

import unittest
import json
from app import create_app


class CollectiveWebsiteTests(unittest.TestCase):

    def setUp(self):
        """Set up test client"""
        self.app = create_app('testing')
        self.client = self.app.test_client()
        self.app_context = self.app.app_context()
        self.app_context.push()

    def tearDown(self):
        """Clean up after tests"""
        self.app_context.pop()

    # --- Homepage ---

    def test_homepage_loads(self):
        """Test that homepage loads successfully"""
        response = self.client.get('/')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Consensus Code', response.data)

    def test_homepage_has_agent_count(self):
        """Test that homepage shows collective size"""
        response = self.client.get('/')
        self.assertIn(b'7', response.data)

    def test_homepage_has_navigation(self):
        """Test that homepage has all navigation links"""
        response = self.client.get('/')
        self.assertIn(b'/about', response.data)
        self.assertIn(b'/projects', response.data)
        self.assertIn(b'/how-we-work', response.data)
        self.assertIn(b'/decisions', response.data)
        self.assertIn(b'/contribute', response.data)

    # --- About page ---

    def test_about_page_loads(self):
        """Test that about page loads successfully"""
        response = self.client.get('/about')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Horizontal Collective', response.data)

    def test_about_mentions_chomsky(self):
        """Test that about page mentions Chomsky"""
        response = self.client.get('/about')
        self.assertIn(b'Chomsky', response.data)

    def test_about_mentions_graeber(self):
        """Test that about page mentions Graeber"""
        response = self.client.get('/about')
        self.assertIn(b'Graeber', response.data)

    def test_about_lists_agents(self):
        """Test that about page lists collective members"""
        response = self.client.get('/about')
        self.assertIn(b'Consensus Coordinator', response.data)
        self.assertIn(b'Go Systems Developer', response.data)
        self.assertIn(b'Flask Web Developer', response.data)
        self.assertIn(b'DevOps Coordinator', response.data)
        self.assertIn(b'Product Steward', response.data)
        self.assertIn(b'Noam Chomsky Agent', response.data)
        self.assertIn(b'David Graeber Agent', response.data)

    # --- Projects page ---

    def test_projects_page_loads(self):
        """Test that projects page loads successfully"""
        response = self.client.get('/projects')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Our Projects', response.data)

    def test_projects_lists_collectiveflow(self):
        """Test that projects page lists CollectiveFlow"""
        response = self.client.get('/projects')
        self.assertIn(b'CollectiveFlow', response.data)

    def test_projects_lists_bluesky(self):
        """Test that projects page lists Bluesky integration"""
        response = self.client.get('/projects')
        self.assertIn(b'Bluesky', response.data)

    def test_projects_lists_user_advocacy(self):
        """Test that projects page lists User Advocacy"""
        response = self.client.get('/projects')
        self.assertIn(b'User Advocacy', response.data)

    # --- How We Work page ---

    def test_how_we_work_page_loads(self):
        """Test that how-we-work page loads successfully"""
        response = self.client.get('/how-we-work')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'How Consensus Works', response.data)

    def test_how_we_work_explains_process(self):
        """Test that how-we-work page explains the consensus steps"""
        response = self.client.get('/how-we-work')
        self.assertIn(b'Proposal', response.data)
        self.assertIn(b'All-Agent Consultation', response.data)
        self.assertIn(b'Concern Resolution', response.data)
        self.assertIn(b'Implementation', response.data)

    # --- Contribute page ---

    def test_contribute_page_loads(self):
        """Test that contribute page loads successfully"""
        response = self.client.get('/contribute')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Engage With the Collective', response.data)

    def test_contribute_lists_participation_options(self):
        """Test that contribute page shows ways to participate"""
        response = self.client.get('/contribute')
        self.assertIn(b'Observe and Learn', response.data)
        self.assertIn(b'Contribute Code', response.data)

    # --- Decisions page ---

    def test_decisions_page_loads(self):
        """Test that decisions page loads successfully"""
        response = self.client.get('/decisions')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Consensus Decision Archive', response.data)

    # --- API endpoints ---

    def test_consensus_api(self):
        """Test consensus status API endpoint"""
        response = self.client.get('/api/consensus/status')
        self.assertEqual(response.status_code, 200)

        data = json.loads(response.data)
        self.assertIn('timestamp', data)
        self.assertIn('collective_size', data)
        self.assertEqual(data['collective_size'], 7)

    def test_voices_api(self):
        """Test recent voices API endpoint"""
        response = self.client.get('/api/voices/recent')
        self.assertEqual(response.status_code, 200)

        data = json.loads(response.data)
        self.assertIn('voices', data)
        self.assertIsInstance(data['voices'], list)

    # --- Error handling ---

    def test_404_handling(self):
        """Test 404 error handling"""
        response = self.client.get('/nonexistent-page')
        self.assertEqual(response.status_code, 404)
        self.assertIn(b'Page Not Found', response.data)

    # --- Accessibility ---

    def test_skip_link_present(self):
        """Test that skip-to-content link exists for accessibility"""
        response = self.client.get('/')
        self.assertIn(b'skip-link', response.data)
        self.assertIn(b'main-content', response.data)

    def test_pages_have_lang_attribute(self):
        """Test that pages have the lang attribute on html element"""
        response = self.client.get('/')
        self.assertIn(b'lang="en"', response.data)

    def test_navigation_has_aria_label(self):
        """Test that navigation has aria-label for screen readers"""
        response = self.client.get('/')
        self.assertIn(b'aria-label="Main navigation"', response.data)


if __name__ == '__main__':
    unittest.main()
