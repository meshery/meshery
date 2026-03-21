package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/utils/events"
)

// ─────────────────────────────────────────────────────────────────────────────
// Mock — no real network calls in tests
// ─────────────────────────────────────────────────────────────────────────────

type mockLLMProvider struct {
	resp string
	err  error
}

func (m *mockLLMProvider) generate(_ context.Context, _ string) (string, error) {
	return m.resp, m.err
}
func (m *mockLLMProvider) name() string { return "mock/test-model" }

// newAITestHandler builds the minimal *Handler needed by the AI handler tests.
func newAITestHandler(t *testing.T) *Handler {
	t.Helper()
	sysID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("generating uuid: %v", err)
	}
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("creating logger: %v", err)
	}
	return &Handler{
		log:          log,
		SystemID:     &sysID,
		EventsBuffer: events.NewEventStreamer(),
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Unit tests
// ─────────────────────────────────────────────────────────────────────────────

func TestParseAIResponse(t *testing.T) {
	tests := []struct {
		name       string
		raw        string
		wantDesc   string
		wantDesign string
	}{
		{
			name:       "well-formed response",
			raw:        "DESCRIPTION: A simple nginx deployment.\nDESIGN:\napiVersion: apps/v1\nkind: Deployment",
			wantDesc:   "A simple nginx deployment.",
			wantDesign: "apiVersion: apps/v1\nkind: Deployment",
		},
		{
			name:       "missing DESIGN tag falls back gracefully",
			raw:        "Just some text from the LLM",
			wantDesc:   "Just some text from the LLM",
			wantDesign: "",
		},
		{
			name:       "empty response",
			raw:        "",
			wantDesc:   "",
			wantDesign: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			desc, design := parseAIResponse(tc.raw)
			if desc != tc.wantDesc {
				t.Errorf("description: got %q, want %q", desc, tc.wantDesc)
			}
			if design != tc.wantDesign {
				t.Errorf("design: got %q, want %q", design, tc.wantDesign)
			}
		})
	}
}

func TestProviderFromRequest_DefaultsToOllama(t *testing.T) {
	p, err := providerFromRequest(AIGenerateRequest{Prompt: "test", Provider: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.name() != "ollama/llama3" {
		t.Errorf("expected ollama/llama3, got %s", p.name())
	}
}

func TestProviderFromRequest_OpenAIRequiresKey(t *testing.T) {
	orig := envLookup
	envLookup = func(_ string) string { return "" }
	defer func() { envLookup = orig }()

	_, err := providerFromRequest(AIGenerateRequest{Provider: "openai"})
	if err == nil {
		t.Fatal("expected error when OPENAI_API_KEY is missing")
	}
}

func TestProviderFromRequest_AnthropicRequiresKey(t *testing.T) {
	orig := envLookup
	envLookup = func(_ string) string { return "" }
	defer func() { envLookup = orig }()

	_, err := providerFromRequest(AIGenerateRequest{Provider: "anthropic"})
	if err == nil {
		t.Fatal("expected error when ANTHROPIC_API_KEY is missing")
	}
}

func TestAIGenerateDesignHandler_Success(t *testing.T) {
	h := newAITestHandler(t)

	// Inject mock provider — no network needed
	orig := providerFactory
	providerFactory = func(_ AIGenerateRequest) (llmProvider, error) {
		return &mockLLMProvider{
			resp: "DESCRIPTION: Nginx deployment on k8s.\nDESIGN:\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: nginx",
		}, nil
	}
	defer func() { providerFactory = orig }()

	body, _ := json.Marshal(AIGenerateRequest{Prompt: "Deploy nginx on Kubernetes"})
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
	if resp.Description == "" {
		t.Error("expected non-empty description")
	}
	if resp.Design == "" {
		t.Error("expected non-empty design/yaml")
	}
}

func TestAIGenerateDesignHandler_EmptyPrompt(t *testing.T) {
	h := newAITestHandler(t)

	body, _ := json.Marshal(AIGenerateRequest{Prompt: "   "})
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/ai/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	h.AIGenerateDesignHandler(rr, req, &models.Preference{}, &models.User{}, nil)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty prompt, got %d", rr.Code)
	}
}