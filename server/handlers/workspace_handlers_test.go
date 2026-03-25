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

type workspaceInvoker func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider)

func TestWorkspaceHandlers_InvalidUUID(t *testing.T) {
	h := newTestHandler(t, nil, "")
	provider := &models.DefaultLocalProvider{}

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		urlVars    map[string]string
		withToken  bool
		invoker    workspaceInvoker
		expectedSC int
	}{
		{
			name:      "given invalid org id when GetWorkspacesHandler return status code 400",
			method:    http.MethodGet,
			path:      "/api/workspaces?orgID=not-a-uuid",
			withToken: true,
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.GetWorkspacesHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid workspace id when GetWorkspaceByIdHandler return status code 400",
			method:  http.MethodGet,
			path:    "/api/workspaces/not-a-uuid?orgID=123e4567-e89b-12d3-a456-426614174000",
			urlVars: map[string]string{"id": "not-a-uuid"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.GetWorkspaceByIdHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given missing org id when GetWorkspaceByIdHandler return status code 400",
			method:  http.MethodGet,
			path:    "/api/workspaces/123e4567-e89b-12d3-a456-426614174000",
			urlVars: map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.GetWorkspaceByIdHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid org id when GetWorkspaceByIdHandler return status code 400",
			method:  http.MethodGet,
			path:    "/api/workspaces/123e4567-e89b-12d3-a456-426614174000?orgID=not-a-uuid",
			urlVars: map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.GetWorkspaceByIdHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:   "given missing organization id when SaveWorkspaceHandler return status code 400",
			method: http.MethodPost,
			path:   "/api/workspaces",
			body:   `{"name":"workspace-1"}`,
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.SaveWorkspaceHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid workspace id when UpdateWorkspaceHandler return status code 400",
			method:  http.MethodPut,
			path:    "/api/workspaces/not-a-uuid",
			body:    `{"name":"workspace-1","organization_id":"123e4567-e89b-12d3-a456-426614174000"}`,
			urlVars: map[string]string{"id": "not-a-uuid"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.UpdateWorkspaceHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given missing organization id when UpdateWorkspaceHandler return status code 400",
			method:  http.MethodPut,
			path:    "/api/workspaces/123e4567-e89b-12d3-a456-426614174000",
			body:    `{"name":"workspace-1"}`,
			urlVars: map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.UpdateWorkspaceHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid workspace id when DeleteWorkspaceHandler return status code 400",
			method:  http.MethodDelete,
			path:    "/api/workspaces/not-a-uuid",
			urlVars: map[string]string{"id": "not-a-uuid"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.DeleteWorkspaceHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid workspace id when GetEnvironmentsOfWorkspaceHandler return status code 400",
			method:  http.MethodGet,
			path:    "/api/workspaces/not-a-uuid/environments",
			urlVars: map[string]string{"id": "not-a-uuid"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.GetEnvironmentsOfWorkspaceHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid environment id when AddEnvironmentToWorkspaceHandler return status code 400",
			method:  http.MethodPost,
			path:    "/api/workspaces/123e4567-e89b-12d3-a456-426614174000/environments/not-a-uuid",
			urlVars: map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000", "environmentID": "not-a-uuid"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.AddEnvironmentToWorkspaceHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:    "given invalid design id when AddDesignToWorkspaceHandler return status code 400",
			method:  http.MethodPost,
			path:    "/api/workspaces/123e4567-e89b-12d3-a456-426614174000/designs/not-a-uuid",
			urlVars: map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000", "designID": "not-a-uuid"},
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.AddDesignToWorkspaceHandler(rr, req, nil, nil, provider)
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
