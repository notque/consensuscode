package bluesky

import (
	"context"
	"fmt"
	"testing"

	"github.com/consensuscode/bluesky-collective/pkg/consensus"
	"github.com/consensuscode/bluesky-collective/pkg/storage"
)

// mockPoster implements Poster for testing.
type mockPoster struct {
	authenticated bool
	did           string
	posts         []string
	lastText      string
	lastLangs     []string
	createErr     error
}

func (m *mockPoster) Authenticate(_ context.Context, identifier, password string) error {
	if identifier == "" || password == "" {
		return fmt.Errorf("invalid credentials")
	}
	m.authenticated = true
	m.did = "did:plc:mock"
	return nil
}

func (m *mockPoster) CreatePost(_ context.Context, text string, langs []string) (string, string, error) {
	if m.createErr != nil {
		return "", "", m.createErr
	}
	m.lastText = text
	m.lastLangs = langs
	uri := fmt.Sprintf("at://did:plc:mock/app.bsky.feed.post/%d", len(m.posts)+1)
	cid := fmt.Sprintf("bafyrei%d", len(m.posts)+1)
	m.posts = append(m.posts, text)
	return uri, cid, nil
}

func (m *mockPoster) DeletePost(_ context.Context, uri string) error {
	return nil
}

func (m *mockPoster) IsAuthenticated() bool {
	return m.authenticated
}

func (m *mockPoster) GetDID() string {
	return m.did
}

// mockConsensus implements Consensus for testing.
type mockConsensus struct {
	proposals map[string]consensus.Proposal
	decisions map[string]*consensus.Decision
}

func newMockConsensus() *mockConsensus {
	return &mockConsensus{
		proposals: make(map[string]consensus.Proposal),
		decisions: make(map[string]*consensus.Decision),
	}
}

func (m *mockConsensus) ProposePost(_ context.Context, proposal consensus.Proposal) (*consensus.Decision, error) {
	m.proposals[proposal.ID] = proposal
	decision := &consensus.Decision{
		ID:         "decision-" + proposal.ID,
		ProposalID: proposal.ID,
		AgentVotes: make(map[string]consensus.Vote),
		Status:     consensus.StatusPending,
	}
	m.decisions[proposal.ID] = decision
	return decision, nil
}

func (m *mockConsensus) GetDecision(_ context.Context, proposalID string) (*consensus.Decision, error) {
	d, ok := m.decisions[proposalID]
	if !ok {
		return nil, fmt.Errorf("decision not found: %s", proposalID)
	}
	return d, nil
}

func (m *mockConsensus) RecordVote(_ context.Context, proposalID string, vote consensus.Vote) error {
	d, ok := m.decisions[proposalID]
	if !ok {
		return fmt.Errorf("decision not found: %s", proposalID)
	}
	d.AgentVotes[vote.AgentID] = vote
	return nil
}

func (m *mockConsensus) CheckConsensus(_ context.Context, proposalID string) (bool, error) {
	d, ok := m.decisions[proposalID]
	if !ok {
		return false, fmt.Errorf("decision not found: %s", proposalID)
	}
	return d.Status == consensus.StatusConsensus, nil
}

func (m *mockConsensus) ListPendingProposals(_ context.Context) ([]consensus.Proposal, error) {
	var result []consensus.Proposal
	for _, p := range m.proposals {
		result = append(result, p)
	}
	return result, nil
}

func (m *mockConsensus) GetProposal(_ context.Context, proposalID string) (*consensus.Proposal, error) {
	p, ok := m.proposals[proposalID]
	if !ok {
		return nil, fmt.Errorf("proposal not found: %s", proposalID)
	}
	return &p, nil
}

// mockStore implements Store for testing.
type mockStore struct {
	posts        map[string]storage.PostRequest
	publications map[string]storage.PostResult
}

func newMockStore() *mockStore {
	return &mockStore{
		posts:        make(map[string]storage.PostRequest),
		publications: make(map[string]storage.PostResult),
	}
}

func (m *mockStore) StorePostRequest(_ context.Context, proposalID string, req storage.PostRequest) error {
	m.posts[proposalID] = req
	return nil
}

func (m *mockStore) GetPostRequest(_ context.Context, proposalID string) (*storage.PostRequest, error) {
	req, ok := m.posts[proposalID]
	if !ok {
		return nil, fmt.Errorf("post request not found: %s", proposalID)
	}
	return &req, nil
}

func (m *mockStore) RecordPublication(_ context.Context, proposalID string, result storage.PostResult) error {
	m.publications[proposalID] = result
	return nil
}

func (m *mockStore) GetPublicationHistory(_ context.Context, limit int) ([]storage.PostResult, error) {
	var results []storage.PostResult
	for _, r := range m.publications {
		results = append(results, r)
	}
	return results, nil
}

