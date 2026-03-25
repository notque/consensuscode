#!/usr/bin/env python3
"""
CollectiveFlow Web Interface

A horizontal, non-hierarchical web interface for the CollectiveFlow consensus system.
This Flask application provides accessible views of proposals and consultations
without authentication or special roles - embodying true collective principles.
"""

import os
import re
import json
import logging
import secrets
import yaml
import uuid
import queue
import threading
import time
from datetime import datetime
from pathlib import Path
from flask import Flask, render_template, jsonify, request, redirect, url_for, flash, session, abort, Response
from storage import get_storage

app = Flask(__name__)

# SECRET_KEY: use environment variable, or generate a random one per process for dev.
# In production, always set SECRET_KEY in the environment.
app.secret_key = os.environ.get('SECRET_KEY') or secrets.token_hex(32)

# CORS: restrict to same-origin by default. If cross-origin API access is needed,
# configure allowed origins explicitly via CORS_ORIGINS environment variable.
cors_origins = os.environ.get('CORS_ORIGINS', '').strip()
if cors_origins:
    from flask_cors import CORS
    CORS(app, origins=cors_origins.split(','))

# ---------- SSE Event Bus ----------
# Teaching note (50/50 teaching/doing):
#
# Server-Sent Events (SSE) is a W3C standard where the server pushes events
# to clients over a long-lived HTTP connection. Unlike WebSockets, SSE is:
#   - Unidirectional (server -> client only), which is perfect for notifications
#   - Built on plain HTTP, so it works through proxies and with curl
#   - Auto-reconnects (the browser's EventSource API handles this natively)
#   - Zero extra dependencies in Flask — just a streaming Response
#
# The pattern below uses an in-memory pub/sub bus. Each SSE client gets its
# own queue. When an event fires, it's pushed to every connected queue.
# This is intentionally simple — no Redis, no external broker. It works for
# a local-only collective where all agents connect to the same process.
#
# Limitation: events are not persisted. If an agent disconnects and
# reconnects, it won't receive events that fired while it was away.
# For catch-up, agents should poll GET /api/proposals on reconnect.
# The Last-Event-Id header support below helps with brief disconnections.


class EventBus:
    """In-memory pub/sub for SSE event distribution.

    Thread-safe. Each subscriber gets an independent queue so slow consumers
    don't block fast ones. Queues have a max size to bound memory — if a
    consumer falls behind, oldest events are silently dropped.
    """

    def __init__(self, maxsize=128):
        self._lock = threading.Lock()
        self._subscribers: list[queue.Queue] = []
        self._event_id = 0
        self._maxsize = maxsize

    def subscribe(self) -> queue.Queue:
        """Register a new subscriber. Returns a queue to read events from."""
        q: queue.Queue = queue.Queue(maxsize=self._maxsize)
        with self._lock:
            self._subscribers.append(q)
        return q

    def unsubscribe(self, q: queue.Queue):
        """Remove a subscriber (call when client disconnects)."""
        with self._lock:
            try:
                self._subscribers.remove(q)
            except ValueError:
                pass

    def publish(self, event_type: str, data: dict):
        """Broadcast an event to all connected subscribers.

        Teaching note: We increment a monotonic event ID so clients can use
        the Last-Event-Id header to detect missed events after reconnection.
        The SSE spec says: if the client reconnects, the browser sends
        Last-Event-Id automatically. We don't do full replay here (that
        would need persistent storage), but the ID lets clients know they
        missed something and should do a full poll.
        """
        with self._lock:
            self._event_id += 1
            event = {
                'id': self._event_id,
                'type': event_type,
                'data': data,
                'timestamp': datetime.now().isoformat(),
            }
            dead_queues = []
            for q in self._subscribers:
                try:
                    q.put_nowait(event)
                except queue.Full:
                    # Consumer is too slow — drop oldest event to make room
                    try:
                        q.get_nowait()
                        q.put_nowait(event)
                    except (queue.Empty, queue.Full):
                        dead_queues.append(q)
            # Clean up any broken queues
            for q in dead_queues:
                try:
                    self._subscribers.remove(q)
                except ValueError:
                    pass


# Single global event bus — shared across all request threads.
event_bus = EventBus()


# ---------- CSRF Protection ----------
# Lightweight token-based CSRF using Flask sessions. No extra dependency needed.

