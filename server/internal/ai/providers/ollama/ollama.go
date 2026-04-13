// Package ollama provides an LLMProvider implementation for Ollama local inference.
package ollama

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
	defaultEndpoint   = "http://localhost:11434"
	defaultModel      = "llama3"
	defaultTimeout    = 120 * time.Second
	chatPath          = "/api/chat"
	embedPath         = "/api/embeddings"
)

// Provider implements ai.LLMProvider for Ollama.
type Provider struct {
	endpoint string
	model    string
	client   *http.Client
}

// Option configures an Ollama Provider.
type Option func(*Provider)

// WithEndpoint overrides the default Ollama server URL.
func WithEndpoint(endpoint string) Option {
	return func(p *Provider) {
		p.endpoint = strings.TrimRight(endpoint, "/")
	}
}

// WithModel overrides the default model name.
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

// New creates an Ollama provider with the given options.
func New(opts ...Option) *Provider {
	p := &Provider{
		endpoint: defaultEndpoint,
		model:    defaultModel,
		client:   &http.Client{Timeout: defaultTimeout},
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

// ------- Ollama API types -------

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ollamaChatResponse struct {
	Model   string        `json:"model"`
	Message ollamaMessage `json:"message"`
	Done    bool          `json:"done"`
}

type ollamaEmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaEmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

// ------- LLMProvider implementation -------

// Complete sends a non-streaming chat request to Ollama.
func (p *Provider) Complete(ctx context.Context, req ai.CompleteRequest) (ai.CompleteResponse, error) {
	msgs := buildMessages(req.Context.SystemPrompt, req.Messages)
	model := resolveModel(req.Context.Model, p.model)

	body := ollamaChatRequest{
		Model:    model,
		Messages: msgs,
		Stream:   false,
	}

	var result ollamaChatResponse
	if err := p.post(ctx, chatPath, body, &result); err != nil {
		return ai.CompleteResponse{}, fmt.Errorf("ollama complete: %w", err)
	}

	return ai.CompleteResponse{
		Content:      result.Message.Content,
		FinishReason: "stop",
	}, nil
}

// Stream sends a streaming chat request to Ollama and returns a StreamReader.
func (p *Provider) Stream(ctx context.Context, req ai.StreamRequest) (ai.StreamReader, error) {
	msgs := buildMessages(req.Context.SystemPrompt, req.Messages)
	model := resolveModel(req.Context.Model, p.model)

	body := ollamaChatRequest{
		Model:    model,
		Messages: msgs,
		Stream:   true,
	}

	bs, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("ollama stream marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint+chatPath, bytes.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("ollama stream request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama stream do: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, &ai.RetryableError{
			Err:        fmt.Errorf("ollama stream: status %d", resp.StatusCode),
			StatusCode: resp.StatusCode,
		}
	}

	return &ollamaStreamReader{scanner: bufio.NewScanner(resp.Body), body: resp.Body}, nil
}

// Embed returns an embedding vector for the given text.
func (p *Provider) Embed(ctx context.Context, req ai.EmbedRequest) (ai.EmbedResponse, error) {
	model := resolveModel(req.Context.Model, p.model)

	body := ollamaEmbedRequest{
		Model:  model,
		Prompt: req.Text,
	}

	var result ollamaEmbedResponse
	if err := p.post(ctx, embedPath, body, &result); err != nil {
		return ai.EmbedResponse{}, fmt.Errorf("ollama embed: %w", err)
	}

	return ai.EmbedResponse{Vector: result.Embedding}, nil
}

// post is a helper that marshals body, POSTs to path, and decodes into out.
func (p *Provider) post(ctx context.Context, path string, body, out interface{}) error {
	bs, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpoint+path, bytes.NewReader(bs))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

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

// ------- stream reader -------

type ollamaStreamReader struct {
	scanner *bufio.Scanner
	body    io.ReadCloser
}

func (r *ollamaStreamReader) Recv() (ai.StreamChunk, error) {
	if !r.scanner.Scan() {
		_ = r.body.Close()
		if err := r.scanner.Err(); err != nil {
			return ai.StreamChunk{}, err
		}
		return ai.StreamChunk{}, io.EOF
	}

	line := r.scanner.Text()
	if line == "" {
		return r.Recv()
	}

	var msg ollamaChatResponse
	if err := json.Unmarshal([]byte(line), &msg); err != nil {
		return ai.StreamChunk{}, fmt.Errorf("ollama stream parse: %w", err)
	}

	chunk := ai.StreamChunk{Content: msg.Message.Content}
	if msg.Done {
		chunk.FinishReason = "stop"
		_ = r.body.Close()
	}

	return chunk, nil
}

// ------- helpers -------

func buildMessages(systemPrompt string, msgs []ai.Message) []ollamaMessage {
	out := make([]ollamaMessage, 0, len(msgs)+1)
	if systemPrompt != "" {
		out = append(out, ollamaMessage{Role: "system", Content: systemPrompt})
	}
	for _, m := range msgs {
		out = append(out, ollamaMessage{Role: string(m.Role), Content: m.Content})
	}
	return out
}

func resolveModel(requestModel, providerDefault string) string {
	if requestModel != "" {
		return requestModel
	}
	return providerDefault
}
