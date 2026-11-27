package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta3/workspace"
	"gorm.io/gorm"
)

// assertMissingOrgIDErrorEnvelope validates that a 400 Bad Request error body
// matches the MeshKit ErrWorkspaceMissingInput envelope and mentions the canonical orgId parameter.
func assertMissingOrgIDErrorEnvelope(t *testing.T, body []byte) {
	t.Helper()
	decoded := decodeErrorBody(t, body)
	if decoded.Code != models.ErrWorkspaceMissingInputCode {
		t.Errorf("expected error code %q, got %q", models.ErrWorkspaceMissingInputCode, decoded.Code)
	}
	if !strings.Contains(strings.ToLower(string(body)), "orgid") {
		t.Errorf("expected 400 body to mention canonical orgId, got %q", string(body))
	}
}

// assertPaginationErrorEnvelope validates that a pagination rejection surfaces
// the expected MeshKit error code inside a structured error envelope.
func assertPaginationErrorEnvelope(t *testing.T, body []byte, wantCode string) {
	t.Helper()
	decoded := decodeErrorBody(t, body)
	if decoded.Code != wantCode {
		t.Errorf("code = %q, want %q (body=%q)", decoded.Code, wantCode, string(body))
	}
}

// testWorkspaceID is the canonical test workspace UUID used across spy fixtures and handler tests.
const testWorkspaceID = "11111111-1111-1111-1111-111111111111"

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
	targetID := uuid.FromStringOrNil(testWorkspaceID)
	return json.Marshal(&workspace.WorkspacePage{
		Page:       0,
		PageSize:   1,
		TotalCount: 1,
		Workspaces: []workspace.AvailableWorkspace{
			{
				ID:   targetID,
				Name: "test-workspace",
			},
		},
	})
}

func (m *workspaceSpyProvider) GetWorkspaceByID(_ *http.Request, _, orgID string) ([]byte, error) {
	m.called = true
	m.observedOrgID = orgID
	targetID := uuid.FromStringOrNil(testWorkspaceID)
	return json.Marshal(&workspace.Workspace{
		ID:   targetID,
		Name: "test-workspace",
	})
}

