# CollectiveFlow API Improvements Proposal

## Current State Analysis

The CollectiveFlow web interface currently provides a minimal REST API with two endpoints:
- `GET /api/proposals` - Returns all proposals with count
- `GET /api/proposal/<id>` - Returns specific proposal

The Flask application has CORS enabled globally and uses YAML file storage with JSON export compatibility.

## Proposed Improvements

### 1. API Documentation (OpenAPI/Swagger)

**Recommendation**: Add OpenAPI 3.0 documentation using `flask-swagger-ui` and `apispec` libraries.

**Benefits for Horizontal Principles**:
- Transparent, self-documenting API reduces knowledge hierarchy
- Any agent/contributor can understand the API without special knowledge
- Eliminates need for "API gatekeepers" who hold documentation knowledge

**Implementation**:
```python
# requirements.txt additions
apispec==6.3.1
apispec-webframeworks==1.0.0
flask-swagger-ui==4.11.1

# In app.py
from apispec import APISpec
from apispec.ext.marshmallow import MarshmallowPlugin
from apispec_webframeworks.flask import FlaskPlugin
from flask_swagger_ui import get_swaggerui_blueprint

# Create OpenAPI spec
spec = APISpec(
    title="CollectiveFlow API",
    version="1.0.0",
    openapi_version="3.0.3",
    info={
        "description": "Horizontal decision-making API - no authentication, no hierarchy",
        "x-principles": [
            "No authentication required - collective transparency",
            "All endpoints equally accessible",
            "Read-only by default to protect consensus process",
            "Write operations require collective participation"
        ]
    },
    servers=[
        {"url": "http://localhost:5000", "description": "Local development"},
    ],
    plugins=[FlaskPlugin(), MarshmallowPlugin()]
)

# Swagger UI setup
SWAGGER_URL = '/api/docs'
API_URL = '/api/spec.json'

swaggerui_blueprint = get_swaggerui_blueprint(
    SWAGGER_URL,
    API_URL,
    config={'app_name': "CollectiveFlow API - Horizontal Decision Making"}
)
app.register_blueprint(swaggerui_blueprint, url_prefix=SWAGGER_URL)

@app.route('/api/spec.json')
def api_spec():
    """Generate OpenAPI specification."""
    return jsonify(spec.to_dict())
```

**API Documentation Principles**:
- Document the *why* of horizontal design, not just the *how*
- Include examples that demonstrate collective participation
- Explain why authentication is intentionally absent
- Show how the API prevents hierarchy creation

---

### 2. Additional Useful Endpoints

