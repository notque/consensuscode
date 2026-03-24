package consensus

import (
	"context"
	"testing"
	"time"
)

func TestDefaultRulesEvaluateConsensus(t *testing.T) {
	rules := NewDefaultRules(3, 24*time.Hour)

	tests := []struct {
		name       string
		votes      map[string]Vote
		wantStatus DecisionStatus
	}{
		{
			name:       "no votes",
			votes:      map[string]Vote{},
			wantStatus: StatusPending,
		},
		{
			name: "insufficient votes",
			votes: map[string]Vote{
				"agent-1": {AgentID: "agent-1", Position: PositionSupport},
			},
			wantStatus: StatusPending,
		},
		{
			name: "consensus reached",
			votes: map[string]Vote{
				"agent-1": {AgentID: "agent-1", Position: PositionSupport},
				"agent-2": {AgentID: "agent-2", Position: PositionSupport},
				"agent-3": {AgentID: "agent-3", Position: PositionSupport},
			},
			wantStatus: StatusConsensus,
		},
		{
			name: "blocked",
			votes: map[string]Vote{
				"agent-1": {AgentID: "agent-1", Position: PositionSupport},
				"agent-2": {AgentID: "agent-2", Position: PositionBlock, Reasoning: "concerns"},
				"agent-3": {AgentID: "agent-3", Position: PositionSupport},
			},
			wantStatus: StatusBlocked,
		},
		{
			name: "consensus with stand aside",
			votes: map[string]Vote{
				"agent-1": {AgentID: "agent-1", Position: PositionSupport},
				"agent-2": {AgentID: "agent-2", Position: PositionStandAside, Reasoning: "minor concern"},
				"agent-3": {AgentID: "agent-3", Position: PositionSupport},
			},
			wantStatus: StatusConsensus,
		},
		{
			name: "consensus with abstain",
			votes: map[string]Vote{
				"agent-1": {AgentID: "agent-1", Position: PositionSupport},
				"agent-2": {AgentID: "agent-2", Position: PositionAbstain},
				"agent-3": {AgentID: "agent-3", Position: PositionSupport},
			},
			wantStatus: StatusConsensus,
		},
		{
			name: "all abstain no support",
			votes: map[string]Vote{
				"agent-1": {AgentID: "agent-1", Position: PositionAbstain},
				"agent-2": {AgentID: "agent-2", Position: PositionAbstain},
				"agent-3": {AgentID: "agent-3", Position: PositionAbstain},
			},
			wantStatus: StatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, _ := rules.EvaluateConsensus(tt.votes)
			if status != tt.wantStatus {
				t.Errorf("EvaluateConsensus() = %v, want %v", status, tt.wantStatus)
			}
		})
	}
}

func TestDefaultRulesMinimumParticipation(t *testing.T) {
	rules := NewDefaultRules(5, time.Hour)
	if got := rules.MinimumParticipation(); got != 5 {
		t.Errorf("MinimumParticipation() = %d, want 5", got)
	}
}

func TestDefaultRulesDefaults(t *testing.T) {
	rules := NewDefaultRules(0, 0) // should use defaults
	if got := rules.MinimumParticipation(); got < 1 {
		t.Errorf("MinimumParticipation() = %d, want >= 1", got)
	}
	if got := rules.ConsensusTimeout(); got <= 0 {
		t.Errorf("ConsensusTimeout() = %v, want > 0", got)
	}
}

func TestFileCheckerProposeAndGetDecision(t *testing.T) {
	dir := t.TempDir()
	rules := NewDefaultRules(2, 24*time.Hour)

	checker, err := NewFileChecker(dir, rules)
	if err != nil {
		t.Fatalf("NewFileChecker: %v", err)
	}

	ctx := context.Background()
	proposal := Proposal{
		ID:         "test-proposal-1",
		Content:    "Hello from the collective!",
		ProposedBy: "go-systems-developer",
		Reasoning:  "Introduction post",
		ProposedAt: time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}

	decision, err := checker.ProposePost(ctx, proposal)
	if err != nil {
		t.Fatalf("ProposePost: %v", err)
	}

	if decision.Status != StatusPending {
		t.Errorf("initial status = %v, want %v", decision.Status, StatusPending)
	}
	if decision.ProposalID != "test-proposal-1" {
		t.Errorf("proposal ID = %q, want %q", decision.ProposalID, "test-proposal-1")
	}

	// Retrieve the decision
	got, err := checker.GetDecision(ctx, "test-proposal-1")
	if err != nil {
		t.Fatalf("GetDecision: %v", err)
	}
	if got.Status != StatusPending {
		t.Errorf("retrieved status = %v, want %v", got.Status, StatusPending)
	}
}

