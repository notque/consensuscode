// Package atproto provides a low-level HTTP client for the AT Protocol (Bluesky).
// It handles authentication, session management, and XRPC calls.
package atproto

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Client is an HTTP client for the AT Protocol XRPC endpoints.
type Client struct {
	httpClient *http.Client
	serviceURL string

	mu      sync.RWMutex
	session *Session
}

// Session holds AT Protocol authentication state.
type Session struct {
	AccessJWT  string `json:"accessJwt"`
	RefreshJWT string `json:"refreshJwt"`
	Handle     string `json:"handle"`
	DID        string `json:"did"`
}

// CreateSessionRequest is the request body for com.atproto.server.createSession.
type CreateSessionRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

// CreateRecordRequest is the request body for com.atproto.repo.createRecord.
type CreateRecordRequest struct {
	Repo       string      `json:"repo"`
	Collection string      `json:"collection"`
	Record     interface{} `json:"record"`
}

// CreateRecordResponse is the response from com.atproto.repo.createRecord.
type CreateRecordResponse struct {
	URI string `json:"uri"`
	CID string `json:"cid"`
}

// DeleteRecordRequest is the request body for com.atproto.repo.deleteRecord.
type DeleteRecordRequest struct {
	Repo       string `json:"repo"`
	Collection string `json:"collection"`
	RKey       string `json:"rkey"`
}

// FeedPost is the record structure for app.bsky.feed.post.
type FeedPost struct {
	Type      string     `json:"$type"`
	Text      string     `json:"text"`
	CreatedAt string     `json:"createdAt"`
	Langs     []string   `json:"langs,omitempty"`
	Reply     *ReplyRef  `json:"reply,omitempty"`
}

// ReplyRef references the parent and root posts for a reply.
type ReplyRef struct {
	Root   StrongRef `json:"root"`
	Parent StrongRef `json:"parent"`
}

// StrongRef is a reference to a specific record version.
type StrongRef struct {
	URI string `json:"uri"`
	CID string `json:"cid"`
}

