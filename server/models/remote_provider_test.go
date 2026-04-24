package models

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
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
