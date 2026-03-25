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
    yaml_path = PROPOSALS_DIR / f"{proposal_id}.yaml"
    
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

        # Validate required fields
        if not comment:
            flash('Comment is required — your reasoning matters to the collective.', 'error')
            return redirect(url_for('proposal_detail', proposal_id=proposal_id) + '#consultation-form')

        if position not in ('support', 'support-with-concerns', 'block'):
            flash('Invalid position selected.', 'error')
            return redirect(url_for('proposal_detail', proposal_id=proposal_id) + '#consultation-form')

        # Build consultation entry
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

        # Append to proposal consultations
        if 'consultations' not in proposal:
            proposal['consultations'] = []
        proposal['consultations'].append(consultation)

        # Add to consensus history
        if 'consensus_history' not in proposal:
            proposal['consensus_history'] = []
        proposal['consensus_history'].append({
            'timestamp': consultation['timestamp'],
            'event': 'consultation_added',
            'actor': contributor,
            'details': f"Position: {position}"
        })

        # Save updated proposal
        yaml_path = PROPOSALS_DIR / f"{proposal_id}.yaml"
        with open(yaml_path, 'w') as f:
            yaml.safe_dump(proposal, f, default_flow_style=False, sort_keys=False)

        # Also update JSON for API compatibility
        json_path = PROPOSALS_DIR / f"{proposal_id}.json"
        if json_path.exists():
            with open(json_path, 'w') as f:
                json.dump(proposal, f, indent=2, default=str)

        flash(f'Consultation input from "{contributor}" recorded. Thank you for participating!', 'success')
        return redirect(url_for('proposal_detail', proposal_id=proposal_id) + '#consultations-heading')

    except Exception as e:
        flash(f'Error submitting consultation: {str(e)}', 'error')
        return redirect(url_for('proposal_detail', proposal_id=proposal_id))


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