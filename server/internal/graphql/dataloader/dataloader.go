package dataloader

import (
	"context"
	"net/http"

	"github.com/meshery/meshery/server/models"
)

type contextKey struct{ name string }

var loadersKey = &contextKey{name: "dataloaders"}

// Loaders holds all request-scoped data loaders.
type Loaders struct {
	ControlPlaneLoader *ControlPlaneBatchLoader
}

// Middleware injects freshly instantiated, request-scoped DataLoaders into the request context.
// Resolvers run concurrently across requests; a per-request cache ensures no cross-request data leakage.
func Middleware(providerGetter func(ctx context.Context) models.Provider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			loaders := &Loaders{
				ControlPlaneLoader: newControlPlaneBatchLoader(ctx, providerGetter),
			}
			r = r.WithContext(context.WithValue(ctx, loadersKey, loaders))
			next.ServeHTTP(w, r)
		})
	}
}

// For extracts the request-scoped Loaders from context.
func For(ctx context.Context) *Loaders {
	loaders, ok := ctx.Value(loadersKey).(*Loaders)
	if !ok {
		return nil
	}
	return loaders
}

// WithLoaders sets the Loaders instance in a given context (useful for testing or sub-contexts).
func WithLoaders(ctx context.Context, loaders *Loaders) context.Context {
	return context.WithValue(ctx, loadersKey, loaders)
}

