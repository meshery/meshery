package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/meshery/meshery/server/models/connections"
)

// NewProviderFromConnection creates and initializes an appropriate Provider implementation
// based on the provided Meshery connection and credential payloads.
func NewProviderFromConnection(ctx context.Context, conn *connections.ConnectionPayload) (Provider, error) {
	if conn == nil {
		return nil, ErrInvalidConfig(fmt.Errorf("connection payload is nil"))
	}

	providerType := ""
	if p, ok := conn.MetaData["provider"].(string); ok {
		providerType = strings.ToLower(p)
	} else if conn.SubType != "" {
		providerType = strings.ToLower(conn.SubType)
	} else {
		return nil, ErrInvalidConfig(fmt.Errorf("unable to determine provider type from connection payload"))
	}

	serverURL := ""
	if url, ok := conn.MetaData["serverURL"].(string); ok {
		serverURL = url
	}

	config := ProviderConfig{
		Provider:  providerType,
		Model:     conn.Model,
		ServerURL: serverURL,
	}

	// Provide fallback endpoints if the user didn't specify one
	if config.ServerURL == "" {
		switch config.Provider {
		case "openai":
			config.ServerURL = "https://api.openai.com"
		case "claude":
			config.ServerURL = "https://api.anthropic.com"
		case "ollama":
			config.ServerURL = "http://localhost:11434"
		default:
			return nil, ErrInvalidConfig(fmt.Errorf("no serverURL provided and no default exists for provider: %s", config.Provider))
		}
	}

	var p Provider
	switch config.Provider {
	case "openai":
		p = &OpenAIProvider{}
	case "claude":
		p = &ClaudeProvider{}
	case "ollama":
		p = &OllamaProvider{}
	default:
		return nil, ErrInvalidConfig(fmt.Errorf("unsupported provider type: %s", config.Provider))
	}

	err := p.Initialize(ctx, config, conn.CredentialSecret)
	if err != nil {
		return nil, err
	}

	return p, nil
}