def generate_csrf_token():
    """Generate or retrieve a CSRF token for the current session."""
    if '_csrf_token' not in session:
        session['_csrf_token'] = secrets.token_hex(32)
    return session['_csrf_token']

# Make csrf_token() available in all templates
app.jinja_env.globals['csrf_token'] = generate_csrf_token

@app.before_request
def csrf_protect():
    """Reject POST/PUT/DELETE requests that lack a valid CSRF token."""
    if request.method in ('POST', 'PUT', 'DELETE'):
        # Skip CSRF check for JSON API endpoints (they use CORS + Accept headers)
        if request.path.startswith('/api/'):
            return
        token = session.get('_csrf_token', None)
        form_token = request.form.get('_csrf_token', None)
        if not token or token != form_token:
            abort(403)


# ---------- Configuration ----------
DATA_DIR = os.environ.get('COLLECTIVEFLOW_DATA', '../data')
PROPOSALS_DIR = Path(DATA_DIR) / 'proposals'

# Storage abstraction: supports YAML (default) or SQLite backends.
# Set STORAGE_BACKEND=sqlite to use SQLite; defaults to YAML for
# backwards compatibility.  The storage object satisfies the
# StorageBackend protocol from storage.py.
storage = get_storage(DATA_DIR)


def load_proposals():
    """Load all proposals via the configured storage backend."""
    return storage.load_proposals()


def get_proposal(proposal_id):
    """Load a specific proposal by ID via the configured storage backend."""
    return storage.get_proposal(proposal_id)


VALID_STATUSES = ['proposed', 'consultation', 'consensus', 'implemented', 'blocked', 'withdrawn']
VALID_URGENCIES = ['low', 'medium', 'high', 'emergency']

# Valid status transitions — each status maps to the statuses it can move to.
# This prevents illogical jumps (e.g., going from 'withdrawn' to 'implemented').
STATUS_TRANSITIONS = {
    'proposed': ['consultation', 'withdrawn'],
    'consultation': ['consensus', 'blocked', 'withdrawn'],
    'consensus': ['implemented', 'blocked', 'withdrawn'],
    'implemented': [],
    'blocked': ['consultation', 'withdrawn'],
    'withdrawn': [],
}


def api_error(message, code, http_status=400):
    """Return a consistent JSON error response.

    Teaching note: Every API should have a single, predictable error shape.
    Clients can always check for the "error" key to know something went wrong,
    and use the "code" key for programmatic handling (e.g., retry logic).
    """
    return jsonify({'error': message, 'code': code}), http_status


def save_proposal(proposal_data):
    """Save a new proposal via the configured storage backend."""
    return storage.save_proposal(proposal_data)


def update_proposal(proposal_id, proposal_data):
    """Update an existing proposal on disk.

    Teaching note: Separating create vs. update keeps each function's
    responsibility clear.  save_proposal() adds default metadata for new
    proposals; update_proposal() writes the already-complete dict back
    without adding extra fields.
    """
    yaml_path = PROPOSALS_DIR / f"{proposal_id}.yaml"
    with open(yaml_path, 'w') as f:
        yaml.safe_dump(proposal_data, f, default_flow_style=False, sort_keys=False)

    json_path = PROPOSALS_DIR / f"{proposal_id}.json"
    with open(json_path, 'w') as f:
        json.dump(proposal_data, f, indent=2, default=str)

@app.route('/')
def index():
    """Home page showing all proposals."""
    proposals = load_proposals()
    
    # Group proposals by status for better organization
    grouped = {
        'consultation': [],
        'proposed': [],
        'consensus': [],
        'implemented': [],
        'blocked': [],
        'withdrawn': []
    }
    
    for proposal in proposals:
        status = proposal.get('status', 'proposed')
        if status in grouped:
            grouped[status].append(proposal)
    
    return render_template('index.html', grouped_proposals=grouped)

@app.route('/proposal/<proposal_id>')
def proposal_detail(proposal_id):
    """Detailed view of a specific proposal."""
    proposal = get_proposal(proposal_id)
    
    if not proposal:
        return "Proposal not found", 404
    
    return render_template('proposal.html', proposal=proposal)