#### 2.1 Filtering and Querying
```python
@app.route('/api/proposals/active')
def api_proposals_active():
    """
    Get proposals requiring collective attention.

    Returns proposals in 'proposed' or 'consultation' status,
    sorted by urgency to help agents prioritize collective work.
    """
    proposals = load_proposals()
    active = [
        p for p in proposals
        if p.get('status') in ['proposed', 'consultation']
    ]

    # Sort by urgency (emergency first)
    urgency_order = {'emergency': 0, 'high': 1, 'medium': 2, 'low': 3}
    active.sort(key=lambda p: urgency_order.get(p.get('urgency', 'medium'), 2))

    return jsonify({
        'active_proposals': active,
        'count': len(active),
        'requires_attention': len([p for p in active if p.get('urgency') in ['emergency', 'high']])
    })

@app.route('/api/proposals/status/<status>')
def api_proposals_by_status(status):
    """
    Filter proposals by status.

    Horizontal transparency: All statuses equally queryable.
    """
    valid_statuses = ['proposed', 'consultation', 'consensus',
                     'implemented', 'blocked', 'withdrawn']

    if status not in valid_statuses:
        return jsonify({
            'error': 'Invalid status',
            'valid_statuses': valid_statuses
        }), 400

    proposals = load_proposals()
    filtered = [p for p in proposals if p.get('status') == status]

    return jsonify({
        'status': status,
        'proposals': filtered,
        'count': len(filtered)
    })

@app.route('/api/proposals/urgency/<urgency>')
def api_proposals_by_urgency(urgency):
    """Filter proposals by urgency level."""
    valid_urgencies = ['low', 'medium', 'high', 'emergency']

    if urgency not in valid_urgencies:
        return jsonify({
            'error': 'Invalid urgency',
            'valid_urgencies': valid_urgencies
        }), 400

    proposals = load_proposals()
    filtered = [p for p in proposals if p.get('urgency') == urgency]

    return jsonify({
        'urgency': urgency,
        'proposals': filtered,
        'count': len(filtered)
    })

@app.route('/api/proposals/search')
def api_proposals_search():
    """
    Search proposals by title/description.

    Query params:
    - q: Search query
    - status: Filter by status
    - urgency: Filter by urgency
    """
    query = request.args.get('q', '').lower()
    status_filter = request.args.get('status')
    urgency_filter = request.args.get('urgency')

    proposals = load_proposals()

    # Apply filters
    if query:
        proposals = [
            p for p in proposals
            if query in p.get('title', '').lower()
            or query in p.get('description', '').lower()
        ]

    if status_filter:
        proposals = [p for p in proposals if p.get('status') == status_filter]

    if urgency_filter:
        proposals = [p for p in proposals if p.get('urgency') == urgency_filter]

    return jsonify({
        'query': query,
        'filters': {
            'status': status_filter,
            'urgency': urgency_filter
        },
        'proposals': proposals,
        'count': len(proposals)
    })
```

#### 2.2 Collective Statistics
```python
@app.route('/api/stats')
def api_stats():
    """
    Collective activity statistics.

    Non-hierarchical metrics that show collective health,
    not individual performance rankings.
    """
    proposals = load_proposals()

    # Status distribution
    status_counts = {}
    for p in proposals:
        status = p.get('status', 'proposed')
        status_counts[status] = status_counts.get(status, 0) + 1

    # Urgency distribution
    urgency_counts = {}
    for p in proposals:
        urgency = p.get('urgency', 'medium')
        urgency_counts[urgency] = urgency_counts.get(urgency, 0) + 1

    # Consensus rate
    total = len(proposals)
    consensus_reached = len([
        p for p in proposals
        if p.get('status') in ['consensus', 'implemented']
    ])
    consensus_rate = (consensus_reached / total * 100) if total > 0 else 0

    # Active consultation count
    active_consultations = len([
        p for p in proposals
        if p.get('status') == 'consultation'
    ])

    # Unique contributors (agents)
    contributors = set()
    for p in proposals:
        contributors.add(p.get('proposer', 'anonymous'))
        for consultation in p.get('consultations', []):
            contributors.add(consultation.get('contributor', 'anonymous'))

    return jsonify({
        'collective_health': {
            'total_proposals': total,
            'active_consultations': active_consultations,
            'consensus_rate': round(consensus_rate, 1),
            'participating_agents': len(contributors)
        },
        'status_distribution': status_counts,
        'urgency_distribution': urgency_counts,
        'generated_at': datetime.now().isoformat()
    })

@app.route('/api/stats/timeline')
def api_stats_timeline():
    """
    Collective activity over time.

    Shows proposal creation and consensus trends without
    creating competitive metrics or performance hierarchies.
    """
    proposals = load_proposals()

    # Group by month
    monthly = {}
    for p in proposals:
        date_str = p.get('date', '')
        if date_str:
            try:
                dt = datetime.fromisoformat(date_str.replace('Z', '+00:00'))
                month_key = dt.strftime('%Y-%m')

                if month_key not in monthly:
                    monthly[month_key] = {
                        'proposals': 0,
                        'consensus': 0,
                        'implemented': 0
                    }

                monthly[month_key]['proposals'] += 1

                if p.get('status') == 'consensus':
                    monthly[month_key]['consensus'] += 1
                elif p.get('status') == 'implemented':
                    monthly[month_key]['implemented'] += 1
                    monthly[month_key]['consensus'] += 1  # Implemented implies consensus

            except:
                continue

    # Convert to sorted list
    timeline = [
        {
            'month': month,
            **data
        }
        for month, data in sorted(monthly.items())
    ]

    return jsonify({
        'timeline': timeline,
        'periods': len(timeline)
    })
```

