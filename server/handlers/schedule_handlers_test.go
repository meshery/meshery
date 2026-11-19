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
	meshkiterrors "github.com/meshery/meshkit/errors"
	"gorm.io/gorm"
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

type mockScheduleProvider struct {
	*models.DefaultLocalProvider
	deleteScheduleFunc func(r *http.Request, id string) ([]byte, error)
	getScheduleFunc    func(r *http.Request, id string) ([]byte, error)
	saveScheduleFunc   func(token string, schedule *models.Schedule) ([]byte, error)
	getSchedulesFunc   func(r *http.Request, page, pagesize, order string) ([]byte, error)
}

func (m *mockScheduleProvider) DeleteSchedule(r *http.Request, id string) ([]byte, error) {
	if m.deleteScheduleFunc != nil {
		return m.deleteScheduleFunc(r, id)
	}
	return m.DefaultLocalProvider.DeleteSchedule(r, id)
}

func (m *mockScheduleProvider) GetSchedule(r *http.Request, id string) ([]byte, error) {
	if m.getScheduleFunc != nil {
		return m.getScheduleFunc(r, id)
	}
	return m.DefaultLocalProvider.GetSchedule(r, id)
}

func (m *mockScheduleProvider) SaveSchedule(token string, schedule *models.Schedule) ([]byte, error) {
	if m.saveScheduleFunc != nil {
		return m.saveScheduleFunc(token, schedule)
	}
	return m.DefaultLocalProvider.SaveSchedule(token, schedule)
}

func (m *mockScheduleProvider) GetSchedules(r *http.Request, page, pagesize, order string) ([]byte, error) {
	if m.getSchedulesFunc != nil {
		return m.getSchedulesFunc(r, page, pagesize, order)
	}
	return m.DefaultLocalProvider.GetSchedules(r, page, pagesize, order)
}

func TestScheduleHandlers_InvalidUUIDValidation(t *testing.T) {
	h := newScheduleTestHandler(t)
	mockProv := &models.DefaultLocalProvider{}

	testCases := []struct {
		name    string
		handler func(w http.ResponseWriter, r *http.Request)
		method  string
		url     string
		idVal   string
	}{
		{
			name: "DeleteScheduleHandler with malformed UUID",
			handler: func(w http.ResponseWriter, r *http.Request) {
				h.DeleteScheduleHandler(w, r, nil, nil, mockProv)
			},
			method: http.MethodDelete,
			url:    "/api/user/schedules/not-a-uuid",
			idVal:  "not-a-uuid",
		},
		{
			name: "DeleteScheduleHandler with nil UUID",
			handler: func(w http.ResponseWriter, r *http.Request) {
				h.DeleteScheduleHandler(w, r, nil, nil, mockProv)
			},
			method: http.MethodDelete,
			url:    "/api/user/schedules/" + uuid.Nil.String(),
			idVal:  uuid.Nil.String(),
		},
		{
			name: "DeleteScheduleHandler with compact nil UUID",
			handler: func(w http.ResponseWriter, r *http.Request) {
				h.DeleteScheduleHandler(w, r, nil, nil, mockProv)
			},
			method: http.MethodDelete,
			url:    "/api/user/schedules/00000000000000000000000000000000",
			idVal:  "00000000000000000000000000000000",
		},
		{
			name: "GetScheduleHandler with malformed UUID",
			handler: func(w http.ResponseWriter, r *http.Request) {
				h.GetScheduleHandler(w, r, nil, nil, mockProv)
			},
			method: http.MethodGet,
			url:    "/api/user/schedules/not-a-uuid",
			idVal:  "not-a-uuid",
		},
		{
			name: "GetScheduleHandler with nil UUID",
			handler: func(w http.ResponseWriter, r *http.Request) {
				h.GetScheduleHandler(w, r, nil, nil, mockProv)
			},
			method: http.MethodGet,
			url:    "/api/user/schedules/" + uuid.Nil.String(),
			idVal:  uuid.Nil.String(),
		},
		{
			name: "GetScheduleHandler with compact nil UUID",
			handler: func(w http.ResponseWriter, r *http.Request) {
				h.GetScheduleHandler(w, r, nil, nil, mockProv)
			},
			method: http.MethodGet,
			url:    "/api/user/schedules/00000000000000000000000000000000",
			idVal:  "00000000000000000000000000000000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tc.idVal})
			rec := httptest.NewRecorder()
			tc.handler(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("%s status = %d, want 400; body: %s", tc.name, rec.Code, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
				t.Errorf("%s expected error to contain %s, got: %s", tc.name, models.ErrInvalidUUIDValueCode, rec.Body.String())
			}
		})
	}
}

