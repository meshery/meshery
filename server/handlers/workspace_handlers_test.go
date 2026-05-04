package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

// workspaceSpyProvider embeds DefaultLocalProvider and records the orgID
// that GetWorkspaces / GetWorkspaceByID are invoked with. It allows the
// handler tests below to verify that the handler extracted the query
// parameter correctly.
type workspaceSpyProvider struct {
	*models.DefaultLocalProvider
	observedOrgID string
	called        bool
}

func newWorkspaceSpyProvider() *workspaceSpyProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &workspaceSpyProvider{DefaultLocalProvider: base}
}

func (m *workspaceSpyProvider) GetWorkspaces(_, _, _, _, _, _, orgID string) ([]byte, error) {
	m.called = true
	m.observedOrgID = orgID
	return []byte(`{"workspaces":[]}`), nil
}

func (m *workspaceSpyProvider) GetWorkspaceByID(_ *http.Request, _, orgID string) ([]byte, error) {
	m.called = true
	m.observedOrgID = orgID
	return []byte(`{}`), nil
}

// TestGetWorkspacesHandler_AcceptsOrgIdAndLegacyOrgID asserts that the
// canonical `orgId` query parameter is preferred, AND that the legacy
// `orgID` spelling is still accepted as a dual-accept fallback during
// the Phase 2 deprecation window (mesheryctl and other legacy clients
// still emit `orgID`). Missing parameter continues to return 400.
func TestGetWorkspacesHandler_AcceptsOrgIdAndLegacyOrgID(t *testing.T) {
	cases := []struct {
		name         string
		rawQuery     string
		wantStatus   int
		wantOrgID    string
		wantProvider bool
	}{
		{
			name:         "canonical orgId is accepted",
			rawQuery:     "orgId=abc",
			wantStatus:   http.StatusOK,
			wantOrgID:    "abc",
			wantProvider: true,
		},
		{
			name:         "legacy orgID is accepted via dual-accept fallback",
			rawQuery:     "orgID=abc",
			wantStatus:   http.StatusOK,
			wantOrgID:    "abc",
			wantProvider: true,
		},
		{
			name:         "canonical orgId wins when both are supplied",
			rawQuery:     "orgId=canonical&orgID=legacy",
			wantStatus:   http.StatusOK,
			wantOrgID:    "canonical",
			wantProvider: true,
		},
		{
			name:         "missing parameter returns 400",
			rawQuery:     "",
			wantStatus:   http.StatusBadRequest,
			wantProvider: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			provider := newWorkspaceSpyProvider()

			req := httptest.NewRequest(http.MethodGet, "/api/workspaces?"+tc.rawQuery, nil)
			req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
			rec := httptest.NewRecorder()

			h.GetWorkspacesHandler(rec, req, nil, nil, provider)

			if rec.Code != tc.wantStatus {
				t.Fatalf("expected status %d, got %d (body=%q)", tc.wantStatus, rec.Code, rec.Body.String())
			}

			if provider.called != tc.wantProvider {
				t.Fatalf("provider called=%v, want %v", provider.called, tc.wantProvider)
			}

			if tc.wantProvider && provider.observedOrgID != tc.wantOrgID {
				t.Fatalf("provider received orgID=%q, want %q", provider.observedOrgID, tc.wantOrgID)
			}

			if tc.wantStatus == http.StatusBadRequest {
				if !strings.Contains(rec.Body.String(), "orgId") {
					t.Errorf("expected 400 body to mention canonical orgId, got %q", rec.Body.String())
				}
			}
		})
	}
}

// TestGetWorkspaceByIdHandler_AcceptsOrgIdAndLegacyOrgID mirrors the coverage
// above for the single-workspace endpoint: canonical `orgId` preferred,
// legacy `orgID` dual-accepted during Phase 2.
func TestGetWorkspaceByIdHandler_AcceptsOrgIdAndLegacyOrgID(t *testing.T) {
	cases := []struct {
		name         string
		rawQuery     string
		wantStatus   int
		wantOrgID    string
		wantProvider bool
	}{
		{
			name:         "canonical orgId is accepted",
			rawQuery:     "orgId=abc",
			wantStatus:   http.StatusOK,
			wantOrgID:    "abc",
			wantProvider: true,
		},
		{
			name:         "legacy orgID is accepted via dual-accept fallback",
			rawQuery:     "orgID=abc",
			wantStatus:   http.StatusOK,
			wantOrgID:    "abc",
			wantProvider: true,
		},
		{
			name:         "canonical orgId wins when both are supplied",
			rawQuery:     "orgId=canonical&orgID=legacy",
			wantStatus:   http.StatusOK,
			wantOrgID:    "canonical",
			wantProvider: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			provider := newWorkspaceSpyProvider()

			req := httptest.NewRequest(http.MethodGet, "/api/workspaces/workspace-1?"+tc.rawQuery, nil)
			req = mux.SetURLVars(req, map[string]string{"id": "workspace-1"})
			rec := httptest.NewRecorder()

			h.GetWorkspaceByIdHandler(rec, req, nil, nil, provider)

			if rec.Code != tc.wantStatus {
				t.Fatalf("expected status %d, got %d (body=%q)", tc.wantStatus, rec.Code, rec.Body.String())
			}

			if provider.called != tc.wantProvider {
				t.Fatalf("provider called=%v, want %v", provider.called, tc.wantProvider)
			}

			if tc.wantProvider && provider.observedOrgID != tc.wantOrgID {
				t.Fatalf("provider received orgID=%q, want %q", provider.observedOrgID, tc.wantOrgID)
			}
		})
	}
}