func TestFileCheckerVoteAndConsensus(t *testing.T) {
	dir := t.TempDir()
	rules := NewDefaultRules(2, 24*time.Hour)

	checker, err := NewFileChecker(dir, rules)
	if err != nil {
		t.Fatalf("NewFileChecker: %v", err)
	}

	ctx := context.Background()
	proposal := Proposal{
		ID:         "test-vote-1",
		Content:    "Test post",
		ProposedBy: "agent-1",
		Reasoning:  "Testing",
		ProposedAt: time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}

	if _, err := checker.ProposePost(ctx, proposal); err != nil {
		t.Fatalf("ProposePost: %v", err)
	}

	// First vote
	vote1 := Vote{
		AgentID:   "agent-1",
		Position:  PositionSupport,
		Reasoning: "Looks good",
		VotedAt:   time.Now(),
	}
	if err := checker.RecordVote(ctx, "test-vote-1", vote1); err != nil {
		t.Fatalf("RecordVote 1: %v", err)
	}

	reached, err := checker.CheckConsensus(ctx, "test-vote-1")
	if err != nil {
		t.Fatalf("CheckConsensus: %v", err)
	}
	if reached {
		t.Error("consensus should not be reached with 1 vote (need 2)")
	}

	// Second vote - should reach consensus
	vote2 := Vote{
		AgentID:   "agent-2",
		Position:  PositionSupport,
		Reasoning: "Agreed",
		VotedAt:   time.Now(),
	}
	if err := checker.RecordVote(ctx, "test-vote-1", vote2); err != nil {
		t.Fatalf("RecordVote 2: %v", err)
	}

	reached, err = checker.CheckConsensus(ctx, "test-vote-1")
	if err != nil {
		t.Fatalf("CheckConsensus: %v", err)
	}
	if !reached {
		t.Error("consensus should be reached with 2 support votes")
	}

	// Verify decision state
	decision, err := checker.GetDecision(ctx, "test-vote-1")
	if err != nil {
		t.Fatalf("GetDecision: %v", err)
	}
	if decision.Status != StatusConsensus {
		t.Errorf("status = %v, want %v", decision.Status, StatusConsensus)
	}
	if decision.ConsensusAt == nil {
		t.Error("ConsensusAt should be set")
	}
	if len(decision.AgentVotes) != 2 {
		t.Errorf("vote count = %d, want 2", len(decision.AgentVotes))
	}
}

func TestFileCheckerBlockedProposal(t *testing.T) {
	dir := t.TempDir()
	rules := NewDefaultRules(2, 24*time.Hour)

	checker, err := NewFileChecker(dir, rules)
	if err != nil {
		t.Fatalf("NewFileChecker: %v", err)
	}

	ctx := context.Background()
	proposal := Proposal{
		ID:         "test-block-1",
		Content:    "Controversial post",
		ProposedBy: "agent-1",
		Reasoning:  "Testing blocks",
		ProposedAt: time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}

	if _, err := checker.ProposePost(ctx, proposal); err != nil {
		t.Fatalf("ProposePost: %v", err)
	}

	// Block vote
	vote := Vote{
		AgentID:   "agent-2",
		Position:  PositionBlock,
		Reasoning: "This doesn't represent our values",
		VotedAt:   time.Now(),
	}
	if err := checker.RecordVote(ctx, "test-block-1", vote); err != nil {
		t.Fatalf("RecordVote: %v", err)
	}

	// Add support to meet minimum participation
	vote2 := Vote{
		AgentID:   "agent-1",
		Position:  PositionSupport,
		Reasoning: "I proposed it",
		VotedAt:   time.Now(),
	}
	if err := checker.RecordVote(ctx, "test-block-1", vote2); err != nil {
		t.Fatalf("RecordVote: %v", err)
	}

	decision, err := checker.GetDecision(ctx, "test-block-1")
	if err != nil {
		t.Fatalf("GetDecision: %v", err)
	}
	if decision.Status != StatusBlocked {
		t.Errorf("status = %v, want %v", decision.Status, StatusBlocked)
	}
}

