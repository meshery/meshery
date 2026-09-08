package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/meshery/meshery/server/models"
)

// ClaudeProvider implements the Provider interface for Anthropic.
type ClaudeProvider struct {
	config     ProviderConfig
	apiKey     string
	httpClient *http.Client
}

func (p *ClaudeProvider) Initialize(ctx context.Context, config ProviderConfig, credentialSecret map[string]interface{}) error {
	p.config = config

	u, err := ValidateURL(config.ServerURL, false)
	if err != nil {
		return err
	}
	if err := EnsureHTTPS(u); err != nil {
		return err
	}

	if credentialSecret == nil {
		return ErrInvalidCredentials(fmt.Errorf("credential secret map is nil"))
	}
	
	// Use Meshery's canonical credential unwrapper
	payload := models.CredentialPayload(credentialSecret)
	if payload == nil {
		p.apiKey = models.CredentialAuthSecret(credentialSecret)
	} else {
		if val, ok := payload["apiKey"].(string); ok && val != "" {
			p.apiKey = val
		} else if val, ok := payload["secret"].(string); ok && val != "" {
			p.apiKey = val
		}
	}
	
	if p.apiKey == "" {
		return ErrInvalidCredentials(fmt.Errorf("API key not found in credential secret"))
	}
	p.httpClient = &http.Client{
		Timeout: 60 * time.Second,
	}

	return nil
}

func (p *ClaudeProvider) Ping(ctx context.Context) error {
	// Anthropic has a /v1/models endpoint as well, or we can just hit a lightweight route.
	// Since /v1/models is standard for Anthropic now, we'll try that. 
	// If it fails with 404, we could try a dummy message, but /v1/models is safe.
	url := fmt.Sprintf("%s/v1/models", strings.TrimRight(p.config.ServerURL, "/"))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return NormalizeError(err)
	}

	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return NormalizeError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ErrAuthFailure(fmt.Errorf("unauthorized (status %d)", resp.StatusCode))
	}

	if resp.StatusCode != http.StatusOK {
		return ErrProviderUnavailable(fmt.Errorf("provider returned status %d", resp.StatusCode))
	}

	return nil
}

func (p *ClaudeProvider) Models(ctx context.Context) ([]ModelInfo, error) {
	knownModels := []ModelInfo{
		{
			ID: "claude-3-5-sonnet-20240620",
			Capabilities: ModelCapabilities{
				MaxContextWindow: 200000,
				SupportsTools:    true,
				SupportsImages:   true,
			},
		},
		{
			ID: "claude-3-opus-20240229",
			Capabilities: ModelCapabilities{
				MaxContextWindow: 200000,
				SupportsTools:    true,
				SupportsImages:   true,
			},
		},
		{
			ID: "claude-3-haiku-20240307",
			Capabilities: ModelCapabilities{
				MaxContextWindow: 200000,
				SupportsTools:    true,
				SupportsImages:   true,
			},
		},
	}
	return knownModels, nil
}

func (p *ClaudeProvider) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	url := fmt.Sprintf("%s/v1/messages", strings.TrimRight(p.config.ServerURL, "/"))

	messages := []map[string]string{
		{
			"role":    "user",
			"content": req.UserPrompt,
		},
	}

	systemStr := req.SystemPrompt
	if req.SchemaContext != "" {
		systemStr += "\n\nContext:\n" + req.SchemaContext
	}

	bodyData := map[string]interface{}{
		"model":      p.config.Model,
		"max_tokens": 4096,
		"messages":   messages,
	}

	if systemStr != "" {
		bodyData["system"] = systemStr
	}

	jsonBytes, err := json.Marshal(bodyData)
	if err != nil {
		return nil, ErrLLMInference(fmt.Errorf("failed to marshal request: %w", err))
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, NormalizeError(err)
	}

	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, NormalizeError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, ErrAuthFailure(fmt.Errorf("unauthorized (status %d)", resp.StatusCode))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrProviderUnavailable(fmt.Errorf("provider returned status %d", resp.StatusCode))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrMalformedResponse(fmt.Errorf("failed to read response body"))
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, ErrMalformedResponse(fmt.Errorf("failed to parse JSON response: %w", err))
	}

	if len(result.Content) == 0 {
		return nil, ErrMalformedResponse(fmt.Errorf("provider returned 0 content blocks"))
	}

	// Find the text block
	var textResponse string
	for _, block := range result.Content {
		if block.Type == "text" {
			textResponse = block.Text
			break
		}
	}

	return &GenerateResponse{
		RawResponse: textResponse,
	}, nil
}