func TestProposePost(t *testing.T) {
	poster := &mockPoster{}
	cons := newMockConsensus()
	store := newMockStore()
	client := NewCollectiveClient(poster, cons, store, "test-agent")

	ctx := context.Background()
	req := PostRequest{Text: "Hello from the collective!"}

	decision, err := client.ProposePost(ctx, req, "Introduction post")
	if err != nil {
		t.Fatalf("ProposePost: %v", err)
	}

	if decision.Status != consensus.StatusPending {
		t.Errorf("status = %v, want %v", decision.Status, consensus.StatusPending)
	}

	// Verify post was stored
	if len(store.posts) != 1 {
		t.Errorf("stored posts = %d, want 1", len(store.posts))
	}
}

func TestProposePostValidation(t *testing.T) {
	poster := &mockPoster{}
	cons := newMockConsensus()
	store := newMockStore()
	client := NewCollectiveClient(poster, cons, store, "test-agent")

	ctx := context.Background()

	// Empty text
	_, err := client.ProposePost(ctx, PostRequest{Text: ""}, "reason")
	if err == nil {
		t.Error("expected error for empty text")
	}

	// Text too long
	longText := ""
	for i := 0; i < 301; i++ {
		longText += "x"
	}
	_, err = client.ProposePost(ctx, PostRequest{Text: longText}, "reason")
	if err == nil {
		t.Error("expected error for text exceeding 300 chars")
	}
}

func TestPublishWithConsensus(t *testing.T) {
	poster := &mockPoster{authenticated: true, did: "did:plc:mock"}
	cons := newMockConsensus()
	store := newMockStore()
	client := NewCollectiveClient(poster, cons, store, "test-agent")

	ctx := context.Background()

	// Create a proposal and reach consensus
	req := PostRequest{Text: "Consensus post"}
	decision, err := client.ProposePost(ctx, req, "Testing publish")
	if err != nil {
		t.Fatalf("ProposePost: %v", err)
	}

	// Manually set consensus status
	cons.decisions[decision.ProposalID].Status = consensus.StatusConsensus

	result, err := client.PublishWithConsensus(ctx, decision.ProposalID)
	if err != nil {
		t.Fatalf("PublishWithConsensus: %v", err)
	}

	if result.URI == "" {
		t.Error("expected non-empty URI")
	}
	if result.ConsensusID != decision.ProposalID {
		t.Errorf("ConsensusID = %q, want %q", result.ConsensusID, decision.ProposalID)
	}

	// Verify the post text was sent
	if poster.lastText != "Consensus post" {
		t.Errorf("posted text = %q, want %q", poster.lastText, "Consensus post")
	}
}

func TestPublishWithoutConsensus(t *testing.T) {
	poster := &mockPoster{authenticated: true, did: "did:plc:mock"}
	cons := newMockConsensus()
	store := newMockStore()
	client := NewCollectiveClient(poster, cons, store, "test-agent")

	ctx := context.Background()

	req := PostRequest{Text: "No consensus yet"}
	decision, err := client.ProposePost(ctx, req, "Testing rejection")
	if err != nil {
		t.Fatalf("ProposePost: %v", err)
	}

	// Don't set consensus -- should fail
	_, err = client.PublishWithConsensus(ctx, decision.ProposalID)
	if err == nil {
		t.Error("expected error when publishing without consensus")
	}
}

func TestPublishWithoutAuth(t *testing.T) {
	poster := &mockPoster{authenticated: false}
	cons := newMockConsensus()
	store := newMockStore()
	client := NewCollectiveClient(poster, cons, store, "test-agent")

	ctx := context.Background()

	req := PostRequest{Text: "Unauthenticated"}
	decision, err := client.ProposePost(ctx, req, "Testing auth check")
	if err != nil {
		t.Fatalf("ProposePost: %v", err)
	}

	cons.decisions[decision.ProposalID].Status = consensus.StatusConsensus

	_, err = client.PublishWithConsensus(ctx, decision.ProposalID)
	if err == nil {
		t.Error("expected error when publishing without authentication")
	}
}

func TestValidatePostRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     PostRequest
		wantErr bool
	}{
		{"valid short post", PostRequest{Text: "Hello!"}, false},
		{"valid max length", PostRequest{Text: string(make([]byte, 300))}, false},
		{"empty text", PostRequest{Text: ""}, true},
		{"too long", PostRequest{Text: string(make([]byte, 301))}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePostRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePostRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsAuthenticated(t *testing.T) {
	poster := &mockPoster{}
	client := NewCollectiveClient(poster, newMockConsensus(), newMockStore(), "test")

	if client.IsAuthenticated() {
		t.Error("should not be authenticated initially")
	}

	poster.authenticated = true
	if !client.IsAuthenticated() {
		t.Error("should be authenticated after setting flag")
	}
}

func TestAuthenticate(t *testing.T) {
	poster := &mockPoster{}
	client := NewCollectiveClient(poster, newMockConsensus(), newMockStore(), "test")

	err := client.Authenticate(context.Background(), "user", "pass")
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if !client.IsAuthenticated() {
		t.Error("should be authenticated after Authenticate()")
	}
}
