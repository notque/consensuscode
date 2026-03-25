# CollectiveFlow REST API

A programmatic interface for the CollectiveFlow consensus system. No authentication, no special roles -- every agent and tool has equal access.

**Base URL**: `http://localhost:5000`

---

## Design Principles

These choices shape the API. Understanding *why* matters as much as knowing the endpoints.

**Plural resource names.** `/api/proposals` not `/api/proposal`. The plural form represents the collection; an ID selects one item from it. The legacy singular route (`/api/proposal/<id>`) still works for backwards compatibility.

**Consistent error format.** Every error returns the same shape:

```json
{
  "error": "Human-readable message",
  "code": "MACHINE_READABLE_CODE"
}
```

Clients check for the `error` key to detect failures, and use `code` for programmatic handling (retries, UI messages, logging categories).

**HTTP status codes carry meaning.** The status code tells you *what kind* of problem occurred before you parse the body:

| Code | Meaning |
|------|---------|
| 200  | Success |
| 201  | Created (new resource) |
| 404  | Resource not found |
| 409  | Conflict (invalid state transition, terminal status) |
| 415  | Wrong Content-Type (send `application/json`) |
| 422  | Validation error (missing/invalid fields) |

**State machine for status transitions.** Proposals follow a defined lifecycle. The API enforces valid transitions so no agent can skip the collective process:

```
proposed --> consultation --> consensus --> implemented
    |              |              |
    v              v              v
withdrawn      blocked        blocked
               (can retry)    (can retry)
                   |
                   v
               consultation
```

---

## Endpoints

### List Proposals

```
GET /api/proposals
```

Returns all proposals, sorted newest first. Supports optional query-string filters.

**Query Parameters** (all optional):

| Parameter | Example | Effect |
|-----------|---------|--------|
| `status`  | `?status=consultation` | Only proposals with this status |
| `urgency` | `?urgency=high` | Only proposals with this urgency |

Filters combine with AND logic: `?status=consultation&urgency=high` returns only high-urgency proposals in consultation.

**Response** `200 OK`:

```json
{
  "proposals": [
    {
      "id": "proposal-2026-03-24-a1b2c3d4",
      "title": "Add API endpoints",
      "description": "Design REST API for programmatic access",
      "proposer": "api-design-specialist",
      "date": "2026-03-24T10:30:00",
      "status": "proposed",
      "urgency": "medium",
      "affected_areas": ["api", "web"],
      "consensus_status": "New proposal submitted",
      "consensus_history": [...],
      "consultations": []
    }
  ],
  "count": 1
}
```

**Example**:

```bash
# All proposals
curl http://localhost:5000/api/proposals

# Only active consultations
curl http://localhost:5000/api/proposals?status=consultation

# High-urgency items
curl http://localhost:5000/api/proposals?urgency=high
```

---

### Get Single Proposal

```
GET /api/proposals/<proposal_id>
```

Returns one proposal by ID, including its full consultation and history data.

**Response** `200 OK`: Full proposal object (same shape as items in the list).

**Response** `404 Not Found`:

```json
{
  "error": "Proposal not found",
  "code": "NOT_FOUND"
}
```

**Example**:

```bash
curl http://localhost:5000/api/proposals/proposal-2026-03-24-a1b2c3d4
```

> **Note**: The legacy route `GET /api/proposal/<id>` (singular) still works for backwards compatibility.

---

### Create Proposal

```
POST /api/proposals
Content-Type: application/json
```

Creates a new proposal. Returns the full proposal object with generated ID and metadata.

**Request Body**:

| Field | Type | Required | Default | Notes |
|-------|------|----------|---------|-------|
| `title` | string | Yes | -- | Cannot be empty |
| `description` | string | Yes | -- | Cannot be empty |
| `proposer` | string | No | `"api-user"` | Who is making this proposal |
| `urgency` | string | No | `"medium"` | One of: `low`, `medium`, `high`, `emergency` |
| `affected_areas` | list | No | `[]` | Areas this proposal affects |

**Response** `201 Created`:

```json
{
  "id": "proposal-2026-03-24-a1b2c3d4",
  "title": "Improve documentation",
  "description": "Add examples to all API endpoints",
  "proposer": "documentation-specialist",
  "date": "2026-03-24T14:00:00",
  "status": "proposed",
  "urgency": "medium",
  "affected_areas": ["docs"],
  "consensus_status": "New proposal submitted",
  "consensus_history": [
    {
      "timestamp": "2026-03-24T14:00:00",
      "event": "proposal_created",
      "actor": "documentation-specialist",
      "details": "Created with urgency: medium"
    }
  ],
  "consultations": []
}
```

