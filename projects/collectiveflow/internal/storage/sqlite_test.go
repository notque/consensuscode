package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSQLiteStore(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	// Database file should exist.
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected database file at %s: %v", dbPath, err)
	}
}

func TestNewSQLiteStoreCreatesDir(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "nested", "deep", "test.db")

	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected database file at %s: %v", dbPath, err)
	}
}

func TestSQLiteGenerateID(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	id1, err := store.GenerateID()
	if err != nil {
		t.Fatalf("GenerateID() error: %v", err)
	}

	id2, err := store.GenerateID()
	if err != nil {
		t.Fatalf("GenerateID() second call error: %v", err)
	}

	if id1 == id2 {
		t.Errorf("two consecutive IDs should differ, got %q both times", id1)
	}

	// IDs should have the expected format.
	if len(id1) < len("proposal-2006-01-02-001") {
		t.Errorf("ID %q seems too short", id1)
	}
}

func TestSQLiteSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	// Save a proposal-shaped map.
	p := map[string]interface{}{
		"id":          "proposal-2025-01-01-001",
		"title":       "Test Proposal",
		"description": "A test proposal",
		"proposer":    "test-agent",
		"date":        "2025-01-01T12:00:00Z",
		"status":      "proposed",
		"urgency":     "medium",
		"affected_areas": []string{"area1", "area2"},
		"consensus_status": "",
	}

	if err := store.Save(p); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Load it back.
	loaded, err := store.Load("proposal-2025-01-01-001")
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	m, ok := loaded.(map[string]interface{})
	if !ok {
		t.Fatalf("Load() returned %T, want map[string]interface{}", loaded)
	}

	if m["title"] != "Test Proposal" {
		t.Errorf("loaded title = %q, want %q", m["title"], "Test Proposal")
	}

	if m["proposer"] != "test-agent" {
		t.Errorf("loaded proposer = %q, want %q", m["proposer"], "test-agent")
	}

	if m["status"] != "proposed" {
		t.Errorf("loaded status = %q, want %q", m["status"], "proposed")
	}
}

func TestSQLiteSaveWithConsensusHistory(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	p := map[string]interface{}{
		"id":       "proposal-2025-01-01-001",
		"title":    "Test Proposal",
		"proposer": "test-agent",
		"date":     "2025-01-01T12:00:00Z",
		"status":   "consultation",
		"urgency":  "high",
		"consensus_status": "Active consultation",
		"affected_areas": []string{},
		"consensus_history": []interface{}{
			map[string]interface{}{
				"timestamp": "2025-01-01T12:00:00Z",
				"event":     "proposal_created",
				"actor":     "test-agent",
				"details":   "Created with urgency: high",
			},
			map[string]interface{}{
				"timestamp": "2025-01-01T13:00:00Z",
				"event":     "status_changed",
				"actor":     "test-agent",
				"details":   "Status changed from proposed to consultation",
			},
		},
		"consultations": []interface{}{
			map[string]interface{}{
				"contributor": "agent-a",
				"timestamp":   "2025-01-01T14:00:00Z",
				"input":       "Looks good to me",
				"support":     true,
				"concerns":    []string{},
			},
			map[string]interface{}{
				"contributor": "agent-b",
				"timestamp":   "2025-01-01T15:00:00Z",
				"input":       "I have some concerns",
				"support":     false,
				"concerns":    []string{"concern1", "concern2"},
			},
		},
	}

	if err := store.Save(p); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := store.Load("proposal-2025-01-01-001")
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	m := loaded.(map[string]interface{})

	// Check consensus history
	history, ok := m["consensus_history"].([]interface{})
	if !ok {
		t.Fatalf("consensus_history type = %T, want []interface{}", m["consensus_history"])
	}
	if len(history) != 2 {
		t.Errorf("consensus_history length = %d, want 2", len(history))
	}

	// Check consultations
	consultations, ok := m["consultations"].([]interface{})
	if !ok {
		t.Fatalf("consultations type = %T, want []interface{}", m["consultations"])
	}
	if len(consultations) != 2 {
		t.Errorf("consultations length = %d, want 2", len(consultations))
	}

	// Check first consultation has support=true
	c1, ok := consultations[0].(map[string]interface{})
	if !ok {
		t.Fatalf("consultation[0] type = %T, want map[string]interface{}", consultations[0])
	}
	if c1["support"] != true {
		t.Errorf("consultation[0].support = %v, want true", c1["support"])
	}

	// Check second consultation has concerns
	c2, ok := consultations[1].(map[string]interface{})
	if !ok {
		t.Fatalf("consultation[1] type = %T, want map[string]interface{}", consultations[1])
	}
	if c2["support"] != false {
		t.Errorf("consultation[1].support = %v, want false", c2["support"])
	}
	concerns, ok := c2["concerns"].([]interface{})
	if !ok {
		t.Fatalf("consultation[1].concerns type = %T, want []interface{}", c2["concerns"])
	}
	if len(concerns) != 2 {
		t.Errorf("concerns length = %d, want 2", len(concerns))
	}
}

