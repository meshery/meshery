package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
)

func TestLoadTestUsingSMPHandler_BoundaryChecks(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	handler := &Handler{
		config: &models.HandlerConfig{},
		log:    log,
	}

	tests := []struct {
		name string
		body string
	}{
		{
			name: "EmptyClients",
			body: `{
				"test": {
					"name": "my-test",
					"duration": "1s",
					"clients": []
				}
			}`,
		},
		{
			name: "EmptyEndpointUrls",
			body: `{
				"test": {
					"name": "my-test",
					"duration": "1s",
					"clients": [
						{ "endpointUrls": [] }
					]
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/perf/profile", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.LoadTestUsingSMPHandler(w, req, nil, nil, nil)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400 for %s, got %d", tt.name, w.Code)
			}
		})
	}
}
