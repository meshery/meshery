package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerformanceProfileHandlers_TypedNilProvider(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	var typedNil *models.DefaultLocalProvider

	t.Run("SavePerformanceProfileHandler with typed nil provider returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/perf/profile", bytes.NewReader([]byte(`{"name":"test-profile"}`)))
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.SavePerformanceProfileHandler(rec, req, nil, nil, typedNil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("GetPerformanceProfilesHandler with typed nil provider returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/perf/profile", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetPerformanceProfilesHandler(rec, req, nil, nil, typedNil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("DeletePerformanceProfileHandler with typed nil provider returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/perf/profile/profile-123", nil)
		req = mux.SetURLVars(req, map[string]string{"performanceProfileId": "profile-123"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.DeletePerformanceProfileHandler(rec, req, nil, nil, typedNil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("GetPerformanceProfileHandler with typed nil provider returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/profile-123", nil)
		req = mux.SetURLVars(req, map[string]string{"performanceProfileId": "profile-123"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetPerformanceProfileHandler(rec, req, nil, nil, typedNil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})
}

func TestPerformanceProfileHandlers_NilLogger_DoesNotPanic(t *testing.T) {
	h := &Handler{
		config: &models.HandlerConfig{},
		log:    nil,
	}

	t.Run("SavePerformanceProfileHandler with nil logger does not panic on nil provider", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/perf/profile", bytes.NewReader([]byte(`{"name":"test-profile"}`)))
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.SavePerformanceProfileHandler(rec, req, nil, nil, nil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("SavePerformanceProfileHandler with nil logger does not panic on malformed body", func(t *testing.T) {
		prov := &models.DefaultLocalProvider{}
		req := httptest.NewRequest(http.MethodPost, "/api/perf/profile", bytes.NewReader([]byte(`invalid-json`)))
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.SavePerformanceProfileHandler(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("SavePerformanceProfileHandler with nil logger does not panic on empty body", func(t *testing.T) {
		prov := &models.DefaultLocalProvider{}
		req := httptest.NewRequest(http.MethodPost, "/api/perf/profile", bytes.NewReader([]byte(``)))
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.SavePerformanceProfileHandler(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GetPerformanceProfilesHandler with nil logger does not panic on nil provider", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/perf/profile", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetPerformanceProfilesHandler(rec, req, nil, nil, nil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("DeletePerformanceProfileHandler with nil logger does not panic on invalid UUID", func(t *testing.T) {
		prov := &models.DefaultLocalProvider{}
		req := httptest.NewRequest(http.MethodDelete, "/api/perf/profile/invalid-uuid", nil)
		req = mux.SetURLVars(req, map[string]string{"performanceProfileId": "invalid-uuid"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.DeletePerformanceProfileHandler(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GetPerformanceProfileHandler with nil logger does not panic on invalid UUID", func(t *testing.T) {
		prov := &models.DefaultLocalProvider{}
		req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/invalid-uuid", nil)
		req = mux.SetURLVars(req, map[string]string{"performanceProfileId": "invalid-uuid"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetPerformanceProfileHandler(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
