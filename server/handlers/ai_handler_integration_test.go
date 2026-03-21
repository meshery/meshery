package handlers

// Integration tests for the AI Adapter.
//
// These tests require a running Ollama instance on localhost:11434.
// They are skipped automatically when Ollama is not available, so they
// never break CI — but they run locally when a developer has Ollama up.
//
// Run manually:
//
//	ollama pull llama3
//	ollama serve
//	go test ./server/handlers/ -run Integration -v

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/utils/events"
)

// skipIfOllamaDown skips the test if Ollama is not reachable.
// This keeps CI green — integration tests only run when Ollama is available.
func skipIfOllamaDown(t *testing.T) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:11434/api/tags", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Skip("Ollama not running — skipping integration test. Start with: ollama serve")
	}
	resp.Body.Close()
}

// newIntegrationHandler builds a real Handler (no mocks) for integration tests.
func newIntegrationHandler(t *testing.T) *Handler {
	t.Helper()
	sysID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("generating uuid: %v", err)
	}
	log, err := logger.New("integration-test", logger.Options{})
	if err != nil {
		t.Fatalf("creating logger: %v", err)
	}
	return &Handler{
		log:          log,
		SystemID:     &sysID,
		EventsBuffer: events.NewEventStreamer(),
	}
}

// Integration Test 1 — Ollama end-to-end

func TestIntegration_OllamaGeneratesDesign(t *testing.T) {
	skipIfOllamaDown(t)

	h := newIntegrationHandler(t)

	body, _ := json.Marshal(AIGenerateRequest{
		Prompt:   "Deploy nginx on Kubernetes. Be brief.",
		Provider: "ollama",
		Model:    "llama3",
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/ai/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	h.AIGenerateDesignHandler(rr, req, &models.Preference{}, &models.User{}, nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp AIGenerateResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}

	// Validate provider field
	if resp.Provider != "ollama/llama3" {
		t.Errorf("provider: got %q, want %q", resp.Provider, "ollama/llama3")
	}

	// Validate description is not empty
	if resp.Description == "" {
		t.Error("expected non-empty description from Ollama")
	}

	t.Logf("   Ollama responded successfully")
	t.Logf("   Provider:    %s", resp.Provider)
	t.Logf("   Description: %.100s...", resp.Description)
	t.Logf("   Design:      %.100s...", resp.Design)
}

// Integration Test 2 — Ollama provider health check

func TestIntegration_ProvidersEndpointShowsOllamaAvailable(t *testing.T) {
	skipIfOllamaDown(t)

	h := newIntegrationHandler(t)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/ai/providers", nil)

	h.AIGetProvidersHandler(rr, req, &models.Preference{}, &models.User{}, nil)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var result struct {
		Providers []struct {
			ID        string `json:"id"`
			Available bool   `json:"available"`
		} `json:"providers"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("decoding providers response: %v", err)
	}

	for _, p := range result.Providers {
		if p.ID == "ollama" && !p.Available {
			t.Error("Ollama is running but reported as unavailable")
		}
		t.Logf("   Provider: %-12s available: %v", p.ID, p.Available)
	}
}

// Integration Test 3 — Ollama timeout handling

func TestIntegration_OllamaTimeoutHandledGracefully(t *testing.T) {
	skipIfOllamaDown(t)

	// Point to a port with nothing listening — simulates timeout
	orig := newOllamaProvider("llama3").endpoint
	_ = orig // just for documentation

	p := &ollamaProvider{
		endpoint: "http://localhost:19999", 
		model:    "llama3",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := p.generate(ctx, "test prompt")
	if err == nil {
		t.Fatal("expected error when Ollama endpoint is unreachable")
	}

	t.Logf(" Timeout handled correctly: %v", err)
}

// Integration Test 4 — Wrong method returns 405

func TestIntegration_WrongHTTPMethod(t *testing.T) {
	h := newIntegrationHandler(t)

	// GET on /api/ai/generate should not be handled by AIGenerateDesignHandler
	// (router enforces POST only, but let's verify handler is robust)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/ai/generate", nil)

	h.AIGenerateDesignHandler(rr, req, &models.Preference{}, &models.User{}, nil)

	// Handler will try to decode empty body → 400 bad request
	// This is acceptable — router level enforces 405 before reaching handler
	if rr.Code == http.StatusOK {
		t.Error("GET with no body should not return 200")
	}

	t.Logf(" GET request correctly rejected with status %d", rr.Code)
}