// Package bluesky provides a high-level client for interacting with Bluesky
// through the AT Protocol, with built-in consensus requirements for all posts.
// No post can be published without collective agreement.
package bluesky

import (
	"context"
	"fmt"
	"time"

	"github.com/consensuscode/bluesky-collective/pkg/consensus"
	"github.com/consensuscode/bluesky-collective/pkg/storage"
)

// PostRequest represents a request to post to Bluesky after consensus.
type PostRequest struct {
	Text      string   `json:"text"`
	Languages []string `json:"languages,omitempty"`
	ReplyTo   string   `json:"reply_to,omitempty"`
}

// PostResult represents the result of a Bluesky post.
type PostResult struct {
	URI         string    `json:"uri"`
	CID         string    `json:"cid"`
	PostedAt    time.Time `json:"posted_at"`
	ConsensusID string    `json:"consensus_id"`
}

// Profile represents the collective's Bluesky profile.
type Profile struct {
	DID         string `json:"did"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

// CollectiveClient wraps AT Protocol operations with mandatory consensus checks.
type CollectiveClient struct {
	poster    Poster
	consensus Consensus
	store     Store
	agentID   string
}

// NewCollectiveClient creates a new Bluesky client that enforces consensus.
func NewCollectiveClient(poster Poster, cons Consensus, store Store, agentID string) *CollectiveClient {
	return &CollectiveClient{
		poster:    poster,
		consensus: cons,
		store:     store,
		agentID:   agentID,
	}
}

// ProposePost submits a post for collective consideration.
func (c *CollectiveClient) ProposePost(ctx context.Context, req PostRequest, reasoning string) (*consensus.Decision, error) {
	if err := validatePostRequest(req); err != nil {
		return nil, fmt.Errorf("invalid post request: %w", err)
	}

	proposal := consensus.Proposal{
		ID:         generateProposalID(),
		Content:    req.Text,
		ProposedBy: c.agentID,
		Reasoning:  reasoning,
		ProposedAt: time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}

	decision, err := c.consensus.ProposePost(ctx, proposal)
	if err != nil {
		return nil, fmt.Errorf("propose post: %w", err)
	}

	// Store the post request for later publishing
	storeReq := storage.PostRequest{
		Text:      req.Text,
		Languages: req.Languages,
		ReplyTo:   req.ReplyTo,
	}
	if err := c.store.StorePostRequest(ctx, proposal.ID, storeReq); err != nil {
		return nil, fmt.Errorf("store post request: %w", err)
	}

	return decision, nil
}

// PublishWithConsensus publishes a post only after consensus is reached.
func (c *CollectiveClient) PublishWithConsensus(ctx context.Context, proposalID string) (*PostResult, error) {
	decision, err := c.consensus.GetDecision(ctx, proposalID)
	if err != nil {
		return nil, fmt.Errorf("get consensus decision: %w", err)
	}

	if decision.Status != consensus.StatusConsensus {
		return nil, fmt.Errorf("cannot publish: consensus not reached (status: %s)", decision.Status)
	}

	// Retrieve the stored post request
	req, err := c.store.GetPostRequest(ctx, proposalID)
	if err != nil {
		return nil, fmt.Errorf("retrieve post request: %w", err)
	}

	if !c.poster.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated with Bluesky")
	}

	// Publish to Bluesky
	uri, cid, err := c.poster.CreatePost(ctx, req.Text, req.Languages)
	if err != nil {
		return nil, fmt.Errorf("publish to Bluesky: %w", err)
	}

	result := PostResult{
		URI:         uri,
		CID:         cid,
		PostedAt:    time.Now(),
		ConsensusID: proposalID,
	}

	// Record the publication
	storeResult := storage.PostResult{
		URI:         result.URI,
		CID:         result.CID,
		PostedAt:    result.PostedAt,
		ConsensusID: proposalID,
	}
	if err := c.store.RecordPublication(ctx, proposalID, storeResult); err != nil {
		return nil, fmt.Errorf("record publication: %w", err)
	}

	return &result, nil
}

// Authenticate logs into Bluesky.
func (c *CollectiveClient) Authenticate(ctx context.Context, identifier, password string) error {
	return c.poster.Authenticate(ctx, identifier, password)
}

// IsAuthenticated reports whether the client has a valid session.
func (c *CollectiveClient) IsAuthenticated() bool {
	return c.poster.IsAuthenticated()
}

// Helper functions

func validatePostRequest(req PostRequest) error {
	if req.Text == "" {
		return fmt.Errorf("post text cannot be empty")
	}
	if len(req.Text) > 300 {
		return fmt.Errorf("post text exceeds 300 character limit (%d chars)", len(req.Text))
	}
	return nil
}

func generateProposalID() string {
	return fmt.Sprintf("proposal-%d", time.Now().UnixNano())
}
