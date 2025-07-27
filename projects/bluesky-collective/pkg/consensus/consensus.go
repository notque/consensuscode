// Package consensus provides interfaces and types for collective decision-making
// on Bluesky posts. All posts require consensus from the collective before
// being published to maintain horizontal accountability.
package consensus

import (
	"context"
	"time"
)

// Decision represents a collective decision on a proposed post
type Decision struct {
	ID          string              `json:"id"`
	ProposalID  string              `json:"proposal_id"`
	AgentVotes  map[string]Vote     `json:"agent_votes"`
	Status      DecisionStatus      `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	ConsensusAt *time.Time          `json:"consensus_at,omitempty"`
}

// Vote represents an individual agent's position on a proposal
type Vote struct {
	AgentID     string    `json:"agent_id"`
	Position    Position  `json:"position"`
	Reasoning   string    `json:"reasoning,omitempty"`
	Concerns    []string  `json:"concerns,omitempty"`
	VotedAt     time.Time `json:"voted_at"`
}

// Position represents an agent's stance on a proposal
type Position string

const (
	PositionSupport Position = "support"
	PositionBlock   Position = "block"
	PositionStandAside Position = "stand_aside"
	PositionAbstain Position = "abstain"
)

// DecisionStatus represents the current state of a consensus decision
type DecisionStatus string

const (
	StatusPending   DecisionStatus = "pending"
	StatusConsensus DecisionStatus = "consensus"
	StatusBlocked   DecisionStatus = "blocked"
	StatusWithdrawn DecisionStatus = "withdrawn"
)

// Proposal represents a proposed Bluesky post awaiting consensus
type Proposal struct {
	ID          string           `json:"id"`
	Content     string           `json:"content"`
	Images      []string         `json:"images,omitempty"`
	ProposedBy  string           `json:"proposed_by"`
	Reasoning   string           `json:"reasoning"`
	ProposedAt  time.Time        `json:"proposed_at"`
	ExpiresAt   time.Time        `json:"expires_at"`
}

// ConsensusChecker defines the interface for checking collective consensus
type ConsensusChecker interface {
	// ProposePost submits a new post for collective consideration
	ProposePost(ctx context.Context, proposal Proposal) (*Decision, error)
	
	// GetDecision retrieves the current consensus state for a proposal
	GetDecision(ctx context.Context, proposalID string) (*Decision, error)
	
	// RecordVote records an agent's vote on a proposal
	RecordVote(ctx context.Context, proposalID string, vote Vote) error
	
	// CheckConsensus evaluates if consensus has been reached
	CheckConsensus(ctx context.Context, proposalID string) (bool, error)
	
	// ListPendingProposals returns all proposals awaiting consensus
	ListPendingProposals(ctx context.Context) ([]Proposal, error)
}

// ConsensusRules defines the interface for consensus evaluation logic
type ConsensusRules interface {
	// EvaluateConsensus determines if votes constitute consensus
	EvaluateConsensus(votes map[string]Vote) (DecisionStatus, string)
	
	// MinimumParticipation returns the minimum number of agents needed
	MinimumParticipation() int
	
	// ConsensusTimeout returns how long to wait for consensus
	ConsensusTimeout() time.Duration
}