package sessions

import (
	"context"
	"sort"
	"sync"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
)

// ConnectionContext is everything a [Factory] needs to build a driver for one
// connection: the connection itself, and the provider plus token required to
// resolve the credential it references.
type ConnectionContext struct {
	Connection *connections.Connection
	Provider   models.Provider
	// Token is the caller's auth token, used to fetch the connection's
	// credential from the provider. Sessions therefore run with the requesting
	// user's own credentials, so the target's native authorization (Kubernetes
	// RBAC, for a k8s connection) applies unchanged.
	Token string
}

// Factory builds drivers for a single connection kind.
type Factory interface {
	// Kind is the connections.Connection.Kind this factory serves, e.g.
	// "kubernetes".
	Kind() string

	// NewDriver builds a driver bound to cc.Connection.
	NewDriver(ctx context.Context, cc ConnectionContext) (Driver, error)
}

// Registry resolves a connection onto the driver that serves its kind.
//
// It is safe for concurrent use: registration happens at init time, lookup on
// every request.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

// NewRegistry returns an empty registry. Most callers want [Default].
func NewRegistry() *Registry {
	return &Registry{factories: make(map[string]Factory)}
}

// Register adds f to the registry, replacing any factory previously registered
// for the same kind.
func (r *Registry) Register(f Factory) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.factories == nil {
		r.factories = make(map[string]Factory)
	}
	r.factories[f.Kind()] = f
}

// Kinds lists the registered connection kinds, sorted, for diagnostics.
func (r *Registry) Kinds() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	kinds := make([]string, 0, len(r.factories))
	for k := range r.factories {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	return kinds
}

// Driver builds the driver for cc.Connection's kind. It returns an
// ErrNoDriver-coded error when the kind has no registered factory, and an
// ErrDriverInit-coded error when the factory cannot build a driver.
func (r *Registry) Driver(ctx context.Context, cc ConnectionContext) (Driver, error) {
	if cc.Connection == nil {
		return nil, ErrInvalidTarget("no connection was resolved for this session request")
	}

	kind := cc.Connection.Kind
	r.mu.RLock()
	factory, ok := r.factories[kind]
	r.mu.RUnlock()
	if !ok {
		return nil, ErrNoDriver(kind)
	}

	driver, err := factory.NewDriver(ctx, cc)
	if err != nil {
		return nil, ErrDriverInit(err, kind)
	}
	return driver, nil
}

// Default is the process-wide registry. Drivers register themselves into it
// from an init function, and a driver package is pulled in with a blank import
// wherever sessions are served — the same shape as database/sql drivers and
// image format decoders. Adding a new resource universe is then one new package
// plus one import line, with no change to the transport or the HTTP handlers.
var Default = NewRegistry()

// Register adds f to the [Default] registry.
func Register(f Factory) { Default.Register(f) }
