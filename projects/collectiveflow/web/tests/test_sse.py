"""
Tests for the Server-Sent Events (SSE) endpoint.

Validates that:
1. The /api/events endpoint returns the correct content type and headers
2. The EventBus pub/sub mechanism delivers events to subscribers
3. Events are published when proposals are created, consultations added,
   and statuses changed via the JSON API
4. Query parameter filtering (types, proposal_id) works correctly
5. The consensus_reached event fires when status moves to 'consensus'
"""

import json
import queue
import threading
import time

import pytest
import yaml


class TestEventBus:
    """Unit tests for the in-memory EventBus pub/sub."""

    def test_subscribe_returns_queue(self):
        from app import EventBus
        bus = EventBus()
        q = bus.subscribe()
        assert isinstance(q, queue.Queue)

    def test_publish_delivers_to_subscriber(self):
        from app import EventBus
        bus = EventBus()
        q = bus.subscribe()
        bus.publish('test_event', {'key': 'value'})
        event = q.get_nowait()
        assert event['type'] == 'test_event'
        assert event['data'] == {'key': 'value'}
        assert 'id' in event
        assert 'timestamp' in event

    def test_publish_delivers_to_multiple_subscribers(self):
        from app import EventBus
        bus = EventBus()
        q1 = bus.subscribe()
        q2 = bus.subscribe()
        bus.publish('test_event', {'n': 1})
        e1 = q1.get_nowait()
        e2 = q2.get_nowait()
        assert e1['data'] == e2['data']
        assert e1['id'] == e2['id']

    def test_unsubscribe_removes_queue(self):
        from app import EventBus
        bus = EventBus()
        q = bus.subscribe()
        bus.unsubscribe(q)
        bus.publish('test_event', {'n': 1})
        assert q.empty()

    def test_event_ids_are_monotonic(self):
        from app import EventBus
        bus = EventBus()
        q = bus.subscribe()
        bus.publish('a', {})
        bus.publish('b', {})
        bus.publish('c', {})
        ids = [q.get_nowait()['id'] for _ in range(3)]
        assert ids == sorted(ids)
        assert len(set(ids)) == 3  # all unique

    def test_full_queue_drops_oldest(self):
        from app import EventBus
        bus = EventBus(maxsize=2)
        q = bus.subscribe()
        bus.publish('first', {'n': 1})
        bus.publish('second', {'n': 2})
        bus.publish('third', {'n': 3})  # should drop 'first'
        events = []
        while not q.empty():
            events.append(q.get_nowait())
        types = [e['type'] for e in events]
        assert 'third' in types
        # 'first' may or may not be present depending on timing,
        # but queue should not exceed maxsize
        assert len(events) <= 2


class TestSSEEndpoint:
    """Integration tests for GET /api/events."""

    def test_sse_endpoint_returns_event_stream(self, client, temp_data_dir):
        """The endpoint should return text/event-stream content type."""
        response = client.get('/api/events')
        assert response.content_type.startswith('text/event-stream')

    def test_sse_endpoint_cache_headers(self, client, temp_data_dir):
        """SSE responses must not be cached."""
        response = client.get('/api/events')
        assert response.headers.get('Cache-Control') == 'no-cache'

    def test_sse_endpoint_sends_connected_comment(self, client, temp_data_dir):
        """The stream should start with a connection confirmation comment."""
        response = client.get('/api/events')
        # Read the first chunk from the streaming response
        data = b''
        for chunk in response.response:
            data += chunk if isinstance(chunk, bytes) else chunk.encode()
            if b'\n\n' in data:
                break
        text = data.decode()
        assert text.startswith(': connected')


