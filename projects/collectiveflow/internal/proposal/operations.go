package proposal

import (
	"fmt"
	"time"
	
	"collectiveflow/internal/storage"
)

// adapter is our storage adapter (initialized on first use)
var adapter *StorageAdapter

// initStore ensures the storage backend is initialized
func initStore() error {
	if adapter != nil {
		return nil
	}
	
	// Use file-based storage for now (per collective consensus)
	fileStore, err := storage.NewFileStore("./data/proposals")
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}
	
	adapter = NewStorageAdapter(fileStore)
	return nil
}

// Create creates a new proposal and stores it
func Create(new New) (*Proposal, error) {
	if err := initStore(); err != nil {
		return nil, err
	}
	
	// Generate ID based on date and sequence
	id, err := adapter.GenerateID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate proposal ID: %w", err)
	}
	
	// Parse urgency level
	urgency := UrgencyMedium // default
	switch new.Urgency {
	case "low":
		urgency = UrgencyLow
	case "medium":
		urgency = UrgencyMedium
	case "high":
		urgency = UrgencyHigh
	case "emergency":
		urgency = UrgencyEmergency
	}
	
	// Create the proposal
	proposal := &Proposal{
		ID:            id,
		Title:         new.Title,
		Description:   new.Description,
		Proposer:      new.Proposer,
		Date:          new.Date,
		Status:        StatusProposed,
		Urgency:       urgency,
		AffectedAreas: new.AffectedAreas,
		ConsensusHistory: []ConsensusEvent{
			{
				Timestamp: new.Date,
				Event:     "proposal_created",
				Actor:     new.Proposer,
				Details:   fmt.Sprintf("Created with urgency: %s", urgency),
			},
		},
	}
	
	// Validate before storing
	if err := proposal.Validate(); err != nil {
		return nil, fmt.Errorf("invalid proposal: %w", err)
	}
	
	// Store the proposal
	if err := adapter.Save(proposal); err != nil {
		return nil, fmt.Errorf("failed to store proposal: %w", err)
	}
	
	// Get the file path for transparency
	proposal.FilePath = adapter.GetFilePath(proposal.ID)
	
	return proposal, nil
}

// List retrieves proposals based on filter criteria
func List(filter ListFilter) ([]*Proposal, error) {
	if err := initStore(); err != nil {
		return nil, err
	}
	
	// Get all proposals
	allProposals, err := adapter.ListAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list proposals: %w", err)
	}
	
	// Apply filters
	var filtered []*Proposal
	for _, p := range allProposals {
		// Skip completed proposals unless ShowAll is true
		if !filter.ShowAll {
			if p.Status == StatusImplemented || p.Status == StatusWithdrawn {
				continue
			}
		}
		
		// Filter by status if specified
		if filter.Status != "" && string(p.Status) != filter.Status {
			continue
		}
		
		// Filter by urgency if specified
		if filter.Urgency != "" && string(p.Urgency) != filter.Urgency {
			continue
		}
		
		filtered = append(filtered, p)
		
		// Apply limit
		if filter.Limit > 0 && len(filtered) >= filter.Limit {
			break
		}
	}
	
	return filtered, nil
}

// Get retrieves a specific proposal by ID
func Get(proposalID string) (*Proposal, error) {
	if err := initStore(); err != nil {
		return nil, err
	}
	
	proposal, err := adapter.Load(proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to load proposal: %w", err)
	}
	
	// Set file path for transparency
	proposal.FilePath = adapter.GetFilePath(proposalID)
	
	return proposal, nil
}

// UpdateProposal updates an existing proposal
func UpdateProposal(proposalID string, updates Update) error {
	if err := initStore(); err != nil {
		return err
	}
	
	// Load existing proposal
	proposal, err := adapter.Load(proposalID)
	if err != nil {
		return fmt.Errorf("failed to load proposal: %w", err)
	}
	
	// Check if updates are allowed in current status
	if proposal.Status == StatusImplemented || proposal.Status == StatusWithdrawn {
		return fmt.Errorf("cannot update proposal in %s status", proposal.Status)
	}
	
	// Apply updates
	if updates.Title != "" {
		proposal.Title = updates.Title
	}
	
	if updates.Description != "" {
		proposal.Description = updates.Description
	}
	
	if updates.Urgency != "" {
		switch updates.Urgency {
		case "low":
			proposal.Urgency = UrgencyLow
		case "medium":
			proposal.Urgency = UrgencyMedium
		case "high":
			proposal.Urgency = UrgencyHigh
		case "emergency":
			proposal.Urgency = UrgencyEmergency
		default:
			return fmt.Errorf("invalid urgency level: %s", updates.Urgency)
		}
	}
	
	// Add update event to history
	event := ConsensusEvent{
		Timestamp: time.Now(),
		Event:     "proposal_updated",
		Actor:     "cli-user", // TODO: Get from context
		Details:   "Proposal details updated",
	}
	proposal.ConsensusHistory = append(proposal.ConsensusHistory, event)
	
	// Validate and save
	if err := proposal.Validate(); err != nil {
		return fmt.Errorf("invalid proposal after update: %w", err)
	}
	
	if err := adapter.Save(proposal); err != nil {
		return fmt.Errorf("failed to save updated proposal: %w", err)
	}
	
	return nil
}