#### 2.3 Consultation Access
```python
@app.route('/api/proposal/<proposal_id>/consultations')
def api_proposal_consultations(proposal_id):
    """
    Get all consultations for a proposal.

    Transparent access to collective input.
    """
    proposal = get_proposal(proposal_id)

    if not proposal:
        return jsonify({'error': 'Proposal not found'}), 404

    consultations = proposal.get('consultations', [])

    # Calculate consensus metrics
    total = len(consultations)
    supporting = len([c for c in consultations if c.get('support', False)])
    concerns = []
    for c in consultations:
        concerns.extend(c.get('concerns', []))

    return jsonify({
        'proposal_id': proposal_id,
        'consultations': consultations,
        'consensus_metrics': {
            'total_consultations': total,
            'supporting': supporting,
            'concerns_raised': len(concerns),
            'unanimous_support': total > 0 and supporting == total
        },
        'blocking_concerns': list(set(concerns))
    })

@app.route('/api/proposal/<proposal_id>/history')
def api_proposal_history(proposal_id):
    """
    Get consensus history for a proposal.

    Complete transparency of decision-making process.
    """
    proposal = get_proposal(proposal_id)

    if not proposal:
        return jsonify({'error': 'Proposal not found'}), 404

    return jsonify({
        'proposal_id': proposal_id,
        'title': proposal.get('title'),
        'status': proposal.get('status'),
        'consensus_history': proposal.get('consensus_history', []),
        'current_status': proposal.get('consensus_status')
    })
```

#### 2.4 Health Check & Metadata
```python
@app.route('/api/health')
def api_health():
    """
    API health check.

    Returns system status and collective principles.
    """
    # Check if data directory is accessible
    data_accessible = PROPOSALS_DIR.exists()
    proposal_count = len(list(PROPOSALS_DIR.glob('*.yaml'))) if data_accessible else 0

    return jsonify({
        'status': 'healthy' if data_accessible else 'degraded',
        'version': '1.0.0',
        'principles': {
            'authentication': 'none - collective transparency',
            'authorization': 'horizontal - no privileged access',
            'data_ownership': 'collective',
            'decision_model': 'consensus-based'
        },
        'storage': {
            'accessible': data_accessible,
            'proposal_count': proposal_count,
            'backend': 'file-based YAML'
        },
        'timestamp': datetime.now().isoformat()
    })

@app.route('/api/schema')
def api_schema():
    """
    Return proposal schema documentation.

    Helps agents understand data structure without hierarchy.
    """
    return jsonify({
        'proposal': {
            'id': 'string (auto-generated)',
            'title': 'string (required)',
            'description': 'string (required)',
            'proposer': 'string (required for transparency)',
            'date': 'ISO 8601 datetime',
            'status': {
                'type': 'enum',
                'values': ['proposed', 'consultation', 'consensus',
                          'implemented', 'blocked', 'withdrawn']
            },
            'urgency': {
                'type': 'enum',
                'values': ['low', 'medium', 'high', 'emergency']
            },
            'affected_areas': 'array of strings',
            'consensus_status': 'string',
            'consensus_history': 'array of events',
            'consultations': 'array of consultation objects'
        },
        'consultation': {
            'contributor': 'string (agent/person name)',
            'timestamp': 'ISO 8601 datetime',
            'input': 'string (detailed input)',
            'concerns': 'array of strings',
            'support': 'boolean'
        },
        'consensus_event': {
            'timestamp': 'ISO 8601 datetime',
            'event': 'string (event type)',
            'actor': 'string (who performed action)',
            'details': 'string (optional details)'
        }
    })
```

---

### 3. Response Format Improvements

