package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileStore(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "proposals")

	store, err := NewFileStore(storePath)
	if err != nil {
		t.Fatalf("NewFileStore() error: %v", err)
	}

	// Directory should be created.
	info, err := os.Stat(storePath)
	if err != nil {
		t.Fatalf("expected directory to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %s to be a directory", storePath)
	}

	_ = store
}

func TestGenerateID(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore() error: %v", err)
	}

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

	// IDs should have the expected prefix.
	if len(id1) < len("proposal-2006-01-02-001") {
		t.Errorf("ID %q seems too short", id1)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore() error: %v", err)
	}

	// Save a proposal-shaped map.
	p := map[string]interface{}{
		"id":     "proposal-2025-01-01-001",
		"title":  "Test Proposal",
		"status": "proposed",
	}

	if err := store.Save(p); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// YAML file should exist.
	yamlPath := filepath.Join(dir, "proposal-2025-01-01-001.yaml")
	if _, err := os.Stat(yamlPath); err != nil {
		t.Fatalf("expected YAML file at %s: %v", yamlPath, err)
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
}

func TestLoadNotFound(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore() error: %v", err)
	}

	_, err = store.Load("nonexistent-id")
	if err == nil {
		t.Fatal("Load() of nonexistent ID should return error")
	}
}

func TestListAll(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore() error: %v", err)
	}

	// Start empty.
	items, err := store.ListAll()
	if err != nil {
		t.Fatalf("ListAll() error: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}

	// Add two proposals.
	for _, id := range []string{"proposal-2025-01-01-001", "proposal-2025-01-01-002"} {
		p := map[string]interface{}{
			"id":     id,
			"title":  "Test " + id,
			"status": "proposed",
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

func TestDelete(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore() error: %v", err)
	}

	id := "proposal-2025-01-01-001"
	p := map[string]interface{}{
		"id":     id,
		"title":  "To Be Deleted",
		"status": "proposed",
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

	// Backup should exist.
	backupDir := filepath.Join(dir, ".deleted")
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		t.Fatalf("expected backup directory: %v", err)
	}
	if len(entries) == 0 {
		t.Error("expected at least one backup file")
	}
}

func TestDeleteNotFound(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore() error: %v", err)
	}

	err = store.Delete("does-not-exist")
	if err == nil {
		t.Error("Delete() of nonexistent ID should return error")
	}
}

func TestGetFilePath(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore() error: %v", err)
	}

	path := store.GetFilePath("proposal-2025-01-01-001")
	expected := filepath.Join(dir, "proposal-2025-01-01-001.yaml")
	if path != expected {
		t.Errorf("GetFilePath() = %q, want %q", path, expected)
	}
}

func TestStorageError(t *testing.T) {
	err := &StorageError{
		Op:   "save",
		Path: "/tmp/test",
		Err:  os.ErrPermission,
	}

	msg := err.Error()
	if msg != "save /tmp/test: permission denied" {
		t.Errorf("Error() = %q, want %q", msg, "save /tmp/test: permission denied")
	}

	// Without path.
	err2 := &StorageError{Op: "init", Err: os.ErrNotExist}
	msg2 := err2.Error()
	if msg2 != "init: file does not exist" {
		t.Errorf("Error() = %q, want %q", msg2, "init: file does not exist")
	}

	// Unwrap.
	if err.Unwrap() != os.ErrPermission {
		t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), os.ErrPermission)
	}
}
