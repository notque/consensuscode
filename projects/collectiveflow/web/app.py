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
from datetime import datetime
from pathlib import Path
from flask import Flask, render_template, jsonify, request, redirect, url_for, flash, session, abort

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
    """Load a specific proposal by ID.

    Security: validates proposal_id to prevent path traversal attacks.
    Only alphanumeric characters, hyphens, and underscores are allowed.
    """
    # Reject any proposal_id that could be used for path traversal
    if not re.match(r'^[a-zA-Z0-9_-]+$', proposal_id):
        return None

    yaml_path = PROPOSALS_DIR / f"{proposal_id}.yaml"

    # Defense in depth: verify resolved path stays within PROPOSALS_DIR
    try:
        yaml_path.resolve().relative_to(PROPOSALS_DIR.resolve())
    except ValueError:
        return None

    if yaml_path.exists():
        with open(yaml_path, 'r') as f:
            return yaml.safe_load(f)

    return None

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

@app.route('/api/proposals')
def api_proposals():
    """API endpoint for proposals list."""
    proposals = load_proposals()
    return jsonify({
        'proposals': proposals,
        'count': len(proposals)
    })

@app.route('/api/proposal/<proposal_id>')
def api_proposal(proposal_id):
    """API endpoint for a specific proposal."""
    proposal = get_proposal(proposal_id)
    
    if not proposal:
        return jsonify({'error': 'Proposal not found'}), 404
    
    return jsonify(proposal)

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

        flash('Proposal submitted successfully!', 'success')
        return redirect(url_for('proposal_detail', proposal_id=proposal_id))

    except Exception as e:
        # Log the real error server-side; show generic message to user
        logging.exception("Error creating proposal")
        flash('Error creating proposal. Please try again.', 'error')
        return redirect(url_for('create_proposal_form'))

@app.after_request
def set_security_headers(response):
    """Add security headers to every response.

    These headers defend against clickjacking, MIME-sniffing, XSS reflection,
    and other common browser-level attacks.
    """
    response.headers['X-Content-Type-Options'] = 'nosniff'
    response.headers['X-Frame-Options'] = 'DENY'
    response.headers['X-XSS-Protection'] = '1; mode=block'
    response.headers['Referrer-Policy'] = 'strict-origin-when-cross-origin'
    response.headers['Permissions-Policy'] = 'geolocation=(), camera=(), microphone=()'
    # Content-Security-Policy: allow Tailwind CDN but restrict everything else
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
    return color_map.get(urgency, 'text-gray-600')

if __name__ == '__main__':
    # Debug mode from environment only — never hardcode True in source
    debug_mode = os.environ.get('FLASK_DEBUG', '0').lower() in ('1', 'true', 'yes')
    app.run(debug=debug_mode, host='0.0.0.0', port=5000)