package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	SMP "github.com/layer5io/service-mesh-performance/spec"
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

	payload := models.PerformanceTestConfigFile{
		Config: &SMP.PerformanceTestConfig{
			Id:       "profile-1",
			Name:     "load-test-1",
			Duration: "1s",
			Clients:  []*SMP.PerformanceTestConfig_Client{},
		},
		ServiceMesh: &SMP.ServiceMesh{
			Type: SMP.ServiceMesh_ISTIO,
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/load-test/smp", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		handler.LoadTestUsingSMPHandler(rec, req, nil, nil, nil)
	}()

	if didPanic {
		t.Fatalf("expected no panic")
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d (body=%q)", http.StatusBadRequest, rec.Code, rec.Body.String())
	}

	var decoded map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
		t.Fatalf("response body should be JSON: %v (body=%q)", err, rec.Body.String())
	}

	wantErr := "no clients provided in the SMP configuration"
	if decoded["error"] != wantErr {
		t.Fatalf("expected error %q, got %q", wantErr, decoded["error"])
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

	payload := models.PerformanceTestConfigFile{
		Config: &SMP.PerformanceTestConfig{
			Id:       "profile-1",
			Name:     "load-test-1",
			Duration: "1s",
			Clients: []*SMP.PerformanceTestConfig_Client{
				{
					EndpointUrls: []string{},
				},
			},
		},
		ServiceMesh: &SMP.ServiceMesh{
			Type: SMP.ServiceMesh_ISTIO,
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/load-test/smp", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		handler.LoadTestUsingSMPHandler(rec, req, nil, nil, nil)
	}()

	if didPanic {
		t.Fatalf("expected no panic")
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d (body=%q)", http.StatusBadRequest, rec.Code, rec.Body.String())
	}

	var decoded map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
		t.Fatalf("response body should be JSON: %v (body=%q)", err, rec.Body.String())
	}

	wantErr := "no endpoint URLs provided in the SMP client configuration"
	if decoded["error"] != wantErr {
		t.Fatalf("expected error %q, got %q", wantErr, decoded["error"])
	}
}

