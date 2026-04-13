// Package openai provides an LLMProvider implementation for the OpenAI API.
package openai

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
	defaultBaseURL   = "https://api.openai.com"
	defaultModel     = "gpt-4o"
	defaultTimeout   = 60 * time.Second
	chatPath         = "/v1/chat/completions"
	embeddingsPath   = "/v1/embeddings"
	embeddingModel   = "text-embedding-3-small"
)

// Provider implements ai.LLMProvider for OpenAI.
type Provider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// Option configures an OpenAI Provider.
type Option func(*Provider)

// WithBaseURL overrides the default OpenAI API base URL (useful for proxies / Azure).
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

// New creates an OpenAI provider. apiKey is required.
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

// ------- OpenAI API types -------

type openaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openaiRequest struct {
	Model       string          `json:"model"`
	Messages    []openaiMessage `json:"messages"`
	Stream      bool            `json:"stream,omitempty"`
	Temperature *float64        `json:"temperature,omitempty"`
	MaxTokens   *int            `json:"max_tokens,omitempty"`
}

type openaiChoice struct {
	Message      *openaiMessage `json:"message,omitempty"`
	Delta        *openaiMessage `json:"delta,omitempty"`
	FinishReason string         `json:"finish_reason"`
}

type openaiUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type openaiResponse struct {
	Choices []openaiChoice `json:"choices"`
	Usage   openaiUsage    `json:"usage"`
}

type openaiEmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type openaiEmbedData struct {
	Embedding []float32 `json:"embedding"`
}

type openaiEmbedResponse struct {
	Data  []openaiEmbedData `json:"data"`
	Usage openaiUsage       `json:"usage"`
}

// ------- LLMProvider implementation -------

// Complete sends a non-streaming chat completion request.
func (p *Provider) Complete(ctx context.Context, req ai.CompleteRequest) (ai.CompleteResponse, error) {
	msgs := buildMessages(req.Context.SystemPrompt, req.Messages)
	model := resolveModel(req.Context.Model, p.model)

	body := openaiRequest{
		Model:    model,
		Messages: msgs,
	}
	if req.Options != nil {
		body.Temperature = req.Options.Temperature
		body.MaxTokens = req.Options.MaxTokens
	}

	var result openaiResponse
	if err := p.post(ctx, chatPath, body, &result); err != nil {
		return ai.CompleteResponse{}, fmt.Errorf("openai complete: %w", err)
	}

	if len(result.Choices) == 0 {
		return ai.CompleteResponse{}, fmt.Errorf("openai complete: no choices returned")
	}

	return ai.CompleteResponse{
		Content:      result.Choices[0].Message.Content,
		FinishReason: result.Choices[0].FinishReason,
		Usage: &ai.TokenUsage{
			InputTokens:  result.Usage.PromptTokens,
			OutputTokens: result.Usage.CompletionTokens,
			TotalTokens:  result.Usage.TotalTokens,
		},
	}, nil
}

// Stream sends a streaming chat completion request using SSE.
func (p *Provider) Stream(ctx context.Context, req ai.StreamRequest) (ai.StreamReader, error) {
	msgs := buildMessages(req.Context.SystemPrompt, req.Messages)
	model := resolveModel(req.Context.Model, p.model)

	body := openaiRequest{
		Model:    model,
		Messages: msgs,
		Stream:   true,
	}
	if req.Options != nil {
		body.Temperature = req.Options.Temperature
		body.MaxTokens = req.Options.MaxTokens
	}

	bs, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("openai stream marshal: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+chatPath, bytes.NewReader(bs))
	if err != nil {
		return nil, fmt.Errorf("openai stream request: %w", err)
	}
	p.setHeaders(httpReq)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openai stream do: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, &ai.RetryableError{
			Err:        fmt.Errorf("openai stream: status %d", resp.StatusCode),
			StatusCode: resp.StatusCode,
		}
	}

	return &openaiStreamReader{scanner: bufio.NewScanner(resp.Body), body: resp.Body}, nil
}

// Embed returns an embedding vector for the given text.
func (p *Provider) Embed(ctx context.Context, req ai.EmbedRequest) (ai.EmbedResponse, error) {
	model := req.Context.Model
	if model == "" {
		model = embeddingModel
	}

	body := openaiEmbedRequest{
		Model: model,
		Input: req.Text,
	}

	var result openaiEmbedResponse
	if err := p.post(ctx, embeddingsPath, body, &result); err != nil {
		return ai.EmbedResponse{}, fmt.Errorf("openai embed: %w", err)
	}

	if len(result.Data) == 0 {
		return ai.EmbedResponse{}, fmt.Errorf("openai embed: no embeddings returned")
	}

	return ai.EmbedResponse{
		Vector: result.Data[0].Embedding,
		Usage: &ai.TokenUsage{
			InputTokens: result.Usage.PromptTokens,
			TotalTokens: result.Usage.TotalTokens,
		},
	}, nil
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
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
}

// ------- stream reader (SSE) -------

type openaiStreamReader struct {
	scanner *bufio.Scanner
	body    io.ReadCloser
}

func (r *openaiStreamReader) Recv() (ai.StreamChunk, error) {
	for r.scanner.Scan() {
		line := r.scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")

		if data == "[DONE]" {
			_ = r.body.Close()
			return ai.StreamChunk{}, io.EOF
		}

		var resp openaiResponse
		if err := json.Unmarshal([]byte(data), &resp); err != nil {
			continue
		}

		if len(resp.Choices) == 0 {
			continue
		}

		choice := resp.Choices[0]
		chunk := ai.StreamChunk{FinishReason: choice.FinishReason}
		if choice.Delta != nil {
			chunk.Content = choice.Delta.Content
		}

		return chunk, nil
	}

	_ = r.body.Close()
	if err := r.scanner.Err(); err != nil {
		return ai.StreamChunk{}, err
	}
	return ai.StreamChunk{}, io.EOF
}

// ------- helpers -------

func buildMessages(systemPrompt string, msgs []ai.Message) []openaiMessage {
	out := make([]openaiMessage, 0, len(msgs)+1)
	if systemPrompt != "" {
		out = append(out, openaiMessage{Role: "system", Content: systemPrompt})
	}
	for _, m := range msgs {
		out = append(out, openaiMessage{Role: string(m.Role), Content: m.Content})
	}
	return out
}

func resolveModel(requestModel, providerDefault string) string {
	if requestModel != "" {
		return requestModel
	}
	return providerDefault
}
