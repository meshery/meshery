package ai

import (
	"context"
	"errors"
)

var (
	// ErrProviderSelectionRequired is returned when a request does not declare a provider.
	ErrProviderSelectionRequired = errors.New("provider selection is required")
)

// Service orchestrates provider lookup and middleware application for AI requests.
type Service struct {
	registry    *Registry
	middlewares []Middleware
}

// NewService creates a service with the given registry and middlewares.
func NewService(registry *Registry, middlewares ...Middleware) *Service {
	if registry == nil {
		registry = NewRegistry()
	}

	copied := append([]Middleware(nil), middlewares...)
	return &Service{
		registry:    registry,
		middlewares: copied,
	}
}

// Registry returns the underlying provider registry.
func (s *Service) Registry() *Registry {
	return s.registry
}

// Complete routes a completion request to the selected provider.
func (s *Service) Complete(ctx context.Context, req CompleteRequest) (CompleteResponse, error) {
	provider, err := s.providerFor(req.Context.Provider)
	if err != nil {
		return CompleteResponse{}, err
	}

	return provider.Complete(ctx, req)
}

// Stream routes a streaming request to the selected provider.
func (s *Service) Stream(ctx context.Context, req StreamRequest) (StreamReader, error) {
	provider, err := s.providerFor(req.Context.Provider)
	if err != nil {
		return nil, err
	}

	return provider.Stream(ctx, req)
}

// Embed routes an embedding request to the selected provider.
func (s *Service) Embed(ctx context.Context, req EmbedRequest) (EmbedResponse, error) {
	provider, err := s.providerFor(req.Context.Provider)
	if err != nil {
		return EmbedResponse{}, err
	}

	return provider.Embed(ctx, req)
}

func (s *Service) providerFor(name string) (LLMProvider, error) {
	if name == "" {
		return nil, ErrProviderSelectionRequired
	}

	provider, err := s.registry.Get(name)
	if err != nil {
		return nil, err
	}

	return ChainMiddleware(provider, s.middlewares...), nil
}
