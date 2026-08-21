package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type scheduleSpyProvider struct {
	*models.DefaultLocalProvider
	saveScheduleFn     func(string, *models.Schedule) ([]byte, error)
	getSchedulesFn     func(*http.Request, string, string, string) ([]byte, error)
	getScheduleFn      func(*http.Request, string) ([]byte, error)
	deleteScheduleFn   func(*http.Request, string) ([]byte, error)
	getProviderTokenFn func(*http.Request) (string, error)
}

func newScheduleSpyProvider() *scheduleSpyProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &scheduleSpyProvider{DefaultLocalProvider: base}
}

func (p *scheduleSpyProvider) SaveSchedule(token string, s *models.Schedule) ([]byte, error) {
	if p.saveScheduleFn != nil {
		return p.saveScheduleFn(token, s)
	}
	return p.DefaultLocalProvider.SaveSchedule(token, s)
}

func (p *scheduleSpyProvider) GetSchedules(req *http.Request, page, pageSize, order string) ([]byte, error) {
	if p.getSchedulesFn != nil {
		return p.getSchedulesFn(req, page, pageSize, order)
	}
	return p.DefaultLocalProvider.GetSchedules(req, page, pageSize, order)
}

func (p *scheduleSpyProvider) GetSchedule(req *http.Request, scheduleID string) ([]byte, error) {
	if p.getScheduleFn != nil {
		return p.getScheduleFn(req, scheduleID)
	}
	return p.DefaultLocalProvider.GetSchedule(req, scheduleID)
}

func (p *scheduleSpyProvider) DeleteSchedule(req *http.Request, scheduleID string) ([]byte, error) {
	if p.deleteScheduleFn != nil {
		return p.deleteScheduleFn(req, scheduleID)
	}
	return p.DefaultLocalProvider.DeleteSchedule(req, scheduleID)
}

func (p *scheduleSpyProvider) GetProviderToken(req *http.Request) (string, error) {
	if p.getProviderTokenFn != nil {
		return p.getProviderTokenFn(req)
	}
	return p.DefaultLocalProvider.GetProviderToken(req)
}

