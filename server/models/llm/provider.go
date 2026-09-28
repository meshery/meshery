package llm

import (
	"context"
)

// ProviderConfig represents the generic configuration for an AI provider.
// It maps closely to the non-secret data inside a meshery Connection.
type ProviderConfig struct {
	Provider     string `json:"provider"`     // e.g., "openai", "claude", "ollama"
	Model        string `json:"model"`        // Selected model name
	ServerURL    string `json:"serverURL"`    // Base URL for the provider API
}

// ModelCapabilities defines the metadata and limits for a specific model.
type ModelCapabilities struct {
	MaxContextWindow int  `json:"maxContextWindow"`
	SupportsTools    bool `json:"supportsTools"`
	SupportsImages   bool `json:"supportsImages"`
}

// ModelInfo describes a single model supported by the provider.
type ModelInfo struct {
	ID           string            `json:"id"`
	Capabilities ModelCapabilities `json:"capabilities"`
}

// GenerateRequest represents a generation request sent to the provider.
type GenerateRequest struct {
	SystemPrompt string
	UserPrompt   string
	SchemaContext string // Future use: Relevant Meshery schema context
}

// GenerateResponse represents the result of a generation request.
type GenerateResponse struct {
	RawResponse string
}

// Provider represents the common abstraction for all AI/LLM providers.
// It handles initialization, connectivity checks, model discovery, and inference.
type Provider interface {
	// Initialize validates the configuration and credential, setting up the provider client.
	Initialize(ctx context.Context, config ProviderConfig, credentialSecret map[string]interface{}) error

	// Ping performs a lightweight connectivity check.
	// It should respect context cancellation and timeouts, and return secret-safe errors.
	Ping(ctx context.Context) error

	// Models returns the list of models available from this provider.
	Models(ctx context.Context) ([]ModelInfo, error)

	// Generate performs a structured inference request.
	Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
}
