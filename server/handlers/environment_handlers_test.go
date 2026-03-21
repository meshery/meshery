package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
)

func TestGetEnvironmentByIDHandler(t *testing.T) {
	h := newTestHandler(t, nil, "")
	provider := &models.DefaultLocalProvider{}

	tests := []struct {
		name       string
		id         string
		expectedSC int
	}{
		{
			name:       "given non valid uuid when GetEnvironmentByIDHandler return status code 400",
			id:         "not-a-uuid",
			expectedSC: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/environments/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			rr := httptest.NewRecorder()

			h.GetEnvironmentByIDHandler(rr, req, nil, nil, provider)

			assert.Equal(t, tt.expectedSC, rr.Code)
		})
	}
}

func TestUpdateEnvironmentHandler(t *testing.T) {
	h := newTestHandler(t, nil, "")
	provider := &models.DefaultLocalProvider{}

	tests := []struct {
		name       string
		id         string
		payload    string
		expectedSC int
	}{
		{
			name:       "given non valid uuid when UpdateEnvironmentHandler return status code 400",
			id:         "not-a-uuid",
			payload:    `{"name":"env-1"}`,
			expectedSC: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/environments/"+tt.id, strings.NewReader(tt.payload))
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			rr := httptest.NewRecorder()

			h.UpdateEnvironmentHandler(rr, req, nil, nil, provider)

			assert.Equal(t, tt.expectedSC, rr.Code)
		})
	}
}

func TestDeleteEnvironmentHandler(t *testing.T) {
	h := newTestHandler(t, nil, "")
	provider := &models.DefaultLocalProvider{}

	tests := []struct {
		name       string
		id         string
		expectedSC int
	}{
		{
			name:       "given non valid uuid when DeleteEnvironmentHandler return status code 400",
			id:         "not-a-uuid",
			expectedSC: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/api/environments/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			rr := httptest.NewRecorder()

			h.DeleteEnvironmentHandler(rr, req, nil, nil, provider)

			assert.Equal(t, tt.expectedSC, rr.Code)
		})
	}
}
