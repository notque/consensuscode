package bluesky

import (
	"context"

	"github.com/consensuscode/bluesky-collective/pkg/atproto"
)

// ATPAdapter wraps an atproto.Client to satisfy the Poster interface.
type ATPAdapter struct {
	client *atproto.Client
}

// NewATPAdapter creates a Poster backed by an AT Protocol client.
func NewATPAdapter(client *atproto.Client) *ATPAdapter {
	return &ATPAdapter{client: client}
}

// Authenticate creates a session with the AT Protocol server.
func (a *ATPAdapter) Authenticate(ctx context.Context, identifier, password string) error {
	return a.client.Authenticate(ctx, identifier, password)
}

// CreatePost publishes a post and returns its URI and CID.
func (a *ATPAdapter) CreatePost(ctx context.Context, text string, langs []string) (string, string, error) {
	resp, err := a.client.CreatePost(ctx, text, langs)
	if err != nil {
		return "", "", err
	}
	return resp.URI, resp.CID, nil
}

// DeletePost removes a post by AT URI.
func (a *ATPAdapter) DeletePost(ctx context.Context, uri string) error {
	return a.client.DeletePost(ctx, uri)
}

// IsAuthenticated reports whether the client holds a valid session.
func (a *ATPAdapter) IsAuthenticated() bool {
	return a.client.IsAuthenticated()
}

// GetDID returns the authenticated user's DID.
func (a *ATPAdapter) GetDID() string {
	return a.client.GetDID()
}