@app.route('/api/proposals', methods=['GET'])
def api_proposals():
    """API endpoint for proposals list.

    Supports optional query-string filters:
        ?status=consultation   — filter by status
        ?urgency=high          — filter by urgency

    Teaching note: GET endpoints should be safe (no side-effects) and
    idempotent. Filtering via query parameters keeps the URL clean and
    cache-friendly.
    """
    proposals = load_proposals()

    # Optional filtering
    status_filter = request.args.get('status')
    if status_filter:
        proposals = [p for p in proposals if p.get('status') == status_filter]

    urgency_filter = request.args.get('urgency')
    if urgency_filter:
        proposals = [p for p in proposals if p.get('urgency') == urgency_filter]

    return jsonify({
        'proposals': proposals,
        'count': len(proposals)
    })


@app.route('/api/proposals/<proposal_id>', methods=['GET'])
def api_proposal_by_id(proposal_id):
    """API endpoint for a specific proposal.

    Teaching note: We use /api/proposals/<id> (plural resource with ID)
    rather than /api/proposal/<id>.  The plural form is the REST convention
    because the ID selects one item *from the collection*.  The old
    /api/proposal/<id> route is kept below for backwards compatibility.
    """
    proposal = get_proposal(proposal_id)

    if not proposal:
        return api_error('Proposal not found', 'NOT_FOUND', 404)

    return jsonify(proposal)


# Backwards-compatible alias — existing clients may use the singular form.
@app.route('/api/proposal/<proposal_id>')
def api_proposal(proposal_id):
    """Legacy endpoint — redirects internally to the canonical route."""
    return api_proposal_by_id(proposal_id)


@app.route('/api/proposals', methods=['POST'])
def api_create_proposal():
    """Create a proposal via JSON API.

    Teaching note: POST to the collection URL is the REST way to create a
    new resource.  We return 201 (Created) with a Location header pointing
    to the new resource — this lets clients find the resource without
    parsing the body.

    Expected JSON body:
        {
            "title": "string (required)",
            "description": "string (required)",
            "proposer": "string (optional, defaults to 'api-user')",
            "urgency": "low|medium|high|emergency (optional, defaults to 'medium')",
            "affected_areas": ["string"] (optional)
        }
    """
    # Require JSON content type
    if not request.is_json:
        return api_error(
            'Request must be JSON (set Content-Type: application/json)',
            'INVALID_CONTENT_TYPE',
            415
        )

    data = request.get_json()

    # --- Input validation ---
    errors = []
    title = (data.get('title') or '').strip()
    description = (data.get('description') or '').strip()

    if not title:
        errors.append('title is required')
    if not description:
        errors.append('description is required')

    urgency = data.get('urgency', 'medium')
    if urgency not in VALID_URGENCIES:
        errors.append(f'urgency must be one of: {", ".join(VALID_URGENCIES)}')

    affected_areas = data.get('affected_areas', [])
    if not isinstance(affected_areas, list):
        errors.append('affected_areas must be a list')

    if errors:
        return api_error(
            '; '.join(errors),
            'VALIDATION_ERROR',
            422
        )

    proposal_data = {
        'title': title,
        'description': description,
        'proposer': (data.get('proposer') or 'api-user').strip(),
        'urgency': urgency,
        'affected_areas': affected_areas,
    }

    proposal_id = save_proposal(proposal_data)
    proposal = get_proposal(proposal_id)

    # Notify SSE subscribers
    event_bus.publish('proposal_created', {
        'proposal_id': proposal_id,
        'title': title,
        'proposer': proposal_data['proposer'],
        'urgency': urgency,
        'affected_areas': affected_areas,
    })

    response = jsonify(proposal)
    response.status_code = 201
    response.headers['Location'] = f'/api/proposals/{proposal_id}'
    return response


