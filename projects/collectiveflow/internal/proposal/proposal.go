// Package proposal provides core data structures and operations for collective proposals.
// This package embodies horizontal decision-making principles - no proposal has
// inherent authority, and all proposals require collective consensus.
package proposal

import (
	"fmt"
	"time"
)

// Proposal represents a collective proposal for consideration.
// All proposals are equal - there are no "priority" or "admin" proposals.
type Proposal struct {
	// ID is the unique identifier for this proposal (format: proposal-YYYY-MM-DD-NNN)
	ID string `yaml:"id" json:"id"`

	// Title is a brief description of the proposal
	Title string `yaml:"title" json:"title"`

	// Description provides detailed information about the proposal
	Description string `yaml:"description" json:"description"`

	// Proposer identifies who submitted this (for transparency, not authority)
	Proposer string `yaml:"proposer" json:"proposer"`

	// Date when the proposal was created
	Date time.Time `yaml:"date" json:"date"`

	// Status of the proposal in the collective process
	Status ProposalStatus `yaml:"status" json:"status"`

	// Urgency level - for collective prioritization, not hierarchy
	Urgency UrgencyLevel `yaml:"urgency" json:"urgency"`

	// AffectedAreas lists which areas/domains this proposal impacts
	AffectedAreas []string `yaml:"affected_areas" json:"affected_areas"`

	// ConsensusStatus tracks the consensus-building process
	ConsensusStatus string `yaml:"consensus_status,omitempty" json:"consensus_status,omitempty"`

	// ConsensusHistory tracks all consensus-related events
	ConsensusHistory []ConsensusEvent `yaml:"consensus_history,omitempty" json:"consensus_history,omitempty"`

	// Consultations tracks input from different agents/members
	Consultations []Consultation `yaml:"consultations,omitempty" json:"consultations,omitempty"`

	// Decision records the collective's final decision
	Decision *Decision `yaml:"decision,omitempty" json:"decision,omitempty"`

	// FilePath where this proposal is stored (for transparency)
	FilePath string `yaml:"-" json:"-"`
}

// ProposalStatus represents the current state of a proposal
type ProposalStatus string

const (
	// StatusProposed - Initial state, awaiting consensus process
	StatusProposed ProposalStatus = "proposed"
	
	// StatusConsultation - Active consultation with collective members
	StatusConsultation ProposalStatus = "consultation"
	
	// StatusConsensus - Consensus reached, awaiting implementation
	StatusConsensus ProposalStatus = "consensus"
	
	// StatusImplemented - Proposal has been implemented
	StatusImplemented ProposalStatus = "implemented"
	
	// StatusWithdrawn - Proposer withdrew the proposal
	StatusWithdrawn ProposalStatus = "withdrawn"
	
	// StatusBlocked - Consensus could not be reached
	StatusBlocked ProposalStatus = "blocked"
)

// UrgencyLevel indicates time-sensitivity for collective consideration
type UrgencyLevel string

const (
	UrgencyLow       UrgencyLevel = "low"
	UrgencyMedium    UrgencyLevel = "medium"
	UrgencyHigh      UrgencyLevel = "high"
	UrgencyEmergency UrgencyLevel = "emergency"
)

// ConsensusEvent records a single event in the consensus process
type ConsensusEvent struct {
	Timestamp time.Time `yaml:"timestamp" json:"timestamp"`
	Event     string    `yaml:"event" json:"event"`
	Actor     string    `yaml:"actor" json:"actor"`
	Details   string    `yaml:"details,omitempty" json:"details,omitempty"`
}

// Consultation represents input from a collective member/agent
type Consultation struct {
	Contributor string    `yaml:"contributor" json:"contributor"`
	Timestamp   time.Time `yaml:"timestamp" json:"timestamp"`
	Input       string    `yaml:"input" json:"input"`
	Concerns    []string  `yaml:"concerns,omitempty" json:"concerns,omitempty"`
	Support     bool      `yaml:"support" json:"support"`
}

