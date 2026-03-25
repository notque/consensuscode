-- CollectiveFlow SQLite Schema
-- Designed by the database-design-specialist agent for the collective.
--
-- Design principles:
--   1. No access control tables. No users table. No permissions. (Horizontal principle.)
--   2. Proposals, consultations, and events are separate tables with foreign keys.
--   3. Every mutation is recorded in the audit_log for full traceability.
--   4. The schema supports the existing Go ProposalStore interface unchanged.
--   5. SQLite-specific: WAL mode for concurrent reads, foreign keys enforced.
--
-- To apply: sqlite3 collectiveflow.db < schema.sql

PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

-- ============================================================
-- PROPOSALS
-- The core entity. One row per proposal, mirroring the Go Proposal struct.
-- ============================================================
CREATE TABLE IF NOT EXISTS proposals (
    id              TEXT PRIMARY KEY,                -- e.g. "proposal-2025-07-26-001"
    title           TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    proposer        TEXT NOT NULL,                   -- who submitted (transparency, not authority)
    created_at      TEXT NOT NULL,                   -- ISO 8601 timestamp
    status          TEXT NOT NULL DEFAULT 'proposed'
                    CHECK (status IN (
                        'proposed', 'consultation', 'consensus',
                        'implemented', 'withdrawn', 'blocked'
                    )),
    urgency         TEXT NOT NULL DEFAULT 'medium'
                    CHECK (urgency IN ('low', 'medium', 'high', 'emergency')),
    consensus_status TEXT NOT NULL DEFAULT '',       -- human-readable summary
    -- Denormalized affected_areas stored as JSON array for simplicity.
    -- Querying by affected area uses json_each() in SQLite.
    affected_areas  TEXT NOT NULL DEFAULT '[]'       -- JSON array of strings
);

-- Indexes for the queries the CLI and web interface actually run.
CREATE INDEX IF NOT EXISTS idx_proposals_status     ON proposals (status);
CREATE INDEX IF NOT EXISTS idx_proposals_urgency    ON proposals (urgency);
CREATE INDEX IF NOT EXISTS idx_proposals_created_at ON proposals (created_at);
CREATE INDEX IF NOT EXISTS idx_proposals_proposer   ON proposals (proposer);

-- ============================================================
-- CONSULTATIONS
-- Input from collective members on a proposal.
-- Each row is one agent's consultation on one proposal.
-- An agent can consult multiple times on the same proposal (the data shows this).
-- ============================================================
CREATE TABLE IF NOT EXISTS consultations (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    proposal_id     TEXT NOT NULL REFERENCES proposals(id) ON DELETE CASCADE,
    contributor     TEXT NOT NULL,                   -- agent name
    created_at      TEXT NOT NULL,                   -- ISO 8601 timestamp
    input           TEXT NOT NULL DEFAULT '',        -- the agent's reasoning
    support         INTEGER NOT NULL DEFAULT 0,      -- 0 = false, 1 = true
    -- Concerns stored as JSON array. Most consultations have zero concerns;
    -- a few have 2-5. JSON array avoids a third table for a rare field.
    concerns        TEXT NOT NULL DEFAULT '[]'       -- JSON array of strings
);

CREATE INDEX IF NOT EXISTS idx_consultations_proposal   ON consultations (proposal_id);
CREATE INDEX IF NOT EXISTS idx_consultations_contributor ON consultations (contributor);
CREATE INDEX IF NOT EXISTS idx_consultations_created_at ON consultations (created_at);

-- ============================================================
-- CONSENSUS HISTORY (event log)
-- Every state change, consultation receipt, and decision is an event.
-- This is the audit trail. Append-only by convention.
-- ============================================================
CREATE TABLE IF NOT EXISTS consensus_events (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    proposal_id     TEXT NOT NULL REFERENCES proposals(id) ON DELETE CASCADE,
    event_type      TEXT NOT NULL,                   -- e.g. "proposal_created", "status_changed"
    actor           TEXT NOT NULL,                   -- who caused this event
    details         TEXT NOT NULL DEFAULT '',        -- human-readable description
    created_at      TEXT NOT NULL                    -- ISO 8601 timestamp
);