@app.route('/api/proposals/<proposal_id>/consultation', methods=['POST'])
def api_add_consultation(proposal_id):
    """Add consultation input to a proposal.

    Teaching note: Consultation is a *sub-resource* of a proposal.  POSTing
    to /proposals/<id>/consultation adds a new entry to the consultations
    list.  This mirrors how the collective works: any agent can contribute
    input at any time, and every contribution is recorded.

    Expected JSON body:
        {
            "contributor": "string (required)",
            "input": "string (required)",
            "support": true|false (optional),
            "concerns": ["string"] (optional)
        }
    """
    if not request.is_json:
        return api_error(
            'Request must be JSON (set Content-Type: application/json)',
            'INVALID_CONTENT_TYPE',
            415
        )

    proposal = get_proposal(proposal_id)
    if not proposal:
        return api_error('Proposal not found', 'NOT_FOUND', 404)

    # Cannot add consultation to terminal statuses
    if proposal.get('status') in ('implemented', 'withdrawn'):
        return api_error(
            f'Cannot add consultation to a proposal with status "{proposal["status"]}"',
            'INVALID_STATE',
            409
        )

    data = request.get_json()

    # --- Input validation ---
    errors = []
    contributor = (data.get('contributor') or '').strip()
    input_text = (data.get('input') or '').strip()

    if not contributor:
        errors.append('contributor is required')
    if not input_text:
        errors.append('input is required')

    concerns = data.get('concerns', [])
    if not isinstance(concerns, list):
        errors.append('concerns must be a list of strings')

    if errors:
        return api_error('; '.join(errors), 'VALIDATION_ERROR', 422)

    now = datetime.now().isoformat()

    consultation_entry = {
        'contributor': contributor,
        'timestamp': now,
        'input': input_text,
        'support': bool(data.get('support', True)),
    }
    if concerns:
        consultation_entry['concerns'] = concerns

    # Append to proposal
    if 'consultations' not in proposal:
        proposal['consultations'] = []
    proposal['consultations'].append(consultation_entry)

    # Record in consensus history
    if 'consensus_history' not in proposal:
        proposal['consensus_history'] = []
    proposal['consensus_history'].append({
        'timestamp': now,
        'event': 'consultation_added',
        'actor': contributor,
        'details': f'{"Support" if consultation_entry["support"] else "Concern"}: {input_text[:80]}'
    })

    update_proposal(proposal_id, proposal)

    # Notify SSE subscribers
    event_bus.publish('consultation_added', {
        'proposal_id': proposal_id,
        'proposal_title': proposal.get('title', ''),
        'contributor': contributor,
        'support': consultation_entry['support'],
        'consultation_count': len(proposal.get('consultations', [])),
    })

    return jsonify({
        'message': 'Consultation added',
        'consultation': consultation_entry,
        'proposal_id': proposal_id,
    }), 201


@app.route('/api/proposals/<proposal_id>/status', methods=['PUT'])
def api_update_status(proposal_id):
    """Advance or change the status of a proposal.

    Teaching note: PUT to a specific sub-resource (/status) is appropriate
    because we are *replacing* the status value.  We enforce valid state
    transitions so proposals follow the collective's process.

    Valid transitions:
        proposed     -> consultation, withdrawn
        consultation -> consensus, blocked, withdrawn
        consensus    -> implemented, blocked, withdrawn
        blocked      -> consultation, withdrawn
        implemented  -> (terminal)
        withdrawn    -> (terminal)

    Expected JSON body:
        {
            "status": "string (required — target status)",
            "actor": "string (required — who is making this change)",
            "reason": "string (optional — rationale for the change)"
        }
    """
    if not request.is_json:
        return api_error(
            'Request must be JSON (set Content-Type: application/json)',
            'INVALID_CONTENT_TYPE',
            415
        )

    proposal = get_proposal(proposal_id)
    if not proposal:
        return api_error('Proposal not found', 'NOT_FOUND', 404)

    data = request.get_json()

    # --- Input validation ---
    errors = []
    new_status = (data.get('status') or '').strip()
    actor = (data.get('actor') or '').strip()

    if not new_status:
        errors.append('status is required')
    elif new_status not in VALID_STATUSES:
        errors.append(f'status must be one of: {", ".join(VALID_STATUSES)}')

    if not actor:
        errors.append('actor is required')

    if errors:
        return api_error('; '.join(errors), 'VALIDATION_ERROR', 422)

    current_status = proposal.get('status', 'proposed')

    # Check valid transition
    allowed = STATUS_TRANSITIONS.get(current_status, [])
    if new_status not in allowed:
        return api_error(
            f'Cannot transition from "{current_status}" to "{new_status}". '
            f'Allowed transitions: {", ".join(allowed) if allowed else "none (terminal status)"}',
            'INVALID_TRANSITION',
            409
        )

    now = datetime.now().isoformat()
    reason = (data.get('reason') or '').strip()

    proposal['status'] = new_status
    proposal['consensus_status'] = f'Status changed to {new_status}'

    if 'consensus_history' not in proposal:
        proposal['consensus_history'] = []
    history_entry = {
        'timestamp': now,
        'event': 'status_changed',
        'actor': actor,
        'details': f'Status changed from {current_status} to {new_status}'
    }
    if reason:
        history_entry['details'] += f' — {reason}'
    proposal['consensus_history'].append(history_entry)

    update_proposal(proposal_id, proposal)

    # Notify SSE subscribers — status_changed for all transitions
    event_bus.publish('status_changed', {
        'proposal_id': proposal_id,
        'proposal_title': proposal.get('title', ''),
        'previous_status': current_status,
        'new_status': new_status,
        'actor': actor,
        'reason': reason,
    })

    # Additional consensus_reached event when a proposal reaches consensus.
    # Teaching note: this is a semantic event on top of the mechanical
    # status_changed event. Agents that only care about "did we agree?" can
    # filter to just consensus_reached without parsing status fields.
    if new_status == 'consensus':
        event_bus.publish('consensus_reached', {
            'proposal_id': proposal_id,
            'proposal_title': proposal.get('title', ''),
            'actor': actor,
            'consultation_count': len(proposal.get('consultations', [])),
        })

    return jsonify({
        'message': f'Status updated to {new_status}',
        'proposal_id': proposal_id,
        'previous_status': current_status,
        'new_status': new_status,
    })