func (m *workspaceSpyProvider) SaveWorkspace(_ *http.Request, _ *workspace.WorkspacePayload, _ string, _ bool) ([]byte, error) {
	m.called = true
	return []byte(`{"id":"` + testWorkspaceID + `"}`), nil
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

			if tc.wantStatus == http.StatusOK {
				if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
					t.Errorf("expected Content-Type application/json, got %q", ct)
				}
				var page workspace.WorkspacePage
				if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
					t.Fatalf("failed to unmarshal workspaces payload: %v (body=%q)", err, rec.Body.String())
				}
				if len(page.Workspaces) != 1 || page.Workspaces[0].Name != "test-workspace" {
					t.Errorf("unexpected workspaces payload: %+v", page)
				}
			}

			if tc.wantStatus == http.StatusBadRequest {
				assertMissingOrgIDErrorEnvelope(t, rec.Body.Bytes())
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

			req := httptest.NewRequest(http.MethodGet, "/api/workspaces/"+testWorkspaceID+"?"+tc.rawQuery, nil)
			req = mux.SetURLVars(req, map[string]string{"id": testWorkspaceID})
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

			if tc.wantStatus == http.StatusOK {
				if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
					t.Errorf("expected Content-Type application/json, got %q", ct)
				}
				var ws workspace.Workspace
				if err := json.Unmarshal(rec.Body.Bytes(), &ws); err != nil {
					t.Fatalf("failed to unmarshal workspace payload: %v (body=%q)", err, rec.Body.String())
				}
				if ws.ID.String() != testWorkspaceID || ws.Name != "test-workspace" {
					t.Errorf("unexpected workspace payload: %+v", ws)
				}
			}

			if tc.wantStatus == http.StatusBadRequest {
				assertMissingOrgIDErrorEnvelope(t, rec.Body.Bytes())
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

func (m *workspaceFailingProvider) GetWorkspaceByID(_ *http.Request, _, _ string) ([]byte, error) {
	return nil, m.err
}

func (m *workspaceFailingProvider) UpdateWorkspace(_ *http.Request, _ *workspace.WorkspaceUpdatePayload, _ string) (*workspace.Workspace, error) {
	return nil, m.err
}

func (m *workspaceFailingProvider) DeleteWorkspace(_ *http.Request, _ string) ([]byte, error) {
	return nil, m.err
}

func (m *workspaceFailingProvider) GetEnvironmentsOfWorkspace(_ *http.Request, _, _, _, _, _, _ string) ([]byte, error) {
	return nil, m.err
}

func (m *workspaceFailingProvider) GetDesignsOfWorkspace(_ *http.Request, _, _, _, _, _, _ string, _ []string) ([]byte, error) {
	return nil, m.err
}

func (m *workspaceFailingProvider) GetViewsOfWorkspace(_ *http.Request, _, _, _, _, _, _ string) ([]byte, error) {
	return nil, m.err
}

func (m *workspaceFailingProvider) GetTeamsOfWorkspace(_ *http.Request, _, _, _, _, _, _ string) ([]byte, error) {
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

// TestGetWorkspacesHandler_InvalidPagination_ReturnsBadRequest asserts that a
// rejected page or pagesize is answered with 400, not a provider-failure status.
func TestGetWorkspacesHandler_InvalidPagination_ReturnsBadRequest(t *testing.T) {
	cases := []struct {
		name     string
		query    string
		provider error
		wantCode string
	}{
		{"page", "page=-1", models.ErrPageNumber(errors.New("bad page")), models.ErrPageNumberCode},
		{"pagesize", "pagesize=abc", models.ErrPageSize(errors.New("bad pagesize")), models.ErrPageSizeCode},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			provider := newWorkspaceFailingProvider(tc.provider)

			req := httptest.NewRequest(http.MethodGet, "/api/workspaces?orgId="+testWorkspaceID+"&"+tc.query, nil)
			req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
			rec := httptest.NewRecorder()

			h.GetWorkspacesHandler(rec, req, nil, nil, provider)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusBadRequest, rec.Body.String())
			}
			assertPaginationErrorEnvelope(t, rec.Body.Bytes(), tc.wantCode)
		})
	}
}

// TestGetWorkspacesHandler_PropagatesProviderStatus covers the read path.
func TestGetWorkspacesHandler_PropagatesProviderStatus(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newWorkspaceFailingProvider(
		models.ErrFetch(errors.New("forbidden"), "Workspaces", http.StatusForbidden),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/workspaces?orgId="+testWorkspaceID, nil)
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

	body := `{"name":"team-space","organizationId":"` + testWorkspaceID + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/workspaces", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.SaveWorkspaceHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var decoded struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("failed to unmarshal created workspace response: %v (body=%q)", err, rec.Body.String())
	}
	if decoded.ID != testWorkspaceID {
		t.Errorf("expected workspace id %q, got %q", testWorkspaceID, decoded.ID)
	}
}

// invalidWorkspaceResourcePaginationCases drives the per-resource
// invalid-pagination tests below: the provider rejects the page or pagesize
// and the handler must answer 400 with that error's code.
var invalidWorkspaceResourcePaginationCases = []struct {
	name     string
	query    string
	provider error
	wantCode string
}{
	{"page", "page=-1", models.ErrPageNumber(errors.New("bad page")), models.ErrPageNumberCode},
	{"pagesize", "pagesize=abc", models.ErrPageSize(errors.New("bad pagesize")), models.ErrPageSizeCode},
}

// runWorkspaceResourceInvalidPagination calls one workspace sub-resource
// handler with each invalid pagination case and asserts 400 plus the code.
func runWorkspaceResourceInvalidPagination(t *testing.T, resource string, call func(h *Handler, w http.ResponseWriter, req *http.Request, provider models.Provider)) {
	t.Helper()
	for _, tc := range invalidWorkspaceResourcePaginationCases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			provider := newWorkspaceFailingProvider(tc.provider)

			req := httptest.NewRequest(http.MethodGet, "/api/workspaces/"+testWorkspaceID+"/"+resource+"?"+tc.query, nil)
			req = mux.SetURLVars(req, map[string]string{"id": testWorkspaceID})
			rec := httptest.NewRecorder()

			call(h, rec, req, provider)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusBadRequest, rec.Body.String())
			}
			assertPaginationErrorEnvelope(t, rec.Body.Bytes(), tc.wantCode)
		})
	}
}

