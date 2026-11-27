package models

import (
	"context"
	stderrors "errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/httputil"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func TestRemoteProviderMarkAllEventsReadForwardsRequest(t *testing.T) {
	const testJWT = "test-user-jwt"
	var gotMethod, gotPath, gotAuthorization, gotAPIKey string
	var gotBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAuthorization = r.Header.Get("Authorization")
		gotAPIKey = r.Header.Get("X-API-Key")
		gotBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"All user-owned events marked as read"}`))
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{{Feature: PersistEvents, Endpoint: "/events"}}
	if err := provider.MarkAllEventsRead(context.Background(), testJWT, uuid.Nil); err != nil {
		t.Fatalf("MarkAllEventsRead returned error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/api/events/mark-all-read" {
		t.Errorf("path = %q, want /api/events/mark-all-read", gotPath)
	}
	if len(gotBody) != 0 {
		t.Errorf("body = %q, want empty body", gotBody)
	}
	if gotAuthorization != "bearer "+testJWT {
		t.Errorf("Authorization = %q, want bearer token", gotAuthorization)
	}
	if gotAPIKey != "" {
		t.Errorf("X-API-Key = %q, want empty for an authenticated user", gotAPIKey)
	}
}

func TestRemoteProviderMarkAllEventsReadPreservesFailureStatus(t *testing.T) {
	for _, status := range []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusTooManyRequests, http.StatusInternalServerError} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(status)
			}))
			defer server.Close()

			provider := newTestRemoteProvider(t, server.URL)
			provider.Capabilities = Capabilities{{Feature: PersistEvents, Endpoint: "/events"}}
			err := provider.MarkAllEventsRead(context.Background(), "test-user-jwt", uuid.Nil)
			if err == nil {
				t.Fatalf("MarkAllEventsRead returned nil, want HTTP %d error", status)
			}
			if got := httputil.StatusForProviderError(err, 0); got != status {
				t.Errorf("provider status = %d, want %d", got, status)
			}
		})
	}
}

func TestRemoteProviderMarkAllEventsReadHonorsCancellation(t *testing.T) {
	requestStarted := make(chan struct{})
	cleanup := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		close(requestStarted)
		select {
		case <-r.Context().Done():
		case <-cleanup:
		}
	}))
	t.Cleanup(func() {
		close(cleanup)
		server.Close()
	})

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{{Feature: PersistEvents, Endpoint: "/events"}}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- provider.MarkAllEventsRead(ctx, "test-user-jwt", uuid.Nil)
	}()

	select {
	case <-requestStarted:
		cancel()
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for remote request")
	}

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("MarkAllEventsRead returned nil after cancellation")
		}
	case <-time.After(time.Second):
		t.Fatal("MarkAllEventsRead did not return after context cancellation")
	}
}

func TestRemoteProviderMarkAllEventsReadCancelsTokenRefresh(t *testing.T) {
	refreshStarted := make(chan struct{})
	cleanup := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/events/mark-all-read":
			w.WriteHeader(http.StatusUnauthorized)
		case "/refresh":
			close(refreshStarted)
			select {
			case <-r.Context().Done():
			case <-cleanup:
			}
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(func() {
		close(cleanup)
		server.Close()
	})

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{{Feature: PersistEvents, Endpoint: "/events"}}
	token := encodeTestToken(t, oauth2.Token{AccessToken: "expired", RefreshToken: "refreshable"})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- provider.MarkAllEventsRead(ctx, token, uuid.Nil)
	}()

	select {
	case <-refreshStarted:
		cancel()
	case <-time.After(time.Second):
		cancel()
		t.Fatal("timed out waiting for token refresh")
	}

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("MarkAllEventsRead returned nil after cancellation")
		}
	case <-time.After(time.Second):
		t.Fatal("MarkAllEventsRead did not return after cancelling token refresh")
	}
}

func TestRemoteProviderMarkAllEventsReadCancelsWhileRefreshInFlight(t *testing.T) {
	markAllStarted := make(chan struct{}, 2)
	refreshStarted := make(chan struct{})
	cleanup := make(chan struct{})
	var refreshCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/events/mark-all-read":
			markAllStarted <- struct{}{}
			w.WriteHeader(http.StatusUnauthorized)
		case "/refresh":
			if refreshCalls.Add(1) == 1 {
				close(refreshStarted)
			}
			select {
			case <-r.Context().Done():
			case <-cleanup:
			}
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(func() {
		close(cleanup)
		server.Close()
	})

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{{Feature: PersistEvents, Endpoint: "/events"}}
	token := encodeTestToken(t, oauth2.Token{AccessToken: "expired", RefreshToken: "refreshable"})

	firstCtx, cancelFirst := context.WithCancel(context.Background())
	defer cancelFirst()
	firstDone := make(chan error, 1)
	go func() {
		firstDone <- provider.MarkAllEventsRead(firstCtx, token, uuid.Nil)
	}()

	select {
	case <-refreshStarted:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for first token refresh")
	}

	secondCtx, cancelSecond := context.WithCancel(context.Background())
	secondDone := make(chan error, 1)
	go func() {
		secondDone <- provider.MarkAllEventsRead(secondCtx, token, uuid.Nil)
	}()
	select {
	case <-markAllStarted:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for first mark-all request")
	}
	select {
	case <-markAllStarted:
		cancelSecond()
	case <-time.After(time.Second):
		cancelSecond()
		t.Fatal("timed out waiting for second mark-all request")
	}

	select {
	case err := <-secondDone:
		if err == nil {
			t.Fatal("second MarkAllEventsRead returned nil after cancellation")
		}
	case <-time.After(time.Second):
		t.Fatal("second MarkAllEventsRead remained blocked behind the in-flight refresh")
	}

	cancelFirst()
	select {
	case err := <-firstDone:
		if err == nil {
			t.Fatal("first MarkAllEventsRead returned nil after cancellation")
		}
	case <-time.After(time.Second):
		t.Fatal("first MarkAllEventsRead did not return after cancellation")
	}
	if got := refreshCalls.Load(); got != 1 {
		t.Fatalf("refresh calls = %d, want 1", got)
	}
}

func TestRemoteProviderMarkAllEventsReadRetriesAfterRefreshLeaderCancellation(t *testing.T) {
	markAllStarted := make(chan struct{}, 2)
	firstRefreshStarted := make(chan struct{})
	cleanup := make(chan struct{})
	var refreshCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/events/mark-all-read":
			if r.Header.Get("Authorization") == "bearer fresh-token" {
				w.WriteHeader(http.StatusOK)
				return
			}
			markAllStarted <- struct{}{}
			w.WriteHeader(http.StatusUnauthorized)
		case "/refresh":
			if refreshCalls.Add(1) == 1 {
				close(firstRefreshStarted)
				select {
				case <-r.Context().Done():
				case <-cleanup:
				}
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"token":"fresh-token"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(func() {
		close(cleanup)
		server.Close()
	})

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{{Feature: PersistEvents, Endpoint: "/events"}}
	token := encodeTestToken(t, oauth2.Token{AccessToken: "expired", RefreshToken: "refreshable"})

	leaderCtx, cancelLeader := context.WithCancel(context.Background())
	defer cancelLeader()
	leaderDone := make(chan error, 1)
	go func() {
		leaderDone <- provider.MarkAllEventsRead(leaderCtx, token, uuid.Nil)
	}()

	select {
	case <-firstRefreshStarted:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for leader token refresh")
	}

	waiterDone := make(chan error, 1)
	go func() {
		waiterDone <- provider.MarkAllEventsRead(context.Background(), token, uuid.Nil)
	}()
	select {
	case <-markAllStarted:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for leader mark-all request")
	}
	select {
	case <-markAllStarted:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for waiter mark-all request")
	}

	cancelLeader()
	select {
	case err := <-leaderDone:
		if err == nil {
			t.Fatal("leader MarkAllEventsRead returned nil after cancellation")
		}
	case <-time.After(time.Second):
		t.Fatal("leader MarkAllEventsRead did not return after cancellation")
	}

	select {
	case err := <-waiterDone:
		if err != nil {
			t.Fatalf("waiter MarkAllEventsRead returned error: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("waiter MarkAllEventsRead did not retry after leader cancellation")
	}
	if got := refreshCalls.Load(); got != 2 {
		t.Fatalf("refresh calls = %d, want 2", got)
	}
}

func TestRemoteProviderMarkAllEventsReadPreservesRefreshFailureStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/events/mark-all-read", "/refresh":
			w.WriteHeader(http.StatusUnauthorized)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{{Feature: PersistEvents, Endpoint: "/events"}}
	token := encodeTestToken(t, oauth2.Token{AccessToken: "expired", RefreshToken: "refreshable"})
	err := provider.MarkAllEventsRead(context.Background(), token, uuid.Nil)
	if err == nil {
		t.Fatal("MarkAllEventsRead returned nil, want refresh failure")
	}
	if got := httputil.StatusForProviderError(err, 0); got != http.StatusUnauthorized {
		t.Fatalf("provider status = %d, want %d", got, http.StatusUnauthorized)
	}
}

// TestRemoteProvider_OutboundOrgIdIsCanonical asserts that RemoteProvider
// emits the canonical camelCase `orgId` query parameter (not the legacy
// PascalCase `orgID`) on every outbound URL construction site identified
// in the identifier-naming migration plan (Agent 2.B).
//
// Sites covered:
//   - GetEnvironments
//   - GetEnvironmentByID
//   - GetWorkspaces
//   - GetWorkspaceByID
//   - GetCatalogMesheryPatterns (q.Add pluralized form)
func TestRemoteProvider_OutboundOrgIdIsCanonical(t *testing.T) {
	const (
		orgID       = "org-123"
		workspaceID = "ws-456"
		envID       = "env-789"
	)

	type call struct {
		path      string
		rawQuery  string
		orgIdVals []string
	}

	cases := []struct {
		name      string
		feature   Feature
		endpoint  string
		invoke    func(t *testing.T, l *RemoteProvider)
		wantOrgID string
		wantPath  string
	}{
		{
			name:     "GetEnvironments uses orgId",
			feature:  PersistEnvironments,
			endpoint: "/environments",
			invoke: func(t *testing.T, l *RemoteProvider) {
				if _, err := l.GetEnvironments("test-token", "0", "10", "", "", "", orgID); err != nil {
					t.Fatalf("GetEnvironments: %v", err)
				}
			},
			wantOrgID: orgID,
			wantPath:  "/environments",
		},
		{
			name:     "GetEnvironmentByID uses orgId",
			feature:  PersistEnvironments,
			endpoint: "/environments",
			invoke: func(t *testing.T, l *RemoteProvider) {
				req, _ := http.NewRequest(http.MethodGet, "/ignored", nil)
				req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: "test-token"})
				if _, err := l.GetEnvironmentByID(req, envID, orgID); err != nil {
					t.Fatalf("GetEnvironmentByID: %v", err)
				}
			},
			wantOrgID: orgID,
			wantPath:  "/environments/" + envID,
		},
		{
			name:     "GetWorkspaces uses orgId",
			feature:  PersistWorkspaces,
			endpoint: "/workspaces",
			invoke: func(t *testing.T, l *RemoteProvider) {
				if _, err := l.GetWorkspaces("test-token", "0", "10", "", "", "", orgID); err != nil {
					t.Fatalf("GetWorkspaces: %v", err)
				}
			},
			wantOrgID: orgID,
			wantPath:  "/workspaces",
		},
		{
			name:     "GetWorkspaceByID uses orgId",
			feature:  PersistWorkspaces,
			endpoint: "/workspaces",
			invoke: func(t *testing.T, l *RemoteProvider) {
				req, _ := http.NewRequest(http.MethodGet, "/ignored", nil)
				req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: "test-token"})
				if _, err := l.GetWorkspaceByID(req, workspaceID, orgID); err != nil {
					t.Fatalf("GetWorkspaceByID: %v", err)
				}
			},
			wantOrgID: orgID,
			wantPath:  "/workspaces/" + workspaceID,
		},
		{
			name:     "GetCatalogMesheryPatterns uses orgId (q.Add plural form)",
			feature:  MesheryPatternsCatalog,
			endpoint: "/patterns/catalog",
			invoke: func(t *testing.T, l *RemoteProvider) {
				if _, err := l.GetCatalogMesheryPatterns(
					"test-token",
					"0", "10", "", "", "",
					nil, nil, nil, nil,
					[]string{orgID}, nil, nil,
				); err != nil {
					t.Fatalf("GetCatalogMesheryPatterns: %v", err)
				}
			},
			wantOrgID: orgID,
			wantPath:  "/patterns/catalog",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var observed atomic.Value // of call
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				observed.Store(call{
					path:      r.URL.Path,
					rawQuery:  r.URL.RawQuery,
					orgIdVals: r.URL.Query()["orgId"],
				})
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))
			}))
			defer server.Close()

			provider := newTestRemoteProvider(t, server.URL)
			provider.Capabilities = Capabilities{
				{Feature: tc.feature, Endpoint: tc.endpoint},
			}

			tc.invoke(t, provider)

			got, ok := observed.Load().(call)
			if !ok {
				t.Fatal("no request was observed by the test server")
			}

			// Canonical form must be present.
			if len(got.orgIdVals) == 0 || got.orgIdVals[0] != tc.wantOrgID {
				t.Errorf("expected orgId=%q in query, got %q (raw=%q)", tc.wantOrgID, got.orgIdVals, got.rawQuery)
			}

			// Legacy form must NOT be present.
			if strings.Contains(got.rawQuery, "orgID=") {
				t.Errorf("found legacy orgID=... in outbound query, expected canonical orgId=. raw=%q", got.rawQuery)
			}

			if tc.wantPath != "" && got.path != tc.wantPath {
				t.Errorf("expected path %q, got %q", tc.wantPath, got.path)
			}
		})
	}
}

// TestRemoteProvider_GetUserByID_SystemInstanceReturnsSentinel pins the
// contract with GetUserByIDHandler: when the caller asks for this Meshery
// instance's own UUID (because a design's user_id is the instance UUID and
// not a real user), the remote provider must signal with
// ErrUserIsSystemInstance so the handler can respond 204 No Content instead
// of conflating it with "user missing" (404) or emitting a plain-text body
// that later crashes the client's JSON parser. This is the signal the Open
// Recents 404/SyntaxError loop traced back to.
func TestRemoteProvider_GetUserByID_SystemInstanceReturnsSentinel(t *testing.T) {
	const instanceID = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	viper.Set("INSTANCE_ID", instanceID)
	t.Cleanup(func() { viper.Set("INSTANCE_ID", "") })

	rp := &RemoteProvider{}
	body, err := rp.GetUserByID(nil, instanceID)

	if body != nil {
		t.Errorf("expected nil body for system instance ID, got %q", string(body))
	}
	if err == nil {
		t.Fatal("expected ErrUserIsSystemInstance, got nil — silent (nil, nil) is the old bug")
	}
	if !stderrors.Is(err, ErrUserIsSystemInstance) {
		t.Errorf("expected err to match ErrUserIsSystemInstance via errors.Is, got %v", err)
	}
}
