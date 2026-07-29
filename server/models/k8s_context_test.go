package models

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	mkerrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/logger"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/meshery/schemas/models/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestFlushMeshSyncDataNilKubernetesServerIDNoPanic verifies that FlushMeshSyncData
// does not panic when K8sContext entries have a nil KubernetesServerID. Before the
// fix, the refcount loop called KubernetesServerID.String() unconditionally,
// triggering a nil dereference whenever a context had no associated server ID.
func TestFlushMeshSyncDataNilKubernetesServerIDNoPanic(t *testing.T) {
	// All contexts have nil KubernetesServerID — the scenario that previously caused
	// a panic. With the nil guard in place, refCount stays 0 and the flush block is
	// skipped entirely, so no provider, broadcast, or logger calls are made.
	k8sctxs := []*K8sContext{
		{ID: "ctx-1", Name: "cluster-a", KubernetesServerID: nil},
		{ID: "ctx-2", Name: "cluster-b", KubernetesServerID: nil},
	}
	ctx := context.WithValue(context.Background(), AllKubeClusterKey, k8sctxs)
	k8sCtx := K8sContext{ID: "ctx-1", Name: "cluster-a", Server: "https://k8s.example.com"}

	// Should complete without panic. If the guard is missing this will panic with a
	// nil dereference on KubernetesServerID.String() inside the refcount loop.
	FlushMeshSyncData(ctx, k8sCtx, nil, nil, "00000000-0000-0000-0000-000000000000", nil, nil)
}

// TestFlushMeshSyncDataMixedNilAndPopulatedServerIDs verifies that the refcount loop
// correctly skips nil entries while still counting non-nil ones, and that no panic
// occurs when both nil and populated KubernetesServerID values are present.
func TestFlushMeshSyncDataMixedNilAndPopulatedServerIDs(t *testing.T) {
	serverID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate UUID: %v", err)
	}

	k8sctxs := []*K8sContext{
		{ID: "ctx-1", Name: "cluster-a", KubernetesServerID: nil},
		{ID: "ctx-2", Name: "cluster-b", KubernetesServerID: &serverID},
		{ID: "ctx-3", Name: "cluster-c", KubernetesServerID: nil},
	}
	ctx := context.WithValue(context.Background(), AllKubeClusterKey, k8sctxs)
	// ctx-1 has nil ServerID — refCount for its sid ("") will be 0, flush skipped.
	k8sCtx := K8sContext{ID: "ctx-1", Name: "cluster-a", Server: "https://k8s.example.com"}

	FlushMeshSyncData(ctx, k8sCtx, nil, nil, "00000000-0000-0000-0000-000000000000", nil, nil)
}

// newMeshSyncTestDB spins up an in-memory SQLite database with the MeshSync
// resource tables migrated, returning a meshkit database handler over it.
func newMeshSyncTestDB(t *testing.T) *database.Handler {
	t.Helper()
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}
	// A ":memory:" database is scoped to a single connection, so a pooled second
	// connection would see none of the migrated tables ("no such table"). Pin the
	// pool to one connection to keep the migrated schema visible for the whole test.
	sqlDB, err := gdb.DB()
	if err != nil {
		t.Fatalf("failed to access underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := gdb.AutoMigrate(
		&meshsyncmodel.KubernetesKeyValue{},
		&meshsyncmodel.KubernetesResource{},
		&meshsyncmodel.KubernetesResourceSpec{},
		&meshsyncmodel.KubernetesResourceStatus{},
		&meshsyncmodel.KubernetesResourceObjectMeta{},
	); err != nil {
		t.Fatalf("failed to migrate meshsync tables: %v", err)
	}
	return &database.Handler{DB: gdb}
}

// seedMeshSyncResource inserts a KubernetesResource together with one child row in
// every dependent table, all sharing the resource's id (the same shape MeshSync
// persists). Children are written as independent rows rather than via associations
// so the linkage under test is explicit and SetID cannot rewrite the ids.
func seedMeshSyncResource(t *testing.T, db *database.Handler, id, clusterID string) {
	t.Helper()
	rows := []interface{}{
		// KubernetesResourceMeta left nil so the BeforeCreate SetID hook is a no-op
		// and the explicit id/cluster_id survive.
		&meshsyncmodel.KubernetesResource{ID: id, ClusterID: clusterID, Kind: "Service", APIVersion: "v1"},
		&meshsyncmodel.KubernetesResourceObjectMeta{ID: id, ClusterID: clusterID, Name: "svc-" + id},
		&meshsyncmodel.KubernetesResourceSpec{ID: id, Attribute: "{}"},
		&meshsyncmodel.KubernetesResourceStatus{ID: id, Attribute: "{}"},
		&meshsyncmodel.KubernetesKeyValue{ID: id, UniqueID: id + "-kv", Kind: meshsyncmodel.KindLabel, Key: "app", Value: "demo"},
	}
	for _, row := range rows {
		if err := db.Create(row).Error; err != nil {
			t.Fatalf("failed to seed %T for %s: %v", row, id, err)
		}
	}
}

