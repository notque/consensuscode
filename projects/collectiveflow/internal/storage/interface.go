// Package storage provides the storage abstraction for CollectiveFlow.
// Following collective consensus, we start with file-based storage but
// maintain flexibility for future database integration if needed.
package storage

// ProposalStore defines the interface for proposal storage.
// This abstraction allows the collective to change storage backends
// through consensus without modifying application logic.
type ProposalStore interface {
	// Save stores or updates a proposal
	Save(p interface{}) error
	
	// Load retrieves a proposal by ID
	Load(id string) (interface{}, error)
	
	// ListAll retrieves all proposals (for filtering in memory)
	ListAll() ([]interface{}, error)
	
	// Delete removes a proposal (rarely used - proposals are historical record)
	Delete(id string) error
	
	// GenerateID creates a new unique proposal ID
	GenerateID() (string, error)
	
	// GetFilePath returns the storage path for transparency
	GetFilePath(id string) string
}

// EventStore will handle event sourcing in the future
type EventStore interface {
	// AppendEvent adds an event to the event log
	AppendEvent(event Event) error
	
	// GetEvents retrieves events for an entity
	GetEvents(entityID string) ([]Event, error)
	
	// GetAllEvents retrieves all events (for replay)
	GetAllEvents() ([]Event, error)
}

// Event represents an event in the event sourcing system
type Event struct {
	ID        string                 `json:"id"`
	Timestamp int64                  `json:"timestamp"`
	Type      string                 `json:"type"`
	EntityID  string                 `json:"entity_id"`
	Actor     string                 `json:"actor"`
	Data      map[string]interface{} `json:"data"`
}

// StorageError represents storage-specific errors
type StorageError struct {
	Op      string // Operation that failed
	Path    string // Path or ID involved
	Err     error  // Underlying error
}

func (e *StorageError) Error() string {
	if e.Path != "" {
		return e.Op + " " + e.Path + ": " + e.Err.Error()
	}
	return e.Op + ": " + e.Err.Error()
}

func (e *StorageError) Unwrap() error {
	return e.Err
}