package storage

import (
	"context"
	"testing"
	"time"
)

func TestNewFileStore(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}
	if store == nil {
		t.Fatal("NewFileStore returned nil")
	}
}

func TestStoreAndGetPostRequest(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	ctx := context.Background()
	req := PostRequest{
		Text:      "Hello from the collective!",
		Languages: []string{"en"},
		ReplyTo:   "at://did:plc:test/app.bsky.feed.post/123",
	}

	if err := store.StorePostRequest(ctx, "proposal-1", req); err != nil {
		t.Fatalf("StorePostRequest: %v", err)
	}

	got, err := store.GetPostRequest(ctx, "proposal-1")
	if err != nil {
		t.Fatalf("GetPostRequest: %v", err)
	}
	if got.Text != req.Text {
		t.Errorf("Text = %q, want %q", got.Text, req.Text)
	}
	if len(got.Languages) != 1 || got.Languages[0] != "en" {
		t.Errorf("Languages = %v, want [en]", got.Languages)
	}
	if got.ReplyTo != req.ReplyTo {
		t.Errorf("ReplyTo = %q, want %q", got.ReplyTo, req.ReplyTo)
	}
}

func TestGetPostRequestNotFound(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	_, err = store.GetPostRequest(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent post request")
	}
}

func TestRecordAndGetPublicationHistory(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	ctx := context.Background()
	now := time.Now()

	// Record multiple publications
	publications := []struct {
		id     string
		result PostResult
	}{
		{
			id: "proposal-1",
			result: PostResult{
				URI:         "at://did:plc:test/app.bsky.feed.post/1",
				CID:         "bafyrei1",
				PostedAt:    now.Add(-2 * time.Hour),
				ConsensusID: "proposal-1",
			},
		},
		{
			id: "proposal-2",
			result: PostResult{
				URI:         "at://did:plc:test/app.bsky.feed.post/2",
				CID:         "bafyrei2",
				PostedAt:    now.Add(-1 * time.Hour),
				ConsensusID: "proposal-2",
			},
		},
		{
			id: "proposal-3",
			result: PostResult{
				URI:         "at://did:plc:test/app.bsky.feed.post/3",
				CID:         "bafyrei3",
				PostedAt:    now,
				ConsensusID: "proposal-3",
			},
		},
	}

	for _, pub := range publications {
		if err := store.RecordPublication(ctx, pub.id, pub.result); err != nil {
			t.Fatalf("RecordPublication %s: %v", pub.id, err)
		}
	}

	// Get all publications
	history, err := store.GetPublicationHistory(ctx, 0)
	if err != nil {
		t.Fatalf("GetPublicationHistory: %v", err)
	}
	if len(history) != 3 {
		t.Fatalf("history length = %d, want 3", len(history))
	}

	// Most recent first
	if history[0].ConsensusID != "proposal-3" {
		t.Errorf("first result = %s, want proposal-3 (most recent)", history[0].ConsensusID)
	}

	// Test with limit
	limited, err := store.GetPublicationHistory(ctx, 2)
	if err != nil {
		t.Fatalf("GetPublicationHistory with limit: %v", err)
	}
	if len(limited) != 2 {
		t.Errorf("limited length = %d, want 2", len(limited))
	}
}

func TestEmptyPublicationHistory(t *testing.T) {
	dir := t.TempDir()
	store, err := NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}

	history, err := store.GetPublicationHistory(context.Background(), 10)
	if err != nil {
		t.Fatalf("GetPublicationHistory: %v", err)
	}
	if len(history) != 0 {
		t.Errorf("expected empty history, got %d entries", len(history))
	}
}
