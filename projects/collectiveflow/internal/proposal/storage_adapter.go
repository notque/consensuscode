package proposal

import (
	"encoding/json"
	"fmt"
	"collectiveflow/internal/storage"
)

// StorageAdapter wraps the generic storage interface to work with Proposal types
type StorageAdapter struct {
	store storage.ProposalStore
}

// NewStorageAdapter creates a new storage adapter
func NewStorageAdapter(store storage.ProposalStore) *StorageAdapter {
	return &StorageAdapter{store: store}
}

// Save stores a proposal
func (a *StorageAdapter) Save(p *Proposal) error {
	return a.store.Save(p)
}

// Load retrieves a proposal by ID
func (a *StorageAdapter) Load(id string) (*Proposal, error) {
	data, err := a.store.Load(id)
	if err != nil {
		return nil, err
	}
	
	// Convert the interface{} back to Proposal
	// First marshal to JSON then unmarshal to Proposal
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert loaded data: %w", err)
	}
	
	var proposal Proposal
	if err := json.Unmarshal(jsonData, &proposal); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proposal: %w", err)
	}
	
	return &proposal, nil
}

// ListAll retrieves all proposals
func (a *StorageAdapter) ListAll() ([]*Proposal, error) {
	dataList, err := a.store.ListAll()
	if err != nil {
		return nil, err
	}
	
	var proposals []*Proposal
	for _, data := range dataList {
		// Convert each interface{} to Proposal
		jsonData, err := json.Marshal(data)
		if err != nil {
			continue // Skip invalid entries
		}
		
		var proposal Proposal
		if err := json.Unmarshal(jsonData, &proposal); err != nil {
			continue // Skip invalid entries
		}
		
		proposals = append(proposals, &proposal)
	}
	
	return proposals, nil
}

// Delete removes a proposal
func (a *StorageAdapter) Delete(id string) error {
	return a.store.Delete(id)
}

// GenerateID creates a new unique proposal ID
func (a *StorageAdapter) GenerateID() (string, error) {
	return a.store.GenerateID()
}

// GetFilePath returns the storage path for transparency
func (a *StorageAdapter) GetFilePath(id string) string {
	return a.store.GetFilePath(id)
}