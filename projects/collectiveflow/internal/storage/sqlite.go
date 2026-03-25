package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// SQLiteStore implements ProposalStore using a SQLite database.
// This backend shares the same database schema used by the Python/Flask
// web interface, so both CLI and web can read/write the same data.
//
// Uses modernc.org/sqlite (pure Go, no CGO) for portability.
type SQLiteStore struct {
	db   *sql.DB
	mu   sync.RWMutex // Protects concurrent access
	path string       // Path to the database file

	// Sequence counter for generating IDs
	sequence int
}

// NewSQLiteStore creates a new SQLite-based storage backend.
// If the database file does not exist, the schema is applied automatically
// from the bundled SQL (matching scripts/schema.sql).
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	// Ensure the parent directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, &StorageError{
			Op:   "create database directory",
			Path: dir,
			Err:  err,
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, &StorageError{
			Op:   "open database",
			Path: dbPath,
			Err:  err,
		}
	}

	// Enable WAL mode and foreign keys (matching schema.sql pragmas)
	if _, err := db.Exec("PRAGMA journal_mode = WAL"); err != nil {
		db.Close()
		return nil, &StorageError{Op: "set WAL mode", Path: dbPath, Err: err}
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, &StorageError{Op: "enable foreign keys", Path: dbPath, Err: err}
	}

	s := &SQLiteStore{
		db:   db,
		path: dbPath,
	}

	// Apply schema if tables don't exist yet
	if err := s.ensureSchema(); err != nil {
		db.Close()
		return nil, err
	}

	// Initialize sequence counter based on existing proposals
	if err := s.initSequence(); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

// schema is the embedded DDL, matching scripts/schema.sql.
// Kept inline so the Go binary is self-contained.
const schema = `
CREATE TABLE IF NOT EXISTS proposals (
    id              TEXT PRIMARY KEY,
    title           TEXT NOT NULL,
    description     TEXT NOT NULL DEFAULT '',
    proposer        TEXT NOT NULL,
    created_at      TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'proposed'
                    CHECK (status IN (
                        'proposed', 'consultation', 'consensus',
                        'implemented', 'withdrawn', 'blocked'
                    )),
    urgency         TEXT NOT NULL DEFAULT 'medium'
                    CHECK (urgency IN ('low', 'medium', 'high', 'emergency')),
    consensus_status TEXT NOT NULL DEFAULT '',
    affected_areas  TEXT NOT NULL DEFAULT '[]'
);

CREATE INDEX IF NOT EXISTS idx_proposals_status     ON proposals (status);
CREATE INDEX IF NOT EXISTS idx_proposals_urgency    ON proposals (urgency);
CREATE INDEX IF NOT EXISTS idx_proposals_created_at ON proposals (created_at);
CREATE INDEX IF NOT EXISTS idx_proposals_proposer   ON proposals (proposer);

CREATE TABLE IF NOT EXISTS consultations (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    proposal_id     TEXT NOT NULL REFERENCES proposals(id) ON DELETE CASCADE,
    contributor     TEXT NOT NULL,
    created_at      TEXT NOT NULL,
    input           TEXT NOT NULL DEFAULT '',
    support         INTEGER NOT NULL DEFAULT 0,
    concerns        TEXT NOT NULL DEFAULT '[]'
);

CREATE INDEX IF NOT EXISTS idx_consultations_proposal   ON consultations (proposal_id);
CREATE INDEX IF NOT EXISTS idx_consultations_contributor ON consultations (contributor);
CREATE INDEX IF NOT EXISTS idx_consultations_created_at ON consultations (created_at);

CREATE TABLE IF NOT EXISTS consensus_events (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    proposal_id     TEXT NOT NULL REFERENCES proposals(id) ON DELETE CASCADE,
    event_type      TEXT NOT NULL,
    actor           TEXT NOT NULL,
    details         TEXT NOT NULL DEFAULT '',
    created_at      TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_proposal   ON consensus_events (proposal_id);
CREATE INDEX IF NOT EXISTS idx_events_type       ON consensus_events (event_type);
CREATE INDEX IF NOT EXISTS idx_events_actor      ON consensus_events (actor);
CREATE INDEX IF NOT EXISTS idx_events_created_at ON consensus_events (created_at);

CREATE TABLE IF NOT EXISTS decisions (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    proposal_id     TEXT NOT NULL UNIQUE REFERENCES proposals(id) ON DELETE CASCADE,
    result          TEXT NOT NULL
                    CHECK (result IN ('approved', 'rejected', 'deferred', 'no_consensus')),
    rationale       TEXT NOT NULL DEFAULT '',
    created_at      TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_decisions_proposal ON decisions (proposal_id);
CREATE INDEX IF NOT EXISTS idx_decisions_result   ON decisions (result);

CREATE TABLE IF NOT EXISTS audit_log (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    table_name      TEXT NOT NULL,
    row_id          TEXT NOT NULL,
    operation       TEXT NOT NULL
                    CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
    old_values      TEXT,
    new_values      TEXT,
    actor           TEXT NOT NULL DEFAULT 'system',
    created_at      TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_audit_table     ON audit_log (table_name);
CREATE INDEX IF NOT EXISTS idx_audit_row       ON audit_log (row_id);
CREATE INDEX IF NOT EXISTS idx_audit_operation ON audit_log (operation);
CREATE INDEX IF NOT EXISTS idx_audit_created   ON audit_log (created_at);

CREATE TABLE IF NOT EXISTS schema_version (
    version     INTEGER PRIMARY KEY,
    applied_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    description TEXT NOT NULL DEFAULT ''
);

INSERT OR IGNORE INTO schema_version (version, description)
VALUES (1, 'Initial schema: proposals, consultations, consensus_events, decisions, audit_log');
`