@app.route('/api/collective/stats', methods=['GET'])
def api_collective_stats():
    """Collective statistics — a read-only snapshot of activity.

    Teaching note: Stats endpoints are great candidates for caching headers
    in production.  Here we keep it simple: compute on every request.  The
    response shape is stable so clients can depend on it.
    """
    proposals = load_proposals()

    status_counts = {}
    for status in VALID_STATUSES:
        status_counts[status] = 0
    for p in proposals:
        s = p.get('status', 'proposed')
        if s in status_counts:
            status_counts[s] += 1

    contributors = set()
    total_consultations = 0
    for p in proposals:
        contributors.add(p.get('proposer', 'anonymous'))
        consultations = p.get('consultations', [])
        total_consultations += len(consultations)
        for c in consultations:
            contributors.add(c.get('contributor', 'anonymous'))

    return jsonify({
        'total_proposals': len(proposals),
        'status_counts': status_counts,
        'total_consultations': total_consultations,
        'contributor_count': len(contributors),
        'contributors': sorted(contributors),
    })


# ---------- SSE Streaming Endpoint ----------

@app.route('/api/events')
def api_events():
    """Server-Sent Events stream for real-time proposal notifications.

    Teaching note (50/50 teaching/doing):

    SSE uses a simple text protocol over HTTP. Each event looks like:

        id: 42
        event: proposal_created
        data: {"proposal_id": "proposal-2026-03-24-abc123", ...}

        (blank line ends the event)

    Clients connect with:
        - Browser: new EventSource('/api/events')
        - curl:    curl -N http://localhost:5000/api/events

    The endpoint supports optional query parameters for filtering:
        ?types=proposal_created,status_changed  — only receive these event types
        ?proposal_id=proposal-2026-03-24-abc    — only events for this proposal

    No authentication required — horizontal principle. Any agent or tool can
    subscribe. This directly addresses the consensus assessment finding that
    agents miss proposals because there's no notification mechanism.

    Implementation: Flask's Response with a generator function creates the
    long-lived connection. We subscribe to the global event_bus and yield
    SSE-formatted strings as events arrive. When the client disconnects,
    the generator's finally block cleans up the subscription.
    """
    # Parse optional filters from query string
    type_filter = set()
    types_param = request.args.get('types', '').strip()
    if types_param:
        type_filter = set(t.strip() for t in types_param.split(',') if t.strip())

    proposal_filter = request.args.get('proposal_id', '').strip() or None

    def stream():
        q = event_bus.subscribe()
        try:
            # Send an initial comment so the client knows the connection is alive.
            # SSE spec: lines starting with ':' are comments, ignored by EventSource
            # but useful for keeping proxies from closing idle connections.
            yield ': connected to CollectiveFlow event stream\n\n'

            while True:
                try:
                    # Block for up to 30 seconds, then send a keepalive comment.
                    # Teaching note: without periodic data, reverse proxies and
                    # load balancers may close "idle" connections (common timeout
                    # is 60s). A 30s heartbeat prevents that.
                    event = q.get(timeout=30)
                except queue.Empty:
                    # No event within 30s — send keepalive
                    yield ': keepalive\n\n'
                    continue

                # Apply filters
                if type_filter and event['type'] not in type_filter:
                    continue
                if proposal_filter:
                    event_proposal = event.get('data', {}).get('proposal_id', '')
                    if event_proposal != proposal_filter:
                        continue

                # Format as SSE
                yield f"id: {event['id']}\n"
                yield f"event: {event['type']}\n"
                yield f"data: {json.dumps(event['data'], default=str)}\n"
                yield '\n'  # Blank line terminates the event

        except GeneratorExit:
            # Client disconnected — clean up
            pass
        finally:
            event_bus.unsubscribe(q)

    return Response(
        stream(),
        mimetype='text/event-stream',
        headers={
            'Cache-Control': 'no-cache',
            'X-Accel-Buffering': 'no',  # Disable nginx buffering if present
            'Connection': 'keep-alive',
        },
    )


