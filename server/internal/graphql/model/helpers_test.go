package model_test

import (
	"context"
	"testing"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshkit/database"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// mockPersister is a lightweight stub implementing models.GenericPersisterGetter.
type mockPersister struct {
	dbHandler *database.Handler
}

func (m *mockPersister) GetGenericPersister() *database.Handler {
	return m.dbHandler
}

func setupTestDB(t *testing.T) *database.Handler {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
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

	return &database.Handler{DB: gdb}
}

func TestHelpers_NilProviderSafety(t *testing.T) {
	ctx := context.Background()

	// 1. GetControlPlaneState with nil provider
	cpState, err := model.GetControlPlaneState(ctx, model.AllMeshType, nil, []string{"c1"})
	assert.NoError(t, err)
	assert.Len(t, cpState, len(model.AllMeshType))
	for _, cp := range cpState {
		assert.Empty(t, cp.Members)
	}

	// GetControlPlaneState with nil provider and empty cids
	cpStateEmptyCids, err := model.GetControlPlaneState(ctx, model.AllMeshType, nil, []string{})
	assert.NoError(t, err)
	assert.Empty(t, cpStateEmptyCids)

	// 2. BatchGetControlPlaneState with nil provider
	batchCp, err := model.BatchGetControlPlaneState(ctx, model.AllMeshType, nil, []string{"c1"})
	assert.NoError(t, err)
	assert.NotNil(t, batchCp)
	assert.Empty(t, batchCp["c1"])

	// 3. GetAddonsState with nil provider
	addons, err := model.GetAddonsState(ctx, model.AllMeshType, nil, []string{"c1"})
	assert.NoError(t, err)
	assert.Empty(t, addons)

	// 4. GetDataPlaneState with nil provider
	dpState, err := model.GetDataPlaneState(ctx, model.AllMeshType, nil, []string{"c1"})
	assert.NoError(t, err)
	assert.Empty(t, dpState)

	// 5. SelectivelyFetchNamespaces with nil provider
	namespaces, err := model.SelectivelyFetchNamespaces([]string{"c1"}, nil)
	assert.NoError(t, err)
	assert.Empty(t, namespaces)
}

func TestHelpers_LightweightStubExecution(t *testing.T) {
	db := setupTestDB(t)
	persister := &mockPersister{dbHandler: db}
	ctx := context.Background()

	// 1. GetControlPlaneState with lightweight stub
	cpState, err := model.GetControlPlaneState(ctx, model.AllMeshType, persister, []string{"cluster-1"})
	assert.NoError(t, err)
	assert.Len(t, cpState, len(model.AllMeshType))
	for _, cp := range cpState {
		assert.Empty(t, cp.Members)
	}

	// 2. BatchGetControlPlaneState with lightweight stub
	batchMap, err := model.BatchGetControlPlaneState(ctx, model.AllMeshType, persister, []string{"cluster-1"})
	assert.NoError(t, err)
	assert.NotNil(t, batchMap)
	assert.Len(t, batchMap["cluster-1"], len(model.AllMeshType))
	for _, cp := range batchMap["cluster-1"] {
		assert.Empty(t, cp.Members)
	}

	// 3. GetAddonsState with lightweight stub
	addons, err := model.GetAddonsState(ctx, model.AllMeshType, persister, []string{"cluster-1"})
	assert.NoError(t, err)
	assert.Empty(t, addons)

	// 4. GetDataPlaneState with lightweight stub
	dp, err := model.GetDataPlaneState(ctx, model.AllMeshType, persister, []string{"cluster-1"})
	assert.NoError(t, err)
	assert.Empty(t, dp)

	// 5. SelectivelyFetchNamespaces with lightweight stub
	ns, err := model.SelectivelyFetchNamespaces([]string{"cluster-1"}, persister)
	assert.NoError(t, err)
	assert.Empty(t, ns)
}