// ensureSchema applies the DDL if the proposals table does not exist.
func (s *SQLiteStore) ensureSchema() error {
	var name string
	err := s.db.QueryRow(
		"SELECT name FROM sqlite_master WHERE type='table' AND name='proposals'",
	).Scan(&name)

	if err == sql.ErrNoRows {
		// Tables don't exist yet; create them.
		if _, execErr := s.db.Exec(schema); execErr != nil {
			return &StorageError{Op: "apply schema", Path: s.path, Err: execErr}
		}
	} else if err != nil {
		return &StorageError{Op: "check schema", Path: s.path, Err: err}
	}

	return nil
}

// initSequence determines the next sequence number by looking at today's proposals.
func (s *SQLiteStore) initSequence() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	prefix := "proposal-" + today + "-"

	rows, err := s.db.Query(
		"SELECT id FROM proposals WHERE id LIKE ?",
		prefix+"%",
	)
	if err != nil {
		return &StorageError{Op: "init sequence", Path: s.path, Err: err}
	}
	defer rows.Close()

	maxSeq := 0
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			continue
		}
		// Parse sequence from "proposal-YYYY-MM-DD-NNN"
		parts := strings.Split(id, "-")
		if len(parts) >= 5 {
			if seq, err := strconv.Atoi(parts[4]); err == nil && seq > maxSeq {
				maxSeq = seq
			}
		}
	}
	if err := rows.Err(); err != nil {
		return &StorageError{Op: "init sequence scan", Path: s.path, Err: err}
	}

	s.sequence = maxSeq
	return nil
}