@app.route('/proposals')
def proposals_list():
    """View showing all proposals grouped by status."""
    proposals = load_proposals()

    # Group proposals by status
    grouped = {
        'consultation': [],
        'proposed': [],
        'consensus': [],
        'implemented': [],
        'blocked': [],
        'withdrawn': []
    }

    for proposal in proposals:
        status = proposal.get('status', 'proposed')
        if status in grouped:
            grouped[status].append(proposal)

    return render_template('proposals.html', grouped_proposals=grouped)

@app.route('/about')
def about():
    """View showing information about the collective and its principles."""
    return render_template('about.html')

@app.route('/dashboard')
def dashboard():
    """Dashboard showing collective-wide statistics, velocity, and participation."""
    proposals = load_proposals()

    # Basic stats
    status_counts = {}
    urgency_counts = {}
    for proposal in proposals:
        status = proposal.get('status', 'proposed')
        status_counts[status] = status_counts.get(status, 0) + 1
        urgency = proposal.get('urgency', 'medium')
        urgency_counts[urgency] = urgency_counts.get(urgency, 0) + 1

    proposed_count = status_counts.get('proposed', 0)
    consultation_count = status_counts.get('consultation', 0)
    consensus_count = status_counts.get('consensus', 0)
    implemented_count = status_counts.get('implemented', 0)

    # Collect all unique contributors
    contributors = set()
    for proposal in proposals:
        contributors.add(proposal.get('proposer', 'anonymous'))
        for consultation in proposal.get('consultations', []):
            contributors.add(consultation.get('contributor', 'anonymous'))

    stats = {
        'total_proposals': len(proposals),
        'active_count': proposed_count + consultation_count,
        'consensus_count': consensus_count,
        'implemented_count': implemented_count,
        'contributor_count': len(contributors),
    }

    # Velocity metrics
    total_history_events = 0
    proposals_with_history = 0
    total_consultations = 0
    for proposal in proposals:
        history = proposal.get('consensus_history', [])
        if history:
            proposals_with_history += 1
            total_history_events += len(history)
        total_consultations += len(proposal.get('consultations', []))

    velocity = {
        'proposals_with_history': proposals_with_history,
        'avg_events_per_proposal': round(total_history_events / proposals_with_history, 1) if proposals_with_history > 0 else 0,
        'avg_consultations': round(total_consultations / len(proposals), 1) if proposals else 0,
    }

    # Agent participation rates
    agent_contributions = {}
    for proposal in proposals:
        for consultation in proposal.get('consultations', []):
            contributor = consultation.get('contributor', 'anonymous')
            agent_contributions[contributor] = agent_contributions.get(contributor, 0) + 1

    max_contributions = max(agent_contributions.values()) if agent_contributions else 1
    participation = sorted(
        [
            {
                'name': name,
                'count': count,
                'percentage': (count / max_contributions * 100) if max_contributions > 0 else 0
            }
            for name, count in agent_contributions.items()
        ],
        key=lambda x: x['count'],
        reverse=True
    )

    # Project status overview (based on affected_areas)
    project_map = {}
    for proposal in proposals:
        areas = proposal.get('affected_areas', [])
        status = proposal.get('status', 'proposed')
        for area in areas:
            if area not in project_map:
                project_map[area] = {'name': area, 'proposal_count': 0, 'statuses': []}
            project_map[area]['proposal_count'] += 1
            if status not in project_map[area]['statuses']:
                project_map[area]['statuses'].append(status)

    projects = sorted(project_map.values(), key=lambda x: x['proposal_count'], reverse=True)

    # Recent activity (from consensus_history across all proposals)
    recent_events = []
    for proposal in proposals:
        for event in proposal.get('consensus_history', []):
            recent_events.append({
                'type': event.get('event', 'update'),
                'proposal_id': proposal.get('id'),
                'proposal_title': proposal.get('title'),
                'timestamp': event.get('timestamp'),
                'actor': event.get('actor'),
            })
    # Sort by timestamp descending, take last 10
    recent_events.sort(key=lambda e: str(e.get('timestamp', '')), reverse=True)
    recent_events = recent_events[:10]

    return render_template(
        'dashboard.html',
        stats=stats,
        status_counts=status_counts,
        urgency_counts=urgency_counts,
        velocity=velocity,
        participation=participation,
        projects=projects,
        recent_events=recent_events,
    )