// UpdateStatus changes the status of a proposal
func UpdateStatus(proposalID string, newStatus ProposalStatus, actor string) error {
	if err := initStore(); err != nil {
		return err
	}
	
	// Load existing proposal
	proposal, err := adapter.Load(proposalID)
	if err != nil {
		return fmt.Errorf("failed to load proposal: %w", err)
	}
	
	// Check if transition is valid
	if !proposal.CanTransitionTo(newStatus) {
		return fmt.Errorf("cannot transition from %s to %s", proposal.Status, newStatus)
	}
	
	// Update status
	oldStatus := proposal.Status
	proposal.Status = newStatus
	
	// Add status change event to history
	event := ConsensusEvent{
		Timestamp: time.Now(),
		Event:     "status_changed",
		Actor:     actor,
		Details:   fmt.Sprintf("Status changed from %s to %s", oldStatus, newStatus),
	}
	proposal.ConsensusHistory = append(proposal.ConsensusHistory, event)
	
	// Update consensus status if relevant
	switch newStatus {
	case StatusConsultation:
		proposal.ConsensusStatus = "Active consultation in progress"
	case StatusConsensus:
		proposal.ConsensusStatus = "Consensus reached"
	case StatusBlocked:
		proposal.ConsensusStatus = "Consensus blocked - concerns need addressing"
	}
	
	// Save updated proposal
	if err := adapter.Save(proposal); err != nil {
		return fmt.Errorf("failed to save proposal with new status: %w", err)
	}
	
	return nil
}

// AddConsultationInput records consultation input from a collective member
func AddConsultationInput(proposalID string, consultation Consultation) error {
	if err := initStore(); err != nil {
		return err
	}
	
	// Load existing proposal
	proposal, err := adapter.Load(proposalID)
	if err != nil {
		return fmt.Errorf("failed to load proposal: %w", err)
	}
	
	// Verify proposal is in consultation status
	if proposal.Status != StatusConsultation {
		return fmt.Errorf("proposal must be in consultation status to add input (current: %s)", proposal.Status)
	}
	
	// Add the consultation
	proposal.AddConsultation(consultation)
	
	// Save updated proposal
	if err := adapter.Save(proposal); err != nil {
		return fmt.Errorf("failed to save proposal with consultation: %w", err)
	}
	
	return nil
}

// RecordDecision records the collective's final decision on a proposal
func RecordDecision(proposalID string, decision Decision) error {
	if err := initStore(); err != nil {
		return err
	}
	
	// Load existing proposal
	proposal, err := adapter.Load(proposalID)
	if err != nil {
		return fmt.Errorf("failed to load proposal: %w", err)
	}
	
	// Verify proposal is in appropriate status
	if proposal.Status != StatusConsensus && proposal.Status != StatusConsultation {
		return fmt.Errorf("proposal must be in consensus or consultation status to record decision")
	}
	
	// Record the decision
	proposal.Decision = &decision
	
	// Update status based on decision
	switch decision.Result {
	case DecisionApproved:
		proposal.Status = StatusConsensus
		proposal.ConsensusStatus = "Approved by collective consensus"
	case DecisionRejected:
		proposal.Status = StatusWithdrawn
		proposal.ConsensusStatus = "Rejected by collective"
	case DecisionDeferred:
		proposal.Status = StatusProposed
		proposal.ConsensusStatus = "Deferred for future consideration"
	case DecisionNoConsensus:
		proposal.Status = StatusBlocked
		proposal.ConsensusStatus = "No consensus reached"
	}
	
	// Add decision event to history
	event := ConsensusEvent{
		Timestamp: decision.Timestamp,
		Event:     "decision_recorded",
		Actor:     "collective",
		Details:   fmt.Sprintf("Decision: %s", decision.Result),
	}
	proposal.ConsensusHistory = append(proposal.ConsensusHistory, event)
	
	// Save updated proposal
	if err := adapter.Save(proposal); err != nil {
		return fmt.Errorf("failed to save proposal with decision: %w", err)
	}
	
	return nil
}