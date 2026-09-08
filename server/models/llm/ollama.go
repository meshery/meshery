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
)

// OllamaProvider implements the Provider interface for local Ollama instances.
type OllamaProvider struct {
	config     ProviderConfig
	httpClient *http.Client
}

func (p *OllamaProvider) Initialize(ctx context.Context, config ProviderConfig, credentialSecret map[string]interface{}) error {
	p.config = config

	// Validate URL without enforcing HTTPS, as Ollama is usually local HTTP
	// Allow private/loopback IPs since this is a local provider
	u, err := ValidateURL(config.ServerURL, true)
	if err != nil {
		return err
	}

	// Double check we are not accidentally hitting cloud endpoints over HTTP.
	// If it's a known public cloud, we shouldn't use HTTP.
	if !isLocalHostOrLAN(u.Hostname()) {
		// Log a warning or strictly enforce HTTPS depending on security requirements.
		// For now, we allow it if it passed ValidateURL, to support custom VPN/WAN setups,
		// but typically Ollama is local.
	}

	p.httpClient = &http.Client{
		Timeout: 300 * time.Second, // Local generation can be slow
	}

	return nil
}

func (p *OllamaProvider) Ping(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/tags", strings.TrimRight(p.config.ServerURL, "/"))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return NormalizeError(err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return NormalizeError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrProviderUnavailable(fmt.Errorf("provider returned status %d", resp.StatusCode))
	}

	return nil
}

func (p *OllamaProvider) Models(ctx context.Context) ([]ModelInfo, error) {
	url := fmt.Sprintf("%s/api/tags", strings.TrimRight(p.config.ServerURL, "/"))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, NormalizeError(err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, NormalizeError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrProviderUnavailable(fmt.Errorf("provider returned status %d", resp.StatusCode))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrMalformedResponse(fmt.Errorf("failed to read response body"))
	}

	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, ErrMalformedResponse(fmt.Errorf("failed to parse JSON response: %v", err))
	}

	var models []ModelInfo
	for _, m := range result.Models {
		models = append(models, ModelInfo{
			ID: m.Name,
			Capabilities: ModelCapabilities{
				MaxContextWindow: 8192, // Default Ollama context window is typically smaller
				SupportsTools:    false, // Varies by model in Ollama
				SupportsImages:   false,
			},
		})
	}

	return models, nil
}

func (p *OllamaProvider) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	url := fmt.Sprintf("%s/api/chat", strings.TrimRight(p.config.ServerURL, "/"))

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
		"stream":   false,
	}

	jsonBytes, err := json.Marshal(bodyData)
	if err != nil {
		return nil, ErrLLMInference(fmt.Errorf("failed to marshal request: %v", err))
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, NormalizeError(err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, NormalizeError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrProviderUnavailable(fmt.Errorf("provider returned status %d", resp.StatusCode))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrMalformedResponse(fmt.Errorf("failed to read response body"))
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, ErrMalformedResponse(fmt.Errorf("failed to parse JSON response: %v", err))
	}

	return &GenerateResponse{
		RawResponse: result.Message.Content,
	}, nil
}
