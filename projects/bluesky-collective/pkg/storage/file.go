// Package storage provides file-based persistence for proposals, votes, and
// publication history. Data is stored as JSON files in a configurable directory,
// keeping with the collective's local-only, simple-tools philosophy.
package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// PostRequest mirrors the bluesky.PostRequest type to avoid circular imports.
type PostRequest struct {
	Text      string   `json:"text"`
	Images    [][]byte `json:"images,omitempty"`
	ReplyTo   string   `json:"reply_to,omitempty"`
	Quote     string   `json:"quote,omitempty"`
	Languages []string `json:"languages,omitempty"`
}

// PostResult mirrors the bluesky.PostResult type to avoid circular imports.
type PostResult struct {
	URI         string    `json:"uri"`
	CID         string    `json:"cid"`
	PostedAt    time.Time `json:"posted_at"`
	ConsensusID string    `json:"consensus_id"`
}

// FileStore persists data as JSON files in a directory tree.
type FileStore struct {
	baseDir string
	mu      sync.RWMutex
}

// NewFileStore creates a new file-based store rooted at baseDir.
// It creates the directory structure if it doesn't exist.
func NewFileStore(baseDir string) (*FileStore, error) {
	dirs := []string{
		filepath.Join(baseDir, "proposals"),
		filepath.Join(baseDir, "posts"),
		filepath.Join(baseDir, "publications"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("create storage directory %s: %w", dir, err)
		}
	}

	return &FileStore{baseDir: baseDir}, nil
}

// StorePostRequest saves a post request associated with a proposal ID.
func (s *FileStore) StorePostRequest(_ context.Context, proposalID string, req PostRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.baseDir, "posts", proposalID+".json")
	return writeJSON(path, req)
}

// GetPostRequest retrieves a stored post request by proposal ID.
func (s *FileStore) GetPostRequest(_ context.Context, proposalID string) (*PostRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := filepath.Join(s.baseDir, "posts", proposalID+".json")
	var req PostRequest
	if err := readJSON(path, &req); err != nil {
		return nil, fmt.Errorf("get post request %s: %w", proposalID, err)
	}
	return &req, nil
}

// RecordPublication saves the result of a published post.
func (s *FileStore) RecordPublication(_ context.Context, proposalID string, result PostResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.baseDir, "publications", proposalID+".json")
	return writeJSON(path, result)
}

// GetPublicationHistory returns the most recent publications, up to limit.
func (s *FileStore) GetPublicationHistory(_ context.Context, limit int) ([]PostResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dir := filepath.Join(s.baseDir, "publications")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read publications directory: %w", err)
	}

	var results []PostResult
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		var result PostResult
		path := filepath.Join(dir, entry.Name())
		if err := readJSON(path, &result); err != nil {
			continue // skip corrupted files
		}
		results = append(results, result)
	}

	// Sort by posted time, most recent first
	sort.Slice(results, func(i, j int) bool {
		return results[i].PostedAt.After(results[j].PostedAt)
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// writeJSON marshals v to JSON and writes it to path.
func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}
	return nil
}

// readJSON reads a JSON file at path and unmarshals it into v.
func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %s: %w", path, err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal json from %s: %w", path, err)
	}
	return nil
}
