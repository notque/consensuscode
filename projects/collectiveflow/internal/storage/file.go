package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	
	"gopkg.in/yaml.v3"
)

// FileStore implements ProposalStore using the filesystem.
// Each proposal is stored as a YAML file for human readability.
// This aligns with our transparency and accessibility principles.
type FileStore struct {
	basePath string
	mu       sync.RWMutex // Protects concurrent access
	
	// Sequence counter for generating IDs
	sequence int
}

// NewFileStore creates a new file-based storage backend
func NewFileStore(basePath string) (*FileStore, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, &StorageError{
			Op:   "create storage directory",
			Path: basePath,
			Err:  err,
		}
	}
	
	fs := &FileStore{
		basePath: basePath,
	}
	
	// Initialize sequence counter based on existing proposals
	if err := fs.initSequence(); err != nil {
		return nil, err
	}
	
	return fs, nil
}

// initSequence determines the next sequence number by scanning existing files
func (fs *FileStore) initSequence() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	entries, err := os.ReadDir(fs.basePath)
	if err != nil {
		return &StorageError{
			Op:   "read directory",
			Path: fs.basePath,
			Err:  err,
		}
	}
	
	maxSeq := 0
	today := time.Now().Format("2006-01-02")
	
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		
		// Parse filename: proposal-YYYY-MM-DD-NNN.yaml
		name := strings.TrimSuffix(entry.Name(), ".yaml")
		parts := strings.Split(name, "-")
		
		if len(parts) < 5 { // proposal-YYYY-MM-DD-NNN
			continue
		}
		
		// Check if it's from today
		datePart := strings.Join(parts[1:4], "-")
		if datePart == today {
			// Extract sequence number
			if seq, err := strconv.Atoi(parts[4]); err == nil && seq > maxSeq {
				maxSeq = seq
			}
		}
	}
	
	fs.sequence = maxSeq
	return nil
}

// Save stores or updates a proposal
func (fs *FileStore) Save(p interface{}) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	// Extract ID from the proposal using reflection
	// This maintains storage independence while working with proposals
	proposalMap := make(map[string]interface{})
	
	// First convert to JSON to get a map
	jsonData, err := json.Marshal(p)
	if err != nil {
		return &StorageError{
			Op:  "marshal proposal",
			Err: err,
		}
	}
	
	if err := json.Unmarshal(jsonData, &proposalMap); err != nil {
		return &StorageError{
			Op:  "extract proposal fields",
			Err: err,
		}
	}
	
	// Get ID from the map
	id, ok := proposalMap["id"].(string)
	if !ok || id == "" {
		return &StorageError{
			Op:  "save proposal",
			Err: fmt.Errorf("proposal must have an ID"),
		}
	}
	
	// Create YAML representation
	data, err := yaml.Marshal(p)
	if err != nil {
		return &StorageError{
			Op:   "marshal proposal",
			Path: id,
			Err:  err,
		}
	}
	
	// Write to file
	filename := fs.getFilename(id)
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return &StorageError{
			Op:   "write file",
			Path: filename,
			Err:  err,
		}
	}
	
	// Also create a JSON copy for potential API access
	jsonFile := strings.TrimSuffix(filename, ".yaml") + ".json"
	formattedJSON, err := json.MarshalIndent(p, "", "  ")
	if err == nil {
		os.WriteFile(jsonFile, formattedJSON, 0644) // Ignore errors for JSON
	}
	
	return nil
}

// Load retrieves a proposal by ID
func (fs *FileStore) Load(id string) (interface{}, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	filename := fs.getFilename(id)
	
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, &StorageError{
				Op:   "load proposal",
				Path: id,
				Err:  fmt.Errorf("proposal not found"),
			}
		}
		return nil, &StorageError{
			Op:   "read file",
			Path: filename,
			Err:  err,
		}
	}
	
	// Return raw data - let the caller unmarshal to their type
	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, &StorageError{
			Op:   "unmarshal proposal",
			Path: id,
			Err:  err,
		}
	}
	
	return result, nil
}

