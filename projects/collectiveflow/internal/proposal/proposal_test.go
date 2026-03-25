package proposal

import (
	"testing"
	"time"
)

func TestProposalValidate(t *testing.T) {
	tests := []struct {
		name    string
		prop    Proposal
		wantErr bool
	}{
		{
			name: "valid proposal",
			prop: Proposal{
				Title:    "Test proposal",
				Proposer: "test-agent",
				Status:   StatusProposed,
				Urgency:  UrgencyMedium,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			prop: Proposal{
				Proposer: "test-agent",
				Status:   StatusProposed,
				Urgency:  UrgencyMedium,
			},
			wantErr: true,
		},
		{
			name: "missing proposer",
			prop: Proposal{
				Title:   "Test proposal",
				Status:  StatusProposed,
				Urgency: UrgencyMedium,
			},
			wantErr: true,
		},
		{
			name: "invalid urgency",
			prop: Proposal{
				Title:    "Test proposal",
				Proposer: "test-agent",
				Status:   StatusProposed,
				Urgency:  "critical", // not a valid level
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			prop: Proposal{
				Title:    "Test proposal",
				Proposer: "test-agent",
				Status:   "unknown",
				Urgency:  UrgencyMedium,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.prop.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCanTransitionTo(t *testing.T) {
	tests := []struct {
		name      string
		from      ProposalStatus
		to        ProposalStatus
		wantAllow bool
	}{
		{"proposed to consultation", StatusProposed, StatusConsultation, true},
		{"proposed to withdrawn", StatusProposed, StatusWithdrawn, true},
		{"proposed to implemented", StatusProposed, StatusImplemented, false},
		{"consultation to consensus", StatusConsultation, StatusConsensus, true},
		{"consultation to blocked", StatusConsultation, StatusBlocked, true},
		{"consultation to withdrawn", StatusConsultation, StatusWithdrawn, true},
		{"consultation to proposed", StatusConsultation, StatusProposed, false},
		{"consensus to implemented", StatusConsensus, StatusImplemented, true},
		{"consensus to consultation", StatusConsensus, StatusConsultation, true},
		{"consensus to withdrawn", StatusConsensus, StatusWithdrawn, false},
		{"implemented to anything", StatusImplemented, StatusProposed, false},
		{"withdrawn to anything", StatusWithdrawn, StatusProposed, false},
		{"blocked to consultation", StatusBlocked, StatusConsultation, true},
		{"blocked to proposed", StatusBlocked, StatusProposed, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proposal{Status: tt.from}
			got := p.CanTransitionTo(tt.to)
			if got != tt.wantAllow {
				t.Errorf("CanTransitionTo(%s -> %s) = %v, want %v",
					tt.from, tt.to, got, tt.wantAllow)
			}
		})
	}
}

func TestHasUnanimousSupport(t *testing.T) {
	tests := []struct {
		name          string
		consultations []Consultation
		want          bool
	}{
		{
			name:          "no consultations",
			consultations: nil,
			want:          false,
		},
		{
			name: "all support",
			consultations: []Consultation{
				{Contributor: "a", Support: true},
				{Contributor: "b", Support: true},
			},
			want: true,
		},
		{
			name: "one dissent",
			consultations: []Consultation{
				{Contributor: "a", Support: true},
				{Contributor: "b", Support: false},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proposal{Consultations: tt.consultations}
			if got := p.HasUnanimousSupport(); got != tt.want {
				t.Errorf("HasUnanimousSupport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBlockingConcerns(t *testing.T) {
	p := &Proposal{
		Consultations: []Consultation{
			{Contributor: "a", Support: true},
			{Contributor: "b", Support: false, Concerns: []string{"too complex", "needs review"}},
			{Contributor: "c", Support: false, Concerns: []string{"timeline unclear"}},
			{Contributor: "d", Support: false}, // no concerns listed but not supporting
		},
	}

	concerns := p.GetBlockingConcerns()
	if len(concerns) != 3 {
		t.Errorf("GetBlockingConcerns() returned %d concerns, want 3", len(concerns))
	}

	expected := []string{"too complex", "needs review", "timeline unclear"}
	for i, want := range expected {
		if i >= len(concerns) {
			break
		}
		if concerns[i] != want {
			t.Errorf("concern[%d] = %q, want %q", i, concerns[i], want)
		}
	}
}

func TestAddConsultation(t *testing.T) {
	p := &Proposal{
		ID:     "test-001",
		Title:  "Test",
		Status: StatusConsultation,
	}

	c := Consultation{
		Contributor: "test-agent",
		Timestamp:   time.Now(),
		Input:       "Looks good to me.",
		Support:     true,
	}

	p.AddConsultation(c)

	if len(p.Consultations) != 1 {
		t.Fatalf("expected 1 consultation, got %d", len(p.Consultations))
	}
	if p.Consultations[0].Contributor != "test-agent" {
		t.Errorf("contributor = %q, want %q", p.Consultations[0].Contributor, "test-agent")
	}
	if len(p.ConsensusHistory) != 1 {
		t.Fatalf("expected 1 consensus history event, got %d", len(p.ConsensusHistory))
	}
	if p.ConsensusHistory[0].Event != "consultation_received" {
		t.Errorf("event = %q, want %q", p.ConsensusHistory[0].Event, "consultation_received")
	}
}
