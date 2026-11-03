package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

// TestHandlers_MissingTokenInContextReturns500 verifies that handlers requiring
// an authentication token in context do not panic when the token is missing, but
// instead return a structured HTTP 500 error via ErrFetchToken.
func TestHandlers_MissingTokenInContextReturns500(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &models.DefaultLocalProvider{}
	provider.Initialize()

	t.Run("FetchResultsHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/results", nil)
		rec := httptest.NewRecorder()

		h.FetchResultsHandler(rec, req, nil, nil, provider)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d (body=%q)", rec.Code, rec.Body.String())
		}
	})

	t.Run("FetchAllResultsHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/perf/results", nil)
		rec := httptest.NewRecorder()

		h.FetchAllResultsHandler(rec, req, nil, nil, provider)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d (body=%q)", rec.Code, rec.Body.String())
		}
	})

	t.Run("GetResultHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/perf/result/11111111-1111-1111-1111-111111111111", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "11111111-1111-1111-1111-111111111111"})
		rec := httptest.NewRecorder()

		h.GetResultHandler(rec, req, nil, nil, provider)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d (body=%q)", rec.Code, rec.Body.String())
		}
	})

	t.Run("GetPerformanceProfilesHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/perf/profile", nil)
		rec := httptest.NewRecorder()

		h.GetPerformanceProfilesHandler(rec, req, nil, nil, provider)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d (body=%q)", rec.Code, rec.Body.String())
		}
	})

	t.Run("GetMesheryPatternsHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/pattern", nil)
		rec := httptest.NewRecorder()

		h.GetMesheryPatternsHandler(rec, req, nil, nil, provider)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d (body=%q)", rec.Code, rec.Body.String())
		}
	})

	t.Run("GetCatalogMesheryPatternsHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/catalog/pattern", nil)
		rec := httptest.NewRecorder()

		h.GetCatalogMesheryPatternsHandler(rec, req, nil, nil, provider)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d (body=%q)", rec.Code, rec.Body.String())
		}
	})

	t.Run("GetMesheryFiltersHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/filter", nil)
		rec := httptest.NewRecorder()

		h.GetMesheryFiltersHandler(rec, req, nil, nil, provider)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d (body=%q)", rec.Code, rec.Body.String())
		}
	})

	t.Run("GetCatalogMesheryFiltersHandler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/catalog/filter", nil)
		rec := httptest.NewRecorder()

		h.GetCatalogMesheryFiltersHandler(rec, req, nil, nil, provider)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d (body=%q)", rec.Code, rec.Body.String())
		}
	})
}
