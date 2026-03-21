package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
)

// ─────────────────────────────────────────────────────────────────────────────
// Request / Response types
// ─────────────────────────────────────────────────────────────────────────────

// AIGenerateRequest is the payload sent to POST /api/ai/generate.
type AIGenerateRequest struct {
	Prompt   string `json:"prompt"`   // e.g. "HA k8s cluster with Prometheus"
	Provider string `json:"provider"` // "ollama" | "openai" | "anthropic" (default: "ollama")
	Model    string `json:"model"`    // optional model override e.g. "llama3", "gpt-4o"
}

// AIGenerateResponse is returned on success.
type AIGenerateResponse struct {
	Provider    string `json:"provider"`    // which LLM was used
	Model       string `json:"model"`       // which model was used
	Description string `json:"description"` // human-readable summary from the LLM
	Design      string `json:"design"`      // Meshery design / YAML manifest
	PatternID   string `json:"pattern_id,omitempty"` // saved pattern ID — open in canvas at /api/pattern/{id}
}

// ─────────────────────────────────────────────────────────────────────────────
// LLM provider abstraction  (BYOM — Bring Your Own Model)
// ─────────────────────────────────────────────────────────────────────────────

// llmProvider is the single interface every backend must satisfy.
// Adding a new LLM means implementing two methods — nothing else changes.
type llmProvider interface {
	generate(ctx context.Context, prompt string) (string, error)
	name() string
}

// ── Ollama (local, private, no API key needed) ────────────────────────────────

type ollamaProvider struct {
	endpoint string
	model    string
}

func newOllamaProvider(model string) *ollamaProvider {
	if model == "" {
		model = "llama3"
	}
	return &ollamaProvider{endpoint: "http://localhost:11434", model: model}
}

func (o *ollamaProvider) name() string { return fmt.Sprintf("ollama/%s", o.model) }

func (o *ollamaProvider) generate(ctx context.Context, prompt string) (string, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"model":  o.model,
		"prompt": prompt,
		"stream": false,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.endpoint+"/api/generate", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("ollama: building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama: request failed (is Ollama running?): %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ollama: reading response: %w", err)
	}
	var result struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("ollama: parsing response: %w", err)
	}
	return result.Response, nil
}

// ── OpenAI ────────────────────────────────────────────────────────────────────

type openAIProvider struct {
	apiKey string
	model  string
}

func newOpenAIProvider(apiKey, model string) *openAIProvider {
	if model == "" {
		model = "gpt-4o"
	}
	return &openAIProvider{apiKey: apiKey, model: model}
}

func (o *openAIProvider) name() string { return fmt.Sprintf("openai/%s", o.model) }

