package ai

import (
	"errors"
	"sort"
	"sync"
)

var (
	// ErrProviderNameRequired is returned when a provider is registered without a name.
	ErrProviderNameRequired = errors.New("provider name is required")
	// ErrProviderNil is returned when a nil provider is registered.
	ErrProviderNil = errors.New("provider is nil")
	// ErrProviderAlreadyRegistered is returned when a provider name is reused.
	ErrProviderAlreadyRegistered = errors.New("provider is already registered")
	// ErrProviderNotFound is returned when a requested provider is not registered.
	ErrProviderNotFound = errors.New("provider not found")
)

// Registry stores AI providers by name.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]LLMProvider
}

// NewRegistry creates an empty provider registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]LLMProvider),
	}
}

// Register stores a provider under a unique name.
func (r *Registry) Register(name string, provider LLMProvider) error {
	if name == "" {
		return ErrProviderNameRequired
	}
	if provider == nil {
		return ErrProviderNil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.providers[name]; exists {
		return ErrProviderAlreadyRegistered
	}

	r.providers[name] = provider
	return nil
}

// Get returns a registered provider by name.
func (r *Registry) Get(name string) (LLMProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	provider, exists := r.providers[name]
	if !exists {
		return nil, ErrProviderNotFound
	}

	return provider, nil
}

// Names returns the registered provider names in sorted order.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

// Has reports whether a provider with the given name exists.
func (r *Registry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.providers[name]
	return exists
}