// Decision represents the collective's final decision on a proposal
type Decision struct {
	Result    DecisionResult `yaml:"result" json:"result"`
	Timestamp time.Time      `yaml:"timestamp" json:"timestamp"`
	Rationale string         `yaml:"rationale" json:"rationale"`
	// No "decider" field - decisions are collective
}

// DecisionResult represents possible collective decision outcomes
type DecisionResult string

const (
	DecisionApproved    DecisionResult = "approved"
	DecisionRejected    DecisionResult = "rejected"
	DecisionDeferred    DecisionResult = "deferred"
	DecisionNoConsensus DecisionResult = "no_consensus"
)

// New represents data for creating a new proposal
type New struct {
	Title         string
	Description   string
	Proposer      string
	Urgency       string
	AffectedAreas []string
	Date          time.Time
}

// Update represents fields that can be updated on a proposal
type Update struct {
	Title       string
	Description string
	Urgency     string
}

// ListFilter provides filtering options for listing proposals
type ListFilter struct {
	Status  string
	Urgency string
	Limit   int
	ShowAll bool
}

// Validate ensures the proposal data is valid
func (p *Proposal) Validate() error {
	if p.Title == "" {
		return fmt.Errorf("proposal title cannot be empty")
	}
	
	if p.Proposer == "" {
		return fmt.Errorf("proposer must be identified for transparency")
	}
	
	// Validate urgency level
	switch p.Urgency {
	case UrgencyLow, UrgencyMedium, UrgencyHigh, UrgencyEmergency:
		// Valid
	default:
		return fmt.Errorf("invalid urgency level: %s", p.Urgency)
	}
	
	// Validate status
	switch p.Status {
	case StatusProposed, StatusConsultation, StatusConsensus, 
	     StatusImplemented, StatusWithdrawn, StatusBlocked:
		// Valid
	default:
		return fmt.Errorf("invalid proposal status: %s", p.Status)
	}
	
	return nil
}

// CanTransitionTo checks if a status transition is valid
func (p *Proposal) CanTransitionTo(newStatus ProposalStatus) bool {
	// Define valid transitions (no hierarchical overrides)
	validTransitions := map[ProposalStatus][]ProposalStatus{
		StatusProposed:     {StatusConsultation, StatusWithdrawn},
		StatusConsultation: {StatusConsensus, StatusBlocked, StatusWithdrawn},
		StatusConsensus:    {StatusImplemented, StatusConsultation}, // Can return to consultation if concerns arise
		StatusImplemented:  {}, // Terminal state
		StatusWithdrawn:    {}, // Terminal state
		StatusBlocked:      {StatusConsultation}, // Can retry consensus
	}
	
	allowed, exists := validTransitions[p.Status]
	if !exists {
		return false
	}
	
	for _, status := range allowed {
		if status == newStatus {
			return true
		}
	}
	
	return false
}

// AddConsultation records input from a collective member
func (p *Proposal) AddConsultation(consultation Consultation) {
	p.Consultations = append(p.Consultations, consultation)
	
	// Add to consensus history
	event := ConsensusEvent{
		Timestamp: consultation.Timestamp,
		Event:     "consultation_received",
		Actor:     consultation.Contributor,
		Details:   fmt.Sprintf("Support: %v", consultation.Support),
	}
	p.ConsensusHistory = append(p.ConsensusHistory, event)
}

// HasUnanimousSupport checks if all consultations support the proposal
func (p *Proposal) HasUnanimousSupport() bool {
	if len(p.Consultations) == 0 {
		return false
	}
	
	for _, consultation := range p.Consultations {
		if !consultation.Support {
			return false
		}
	}
	
	return true
}

// GetBlockingConcerns returns all concerns that are blocking consensus
func (p *Proposal) GetBlockingConcerns() []string {
	var concerns []string
	
	for _, consultation := range p.Consultations {
		if !consultation.Support && len(consultation.Concerns) > 0 {
			concerns = append(concerns, consultation.Concerns...)
		}
	}
	
	return concerns
}