// TestWorkspacePayloadWire_UnmarshalJSON exercises the dual-accept body
// contract for POST /api/workspaces. The schemas-generated struct tags
// OrganizationID as json:"organization_id", but canonical in-repo consumers
// now emit `organizationId`. Go's case-insensitive tag fallback does NOT
// match across an underscore boundary, so the wrapper's UnmarshalJSON must
// intercept both spellings during the Phase 2 deprecation window. Canonical
// must win when both are supplied. Table is kept narrow on purpose: the
// wrapper only routes OrganizationID, so we only assert that field here.
func TestWorkspacePayloadWire_UnmarshalJSON(t *testing.T) {
	const (
		canonicalUUID = "11111111-1111-1111-1111-111111111111"
		legacyUUID    = "22222222-2222-2222-2222-222222222222"
	)

	cases := []struct {
		name    string
		body    string
		wantOrg string
	}{
		{
			name:    "canonical organizationId only",
			body:    `{"name":"ws","organizationId":"` + canonicalUUID + `"}`,
			wantOrg: canonicalUUID,
		},
		{
			name:    "legacy organization_id only",
			body:    `{"name":"ws","organization_id":"` + legacyUUID + `"}`,
			wantOrg: legacyUUID,
		},
		{
			name:    "both supplied, canonical wins",
			body:    `{"name":"ws","organizationId":"` + canonicalUUID + `","organization_id":"` + legacyUUID + `"}`,
			wantOrg: canonicalUUID,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var wire workspacePayloadWire
			if err := json.Unmarshal([]byte(tc.body), &wire); err != nil {
				t.Fatalf("unexpected unmarshal error: %v", err)
			}
			if got := wire.OrganizationID.String(); got != tc.wantOrg {
				t.Fatalf("OrganizationID = %q, want %q", got, tc.wantOrg)
			}
			if wire.Name != "ws" {
				t.Fatalf("Name = %q, want %q", wire.Name, "ws")
			}
		})
	}
}

// TestWorkspaceUpdatePayloadWire_UnmarshalJSON mirrors the payload coverage
// for the PUT /api/workspaces/{id} path: the update wrapper must route both
// `organizationId` (canonical) and `organization_id` (legacy) onto the
// underlying schemas-generated OrganizationID field, with canonical taking
// precedence when both are present.
func TestWorkspaceUpdatePayloadWire_UnmarshalJSON(t *testing.T) {
	const (
		canonicalUUID = "33333333-3333-3333-3333-333333333333"
		legacyUUID    = "44444444-4444-4444-4444-444444444444"
	)

	cases := []struct {
		name    string
		body    string
		wantOrg string
	}{
		{
			name:    "canonical organizationId only",
			body:    `{"name":"ws","organizationId":"` + canonicalUUID + `"}`,
			wantOrg: canonicalUUID,
		},
		{
			name:    "legacy organization_id only",
			body:    `{"name":"ws","organization_id":"` + legacyUUID + `"}`,
			wantOrg: legacyUUID,
		},
		{
			name:    "both supplied, canonical wins",
			body:    `{"name":"ws","organizationId":"` + canonicalUUID + `","organization_id":"` + legacyUUID + `"}`,
			wantOrg: canonicalUUID,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var wire workspaceUpdatePayloadWire
			if err := json.Unmarshal([]byte(tc.body), &wire); err != nil {
				t.Fatalf("unexpected unmarshal error: %v", err)
			}
			if got := wire.OrganizationID.String(); got != tc.wantOrg {
				t.Fatalf("OrganizationID = %q, want %q", got, tc.wantOrg)
			}
		})
	}
}
