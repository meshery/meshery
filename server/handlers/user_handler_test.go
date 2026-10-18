package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
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
// failure path. The handler must emit a structured JSON error body so the
// client does not choke on the meshkit "Unable to fetch data..." string that
// starts with the letter 'U'.
//
// The status is 502, not the 404 this asserted previously: a provider error
// that carries no HTTP status of its own is an upstream failure, and claiming
// "not found" invents a fact about the user record we have no evidence for.
// See providerStatus in handlers/utils.go.
func TestGetUserByIDHandler_ProviderErrorReturnsJSON(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &userSpyProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}, err: fmt.Errorf("remote blew up")}
	provider.Initialize()

	req := httptest.NewRequest(http.MethodGet, "/api/user/profile/00000000-0000-0000-0000-000000000001", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000001"})
	rec := httptest.NewRecorder()

	h.GetUserByIDHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	assertJSONErrorEnvelope(t, rec.Result())
}

// TestGetUserByIDHandler_PropagatesProviderStatus is the counterpart: when the
// provider DID report a status, that status must reach the client verbatim
// rather than being flattened to a fabricated 404.
func TestGetUserByIDHandler_PropagatesProviderStatus(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &userSpyProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		err:                  models.ErrFetch(fmt.Errorf("forbidden"), "User", http.StatusForbidden),
	}
	provider.Initialize()

	req := httptest.NewRequest(http.MethodGet, "/api/user/profile/00000000-0000-0000-0000-000000000001", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000001"})
	rec := httptest.NewRecorder()

	h.GetUserByIDHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d (body=%q)", rec.Code, rec.Body.String())
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

// TestUserPrefsHandler_RejectedUpdateDoesNotMutateStoredPreferences verifies that
// when an invalid preference update is rejected with HTTP 400, the stored in-memory
// preferences remain unchanged, and subsequent valid updates succeed.
func TestUserPrefsHandler_RejectedUpdateDoesNotMutateStoredPreferences(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	persister, err := models.NewMapPreferencePersister()
	if err != nil {
		t.Fatalf("failed to create persister: %v", err)
	}

	provider := &models.DefaultLocalProvider{
		MapPreferencePersister: persister,
	}
	provider.Initialize()

	userID := "00000000-0000-0000-0000-000000000001"
	user := &models.User{
		ID: uuid.FromStringOrNil(userID),
	}

	// 1. Establish initial valid preferences
	initialPref := models.NewDefaultPreference()
	initialPref.LoadTestPreferences = &models.LoadTestPreferences{
		QueriesPerSecond:   5,
		Duration:           "30s",
		ConcurrentRequests: 1,
		LoadGenerator:      "fortio",
	}
	initialPref.AnonymousUsageStats = true

	if err := provider.RecordPreferences(nil, userID, initialPref); err != nil {
		t.Fatalf("failed to record initial preferences: %v", err)
	}

	// Read initial preference as SessionInjectorMiddleware does
	prefObj, err := provider.ReadFromPersister(userID)
	if err != nil {
		t.Fatalf("failed to read preferences: %v", err)
	}

	// 2. Send an invalid update with QPS = -1 (must reject with HTTP 400)
	invalidBody := `{"loadTestPrefs": {"qps": -1}}`
	req := httptest.NewRequest(http.MethodPost, "/api/user/prefs", strings.NewReader(invalidBody))
	rec := httptest.NewRecorder()

	h.UserPrefsHandler(rec, req, prefObj, user, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected HTTP 400 for invalid QPS, got %d (body=%q)", rec.Code, rec.Body.String())
	}

	// Verify persisted preferences were NOT mutated by the rejected request
	persistedAfterRejection, err := provider.ReadFromPersister(userID)
	if err != nil {
		t.Fatalf("failed to read persisted preferences after rejection: %v", err)
	}
	if persistedAfterRejection.LoadTestPreferences == nil || persistedAfterRejection.LoadTestPreferences.QueriesPerSecond != 5 {
		t.Errorf("rejected update modified stored QPS: got %v, want 5", persistedAfterRejection.LoadTestPreferences)
	}

	// 3. Send a subsequent unrelated valid update (must succeed with HTTP 200)
	validBody := `{"anonymousUsageStats": false}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/user/prefs", strings.NewReader(validBody))
	rec2 := httptest.NewRecorder()

	prefObj2, err := provider.ReadFromPersister(userID)
	if err != nil {
		t.Fatalf("failed to read preferences for second request: %v", err)
	}

	h.UserPrefsHandler(rec2, req2, prefObj2, user, provider)

	if rec2.Code != http.StatusOK {
		t.Fatalf("subsequent valid update failed: got HTTP %d, want 200 (body=%q)", rec2.Code, rec2.Body.String())
	}

	// Verify updated preference was persisted and QPS remains intact
	persistedFinal, err := provider.ReadFromPersister(userID)
	if err != nil {
		t.Fatalf("failed to read final persisted preferences: %v", err)
	}
	if persistedFinal.AnonymousUsageStats != false {
		t.Errorf("expected AnonymousUsageStats to be false, got true")
	}
	if persistedFinal.LoadTestPreferences == nil || persistedFinal.LoadTestPreferences.QueriesPerSecond != 5 {
		t.Errorf("expected LoadTestPreferences.QueriesPerSecond to remain 5, got %v", persistedFinal.LoadTestPreferences)
	}
}