**Current Issues**:
- Inconsistent response structures
- No metadata about the response itself
- Missing pagination support
- No HATEOAS links

**Proposed Standard Response Format**:
```python
def api_response(data=None, error=None, status=200, meta=None):
    """
    Standard API response format.

    Consistent structure helps all agents interpret responses equally.
    """
    response = {
        'timestamp': datetime.now().isoformat(),
        'status': 'success' if not error else 'error'
    }

    if error:
        response['error'] = error

    if data is not None:
        response['data'] = data

    if meta:
        response['meta'] = meta

    # Add horizontal principles reminder
    response['_principles'] = {
        'access': 'equal for all',
        'transparency': 'complete data visibility',
        'collective': 'no individual authority'
    }

    return jsonify(response), status

# Usage examples:
@app.route('/api/proposals')
def api_proposals():
    """API endpoint for proposals list."""
    try:
        proposals = load_proposals()

        # Pagination support
        page = request.args.get('page', 1, type=int)
        per_page = request.args.get('per_page', 20, type=int)
        per_page = min(per_page, 100)  # Max 100 per page

        total = len(proposals)
        start = (page - 1) * per_page
        end = start + per_page

        paginated = proposals[start:end]

        return api_response(
            data=paginated,
            meta={
                'pagination': {
                    'page': page,
                    'per_page': per_page,
                    'total': total,
                    'total_pages': (total + per_page - 1) // per_page
                },
                'count': len(paginated)
            }
        )
    except Exception as e:
        return api_response(error=str(e), status=500)
```

**HATEOAS Links** (Hypermedia as the Engine of Application State):
```python
def add_proposal_links(proposal):
    """Add HATEOAS links to proposal object."""
    proposal_id = proposal.get('id')

    links = {
        'self': f'/api/proposal/{proposal_id}',
        'consultations': f'/api/proposal/{proposal_id}/consultations',
        'history': f'/api/proposal/{proposal_id}/history',
        'web_view': f'/proposal/{proposal_id}'
    }

    # Add action links based on status
    status = proposal.get('status')
    if status == 'consultation':
        links['add_consultation'] = f'/api/proposal/{proposal_id}/consultations'

    proposal['_links'] = links
    return proposal
```

---

### 4. Error Handling Enhancements

**Current Issues**:
- Generic error messages
- No error codes for programmatic handling
- Missing validation details

**Proposed Error Handling**:
```python
# Error response schema
class APIError(Exception):
    """Base API error with structured details."""

    def __init__(self, message, code=None, details=None, status=400):
        self.message = message
        self.code = code or 'API_ERROR'
        self.details = details or {}
        self.status = status
        super().__init__(self.message)

    def to_dict(self):
        return {
            'error': {
                'message': self.message,
                'code': self.code,
                'details': self.details
            },
            'timestamp': datetime.now().isoformat(),
            '_principles': {
                'transparency': 'Errors shown clearly to prevent knowledge hierarchy',
                'horizontal': 'All users see same error information'
            }
        }

# Common error types
class NotFoundError(APIError):
    def __init__(self, resource, resource_id):
        super().__init__(
            message=f"{resource} not found",
            code='NOT_FOUND',
            details={'resource': resource, 'id': resource_id},
            status=404
        )

class ValidationError(APIError):
    def __init__(self, field_errors):
        super().__init__(
            message="Validation failed",
            code='VALIDATION_ERROR',
            details={'fields': field_errors},
            status=400
        )

class InvalidStatusError(APIError):
    def __init__(self, current, attempted, allowed):
        super().__init__(
            message="Invalid status transition",
            code='INVALID_TRANSITION',
            details={
                'current_status': current,
                'attempted_status': attempted,
                'allowed_transitions': allowed
            },
            status=400
        )

# Error handler
@app.errorhandler(APIError)
def handle_api_error(error):
    return jsonify(error.to_dict()), error.status

@app.errorhandler(404)
def handle_404(error):
    return jsonify({
        'error': {
            'message': 'Endpoint not found',
            'code': 'NOT_FOUND'
        },
        'available_endpoints': {
            'proposals': '/api/proposals',
            'proposal': '/api/proposal/<id>',
            'active': '/api/proposals/active',
            'stats': '/api/stats',
            'docs': '/api/docs'
        },
        'timestamp': datetime.now().isoformat()
    }), 404

@app.errorhandler(500)
def handle_500(error):
    return jsonify({
        'error': {
            'message': 'Internal server error',
            'code': 'INTERNAL_ERROR',
            'details': 'Please check server logs for details'
        },
        'timestamp': datetime.now().isoformat()
    }), 500

# Usage in endpoints
@app.route('/api/proposal/<proposal_id>')
def api_proposal(proposal_id):
    proposal = get_proposal(proposal_id)

    if not proposal:
        raise NotFoundError('proposal', proposal_id)

    return api_response(data=add_proposal_links(proposal))
```

