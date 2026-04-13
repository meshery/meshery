package ai

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrProviderSelectionRequired is returned when a request does not declare a provider.
	ErrProviderSelectionRequired = errors.New("provider selection is required")
)

// wrappedEntry holds a provider already wrapped with the service middlewares.
type wrappedEntry struct {
	provider LLMProvider
}

// Service orchestrates provider lookup and middleware application for AI requests.
// Providers are wrapped with middlewares once on first use and cached thereafter.
type Service struct {
	registry    *Registry
	middlewares []Middleware

	mu      sync.RWMutex
	wrapped map[string]wrappedEntry
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
		wrapped:     make(map[string]wrappedEntry),
	}
}

// Registry returns the underlying provider registry.
func (s *Service) Registry() *Registry {
	return s.registry
}

// ProviderNames returns the names of all registered providers.
func (s *Service) ProviderNames() []string {
	return s.registry.Names()
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

// providerFor returns the middleware-wrapped provider for the given name.
// The wrapped provider is cached so middlewares are applied only once per provider.
func (s *Service) providerFor(name string) (LLMProvider, error) {
	if name == "" {
		return nil, ErrProviderSelectionRequired
	}

	// Fast path: already wrapped and cached.
	s.mu.RLock()
	if entry, ok := s.wrapped[name]; ok {
		s.mu.RUnlock()
		return entry.provider, nil
	}
	s.mu.RUnlock()

	// Slow path: look up the raw provider and wrap it.
	raw, err := s.registry.Get(name)
	if err != nil {
		return nil, err
	}

	wrapped := ChainMiddleware(raw, s.middlewares...)

	s.mu.Lock()
	// Double-checked locking: another goroutine may have populated the cache.
	if entry, ok := s.wrapped[name]; ok {
		s.mu.Unlock()
		return entry.provider, nil
	}
	s.wrapped[name] = wrappedEntry{provider: wrapped}
	s.mu.Unlock()

	return wrapped, nil
}