// GetProfileResponse is the response from app.bsky.actor.getProfile.
type GetProfileResponse struct {
	DID         string `json:"did"`
	Handle      string `json:"handle"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	Banner      string `json:"banner"`
}

// FeedViewPost is a single post in a feed response.
type FeedViewPost struct {
	Post PostView `json:"post"`
}

// PostView is the view of a single post.
type PostView struct {
	URI       string          `json:"uri"`
	CID       string          `json:"cid"`
	Author    ActorView       `json:"author"`
	Record    json.RawMessage `json:"record"`
	IndexedAt string          `json:"indexedAt"`
}

// ActorView is the view of a user/actor.
type ActorView struct {
	DID         string `json:"did"`
	Handle      string `json:"handle"`
	DisplayName string `json:"displayName"`
}

// GetAuthorFeedResponse is the response from app.bsky.feed.getAuthorFeed.
type GetAuthorFeedResponse struct {
	Feed   []FeedViewPost `json:"feed"`
	Cursor string         `json:"cursor"`
}

// XRPCError is an error returned by the XRPC API.
type XRPCError struct {
	StatusCode int
	ErrorType  string `json:"error"`
	Message    string `json:"message"`
}

func (e *XRPCError) Error() string {
	return fmt.Sprintf("xrpc error %d: %s - %s", e.StatusCode, e.ErrorType, e.Message)
}

// NewClient creates a new AT Protocol client targeting the given service URL.
func NewClient(serviceURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		serviceURL: serviceURL,
	}
}

// NewClientWithHTTP creates a new AT Protocol client with a custom http.Client.
// This is useful for testing.
func NewClientWithHTTP(serviceURL string, httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		serviceURL: serviceURL,
	}
}

// Authenticate creates a session with the AT Protocol server.
func (c *Client) Authenticate(ctx context.Context, identifier, password string) error {
	reqBody := CreateSessionRequest{
		Identifier: identifier,
		Password:   password,
	}

	var session Session
	if err := c.doRequest(ctx, http.MethodPost, "com.atproto.server.createSession", reqBody, &session, false); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	c.mu.Lock()
	c.session = &session
	c.mu.Unlock()

	return nil
}

// CreatePost publishes a new post (app.bsky.feed.post record).
func (c *Client) CreatePost(ctx context.Context, text string, langs []string) (*CreateRecordResponse, error) {
	session, err := c.getSession()
	if err != nil {
		return nil, err
	}

	record := FeedPost{
		Type:      "app.bsky.feed.post",
		Text:      text,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Langs:     langs,
	}

	reqBody := CreateRecordRequest{
		Repo:       session.DID,
		Collection: "app.bsky.feed.post",
		Record:     record,
	}

	var resp CreateRecordResponse
	if err := c.doRequest(ctx, http.MethodPost, "com.atproto.repo.createRecord", reqBody, &resp, true); err != nil {
		return nil, fmt.Errorf("create post failed: %w", err)
	}

	return &resp, nil
}

// CreateReply publishes a reply to an existing post.
func (c *Client) CreateReply(ctx context.Context, text string, parentURI, parentCID, rootURI, rootCID string, langs []string) (*CreateRecordResponse, error) {
	session, err := c.getSession()
	if err != nil {
		return nil, err
	}

	record := FeedPost{
		Type:      "app.bsky.feed.post",
		Text:      text,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		Langs:     langs,
		Reply: &ReplyRef{
			Root:   StrongRef{URI: rootURI, CID: rootCID},
			Parent: StrongRef{URI: parentURI, CID: parentCID},
		},
	}

	reqBody := CreateRecordRequest{
		Repo:       session.DID,
		Collection: "app.bsky.feed.post",
		Record:     record,
	}

	var resp CreateRecordResponse
	if err := c.doRequest(ctx, http.MethodPost, "com.atproto.repo.createRecord", reqBody, &resp, true); err != nil {
		return nil, fmt.Errorf("create reply failed: %w", err)
	}

	return &resp, nil
}

// DeletePost deletes a post by its AT URI.
func (c *Client) DeletePost(ctx context.Context, uri string) error {
	session, err := c.getSession()
	if err != nil {
		return err
	}

	rkey, err := rkeyFromURI(uri)
	if err != nil {
		return err
	}

	reqBody := DeleteRecordRequest{
		Repo:       session.DID,
		Collection: "app.bsky.feed.post",
		RKey:       rkey,
	}

	if err := c.doRequest(ctx, http.MethodPost, "com.atproto.repo.deleteRecord", reqBody, nil, true); err != nil {
		return fmt.Errorf("delete post failed: %w", err)
	}

	return nil
}

// GetProfile retrieves a user's profile.
func (c *Client) GetProfile(ctx context.Context, actor string) (*GetProfileResponse, error) {
	url := fmt.Sprintf("%s/xrpc/app.bsky.actor.getProfile?actor=%s", c.serviceURL, actor)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	c.mu.RLock()
	if c.session != nil {
		req.Header.Set("Authorization", "Bearer "+c.session.AccessJWT)
	}
	c.mu.RUnlock()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseXRPCError(resp)
	}

	var profile GetProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &profile, nil
}

// GetAuthorFeed retrieves posts from a user's feed.
func (c *Client) GetAuthorFeed(ctx context.Context, actor string, limit int) (*GetAuthorFeedResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	url := fmt.Sprintf("%s/xrpc/app.bsky.feed.getAuthorFeed?actor=%s&limit=%d", c.serviceURL, actor, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	c.mu.RLock()
	if c.session != nil {
		req.Header.Set("Authorization", "Bearer "+c.session.AccessJWT)
	}
	c.mu.RUnlock()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, parseXRPCError(resp)
	}

	var feed GetAuthorFeedResponse
	if err := json.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &feed, nil
}

// IsAuthenticated reports whether the client holds a valid session.
func (c *Client) IsAuthenticated() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.session != nil && c.session.AccessJWT != ""
}

// GetDID returns the authenticated user's DID, or empty string if not authenticated.
func (c *Client) GetDID() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.session == nil {
		return ""
	}
	return c.session.DID
}

// GetHandle returns the authenticated user's handle, or empty string if not authenticated.
func (c *Client) GetHandle() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.session == nil {
		return ""
	}
	return c.session.Handle
}

// doRequest performs an XRPC request.
func (c *Client) doRequest(ctx context.Context, method, nsid string, body interface{}, result interface{}, auth bool) error {
	url := fmt.Sprintf("%s/xrpc/%s", c.serviceURL, nsid)

	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if auth {
		c.mu.RLock()
		session := c.session
		c.mu.RUnlock()

		if session == nil {
			return fmt.Errorf("not authenticated: call Authenticate first")
		}
		req.Header.Set("Authorization", "Bearer "+session.AccessJWT)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseXRPCError(resp)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// getSession returns the current session or an error if not authenticated.
func (c *Client) getSession() (*Session, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.session == nil {
		return nil, fmt.Errorf("not authenticated: call Authenticate first")
	}
	return c.session, nil
}

// parseXRPCError parses an error response from the XRPC API.
func parseXRPCError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	xrpcErr := &XRPCError{StatusCode: resp.StatusCode}
	if err := json.Unmarshal(body, xrpcErr); err != nil {
		xrpcErr.Message = string(body)
	}
	return xrpcErr
}

// rkeyFromURI extracts the record key from an AT URI.
// AT URIs have the form: at://did:plc:xxx/collection/rkey
func rkeyFromURI(uri string) (string, error) {
	// Minimal parsing: find the last path segment
	// at://did:plc:xxx/app.bsky.feed.post/3abcdef
	for i := len(uri) - 1; i >= 0; i-- {
		if uri[i] == '/' {
			return uri[i+1:], nil
		}
	}
	return "", fmt.Errorf("invalid AT URI: %s", uri)
}
