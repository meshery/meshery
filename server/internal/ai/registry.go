package ai

import (
	"fmt"
	"sort"
	"sync"
)

// Registry holds named providers and supports concurrent access.
// Providers are registered at startup and looked up per-request.
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

// Register adds a provider under its Name(). Overwrites if the name
// already exists — lets callers re-register with updated config.
func (r *Registry) Register(p LLMProvider) {
	if p == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[p.Name()] = p
}

// Get returns the provider registered under name, or an error if none exists.
func (r *Registry) Get(name string) (LLMProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("ai: provider %q not registered", name)
	}
	return p, nil
}

// Has reports whether a provider is registered under name.
func (r *Registry) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.providers[name]
	return ok
}

// Names returns all registered provider names in sorted order.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.providers))
	for n := range r.providers {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
