package models

import (
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/viper"
)

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
		name       string
		feature    Feature
		endpoint   string
		invoke     func(t *testing.T, l *RemoteProvider)
		wantOrgID  string
		wantPath   string
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

// TestRemoteProvider_GetUserDetails_CachingAndDeduplication verifies that
// GetUserDetails implements in-memory TTL caching and singleflight request
// deduplication (#22796):
//  1. 10 concurrent requests for the same token trigger exactly 1 outbound call.
//  2. Subsequent calls within TTL are served from memory (0 outbound calls).
//  3. After TTL expires, the next request refreshes from remote provider.
//  4. Logging out evicts the cached entry so the next request refreshes.
func TestRemoteProvider_GetUserDetails_CachingAndDeduplication(t *testing.T) {
	var outboundCalls int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/logout" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message": "logged out"}`))
			return
		}

		if r.URL.Path == "/api/user/profile" {
			atomic.AddInt32(&outboundCalls, 1)
			time.Sleep(30 * time.Millisecond) // Simulate WAN latency

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id": "11111111-2222-3333-4444-555555555555",
				"firstName": "Test",
				"lastName": "User",
				"email": "test@meshery.io"
			}`))
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{
		{Feature: UsersProfile, Endpoint: "/api/user/profile"},
	}
	provider.UserCacheTTL = 60 * time.Millisecond

	const testToken = "test-auth-jwt-token"
	req, _ := http.NewRequest(http.MethodGet, "/api/dashboard", nil)
	req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: testToken})

	// Phase 1: 10 concurrent requests with the same token
	const concurrentCount = 10
	var wg sync.WaitGroup
	wg.Add(concurrentCount)

	for i := 0; i < concurrentCount; i++ {
		go func() {
			defer wg.Done()
			user, err := provider.GetUserDetails(req)
			if err != nil {
				t.Errorf("GetUserDetails failed: %v", err)
				return
			}
			if user == nil || user.FirstName != "Test" {
				t.Errorf("expected user with FirstName 'Test', got %+v", user)
			}
		}()
	}
	wg.Wait()

	// Assert: Singleflight deduplicated all concurrent requests to exactly 1 outbound call
	if calls := atomic.LoadInt32(&outboundCalls); calls != 1 {
		t.Fatalf("expected exactly 1 outbound HTTP call for 10 concurrent requests, got %d", calls)
	}

	// Phase 2: Sequential cached call (served from memory)
	user, err := provider.GetUserDetails(req)
	if err != nil || user == nil {
		t.Fatalf("cached GetUserDetails failed: %v", err)
	}
	if calls := atomic.LoadInt32(&outboundCalls); calls != 1 {
		t.Fatalf("expected outbound calls to remain 1 on cache hit, got %d", calls)
	}

	// Phase 3: TTL expiry
	time.Sleep(80 * time.Millisecond) // Wait for UserCacheTTL (60ms) to elapse

	userExpired, err := provider.GetUserDetails(req)
	if err != nil || userExpired == nil {
		t.Fatalf("GetUserDetails after TTL failed: %v", err)
	}
	if calls := atomic.LoadInt32(&outboundCalls); calls != 2 {
		t.Fatalf("expected outbound calls to be 2 after TTL expiry, got %d", calls)
	}

	// Phase 4: Logout invalidation
	rec := httptest.NewRecorder()
	logoutReq, _ := http.NewRequest(http.MethodPost, "/logout", nil)
	logoutReq.AddCookie(&http.Cookie{Name: TokenCookieName, Value: testToken})
	logoutReq.AddCookie(&http.Cookie{Name: "session_cookie", Value: "test-session"})

	if err := provider.Logout(rec, logoutReq); err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	// Next call must refresh from remote provider because cache entry was evicted
	userAfterLogout, err := provider.GetUserDetails(req)
	if err != nil || userAfterLogout == nil {
		t.Fatalf("GetUserDetails after logout failed: %v", err)
	}
	if calls := atomic.LoadInt32(&outboundCalls); calls != 3 {
		t.Fatalf("expected outbound calls to be 3 after logout eviction, got %d", calls)
	}
}
