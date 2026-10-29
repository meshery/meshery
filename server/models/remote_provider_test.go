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
	for _, status := range []int{http.StatusUnauthorized, http.StatusTooManyRequests, http.StatusInternalServerError} {
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
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		close(requestStarted)
		<-r.Context().Done()
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{{Feature: PersistEvents, Endpoint: "/events"}}
	ctx, cancel := context.WithCancel(context.Background())
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
