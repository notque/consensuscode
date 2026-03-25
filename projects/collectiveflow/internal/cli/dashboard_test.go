package cli

import (
	"testing"
	"time"

	"collectiveflow/internal/proposal"
)

func TestPrintStatusSummary(t *testing.T) {
	// Verify it does not panic on empty input.
	printStatusSummary(nil)
	printStatusSummary([]*proposal.Proposal{})
}

func TestPrintRecentActivity(t *testing.T) {
	now := time.Now()

	proposals := []*proposal.Proposal{
		{
			ID:    "p-001",
			Title: "Recent Proposal",
			Date:  now.Add(-1 * time.Hour),
			ConsensusHistory: []proposal.ConsensusEvent{
				{
					Timestamp: now.Add(-30 * time.Minute),
					Event:     "status_changed",
					Actor:     "test-agent",
				},
			},
		},
		{
			ID:    "p-old",
			Title: "Old Proposal",
			Date:  now.Add(-60 * 24 * time.Hour), // 60 days ago
		},
	}

	// Should not panic.
	printRecentActivity(proposals, 7)
}

func TestPrintRecentActivityEmpty(t *testing.T) {
	// Should not panic with no proposals.
	printRecentActivity(nil, 30)
}

func TestPrintMissingInput(t *testing.T) {
	proposals := []*proposal.Proposal{
		{
			ID:     "p-001",
			Title:  "Active Proposal",
			Status: proposal.StatusConsultation,
			Consultations: []proposal.Consultation{
				{Contributor: "go-systems-developer", Support: true},
			},
		},
		{
			ID:     "p-002",
			Title:  "Proposed Only",
			Status: proposal.StatusProposed,
		},
	}

	// Should not panic. The first proposal is in consultation with only
	// one contributor, so many agents should be listed as missing.
	printMissingInput(proposals)
}

func TestPrintMissingInputNoneActive(t *testing.T) {
	proposals := []*proposal.Proposal{
		{
			ID:     "p-001",
			Title:  "Already Implemented",
			Status: proposal.StatusImplemented,
		},
	}

	// Should not panic.
	printMissingInput(proposals)
}

func TestKnownAgentsNotEmpty(t *testing.T) {
	if len(knownAgents) == 0 {
		t.Error("knownAgents should not be empty")
	}
}
