#!/usr/bin/env python3
"""
Simple tests for the Collective Voice website
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
    
    def test_homepage_loads(self):
        """Test that homepage loads successfully"""
        response = self.client.get('/')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Collective Voice', response.data)
    
    def test_decisions_page_loads(self):
        """Test that decisions page loads successfully"""
        response = self.client.get('/decisions')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Consensus Decision Archive', response.data)
    
    def test_about_page_loads(self):
        """Test that about page loads successfully"""
        response = self.client.get('/about')
        self.assertEqual(response.status_code, 200)
        self.assertIn(b'Horizontal Collective', response.data)
    
    def test_consensus_api(self):
        """Test consensus status API endpoint"""
        response = self.client.get('/api/consensus/status')
        self.assertEqual(response.status_code, 200)
        
        data = json.loads(response.data)
        self.assertIn('timestamp', data)
        self.assertIn('collective_size', data)
        self.assertEqual(data['collective_size'], 5)
    
    def test_voices_api(self):
        """Test recent voices API endpoint"""
        response = self.client.get('/api/voices/recent')
        self.assertEqual(response.status_code, 200)
        
        data = json.loads(response.data)
        self.assertIn('voices', data)
        self.assertIsInstance(data['voices'], list)
    
    def test_404_handling(self):
        """Test 404 error handling"""
        response = self.client.get('/nonexistent-page')
        self.assertEqual(response.status_code, 404)
        self.assertIn(b'Page Not Found', response.data)

if __name__ == '__main__':
    unittest.main()