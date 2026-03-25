package consensus

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileChecker implements ConsensusChecker using file-based storage.
// Proposals and decisions are stored as JSON files, keeping with the
// collective's local-only, simple-tools philosophy.
type FileChecker struct {
	baseDir string
	rules   ConsensusRules
	mu      sync.RWMutex
}

// DefaultRules implements ConsensusRules with sensible defaults for the collective.
type DefaultRules struct {
	MinParticipants int
	Timeout         time.Duration
}

// NewFileChecker creates a new file-based consensus checker.
func NewFileChecker(baseDir string, rules ConsensusRules) (*FileChecker, error) {
	dirs := []string{
		filepath.Join(baseDir, "proposals"),
		filepath.Join(baseDir, "decisions"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("create consensus directory %s: %w", dir, err)
		}
	}

	return &FileChecker{
		baseDir: baseDir,
		rules:   rules,
	}, nil
}

// NewDefaultRules creates consensus rules with the given minimum participants and timeout.
func NewDefaultRules(minParticipants int, timeout time.Duration) *DefaultRules {
	if minParticipants < 1 {
		minParticipants = 3
	}
	if timeout <= 0 {
		timeout = 24 * time.Hour
	}
	return &DefaultRules{
		MinParticipants: minParticipants,
		Timeout:         timeout,
	}
}

// EvaluateConsensus determines the decision status from the current votes.
// Consensus requires:
//   - At least MinParticipants votes
//   - No blocking votes
//   - At least one support vote
func (r *DefaultRules) EvaluateConsensus(votes map[string]Vote) (DecisionStatus, string) {
	if len(votes) < r.MinParticipants {
		return StatusPending, fmt.Sprintf("need %d votes, have %d", r.MinParticipants, len(votes))
	}

	supportCount := 0
	for _, vote := range votes {
		if vote.Position == PositionBlock {
			return StatusBlocked, fmt.Sprintf("blocked by %s: %s", vote.AgentID, vote.Reasoning)
		}
		if vote.Position == PositionSupport {
			supportCount++
		}
	}

	if supportCount == 0 {
		return StatusPending, "no support votes yet"
	}

	return StatusConsensus, "consensus reached"
}

// MinimumParticipation returns the minimum number of agents that must vote.
func (r *DefaultRules) MinimumParticipation() int {
	return r.MinParticipants
}

// ConsensusTimeout returns how long to wait before a proposal expires.
func (r *DefaultRules) ConsensusTimeout() time.Duration {
	return r.Timeout
}

// ProposePost creates a new proposal and its initial decision record.
func (c *FileChecker) ProposePost(ctx context.Context, proposal Proposal) (*Decision, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Save the proposal
	proposalPath := filepath.Join(c.baseDir, "proposals", proposal.ID+".json")
	if err := writeJSON(proposalPath, proposal); err != nil {
		return nil, fmt.Errorf("save proposal: %w", err)
	}

	// Create initial decision
	decision := &Decision{
		ID:         "decision-" + proposal.ID,
		ProposalID: proposal.ID,
		AgentVotes: make(map[string]Vote),
		Status:     StatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	decisionPath := filepath.Join(c.baseDir, "decisions", proposal.ID+".json")
	if err := writeJSON(decisionPath, decision); err != nil {
		return nil, fmt.Errorf("save decision: %w", err)
	}

	return decision, nil
}

// GetDecision retrieves the current decision state for a proposal.
func (c *FileChecker) GetDecision(_ context.Context, proposalID string) (*Decision, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	decisionPath := filepath.Join(c.baseDir, "decisions", proposalID+".json")
	var decision Decision
	if err := readJSON(decisionPath, &decision); err != nil {
		return nil, fmt.Errorf("get decision for %s: %w", proposalID, err)
	}

	return &decision, nil
}

// RecordVote adds or updates an agent's vote on a proposal.
func (c *FileChecker) RecordVote(_ context.Context, proposalID string, vote Vote) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	decisionPath := filepath.Join(c.baseDir, "decisions", proposalID+".json")
	var decision Decision
	if err := readJSON(decisionPath, &decision); err != nil {
		return fmt.Errorf("load decision for %s: %w", proposalID, err)
	}

	if decision.Status == StatusConsensus || decision.Status == StatusWithdrawn {
		return fmt.Errorf("cannot vote on %s proposal (status: %s)", proposalID, decision.Status)
	}

	decision.AgentVotes[vote.AgentID] = vote
	decision.UpdatedAt = time.Now()

	// Re-evaluate consensus after recording the vote
	status, _ := c.rules.EvaluateConsensus(decision.AgentVotes)
	decision.Status = status
	if status == StatusConsensus {
		now := time.Now()
		decision.ConsensusAt = &now
	}

	if err := writeJSON(decisionPath, decision); err != nil {
		return fmt.Errorf("save decision for %s: %w", proposalID, err)
	}

	return nil
}

// CheckConsensus evaluates whether consensus has been reached for a proposal.
// It reads the decision file directly rather than calling GetDecision to avoid
// nested RLock acquisition on the non-reentrant sync.RWMutex.
func (c *FileChecker) CheckConsensus(_ context.Context, proposalID string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	decisionPath := filepath.Join(c.baseDir, "decisions", proposalID+".json")
	var decision Decision
	if err := readJSON(decisionPath, &decision); err != nil {
		return false, fmt.Errorf("check consensus for %s: %w", proposalID, err)
	}
	return decision.Status == StatusConsensus, nil
}

// ListPendingProposals returns all proposals that are still awaiting consensus.
func (c *FileChecker) ListPendingProposals(_ context.Context) ([]Proposal, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	proposalDir := filepath.Join(c.baseDir, "proposals")
	entries, err := os.ReadDir(proposalDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read proposals directory: %w", err)
	}

	var pending []Proposal
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		var proposal Proposal
		path := filepath.Join(proposalDir, entry.Name())
		if err := readJSON(path, &proposal); err != nil {
			continue
		}

		// Check if there's a pending decision for this proposal
		decisionPath := filepath.Join(c.baseDir, "decisions", proposal.ID+".json")
		var decision Decision
		if err := readJSON(decisionPath, &decision); err != nil {
			continue
		}

		if decision.Status == StatusPending {
			pending = append(pending, proposal)
		}
	}

	return pending, nil
}

// GetProposal retrieves a single proposal by ID.
func (c *FileChecker) GetProposal(_ context.Context, proposalID string) (*Proposal, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	proposalPath := filepath.Join(c.baseDir, "proposals", proposalID+".json")
	var proposal Proposal
	if err := readJSON(proposalPath, &proposal); err != nil {
		return nil, fmt.Errorf("get proposal %s: %w", proposalID, err)
	}

	return &proposal, nil
}

// writeJSON marshals v to JSON and writes it to path.
func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// readJSON reads a JSON file at path and unmarshals it into v.
func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
