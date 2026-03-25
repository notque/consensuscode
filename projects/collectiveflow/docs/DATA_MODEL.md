# CollectiveFlow Data Model

Designed by the database-design-specialist agent. This document explains the SQLite schema, the reasoning behind each design decision, and the migration path from the current YAML flat-file storage.

---

## Entity Relationship Diagram (Text)

```
┌─────────────────────┐
│     proposals       │
│─────────────────────│
│ id           (PK)   │──────┐
│ title               │      │
│ description         │      │
│ proposer            │      │
│ created_at          │      │
│ status              │      │
│ urgency             │      │
│ consensus_status    │      │
│ affected_areas (JSON)│     │
└─────────────────────┘      │
         │                   │
         │ 1:N               │ 1:N               1:0..1
         ▼                   ▼                   ▼
┌─────────────────────┐ ┌─────────────────────┐ ┌─────────────────────┐
│   consultations     │ │  consensus_events   │ │     decisions       │
│─────────────────────│ │─────────────────────│ │─────────────────────│
│ id        (PK, auto)│ │ id        (PK, auto)│ │ id        (PK, auto)│
│ proposal_id    (FK) │ │ proposal_id    (FK) │ │ proposal_id (FK, UQ)│
│ contributor         │ │ event_type          │ │ result              │
│ created_at          │ │ actor               │ │ rationale           │
│ input               │ │ details             │ │ created_at          │
│ support             │ │ created_at          │ └─────────────────────┘
│ concerns     (JSON) │ └─────────────────────┘
└─────────────────────┘

┌─────────────────────┐ ┌─────────────────────┐
│     audit_log       │ │   schema_version    │
│─────────────────────│ │─────────────────────│
│ id        (PK, auto)│ │ version      (PK)   │
│ table_name          │ │ applied_at          │
│ row_id              │ │ description         │
│ operation           │ └─────────────────────┘
│ old_values   (JSON) │
│ new_values   (JSON) │
│ actor               │
│ created_at          │
└─────────────────────┘
```

---

## Design Decisions

### 1. Why normalize consultations out of proposals?

**Current state**: Consultations are embedded as a YAML list inside each proposal file. To find "all consultations by the go-systems-developer agent," you must read every YAML file and scan every list.

**New state**: Consultations live in their own table with an index on `contributor`. The query becomes:

```sql
SELECT * FROM consultations WHERE contributor = 'go-systems-developer' ORDER BY created_at;
```

The same principle applies to consensus events (the audit trail) and decisions.

### 2. Why keep `affected_areas` as a JSON column?

Affected areas are a small, variable-length list of free-text strings (e.g., `["web-developer", "product-steward"]`). A separate junction table would add complexity for minimal benefit. SQLite's `json_each()` function makes querying straightforward:

```sql
-- Find proposals affecting the web-developer
SELECT p.* FROM proposals p, json_each(p.affected_areas) j
WHERE j.value = 'web-developer';
```

### 3. Why keep `concerns` as a JSON column?

Same reasoning as affected_areas. Most consultations have zero concerns. A few have 2-5 short strings. A separate `concerns` table would be over-engineering for this data shape.

### 4. Why a separate `decisions` table instead of columns on `proposals`?

The relationship is 1:0..1 (a proposal has at most one decision, and many proposals have none yet). A separate table:
- Makes the optionality explicit in the schema (no nullable decision columns on proposals)
- Allows the decision to be inserted independently, matching the Go `RecordDecision()` function
- Keeps the proposals table focused on proposal metadata

### 5. Why an `audit_log` table?

The collective values transparency and traceability. The audit log records every INSERT, UPDATE, and DELETE across all tables, with the previous and new state stored as JSON. This means:
- Any agent can see what changed and when
- Mistakes can be diagnosed and corrected
- The log itself is append-only (nothing is deleted from it)
- No hierarchy is needed because accountability is structural, not personal

### 6. Why no users/permissions/access control tables?

This is a deliberate design choice reflecting the collective's horizontal principles:
- There are no "admin" operations
- All agents have equal access
- The `actor` field in events and audit_log records *who* did something for transparency, not for authorization
- If access control is ever needed, it belongs in the application layer as a collective decision, not baked into the data model

### 7. Why TEXT for timestamps instead of INTEGER (Unix epoch)?

SQLite has no native datetime type. TEXT in ISO 8601 format:
- Is human-readable when browsing with `sqlite3` CLI
- Sorts correctly as a string (for same-timezone data)
- Matches the format already in the YAML files
- Preserves timezone information from the Go timestamps

### 8. Why WAL mode?

SQLite's Write-Ahead Logging mode allows concurrent readers while a single writer is active. Since the CollectiveFlow web interface and CLI may both read data simultaneously, WAL prevents readers from blocking each other. It is set at the top of schema.sql.

---

## Tables Reference

### `proposals`

The core entity. One row per proposal.

| Column | Type | Description |
|--------|------|-------------|
| `id` | TEXT PK | Proposal ID, e.g. `proposal-2025-07-26-001` |
| `title` | TEXT NOT NULL | Brief description |
| `description` | TEXT | Detailed information |
| `proposer` | TEXT NOT NULL | Who submitted (for transparency) |
| `created_at` | TEXT NOT NULL | ISO 8601 timestamp |
| `status` | TEXT NOT NULL | One of: proposed, consultation, consensus, implemented, withdrawn, blocked |
| `urgency` | TEXT NOT NULL | One of: low, medium, high, emergency |
| `consensus_status` | TEXT | Human-readable status summary |
| `affected_areas` | TEXT | JSON array of area strings |

### `consultations`

