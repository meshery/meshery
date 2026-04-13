// Package anthropic provides an LLMProvider implementation for the Anthropic Messages API.
package anthropic

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/meshery/meshery/server/internal/ai"
)

const (
	defaultBaseURL  = "https://api.anthropic.com"
	defaultModel    = "claude-3-5-sonnet-20241022"
	defaultTimeout  = 60 * time.Second
	messagesPath    = "/v1/messages"
	anthropicVersion = "2023-06-01"
)

// Provider implements ai.LLMProvider for Anthropic.
type Provider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// Option configures an Anthropic Provider.
type Option func(*Provider)

// WithBaseURL overrides the default Anthropic API base URL.
func WithBaseURL(url string) Option {
	return func(p *Provider) {
		p.baseURL = strings.TrimRight(url, "/")
	}
}

// WithModel overrides the default chat model.
func WithModel(model string) Option {
	return func(p *Provider) {
		p.model = model
	}
}

// WithHTTPClient replaces the default HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(p *Provider) {
		p.client = c
	}
}

// New creates an Anthropic provider. apiKey is required.
func New(apiKey string, opts ...Option) *Provider {
	p := &Provider{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		model:   defaultModel,
		client:  &http.Client{Timeout: defaultTimeout},
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

// ------- Anthropic API types -------

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	Messages  []anthropicMessage `json:"messages"`
	System    string             `json:"system,omitempty"`
	MaxTokens int                `json:"max_tokens"`
	Stream    bool               `json:"stream,omitempty"`

	// Temperature uses a pointer so zero-value is not sent when not configured.
	Temperature *float64 `json:"temperature,omitempty"`
}

type anthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type anthropicResponse struct {
	Content    []anthropicContent `json:"content"`
	StopReason string             `json:"stop_reason"`
	Usage      anthropicUsage     `json:"usage"`
}

// ------- SSE event types for streaming -------

// anthropicStreamEvent represents a server-sent event from the Anthropic streaming API.
type anthropicStreamEvent struct {
	Type  string          `json:"type"`
	Delta *anthropicDelta `json:"delta,omitempty"`
	Usage *anthropicUsage `json:"usage,omitempty"`
}

type anthropicDelta struct {
	Type       string `json:"type"`
	Text       string `json:"text,omitempty"`
	StopReason string `json:"stop_reason,omitempty"`
}

// ------- LLMProvider implementation -------

// Complete sends a non-streaming Messages request.
func (p *Provider) Complete(ctx context.Context, req ai.CompleteRequest) (ai.CompleteResponse, error) {
	msgs, system := buildMessages(req.Context.SystemPrompt, req.Messages)
	model := resolveModel(req.Context.Model, p.model)

	maxTokens := 4096
	body := anthropicRequest{
		Model:     model,
		Messages:  msgs,
		System:    system,
		MaxTokens: maxTokens,
	}
	if req.Options != nil {
		body.Temperature = req.Options.Temperature
		if req.Options.MaxTokens != nil {
			body.MaxTokens = *req.Options.MaxTokens
		}
	}

	var result anthropicResponse
	if err := p.post(ctx, messagesPath, body, &result); err != nil {
		return ai.CompleteResponse{}, fmt.Errorf("anthropic complete: %w", err)
	}

	if len(result.Content) == 0 {
		return ai.CompleteResponse{}, fmt.Errorf("anthropic complete: no content returned")
	}

	// Collect all text blocks.
	var sb strings.Builder
	for _, block := range result.Content {
		if block.Type == "text" {
			sb.WriteString(block.Text)
		}
	}

	return ai.CompleteResponse{
		Content:      sb.String(),
		FinishReason: result.StopReason,
		Usage: &ai.TokenUsage{
			InputTokens:  result.Usage.InputTokens,
			OutputTokens: result.Usage.OutputTokens,
			TotalTokens:  result.Usage.InputTokens + result.Usage.OutputTokens,
		},
	}, nil
}

// Stream sends a streaming Messages request using SSE.
func (p *Provider) Stream(ctx context.Context, req ai.StreamRequest) (ai.StreamReader, error) {
	msgs, system := buildMessages(req.Context.SystemPrompt, req.Messages)
	model := resolveModel(req.Context.Model, p.model)

	maxTokens := 4096
	body := anthropicRequest{
		Model:     model,
		Messages:  msgs,
		System:    system,
		MaxTokens: maxTokens,
		Stream:    true,
	}
	if req.Options != nil {
		body.Temperature = req.Options.Temperature
		if req.Options.MaxTokens != nil {
			body.MaxTokens = *req.Options.MaxTokens
		}
	}

	bs, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("anthropic stream marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+messagesPath, bytes.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("anthropic stream request: %w", err)
	}
	p.setHeaders(httpReq)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("anthropic stream do: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, &ai.RetryableError{
			Err:        fmt.Errorf("anthropic stream: status %d", resp.StatusCode),
			StatusCode: resp.StatusCode,
		}
	}

	return &anthropicStreamReader{scanner: bufio.NewScanner(resp.Body), body: resp.Body}, nil
}

