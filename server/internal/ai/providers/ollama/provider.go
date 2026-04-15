package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/meshery/meshery/server/internal/ai"
)

const defaultBaseURL = "http://localhost:11434"

// Config holds connection parameters for an Ollama instance.
type Config struct {
	BaseURL string // defaults to http://localhost:11434
	Model   string // e.g. "llama3"
}

type provider struct {
	baseURL string
	model   string
	client  *http.Client
}

// New creates an Ollama provider. The extended timeout accounts for
// cold-start model loading which can take 10-30s on first request.
func New(cfg Config) *provider {
	base := cfg.BaseURL
	if base == "" {
		base = defaultBaseURL
	}
	return &provider{
		baseURL: base,
		model:   cfg.Model,
		client:  &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *provider) Name() string { return "ollama" }

func (p *provider) resolveModel(reqModel string) (string, error) {
	if reqModel != "" {
		return reqModel, nil
	}
	if p.model != "" {
		return p.model, nil
	}
	return "", fmt.Errorf("ollama: no model specified in request or provider config")
}

// Complete sends a chat completion request with streaming disabled.
func (p *provider) Complete(ctx context.Context, req ai.CompleteRequest) (ai.CompleteResponse, error) {
	model, err := p.resolveModel(req.Model)
	if err != nil {
		return ai.CompleteResponse{}, err
	}

	payload := chatRequest{
		Model:    model,
		Messages: toOllamaMessages(req.Messages),
		Stream:   boolPtr(false),
	}
	opts := buildOptions(req)
	if opts != nil {
		payload.Options = opts
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return ai.CompleteResponse{}, fmt.Errorf("ollama: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return ai.CompleteResponse{}, fmt.Errorf("ollama: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return ai.CompleteResponse{}, fmt.Errorf("ollama: request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return ai.CompleteResponse{}, fmt.Errorf("ollama: %d — %s", resp.StatusCode, string(b))
	}

	var cr chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return ai.CompleteResponse{}, fmt.Errorf("ollama: decode response: %w", err)
	}

	return ai.CompleteResponse{
		Content:      cr.Message.Content,
		Model:        cr.Model,
		FinishReason: finishReason(cr.DoneReason),
		Usage: ai.TokenUsage{
			PromptTokens:     cr.PromptEvalCount,
			CompletionTokens: cr.EvalCount,
			TotalTokens:      cr.PromptEvalCount + cr.EvalCount,
		},
	}, nil
}

// Stream returns the raw response body for NDJSON streaming.
// Each line is a JSON object with a "done" field. The caller
// is responsible for closing the returned ReadCloser.
func (p *provider) Stream(ctx context.Context, req ai.CompleteRequest) (io.ReadCloser, error) {
	model, err := p.resolveModel(req.Model)
	if err != nil {
		return nil, err
	}

	payload := chatRequest{
		Model:    model,
		Messages: toOllamaMessages(req.Messages),
		Stream:   boolPtr(true),
	}
	opts := buildOptions(req)
	if opts != nil {
		payload.Options = opts
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("ollama: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ollama: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama: stream request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		_ = resp.Body.Close()
		return nil, fmt.Errorf("ollama: %d — %s", resp.StatusCode, string(b))
	}
	return resp.Body, nil
}

// Embed generates a vector embedding via Ollama's /api/embeddings endpoint.
func (p *provider) Embed(ctx context.Context, req ai.EmbedRequest) (ai.EmbedResponse, error) {
	model, err := p.resolveModel(req.Model)
	if err != nil {
		return ai.EmbedResponse{}, err
	}

	payload := embedRequest{Model: model, Prompt: req.Text}
	body, err := json.Marshal(payload)
	if err != nil {
		return ai.EmbedResponse{}, fmt.Errorf("ollama: marshal embed request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/api/embeddings", bytes.NewReader(body))
	if err != nil {
		return ai.EmbedResponse{}, fmt.Errorf("ollama: build embed request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return ai.EmbedResponse{}, fmt.Errorf("ollama: embed request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return ai.EmbedResponse{}, fmt.Errorf("ollama: embed %d — %s", resp.StatusCode, string(b))
	}

	var er embedResponse
	if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
		return ai.EmbedResponse{}, fmt.Errorf("ollama: decode embed response: %w", err)
	}

	return ai.EmbedResponse{Embedding: er.Embedding}, nil
}

// Models queries /api/tags and returns the names of locally available models.
func (p *provider) Models(ctx context.Context) ([]string, error) {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, p.baseURL+"/api/tags", nil)
	if err != nil {
		return nil, fmt.Errorf("ollama: build tags request: %w", err)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama: tags request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("ollama: tags %d — %s", resp.StatusCode, string(b))
	}

	var tr tagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, fmt.Errorf("ollama: decode tags: %w", err)
	}

	names := make([]string, len(tr.Models))
	for i, m := range tr.Models {
		names[i] = m.Name
	}
	return names, nil
}

// --- wire types matching the Ollama HTTP API ---

type chatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   *bool           `json:"stream,omitempty"`
	Options  *ollamaOptions  `json:"options,omitempty"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaOptions struct {
	NumPredict int     `json:"num_predict,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

type chatResponse struct {
	Model           string        `json:"model"`
	Message         ollamaMessage `json:"message"`
	Done            bool          `json:"done"`
	DoneReason      string        `json:"done_reason"`
	PromptEvalCount int           `json:"prompt_eval_count"`
	EvalCount       int           `json:"eval_count"`
}

type embedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type embedResponse struct {
	Embedding []float64 `json:"embedding"`
}

type tagsResponse struct {
	Models []struct {
		Name string `json:"name"`
	} `json:"models"`
}

func buildOptions(req ai.CompleteRequest) *ollamaOptions {
	var opts ollamaOptions
	set := false
	if req.MaxTokens != nil {
		opts.NumPredict = *req.MaxTokens
		set = true
	}
	if req.Temperature != nil {
		opts.Temperature = *req.Temperature
		set = true
	}
	if !set {
		return nil
	}
	return &opts
}

func toOllamaMessages(msgs []ai.Message) []ollamaMessage {
	out := make([]ollamaMessage, len(msgs))
	for i, m := range msgs {
		out[i] = ollamaMessage{Role: m.Role, Content: m.Content}
	}
	return out
}

func finishReason(reason string) string {
	if reason == "" {
		return "stop"
	}
	return reason
}

func boolPtr(v bool) *bool { return &v }
