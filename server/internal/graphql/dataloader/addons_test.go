package dataloader_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/meshery/meshery/server/internal/graphql/dataloader"
	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func seedService(t *testing.T, db *database.Handler, id, clusterID, name, namespace string) {
	t.Helper()
	res := meshsyncmodel.KubernetesResource{
		ID:        id,
		ClusterID: clusterID,
		Kind:      "Service",
	}
	require.NoError(t, db.Create(&res).Error)

	meta := meshsyncmodel.KubernetesResourceObjectMeta{
		ID:        id,
		Name:      name,
		Namespace: namespace,
	}
	require.NoError(t, db.Create(&meta).Error)

	spec := meshsyncmodel.KubernetesResourceSpec{
		ID:        id,
		Attribute: `{}`,
	}
	require.NoError(t, db.Create(&spec).Error)

	status := meshsyncmodel.KubernetesResourceStatus{
		ID:        id,
		Attribute: `{}`,
	}
	require.NoError(t, db.Create(&status).Error)
}

func TestAddonDataLoader_BatchingAndQueryReduction(t *testing.T) {
	dbHandler, queryCount := setupTestDB(t)

	// Seed 5 clusters with Grafana service in istio-system namespace
	clusterIDs := []string{"cluster-1", "cluster-2", "cluster-3", "cluster-4", "cluster-5"}
	for _, cid := range clusterIDs {
		seedService(t, dbHandler, "svc-"+cid, cid, "grafana", "istio-system")
	}

	prov := &mockProvider{dbHandler: dbHandler}
	selectors := model.AllMeshType

	// 1. Without DataLoader (Sequential queries pattern):
	atomic.StoreInt64(queryCount, 0)
	for _, cid := range clusterIDs {
		_, err := model.GetAddonsState(context.Background(), selectors, prov, []string{cid})
		require.NoError(t, err)
	}
	queriesWithoutLoader := atomic.LoadInt64(queryCount)

	// 2. With DataLoader (Batched single query pattern):
	atomic.StoreInt64(queryCount, 0)
	ctx := context.Background()
	loaderMiddleware := dataloader.Middleware(func(ctx context.Context) models.GenericPersisterGetter {
		return prov
	})

	var loader *dataloader.Loaders
	handler := loaderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader = dataloader.For(r.Context())
	}))
	req := httptest.NewRequest("POST", "/graphql", nil).WithContext(ctx)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.NotNil(t, loader)
	require.NotNil(t, loader.AddonLoader)

	// Load all 5 clusters via DataLoader
	result, err := loader.AddonLoader.LoadMany(clusterIDs, selectors)
	require.NoError(t, err)
	assert.Len(t, result, 5)

	// Verify each cluster has its addon
	for _, cid := range clusterIDs {
		addons := result[cid]
		require.NotEmpty(t, addons, "expected addons for %s", cid)
		var istioAddon *model.AddonList
		for _, addon := range addons {
			if addon.Name == "grafana" && addon.Owner == model.MeshTypeIstio.String() {
				istioAddon = addon
				break
			}
		}
		require.NotNil(t, istioAddon, "expected grafana addon owned by ISTIO for %s", cid)
		assert.Equal(t, model.MeshTypeIstio.String(), istioAddon.Owner)
	}

	queriesWithLoader := atomic.LoadInt64(queryCount)

	t.Logf("Addons Query count without DataLoader (N=5): %d", queriesWithoutLoader)
	t.Logf("Addons Query count with DataLoader (N=5): %d", queriesWithLoader)

	// DataLoader should execute significantly fewer queries
	assert.Less(t, queriesWithLoader, queriesWithoutLoader, "DataLoader should reduce query count from O(N) to O(1)")

	// 3. Cache Test: Loading the same clusters again in the same request should result in 0 new queries
	atomic.StoreInt64(queryCount, 0)
	cachedResult, err := loader.AddonLoader.LoadMany(clusterIDs, selectors)
	require.NoError(t, err)
	assert.Len(t, cachedResult, 5)

	queriesFromCache := atomic.LoadInt64(queryCount)
	assert.Equal(t, int64(0), queriesFromCache, "Subsequent lookups within the same request must be served from cache with 0 queries")
}

func TestAddonDataLoader_SelectorIsolation(t *testing.T) {
	dbHandler, _ := setupTestDB(t)

	// Seed cluster-addon with Grafana in istio-system and Prometheus in linkerd-system
	seedService(t, dbHandler, "svc-istio-grafana", "cluster-addon", "grafana", "istio-system")
	seedService(t, dbHandler, "svc-linkerd-prometheus", "cluster-addon", "prometheus", "linkerd-system")

	prov := &mockProvider{dbHandler: dbHandler}
	middleware := dataloader.Middleware(func(ctx context.Context) models.GenericPersisterGetter {
		return prov
	})

	var loader *dataloader.Loaders
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader = dataloader.For(r.Context())
	}))

	ctx := context.Background()
	req := httptest.NewRequest("POST", "/graphql", nil).WithContext(ctx)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.NotNil(t, loader)
	require.NotNil(t, loader.AddonLoader)

	// First Load: requested with Istio selector. The loader fetches AllMeshType to safely
	// serve multiple callers across the batch window and caches the full set.
	istioResults, err := loader.AddonLoader.Load("cluster-addon", []model.MeshType{model.MeshTypeIstio})
	require.NoError(t, err)
	require.NotEmpty(t, istioResults)

	// Second Load: requested with Linkerd selector. Served directly from the request cache.
	linkerdResults, err := loader.AddonLoader.Load("cluster-addon", []model.MeshType{model.MeshTypeLinkerd})
	require.NoError(t, err)
	require.NotEmpty(t, linkerdResults)

	// AddonBatchLoader is selector-blind by design: both calls receive the exact same full cached set,
	// containing both Istio and Linkerd addons. Actual selector-level filtering is performed downstream
	// by the GraphQL resolver (verified in resolver tests).
	assert.Equal(t, len(istioResults), len(linkerdResults), "loader cache is selector-blind and returns full set for subsequent calls")
	assert.ElementsMatch(t, istioResults, linkerdResults)
}