func TestScheduleHandlers_NotFoundAndStatusPropagation(t *testing.T) {
	h := newScheduleTestHandler(t)
	validID := uuid.Must(uuid.NewV4()).String()

	t.Run("DeleteScheduleHandler returns 404 when gorm.ErrRecordNotFound", func(t *testing.T) {
		prov := &mockScheduleProvider{
			DefaultLocalProvider: &models.DefaultLocalProvider{},
			deleteScheduleFunc: func(_ *http.Request, _ string) ([]byte, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		req := httptest.NewRequest(http.MethodDelete, "/api/user/schedules/"+validID, nil)
		req = mux.SetURLVars(req, map[string]string{"id": validID})
		rec := httptest.NewRecorder()
		h.DeleteScheduleHandler(rec, req, nil, nil, prov)

		if rec.Code != http.StatusNotFound {
			t.Errorf("DeleteScheduleHandler status = %d, want 404; body: %s", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), models.ErrResultNotFoundCode) {
			t.Errorf("DeleteScheduleHandler expected error to contain %s, got: %s", models.ErrResultNotFoundCode, rec.Body.String())
		}
	})

	t.Run("GetScheduleHandler returns 404 when ErrResultNotFound", func(t *testing.T) {
		prov := &mockScheduleProvider{
			DefaultLocalProvider: &models.DefaultLocalProvider{},
			getScheduleFunc: func(_ *http.Request, _ string) ([]byte, error) {
				return nil, models.ErrResultNotFound(errors.New("not found"))
			},
		}
		req := httptest.NewRequest(http.MethodGet, "/api/user/schedules/"+validID, nil)
		req = mux.SetURLVars(req, map[string]string{"id": validID})
		rec := httptest.NewRecorder()
		h.GetScheduleHandler(rec, req, nil, nil, prov)

		if rec.Code != http.StatusNotFound {
			t.Errorf("GetScheduleHandler status = %d, want 404; body: %s", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), models.ErrResultNotFoundCode) {
			t.Errorf("GetScheduleHandler expected error to contain %s, got: %s", models.ErrResultNotFoundCode, rec.Body.String())
		}
	})

	t.Run("GetScheduleHandler returns 404 when gorm.ErrRecordNotFound", func(t *testing.T) {
		prov := &mockScheduleProvider{
			DefaultLocalProvider: &models.DefaultLocalProvider{},
			getScheduleFunc: func(_ *http.Request, _ string) ([]byte, error) {
				return nil, gorm.ErrRecordNotFound
			},
		}
		req := httptest.NewRequest(http.MethodGet, "/api/user/schedules/"+validID, nil)
		req = mux.SetURLVars(req, map[string]string{"id": validID})
		rec := httptest.NewRecorder()
		h.GetScheduleHandler(rec, req, nil, nil, prov)

		if rec.Code != http.StatusNotFound {
			t.Errorf("GetScheduleHandler status = %d, want 404; body: %s", rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), models.ErrResultNotFoundCode) {
			t.Errorf("GetScheduleHandler expected error to contain %s, got: %s", models.ErrResultNotFoundCode, rec.Body.String())
		}
	})

	t.Run("GetSchedulesHandler returns 400 for pagination error", func(t *testing.T) {
		prov := &mockScheduleProvider{
			DefaultLocalProvider: &models.DefaultLocalProvider{},
			getSchedulesFunc: func(_ *http.Request, _, _, _ string) ([]byte, error) {
				return nil, meshkiterrors.New(models.ErrPageSizeCode, meshkiterrors.Alert, []string{"page size invalid"}, nil, nil, nil)
			},
		}
		req := httptest.NewRequest(http.MethodGet, "/api/user/schedules?pagesize=invalid", nil)
		rec := httptest.NewRecorder()
		h.GetSchedulesHandler(rec, req, nil, nil, prov)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("GetSchedulesHandler status = %d, want 400; body: %s", rec.Code, rec.Body.String())
		}
	})
}
