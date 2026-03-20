package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func TestEnvironmentAndWorkspaceHandlers_InvalidIDReturnsBadRequest(t *testing.T) {
	h := newTestHandler(t, nil, "")
	provider := &models.DefaultLocalProvider{}

	tests := []struct {
		name          string
		method        string
		path          string
		id            string
		invokeHandler func(*httptest.ResponseRecorder, *http.Request)
	}{
		{
			name:   "get environment by id with malformed id",
			method: http.MethodGet,
			path:   "/api/environments/abc",
			id:     "abc",
			invokeHandler: func(rr *httptest.ResponseRecorder, req *http.Request) {
				h.GetEnvironmentByIDHandler(rr, req, nil, nil, provider)
			},
		},
		{
			name:   "delete environment with malformed id",
			method: http.MethodDelete,
			path:   "/api/environments/not-a-uuid",
			id:     "not-a-uuid",
			invokeHandler: func(rr *httptest.ResponseRecorder, req *http.Request) {
				h.DeleteEnvironmentHandler(rr, req, nil, nil, provider)
			},
		},
		{
			name:   "update environment with empty id",
			method: http.MethodPut,
			path:   "/api/environments/",
			id:     "",
			invokeHandler: func(rr *httptest.ResponseRecorder, req *http.Request) {
				h.UpdateEnvironmentHandler(rr, req, nil, nil, provider)
			},
		},
		{
			name:   "get workspace by id with malformed id",
			method: http.MethodGet,
			path:   "/api/workspaces/abc",
			id:     "abc",
			invokeHandler: func(rr *httptest.ResponseRecorder, req *http.Request) {
				h.GetWorkspaceByIdHandler(rr, req, nil, nil, provider)
			},
		},
		{
			name:   "delete workspace with malformed id",
			method: http.MethodDelete,
			path:   "/api/workspaces/not-a-uuid",
			id:     "not-a-uuid",
			invokeHandler: func(rr *httptest.ResponseRecorder, req *http.Request) {
				h.DeleteWorkspaceHandler(rr, req, nil, nil, provider)
			},
		},
		{
			name:   "update workspace with empty id",
			method: http.MethodPut,
			path:   "/api/workspaces/",
			id:     "",
			invokeHandler: func(rr *httptest.ResponseRecorder, req *http.Request) {
				h.UpdateWorkspaceHandler(rr, req, nil, nil, provider)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			rr := httptest.NewRecorder()

			tt.invokeHandler(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
			}
		})
	}
}