func TestSaveScheduleHandler(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())
	validBody := fmt.Sprintf(`{"name":"nightly","userId":"%s","cronExpression":"0 0 * * *"}`, userID)

	tests := []struct {
		name       string
		body       string
		setup      func(*scheduleSpyProvider)
		wantStatus int
		wantBody   string
	}{
		{
			name:       "invalid request body",
			body:       "{",
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "provider token error",
			body: validBody,
			setup: func(p *scheduleSpyProvider) {
				p.getProviderTokenFn = func(*http.Request) (string, error) {
					return "", fmt.Errorf("no token")
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "provider save error",
			body: validBody,
			setup: func(p *scheduleSpyProvider) {
				p.saveScheduleFn = func(string, *models.Schedule) ([]byte, error) {
					return nil, fmt.Errorf("save failed")
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "success",
			body: validBody,
			setup: func(p *scheduleSpyProvider) {
				p.saveScheduleFn = func(_ string, s *models.Schedule) ([]byte, error) {
					return []byte(fmt.Sprintf(`{"name":%q}`, s.Name)), nil
				}
			},
			wantStatus: http.StatusOK,
			wantBody:   `"name":"nightly"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			p := newScheduleSpyProvider()
			if tc.setup != nil {
				tc.setup(p)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/user/schedules", bytes.NewBufferString(tc.body))
			rec := httptest.NewRecorder()
			h.SaveScheduleHandler(rec, req, nil, nil, p)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantBody != "" {
				assert.Contains(t, rec.Body.String(), tc.wantBody)
			}
			if tc.wantStatus == http.StatusOK {
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			}
		})
	}
}

func TestGetSchedulesHandler(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(*scheduleSpyProvider)
		wantStatus int
		wantBody   string
	}{
		{
			name: "provider error",
			setup: func(p *scheduleSpyProvider) {
				p.getSchedulesFn = func(*http.Request, string, string, string) ([]byte, error) {
					return nil, fmt.Errorf("query failed")
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "success",
			setup: func(p *scheduleSpyProvider) {
				p.getSchedulesFn = func(_ *http.Request, page, pageSize, order string) ([]byte, error) {
					require.Equal(t, "1", page)
					require.Equal(t, "10", pageSize)
					require.Equal(t, "updated_at desc", order)
					return []byte(`{"schedules":[]}`), nil
				}
			},
			wantStatus: http.StatusOK,
			wantBody:   `"schedules"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			p := newScheduleSpyProvider()
			tc.setup(p)

			req := httptest.NewRequest(http.MethodGet, "/api/user/schedules?page=1&pagesize=10&order=updated_at+desc", nil)
			rec := httptest.NewRecorder()
			h.GetSchedulesHandler(rec, req, nil, nil, p)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantBody != "" {
				assert.Contains(t, rec.Body.String(), tc.wantBody)
			}
			if tc.wantStatus == http.StatusOK {
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			}
		})
	}
}

func TestDeleteScheduleHandler(t *testing.T) {
	scheduleID := uuid.Must(uuid.NewV4())

	tests := []struct {
		name       string
		setup      func(*scheduleSpyProvider)
		wantStatus int
	}{
		{
			name: "provider error",
			setup: func(p *scheduleSpyProvider) {
				p.deleteScheduleFn = func(*http.Request, string) ([]byte, error) {
					return nil, fmt.Errorf("delete failed")
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "success",
			setup: func(p *scheduleSpyProvider) {
				p.deleteScheduleFn = func(_ *http.Request, id string) ([]byte, error) {
					require.Equal(t, scheduleID.String(), id)
					return []byte(`{"deleted":true}`), nil
				}
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			p := newScheduleSpyProvider()
			tc.setup(p)

			req := httptest.NewRequest(http.MethodDelete, "/api/user/schedules/"+scheduleID.String(), nil)
			req = mux.SetURLVars(req, map[string]string{"id": scheduleID.String()})
			rec := httptest.NewRecorder()
			h.DeleteScheduleHandler(rec, req, nil, nil, p)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantStatus == http.StatusOK {
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
				assert.Contains(t, rec.Body.String(), `"deleted":true`)
			}
		})
	}
}

func TestGetScheduleHandler(t *testing.T) {
	scheduleID := uuid.Must(uuid.NewV4())

	tests := []struct {
		name       string
		setup      func(*scheduleSpyProvider)
		wantStatus int
	}{
		{
			name: "provider error",
			setup: func(p *scheduleSpyProvider) {
				p.getScheduleFn = func(*http.Request, string) ([]byte, error) {
					return nil, fmt.Errorf("not found")
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "success",
			setup: func(p *scheduleSpyProvider) {
				p.getScheduleFn = func(_ *http.Request, id string) ([]byte, error) {
					require.Equal(t, scheduleID.String(), id)
					return []byte(`{"name":"nightly"}`), nil
				}
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			p := newScheduleSpyProvider()
			tc.setup(p)

			req := httptest.NewRequest(http.MethodGet, "/api/user/schedules/"+scheduleID.String(), nil)
			req = mux.SetURLVars(req, map[string]string{"id": scheduleID.String()})
			rec := httptest.NewRecorder()
			h.GetScheduleHandler(rec, req, nil, nil, p)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantStatus == http.StatusOK {
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
				assert.Contains(t, rec.Body.String(), `"name":"nightly"`)
			}
		})
	}
}

func TestSaveScheduleHandler_PassesDecodedScheduleToProvider(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newScheduleSpyProvider()
	userID := uuid.Must(uuid.NewV4())

	var observed *models.Schedule
	p.saveScheduleFn = func(_ string, s *models.Schedule) ([]byte, error) {
		observed = s
		return []byte(`{}`), nil
	}

	body := fmt.Sprintf(`{"name":"weekly","userId":"%s","cronExpression":"0 0 * * 0"}`, userID)
	req := httptest.NewRequest(http.MethodPost, "/api/user/schedules", bytes.NewBufferString(body))
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "token"))
	rec := httptest.NewRecorder()
	h.SaveScheduleHandler(rec, req, nil, nil, p)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, observed)
	assert.Equal(t, "weekly", observed.Name)
	assert.Equal(t, "0 0 * * 0", observed.CronExpression)
}