---

### 5. API Versioning Strategy

**Recommendation**: URL-based versioning with clear transition path

**Rationale for Horizontal Principles**:
- Transparent version visibility in URL
- No hidden version negotiation
- Easy to understand and use without special knowledge
- All versions equally accessible during transition

**Implementation**:
```python
# Version 1 (current)
@app.route('/api/v1/proposals')
def api_v1_proposals():
    """Version 1 API - current implementation."""
    return api_proposals()

# Default to latest version
@app.route('/api/proposals')
def api_proposals_latest():
    """Latest version (defaults to v1)."""
    return api_v1_proposals()

# Version info endpoint
@app.route('/api/versions')
def api_versions():
    """
    Show available API versions.

    Transparent versioning - no hidden deprecation.
    """
    return jsonify({
        'current': 'v1',
        'available': ['v1'],
        'deprecated': [],
        'sunset': {},  # Version: sunset date
        'endpoints': {
            'v1': '/api/v1/',
            'latest': '/api/'
        },
        'principles': {
            'versioning_strategy': 'URL-based for transparency',
            'deprecation_policy': 'Collective decision with notice period',
            'breaking_changes': 'New version only, old version maintained'
        }
    })

# Future version example
@app.route('/api/v2/proposals')
def api_v2_proposals():
    """
    Version 2 API with enhanced features.

    Both v1 and v2 available during transition period.
    """
    # New implementation with backward-incompatible changes
    pass
```

**Version Deprecation Process** (Horizontal):
1. Proposal in CollectiveFlow to deprecate version
2. Collective consensus on timeline
3. Public announcement with long notice period
4. Both versions run in parallel
5. Old version sunset only after consensus

---

### 6. CORS Configuration Review

**Current Implementation**:
```python
CORS(app)  # Enable cross-origin requests for API compatibility
```

**Issues**:
- Allows all origins (security risk)
- No credential handling configuration
- Missing preflight optimization

**Recommended Configuration**:
```python
from flask_cors import CORS

# Horizontal CORS - Open by design but with sensible defaults
CORS(app, resources={
    r"/api/*": {
        "origins": os.environ.get('CORS_ORIGINS', '*').split(','),
        "methods": ["GET", "HEAD", "OPTIONS"],
        "allow_headers": ["Content-Type", "Accept"],
        "expose_headers": ["X-Total-Count", "X-Page", "X-Per-Page"],
        "max_age": 3600,
        "supports_credentials": False  # No auth = no credentials needed
    }
})

# Add CORS info to API metadata
@app.route('/api/cors-policy')
def api_cors_policy():
    """
    Explain CORS configuration.

    Transparency about access controls.
    """
    return jsonify({
        'policy': 'Open by default',
        'principles': {
            'access': 'No restrictions on origins for read operations',
            'methods': 'GET, HEAD, OPTIONS for collective transparency',
            'credentials': 'Not needed - no authentication system',
            'security': 'Horizontal access without barriers'
        },
        'configuration': {
            'allowed_origins': os.environ.get('CORS_ORIGINS', '*'),
            'allowed_methods': ['GET', 'HEAD', 'OPTIONS'],
            'max_age': 3600
        },
        'note': 'Write operations future: require collective process, not authentication'
    })
```

