package resolver

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/internal/graphql/dataloader"
	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type mockResolverProvider struct {
	dbHandler *database.Handler
}

func (m *mockResolverProvider) GetGenericPersister() *database.Handler {
	return m.dbHandler
}

func setupResolverTestDB(t *testing.T) *database.Handler {
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

func seedResolverService(t *testing.T, db *database.Handler, id, clusterID, name, namespace string) {
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
}

func newTestResolver(t *testing.T) *Resolver {
	t.Helper()
	log, err := logger.New("test-resolver", logger.Options{
		Format:   logger.SyslogLogFormat,
		LogLevel: int(logrus.ErrorLevel),
		Output:   io.Discard,
	})
	require.NoError(t, err)
	return &Resolver{Log: log}
}

// TestGetAvailableAddons_SelectorFiltering verifies that resolver-side filtering
// properly isolates addons by requested mesh type, even though the underlying DataLoader
// loads and caches the full AllMeshType set.
func TestGetAvailableAddons_SelectorFiltering(t *testing.T) {
	dbHandler := setupResolverTestDB(t)

	// Seed cluster-1 with Grafana (Istio) and Prometheus (Linkerd)
	seedResolverService(t, dbHandler, "svc-istio-grafana", "cluster-1", "grafana", "istio-system")
	seedResolverService(t, dbHandler, "svc-linkerd-prom", "cluster-1", "prometheus", "linkerd-system")

	prov := &mockResolverProvider{dbHandler: dbHandler}
	middleware := dataloader.Middleware(func(ctx context.Context) models.GenericPersisterGetter {
		return prov
	})

	var requestLoaders *dataloader.Loaders
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestLoaders = dataloader.For(r.Context())
	}))

	ctx := context.Background()
	req := httptest.NewRequest("POST", "/graphql", nil).WithContext(ctx)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.NotNil(t, requestLoaders)
	reqCtx := dataloader.WithLoaders(ctx, requestLoaders)

	r := newTestResolver(t)

	// Case 1: Filter by Istio only -> must only return Grafana (Owner: ISTIO), Prometheus must be excluded
	istioType := model.MeshTypeIstio
	istioFilter := &model.ServiceMeshFilter{
		Type:          &istioType,
		K8sClusterIDs: []string{"cluster-1"},
	}
	istioAddons, err := r.getAvailableAddons(reqCtx, prov, istioFilter)
	require.NoError(t, err)
	require.Len(t, istioAddons, 1)
	assert.Equal(t, "grafana", istioAddons[0].Name)
	assert.Equal(t, model.MeshTypeIstio.String(), istioAddons[0].Owner)

	// Case 2: Filter by Linkerd only -> must only return Prometheus (Owner: LINKERD), Grafana must be excluded
	linkerdType := model.MeshTypeLinkerd
	linkerdFilter := &model.ServiceMeshFilter{
		Type:          &linkerdType,
		K8sClusterIDs: []string{"cluster-1"},
	}
	linkerdAddons, err := r.getAvailableAddons(reqCtx, prov, linkerdFilter)
	require.NoError(t, err)
	require.Len(t, linkerdAddons, 1)
	assert.Equal(t, "prometheus", linkerdAddons[0].Name)
	assert.Equal(t, model.MeshTypeLinkerd.String(), linkerdAddons[0].Owner)

	// Case 3: Filter AllMesh -> returns both services across AllMesh, Istio, and Linkerd selectors
	allMeshType := model.MeshTypeAllMesh
	allFilter := &model.ServiceMeshFilter{
		Type:          &allMeshType,
		K8sClusterIDs: []string{"cluster-1"},
	}
	allAddons, err := r.getAvailableAddons(reqCtx, prov, allFilter)
	require.NoError(t, err)
	// ALL_MESH matches both (2), ISTIO matches grafana (1), LINKERD matches prometheus (1) -> total 4
	require.Len(t, allAddons, 4)
	owners := make([]string, 0, len(allAddons))
	for _, a := range allAddons {
		owners = append(owners, a.Owner)
	}
	assert.Equal(t, []string{
		model.MeshTypeAllMesh.String(),
		model.MeshTypeAllMesh.String(),
		model.MeshTypeIstio.String(),
		model.MeshTypeLinkerd.String(),
	}, owners)
}