@app.route('/collective')
def collective_view():
    """View showing the collective's current state and activity."""
    proposals = load_proposals()
    
    # Calculate collective statistics
    stats = {
        'total_proposals': len(proposals),
        'active_consultations': len([p for p in proposals if p.get('status') == 'consultation']),
        'implemented': len([p for p in proposals if p.get('status') == 'implemented']),
        'contributors': set()
    }
    
    # Collect all unique contributors
    for proposal in proposals:
        stats['contributors'].add(proposal.get('proposer', 'anonymous'))
        for consultation in proposal.get('consultations', []):
            stats['contributors'].add(consultation.get('contributor', 'anonymous'))
    
    stats['contributors'] = list(stats['contributors'])
    stats['contributor_count'] = len(stats['contributors'])
    
    # Recent activity
    recent_events = []
    for proposal in proposals[:5]:  # Last 5 proposals
        event = {
            'type': 'proposal_created',
            'proposal_id': proposal.get('id'),
            'proposal_title': proposal.get('title'),
            'timestamp': proposal.get('date'),
            'actor': proposal.get('proposer')
        }
        recent_events.append(event)
    
    return render_template('collective.html', stats=stats, recent_events=recent_events)

@app.route('/create')
def create_proposal_form():
    """Show proposal creation form."""
    return render_template('create_proposal.html')

@app.route('/create', methods=['POST'])
def create_proposal():
    """Handle proposal creation."""
    VALID_URGENCIES = {'low', 'medium', 'high', 'emergency'}
    VALID_AREAS = {
        'infrastructure', 'web-interface', 'cli-tool', 'consensus-process',
        'documentation', 'external-communication', 'agent-coordination',
        'testing', 'other'
    }

    try:
        # Collect form data
        raw_urgency = request.form.get('urgency', 'medium')
        raw_areas = request.form.getlist('affected_areas')

        proposal_data = {
            'title': request.form.get('title', '').strip()[:200],
            'description': request.form.get('description', '').strip()[:5000],
            'proposer': request.form.get('proposer', 'anonymous').strip()[:100],
            'urgency': raw_urgency if raw_urgency in VALID_URGENCIES else 'medium',
            'affected_areas': [a for a in raw_areas if a in VALID_AREAS]
        }

        # Basic validation
        if not proposal_data['title']:
            flash('Title is required', 'error')
            return redirect(url_for('create_proposal_form'))

        if not proposal_data['description']:
            flash('Description is required', 'error')
            return redirect(url_for('create_proposal_form'))

        # Save proposal
        proposal_id = save_proposal(proposal_data)

        # Notify SSE subscribers (same event shape as the API route)
        event_bus.publish('proposal_created', {
            'proposal_id': proposal_id,
            'title': proposal_data['title'],
            'proposer': proposal_data['proposer'],
            'urgency': proposal_data['urgency'],
            'affected_areas': proposal_data['affected_areas'],
        })

        flash('Proposal submitted successfully!', 'success')
        return redirect(url_for('proposal_detail', proposal_id=proposal_id))

    except Exception as e:
        # Log the real error server-side; show generic message to user
        logging.exception("Error creating proposal")
        flash('Error creating proposal. Please try again.', 'error')
        return redirect(url_for('create_proposal_form'))