**Environment-Based Configuration**:
```bash
# .env file
CORS_ORIGINS=http://localhost:3000,http://localhost:5173,https://collective.example.com

# For production, be explicit:
CORS_ORIGINS=https://collectiveflow.example.com,https://collective.example.org
```

---

## Complete Enhanced API Structure

```
GET  /api/health                      - Health check & principles
GET  /api/versions                    - API version info
GET  /api/schema                      - Data schema documentation
GET  /api/cors-policy                 - CORS configuration info

GET  /api/proposals                   - List all (paginated)
GET  /api/proposals/active            - Active proposals needing attention
GET  /api/proposals/status/<status>   - Filter by status
GET  /api/proposals/urgency/<urgency> - Filter by urgency
GET  /api/proposals/search            - Search with filters
GET  /api/proposal/<id>               - Specific proposal
GET  /api/proposal/<id>/consultations - Proposal consultations
GET  /api/proposal/<id>/history       - Consensus history

GET  /api/stats                       - Collective statistics
GET  /api/stats/timeline              - Activity over time

GET  /api/docs                        - Swagger UI documentation
GET  /api/spec.json                   - OpenAPI specification
```

---

## Implementation Priority

### Phase 1: Foundation (Immediate)
1. ✅ Standard response format with `api_response()`
2. ✅ Enhanced error handling with structured errors
3. ✅ Basic filtering endpoints (active, status, urgency)
4. ✅ Health check endpoint

### Phase 2: Documentation (Week 1)
1. ✅ OpenAPI/Swagger integration
2. ✅ Schema documentation endpoint
3. ✅ CORS policy documentation
4. ✅ Version information endpoint

### Phase 3: Enhanced Features (Week 2)
1. ✅ Search endpoint with multiple filters
2. ✅ Consultation-specific endpoints
3. ✅ Statistics and timeline endpoints
4. ✅ HATEOAS links in responses

### Phase 4: Polish (Week 3)
1. ✅ Pagination for all list endpoints
2. ✅ Rate limiting (gentle, transparent)
3. ✅ Response caching headers
4. ✅ API versioning structure

---

## Horizontal Design Principles Maintained

### ✅ No Authentication
- All endpoints public by design
- Transparency over access control
- Trust-based collective model

### ✅ Equal Access
- Same data available to all
- No privileged endpoints
- No "admin" routes

### ✅ Transparent Errors
- Clear error messages for all
- No hidden system details
- Help everyone debug equally

### ✅ Self-Documenting
- OpenAPI specification
- Schema documentation
- Inline principle explanations

### ✅ Collective-First Design
- Endpoints reflect consensus process
- Statistics show collective health, not rankings
- Filters help collective prioritization

---

## Testing Strategy

### Unit Tests
```python
# tests/test_api.py
import pytest
from app import app

@pytest.fixture
def client():
    app.config['TESTING'] = True
    with app.test_client() as client:
        yield client

def test_health_endpoint(client):
    """Health endpoint returns system status."""
    rv = client.get('/api/health')
    assert rv.status_code == 200
    data = rv.get_json()
    assert 'status' in data
    assert 'principles' in data

def test_proposals_list(client):
    """Proposals endpoint returns paginated list."""
    rv = client.get('/api/proposals')
    assert rv.status_code == 200
    data = rv.get_json()
    assert 'data' in data
    assert 'meta' in data
    assert 'pagination' in data['meta']

def test_proposal_not_found(client):
    """Missing proposal returns structured error."""
    rv = client.get('/api/proposal/nonexistent')
    assert rv.status_code == 404
    data = rv.get_json()
    assert 'error' in data
    assert data['error']['code'] == 'NOT_FOUND'

def test_cors_headers(client):
    """CORS headers present on API routes."""
    rv = client.get('/api/proposals')
    assert 'Access-Control-Allow-Origin' in rv.headers
```