// TestGetAvailableAddons_MeshTypeGroupedOrdering verifies that the returned addons
// are ordered grouped by mesh type (selectors outer loop) across multiple clusters.
func TestGetAvailableAddons_MeshTypeGroupedOrdering(t *testing.T) {
	dbHandler := setupResolverTestDB(t)

	// Seed cluster-1 and cluster-2 with both Istio (grafana) and Linkerd (prometheus) services
	seedResolverService(t, dbHandler, "svc-c1-istio", "c1", "grafana", "istio-system")
	seedResolverService(t, dbHandler, "svc-c1-linkerd", "c1", "prometheus", "linkerd-system")
	seedResolverService(t, dbHandler, "svc-c2-istio", "c2", "grafana", "istio-system")
	seedResolverService(t, dbHandler, "svc-c2-linkerd", "c2", "prometheus", "linkerd-system")

	prov := &mockResolverProvider{dbHandler: dbHandler}
	middleware := dataloader.Middleware(func(ctx context.Context) models.GenericPersisterGetter {
		return prov
	})

	var requestLoaders *dataloader.Loaders
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestLoaders = dataloader.For(r.Context())
	}))

	ctx := context.Background()
	req := httptest.NewRequest("POST", "/graphql", nil).WithContext(ctx)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	require.NotNil(t, requestLoaders)
	reqCtx := dataloader.WithLoaders(ctx, requestLoaders)

	r := newTestResolver(t)

	filter := &model.ServiceMeshFilter{
		K8sClusterIDs: []string{"c1", "c2"},
	}
	addons, err := r.getAvailableAddons(reqCtx, prov, filter)
	require.NoError(t, err)
	// ALL_MESH (4) + ISTIO (2) + LINKERD (2) = 8
	require.Len(t, addons, 8)

	// The addons must be grouped by mesh type (all ALL_MESH first, then ISTIO, then LINKERD)
	var seenIstio, seenLinkerd bool
	for _, a := range addons {
		if a.Owner == model.MeshTypeIstio.String() {
			seenIstio = true
		}
		if a.Owner == model.MeshTypeLinkerd.String() {
			seenLinkerd = true
		}
		if seenIstio && a.Owner == model.MeshTypeAllMesh.String() {
			t.Errorf("ALL_MESH must precede ISTIO in mesh-type grouped ordering")
		}
		if seenLinkerd && (a.Owner == model.MeshTypeIstio.String() || a.Owner == model.MeshTypeAllMesh.String()) {
			t.Errorf("ISTIO and ALL_MESH must precede LINKERD in mesh-type grouped ordering")
		}
	}
}

// TestGetAvailableAddons_FallbackWithoutLoader verifies clean fallback to GetAddonsState
// when dataloader is absent from context.
func TestGetAvailableAddons_FallbackWithoutLoader(t *testing.T) {
	dbHandler := setupResolverTestDB(t)
	seedResolverService(t, dbHandler, "svc-fb-grafana", "c-fb", "grafana", "istio-system")

	prov := &mockResolverProvider{dbHandler: dbHandler}
	r := newTestResolver(t)

	// Context with NO dataloader
	ctx := context.Background()
	filter := &model.ServiceMeshFilter{
		K8sClusterIDs: []string{"c-fb"},
	}
	addons, err := r.getAvailableAddons(ctx, prov, filter)
	require.NoError(t, err)
	// ALL_MESH (1) + ISTIO (1) = 2
	require.Len(t, addons, 2)
	assert.Equal(t, "grafana", addons[0].Name)
	assert.Equal(t, model.MeshTypeAllMesh.String(), addons[0].Owner)
	assert.Equal(t, "grafana", addons[1].Name)
	assert.Equal(t, model.MeshTypeIstio.String(), addons[1].Owner)
}
