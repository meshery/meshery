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

func TestUpdateWorkspaceHandler(t *testing.T) {
	h := newTestHandler(t, nil, "")
	provider := &models.DefaultLocalProvider{}

	tests := []struct {
		name       string
		id         string
		payload    string
		expectedSC int
	}{
		{
			name:       "given non valid uuid when UpdateWorkspaceHandler return status code 400",
			id:         "not-a-uuid",
			payload:    `{"name":"workspace-1"}`,
			expectedSC: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/workspaces/"+tt.id, strings.NewReader(tt.payload))
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			rr := httptest.NewRecorder()

			h.UpdateWorkspaceHandler(rr, req, nil, nil, provider)

			assert.Equal(t, tt.expectedSC, rr.Code)
		})
	}
}
