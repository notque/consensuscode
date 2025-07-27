package bluesky

import (
	"context"
	"time"
)

// ATPClient defines the low-level AT Protocol operations
type ATPClient interface {
	// CreatePost creates a post on Bluesky
	CreatePost(ctx context.Context, req PostRequest) (uri string, cid string, error error)
	
	// DeletePost removes a post from Bluesky
	DeletePost(ctx context.Context, uri string) error
	
	// GetPost retrieves a specific post
	GetPost(ctx context.Context, uri string) (*Post, error)
	
	// Authenticate performs authentication with Bluesky
	Authenticate(ctx context.Context, identifier, password string) error
}

// Storage defines the interface for storing consensus and post data
type Storage interface {
	// StorePostRequest saves a post request for later publishing
	StorePostRequest(ctx context.Context, proposalID string, req PostRequest) error
	
	// GetPostRequest retrieves a stored post request
	GetPostRequest(ctx context.Context, proposalID string) (*PostRequest, error)
	
	// RecordPublication records that a post has been published
	RecordPublication(ctx context.Context, proposalID string, result *PostResult) error
	
	// GetPublicationHistory retrieves publication history
	GetPublicationHistory(ctx context.Context, limit int) ([]PostResult, error)
}

// Post represents a Bluesky post
type Post struct {
	URI       string                 `json:"uri"`
	CID       string                 `json:"cid"`
	Author    Author                 `json:"author"`
	Record    Record                 `json:"record"`
	CreatedAt time.Time              `json:"created_at"`
}

// Author represents a post author
type Author struct {
	DID         string `json:"did"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
}

// Record represents the post content record
type Record struct {
	Text      string   `json:"text"`
	CreatedAt string   `json:"created_at"`
	Languages []string `json:"langs,omitempty"`
}