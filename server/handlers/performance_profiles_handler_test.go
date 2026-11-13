package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func TestSavePerformanceProfileHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	body := []byte(`{"name":"test-profile"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/perf/profile", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.SavePerformanceProfileHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got %s", rec.Body.String())
	}
}

func TestGetPerformanceProfilesHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile", nil)
	rec := httptest.NewRecorder()

	h.GetPerformanceProfilesHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got %s", rec.Body.String())
	}
}

func TestDeletePerformanceProfileHandler_MissingID_ReturnsBadRequest(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodDelete, "/api/perf/profile/", nil)
	rec := httptest.NewRecorder()

	h.DeletePerformanceProfileHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "performanceProfileId is required") {
		t.Errorf("expected 'performanceProfileId is required', got %s", rec.Body.String())
	}
}

func TestDeletePerformanceProfileHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodDelete, "/api/perf/profile/profile-123", nil)
	req = mux.SetURLVars(req, map[string]string{"performanceProfileId": "profile-123"})
	rec := httptest.NewRecorder()

	h.DeletePerformanceProfileHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got %s", rec.Body.String())
	}
}

func TestGetPerformanceProfileHandler_MissingID_ReturnsBadRequest(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/", nil)
	rec := httptest.NewRecorder()

	h.GetPerformanceProfileHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "performanceProfileId is required") {
		t.Errorf("expected 'performanceProfileId is required', got %s", rec.Body.String())
	}
}

func TestGetPerformanceProfileHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/profile-123", nil)
	req = mux.SetURLVars(req, map[string]string{"performanceProfileId": "profile-123"})
	rec := httptest.NewRecorder()

	h.GetPerformanceProfileHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got %s", rec.Body.String())
	}
}