### Integration Tests
```python
def test_proposal_workflow(client):
    """Test complete proposal access workflow."""
    # List proposals
    rv = client.get('/api/proposals')
    proposals = rv.get_json()['data']

    if proposals:
        # Get specific proposal
        proposal_id = proposals[0]['id']
        rv = client.get(f'/api/proposal/{proposal_id}')
        assert rv.status_code == 200

        # Get consultations
        rv = client.get(f'/api/proposal/{proposal_id}/consultations')
        assert rv.status_code == 200

        # Get history
        rv = client.get(f'/api/proposal/{proposal_id}/history')
        assert rv.status_code == 200
```

---

## Security Considerations

### Current Security Posture
- ✅ No authentication = no credential theft risk
- ✅ Read-only API = no data tampering via API
- ✅ YAML storage = easy to backup and audit
- ⚠️ CORS fully open = potential for client-side attacks
- ⚠️ No rate limiting = potential for abuse

### Recommended Security Enhancements

#### Rate Limiting (Gentle & Transparent)
```python
from flask_limiter import Limiter
from flask_limiter.util import get_remote_address

limiter = Limiter(
    app=app,
    key_func=get_remote_address,
    default_limits=["1000 per hour", "100 per minute"],
    storage_uri="memory://",
    strategy="fixed-window"
)

# More generous limits reflect collective trust
@app.route('/api/proposals')
@limiter.limit("500 per hour")
def api_proposals():
    """Generous rate limit - trust-based."""
    pass

# Rate limit info endpoint
@app.route('/api/rate-limits')
def api_rate_limits():
    """
    Explain rate limiting policy.

    Transparent about technical constraints.
    """
    return jsonify({
        'policy': 'Gentle limits to prevent abuse while allowing collective access',
        'limits': {
            'default': '1000 requests/hour, 100 requests/minute',
            'per_endpoint': {
                '/api/proposals': '500 requests/hour',
                '/api/stats': '200 requests/hour'
            }
        },
        'principles': {
            'trust': 'Limits assume good faith',
            'transparency': 'Rate limit info publicly available',
            'generous': 'High limits for collective work'
        },
        'headers': {
            'X-RateLimit-Limit': 'Total requests allowed',
            'X-RateLimit-Remaining': 'Requests remaining',
            'X-RateLimit-Reset': 'Unix timestamp when limit resets'
        }
    })
```

#### Input Validation
```python
from werkzeug.exceptions import BadRequest

def validate_pagination(page, per_page):
    """Validate pagination parameters."""
    if page < 1:
        raise ValidationError({'page': 'Must be >= 1'})

    if per_page < 1 or per_page > 100:
        raise ValidationError({'per_page': 'Must be between 1 and 100'})

def validate_status(status):
    """Validate status parameter."""
    valid = ['proposed', 'consultation', 'consensus',
             'implemented', 'blocked', 'withdrawn']

    if status not in valid:
        raise ValidationError({
            'status': f'Must be one of: {", ".join(valid)}'
        })
```

#### Response Headers
```python
@app.after_request
def add_security_headers(response):
    """Add security headers to all responses."""
    # Prevent clickjacking
    response.headers['X-Frame-Options'] = 'SAMEORIGIN'

    # XSS protection
    response.headers['X-Content-Type-Options'] = 'nosniff'

    # Referrer policy
    response.headers['Referrer-Policy'] = 'strict-origin-when-cross-origin'

    # Content Security Policy for API
    if request.path.startswith('/api/'):
        response.headers['Content-Security-Policy'] = "default-src 'none'"

    # Add horizontal principle header
    response.headers['X-Collective-Principles'] = 'horizontal,transparent,consensus'

    return response
```

---

## Monitoring & Observability

