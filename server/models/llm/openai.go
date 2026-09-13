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

// OpenAIProvider implements the Provider interface for OpenAI and OpenAI-compatible endpoints.
type OpenAIProvider struct {
	config     ProviderConfig
	apiKey     string
	httpClient *http.Client
}

// Initialize validates the configuration and credential, setting up the provider client.
func (p *OpenAIProvider) Initialize(ctx context.Context, config ProviderConfig, credentialSecret map[string]interface{}) error {
	p.config = config

	// Validate URL (OpenAI is a hosted provider, so block local IPs)
	u, err := ValidateURL(config.ServerURL, false)
	if err != nil {
		return err
	}
	if err := EnsureHTTPS(u); err != nil {
		return err
	}

	// Resolve API Key from credentials
	if credentialSecret == nil {
		return ErrInvalidCredentials(fmt.Errorf("credential secret map is nil"))
	}
	
	// Use Meshery's canonical credential unwrapper
	var apiKey string
	payload := models.CredentialPayload(credentialSecret)
	if payload == nil {
		// Fallback to seeing if it's a bare string legacy secret
		apiKey = models.CredentialAuthSecret(credentialSecret)
	} else {
		if val, ok := payload["apiKey"].(string); ok && val != "" {
			apiKey = val
		} else if val, ok := payload["secret"].(string); ok && val != "" {
			apiKey = val
		}
	}
	
	if apiKey == "" {
		return ErrInvalidCredentials(fmt.Errorf("API key not found in credential secret"))
	}
	
	p.apiKey = apiKey
	p.httpClient = &http.Client{
		Timeout: 60 * time.Second, // Bounded request timeout
	}

	return nil
}

// Ping performs a lightweight connectivity check.
func (p *OpenAIProvider) Ping(ctx context.Context) error {
	// OpenAI has a lightweight /models endpoint
	url := fmt.Sprintf("%s/v1/models", strings.TrimRight(p.config.ServerURL, "/"))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return NormalizeError(err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return NormalizeError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrAuthFailure(fmt.Errorf("unauthorized (401)"))
	}

	if resp.StatusCode != http.StatusOK {
		return ErrProviderUnavailable(fmt.Errorf("provider returned status %d", resp.StatusCode))
	}

	return nil
}

// Models returns the list of models available from this provider.
func (p *OpenAIProvider) Models(ctx context.Context) ([]ModelInfo, error) {
	// For MVP, we define the known OpenAI models with their capabilities.
	// We can fetch from /v1/models, but that doesn't return context window limits dynamically.
	knownModels := []ModelInfo{
		{
			ID: "gpt-4o",
			Capabilities: ModelCapabilities{
				MaxContextWindow: 128000,
				SupportsTools:    true,
				SupportsImages:   true,
			},
		},
		{
			ID: "gpt-4-turbo",
			Capabilities: ModelCapabilities{
				MaxContextWindow: 128000,
				SupportsTools:    true,
				SupportsImages:   true,
			},
		},
		{
			ID: "gpt-3.5-turbo",
			Capabilities: ModelCapabilities{
				MaxContextWindow: 16385,
				SupportsTools:    true,
				SupportsImages:   false,
			},
		},
	}
	return knownModels, nil
}

// Generate performs a structured inference request.
func (p *OpenAIProvider) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	url := fmt.Sprintf("%s/v1/chat/completions", strings.TrimRight(p.config.ServerURL, "/"))

	// Construct request body
	messages := []map[string]string{}
	
	if req.SystemPrompt != "" || req.SchemaContext != "" {
		sysContent := req.SystemPrompt
		if req.SchemaContext != "" {
			sysContent += "\n\nContext:\n" + req.SchemaContext
		}
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": sysContent,
		})
	}

	messages = append(messages, map[string]string{
		"role":    "user",
		"content": req.UserPrompt,
	})

	bodyData := map[string]interface{}{
		"model":    p.config.Model,
		"messages": messages,
		// Temperature, max_tokens, etc. could be added here
	}

	jsonBytes, err := json.Marshal(bodyData)
	if err != nil {
		return nil, ErrLLMInference(fmt.Errorf("failed to marshal request: %w", err))
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, NormalizeError(err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, NormalizeError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrAuthFailure(fmt.Errorf("unauthorized (401)"))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrProviderUnavailable(fmt.Errorf("provider returned status %d", resp.StatusCode))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrMalformedResponse(fmt.Errorf("failed to read response body"))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, ErrMalformedResponse(fmt.Errorf("failed to parse JSON response: %w", err))
	}

	if len(result.Choices) == 0 {
		return nil, ErrMalformedResponse(fmt.Errorf("provider returned 0 choices"))
	}

	return &GenerateResponse{
		RawResponse: result.Choices[0].Message.Content,
	}, nil
}
