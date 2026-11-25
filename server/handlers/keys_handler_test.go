package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	"github.com/meshery/meshery/server/models"
)

// recordingKeysProvider scripts GetUsersKeys and records the arguments the
// handler passed through, so the tests can assert both the response mapping and
// that validation did not alter a valid request. called distinguishes "the
// handler never reached the provider" from "the handler called it with an empty
// organization ID", which a gotOrgID == "" assertion cannot tell apart.
type recordingKeysProvider struct {
	*models.DefaultLocalProvider
	resp     []byte
	err      error
	called   bool
	gotToken string
	gotOrgID string
}

func newRecordingKeysProvider(resp []byte, err error) *recordingKeysProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &recordingKeysProvider{DefaultLocalProvider: base, resp: resp, err: err}
}

func (p *recordingKeysProvider) GetUsersKeys(token, _, _, _, _, _ string, orgID string) ([]byte, error) {
	p.called = true
	p.gotToken = token
	p.gotOrgID = orgID
	return p.resp, p.err
}

// TestGetUsersKeys_NilProviderIsServerError pins the guard: a nil provider used
// to be dereferenced directly, so a provider that failed to initialize took the
// whole server down with a SIGSEGV instead of failing the one request.
func TestGetUsersKeys_NilProviderIsServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/orgs/"+uuid.Must(uuid.NewV4()).String()+"/users/keys", nil))
	rec := httptest.NewRecorder()

	h.GetUsersKeys(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusInternalServerError, rec.Body.String())
	}
	if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != ErrNilProviderCode {
		t.Errorf("code = %q, want %q", decoded.Code, ErrNilProviderCode)
	}
}

// TestGetUsersKeys_InvalidOrganizationIDIsBadRequest verifies that a malformed
// organization identifier is rejected at the boundary instead of being handed
// to the provider, where it surfaced as an unclassified downstream error.
func TestGetUsersKeys_InvalidOrganizationIDIsBadRequest(t *testing.T) {
	for _, orgID := range []string{"not-a-uuid", "12345", "abc/def", uuid.Nil.String()} {
		t.Run(orgID, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			provider := newRecordingKeysProvider([]byte(`{"keys":[]}`), nil)

			req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/orgs/"+orgID+"/users/keys", nil))
			req = mux.SetURLVars(req, map[string]string{"orgID": orgID})
			rec := httptest.NewRecorder()

			h.GetUsersKeys(rec, req, nil, nil, provider)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusBadRequest, rec.Body.String())
			}
			if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
				t.Errorf("code = %q, want %q", decoded.Code, models.ErrInvalidUUIDValueCode)
			} else if orgID == uuid.Nil.String() && !strings.Contains(rec.Body.String(), "must not be the nil UUID") {
				t.Errorf("body = %q, want it to contain %q", rec.Body.String(), "must not be the nil UUID")
			}
			if provider.called {
				t.Error("provider was called, want no call at all")
			}
		})
	}
}

// TestGetUsersKeys_ValidOrganizationIDReachesProvider keeps the validation
// honest: a well-formed identifier must pass through untouched, together with
// the request token.
func TestGetUsersKeys_ValidOrganizationIDReachesProvider(t *testing.T) {
	tests := []struct {
		name  string
		orgID string
	}{
		{name: "organization id", orgID: uuid.Must(uuid.NewV4()).String()},
		{name: "no organization id", orgID: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			provider := newRecordingKeysProvider([]byte(`{"keys":[]}`), nil)

			req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/orgs/"+tc.orgID+"/users/keys", nil))
			req = mux.SetURLVars(req, map[string]string{"orgID": tc.orgID})
			rec := httptest.NewRecorder()

			h.GetUsersKeys(rec, req, nil, nil, provider)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusOK, rec.Body.String())
			}
			if got := rec.Body.String(); got != `{"keys":[]}` {
				t.Errorf("body = %q, want the provider payload verbatim", got)
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("content-type = %q, want application/json", ct)
			}
			if provider.gotOrgID != tc.orgID {
				t.Errorf("provider got orgID %q, want %q", provider.gotOrgID, tc.orgID)
			}
			if provider.gotToken != "test-token" {
				t.Errorf("provider got token %q, want %q", provider.gotToken, "test-token")
			}
		})
	}
}

// TestGetUsersKeys_NotFoundIs404 pins the status this endpoint returns when the
// organization does not exist, for both shapes the provider layer reports it.
//
// No branch in the handler is needed for this. models.ErrResultNotFound tags the
// error with httputil.WithProviderStatus(..., http.StatusNotFound)
// (server/models/error.go), and providerStatus (server/handlers/utils.go) reads
// that tag back through httputil.StatusForProviderError, so the status survives
// the ErrGetUsersKeys envelope. The test pins that so a later change to either
// piece cannot silently turn a missing organization into a 502.
func TestGetUsersKeys_NotFoundIs404(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "local provider reports the record is missing", err: models.ErrResultNotFound(errors.New("organization not found"))},
		{name: "remote provider answers 404", err: models.ErrFetch(errors.New("failed to fetch keys"), "not found", http.StatusNotFound)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			provider := newRecordingKeysProvider(nil, tc.err)

			orgID := uuid.Must(uuid.NewV4()).String()
			req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/orgs/"+orgID+"/users/keys", nil))
			req = mux.SetURLVars(req, map[string]string{"orgID": orgID})
			rec := httptest.NewRecorder()

			h.GetUsersKeys(rec, req, nil, nil, provider)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusNotFound, rec.Body.String())
			}
		})
	}
}

// TestGetUsersKeys_NilLogger_DoesNotPanic verifies that the handler guards
// against nil logger dereferences across all error and success branches.
func TestGetUsersKeys_NilLogger_DoesNotPanic(t *testing.T) {
	h := &Handler{} // nil logger

	t.Run("missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/identity/orgs/users/keys", nil)
		rec := httptest.NewRecorder()
		h.GetUsersKeys(rec, req, nil, nil, nil)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})

	t.Run("nil provider", func(t *testing.T) {
		req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/orgs/users/keys", nil))
		rec := httptest.NewRecorder()
		h.GetUsersKeys(rec, req, nil, nil, nil)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})

	t.Run("invalid orgID", func(t *testing.T) {
		req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/orgs/invalid/users/keys", nil))
		req = mux.SetURLVars(req, map[string]string{"orgID": "invalid"})
		rec := httptest.NewRecorder()
		provider := newRecordingKeysProvider([]byte(`{"keys":[]}`), nil)
		h.GetUsersKeys(rec, req, nil, nil, provider)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("zero orgID", func(t *testing.T) {
		req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/orgs/"+uuid.Nil.String()+"/users/keys", nil))
		req = mux.SetURLVars(req, map[string]string{"orgID": uuid.Nil.String()})
		rec := httptest.NewRecorder()
		provider := newRecordingKeysProvider([]byte(`{"keys":[]}`), nil)
		h.GetUsersKeys(rec, req, nil, nil, provider)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("provider error", func(t *testing.T) {
		req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/orgs/users/keys", nil))
		rec := httptest.NewRecorder()
		provider := newRecordingKeysProvider(nil, errors.New("upstream failed"))
		h.GetUsersKeys(rec, req, nil, nil, provider)
		if rec.Code != http.StatusBadGateway {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadGateway)
		}
	})

	t.Run("success", func(t *testing.T) {
		req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/orgs/users/keys", nil))
		rec := httptest.NewRecorder()
		provider := newRecordingKeysProvider([]byte(`{"keys":[]}`), nil)
		h.GetUsersKeys(rec, req, nil, nil, provider)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})
}
