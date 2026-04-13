package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/meshery/meshery/server/internal/ai"
	"github.com/meshery/meshery/server/internal/ai/design"
	"github.com/meshery/meshery/server/internal/ai/providers"
	"github.com/meshery/meshery/server/models"
)

// ------- request / response types -------

// aiGenerateRequest is the request body for POST /api/ai/generate.
type aiGenerateRequest struct {
	Prompt   string                  `json:"prompt"`
	Provider string                  `json:"provider"`
	Model    string                  `json:"model,omitempty"`
	Options  *ai.GenerationOptions   `json:"options,omitempty"`
}

// aiGenerateResponse is the response body for POST /api/ai/generate.
type aiGenerateResponse struct {
	Content      string          `json:"content"`
	FinishReason string          `json:"finish_reason,omitempty"`
	Usage        *ai.TokenUsage  `json:"usage,omitempty"`
}

// aiStreamRequest is the request body for POST /api/ai/stream.
type aiStreamRequest struct {
	Prompt   string                `json:"prompt"`
	Provider string                `json:"provider"`
	Model    string                `json:"model,omitempty"`
	Options  *ai.GenerationOptions `json:"options,omitempty"`
}

// aiProviderInfo is returned as part of GET /api/ai/providers.
type aiProviderInfo struct {
	Name    string `json:"name"`
	Default bool   `json:"default,omitempty"`
}

// aiProvidersResponse is the response body for GET /api/ai/providers.
type aiProvidersResponse struct {
	Providers []aiProviderInfo `json:"providers"`
}

// ------- handlers -------

