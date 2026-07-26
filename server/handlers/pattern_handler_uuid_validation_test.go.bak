package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

// patternSpyProvider embeds DefaultLocalProvider and tracks whether
// GetMesheryPattern was called. Tests use this to assert the provider
// is never reached when the path ID is invalid.
type patternSpyProvider struct {
	*models.DefaultLocalProvider
	called bool
}

func (p *patternSpyProvider) GetMesheryPattern(_ *http.Request, _ string, _ string) ([]byte, error) {
	p.called = true
	return nil, nil
}

func newPatternSpyProvider() *patternSpyProvider {
	sp := &patternSpyProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	sp.Initialize()
	return sp
}

// TestGetMesheryPatternHandler_InvalidIDReturns400 asserts that
// GetMesheryPatternHandler returns 400 Bad Request for invalid design
// IDs without invoking the provider.
func TestGetMesheryPatternHandler_InvalidIDReturns400(t *testing.T) {
	cases := []struct {
		name string
		id   string
	}{
		{"undefined string", "undefined"},
		{"empty string", ""},
		{"nil UUID", "00000000-0000-0000-0000-000000000000"},
		{"garbage string", "not-a-uuid"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			spy := newPatternSpyProvider()

			req := httptest.NewRequest(http.MethodGet, "/api/pattern/"+tc.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})
			rec := httptest.NewRecorder()

			h.GetMesheryPatternHandler(rec, req, nil, &models.User{}, spy)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d (body=%q)", rec.Code, rec.Body.String())
			}
			if spy.called {
				t.Error("provider.GetMesheryPattern was called; expected early return on invalid ID")
			}
		})
	}
}

// TestGetMesheryPatternHandler_ValidIDProceedsToProvider asserts that
// a well-formed, non-nil UUID passes the guard and reaches the provider.
// (Provider returns nil here; the test only checks the guard didn't fire.)
func TestGetMesheryPatternHandler_ValidIDProceedsToProvider(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	spy := newPatternSpyProvider()

	validID := "550e8400-e29b-41d4-a716-446655440000"
	req := httptest.NewRequest(http.MethodGet, "/api/pattern/"+validID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": validID})
	rec := httptest.NewRecorder()

	h.GetMesheryPatternHandler(rec, req, nil, &models.User{}, spy)

	// 400 must NOT be the UUID-guard 400 — if we got here, the guard passed.
	if rec.Code == http.StatusBadRequest {
		t.Errorf("valid UUID should not trigger 400 guard (body=%q)", rec.Body.String())
	}
}

// TestCloneMesheryPatternHandler_InvalidIDReturns400 asserts that
// CloneMesheryPatternHandler returns 400 Bad Request for invalid design
// IDs without invoking the provider.
func TestCloneMesheryPatternHandler_InvalidIDReturns400(t *testing.T) {
	cases := []struct {
		name string
		id   string
	}{
		{"undefined string", "undefined"},
		{"empty string", ""},
		{"nil UUID", "00000000-0000-0000-0000-000000000000"},
		{"garbage string", "not-a-uuid"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			spy := newPatternSpyProvider()

			req := httptest.NewRequest(http.MethodPost, "/api/pattern/clone/"+tc.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})
			rec := httptest.NewRecorder()

			h.CloneMesheryPatternHandler(rec, req, nil, &models.User{}, spy)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d (body=%q)", rec.Code, rec.Body.String())
			}
			if spy.called {
				t.Error("provider.GetMesheryPattern was called; expected early return on invalid ID")
			}
		})
	}
}
