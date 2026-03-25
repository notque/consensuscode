#!/usr/bin/env python3
"""
CollectiveFlow Web Interface

A horizontal, non-hierarchical web interface for the CollectiveFlow consensus system.
This Flask application provides accessible views of proposals and consultations
without authentication or special roles - embodying true collective principles.
"""

import os
import json
import yaml
import uuid
from datetime import datetime
from pathlib import Path
from flask import Flask, render_template, jsonify, request, redirect, url_for, flash
from flask_cors import CORS

app = Flask(__name__)
app.secret_key = os.environ.get('SECRET_KEY', 'collective-flow-dev-key')  # For flash messages
CORS(app)  # Enable cross-origin requests for API compatibility

# Configuration from environment or defaults
DATA_DIR = os.environ.get('COLLECTIVEFLOW_DATA', '../data')
PROPOSALS_DIR = Path(DATA_DIR) / 'proposals'

def load_proposals():
    """Load all proposals from the data directory."""
    proposals = []
    
    if not PROPOSALS_DIR.exists():
        return proposals
    
    for yaml_file in PROPOSALS_DIR.glob('*.yaml'):
        try:
            with open(yaml_file, 'r') as f:
                proposal = yaml.safe_load(f)
                if proposal:
                    proposals.append(proposal)
        except Exception as e:
            print(f"Error loading {yaml_file}: {e}")
    
    # Sort by date, newest first
    # Normalize dates to strings for consistent comparison (YAML may parse dates as datetime objects)
    def sort_key(p):
        date = p.get('date', '')
        if isinstance(date, datetime):
            return date.isoformat()
        return str(date)

    proposals.sort(key=sort_key, reverse=True)
    return proposals

def get_proposal(proposal_id):
    """Load a specific proposal by ID."""
    try:
        yaml_path = PROPOSALS_DIR / f"{proposal_id}.yaml"

        if yaml_path.exists():
            with open(yaml_path, 'r') as f:
                return yaml.safe_load(f)
    except OSError:
        # Handles edge cases like excessively long IDs that exceed
        # filesystem path limits.
        pass

    return None

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
    """Save a new proposal to the data directory."""
    # Ensure proposals directory exists
    PROPOSALS_DIR.mkdir(parents=True, exist_ok=True)

    # Generate unique ID if not provided
    if 'id' not in proposal_data:
        proposal_data['id'] = f"proposal-{datetime.now().strftime('%Y-%m-%d')}-{str(uuid.uuid4())[:8]}"

    # Add metadata
    proposal_data['date'] = datetime.now().isoformat()
    proposal_data['status'] = 'proposed'
    proposal_data['consensus_status'] = 'New proposal submitted'
    proposal_data['consensus_history'] = [{
        'timestamp': proposal_data['date'],
        'event': 'proposal_created',
        'actor': proposal_data.get('proposer', 'web-user'),
        'details': f"Created with urgency: {proposal_data.get('urgency', 'medium')}"
    }]
    proposal_data['consultations'] = []

    # Save to YAML file
    yaml_path = PROPOSALS_DIR / f"{proposal_data['id']}.yaml"
    with open(yaml_path, 'w') as f:
        yaml.safe_dump(proposal_data, f, default_flow_style=False, sort_keys=False)

    # Also save JSON for API compatibility
    json_path = PROPOSALS_DIR / f"{proposal_data['id']}.json"
    with open(json_path, 'w') as f:
        json.dump(proposal_data, f, indent=2)

    return proposal_data['id']


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
    try:
        # Collect form data
        proposal_data = {
            'title': request.form.get('title', '').strip(),
            'description': request.form.get('description', '').strip(),
            'proposer': request.form.get('proposer', 'anonymous').strip(),
            'urgency': request.form.get('urgency', 'medium'),
            'affected_areas': request.form.getlist('affected_areas')
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
        
        flash(f'Proposal "{proposal_data["title"]}" submitted successfully!', 'success')
        return redirect(url_for('proposal_detail', proposal_id=proposal_id))
        
    except Exception as e:
        flash(f'Error creating proposal: {str(e)}', 'error')
        return redirect(url_for('create_proposal_form'))

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
    return color_map.get(urgency, 'text-gray-600')

if __name__ == '__main__':
    # Run in development mode
    app.run(debug=True, host='0.0.0.0', port=5000)