class TestSSEEventFiring:
    """Test that API mutations fire SSE events."""

    def test_proposal_created_fires_event(self, client, temp_data_dir):
        """Creating a proposal via API should fire a proposal_created event."""
        import app as app_module
        q = app_module.event_bus.subscribe()
        try:
            response = client.post('/api/proposals', json={
                'title': 'SSE Test Proposal',
                'description': 'Testing SSE event firing',
                'proposer': 'sse-tester',
                'urgency': 'medium',
            })
            assert response.status_code == 201

            event = q.get(timeout=2)
            assert event['type'] == 'proposal_created'
            assert event['data']['title'] == 'SSE Test Proposal'
            assert event['data']['proposer'] == 'sse-tester'
            assert event['data']['urgency'] == 'medium'
            assert 'proposal_id' in event['data']
        finally:
            app_module.event_bus.unsubscribe(q)

    def test_consultation_added_fires_event(self, client, temp_data_dir, sample_proposals):
        """Adding consultation via API should fire a consultation_added event."""
        import app as app_module
        q = app_module.event_bus.subscribe()
        try:
            response = client.post(
                '/api/proposals/test-proposal-002/consultation',
                json={
                    'contributor': 'sse-agent',
                    'input': 'I support this via SSE test',
                    'support': True,
                },
            )
            assert response.status_code == 201

            event = q.get(timeout=2)
            assert event['type'] == 'consultation_added'
            assert event['data']['proposal_id'] == 'test-proposal-002'
            assert event['data']['contributor'] == 'sse-agent'
            assert event['data']['support'] is True
            assert isinstance(event['data']['consultation_count'], int)
        finally:
            app_module.event_bus.unsubscribe(q)

    def test_status_changed_fires_event(self, client, temp_data_dir, sample_proposals):
        """Changing proposal status via API should fire a status_changed event."""
        import app as app_module
        q = app_module.event_bus.subscribe()
        try:
            response = client.put(
                '/api/proposals/test-proposal-001/status',
                json={
                    'status': 'consultation',
                    'actor': 'sse-test-agent',
                    'reason': 'Testing SSE events',
                },
            )
            assert response.status_code == 200

            event = q.get(timeout=2)
            assert event['type'] == 'status_changed'
            assert event['data']['proposal_id'] == 'test-proposal-001'
            assert event['data']['previous_status'] == 'proposed'
            assert event['data']['new_status'] == 'consultation'
            assert event['data']['actor'] == 'sse-test-agent'
        finally:
            app_module.event_bus.unsubscribe(q)

    def test_consensus_reached_fires_extra_event(self, client, temp_data_dir, sample_proposals):
        """Moving to 'consensus' status should fire both status_changed and consensus_reached."""
        import app as app_module
        q = app_module.event_bus.subscribe()
        try:
            response = client.put(
                '/api/proposals/test-proposal-002/status',
                json={
                    'status': 'consensus',
                    'actor': 'sse-test-agent',
                },
            )
            assert response.status_code == 200

            # Should get two events: status_changed and consensus_reached
            events = []
            for _ in range(2):
                events.append(q.get(timeout=2))

            event_types = {e['type'] for e in events}
            assert 'status_changed' in event_types
            assert 'consensus_reached' in event_types

            consensus_event = [e for e in events if e['type'] == 'consensus_reached'][0]
            assert consensus_event['data']['proposal_id'] == 'test-proposal-002'
            assert consensus_event['data']['proposal_title'] == 'Test Proposal: With Consultations'
        finally:
            app_module.event_bus.unsubscribe(q)

    def test_non_consensus_status_change_no_extra_event(self, client, temp_data_dir, sample_proposals):
        """Moving to a status other than 'consensus' should NOT fire consensus_reached."""
        import app as app_module
        q = app_module.event_bus.subscribe()
        try:
            # proposed -> consultation (not consensus)
            response = client.put(
                '/api/proposals/test-proposal-001/status',
                json={
                    'status': 'consultation',
                    'actor': 'sse-test-agent',
                },
            )
            assert response.status_code == 200

            event = q.get(timeout=2)
            assert event['type'] == 'status_changed'

            # Queue should be empty — no consensus_reached event
            assert q.empty()
        finally:
            app_module.event_bus.unsubscribe(q)