// HandleAIGenerate handles POST /api/ai/generate.
// It sends a single-turn completion request to the configured AI provider and
// returns the full response body once the provider has finished generating.
//
// swagger:route POST /api/ai/generate idGetAIGenerate
// Returns an AI-generated response for the supplied natural-language prompt.
// responses:
//
//	200: aiGenerateResponseWrapper
func (h *Handler) HandleAIGenerate(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	if h.aiService == nil {
		h.log.Error("AI service is not configured")
		http.Error(w, "AI service is not configured", http.StatusServiceUnavailable)
		return
	}

	var req aiGenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(ErrAIInvalidRequest(err))
		http.Error(w, ErrAIInvalidRequest(err).Error(), http.StatusBadRequest)
		return
	}

	if req.Prompt == "" || req.Provider == "" {
		err := fmt.Errorf("prompt and provider are required fields")
		h.log.Error(ErrAIInvalidRequest(err))
		http.Error(w, ErrAIInvalidRequest(err).Error(), http.StatusBadRequest)
		return
	}

	completeReq := ai.CompleteRequest{
		Context: ai.RequestContext{
			Provider: req.Provider,
			Model:    req.Model,
		},
		Messages: []ai.Message{
			{Role: ai.RoleUser, Content: req.Prompt},
		},
		Options: req.Options,
	}

	resp, err := h.aiService.Complete(r.Context(), completeReq)
	if err != nil {
		h.log.Error(ErrAIComplete(err))
		http.Error(w, ErrAIComplete(err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(aiGenerateResponse{
		Content:      resp.Content,
		FinishReason: resp.FinishReason,
		Usage:        resp.Usage,
	}); err != nil {
		h.log.Error(ErrWriteResponse(err))
	}
}

// HandleAIStream handles POST /api/ai/stream.
// It opens a streaming chat completion and forwards the tokens to the caller
// using Server-Sent Events (text/event-stream).
//
// swagger:route POST /api/ai/stream idGetAIStream
// Streams an AI-generated response for the supplied natural-language prompt.
func (h *Handler) HandleAIStream(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	if h.aiService == nil {
		h.log.Error("AI service is not configured")
		http.Error(w, "AI service is not configured", http.StatusServiceUnavailable)
		return
	}

	var req aiStreamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(ErrAIInvalidRequest(err))
		http.Error(w, ErrAIInvalidRequest(err).Error(), http.StatusBadRequest)
		return
	}

	if req.Prompt == "" || req.Provider == "" {
		err := fmt.Errorf("prompt and provider are required fields")
		h.log.Error(ErrAIInvalidRequest(err))
		http.Error(w, ErrAIInvalidRequest(err).Error(), http.StatusBadRequest)
		return
	}

	streamReq := ai.StreamRequest{
		Context: ai.RequestContext{
			Provider: req.Provider,
			Model:    req.Model,
		},
		Messages: []ai.Message{
			{Role: ai.RoleUser, Content: req.Prompt},
		},
		Options: req.Options,
	}

	reader, err := h.aiService.Stream(r.Context(), streamReq)
	if err != nil {
		h.log.Error(ErrAIStream(err))
		http.Error(w, ErrAIStream(err).Error(), http.StatusInternalServerError)
		return
	}

	// Write SSE headers.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, canFlush := w.(http.Flusher)

	for {
		chunk, err := reader.Recv()
		if err == io.EOF {
			fmt.Fprint(w, "data: [DONE]\n\n")
			if canFlush {
				flusher.Flush()
			}
			return
		}
		if err != nil {
			h.log.Error(ErrAIStream(err))
			fmt.Fprintf(w, "data: {\"error\":%q}\n\n", err.Error())
			if canFlush {
				flusher.Flush()
			}
			return
		}

		bs, _ := json.Marshal(chunk)
		fmt.Fprintf(w, "data: %s\n\n", bs)
		if canFlush {
			flusher.Flush()
		}
	}
}

// HandleAIProviders handles GET /api/ai/providers.
// It returns the list of AI provider names registered with the AI service.
//
// swagger:route GET /api/ai/providers idGetAIProviders
// Returns the list of registered AI providers.
// responses:
//
//	200: aiProvidersResponseWrapper
func (h *Handler) HandleAIProviders(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	if h.aiService == nil {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(aiProvidersResponse{Providers: []aiProviderInfo{}})
		return
	}

	names := h.aiService.ProviderNames()
	infos := make([]aiProviderInfo, 0, len(names))
	for _, name := range names {
		infos = append(infos, aiProviderInfo{Name: name})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(aiProvidersResponse{Providers: infos}); err != nil {
		h.log.Error(ErrWriteResponse(err))
	}
}

// HandleAIHealth handles GET /api/ai/health.
// It returns HTTP 200 if the AI service is configured, 503 otherwise.
//
// swagger:route GET /api/ai/health idGetAIHealth
// Returns the health status of the AI service.
func (h *Handler) HandleAIHealth(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	type healthResponse struct {
		Status    string   `json:"status"`
		Providers []string `json:"providers,omitempty"`
	}

	if h.aiService == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(healthResponse{Status: "unavailable"})
		return
	}

	names := h.aiService.ProviderNames()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(healthResponse{
		Status:    "ok",
		Providers: names,
	})
}

// HandleAIDesignGenerate handles POST /api/ai/design/generate.
// It runs the full NL → LLM → validate → re-prompt pipeline and returns the
// resulting Meshery PatternFile JSON.
//
// swagger:route POST /api/ai/design/generate idGetAIDesignGenerate
// Generates a Meshery Design from a natural-language prompt.
// responses:
//
//	200: aiDesignGenerateResponseWrapper
func (h *Handler) HandleAIDesignGenerate(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	if h.aiService == nil {
		http.Error(w, "AI service is not configured", http.StatusServiceUnavailable)
		return
	}

	var req struct {
		Prompt   string `json:"prompt"`
		Provider string `json:"provider"`
		Model    string `json:"model,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(ErrAIInvalidRequest(err))
		http.Error(w, ErrAIInvalidRequest(err).Error(), http.StatusBadRequest)
		return
	}
	if req.Prompt == "" || req.Provider == "" {
		err := fmt.Errorf("prompt and provider are required fields")
		h.log.Error(ErrAIInvalidRequest(err))
		http.Error(w, ErrAIInvalidRequest(err).Error(), http.StatusBadRequest)
		return
	}

	gen := design.NewGenerator(h.aiService, h.registryManager)
	result, err := gen.Generate(r.Context(), design.GenerateRequest{
		Prompt:   req.Prompt,
		Provider: req.Provider,
		Model:    req.Model,
	})
	if err != nil {
		h.log.Error(ErrAIComplete(err))
		http.Error(w, ErrAIComplete(err).Error(), http.StatusInternalServerError)
		return
	}
	if result.Design == nil {
		http.Error(w, "AI did not generate a valid design", http.StatusUnprocessableEntity)
		return
	}

	resp := struct {
		Design   interface{} `json:"design"`
		Warnings []string    `json:"warnings,omitempty"`
	}{
		Design:   result.Design,
		Warnings: result.Warnings,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error(ErrWriteResponse(err))
	}
}

// ------- AI service initialisation helpers -------

// InitAIService configures the AI service on the handler from a slice of
// provider configs. Providers that are not configured are silently skipped.
// Call this once after NewHandlerInstance if you want AI endpoints enabled.
func (h *Handler) InitAIService(cfgs []providers.Config) {
	if len(cfgs) == 0 {
		return
	}

	reg := ai.NewRegistry()
	for _, cfg := range cfgs {
		p, err := providers.New(cfg)
		if err != nil {
			h.log.Error(fmt.Errorf("ai: skipping provider %q: %w", cfg.Type, err))
			continue
		}
		if regErr := reg.Register(string(cfg.Type), p); regErr != nil {
			h.log.Error(fmt.Errorf("ai: failed to register provider %q: %w", cfg.Type, regErr))
		}
	}

	h.aiService = ai.NewService(reg,
		ai.LoggingMiddleware(h.log),
		ai.RetryMiddleware(ai.DefaultRetryConfig),
	)
}
