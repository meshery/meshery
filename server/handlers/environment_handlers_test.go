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

type invokerHandler func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider)

func TestEnvironmentHandlers_InvalidUUID(t *testing.T) {
	h := newTestHandler(t, nil, "")
	provider := &models.DefaultLocalProvider{}
	invalidID := "not-a-uuid"

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		invoker    invokerHandler
		expectedSC int
	}{
		{
			name:   "given non valid uuid when GetEnvironmentByIDHandler return status code 400",
			method: http.MethodGet,
			path:   "/api/environments/" + invalidID,
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.GetEnvironmentByIDHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:   "given non valid uuid when UpdateEnvironmentHandler return status code 400",
			method: http.MethodPut,
			path:   "/api/environments/" + invalidID,
			body:   `{"name":"env-1"}`,
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.UpdateEnvironmentHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
		{
			name:   "given non valid uuid when DeleteEnvironmentHandler return status code 400",
			method: http.MethodDelete,
			path:   "/api/environments/" + invalidID,
			invoker: func(rr *httptest.ResponseRecorder, req *http.Request, provider *models.DefaultLocalProvider) {
				h.DeleteEnvironmentHandler(rr, req, nil, nil, provider)
			},
			expectedSC: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader *strings.Reader
			if tt.body != "" {
				bodyReader = strings.NewReader(tt.body)
			} else {
				bodyReader = strings.NewReader("")
			}

			req := httptest.NewRequest(tt.method, tt.path, bodyReader)
			req = mux.SetURLVars(req, map[string]string{"id": invalidID})
			rr := httptest.NewRecorder()

			tt.invoker(rr, req, provider)
			assert.Equal(t, tt.expectedSC, rr.Code)
		})
	}
}
