// Package providers contains a factory for creating LLMProvider instances from
// a generic configuration, and the concrete provider implementations.
package providers

import (
	"fmt"
	"strings"

	"github.com/meshery/meshery/server/internal/ai"
	"github.com/meshery/meshery/server/internal/ai/providers/anthropic"
	"github.com/meshery/meshery/server/internal/ai/providers/ollama"
	"github.com/meshery/meshery/server/internal/ai/providers/openai"
)

// ProviderType identifies a supported LLM backend.
type ProviderType string

const (
	ProviderOpenAI    ProviderType = "openai"
	ProviderOllama    ProviderType = "ollama"
	ProviderAnthropic ProviderType = "anthropic"
)

// Config carries all the information needed to instantiate a provider.
type Config struct {
	// Type is the provider backend ("openai", "ollama", "anthropic").
	Type ProviderType `json:"type"`

	// APIKey is used by cloud providers (OpenAI, Anthropic). Ignored by Ollama.
	APIKey string `json:"api_key,omitempty"`

	// Endpoint overrides the default API base URL.
	// For Ollama this defaults to http://localhost:11434.
	// For OpenAI / Anthropic this is optional (useful for proxies or Azure).
	Endpoint string `json:"endpoint,omitempty"`

	// Model overrides the provider's default model.
	Model string `json:"model,omitempty"`
}

// New creates an LLMProvider from config. Returns an error if the provider
// type is unknown or if required fields (e.g. APIKey) are missing.
func New(cfg Config) (ai.LLMProvider, error) {
	switch ProviderType(strings.ToLower(string(cfg.Type))) {
	case ProviderOpenAI:
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("openai provider requires an api_key")
		}
		opts := []openai.Option{}
		if cfg.Endpoint != "" {
			opts = append(opts, openai.WithBaseURL(cfg.Endpoint))
		}
		if cfg.Model != "" {
			opts = append(opts, openai.WithModel(cfg.Model))
		}
		return openai.New(cfg.APIKey, opts...), nil

	case ProviderOllama:
		opts := []ollama.Option{}
		if cfg.Endpoint != "" {
			opts = append(opts, ollama.WithEndpoint(cfg.Endpoint))
		}
		if cfg.Model != "" {
			opts = append(opts, ollama.WithModel(cfg.Model))
		}
		return ollama.New(opts...), nil

	case ProviderAnthropic:
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("anthropic provider requires an api_key")
		}
		opts := []anthropic.Option{}
		if cfg.Endpoint != "" {
			opts = append(opts, anthropic.WithBaseURL(cfg.Endpoint))
		}
		if cfg.Model != "" {
			opts = append(opts, anthropic.WithModel(cfg.Model))
		}
		return anthropic.New(cfg.APIKey, opts...), nil

	default:
		return nil, fmt.Errorf("unknown provider type %q; supported: openai, ollama, anthropic", cfg.Type)
	}
}