CREATE INDEX IF NOT EXISTS idx_events_proposal   ON consensus_events (proposal_id);
CREATE INDEX IF NOT EXISTS idx_events_type       ON consensus_events (event_type);
CREATE INDEX IF NOT EXISTS idx_events_actor      ON consensus_events (actor);
CREATE INDEX IF NOT EXISTS idx_events_created_at ON consensus_events (created_at);

-- ============================================================
-- DECISIONS
-- The collective's final decision on a proposal. At most one per proposal.
-- Kept separate from proposals to make the 1:0..1 relationship explicit.
-- ============================================================
CREATE TABLE IF NOT EXISTS decisions (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    proposal_id     TEXT NOT NULL UNIQUE REFERENCES proposals(id) ON DELETE CASCADE,
    result          TEXT NOT NULL
                    CHECK (result IN ('approved', 'rejected', 'deferred', 'no_consensus')),
    rationale       TEXT NOT NULL DEFAULT '',
    created_at      TEXT NOT NULL                    -- ISO 8601 timestamp
    -- No "decider" column. Decisions are collective.
);

CREATE INDEX IF NOT EXISTS idx_decisions_proposal ON decisions (proposal_id);
CREATE INDEX IF NOT EXISTS idx_decisions_result   ON decisions (result);

-- ============================================================
-- AUDIT LOG
-- Every write operation is recorded here for full traceability.
-- This table is append-only. Nothing is ever deleted from it.
-- ============================================================
CREATE TABLE IF NOT EXISTS audit_log (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    table_name      TEXT NOT NULL,                   -- which table was modified
    row_id          TEXT NOT NULL,                   -- primary key of the modified row
    operation       TEXT NOT NULL                    -- INSERT, UPDATE, DELETE
                    CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
    old_values      TEXT,                            -- JSON of previous state (NULL for INSERT)
    new_values      TEXT,                            -- JSON of new state (NULL for DELETE)
    actor           TEXT NOT NULL DEFAULT 'system',  -- who performed the action
    created_at      TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_audit_table     ON audit_log (table_name);
CREATE INDEX IF NOT EXISTS idx_audit_row       ON audit_log (row_id);
CREATE INDEX IF NOT EXISTS idx_audit_operation ON audit_log (operation);
CREATE INDEX IF NOT EXISTS idx_audit_created   ON audit_log (created_at);

-- ============================================================
-- SCHEMA METADATA
-- Tracks schema version so migration scripts can evolve safely.
-- ============================================================
CREATE TABLE IF NOT EXISTS schema_version (
    version     INTEGER PRIMARY KEY,
    applied_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    description TEXT NOT NULL DEFAULT ''
);

INSERT OR IGNORE INTO schema_version (version, description)
VALUES (1, 'Initial schema: proposals, consultations, consensus_events, decisions, audit_log');

-- ============================================================
-- VIEWS for common queries
-- These replace the in-memory filtering the Go code currently does.
-- ============================================================

-- Active proposals (not implemented or withdrawn)
CREATE VIEW IF NOT EXISTS active_proposals AS
SELECT p.*,
       (SELECT COUNT(*) FROM consultations c WHERE c.proposal_id = p.id) AS consultation_count,
       (SELECT COUNT(*) FROM consultations c WHERE c.proposal_id = p.id AND c.support = 1) AS support_count
FROM proposals p
WHERE p.status NOT IN ('implemented', 'withdrawn')
ORDER BY p.created_at DESC;

-- Full proposal detail with decision
CREATE VIEW IF NOT EXISTS proposal_details AS
SELECT p.*,
       d.result      AS decision_result,
       d.rationale   AS decision_rationale,
       d.created_at  AS decision_timestamp
FROM proposals p
LEFT JOIN decisions d ON d.proposal_id = p.id;

-- Consultation summary per agent
CREATE VIEW IF NOT EXISTS agent_activity AS
SELECT contributor,
       COUNT(*)                                     AS total_consultations,
       SUM(CASE WHEN support = 1 THEN 1 ELSE 0 END) AS supports,
       SUM(CASE WHEN support = 0 THEN 1 ELSE 0 END) AS blocks,
       MIN(created_at)                              AS first_consultation,
       MAX(created_at)                              AS last_consultation
FROM consultations
GROUP BY contributor;
