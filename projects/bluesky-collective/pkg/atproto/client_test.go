package atproto

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("https://bsky.social")
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
	if c.serviceURL != "https://bsky.social" {
		t.Errorf("serviceURL = %q, want %q", c.serviceURL, "https://bsky.social")
	}
	if c.IsAuthenticated() {
		t.Error("new client should not be authenticated")
	}
}

func TestAuthenticate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/xrpc/com.atproto.server.createSession" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CreateSessionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if req.Identifier != "test.bsky.social" || req.Password != "test-password" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "AuthenticationRequired",
				"message": "Invalid credentials",
			})
			return
		}

		json.NewEncoder(w).Encode(Session{
			AccessJWT:  "test-access-jwt",
			RefreshJWT: "test-refresh-jwt",
			Handle:     "test.bsky.social",
			DID:        "did:plc:testdid123",
		})
	}))
	defer server.Close()

	c := NewClient(server.URL)

	// Test successful authentication
	err := c.Authenticate(context.Background(), "test.bsky.social", "test-password")
	if err != nil {
		t.Fatalf("Authenticate failed: %v", err)
	}

	if !c.IsAuthenticated() {
		t.Error("client should be authenticated after successful login")
	}
	if got := c.GetDID(); got != "did:plc:testdid123" {
		t.Errorf("GetDID() = %q, want %q", got, "did:plc:testdid123")
	}
	if got := c.GetHandle(); got != "test.bsky.social" {
		t.Errorf("GetHandle() = %q, want %q", got, "test.bsky.social")
	}
}

func TestAuthenticateFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "AuthenticationRequired",
			"message": "Invalid identifier or password",
		})
	}))
	defer server.Close()

	c := NewClient(server.URL)
	err := c.Authenticate(context.Background(), "bad", "creds")
	if err == nil {
		t.Fatal("expected authentication to fail")
	}
	if c.IsAuthenticated() {
		t.Error("client should not be authenticated after failed login")
	}
}

func TestCreatePost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/xrpc/com.atproto.server.createSession":
			json.NewEncoder(w).Encode(Session{
				AccessJWT: "test-jwt",
				DID:       "did:plc:test",
				Handle:    "test.bsky.social",
			})
		case "/xrpc/com.atproto.repo.createRecord":
			auth := r.Header.Get("Authorization")
			if auth != "Bearer test-jwt" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			var req CreateRecordRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}

			if req.Collection != "app.bsky.feed.post" {
				t.Errorf("collection = %q, want %q", req.Collection, "app.bsky.feed.post")
			}

			json.NewEncoder(w).Encode(CreateRecordResponse{
				URI: "at://did:plc:test/app.bsky.feed.post/abc123",
				CID: "bafyreiabc123",
			})
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := NewClient(server.URL)
	ctx := context.Background()

	// Must authenticate first
	if err := c.Authenticate(ctx, "test", "pass"); err != nil {
		t.Fatalf("auth: %v", err)
	}

	resp, err := c.CreatePost(ctx, "Hello from the collective!", nil)
	if err != nil {
		t.Fatalf("CreatePost: %v", err)
	}

	if resp.URI != "at://did:plc:test/app.bsky.feed.post/abc123" {
		t.Errorf("URI = %q, want at://did:plc:test/app.bsky.feed.post/abc123", resp.URI)
	}
	if resp.CID != "bafyreiabc123" {
		t.Errorf("CID = %q, want bafyreiabc123", resp.CID)
	}
}

func TestCreatePostWithoutAuth(t *testing.T) {
	c := NewClient("https://bsky.social")
	_, err := c.CreatePost(context.Background(), "test", nil)
	if err == nil {
		t.Fatal("expected error when creating post without authentication")
	}
}

func TestDeletePost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/xrpc/com.atproto.server.createSession":
			json.NewEncoder(w).Encode(Session{
				AccessJWT: "test-jwt",
				DID:       "did:plc:test",
			})
		case "/xrpc/com.atproto.repo.deleteRecord":
			var req DeleteRecordRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "bad request", http.StatusBadRequest)
				return
			}
			if req.RKey != "abc123" {
				t.Errorf("rkey = %q, want %q", req.RKey, "abc123")
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{})
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := NewClient(server.URL)
	ctx := context.Background()
	if err := c.Authenticate(ctx, "test", "pass"); err != nil {
		t.Fatalf("auth: %v", err)
	}

	err := c.DeletePost(ctx, "at://did:plc:test/app.bsky.feed.post/abc123")
	if err != nil {
		t.Fatalf("DeletePost: %v", err)
	}
}

func TestGetProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/xrpc/app.bsky.actor.getProfile" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		actor := r.URL.Query().Get("actor")
		json.NewEncoder(w).Encode(GetProfileResponse{
			DID:         "did:plc:test",
			Handle:      actor,
			DisplayName: "Test User",
			Description: "A test profile",
		})
	}))
	defer server.Close()

	c := NewClient(server.URL)
	profile, err := c.GetProfile(context.Background(), "test.bsky.social")
	if err != nil {
		t.Fatalf("GetProfile: %v", err)
	}
	if profile.Handle != "test.bsky.social" {
		t.Errorf("Handle = %q, want %q", profile.Handle, "test.bsky.social")
	}
	if profile.DisplayName != "Test User" {
		t.Errorf("DisplayName = %q, want %q", profile.DisplayName, "Test User")
	}
}

func TestGetAuthorFeed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/xrpc/app.bsky.feed.getAuthorFeed" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		record, _ := json.Marshal(map[string]string{
			"text":      "Hello world",
			"createdAt": "2025-01-01T00:00:00Z",
		})

		json.NewEncoder(w).Encode(GetAuthorFeedResponse{
			Feed: []FeedViewPost{
				{
					Post: PostView{
						URI:    "at://did:plc:test/app.bsky.feed.post/1",
						CID:    "bafyrei1",
						Author: ActorView{DID: "did:plc:test", Handle: "test.bsky.social"},
						Record: record,
					},
				},
			},
		})
	}))
	defer server.Close()

	c := NewClient(server.URL)
	feed, err := c.GetAuthorFeed(context.Background(), "test.bsky.social", 10)
	if err != nil {
		t.Fatalf("GetAuthorFeed: %v", err)
	}
	if len(feed.Feed) != 1 {
		t.Fatalf("feed length = %d, want 1", len(feed.Feed))
	}
	if feed.Feed[0].Post.URI != "at://did:plc:test/app.bsky.feed.post/1" {
		t.Errorf("post URI = %q", feed.Feed[0].Post.URI)
	}
}

func TestRkeyFromURI(t *testing.T) {
	tests := []struct {
		uri  string
		want string
	}{
		{"at://did:plc:test/app.bsky.feed.post/abc123", "abc123"},
		{"at://did:plc:test/app.bsky.feed.post/3k5abc", "3k5abc"},
	}

	for _, tt := range tests {
		got, err := rkeyFromURI(tt.uri)
		if err != nil {
			t.Errorf("rkeyFromURI(%q) error: %v", tt.uri, err)
			continue
		}
		if got != tt.want {
			t.Errorf("rkeyFromURI(%q) = %q, want %q", tt.uri, got, tt.want)
		}
	}
}

func TestXRPCError(t *testing.T) {
	err := &XRPCError{
		StatusCode: 401,
		ErrorType:  "AuthenticationRequired",
		Message:    "Invalid credentials",
	}

	s := err.Error()
	if s == "" {
		t.Error("XRPCError.Error() returned empty string")
	}
}