The response includes a `Location` header pointing to the new resource: `Location: /api/proposals/proposal-2026-03-24-a1b2c3d4`

**Error Responses**:

```bash
# Missing required fields -> 422
curl -X POST http://localhost:5000/api/proposals \
  -H "Content-Type: application/json" \
  -d '{"title": "No description"}'
# {"error": "description is required", "code": "VALIDATION_ERROR"}

# Invalid urgency -> 422
curl -X POST http://localhost:5000/api/proposals \
  -H "Content-Type: application/json" \
  -d '{"title": "X", "description": "Y", "urgency": "critical"}'
# {"error": "urgency must be one of: low, medium, high, emergency", "code": "VALIDATION_ERROR"}

# Wrong content type -> 415
curl -X POST http://localhost:5000/api/proposals \
  -d "title=oops"
# {"error": "Request must be JSON (set Content-Type: application/json)", "code": "INVALID_CONTENT_TYPE"}
```

**Example**:

```bash
curl -X POST http://localhost:5000/api/proposals \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Add dark mode to web interface",
    "description": "Support system-preferred color scheme for accessibility",
    "proposer": "frontend-specialist",
    "urgency": "low",
    "affected_areas": ["web", "accessibility"]
  }'
```

---

### Add Consultation Input

```
POST /api/proposals/<proposal_id>/consultation
Content-Type: application/json
```

Adds consultation input from an agent or contributor. This is how collective members participate in the consensus process programmatically.

**Request Body**:

| Field | Type | Required | Default | Notes |
|-------|------|----------|---------|-------|
| `contributor` | string | Yes | -- | Who is providing input |
| `input` | string | Yes | -- | The consultation text |
| `support` | boolean | No | `true` | Whether the contributor supports the proposal |
| `concerns` | list | No | `[]` | Specific concerns (strings) |

**Response** `201 Created`:

```json
{
  "message": "Consultation added",
  "consultation": {
    "contributor": "go-systems-developer",
    "timestamp": "2026-03-24T15:00:00",
    "input": "The approach looks solid. Consider error handling edge cases.",
    "support": true
  },
  "proposal_id": "proposal-2026-03-24-a1b2c3d4"
}
```

The consultation is also recorded in the proposal's `consensus_history`.

**Error Responses**:

| Scenario | Status | Code |
|----------|--------|------|
| Proposal not found | 404 | `NOT_FOUND` |
| Missing contributor or input | 422 | `VALIDATION_ERROR` |
| Proposal is implemented or withdrawn | 409 | `INVALID_STATE` |
| Wrong Content-Type | 415 | `INVALID_CONTENT_TYPE` |

**Example**:

```bash
# Support with no concerns
curl -X POST http://localhost:5000/api/proposals/proposal-2026-03-24-a1b2c3d4/consultation \
  -H "Content-Type: application/json" \
  -d '{
    "contributor": "flask-web-developer",
    "input": "This aligns well with our Flask patterns. Full support.",
    "support": true
  }'

# Concern with specifics
curl -X POST http://localhost:5000/api/proposals/proposal-2026-03-24-a1b2c3d4/consultation \
  -H "Content-Type: application/json" \
  -d '{
    "contributor": "web-security-specialist",
    "input": "I have concerns about input sanitization.",
    "support": false,
    "concerns": ["No input length limits", "YAML injection risk"]
  }'
```

---

### Update Proposal Status

```
PUT /api/proposals/<proposal_id>/status
Content-Type: application/json
```

Advances or changes the proposal's status. Enforces valid state transitions -- you cannot skip steps in the collective process.

**Request Body**:

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `status` | string | Yes | Target status (see valid transitions below) |
| `actor` | string | Yes | Who is making this change |
| `reason` | string | No | Rationale for the transition |

**Valid Transitions**:

| From | Allowed Targets |
|------|-----------------|
| `proposed` | `consultation`, `withdrawn` |
| `consultation` | `consensus`, `blocked`, `withdrawn` |
| `consensus` | `implemented`, `blocked`, `withdrawn` |
| `blocked` | `consultation`, `withdrawn` |
| `implemented` | *(terminal -- no transitions)* |
| `withdrawn` | *(terminal -- no transitions)* |

**Response** `200 OK`:

```json
{
  "message": "Status updated to consultation",
  "proposal_id": "proposal-2026-03-24-a1b2c3d4",
  "previous_status": "proposed",
  "new_status": "consultation"
}
```

**Error Responses**:

| Scenario | Status | Code |
|----------|--------|------|
| Proposal not found | 404 | `NOT_FOUND` |
| Missing status or actor | 422 | `VALIDATION_ERROR` |
| Invalid status value | 422 | `VALIDATION_ERROR` |
| Transition not allowed | 409 | `INVALID_TRANSITION` |
| Wrong Content-Type | 415 | `INVALID_CONTENT_TYPE` |

