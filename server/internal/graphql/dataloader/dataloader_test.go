package dataloader_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/meshery/meshery/server/internal/graphql/dataloader"
	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
)

func TestDataLoaderMiddleware_PerRequestIsolation(t *testing.T) {
	var loader1, loader2 *dataloader.Loaders

	middleware := dataloader.Middleware(func(ctx context.Context) models.Provider {
		return nil
	})

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := dataloader.For(r.Context())
		if loader1 == nil {
			loader1 = l
		} else {
			loader2 = l
		}
	}))

	// Request 1
	req1 := httptest.NewRequest("POST", "/graphql", nil)
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)

	// Request 2
	req2 := httptest.NewRequest("POST", "/graphql", nil)
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	assert.NotNil(t, loader1, "loader1 should not be nil")
	assert.NotNil(t, loader2, "loader2 should not be nil")
	assert.NotEqual(t, loader1, loader2, "each request must receive its own distinct DataLoader instance")
	assert.NotNil(t, loader1.ControlPlaneLoader)
	assert.NotNil(t, loader2.ControlPlaneLoader)
	assert.NotEqual(t, loader1.ControlPlaneLoader, loader2.ControlPlaneLoader, "loaders must not be shared across requests")
}

func TestControlPlaneBatchLoader_EmptyClusterIDs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	loaders := &dataloader.Loaders{
		ControlPlaneLoader: nil,
	}
	ctxWithLoaders := dataloader.WithLoaders(ctx, loaders)

	retrieved := dataloader.For(ctxWithLoaders)
	assert.Equal(t, loaders, retrieved)
}

func TestControlPlaneBatchLoader_LoadEmptyKey(t *testing.T) {
	middleware := dataloader.Middleware(func(ctx context.Context) models.Provider {
		return nil
	})

	var loader *dataloader.Loaders
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader = dataloader.For(r.Context())
	}))

	req := httptest.NewRequest("POST", "/graphql", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	res, err := loader.ControlPlaneLoader.Load("", []model.MeshType{model.MeshTypeIstio})
	assert.NoError(t, err)
	assert.Nil(t, res)

	resMany, err := loader.ControlPlaneLoader.LoadMany([]string{}, []model.MeshType{model.MeshTypeIstio})
	assert.NoError(t, err)
	assert.Empty(t, resMany)
}

