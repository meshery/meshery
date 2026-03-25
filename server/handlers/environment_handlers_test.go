package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
)

type environmentInvoker func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider)

func TestEnvironmentHandlers_InvalidUUID(t *testing.T) {
	h := newTestHandler(t, nil, "")
	provider := &models.DefaultLocalProvider{}

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		urlVars    map[string]string
		withToken  bool
		invoker    environmentInvoker
		expectedSC int
	}{
		{
			name:      "given invalid org id when GetEnvironments return status code 400",
			method:    http.MethodGet,
			path:      "/api/environments?orgID=not-a-uuid",
			withToken: true,
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.GetEnvironments(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid environment id when GetEnvironmentByIDHandler return status code 400",
			method:  http.MethodGet,
			path:    "/api/environments/not-a-uuid",
			urlVars: map[string]string{"id": "not-a-uuid"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.GetEnvironmentByIDHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid org id when GetEnvironmentByIDHandler return status code 400",
			method:  http.MethodGet,
			path:    "/api/environments/123e4567-e89b-12d3-a456-426614174000?orgID=not-a-uuid",
			urlVars: map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.GetEnvironmentByIDHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:   "given invalid organization id when SaveEnvironment return status code 400",
			method: http.MethodPost,
			path:   "/api/environments",
			body:   `{"name":"env-1","organization_id":"not-a-uuid"}`,
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.SaveEnvironment(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid environment id when UpdateEnvironmentHandler return status code 400",
			method:  http.MethodPut,
			path:    "/api/environments/not-a-uuid",
			body:    `{"name":"env-1","organization_id":"123e4567-e89b-12d3-a456-426614174000"}`,
			urlVars: map[string]string{"id": "not-a-uuid"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.UpdateEnvironmentHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid organization id when UpdateEnvironmentHandler return status code 400",
			method:  http.MethodPut,
			path:    "/api/environments/123e4567-e89b-12d3-a456-426614174000",
			body:    `{"name":"env-1","organization_id":"not-a-uuid"}`,
			urlVars: map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.UpdateEnvironmentHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid environment id when DeleteEnvironmentHandler return status code 400",
			method:  http.MethodDelete,
			path:    "/api/environments/not-a-uuid",
			urlVars: map[string]string{"id": "not-a-uuid"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.DeleteEnvironmentHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid connection id when AddConnectionToEnvironmentHandler return status code 400",
			method:  http.MethodPost,
			path:    "/api/environments/123e4567-e89b-12d3-a456-426614174000/connections/not-a-uuid",
			urlVars: map[string]string{"environmentID": "123e4567-e89b-12d3-a456-426614174000", "connectionID": "not-a-uuid"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.AddConnectionToEnvironmentHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid environment id when GetConnectionsOfEnvironmentHandler return status code 400",
			method:  http.MethodGet,
			path:    "/api/environments/not-a-uuid/connections",
			urlVars: map[string]string{"environmentID": "not-a-uuid"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.GetConnectionsOfEnvironmentHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			if tt.withToken {
				req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "token"))
			}
			if tt.urlVars != nil {
				req = mux.SetURLVars(req, tt.urlVars)
			}

			rr := httptest.NewRecorder()
			tt.invoker(rr, req, provider)

			assert.Equal(t, tt.expectedSC, rr.Code)
		})
	}
}
