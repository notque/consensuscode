// Package bluesky provides a client for interacting with the Bluesky AT Protocol
// with built-in consensus requirements for all posts.
package bluesky

import (
	"context"
	"fmt"
	"time"

	"github.com/consensuscode/bluesky-collective/pkg/consensus"
)

// PostRequest represents a request to post to Bluesky after consensus
type PostRequest struct {
	Text      string   `json:"text"`
	Images    [][]byte `json:"images,omitempty"`
	ReplyTo   string   `json:"reply_to,omitempty"`
	Quote     string   `json:"quote,omitempty"`
	Languages []string `json:"languages,omitempty"`
}

// PostResult represents the result of a Bluesky post
type PostResult struct {
	URI       string    `json:"uri"`
	CID       string    `json:"cid"`
	PostedAt  time.Time `json:"posted_at"`
	ConsensusID string  `json:"consensus_id"`
}

// Client defines the interface for Bluesky operations with consensus
type Client interface {
	// ProposePost submits a post for collective consensus
	ProposePost(ctx context.Context, req PostRequest, reasoning string) (*consensus.Decision, error)
	
	// PublishWithConsensus publishes a post after consensus is reached
	PublishWithConsensus(ctx context.Context, proposalID string) (*PostResult, error)
	
	// GetProfile retrieves the collective's profile information
	GetProfile(ctx context.Context) (*Profile, error)
	
	// UpdateProfile updates profile after consensus
	UpdateProfile(ctx context.Context, profile Profile) error
}

// Profile represents the collective's Bluesky profile
type Profile struct {
	DID         string `json:"did"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Avatar      []byte `json:"avatar,omitempty"`
	Banner      []byte `json:"banner,omitempty"`
}

// ConsensusClient implements Client with mandatory consensus checks
type ConsensusClient struct {
	atpClient  ATPClient
	consensus  consensus.ConsensusChecker
	storage    Storage
}

// NewConsensusClient creates a new Bluesky client with consensus requirements
func NewConsensusClient(atpClient ATPClient, consensus consensus.ConsensusChecker, storage Storage) *ConsensusClient {
	return &ConsensusClient{
		atpClient: atpClient,
		consensus: consensus,
		storage:   storage,
	}
}

// ProposePost submits a post for collective consideration
func (c *ConsensusClient) ProposePost(ctx context.Context, req PostRequest, reasoning string) (*consensus.Decision, error) {
	// Validate the post content
	if err := validatePostRequest(req); err != nil {
		return nil, fmt.Errorf("invalid post request: %w", err)
	}
	
	// Create a proposal for collective review
	proposal := consensus.Proposal{
		ID:         generateProposalID(),
		Content:    req.Text,
		ProposedBy: getCurrentAgentID(ctx),
		Reasoning:  reasoning,
		ProposedAt: time.Now(),
		ExpiresAt:  time.Now().Add(24 * time.Hour), // 24 hour consensus period
	}
	
	// Submit to consensus system
	decision, err := c.consensus.ProposePost(ctx, proposal)
	if err != nil {
		return nil, fmt.Errorf("failed to propose post: %w", err)
	}
	
	// Store the post request for later publishing
	if err := c.storage.StorePostRequest(ctx, proposal.ID, req); err != nil {
		return nil, fmt.Errorf("failed to store post request: %w", err)
	}
	
	return decision, nil
}

// PublishWithConsensus publishes a post only after consensus is reached
func (c *ConsensusClient) PublishWithConsensus(ctx context.Context, proposalID string) (*PostResult, error) {
	// Check consensus status
	decision, err := c.consensus.GetDecision(ctx, proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get consensus decision: %w", err)
	}
	
	if decision.Status != consensus.StatusConsensus {
		return nil, fmt.Errorf("cannot publish: consensus not reached (status: %s)", decision.Status)
	}
	
	// Retrieve the stored post request
	req, err := c.storage.GetPostRequest(ctx, proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve post request: %w", err)
	}
	
	// Publish to Bluesky
	uri, cid, err := c.atpClient.CreatePost(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to publish to Bluesky: %w", err)
	}
	
	result := &PostResult{
		URI:         uri,
		CID:         cid,
		PostedAt:    time.Now(),
		ConsensusID: proposalID,
	}
	
	// Record the publication
	if err := c.storage.RecordPublication(ctx, proposalID, result); err != nil {
		return nil, fmt.Errorf("failed to record publication: %w", err)
	}
	
	return result, nil
}

// Helper functions

func validatePostRequest(req PostRequest) error {
	if req.Text == "" {
		return fmt.Errorf("post text cannot be empty")
	}
	
	// Bluesky has a 300 character limit
	if len(req.Text) > 300 {
		return fmt.Errorf("post text exceeds 300 character limit")
	}
	
	// Validate image count (max 4 images per post)
	if len(req.Images) > 4 {
		return fmt.Errorf("maximum 4 images allowed per post")
	}
	
	return nil
}

func generateProposalID() string {
	// Implementation would generate a unique ID
	return fmt.Sprintf("proposal-%d", time.Now().UnixNano())
}

func getCurrentAgentID(ctx context.Context) string {
	// Implementation would extract agent ID from context
	return "go-systems-developer"
}