package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta3/workspace"
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

func (m *workspaceSpyProvider) SaveWorkspace(_ *http.Request, _ *workspace.WorkspacePayload, _ string, _ bool) ([]byte, error) {
	m.called = true
	return []byte(`{"id":"11111111-1111-1111-1111-111111111111"}`), nil
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

// workspaceFailingProvider embeds DefaultLocalProvider and fails every
// workspace call with a caller-supplied error.
type workspaceFailingProvider struct {
	*models.DefaultLocalProvider
	err error
}

func newWorkspaceFailingProvider(err error) *workspaceFailingProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &workspaceFailingProvider{DefaultLocalProvider: base, err: err}
}

func (m *workspaceFailingProvider) SaveWorkspace(_ *http.Request, _ *workspace.WorkspacePayload, _ string, _ bool) ([]byte, error) {
	return nil, m.err
}

func (m *workspaceFailingProvider) GetWorkspaces(_, _, _, _, _, _, _ string) ([]byte, error) {
	return nil, m.err
}

func (m *workspaceFailingProvider) DeleteWorkspace(_ *http.Request, _ string) ([]byte, error) {
	return nil, m.err
}

// TestSaveWorkspaceHandler_PropagatesProviderStatus mirrors the environment
// coverage: the workspace handlers carried the identical defect (ErrGetResult
// + hardcoded 404 on every provider failure), so they get the identical
// contract - real status out, workspace-specific code out.
func TestSaveWorkspaceHandler_PropagatesProviderStatus(t *testing.T) {
	cases := []struct {
		name        string
		providerErr error
		wantStatus  int
	}{
		{
			name:        "provider 403 surfaces as 403, not 404",
			providerErr: models.ErrPost(errors.New("failed to save the workspace"), "Workspace", http.StatusForbidden),
			wantStatus:  http.StatusForbidden,
		},
		{
			name:        "provider 409 surfaces as 409",
			providerErr: models.ErrPost(errors.New("workspace already exists"), "Workspace", http.StatusConflict),
			wantStatus:  http.StatusConflict,
		},
		{
			name:        "unreachable provider defaults to 502, never 404",
			providerErr: models.ErrUnreachableRemoteProvider(errors.New("dial tcp: connection refused")),
			wantStatus:  http.StatusBadGateway,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			provider := newWorkspaceFailingProvider(tc.providerErr)

			body := `{"name":"team-space","description":"","organizationId":"11111111-1111-1111-1111-111111111111"}`
			req := httptest.NewRequest(http.MethodPost, "/api/workspaces", strings.NewReader(body))
			rec := httptest.NewRecorder()

			h.SaveWorkspaceHandler(rec, req, nil, nil, provider)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if rec.Code == http.StatusCreated {
				t.Fatal("a failed create must never answer 201")
			}

			var decoded struct {
				Error                string   `json:"error"`
				Code                 string   `json:"code"`
				SuggestedRemediation []string `json:"suggestedRemediation"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
				t.Fatalf("error body did not parse as JSON: %v (body=%q)", err, rec.Body.String())
			}
			if decoded.Code != ErrSaveWorkspaceCode {
				t.Errorf("code = %q, want %q", decoded.Code, ErrSaveWorkspaceCode)
			}
			if decoded.Code == ErrGetResultCode {
				t.Errorf("code regressed to the performance-results code %s", ErrGetResultCode)
			}
			if strings.Contains(decoded.Error, "unable to get result") {
				t.Errorf("message regressed to the results wording: %q", decoded.Error)
			}
			if len(decoded.SuggestedRemediation) == 0 {
				t.Error("expected workspace-specific remediation on the wire")
			}
		})
	}
}

// TestGetWorkspacesHandler_PropagatesProviderStatus covers the read path.
func TestGetWorkspacesHandler_PropagatesProviderStatus(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newWorkspaceFailingProvider(
		models.ErrFetch(errors.New("forbidden"), "Workspaces", http.StatusForbidden),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/workspaces?orgId=11111111-1111-1111-1111-111111111111", nil)
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
	rec := httptest.NewRecorder()

	h.GetWorkspacesHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusForbidden, rec.Body.String())
	}

	var decoded struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("error body did not parse as JSON: %v", err)
	}
	if decoded.Code != ErrGetWorkspacesCode {
		t.Errorf("code = %q, want %q", decoded.Code, ErrGetWorkspacesCode)
	}
}

// TestSaveWorkspaceHandler_SetsContentTypeOnSuccess guards an adjacent bug
// found while fixing the above: the success path called w.WriteHeader before
// w.Header().Set, so Content-Type was silently dropped. RTK Query's default
// baseQuery dispatches on Content-Type, so the 201 body arrived unparsed.
func TestSaveWorkspaceHandler_SetsContentTypeOnSuccess(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newWorkspaceSpyProvider()

	body := `{"name":"team-space","organizationId":"11111111-1111-1111-1111-111111111111"}`
	req := httptest.NewRequest(http.MethodPost, "/api/workspaces", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.SaveWorkspaceHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
}
