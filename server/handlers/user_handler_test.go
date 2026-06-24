package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

// userSpyProvider embeds DefaultLocalProvider so we pick up a full Provider
// implementation, then overrides only GetUserByID with a scripted response.
// The handler test uses this to drive each error path independently.
type userSpyProvider struct {
	*models.DefaultLocalProvider
	resp []byte
	err  error
}

func (m *userSpyProvider) GetUserByID(_ *http.Request, _ string) ([]byte, error) {
	return m.resp, m.err
}

// TestGetUserByIDHandler_InvalidUUIDReturnsJSON pins the 400 response to a
// JSON envelope so the client's RTK Query parser sees a structured error
// object instead of a plain-text body (which previously crashed the JSON
// parser with "Unexpected token ... is not valid JSON").
func TestGetUserByIDHandler_InvalidUUIDReturnsJSON(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &userSpyProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	provider.Initialize()

	req := httptest.NewRequest(http.MethodGet, "/api/user/profile/not-a-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-uuid"})
	rec := httptest.NewRecorder()

	h.GetUserByIDHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	assertJSONErrorEnvelope(t, rec.Result())
}

// TestGetUserByIDHandler_ProviderErrorReturnsJSON covers the upstream-fetch
// failure path. The handler must still emit a structured JSON error body on
// 404 so the client does not choke on the meshkit "Unable to fetch data..."
// string that starts with the letter 'U'.
func TestGetUserByIDHandler_ProviderErrorReturnsJSON(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &userSpyProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}, err: fmt.Errorf("remote blew up")}
	provider.Initialize()

	req := httptest.NewRequest(http.MethodGet, "/api/user/profile/00000000-0000-0000-0000-000000000001", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000001"})
	rec := httptest.NewRecorder()

	h.GetUserByIDHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	assertJSONErrorEnvelope(t, rec.Result())
}

// TestGetUserByIDHandler_NilRespReturnsJSON exercises the "user not found"
// branch (provider returns (nil, nil) for a non-system ID). Both status and
// body shape must be locked in to JSON so the client behaviour is predictable.
func TestGetUserByIDHandler_NilRespReturnsJSON(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &userSpyProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	provider.Initialize()

	req := httptest.NewRequest(http.MethodGet, "/api/user/profile/00000000-0000-0000-0000-000000000002", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000002"})
	rec := httptest.NewRecorder()

	h.GetUserByIDHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	assertJSONErrorEnvelope(t, rec.Result())
}

// TestGetUserByIDHandler_SystemInstanceReturns204 verifies that when the
// provider signals ErrUserIsSystemInstance (the requested ID is this
// Meshery instance's own UUID, not a real user), the handler responds with
// 204 No Content. This lets the UI render a "system" placeholder without
// the row re-triggering a failed fetch forever — the original bug surface.
func TestGetUserByIDHandler_SystemInstanceReturns204(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &userSpyProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		err:                  models.ErrUserIsSystemInstance,
	}
	provider.Initialize()

	req := httptest.NewRequest(http.MethodGet, "/api/user/profile/00000000-0000-0000-0000-000000000003", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000003"})
	rec := httptest.NewRecorder()

	h.GetUserByIDHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for system instance ID, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	if body := rec.Body.Bytes(); len(body) != 0 {
		t.Errorf("expected empty body for 204, got %q", string(body))
	}
}

// assertJSONErrorEnvelope confirms that the handler error contract held:
// Content-Type is application/json (so RTK Query routes to its JSON handler)
// and the body parses to a JSON object with a non-empty "error" field. Some
// paths now use writeMeshkitError, which may additionally serialize MeshKit
// metadata fields like probable_cause as arrays.
func assertJSONErrorEnvelope(t *testing.T, resp *http.Response) {
	t.Helper()
	defer func() { _ = resp.Body.Close() }()

	if ct := resp.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("expected Content-Type application/json; charset=utf-8, got %q", ct)
	}
	if nosniff := resp.Header.Get("X-Content-Type-Options"); nosniff != "nosniff" {
		t.Errorf("expected X-Content-Type-Options: nosniff, got %q", nosniff)
	}

	var decoded map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("expected body to parse as JSON, got %v", err)
	}
	msg, ok := decoded["error"].(string)
	if !ok || msg == "" {
		t.Errorf("expected non-empty error field, got %+v", decoded)
	}
}
