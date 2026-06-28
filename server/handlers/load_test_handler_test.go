package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
)

func TestLoadTestUsingSMPHandler_EmptyClients(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	handler := &Handler{
		config: &models.HandlerConfig{},
		log:    log,
	}

	body := strings.NewReader(`{
		"test": {
			"name": "my-test",
			"duration": "1s",
			"clients": []
		}
	}`)

	req := httptest.NewRequest(http.MethodPost, "/perf/profile", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.LoadTestUsingSMPHandler(w, req, nil, nil, nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty clients, got %d", w.Code)
	}
}

func TestLoadTestUsingSMPHandler_EmptyEndpointUrls(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	handler := &Handler{
		config: &models.HandlerConfig{},
		log:    log,
	}

	body := strings.NewReader(`{
		"test": {
			"name": "my-test",
			"duration": "1s",
			"clients": [
				{ "endpointUrls": [] }
			]
		}
	}`)

	req := httptest.NewRequest(http.MethodPost, "/perf/profile", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.LoadTestUsingSMPHandler(w, req, nil, nil, nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty endpointUrls, got %d", w.Code)
	}
}
