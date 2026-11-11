package dataloader_test

import (
	"context"
	"net/http"
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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type mockProvider struct {
	models.DefaultLocalProvider
	dbHandler *database.Handler
}

func (m *mockProvider) GetGenericPersister() *database.Handler {
	return m.dbHandler
}

type queryCountingLogger struct {
	count *int64
}

func (l *queryCountingLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}
func (l *queryCountingLogger) Info(ctx context.Context, s string, i ...interface{})  {}
func (l *queryCountingLogger) Warn(ctx context.Context, s string, i ...interface{})  {}
func (l *queryCountingLogger) Error(ctx context.Context, s string, i ...interface{}) {}
func (l *queryCountingLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	atomic.AddInt64(l.count, 1)
}

func setupTestDB(t *testing.T) (*database.Handler, *int64) {
	t.Helper()
	var queryCount int64
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: &queryCountingLogger{count: &queryCount},
	})
	require.NoError(t, err)

	sqlDB, err := gdb.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)

	err = gdb.AutoMigrate(
		&meshsyncmodel.KubernetesKeyValue{},
		&meshsyncmodel.KubernetesResource{},
		&meshsyncmodel.KubernetesResourceSpec{},
		&meshsyncmodel.KubernetesResourceStatus{},
		&meshsyncmodel.KubernetesResourceObjectMeta{},
	)
	require.NoError(t, err)

	return &database.Handler{DB: gdb}, &queryCount
}

func seedPod(t *testing.T, db *database.Handler, id, clusterID, name, namespace, containerImage string) {
	t.Helper()
	res := meshsyncmodel.KubernetesResource{
		ID:        id,
		ClusterID: clusterID,
		Kind:      "Pod",
	}
	require.NoError(t, db.Create(&res).Error)

	meta := meshsyncmodel.KubernetesResourceObjectMeta{
		ID:           id,
		Name:         name,
		Namespace:    namespace,
		GenerateName: name + "-123",
	}
	require.NoError(t, db.Create(&meta).Error)

	spec := meshsyncmodel.KubernetesResourceSpec{
		ID:        id,
		Attribute: `{"containers":[{"name":"istiod","image":"` + containerImage + `"}]}`,
	}
	require.NoError(t, db.Create(&spec).Error)
}

func TestControlPlaneDataLoader_BatchingAndQueryReduction(t *testing.T) {
	dbHandler, queryCount := setupTestDB(t)

	// Seed 5 clusters with Istio pods
	clusterIDs := []string{"cluster-1", "cluster-2", "cluster-3", "cluster-4", "cluster-5"}
	for _, cid := range clusterIDs {
		seedPod(t, dbHandler, "pod-"+cid, cid, "istiod-"+cid, "istio-system", "docker.io/istio/pilot:1.20.0")
	}

	prov := &mockProvider{dbHandler: dbHandler}
	selectors := []model.MeshType{model.MeshTypeIstio}

	// 1. Without DataLoader (Legacy N queries pattern):
	// Loops over each cluster ID and runs separate queries for each.
	atomic.StoreInt64(queryCount, 0)
	for _, cid := range clusterIDs {
		_, err := model.GetControlPlaneState(context.Background(), selectors, prov, []string{cid})
		require.NoError(t, err)
	}
	queriesWithoutLoader := atomic.LoadInt64(queryCount)

	// 2. With DataLoader (Batched single query pattern):
	atomic.StoreInt64(queryCount, 0)
	ctx := context.Background()
	loaderMiddleware := dataloader.Middleware(func(ctx context.Context) models.Provider {
		return prov
	})

	var loader *dataloader.Loaders
	// Simulate request lifecycle
	handler := loaderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader = dataloader.For(r.Context())
	}))
	req := &http.Request{}
	req = req.WithContext(ctx)
	handler.ServeHTTP(nil, req)

	require.NotNil(t, loader)
	require.NotNil(t, loader.ControlPlaneLoader)

	// Load all 5 clusters via DataLoader
	result, err := loader.ControlPlaneLoader.LoadMany(clusterIDs, selectors)
	require.NoError(t, err)
	assert.Len(t, result, 5)

	// Verify each cluster has its control plane
	for _, cid := range clusterIDs {
		cps := result[cid]
		require.NotEmpty(t, cps, "expected control planes for %s", cid)
		assert.Equal(t, "istio", cps[0].Name)
		require.NotEmpty(t, cps[0].Members)
		assert.Equal(t, "istiod-"+cid, cps[0].Members[0].Name)
	}

	queriesWithLoader := atomic.LoadInt64(queryCount)

	t.Logf("Query count without DataLoader (N=5): %d", queriesWithoutLoader)
	t.Logf("Query count with DataLoader (N=5): %d", queriesWithLoader)

	// DataLoader should execute significantly fewer queries (1 batched query vs 5 separate queries)
	assert.Less(t, queriesWithLoader, queriesWithoutLoader, "DataLoader should reduce query count from O(N) to O(1)")

	// 3. Cache Test: Loading the same clusters again in the same request should result in 0 new queries
	atomic.StoreInt64(queryCount, 0)
	cachedResult, err := loader.ControlPlaneLoader.LoadMany(clusterIDs, selectors)
	require.NoError(t, err)
	assert.Len(t, cachedResult, 5)

	queriesFromCache := atomic.LoadInt64(queryCount)
	assert.Equal(t, int64(0), queriesFromCache, "Subsequent lookups within the same request must be served from cache with 0 queries")
}