// ListAll retrieves all proposals
func (fs *FileStore) ListAll() ([]interface{}, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	entries, err := os.ReadDir(fs.basePath)
	if err != nil {
		return nil, &StorageError{
			Op:   "list directory",
			Path: fs.basePath,
			Err:  err,
		}
	}
	
	var proposals []interface{}
	
	for _, entry := range entries {
		// Skip directories and non-YAML files
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		
		// Read and parse the file
		filename := filepath.Join(fs.basePath, entry.Name())
		data, err := os.ReadFile(filename)
		if err != nil {
			// Log error but continue with other files
			continue
		}
		
		var p map[string]interface{}
		if err := yaml.Unmarshal(data, &p); err != nil {
			// Log error but continue with other files
			continue
		}
		
		proposals = append(proposals, p)
	}
	
	// Sort by date (newest first) - extract date field for sorting
	sort.Slice(proposals, func(i, j int) bool {
		dateI, _ := proposals[i].(map[string]interface{})["date"].(string)
		dateJ, _ := proposals[j].(map[string]interface{})["date"].(string)
		return dateI > dateJ // String comparison works for ISO dates
	})
	
	return proposals, nil
}

// Delete removes a proposal (rarely used)
func (fs *FileStore) Delete(id string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	filename := fs.getFilename(id)
	
	// Create backup before deletion
	backupDir := filepath.Join(fs.basePath, ".deleted")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return &StorageError{
			Op:   "create backup directory",
			Path: backupDir,
			Err:  err,
		}
	}
	
	// Read the file for backup
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &StorageError{
				Op:   "delete proposal",
				Path: id,
				Err:  fmt.Errorf("proposal not found"),
			}
		}
		return &StorageError{
			Op:   "read file for backup",
			Path: filename,
			Err:  err,
		}
	}
	
	// Write backup with timestamp
	backupName := fmt.Sprintf("%s.deleted.%d.yaml", id, time.Now().Unix())
	backupPath := filepath.Join(backupDir, backupName)
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		// Continue with deletion even if backup fails
	}
	
	// Delete the original file
	if err := os.Remove(filename); err != nil {
		return &StorageError{
			Op:   "remove file",
			Path: filename,
			Err:  err,
		}
	}
	
	// Also remove JSON copy if it exists
	jsonFile := strings.TrimSuffix(filename, ".yaml") + ".json"
	os.Remove(jsonFile) // Ignore errors
	
	return nil
}

// GenerateID creates a new unique proposal ID
func (fs *FileStore) GenerateID() (string, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	
	// Increment sequence
	fs.sequence++
	
	// Format: proposal-YYYY-MM-DD-NNN
	today := time.Now().Format("2006-01-02")
	id := fmt.Sprintf("proposal-%s-%03d", today, fs.sequence)
	
	// Verify it doesn't exist (defensive programming)
	filename := fs.getFilename(id)
	if _, err := os.Stat(filename); err == nil {
		// File exists, try next sequence
		return fs.GenerateID()
	}
	
	return id, nil
}

// GetFilePath returns the storage path for transparency
func (fs *FileStore) GetFilePath(id string) string {
	return fs.getFilename(id)
}

// getFilename returns the full path for a proposal file
func (fs *FileStore) getFilename(id string) string {
	return filepath.Join(fs.basePath, id+".yaml")
}

// Backup creates a backup of all proposals
func (fs *FileStore) Backup(backupPath string) error {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	
	// Create backup directory
	backupDir := filepath.Join(backupPath, time.Now().Format("2006-01-02-150405"))
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return &StorageError{
			Op:   "create backup directory",
			Path: backupDir,
			Err:  err,
		}
	}
	
	// Copy all YAML files
	entries, err := os.ReadDir(fs.basePath)
	if err != nil {
		return &StorageError{
			Op:   "read directory for backup",
			Path: fs.basePath,
			Err:  err,
		}
	}
	
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		
		src := filepath.Join(fs.basePath, entry.Name())
		dst := filepath.Join(backupDir, entry.Name())
		
		if err := copyFile(src, dst); err != nil {
			// Log error but continue with other files
			continue
		}
	}
	
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	
	_, err = io.Copy(out, in)
	return err
}