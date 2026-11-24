package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScheduleHandlers_TypedNilProvider(t *testing.T) {
	h := newScheduleTestHandler(t)
	var typedNilProv *models.DefaultLocalProvider

	t.Run("SaveScheduleHandler typed nil provider returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/user/schedules", strings.NewReader(`{"cron_expression":"* * * * *"}`))
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.SaveScheduleHandler(rec, req, nil, nil, typedNilProv)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("GetSchedulesHandler typed nil provider returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/schedules", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetSchedulesHandler(rec, req, nil, nil, typedNilProv)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("DeleteScheduleHandler typed nil provider returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/user/schedules/sched-1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "sched-1"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.DeleteScheduleHandler(rec, req, nil, nil, typedNilProv)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("GetScheduleHandler typed nil provider returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/schedules/sched-1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "sched-1"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetScheduleHandler(rec, req, nil, nil, typedNilProv)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})
}

func TestScheduleHandlers_NilLogger_DoesNotPanic(t *testing.T) {
	h := &Handler{
		config: &models.HandlerConfig{},
		log:    nil,
	}

	t.Run("SaveScheduleHandler nil logger does not panic on nil provider", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/user/schedules", strings.NewReader(`{"cron_expression":"* * * * *"}`))
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.SaveScheduleHandler(rec, req, nil, nil, nil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("SaveScheduleHandler nil logger does not panic on malformed body", func(t *testing.T) {
		prov := &models.DefaultLocalProvider{}
		req := httptest.NewRequest(http.MethodPost, "/api/user/schedules", strings.NewReader(`invalid-json`))
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.SaveScheduleHandler(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GetSchedulesHandler nil logger does not panic on nil provider", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/schedules", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetSchedulesHandler(rec, req, nil, nil, nil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("DeleteScheduleHandler nil logger does not panic on missing id", func(t *testing.T) {
		prov := &models.DefaultLocalProvider{}
		req := httptest.NewRequest(http.MethodDelete, "/api/user/schedules/", nil)
		req = mux.SetURLVars(req, map[string]string{"id": ""})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.DeleteScheduleHandler(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GetScheduleHandler nil logger does not panic on missing id", func(t *testing.T) {
		prov := &models.DefaultLocalProvider{}
		req := httptest.NewRequest(http.MethodGet, "/api/user/schedules/", nil)
		req = mux.SetURLVars(req, map[string]string{"id": ""})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetScheduleHandler(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
