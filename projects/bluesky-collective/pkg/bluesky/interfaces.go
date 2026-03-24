package bluesky

import (
	"context"
	"time"

	"github.com/consensuscode/bluesky-collective/pkg/consensus"
	"github.com/consensuscode/bluesky-collective/pkg/storage"
)

// Poster defines the AT Protocol posting operations needed by the bluesky package.
type Poster interface {
	// Authenticate creates a session with the AT Protocol server.
	Authenticate(ctx context.Context, identifier, password string) error

	// CreatePost publishes a post.
	CreatePost(ctx context.Context, text string, langs []string) (uri string, cid string, err error)

	// DeletePost removes a post by AT URI.
	DeletePost(ctx context.Context, uri string) error

	// IsAuthenticated reports whether the client holds a valid session.
	IsAuthenticated() bool

	// GetDID returns the authenticated user's DID.
	GetDID() string
}

// Store defines the storage operations needed by the bluesky package.
type Store interface {
	StorePostRequest(ctx context.Context, proposalID string, req storage.PostRequest) error
	GetPostRequest(ctx context.Context, proposalID string) (*storage.PostRequest, error)
	RecordPublication(ctx context.Context, proposalID string, result storage.PostResult) error
	GetPublicationHistory(ctx context.Context, limit int) ([]storage.PostResult, error)
}

// Consensus defines the consensus operations needed by the bluesky package.
type Consensus interface {
	ProposePost(ctx context.Context, proposal consensus.Proposal) (*consensus.Decision, error)
	GetDecision(ctx context.Context, proposalID string) (*consensus.Decision, error)
	RecordVote(ctx context.Context, proposalID string, vote consensus.Vote) error
	CheckConsensus(ctx context.Context, proposalID string) (bool, error)
	ListPendingProposals(ctx context.Context) ([]consensus.Proposal, error)
	GetProposal(ctx context.Context, proposalID string) (*consensus.Proposal, error)
}

// Post represents a Bluesky post.
type Post struct {
	URI       string    `json:"uri"`
	CID       string    `json:"cid"`
	Author    Author    `json:"author"`
	Record    Record    `json:"record"`
	CreatedAt time.Time `json:"created_at"`
}

// Author represents a post author.
type Author struct {
	DID         string `json:"did"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
}

// Record represents the post content record.
type Record struct {
	Text      string   `json:"text"`
	CreatedAt string   `json:"created_at"`
	Languages []string `json:"langs,omitempty"`
}
