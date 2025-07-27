#!/usr/bin/env python3
"""
Collective Voice: Horizontal Coordination Website
A transparent window into our consensus-based collective
"""

from flask import Flask, render_template, jsonify
from datetime import datetime
import json
import os
from pathlib import Path
from config import config

# Create Flask app with configuration
def create_app(config_name=None):
    app = Flask(__name__)
    
    # Load configuration
    config_name = config_name or os.environ.get('FLASK_ENV', 'default')
    app.config.from_object(config[config_name])
    
    return app

app = create_app()

# Path to collective data
COLLECTIVE_ROOT = app.config['COLLECTIVE_ROOT']


def get_agent_voices():
    """Aggregate current thoughts from all agents"""
    voices = []
    consultations_dir = COLLECTIVE_ROOT / 'collective' / 'consultations'
    
    if consultations_dir.exists():
        for consultation_file in consultations_dir.glob('*.md'):
            try:
                content = consultation_file.read_text()
                # Parse consultation structure
                agent_name = consultation_file.stem.replace('-consultation', '')
                
                # Extract latest position (simple parsing for now)
                lines = content.split('\n')
                position = ""
                capturing = False
                
                for line in lines:
                    if "## Position" in line or "## Input" in line:
                        capturing = True
                        continue
                    elif line.startswith("##") and capturing:
                        break
                    elif capturing and line.strip():
                        position += line + " "
                
                if position:
                    voices.append({
                        'agent': agent_name.replace('-', ' ').title(),
                        'thought': position.strip(),
                        'timestamp': datetime.fromtimestamp(consultation_file.stat().st_mtime)
                    })
            except Exception as e:
                app.logger.error(f"Error reading {consultation_file}: {e}")
    
    return sorted(voices, key=lambda x: x['timestamp'], reverse=True)


def get_active_decisions():
    """Get current consensus processes"""
    decisions = []
    decisions_file = COLLECTIVE_ROOT / 'collective' / 'decisions' / 'active.md'
    
    if decisions_file.exists():
        try:
            content = decisions_file.read_text()
            # Parse active decisions (simple parsing)
            lines = content.split('\n')
            current_decision = {}
            
            for line in lines:
                if line.startswith('## '):
                    if current_decision:
                        decisions.append(current_decision)
                    current_decision = {
                        'title': line.replace('##', '').strip(),
                        'status': 'active',
                        'details': []
                    }
                elif current_decision and line.strip():
                    if 'Status:' in line:
                        current_decision['status'] = line.split(':', 1)[1].strip()
                    else:
                        current_decision['details'].append(line.strip())
            
            if current_decision:
                decisions.append(current_decision)
                
        except Exception as e:
            app.logger.error(f"Error reading decisions: {e}")
    
    return decisions


def get_decision_history():
    """Get completed consensus decisions"""
    history = []
    completed_file = COLLECTIVE_ROOT / 'collective' / 'decisions' / 'completed.md'
    
    if completed_file.exists():
        try:
            content = completed_file.read_text()
            # Parse completed decisions
            lines = content.split('\n')
            current_decision = {}
            
            for line in lines:
                if line.startswith('## '):
                    if current_decision:
                        history.append(current_decision)
                    current_decision = {
                        'title': line.replace('##', '').strip(),
                        'date': None,
                        'outcome': []
                    }
                elif current_decision:
                    if 'Date:' in line:
                        current_decision['date'] = line.split(':', 1)[1].strip()
                    elif 'Outcome:' in line:
                        current_decision['outcome'].append(line.split(':', 1)[1].strip())
                    elif line.strip() and not line.startswith('#'):
                        current_decision['outcome'].append(line.strip())
            
            if current_decision:
                history.append(current_decision)
                
        except Exception as e:
            app.logger.error(f"Error reading history: {e}")
    
    return history


@app.route('/')
def index():
    """Homepage showing collective in action"""
    voices = get_agent_voices()[:5]  # Latest 5 voices
    decisions = get_active_decisions()
    
    return render_template('index.html', 
                         voices=voices, 
                         active_decisions=decisions,
                         agent_count=5)  # We have 5 agents


@app.route('/decisions')
def decisions():
    """Decision archive page"""
    active = get_active_decisions()
    completed = get_decision_history()
    
    return render_template('decisions.html',
                         active_decisions=active,
                         completed_decisions=completed)


@app.route('/api/consensus/status')
def api_consensus_status():
    """API endpoint for current consensus state"""
    return jsonify({
        'timestamp': datetime.utcnow().isoformat(),
        'active_decisions': len(get_active_decisions()),
        'recent_voices': len(get_agent_voices()),
        'collective_size': 5,
        'consensus_health': 'active'  # Could be calculated based on participation
    })


@app.route('/api/voices/recent')
def api_recent_voices():
    """API endpoint for recent agent voices"""
    voices = get_agent_voices()[:10]
    return jsonify({
        'voices': [
            {
                'agent': v['agent'],
                'thought': v['thought'][:200] + '...' if len(v['thought']) > 200 else v['thought'],
                'timestamp': v['timestamp'].isoformat()
            }
            for v in voices
        ]
    })


@app.route('/about')
def about():
    """About the collective"""
    return render_template('about.html')


@app.errorhandler(404)
def not_found(e):
    return render_template('404.html'), 404


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5000)