@app.route('/proposal/<proposal_id>/consult', methods=['POST'])
def add_consultation(proposal_id):
    """Handle consultation input submission for a proposal."""
    proposal = get_proposal(proposal_id)

    if not proposal:
        return "Proposal not found", 404

    try:
        contributor = request.form.get('contributor', '').strip() or 'anonymous'
        position = request.form.get('position', 'support')
        comment = request.form.get('comment', '').strip()
        concerns = request.form.get('concerns', '').strip()
        suggestions = request.form.get('suggestions', '').strip()

        if not comment:
            flash('Comment is required — your reasoning matters to the collective.', 'error')
            return redirect(url_for('proposal_detail', proposal_id=proposal_id) + '#consultation-form')

        if position not in ('support', 'support-with-concerns', 'block'):
            flash('Invalid position selected.', 'error')
            return redirect(url_for('proposal_detail', proposal_id=proposal_id) + '#consultation-form')

        consultation = {
            'contributor': contributor,
            'position': position,
            'support': position in ('support', 'support-with-concerns'),
            'input': comment,
            'timestamp': datetime.now().isoformat(),
        }

        if concerns:
            consultation['concerns'] = [c.strip() for c in concerns.split('\n') if c.strip()]
        if suggestions:
            consultation['suggestions'] = [s.strip() for s in suggestions.split('\n') if s.strip()]

        if 'consultations' not in proposal:
            proposal['consultations'] = []
        proposal['consultations'].append(consultation)

        if 'consensus_history' not in proposal:
            proposal['consensus_history'] = []
        proposal['consensus_history'].append({
            'timestamp': consultation['timestamp'],
            'event': 'consultation_added',
            'actor': contributor,
            'details': f"Position: {position}"
        })

        yaml_path = PROPOSALS_DIR / f"{proposal_id}.yaml"
        with open(yaml_path, 'w') as f:
            yaml.safe_dump(proposal, f, default_flow_style=False, sort_keys=False)

        json_path = PROPOSALS_DIR / f"{proposal_id}.json"
        if json_path.exists():
            with open(json_path, 'w') as f:
                json.dump(proposal, f, indent=2, default=str)

        # Notify SSE subscribers (same event shape as the API route)
        event_bus.publish('consultation_added', {
            'proposal_id': proposal_id,
            'proposal_title': proposal.get('title', ''),
            'contributor': contributor,
            'support': consultation['support'],
            'consultation_count': len(proposal.get('consultations', [])),
        })

        flash(f'Consultation input from "{contributor}" recorded. Thank you for participating!', 'success')
        return redirect(url_for('proposal_detail', proposal_id=proposal_id) + '#consultations-heading')

    except Exception as e:
        flash(f'Error submitting consultation: {str(e)}', 'error')
        return redirect(url_for('proposal_detail', proposal_id=proposal_id))

@app.after_request
def set_security_headers(response):
    """Add security headers to every response."""
    response.headers['X-Content-Type-Options'] = 'nosniff'
    response.headers['X-Frame-Options'] = 'DENY'
    response.headers['X-XSS-Protection'] = '1; mode=block'
    response.headers['Referrer-Policy'] = 'strict-origin-when-cross-origin'
    response.headers['Permissions-Policy'] = 'geolocation=(), camera=(), microphone=()'
    response.headers['Content-Security-Policy'] = (
        "default-src 'self'; "
        "script-src 'self' 'unsafe-inline' https://cdn.tailwindcss.com; "
        "style-src 'self' 'unsafe-inline'; "
        "img-src 'self' data:; "
        "font-src 'self'; "
        "connect-src 'self'; "
        "frame-ancestors 'none'"
    )
    return response


@app.template_filter('humanize_date')
def humanize_date(date_str):
    """Convert ISO date string to human-readable format."""
    try:
        if isinstance(date_str, str):
            # Parse ISO format
            dt = datetime.fromisoformat(date_str.replace('Z', '+00:00'))
        else:
            dt = date_str
        
        # Format nicely
        return dt.strftime('%B %d, %Y at %I:%M %p')
    except:
        return date_str

@app.template_filter('status_emoji')
def status_emoji(status):
    """Return an emoji representing the proposal status."""
    emoji_map = {
        'proposed': '💡',
        'consultation': '🗣️',
        'consensus': '🤝',
        'implemented': '✅',
        'blocked': '🚫',
        'withdrawn': '↩️'
    }
    return emoji_map.get(status, '📄')

@app.template_filter('urgency_color')
def urgency_color(urgency):
    """Return a CSS class for urgency level."""
    color_map = {
        'low': 'text-green-600',
        'medium': 'text-yellow-600',
        'high': 'text-orange-600',
        'emergency': 'text-red-600'
    }
    if not isinstance(urgency, str):
        return 'text-gray-600'
    return color_map.get(urgency, 'text-gray-600')

if __name__ == '__main__':
    # Debug mode from environment only — never hardcode True in source
    debug_mode = os.environ.get('FLASK_DEBUG', '0').lower() in ('1', 'true', 'yes')
    app.run(debug=debug_mode, host='0.0.0.0', port=5000)