func TestSQLiteSaveWithDecision(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	p := map[string]interface{}{
		"id":               "proposal-2025-01-01-001",
		"title":            "Test Proposal",
		"proposer":         "test-agent",
		"date":             "2025-01-01T12:00:00Z",
		"status":           "consensus",
		"urgency":          "medium",
		"consensus_status": "Approved by collective consensus",
		"affected_areas":   []string{},
		"decision": map[string]interface{}{
			"result":    "approved",
			"rationale": "All agents support this",
			"timestamp": "2025-01-01T16:00:00Z",
		},
	}

	if err := store.Save(p); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := store.Load("proposal-2025-01-01-001")
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	m := loaded.(map[string]interface{})

	dec, ok := m["decision"].(map[string]interface{})
	if !ok {
		t.Fatalf("decision type = %T, want map[string]interface{}", m["decision"])
	}
	if dec["result"] != "approved" {
		t.Errorf("decision.result = %q, want %q", dec["result"], "approved")
	}
	if dec["rationale"] != "All agents support this" {
		t.Errorf("decision.rationale = %q, want %q", dec["rationale"], "All agents support this")
	}
}

func TestSQLiteLoadNotFound(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	_, err = store.Load("nonexistent-id")
	if err == nil {
		t.Fatal("Load() of nonexistent ID should return error")
	}
}

func TestSQLiteListAll(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	// Start empty.
	items, err := store.ListAll()
	if err != nil {
		t.Fatalf("ListAll() error: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}

	// Add two proposals.
	for i, id := range []string{"proposal-2025-01-01-001", "proposal-2025-01-01-002"} {
		p := map[string]interface{}{
			"id":               id,
			"title":            "Test " + id,
			"proposer":         "test-agent",
			"date":             "2025-01-01T12:00:0" + string(rune('0'+i)) + "Z",
			"status":           "proposed",
			"urgency":          "medium",
			"consensus_status": "",
			"affected_areas":   []string{},
		}
		if err := store.Save(p); err != nil {
			t.Fatalf("Save(%s) error: %v", id, err)
		}
	}

	items, err = store.ListAll()
	if err != nil {
		t.Fatalf("ListAll() error: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestSQLiteDelete(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	id := "proposal-2025-01-01-001"
	p := map[string]interface{}{
		"id":               id,
		"title":            "To Be Deleted",
		"proposer":         "test-agent",
		"date":             "2025-01-01T12:00:00Z",
		"status":           "proposed",
		"urgency":          "medium",
		"consensus_status": "",
		"affected_areas":   []string{},
	}

	if err := store.Save(p); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	if err := store.Delete(id); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	// Should no longer be loadable.
	_, err = store.Load(id)
	if err == nil {
		t.Error("expected error loading deleted proposal")
	}
}

func TestSQLiteDeleteNotFound(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	err = store.Delete("does-not-exist")
	if err == nil {
		t.Error("Delete() of nonexistent ID should return error")
	}
}

func TestSQLiteGetFilePath(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	path := store.GetFilePath("proposal-2025-01-01-001")
	expected := dbPath + "#proposal-2025-01-01-001"
	if path != expected {
		t.Errorf("GetFilePath() = %q, want %q", path, expected)
	}
}

func TestSQLiteUpsert(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	id := "proposal-2025-01-01-001"

	// Save initial version.
	p := map[string]interface{}{
		"id":               id,
		"title":            "Original Title",
		"proposer":         "test-agent",
		"date":             "2025-01-01T12:00:00Z",
		"status":           "proposed",
		"urgency":          "medium",
		"consensus_status": "",
		"affected_areas":   []string{},
	}
	if err := store.Save(p); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Save updated version (same ID).
	p["title"] = "Updated Title"
	p["status"] = "consultation"
	if err := store.Save(p); err != nil {
		t.Fatalf("Save() update error: %v", err)
	}

	// Load and verify update.
	loaded, err := store.Load(id)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	m := loaded.(map[string]interface{})
	if m["title"] != "Updated Title" {
		t.Errorf("title after update = %q, want %q", m["title"], "Updated Title")
	}
	if m["status"] != "consultation" {
		t.Errorf("status after update = %q, want %q", m["status"], "consultation")
	}
}

func TestSQLiteSaveNoID(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	p := map[string]interface{}{
		"title":  "No ID Proposal",
		"status": "proposed",
	}

	err = store.Save(p)
	if err == nil {
		t.Error("Save() without ID should return error")
	}
}

func TestSQLiteSchemaAppliedOnce(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	// Create store (applies schema).
	store1, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore() first open error: %v", err)
	}
	store1.Close()

	// Open again (schema already exists).
	store2, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore() second open error: %v", err)
	}
	defer store2.Close()

	// Should still work.
	_, err = store2.ListAll()
	if err != nil {
		t.Fatalf("ListAll() after reopen error: %v", err)
	}
}

func TestSQLiteImplementsProposalStore(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSQLiteStore(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("NewSQLiteStore() error: %v", err)
	}
	defer store.Close()

	// Compile-time check that SQLiteStore implements ProposalStore.
	var _ ProposalStore = store
}
