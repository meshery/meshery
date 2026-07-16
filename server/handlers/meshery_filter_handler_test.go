package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

// filterSpyProvider embeds DefaultLocalProvider and tracks whether
// provider methods were called, so tests can assert that invalid
// filter IDs are rejected before reaching the provider layer.
type filterSpyProvider struct {
	*models.DefaultLocalProvider
	called bool
}

func (m *filterSpyProvider) GetMesheryFilterFile(_ *http.Request, _ string) ([]byte, error) {
	m.called = true
	return nil, nil
}

func (m *filterSpyProvider) DeleteMesheryFilter(_ *http.Request, _ string) ([]byte, error) {
	m.called = true
	return nil, nil
}

func (m *filterSpyProvider) GetMesheryFilter(_ *http.Request, _ string) ([]byte, error) {
	m.called = true
	return nil, nil
}

func newFilterSpyProvider() *filterSpyProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &filterSpyProvider{DefaultLocalProvider: base}
}

func TestGetMesheryFilterFileHandler_InvalidUUIDReturns400(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newFilterSpyProvider()

	req := httptest.NewRequest(http.MethodGet, "/api/filter/file/not-a-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-uuid"})
	rec := httptest.NewRecorder()

	h.GetMesheryFilterFileHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if provider.called {
		t.Error("provider should not be called for invalid UUID")
	}
}

func TestGetMesheryFilterFileHandler_NilUUIDReturns400(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newFilterSpyProvider()

	req := httptest.NewRequest(http.MethodGet, "/api/filter/file/00000000-0000-0000-0000-000000000000", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000000"})
	rec := httptest.NewRecorder()

	h.GetMesheryFilterFileHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if provider.called {
		t.Error("provider should not be called for nil UUID")
	}
}

func TestDeleteMesheryFilterHandler_InvalidUUIDReturns400(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newFilterSpyProvider()

	req := httptest.NewRequest(http.MethodDelete, "/api/filter/not-a-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-uuid"})
	rec := httptest.NewRecorder()

	h.DeleteMesheryFilterHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if provider.called {
		t.Error("provider should not be called for invalid UUID")
	}
}

func TestDeleteMesheryFilterHandler_NilUUIDReturns400(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newFilterSpyProvider()

	req := httptest.NewRequest(http.MethodDelete, "/api/filter/00000000-0000-0000-0000-000000000000", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000000"})
	rec := httptest.NewRecorder()

	h.DeleteMesheryFilterHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if provider.called {
		t.Error("provider should not be called for nil UUID")
	}
}

func TestGetMesheryFilterHandler_InvalidUUIDReturns400(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newFilterSpyProvider()

	req := httptest.NewRequest(http.MethodGet, "/api/filter/not-a-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "not-a-uuid"})
	rec := httptest.NewRecorder()

	h.GetMesheryFilterHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if provider.called {
		t.Error("provider should not be called for invalid UUID")
	}
}

func TestGetMesheryFilterHandler_NilUUIDReturns400(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newFilterSpyProvider()

	req := httptest.NewRequest(http.MethodGet, "/api/filter/00000000-0000-0000-0000-000000000000", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000000"})
	rec := httptest.NewRecorder()

	h.GetMesheryFilterHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
	if provider.called {
		t.Error("provider should not be called for nil UUID")
	}
}