// Embed is not natively supported by the Anthropic API. This returns an error
// directing callers to use a dedicated embedding provider (e.g., OpenAI).
func (p *Provider) Embed(_ context.Context, _ ai.EmbedRequest) (ai.EmbedResponse, error) {
	return ai.EmbedResponse{}, fmt.Errorf("anthropic provider does not support embeddings; use a dedicated embedding provider")
}

// post marshals body, POSTs to path, and decodes the response into out.
func (p *Provider) post(ctx context.Context, path string, body, out interface{}) error {
	bs, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+path, bytes.NewReader(bs))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	p.setHeaders(req)

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &ai.RetryableError{
			Err:        fmt.Errorf("status %d", resp.StatusCode),
			StatusCode: resp.StatusCode,
		}
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func (p *Provider) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)
}

// ------- stream reader (SSE) -------

type anthropicStreamReader struct {
	scanner      *bufio.Scanner
	body         io.ReadCloser
	pendingEvent string // event type from "event:" line
}

// Recv reads the next SSE event from the Anthropic streaming response.
//
// Anthropic SSE format:
//
//	event: content_block_delta
//	data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"Hello"}}
//
//	event: message_stop
//	data: {"type":"message_stop"}
func (r *anthropicStreamReader) Recv() (ai.StreamChunk, error) {
	for r.scanner.Scan() {
		line := r.scanner.Text()

		if strings.HasPrefix(line, "event: ") {
			r.pendingEvent = strings.TrimPrefix(line, "event: ")
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")

		var event anthropicStreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}

		switch event.Type {
		case "content_block_delta":
			if event.Delta != nil && event.Delta.Type == "text_delta" {
				return ai.StreamChunk{Content: event.Delta.Text}, nil
			}
		case "message_delta":
			if event.Delta != nil && event.Delta.StopReason != "" {
				return ai.StreamChunk{FinishReason: event.Delta.StopReason}, nil
			}
		case "message_stop":
			_ = r.body.Close()
			return ai.StreamChunk{}, io.EOF
		}
	}

	_ = r.body.Close()
	if err := r.scanner.Err(); err != nil {
		return ai.StreamChunk{}, err
	}
	return ai.StreamChunk{}, io.EOF
}

// ------- helpers -------

// buildMessages splits the system prompt out of the messages slice because
// the Anthropic API carries the system prompt as a top-level field, not as
// a "system" role message.
func buildMessages(systemPrompt string, msgs []ai.Message) ([]anthropicMessage, string) {
	out := make([]anthropicMessage, 0, len(msgs))
	for _, m := range msgs {
		// Anthropic does not accept "system" role inside the messages array.
		if m.Role == ai.RoleSystem {
			continue
		}
		out = append(out, anthropicMessage{Role: string(m.Role), Content: m.Content})
	}
	return out, systemPrompt
}

func resolveModel(requestModel, providerDefault string) string {
	if requestModel != "" {
		return requestModel
	}
	return providerDefault
}
