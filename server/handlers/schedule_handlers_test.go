package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func newScheduleTestHandler(t *testing.T) *Handler {
	t.Helper()
	return &Handler{
		config: &models.HandlerConfig{},
		log:    newTestLogger(t),
	}
}

func TestScheduleHandlers_NilProvider(t *testing.T) {
	h := newScheduleTestHandler(t)

	// SaveScheduleHandler with nil provider
	{
		req := httptest.NewRequest(http.MethodPost, "/api/user/schedules", strings.NewReader(`{"cron_expression":"* * * * *"}`))
		rec := httptest.NewRecorder()
		h.SaveScheduleHandler(rec, req, nil, nil, nil)
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("SaveScheduleHandler status = %d, want 500; body: %s", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), ErrNilProviderCode) {
			t.Errorf("SaveScheduleHandler expected error to contain %s, got: %s", ErrNilProviderCode, rec.Body.String())
		}
	}

	// GetSchedulesHandler with nil provider
	{
		req := httptest.NewRequest(http.MethodGet, "/api/user/schedules", nil)
		rec := httptest.NewRecorder()
		h.GetSchedulesHandler(rec, req, nil, nil, nil)
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("GetSchedulesHandler status = %d, want 500; body: %s", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), ErrNilProviderCode) {
			t.Errorf("GetSchedulesHandler expected error to contain %s, got: %s", ErrNilProviderCode, rec.Body.String())
		}
	}

	// DeleteScheduleHandler with nil provider
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/user/schedules/sched-1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "sched-1"})
		rec := httptest.NewRecorder()
		h.DeleteScheduleHandler(rec, req, nil, nil, nil)
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("DeleteScheduleHandler status = %d, want 500; body: %s", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), ErrNilProviderCode) {
			t.Errorf("DeleteScheduleHandler expected error to contain %s, got: %s", ErrNilProviderCode, rec.Body.String())
		}
	}

	// GetScheduleHandler with nil provider
	{
		req := httptest.NewRequest(http.MethodGet, "/api/user/schedules/sched-1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "sched-1"})
		rec := httptest.NewRecorder()
		h.GetScheduleHandler(rec, req, nil, nil, nil)
		if rec.Code != http.StatusInternalServerError {
			t.Errorf("GetScheduleHandler status = %d, want 500; body: %s", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), ErrNilProviderCode) {
			t.Errorf("GetScheduleHandler expected error to contain %s, got: %s", ErrNilProviderCode, rec.Body.String())
		}
	}
}

func TestScheduleHandlers_EmptyIDValidation(t *testing.T) {
	h := newScheduleTestHandler(t)
	mockProv := &models.DefaultLocalProvider{}

	// DeleteScheduleHandler with empty ID
	{
		req := httptest.NewRequest(http.MethodDelete, "/api/user/schedules/", nil)
		req = mux.SetURLVars(req, map[string]string{"id": ""})
		rec := httptest.NewRecorder()
		h.DeleteScheduleHandler(rec, req, nil, nil, mockProv)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("DeleteScheduleHandler status = %d, want 400 for empty id; body: %s", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), ErrMissingRouteVariableCode) {
			t.Errorf("DeleteScheduleHandler expected error to contain %s, got: %s", ErrMissingRouteVariableCode, rec.Body.String())
		}
	}

	// GetScheduleHandler with empty ID
	{
		req := httptest.NewRequest(http.MethodGet, "/api/user/schedules/", nil)
		req = mux.SetURLVars(req, map[string]string{"id": ""})
		rec := httptest.NewRecorder()
		h.GetScheduleHandler(rec, req, nil, nil, mockProv)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("GetScheduleHandler status = %d, want 400 for empty id; body: %s", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), ErrMissingRouteVariableCode) {
			t.Errorf("GetScheduleHandler expected error to contain %s, got: %s", ErrMissingRouteVariableCode, rec.Body.String())
		}
	}

	// SaveScheduleHandler with empty body (io.EOF from empty reader)
	{
		req := httptest.NewRequest(http.MethodPost, "/api/user/schedules", strings.NewReader(""))
		rec := httptest.NewRecorder()
		h.SaveScheduleHandler(rec, req, nil, nil, mockProv)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("SaveScheduleHandler status = %d, want 400 for empty body; body: %s", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
			t.Errorf("SaveScheduleHandler expected error to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
		}
	}

	// SaveScheduleHandler with null JSON body (parsedBody == nil guard)
	{
		req := httptest.NewRequest(http.MethodPost, "/api/user/schedules", strings.NewReader("null"))
		rec := httptest.NewRecorder()
		h.SaveScheduleHandler(rec, req, nil, nil, mockProv)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("SaveScheduleHandler status = %d, want 400 for null body; body: %s", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
			t.Errorf("SaveScheduleHandler expected error to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
		}
	}
}