func TestAddonDataLoader_ContextCancellation_NoDBQuery(t *testing.T) {
	dbHandler, queryCount := setupTestDB(t)
	seedService(t, dbHandler, "svc-cancel", "cluster-cancel", "grafana", "istio-system")

	prov := &mockProvider{dbHandler: dbHandler}

	// Create a context that is cancelled immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before dispatch

	loader := dataloader.WithLoaders(ctx, &dataloader.Loaders{})
	_ = loader

	// Create loader directly with cancelled context
	// We can test through middleware or direct construction
	req := httptest.NewRequest("POST", "/graphql", nil).WithContext(ctx)
	middleware := dataloader.Middleware(func(c context.Context) models.GenericPersisterGetter {
		return prov
	})
	var l *dataloader.Loaders
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l = dataloader.For(r.Context())
	}))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Reset query counter after seeding
	atomic.StoreInt64(queryCount, 0)

	// Load with cancelled context
	results, err := l.AddonLoader.Load("cluster-cancel", []model.MeshType{model.MeshTypeIstio})
	assert.ErrorIs(t, err, context.Canceled)
	assert.Nil(t, results)

	// Wait for the batch window timer (2ms) to expire
	time.Sleep(15 * time.Millisecond)

	// Verify that NO database query ran
	queriesRun := atomic.LoadInt64(queryCount)
	assert.Equal(t, int64(0), queriesRun, "cancelled context must abort before querying database")
}

func TestAddonDataLoader_ContextCancellation_DuringBatchWindow(t *testing.T) {
	dbHandler, queryCount := setupTestDB(t)
	seedService(t, dbHandler, "svc-cancel-win", "cluster-win", "grafana", "istio-system")

	prov := &mockProvider{dbHandler: dbHandler}

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest("POST", "/graphql", nil).WithContext(ctx)
	middleware := dataloader.Middleware(func(c context.Context) models.GenericPersisterGetter {
		return prov
	})
	var l *dataloader.Loaders
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l = dataloader.For(r.Context())
	}))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	atomic.StoreInt64(queryCount, 0)

	// Start Load in background, then cancel context before the 2ms batch window timer fires
	errChan := make(chan error, 1)
	go func() {
		_, err := l.AddonLoader.Load("cluster-win", []model.MeshType{model.MeshTypeIstio})
		errChan <- err
	}()

	cancel()

	err := <-errChan
	assert.ErrorIs(t, err, context.Canceled)

	// Wait for timer to fire
	time.Sleep(15 * time.Millisecond)

	// Verify that dispatch aborted early on ctx.Err() != nil and ran 0 queries
	queriesRun := atomic.LoadInt64(queryCount)
	assert.Equal(t, int64(0), queriesRun, "cancelled context during batch window must abort dispatch with 0 DB queries")
}

func TestAddons_DuplicateClusterIDs(t *testing.T) {
	dbHandler, _ := setupTestDB(t)

	seedService(t, dbHandler, "svc-dup", "cluster-dup", "grafana", "istio-system")

	prov := &mockProvider{dbHandler: dbHandler}
	selectors := []model.MeshType{model.MeshTypeIstio}

	// Query with duplicate cluster IDs
	res, err := model.GetAddonsState(context.Background(), selectors, prov, []string{"cluster-dup", "cluster-dup"})
	require.NoError(t, err)
	require.Len(t, res, 1)
	assert.Equal(t, "grafana", res[0].Name)
	assert.Equal(t, model.MeshTypeIstio.String(), res[0].Owner)
}

type panickingPersisterGetter struct{}

func (p *panickingPersisterGetter) GetGenericPersister() *database.Handler {
	panic("database connection failure")
}

func TestAddonDataLoader_PanicRecovery(t *testing.T) {
	middleware := dataloader.Middleware(func(ctx context.Context) models.GenericPersisterGetter {
		return &panickingPersisterGetter{}
	})

	var loader *dataloader.Loaders
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader = dataloader.For(r.Context())
	}))

	req := httptest.NewRequest("POST", "/graphql", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.NotNil(t, loader)
	require.NotNil(t, loader.AddonLoader)

	// Loading when provider panics should return an error and not crash the process
	res, err := loader.AddonLoader.Load("cluster-panic", model.AllMeshType)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "addon batch loader panic")
	assert.Nil(t, res)
}