// Save stores or updates a proposal.
// The proposal is passed as interface{} to satisfy the ProposalStore interface.
// It is serialized via JSON round-trip to extract fields.
func (s *SQLiteStore) Save(p interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Convert to a map via JSON round-trip
	jsonData, err := json.Marshal(p)
	if err != nil {
		return &StorageError{Op: "marshal proposal", Err: err}
	}

	var m map[string]interface{}
	if err := json.Unmarshal(jsonData, &m); err != nil {
		return &StorageError{Op: "extract proposal fields", Err: err}
	}

	id, _ := m["id"].(string)
	if id == "" {
		return &StorageError{Op: "save proposal", Err: fmt.Errorf("proposal must have an ID")}
	}

	title, _ := m["title"].(string)
	description, _ := m["description"].(string)
	proposer, _ := m["proposer"].(string)
	status, _ := m["status"].(string)
	urgency, _ := m["urgency"].(string)
	consensusStatus, _ := m["consensus_status"].(string)

	// Handle date: the Proposal struct uses "date" while the DB uses "created_at"
	createdAt := ""
	if dateStr, ok := m["date"].(string); ok && dateStr != "" {
		createdAt = dateStr
	}
	if createdAt == "" {
		createdAt = time.Now().Format(time.RFC3339)
	}

	// affected_areas as JSON array
	affectedAreas := "[]"
	if areas, ok := m["affected_areas"]; ok && areas != nil {
		if areaBytes, err := json.Marshal(areas); err == nil {
			affectedAreas = string(areaBytes)
		}
	}

	tx, err := s.db.Begin()
	if err != nil {
		return &StorageError{Op: "begin transaction", Path: id, Err: err}
	}
	defer tx.Rollback()

	// Upsert proposal
	_, err = tx.Exec(`
		INSERT INTO proposals (id, title, description, proposer, created_at, status, urgency, consensus_status, affected_areas)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			description = excluded.description,
			proposer = excluded.proposer,
			created_at = excluded.created_at,
			status = excluded.status,
			urgency = excluded.urgency,
			consensus_status = excluded.consensus_status,
			affected_areas = excluded.affected_areas
	`, id, title, description, proposer, createdAt, status, urgency, consensusStatus, affectedAreas)
	if err != nil {
		return &StorageError{Op: "upsert proposal", Path: id, Err: err}
	}

	// Sync consensus_history -> consensus_events table
	if history, ok := m["consensus_history"]; ok {
		if events, ok := history.([]interface{}); ok {
			// Delete existing events for this proposal, then re-insert.
			// This ensures the Go round-trip produces the same result as
			// loading from YAML and saving back.
			if _, err := tx.Exec("DELETE FROM consensus_events WHERE proposal_id = ?", id); err != nil {
				return &StorageError{Op: "clear consensus events", Path: id, Err: err}
			}
			for _, e := range events {
				ev, ok := e.(map[string]interface{})
				if !ok {
					continue
				}
				eventType, _ := ev["event"].(string)
				actor, _ := ev["actor"].(string)
				details, _ := ev["details"].(string)
				timestamp := ""
				if ts, ok := ev["timestamp"].(string); ok {
					timestamp = ts
				}
				if timestamp == "" {
					timestamp = time.Now().Format(time.RFC3339)
				}
				_, err := tx.Exec(`
					INSERT INTO consensus_events (proposal_id, event_type, actor, details, created_at)
					VALUES (?, ?, ?, ?, ?)
				`, id, eventType, actor, details, timestamp)
				if err != nil {
					return &StorageError{Op: "insert consensus event", Path: id, Err: err}
				}
			}
		}
	}

	// Sync consultations
	if consults, ok := m["consultations"]; ok {
		if consultList, ok := consults.([]interface{}); ok {
			if _, err := tx.Exec("DELETE FROM consultations WHERE proposal_id = ?", id); err != nil {
				return &StorageError{Op: "clear consultations", Path: id, Err: err}
			}
			for _, c := range consultList {
				cm, ok := c.(map[string]interface{})
				if !ok {
					continue
				}
				contributor, _ := cm["contributor"].(string)
				input, _ := cm["input"].(string)
				support := false
				if s, ok := cm["support"].(bool); ok {
					support = s
				}
				supportInt := 0
				if support {
					supportInt = 1
				}
				timestamp := ""
				if ts, ok := cm["timestamp"].(string); ok {
					timestamp = ts
				}
				if timestamp == "" {
					timestamp = time.Now().Format(time.RFC3339)
				}
				concerns := "[]"
				if cc, ok := cm["concerns"]; ok && cc != nil {
					if cb, err := json.Marshal(cc); err == nil {
						concerns = string(cb)
					}
				}
				_, err := tx.Exec(`
					INSERT INTO consultations (proposal_id, contributor, created_at, input, support, concerns)
					VALUES (?, ?, ?, ?, ?, ?)
				`, id, contributor, timestamp, input, supportInt, concerns)
				if err != nil {
					return &StorageError{Op: "insert consultation", Path: id, Err: err}
				}
			}
		}
	}

	// Sync decision
	if dec, ok := m["decision"]; ok && dec != nil {
		dm, ok := dec.(map[string]interface{})
		if ok {
			result, _ := dm["result"].(string)
			rationale, _ := dm["rationale"].(string)
			timestamp := ""
			if ts, ok := dm["timestamp"].(string); ok {
				timestamp = ts
			}
			if timestamp == "" {
				timestamp = time.Now().Format(time.RFC3339)
			}
			// Upsert decision (unique on proposal_id)
			_, err := tx.Exec(`
				INSERT INTO decisions (proposal_id, result, rationale, created_at)
				VALUES (?, ?, ?, ?)
				ON CONFLICT(proposal_id) DO UPDATE SET
					result = excluded.result,
					rationale = excluded.rationale,
					created_at = excluded.created_at
			`, id, result, rationale, timestamp)
			if err != nil {
				return &StorageError{Op: "upsert decision", Path: id, Err: err}
			}
		}
	}

	// Audit log
	newValues, _ := json.Marshal(map[string]string{"id": id, "title": title})
	_, _ = tx.Exec(`
		INSERT INTO audit_log (table_name, row_id, operation, new_values, actor)
		VALUES ('proposals', ?, 'INSERT', ?, ?)
	`, id, string(newValues), proposer)

	if err := tx.Commit(); err != nil {
		return &StorageError{Op: "commit transaction", Path: id, Err: err}
	}

	return nil
}