### Request Logging
```python
import logging
from datetime import datetime

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s',
    handlers=[
        logging.FileHandler('logs/api.log'),
        logging.StreamHandler()
    ]
)

@app.before_request
def log_request():
    """Log all API requests for transparency."""
    if request.path.startswith('/api/'):
        logging.info(f"API Request: {request.method} {request.path} from {request.remote_addr}")

@app.after_request
def log_response(response):
    """Log API responses."""
    if request.path.startswith('/api/'):
        logging.info(f"API Response: {response.status_code} for {request.path}")
    return response
```

### Metrics Endpoint
```python
# Simple in-memory metrics (for horizontal transparency)
request_counts = {}

@app.before_request
def track_request():
    """Track request metrics."""
    endpoint = request.endpoint or 'unknown'
    request_counts[endpoint] = request_counts.get(endpoint, 0) + 1

@app.route('/api/metrics')
def api_metrics():
    """
    Expose API metrics transparently.

    Shows collective usage patterns without tracking individuals.
    """
    return jsonify({
        'endpoints': request_counts,
        'principles': {
            'aggregation': 'Endpoint counts only, no individual tracking',
            'transparency': 'All metrics publicly visible',
            'collective': 'Shows collective usage, not user behavior'
        },
        'note': 'Metrics reset on server restart (no persistent tracking)'
    })
```

---

## Documentation Files

### API README
```markdown
# CollectiveFlow API

## Horizontal API Design

This API embodies libertarian socialist principles:
- **No authentication**: Collective transparency
- **Equal access**: All endpoints available to all
- **Self-documenting**: OpenAPI specification at /api/docs
- **Read-only**: Write operations require collective process

## Quick Start

### List All Proposals
```bash
curl http://localhost:5000/api/proposals
```

### Get Active Proposals
```bash
curl http://localhost:5000/api/proposals/active
```

### Search Proposals
```bash
curl "http://localhost:5000/api/proposals/search?q=consensus&status=consultation"
```

### View Documentation
Visit http://localhost:5000/api/docs for interactive Swagger UI

## Design Principles

1. **Transparency**: All data visible to all
2. **Horizontal Access**: No privileged endpoints
3. **Collective First**: Endpoints reflect consensus process
4. **Self-Documenting**: API explains itself
5. **Trust-Based**: Generous rate limits, no heavy auth

## Response Format

All API responses follow this structure:

```json
{
  "timestamp": "2025-01-05T10:00:00Z",
  "status": "success",
  "data": { ... },
  "meta": { ... },
  "_principles": {
    "access": "equal for all",
    "transparency": "complete data visibility",
    "collective": "no individual authority"
  }
}
```

## Error Handling

Errors are transparent and helpful:

```json
{
  "timestamp": "2025-01-05T10:00:00Z",
  "status": "error",
  "error": {
    "message": "Validation failed",
    "code": "VALIDATION_ERROR",
    "details": {
      "fields": {
        "status": "Must be one of: proposed, consultation, consensus"
      }
    }
  }
}
```

## Rate Limits

Generous limits to support collective work:
- Default: 1000 requests/hour
- High-traffic endpoints: 500 requests/hour
- Search: 200 requests/hour

Check headers:
- `X-RateLimit-Limit`: Total allowed
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Reset timestamp
```

---

## Summary

This proposal enhances the CollectiveFlow API while maintaining strict adherence to horizontal, consensus-based principles:

1. **OpenAPI Documentation**: Self-documenting API reduces knowledge hierarchy
2. **Enhanced Endpoints**: Better filtering, search, and statistics for collective work
3. **Standardized Responses**: Consistent format helps all agents equally
4. **Transparent Errors**: Clear error messages without hidden system details
5. **URL-Based Versioning**: Transparent version management
6. **Thoughtful CORS**: Open by default with sensible security

All improvements prioritize:
- ✅ Transparency over security theater
- ✅ Collective access over individual privileges
- ✅ Self-documentation over hidden knowledge
- ✅ Trust-based design over restrictive controls

The API remains authentication-free and horizontally accessible while adding features that help the collective work more effectively.