func TestGetEnvironmentsOfWorkspaceHandler_InvalidPagination_ReturnsBadRequest(t *testing.T) {
	runWorkspaceResourceInvalidPagination(t, "environments", func(h *Handler, w http.ResponseWriter, req *http.Request, provider models.Provider) {
		h.GetEnvironmentsOfWorkspaceHandler(w, req, nil, nil, provider)
	})
}

func TestGetDesignsOfWorkspaceHandler_InvalidPagination_ReturnsBadRequest(t *testing.T) {
	runWorkspaceResourceInvalidPagination(t, "designs", func(h *Handler, w http.ResponseWriter, req *http.Request, provider models.Provider) {
		h.GetDesignsOfWorkspaceHandler(w, req, nil, nil, provider)
	})
}

func TestGetViewsOfWorkspaceHandler_InvalidPagination_ReturnsBadRequest(t *testing.T) {
	runWorkspaceResourceInvalidPagination(t, "views", func(h *Handler, w http.ResponseWriter, req *http.Request, provider models.Provider) {
		h.GetViewsOfWorkspaceHandler(w, req, nil, nil, provider)
	})
}

func TestGetTeamsOfWorkspaceHandler_InvalidPagination_ReturnsBadRequest(t *testing.T) {
	runWorkspaceResourceInvalidPagination(t, "teams", func(h *Handler, w http.ResponseWriter, req *http.Request, provider models.Provider) {
		h.GetTeamsOfWorkspaceHandler(w, req, nil, nil, provider)
	})
}

// TestWorkspaceHandlers_MissingWorkspace_ReturnsNotFound pins the 404 contract
// for a workspace that does not exist: the persister reports ErrResultNotFound
// and the handlers must surface it as 404 with a MeshKit code on the wire.
func TestWorkspaceHandlers_MissingWorkspace_ReturnsNotFound(t *testing.T) {
	const missingID = "33333333-3333-3333-3333-333333333333"
	const orgID = "11111111-1111-1111-1111-111111111111"

	assertNotFound := func(t *testing.T, rec *httptest.ResponseRecorder) {
		t.Helper()
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusNotFound, rec.Body.String())
		}
		var decoded struct {
			Code string `json:"code"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
			t.Fatalf("error body did not parse as JSON: %v (body=%q)", err, rec.Body.String())
		}
		if decoded.Code == "" {
			t.Errorf("expected a non-empty error code (body=%q)", rec.Body.String())
		}
	}

	t.Run("GetWorkspaceByIdHandler", func(t *testing.T) {
		h := newTestHandler(t, map[string]models.Provider{}, "")
		provider := newWorkspaceFailingProvider(models.ErrResultNotFound(gorm.ErrRecordNotFound))

		req := httptest.NewRequest(http.MethodGet, "/api/workspaces/"+missingID+"?orgId="+orgID, nil)
		req = mux.SetURLVars(req, map[string]string{"id": missingID})
		rec := httptest.NewRecorder()

		h.GetWorkspaceByIdHandler(rec, req, nil, nil, provider)

		assertNotFound(t, rec)
	})

	t.Run("UpdateWorkspaceHandler", func(t *testing.T) {
		h := newTestHandler(t, map[string]models.Provider{}, "")
		provider := newWorkspaceFailingProvider(models.ErrResultNotFound(gorm.ErrRecordNotFound))

		req := httptest.NewRequest(http.MethodPut, "/api/workspaces/"+missingID, strings.NewReader(`{"name":"x"}`))
		req = mux.SetURLVars(req, map[string]string{"id": missingID})
		rec := httptest.NewRecorder()

		h.UpdateWorkspaceHandler(rec, req, nil, nil, provider)

		assertNotFound(t, rec)
	})
}