func (o *openAIProvider) generate(ctx context.Context, prompt string) (string, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"model": o.model,
		"messages": []map[string]string{
			{"role": "system", "content": mesherySystemPrompt},
			{"role": "user", "content": prompt},
		},
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("openai: building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai: request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("openai: reading response: %w", err)
	}
	var result struct {
		Choices []struct {
			Message struct{ Content string `json:"content"` } `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(raw, &result); err != nil || len(result.Choices) == 0 {
		return "", fmt.Errorf("openai: unexpected response format")
	}
	return result.Choices[0].Message.Content, nil
}

// ── Anthropic ─────────────────────────────────────────────────────────────────

type anthropicProvider struct {
	apiKey string
	model  string
}

func newAnthropicProvider(apiKey, model string) *anthropicProvider {
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}
	return &anthropicProvider{apiKey: apiKey, model: model}
}

func (a *anthropicProvider) name() string { return fmt.Sprintf("anthropic/%s", a.model) }

func (a *anthropicProvider) generate(ctx context.Context, prompt string) (string, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"model":      a.model,
		"max_tokens": 2048,
		"system":     mesherySystemPrompt,
		"messages":   []map[string]string{{"role": "user", "content": prompt}},
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.anthropic.com/v1/messages", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("anthropic: building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic: request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("anthropic: reading response: %w", err)
	}
	var result struct {
		Content []struct{ Text string `json:"text"` } `json:"content"`
	}
	if err := json.Unmarshal(raw, &result); err != nil || len(result.Content) == 0 {
		return "", fmt.Errorf("anthropic: unexpected response format")
	}
	return result.Content[0].Text, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// System prompt  (context-window management — referenced in issue #17097)
// ─────────────────────────────────────────────────────────────────────────────

const mesherySystemPrompt = `You are Meshery's AI assistant for cloud-native infrastructure.

Given a natural-language description, respond with exactly:
DESCRIPTION: <2-3 sentence summary>
DESIGN:
<valid Kubernetes YAML manifest using Meshery schema v0.9.0>

Always include apiVersion, kind, metadata, and spec. No extra prose.`

// ─────────────────────────────────────────────────────────────────────────────
// Provider factory
// ─────────────────────────────────────────────────────────────────────────────

// providerFactory is a package-level var so tests can swap it out without
// spinning up real network services.
var providerFactory = providerFromRequest

// providerFromRequest selects and initialises the correct llmProvider.
// Credentials come from environment variables — never from the request body.
func providerFromRequest(req AIGenerateRequest) (llmProvider, error) {
	switch strings.ToLower(req.Provider) {
	case "openai":
		key := envLookup("OPENAI_API_KEY")
		if key == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
		}
		return newOpenAIProvider(key, req.Model), nil
	case "anthropic":
		key := envLookup("ANTHROPIC_API_KEY")
		if key == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
		}
		return newAnthropicProvider(key, req.Model), nil
	default: // "ollama" or empty → local, no key needed
		return newOllamaProvider(req.Model), nil
	}
}

// envLookup is a var so tests can override it without touching os.
var envLookup = func(key string) string { return os.Getenv(key) }

// ─────────────────────────────────────────────────────────────────────────────
// HTTP handlers
// ─────────────────────────────────────────────────────────────────────────────

// AIGenerateDesignHandler handles POST /api/ai/generate.
//
// Accepts a natural-language prompt, forwards it to the requested LLM
// provider, and returns a Meshery-compatible design + YAML manifest.
//
// swagger:route POST /api/ai/generate idAIGenerateDesign
// responses:
//
//	200: AIGenerateResponse
//	400: description: bad request
//	500: description: internal server error
func (h *Handler) AIGenerateDesignHandler(
	w http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	// ── parse request ─────────────────────────────────────────────────────────
	var req AIGenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(fmt.Errorf("ai handler: decoding request body: %w", err))
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Prompt) == "" {
		http.Error(w, `{"error":"prompt is required"}`, http.StatusBadRequest)
		return
	}

	// ── pick provider ─────────────────────────────────────────────────────────
	llm, err := providerFactory(req)
	if err != nil {
		h.log.Error(fmt.Errorf("ai handler: selecting provider: %w", err))
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	h.log.Info(fmt.Sprintf("AI design request — provider: %s  prompt: %.80s…", llm.name(), req.Prompt))

	// ── publish event so Meshery UI can show a notification ───────────────────
	eventBuilder := events.NewEvent().
		ActedUpon(*h.SystemID).
		FromUser(user.ID).
		WithSeverity(events.Informational).
		WithDescription(fmt.Sprintf("Generating infrastructure design via %s", llm.name())).
		WithMetadata(map[string]interface{}{"prompt": req.Prompt})
	h.EventsBuffer.Publish(eventBuilder.Build())

	// ── call LLM with a 60-second timeout ─────────────────────────────────────
	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()

	raw, err := llm.generate(ctx, req.Prompt)
	if err != nil {
		h.log.Error(fmt.Errorf("ai handler: generating design: %w", err))
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// ── parse LLM output ──────────────────────────────────────────────────────
	description, design := parseAIResponse(raw)

	// ── save as Meshery pattern so it renders on the canvas ───────────────────
	// Satisfies acceptance test: "natural language queries result in a rendered design"
	// Only attempt to save if a real provider is available (skipped in unit tests)
	var patternID string
	if provider != nil {
		patternID = saveAsPattern(h, r, provider, user, req.Prompt, design)
	}

	// ── publish success event ─────────────────────────────────────────────────
	successEvent := events.NewEvent().
		ActedUpon(*h.SystemID).
		FromUser(user.ID).
		WithSeverity(events.Informational).
		WithDescription(fmt.Sprintf("AI design generated and saved via %s", llm.name())).
		WithMetadata(map[string]interface{}{
			"prompt":     req.Prompt,
			"pattern_id": patternID,
		})
	h.EventsBuffer.Publish(successEvent.Build())

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(AIGenerateResponse{
		Provider:    llm.name(),
		Model:       req.Model,
		Description: description,
		Design:      design,
		PatternID:   patternID,
	}); err != nil {
		h.log.Error(fmt.Errorf("ai handler: writing response: %w", err))
	}
}

// AIGetProvidersHandler handles GET /api/ai/providers.
// Returns the list of available LLM providers so the UI can populate a picker.
func (h *Handler) AIGetProvidersHandler(
	w http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	type providerInfo struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
		Local       bool   `json:"local"`
		Available   bool   `json:"available"`
	}

	providers := []providerInfo{
		{ID: "ollama", DisplayName: "Ollama (local)", Local: true, Available: isOllamaRunning()},
		{ID: "openai", DisplayName: "OpenAI", Local: false, Available: envLookup("OPENAI_API_KEY") != ""},
		{ID: "anthropic", DisplayName: "Anthropic Claude", Local: false, Available: envLookup("ANTHROPIC_API_KEY") != ""},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"providers": providers})
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// parseAIResponse splits the LLM output into a human-readable description
// and the YAML design block.
func parseAIResponse(raw string) (description, design string) {
	parts := strings.SplitN(raw, "DESIGN:", 2)
	if len(parts) == 2 {
		desc := strings.TrimPrefix(parts[0], "DESCRIPTION:")
		return strings.TrimSpace(desc), strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(raw), ""
}

// isOllamaRunning does a lightweight health-check against the local Ollama API.
func isOllamaRunning() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:11434/api/tags", nil)
	if err != nil {
		return false
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// saveAsPattern saves the AI-generated YAML design as a Meshery pattern so it
// appears on the visual design canvas. Returns the pattern ID on success, or
// an empty string if saving fails (non-fatal — the design is still returned).
func saveAsPattern(
	h *Handler,
	r *http.Request,
	provider models.Provider,
	user *models.User,
	prompt, yamlDesign string,
) string {
	// Guard against nil provider or empty design — non-fatal, just skip saving
	if provider == nil || strings.TrimSpace(yamlDesign) == "" {
		return ""
	}

	// Derive a friendly name from the first 60 chars of the prompt
	name := prompt
	if len(name) > 60 {
		name = name[:60] + "..."
	}
	name = "AI: " + name

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.log.Warn(fmt.Errorf("ai handler: getting provider token for pattern save: %w", err))
		return ""
	}

	pattern := &models.MesheryPattern{
		Name:        name,
		PatternFile: yamlDesign,
	}

	savedBytes, err := provider.SaveMesheryPattern(token, pattern)
	if err != nil {
		h.log.Warn(fmt.Errorf("ai handler: saving pattern to canvas: %w", err))
		return ""
	}

	// SaveMesheryPattern returns the saved pattern as JSON — extract the ID
	var saved []models.MesheryPattern
	if err := json.Unmarshal(savedBytes, &saved); err != nil || len(saved) == 0 {
		return ""
	}
	if saved[0].ID == nil {
		return ""
	}
	return saved[0].ID.String()
}