// Load retrieves a proposal by ID.
// Returns a map[string]interface{} matching the shape the YAML store returns,
// so StorageAdapter can unmarshal it identically.
func (s *SQLiteStore) Load(id string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	row := s.db.QueryRow("SELECT * FROM proposals WHERE id = ?", id)

	var (
		pID, title, description, proposer, createdAt string
		status, urgency, consensusStatus, affectedAreas string
	)
	err := row.Scan(&pID, &title, &description, &proposer, &createdAt,
		&status, &urgency, &consensusStatus, &affectedAreas)
	if err == sql.ErrNoRows {
		return nil, &StorageError{
			Op:   "load proposal",
			Path: id,
			Err:  fmt.Errorf("proposal not found"),
		}
	}
	if err != nil {
		return nil, &StorageError{Op: "scan proposal", Path: id, Err: err}
	}

	result := map[string]interface{}{
		"id":               pID,
		"title":            title,
		"description":      description,
		"proposer":         proposer,
		"date":             createdAt,
		"status":           status,
		"urgency":          urgency,
		"consensus_status": consensusStatus,
	}

	// Parse affected_areas JSON
	var areas []interface{}
	if err := json.Unmarshal([]byte(affectedAreas), &areas); err == nil {
		result["affected_areas"] = areas
	} else {
		result["affected_areas"] = []interface{}{}
	}

	// Load consensus history
	eventRows, err := s.db.Query(`
		SELECT event_type, actor, details, created_at
		FROM consensus_events
		WHERE proposal_id = ?
		ORDER BY created_at ASC, id ASC
	`, id)
	if err == nil {
		defer eventRows.Close()
		var events []interface{}
		for eventRows.Next() {
			var eventType, actor, details, eventTime string
			if err := eventRows.Scan(&eventType, &actor, &details, &eventTime); err == nil {
				events = append(events, map[string]interface{}{
					"timestamp": eventTime,
					"event":     eventType,
					"actor":     actor,
					"details":   details,
				})
			}
		}
		result["consensus_history"] = events
	}

	// Load consultations
	consultRows, err := s.db.Query(`
		SELECT contributor, created_at, input, support, concerns
		FROM consultations
		WHERE proposal_id = ?
		ORDER BY created_at ASC, id ASC
	`, id)
	if err == nil {
		defer consultRows.Close()
		var consultations []interface{}
		for consultRows.Next() {
			var contributor, cTime, inputText, concernsJSON string
			var supportInt int
			if err := consultRows.Scan(&contributor, &cTime, &inputText, &supportInt, &concernsJSON); err == nil {
				c := map[string]interface{}{
					"contributor": contributor,
					"timestamp":   cTime,
					"input":       inputText,
					"support":     supportInt == 1,
				}
				var concerns []interface{}
				if err := json.Unmarshal([]byte(concernsJSON), &concerns); err == nil && len(concerns) > 0 {
					c["concerns"] = concerns
				}
				consultations = append(consultations, c)
			}
		}
		result["consultations"] = consultations
	}

	// Load decision
	var decResult, decRationale, decTime string
	decErr := s.db.QueryRow(`
		SELECT result, rationale, created_at
		FROM decisions
		WHERE proposal_id = ?
	`, id).Scan(&decResult, &decRationale, &decTime)
	if decErr == nil {
		result["decision"] = map[string]interface{}{
			"result":    decResult,
			"rationale": decRationale,
			"timestamp": decTime,
		}
	}

	return result, nil
}