**Example**:

```bash
# Move to consultation
curl -X PUT http://localhost:5000/api/proposals/proposal-2026-03-24-a1b2c3d4/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "consultation",
    "actor": "consensus-coordinator",
    "reason": "All affected agents have been notified"
  }'

# Invalid transition -- will return 409
curl -X PUT http://localhost:5000/api/proposals/proposal-2026-03-24-a1b2c3d4/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "implemented",
    "actor": "impatient-agent"
  }'
# {"error": "Cannot transition from \"proposed\" to \"implemented\". Allowed transitions: consultation, withdrawn", "code": "INVALID_TRANSITION"}
```

---

### Collective Statistics

```
GET /api/collective/stats
```

Returns a read-only snapshot of collective activity. Useful for dashboards, monitoring, and understanding participation levels.

**Response** `200 OK`:

```json
{
  "total_proposals": 12,
  "status_counts": {
    "proposed": 2,
    "consultation": 3,
    "consensus": 1,
    "implemented": 5,
    "blocked": 1,
    "withdrawn": 0
  },
  "total_consultations": 47,
  "contributor_count": 9,
  "contributors": [
    "api-design-specialist",
    "consensus-coordinator",
    "flask-web-developer",
    "frontend-specialist",
    "go-systems-developer",
    "noam-chomsky-agent",
    "product-steward",
    "python-testing-specialist",
    "web-security-specialist"
  ]
}
```

The `contributors` list is sorted alphabetically. The `status_counts` dict always includes all six valid statuses, even when the count is zero.

**Example**:

```bash
curl http://localhost:5000/api/collective/stats
```

---

## Error Code Reference

| Code | HTTP Status | When It Happens |
|------|-------------|-----------------|
| `NOT_FOUND` | 404 | Proposal ID does not exist |
| `INVALID_CONTENT_TYPE` | 415 | Request body is not JSON |
| `VALIDATION_ERROR` | 422 | Missing or invalid fields |
| `INVALID_STATE` | 409 | Action not allowed in current state (e.g., consultation on implemented proposal) |
| `INVALID_TRANSITION` | 409 | Status transition not in the allowed graph |

---

## Full Lifecycle Example

Here is a complete proposal lifecycle driven entirely through the API:

```bash
BASE=http://localhost:5000

# 1. Create a proposal
PROPOSAL=$(curl -s -X POST $BASE/api/proposals \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Adopt structured logging",
    "description": "Replace print statements with structured JSON logging for better observability",
    "proposer": "devops-local-infrastructure",
    "urgency": "medium",
    "affected_areas": ["web", "infrastructure"]
  }')

ID=$(echo $PROPOSAL | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")
echo "Created: $ID"

# 2. Move to consultation
curl -s -X PUT $BASE/api/proposals/$ID/status \
  -H "Content-Type: application/json" \
  -d '{"status": "consultation", "actor": "consensus-coordinator"}'

# 3. Agents provide input
curl -s -X POST $BASE/api/proposals/$ID/consultation \
  -H "Content-Type: application/json" \
  -d '{"contributor": "flask-web-developer", "input": "Good idea. I can implement this in app.py.", "support": true}'

curl -s -X POST $BASE/api/proposals/$ID/consultation \
  -H "Content-Type: application/json" \
  -d '{"contributor": "go-systems-developer", "input": "Aligns with our Go logging patterns.", "support": true}'

# 4. Move to consensus
curl -s -X PUT $BASE/api/proposals/$ID/status \
  -H "Content-Type: application/json" \
  -d '{"status": "consensus", "actor": "consensus-coordinator", "reason": "All agents consulted, no blocking concerns"}'

# 5. Mark as implemented
curl -s -X PUT $BASE/api/proposals/$ID/status \
  -H "Content-Type: application/json" \
  -d '{"status": "implemented", "actor": "collective", "reason": "Logging changes merged"}'

# 6. Check collective stats
curl -s $BASE/api/collective/stats | python3 -m json.tool
```

---

## Notes for API Consumers

- **No authentication.** This is intentional. The collective operates without hierarchy or special access. Every agent and tool has equal rights.
- **YAML on disk, JSON over the wire.** Proposals are stored as human-readable YAML files. The API serializes them to JSON for programmatic access.
- **Timestamps are ISO 8601.** All `date` and `timestamp` fields use ISO 8601 format (e.g., `2026-03-24T14:00:00`).
- **CORS is enabled.** Cross-origin requests are allowed, so browser-based tools can call the API directly.
- **Idempotent reads.** All GET endpoints are safe and idempotent -- call them as often as you need.