// TestFlushMeshSyncResourcesForCluster is a regression test for the stale
// "objects" table-name subqueries in FlushMeshSyncData. It asserts that flushing a
// cluster removes the parent resource and every child row (spec, status,
// object-meta, key-value) for that cluster, while leaving another cluster's rows
// untouched. Before the fix the child cleanup targeted a non-existent "objects"
// table, so the flush errored out and MeshSync inventory was never removed on
// cluster deletion.
func TestFlushMeshSyncResourcesForCluster(t *testing.T) {
	db := newMeshSyncTestDB(t)

	const (
		flushedCluster = "cluster-to-flush"
		keptCluster    = "cluster-to-keep"
		flushedID      = "res-flushed"
		keptID         = "res-kept"
	)
	seedMeshSyncResource(t, db, flushedID, flushedCluster)
	seedMeshSyncResource(t, db, keptID, keptCluster)

	if err := FlushMeshSyncResourcesForCluster(db, flushedCluster); err != nil {
		t.Fatalf("FlushMeshSyncResourcesForCluster returned error: %v", err)
	}

	assertCount := func(model interface{}, where string, arg interface{}, want int64, label string) {
		t.Helper()
		var got int64
		if err := db.Model(model).Where(where, arg).Count(&got).Error; err != nil {
			t.Fatalf("counting %s: %v", label, err)
		}
		if got != want {
			t.Errorf("%s: got %d rows, want %d", label, got, want)
		}
	}

	// The flushed cluster must have no parent and no child rows left behind.
	assertCount(&meshsyncmodel.KubernetesResource{}, "cluster_id = ?", flushedCluster, 0, "resources (flushed cluster)")
	assertCount(&meshsyncmodel.KubernetesResourceObjectMeta{}, "id = ?", flushedID, 0, "object meta (flushed cluster)")
	assertCount(&meshsyncmodel.KubernetesResourceSpec{}, "id = ?", flushedID, 0, "spec (flushed cluster)")
	assertCount(&meshsyncmodel.KubernetesResourceStatus{}, "id = ?", flushedID, 0, "status (flushed cluster)")
	assertCount(&meshsyncmodel.KubernetesKeyValue{}, "id = ?", flushedID, 0, "key value (flushed cluster)")

	// The other cluster's parent and child rows must be untouched.
	assertCount(&meshsyncmodel.KubernetesResource{}, "cluster_id = ?", keptCluster, 1, "resources (kept cluster)")
	assertCount(&meshsyncmodel.KubernetesResourceObjectMeta{}, "id = ?", keptID, 1, "object meta (kept cluster)")
	assertCount(&meshsyncmodel.KubernetesResourceSpec{}, "id = ?", keptID, 1, "spec (kept cluster)")
	assertCount(&meshsyncmodel.KubernetesResourceStatus{}, "id = ?", keptID, 1, "status (kept cluster)")
	assertCount(&meshsyncmodel.KubernetesKeyValue{}, "id = ?", keptID, 1, "key value (kept cluster)")
}

// TestFlushMeshSyncResourcesForClusterNilHandler verifies the helper reports the
// structured MeshKit error (ErrEmptyMeshSyncHandler) instead of panicking, both when
// the handler itself is nil and when only its embedded *gorm.DB is nil.
func TestFlushMeshSyncResourcesForClusterNilHandler(t *testing.T) {
	cases := map[string]*database.Handler{
		"nil handler":     nil,
		"nil embedded DB": {DB: nil},
	}
	for name, db := range cases {
		t.Run(name, func(t *testing.T) {
			err := FlushMeshSyncResourcesForCluster(db, "any-cluster")
			if err == nil {
				t.Fatal("expected an error for an unusable database handler, got nil")
			}
			if code := mkerrors.GetCode(err); code != ErrEmptyMeshSyncHandlerCode {
				t.Errorf("expected error code %s, got %s", ErrEmptyMeshSyncHandlerCode, code)
			}
		})
	}
}