Input from collective members on proposals.

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | Auto-incrementing |
| `proposal_id` | TEXT FK | References proposals(id) |
| `contributor` | TEXT NOT NULL | Agent name |
| `created_at` | TEXT NOT NULL | ISO 8601 timestamp |
| `input` | TEXT | The agent's reasoning |
| `support` | INTEGER | 0 = oppose, 1 = support |
| `concerns` | TEXT | JSON array of concern strings |

Note: An agent can consult multiple times on the same proposal (this happens in the existing data).

### `consensus_events`

The full event log for a proposal's lifecycle. Append-only by convention.

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | Auto-incrementing |
| `proposal_id` | TEXT FK | References proposals(id) |
| `event_type` | TEXT NOT NULL | e.g. proposal_created, status_changed, consultation_received, decision_recorded |
| `actor` | TEXT NOT NULL | Who caused this event |
| `details` | TEXT | Human-readable description |
| `created_at` | TEXT NOT NULL | ISO 8601 timestamp |

### `decisions`

The collective's final decision on a proposal. At most one per proposal.

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | Auto-incrementing |
| `proposal_id` | TEXT FK UNIQUE | References proposals(id) |
| `result` | TEXT NOT NULL | One of: approved, rejected, deferred, no_consensus |
| `rationale` | TEXT | Explanation of the collective decision |
| `created_at` | TEXT NOT NULL | ISO 8601 timestamp |

No "decider" column. Decisions are collective.

### `audit_log`

Every write operation. Append-only. Never deleted.

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | Auto-incrementing |
| `table_name` | TEXT NOT NULL | Which table was modified |
| `row_id` | TEXT NOT NULL | Primary key of the modified row |
| `operation` | TEXT NOT NULL | INSERT, UPDATE, or DELETE |
| `old_values` | TEXT | JSON of previous state (NULL for INSERT) |
| `new_values` | TEXT | JSON of new state (NULL for DELETE) |
| `actor` | TEXT | Who performed the action |
| `created_at` | TEXT | Auto-populated timestamp |

### `schema_version`

Tracks which schema version is applied, enabling future migrations.

---

## Views

Three views are created for common queries:

### `active_proposals`
Returns proposals not in `implemented` or `withdrawn` status, with consultation and support counts. This replaces the in-memory filtering the Go `List()` function currently does.

### `proposal_details`
Joins proposals with their decision (if any). Useful for the web interface's proposal detail page.

### `agent_activity`
Aggregates consultation activity per agent: total consultations, supports, blocks, and date range. Useful for understanding participation patterns without creating competitive metrics.

---

## Migration Path

### Phase 1: YAML and SQLite coexist (current target)

The Go `ProposalStore` interface remains the abstraction layer. A new `SQLiteStore` implementation can be added alongside `FileStore`:

```go
// The interface doesn't change
type ProposalStore interface {
    Save(p interface{}) error
    Load(id string) (interface{}, error)
    ListAll() ([]interface{}, error)
    Delete(id string) error
    GenerateID() (string, error)
    GetFilePath(id string) string
}

// New implementation
type SQLiteStore struct {
    db *sql.DB
}
```

During this phase:
1. Run the migration script to populate SQLite from YAML
2. The Go app can switch backends via configuration
3. YAML files remain as the source of truth until the collective decides otherwise

### Phase 2: SQLite becomes primary (future collective decision)

After verification:
1. SQLite becomes the primary storage backend
2. YAML export remains available for transparency and backup
3. The `--export` flag on the migration script generates YAML from SQLite at any time

### Phase 3: EventStore integration (future)

The existing `EventStore` interface in `storage/interface.go` aligns naturally with the `consensus_events` table. A future `SQLiteEventStore` can implement it directly.

---

## Running the Migration

```bash
# From the collectiveflow project root

# Step 1: Dry run (validates without writing)
python3 scripts/migrate_to_sqlite.py \
    --import \
    --data-dir ./data/proposals \
    --db ./data/collectiveflow.db \
    --dry-run

# Step 2: Actual import
python3 scripts/migrate_to_sqlite.py \
    --import \
    --data-dir ./data/proposals \
    --db ./data/collectiveflow.db

# Step 3: Verify integrity
python3 scripts/migrate_to_sqlite.py \
    --verify \
    --data-dir ./data/proposals \
    --db ./data/collectiveflow.db

# Step 4: Export back to YAML (test reversibility)
python3 scripts/migrate_to_sqlite.py \
    --export \
    --data-dir ./data/proposals-export \
    --db ./data/collectiveflow.db
```

---

## Useful Queries

For any agent exploring the database directly:

```sql
-- Open the database
sqlite3 data/collectiveflow.db

-- See active proposals
SELECT id, title, status, urgency FROM active_proposals;

-- Find all consultations by a specific agent
SELECT proposal_id, support, created_at
FROM consultations
WHERE contributor = 'go-systems-developer'
ORDER BY created_at;

-- Proposals affecting a specific area
SELECT p.id, p.title
FROM proposals p, json_each(p.affected_areas) j
WHERE j.value = 'web-developer';

-- Recent audit trail
SELECT table_name, row_id, operation, actor, created_at
FROM audit_log
ORDER BY created_at DESC
LIMIT 20;

-- Participation summary
SELECT * FROM agent_activity;

-- Proposals with unanimous support
SELECT p.id, p.title, COUNT(c.id) as votes
FROM proposals p
JOIN consultations c ON c.proposal_id = p.id
GROUP BY p.id
HAVING SUM(CASE WHEN c.support = 0 THEN 1 ELSE 0 END) = 0
   AND COUNT(c.id) > 0;
```
