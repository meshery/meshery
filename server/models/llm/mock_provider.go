package llm

import (
	"context"
	"errors"
)

// MockProvider is a deterministic provider used for testing.
type MockProvider struct {
	ShouldFailInitialize bool
	ShouldFailPing       bool
	ShouldFailModels     bool
	ShouldFailGenerate   bool

	MockResponse string
}

func (m *MockProvider) Initialize(ctx context.Context, config ProviderConfig, credentialSecret map[string]interface{}) error {
	if m.ShouldFailInitialize {
		return ErrInvalidConfig(errors.New("mock initialization failure"))
	}
	return nil
}

func (m *MockProvider) Ping(ctx context.Context) error {
	if m.ShouldFailPing {
		return ErrProviderUnavailable(errors.New("mock ping failure"))
	}
	return nil
}

func (m *MockProvider) Models(ctx context.Context) ([]ModelInfo, error) {
	if m.ShouldFailModels {
		return nil, ErrProviderUnavailable(errors.New("mock models failure"))
	}
	return []ModelInfo{
		{
			ID: "mock-model",
			Capabilities: ModelCapabilities{
				MaxContextWindow: 1000,
				SupportsTools:    false,
				SupportsImages:   false,
			},
		},
	}, nil
}

func (m *MockProvider) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	if m.ShouldFailGenerate {
		return nil, ErrLLMInference(errors.New("mock generation failure"))
	}
	
	resp := m.MockResponse
	if resp == "" {
		resp = "Mock generation success"
	}
	
	return &GenerateResponse{
		RawResponse: resp,
	}, nil
}