// ListAll retrieves all proposals, sorted newest-first by created_at.
func (s *SQLiteStore) ListAll() ([]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT id FROM proposals ORDER BY created_at DESC")
	if err != nil {
		return nil, &StorageError{Op: "list proposals", Path: s.path, Err: err}
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			continue
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, &StorageError{Op: "list proposals scan", Path: s.path, Err: err}
	}

	// Load each proposal individually to get full nested data.
	// We release the read lock and re-acquire per-load because Load also locks.
	s.mu.RUnlock()
	defer s.mu.RLock()

	var proposals []interface{}
	for _, id := range ids {
		p, err := s.Load(id)
		if err != nil {
			continue // Skip entries that fail to load
		}
		proposals = append(proposals, p)
	}

	return proposals, nil
}

// Delete removes a proposal and all associated data (cascaded by foreign keys).
func (s *SQLiteStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check that the proposal exists first
	var exists string
	err := s.db.QueryRow("SELECT id FROM proposals WHERE id = ?", id).Scan(&exists)
	if err == sql.ErrNoRows {
		return &StorageError{
			Op:   "delete proposal",
			Path: id,
			Err:  fmt.Errorf("proposal not found"),
		}
	}
	if err != nil {
		return &StorageError{Op: "check proposal exists", Path: id, Err: err}
	}

	// Record in audit log before deleting
	_, _ = s.db.Exec(`
		INSERT INTO audit_log (table_name, row_id, operation, actor)
		VALUES ('proposals', ?, 'DELETE', 'system')
	`, id)

	// Delete (foreign key cascades handle related rows)
	if _, err := s.db.Exec("DELETE FROM proposals WHERE id = ?", id); err != nil {
		return &StorageError{Op: "delete proposal", Path: id, Err: err}
	}

	return nil
}

// GenerateID creates a new unique proposal ID in the format proposal-YYYY-MM-DD-NNN.
func (s *SQLiteStore) GenerateID() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sequence++
	today := time.Now().Format("2006-01-02")
	id := fmt.Sprintf("proposal-%s-%03d", today, s.sequence)

	// Verify uniqueness
	var exists string
	err := s.db.QueryRow("SELECT id FROM proposals WHERE id = ?", id).Scan(&exists)
	if err == nil {
		// ID already exists, try next
		s.mu.Unlock()
		defer s.mu.Lock()
		return s.GenerateID()
	}

	return id, nil
}

// GetFilePath returns the database file path for transparency.
func (s *SQLiteStore) GetFilePath(id string) string {
	return s.path + "#" + id
}

// Close closes the database connection.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