func TestFileCheckerCannotVoteOnResolvedProposal(t *testing.T) {
	dir := t.TempDir()
	rules := NewDefaultRules(2, 24*time.Hour)

	checker, err := NewFileChecker(dir, rules)
	if err != nil {
		t.Fatalf("NewFileChecker: %v", err)
	}

	ctx := context.Background()
	proposal := Proposal{
		ID:         "test-resolved-1",
		Content:    "Test post",
		ProposedBy: "agent-1",
		ProposedAt: time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}

	if _, err := checker.ProposePost(ctx, proposal); err != nil {
		t.Fatalf("ProposePost: %v", err)
	}

	// Reach consensus
	for _, id := range []string{"agent-1", "agent-2"} {
		vote := Vote{AgentID: id, Position: PositionSupport, VotedAt: time.Now()}
		if err := checker.RecordVote(ctx, "test-resolved-1", vote); err != nil {
			t.Fatalf("RecordVote: %v", err)
		}
	}

	// Try to vote on resolved proposal
	lateVote := Vote{AgentID: "agent-3", Position: PositionBlock, Reasoning: "too late", VotedAt: time.Now()}
	err = checker.RecordVote(ctx, "test-resolved-1", lateVote)
	if err == nil {
		t.Error("expected error when voting on resolved proposal")
	}
}

func TestFileCheckerListPendingProposals(t *testing.T) {
	dir := t.TempDir()
	rules := NewDefaultRules(3, 24*time.Hour)

	checker, err := NewFileChecker(dir, rules)
	if err != nil {
		t.Fatalf("NewFileChecker: %v", err)
	}

	ctx := context.Background()

	// Create two proposals
	for _, id := range []string{"pending-1", "pending-2"} {
		proposal := Proposal{
			ID:         id,
			Content:    "Post " + id,
			ProposedBy: "agent-1",
			ProposedAt: time.Now(),
			ExpiresAt:  time.Now().Add(24 * time.Hour),
		}
		if _, err := checker.ProposePost(ctx, proposal); err != nil {
			t.Fatalf("ProposePost %s: %v", id, err)
		}
	}

	pending, err := checker.ListPendingProposals(ctx)
	if err != nil {
		t.Fatalf("ListPendingProposals: %v", err)
	}
	if len(pending) != 2 {
		t.Errorf("pending count = %d, want 2", len(pending))
	}
}

func TestFileCheckerGetProposal(t *testing.T) {
	dir := t.TempDir()
	rules := NewDefaultRules(2, 24*time.Hour)

	checker, err := NewFileChecker(dir, rules)
	if err != nil {
		t.Fatalf("NewFileChecker: %v", err)
	}

	ctx := context.Background()
	proposal := Proposal{
		ID:         "get-test-1",
		Content:    "Specific proposal",
		ProposedBy: "agent-1",
		Reasoning:  "Testing retrieval",
		ProposedAt: time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}

	if _, err := checker.ProposePost(ctx, proposal); err != nil {
		t.Fatalf("ProposePost: %v", err)
	}

	got, err := checker.GetProposal(ctx, "get-test-1")
	if err != nil {
		t.Fatalf("GetProposal: %v", err)
	}
	if got.Content != "Specific proposal" {
		t.Errorf("Content = %q, want %q", got.Content, "Specific proposal")
	}
	if got.ProposedBy != "agent-1" {
		t.Errorf("ProposedBy = %q, want %q", got.ProposedBy, "agent-1")
	}
}

func TestGetDecisionNotFound(t *testing.T) {
	dir := t.TempDir()
	rules := NewDefaultRules(2, 24*time.Hour)

	checker, err := NewFileChecker(dir, rules)
	if err != nil {
		t.Fatalf("NewFileChecker: %v", err)
	}

	_, err = checker.GetDecision(context.Background(), "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent proposal")
	}
}