// TestFlushMeshSyncResourcesForClusterRollsBackOnError verifies the cleanup is
// atomic: if a delete fails partway through, the whole flush rolls back rather
// than leaving the cluster partially flushed. The object-meta delete is the last
// child, so dropping its table makes that delete fail after the earlier child
// deletes have run inside the transaction; all of them must be rolled back.
func TestFlushMeshSyncResourcesForClusterRollsBackOnError(t *testing.T) {
	db := newMeshSyncTestDB(t)

	const cluster = "cluster-atomic"
	const id = "res-atomic"
	seedMeshSyncResource(t, db, id, cluster)

	if err := db.Migrator().DropTable(&meshsyncmodel.KubernetesResourceObjectMeta{}); err != nil {
		t.Fatalf("failed to drop object-meta table: %v", err)
	}

	if err := FlushMeshSyncResourcesForCluster(db, cluster); err == nil {
		t.Fatal("expected an error when a child delete fails, got nil")
	}

	// The transaction must have rolled back, so the parent and the child rows that
	// were deleted before the failing delete are all still present.
	assertCount := func(model interface{}, where string, arg interface{}, want int64, label string) {
		t.Helper()
		var got int64
		if err := db.Model(model).Where(where, arg).Count(&got).Error; err != nil {
			t.Fatalf("counting %s: %v", label, err)
		}
		if got != want {
			t.Errorf("%s: got %d rows, want %d", label, got, want)
		}
	}
	assertCount(&meshsyncmodel.KubernetesResource{}, "cluster_id = ?", cluster, 1, "resources (rolled back)")
	assertCount(&meshsyncmodel.KubernetesResourceSpec{}, "id = ?", id, 1, "spec (rolled back)")
	assertCount(&meshsyncmodel.KubernetesResourceStatus{}, "id = ?", id, 1, "status (rolled back)")
	assertCount(&meshsyncmodel.KubernetesKeyValue{}, "id = ?", id, 1, "key value (rolled back)")
}

// multiContextKubeconfig builds a kubeconfig with n contexts (ctx-0..ctx-n-1),
// current-context set to the first, each pointing at a distinct unreachable
// loopback server with token auth (no cert files, so the kube handler builds and
// only the API-server lookup fails). Contexts are surfaced as unreachable.
func multiContextKubeconfig(n int) []byte {
	var clusters, contexts, users string
	for i := 0; i < n; i++ {
		clusters += fmt.Sprintf("- cluster:\n    server: https://127.0.0.1:%d\n  name: cluster-%d\n", 59900+i, i)
		contexts += fmt.Sprintf("- context:\n    cluster: cluster-%d\n    user: user-%d\n  name: ctx-%d\n", i, i, i)
		users += fmt.Sprintf("- name: user-%d\n  user:\n    token: token-%d\n", i, i)
	}
	return []byte(fmt.Sprintf(
		"apiVersion: v1\nkind: Config\ncurrent-context: ctx-0\nclusters:\n%scontexts:\n%susers:\n%s",
		clusters, contexts, users,
	))
}

// TestK8sContextsFromKubeconfigDiscoversAllContexts guards against the regression
// where discovery enumerated contexts via kubernetes.ProcessConfig, whose
// clientcmd MinifyConfig pass prunes every context except current-context. That
// made importing a multi-context kubeconfig surface only the current context.
// Enumerating from the un-minified config must return every context.
func TestK8sContextsFromKubeconfigDiscoversAllContexts(t *testing.T) {
	log, err := logger.New("test", logger.Options{Format: logger.JsonLogFormat})
	if err != nil {
		t.Fatalf("failed to build logger: %v", err)
	}
	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))

	const wantContexts = 3
	kubeconfig := multiContextKubeconfig(wantContexts)
	eventMetadata := map[string]interface{}{}

	// includeUnreachable=true mirrors the import wizard: unreachable contexts are
	// still returned (flagged Reachable=false) so the user can register them.
	got := K8sContextsFromKubeconfigWithOptions(nil, uuid.Must(uuid.NewV4()).String(), nil, kubeconfig, &instanceID, eventMetadata, log, true)

	if len(got) != wantContexts {
		var names []string
		for _, kc := range got {
			names = append(names, kc.Name)
		}
		t.Fatalf("discovered %d contexts %v, want %d (all contexts in the kubeconfig, not just current-context)", len(got), names, wantContexts)
	}

	var names []string
	for _, kc := range got {
		names = append(names, kc.Name)
	}
	sort.Strings(names)
	for i, name := range names {
		if want := fmt.Sprintf("ctx-%d", i); name != want {
			t.Errorf("context[%d] = %q, want %q", i, name, want)
		}
	}
}
