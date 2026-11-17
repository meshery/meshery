package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/internal/sql"
	"github.com/meshery/meshkit/database"
	mkerrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
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

type flushFakeProvider struct {
	Provider
	persistedEvents []events.Event
}

func (f *flushFakeProvider) PersistSystemEvent(event events.Event) (PersistedEvent, error) {
	f.persistedEvents = append(f.persistedEvents, event)
	return PersistedEvent{Event: &event, Disposition: EventCreated}, nil
}

func (f *flushFakeProvider) LoadAllK8sContext(string) ([]*K8sContext, error) {
	return nil, nil
}

func (f *flushFakeProvider) GetK8sContext(string, string) (K8sContext, error) {
	return K8sContext{}, nil
}

// TestFlushMeshSyncDataEmptyUserIDAndNilMesheryInstanceID verifies that FlushMeshSyncData
// does not panic when userID is empty or mesheryInstanceID is nil, and handles event persistence safely.
func TestFlushMeshSyncDataEmptyUserIDAndNilMesheryInstanceID(t *testing.T) {
	k8sCtx := K8sContext{ID: "ctx-1", Name: "cluster-a", Server: "https://k8s.example.com"}

	t.Run("empty context list and empty user ID persists event without owner", func(t *testing.T) {
		fp := &flushFakeProvider{}
		FlushMeshSyncData(context.Background(), k8sCtx, fp, nil, "", nil, nil)
		if len(fp.persistedEvents) != 1 {
			t.Fatalf("expected 1 persisted event, got %d", len(fp.persistedEvents))
		}
		if fp.persistedEvents[0].Category != "meshsync" || fp.persistedEvents[0].Action != "flush" {
			t.Errorf("unexpected event category/action: %s/%s", fp.persistedEvents[0].Category, fp.persistedEvents[0].Action)
		}
	})

	t.Run("valid user ID and nil mesheryInstanceID with cluster list", func(t *testing.T) {
		k8sctxs := []*K8sContext{
			{ID: "ctx-1", Name: "cluster-a", KubernetesServerID: nil},
		}
		ctx := context.WithValue(context.Background(), AllKubeClusterKey, k8sctxs)
		fp := &flushFakeProvider{}
		validUser := uuid.Must(uuid.NewV4()).String()
		FlushMeshSyncData(ctx, k8sCtx, fp, nil, validUser, nil, nil)
		if len(fp.persistedEvents) != 0 {
			t.Fatalf("expected no persisted events, got %d", len(fp.persistedEvents))
		}
	})

	t.Run("nil provider does not panic", func(t *testing.T) {
		FlushMeshSyncData(context.Background(), k8sCtx, nil, nil, "", nil, nil)
	})
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

func TestK8sContextGenerateIDNilMesheryInstanceID(t *testing.T) {
	ctxWithoutInstanceID := K8sContext{
		Name:    "minikube",
		Cluster: map[string]interface{}{"server": "https://127.0.0.1:6443"},
		Auth:    map[string]interface{}{"token": "test-token"},
	}

	idWithoutInstance, err := K8sContextGenerateID(ctxWithoutInstanceID)
	if err != nil {
		t.Fatalf("K8sContextGenerateID returned error for nil MesheryInstanceID: %v", err)
	}
	if idWithoutInstance == "" {
		t.Fatal("expected non-empty deterministic ID for context with nil MesheryInstanceID")
	}

	idRepeat, err := K8sContextGenerateID(ctxWithoutInstanceID)
	if err != nil {
		t.Fatalf("K8sContextGenerateID repeat call returned error: %v", err)
	}
	if idRepeat != idWithoutInstance {
		t.Fatalf("expected deterministic ID %q, got %q", idWithoutInstance, idRepeat)
	}

	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctxWithInstanceID := ctxWithoutInstanceID
	ctxWithInstanceID.MesheryInstanceID = &instanceID

	idWithInstance, err := K8sContextGenerateID(ctxWithInstanceID)
	if err != nil {
		t.Fatalf("K8sContextGenerateID returned error for populated MesheryInstanceID: %v", err)
	}
	if idWithInstance == idWithoutInstance {
		t.Fatalf("expected ID with MesheryInstanceID (%q) to differ from nil MesheryInstanceID (%q)", idWithInstance, idWithoutInstance)
	}

	// RedactCredentialsForContext clears MesheryInstanceID to nil; calling
	// K8sContextGenerateID afterward must remain panic-free.
	redacted := RedactCredentialsForContext(&ctxWithInstanceID)
	if redacted.MesheryInstanceID != nil {
		t.Fatalf("expected RedactCredentialsForContext to nil MesheryInstanceID, got %v", redacted.MesheryInstanceID)
	}
	if _, err := K8sContextGenerateID(redacted); err != nil {
		t.Fatalf("K8sContextGenerateID returned error after RedactCredentialsForContext: %v", err)
	}
}

type flushDatabaseProvider struct {
	flushFakeProvider
	db *database.Handler
}

func (p *flushDatabaseProvider) GetGenericPersister() *database.Handler { return p.db }

func TestFlushMeshSyncDataSharedAndLastContext(t *testing.T) {
	serverID := uuid.Must(uuid.NewV4())
	for _, tc := range []struct {
		name   string
		shared bool
		want   int64
	}{{"last context", false, 0}, {"shared server", true, 1}} {
		t.Run(tc.name, func(t *testing.T) {
			db := newMeshSyncTestDB(t)
			seedMeshSyncResource(t, db, "resource", serverID.String())
			current := &K8sContext{ID: uuid.Must(uuid.NewV4()).String(), KubernetesServerID: &serverID}
			contexts := []*K8sContext{current}
			if tc.shared {
				contexts = append(contexts, &K8sContext{ID: uuid.Must(uuid.NewV4()).String(), KubernetesServerID: &serverID})
			}
			ctx := context.WithValue(context.Background(), AllKubeClusterKey, contexts)
			FlushMeshSyncData(ctx, *current, &flushDatabaseProvider{db: db}, nil, "", nil, nil)
			var count int64
			if err := db.Model(&meshsyncmodel.KubernetesResource{}).Where("cluster_id = ?", serverID.String()).Count(&count).Error; err != nil {
				t.Fatal(err)
			}
			if count != tc.want {
				t.Fatalf("remaining resources = %d, want %d", count, tc.want)
			}
		})
	}
}

type flushMockContextLoaderProvider struct {
	flushDatabaseProvider
	allContexts []*K8sContext
}

func (p *flushMockContextLoaderProvider) LoadAllK8sContext(token string) ([]*K8sContext, error) {
	return p.allContexts, nil
}

// TestFlushMeshSyncDataWithoutContextAllKubeClusterKey verifies that when AllKubeClusterKey
// is absent from the request context (as in DeleteContext / DeleteConnection background teardown),
// FlushMeshSyncData falls back to k8sContext.KubernetesServerID and provider context loading
// to reliably flush MeshSync resources instead of aborting with a false error.
func TestFlushMeshSyncDataWithoutContextAllKubeClusterKey(t *testing.T) {
	serverID := uuid.Must(uuid.NewV4())
	ctxID := uuid.Must(uuid.NewV4()).String()

	t.Run("flushes when KubernetesServerID is present without AllKubeClusterKey in context", func(t *testing.T) {
		db := newMeshSyncTestDB(t)
		seedMeshSyncResource(t, db, "resource-1", serverID.String())

		current := K8sContext{
			ID:                 ctxID,
			Name:               "deleted-cluster",
			Server:             "https://k8s.example.com",
			KubernetesServerID: &serverID,
		}

		fp := &flushDatabaseProvider{db: db}
		// Passing an empty background context without AllKubeClusterKey
		FlushMeshSyncData(context.Background(), current, fp, nil, "", nil, nil)

		var count int64
		if err := db.Model(&meshsyncmodel.KubernetesResource{}).Where("cluster_id = ?", serverID.String()).Count(&count).Error; err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Fatalf("expected 0 resources remaining after flush, got %d", count)
		}
		if len(fp.persistedEvents) != 1 {
			t.Fatalf("expected 1 persisted success event, got %d", len(fp.persistedEvents))
		}
		if fp.persistedEvents[0].Severity != events.Informational {
			t.Errorf("expected Informational event, got %s", fp.persistedEvents[0].Severity)
		}
	})

	t.Run("preserves resources when provider reports multiple active contexts sharing the server", func(t *testing.T) {
		db := newMeshSyncTestDB(t)
		seedMeshSyncResource(t, db, "shared-resource", serverID.String())

		current := K8sContext{
			ID:                 ctxID,
			Name:               "deleted-cluster-1",
			Server:             "https://k8s.example.com",
			KubernetesServerID: &serverID,
		}

		otherCtx := &K8sContext{
			ID:                 uuid.Must(uuid.NewV4()).String(),
			Name:               "active-cluster-2",
			Server:             "https://k8s.example.com",
			KubernetesServerID: &serverID,
		}

		mockProvider := &flushMockContextLoaderProvider{
			flushDatabaseProvider: flushDatabaseProvider{db: db},
			allContexts:           []*K8sContext{&current, otherCtx},
		}

		// Empty context without AllKubeClusterKey
		FlushMeshSyncData(context.Background(), current, mockProvider, nil, "", nil, nil)

		var count int64
		if err := db.Model(&meshsyncmodel.KubernetesResource{}).Where("cluster_id = ?", serverID.String()).Count(&count).Error; err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("expected 1 resource preserved for shared server, got %d", count)
		}
	})
}

func TestK8sContext_RedactCredentials(t *testing.T) {
	ctx := &K8sContext{
		ID:             "test-id",
		Name:           "test-name",
		Server:         "https://cluster.example.com",
		ConnectionID:   "conn-123",
		DeploymentType: "out_cluster",
		Version:        "v1.30.0",
		Auth: sql.Map{
			"client-key-data":         "super-secret-key",
			"token":                   "secret-token",
			"client-certificate-data": "client-cert",
		},
		Cluster: sql.Map{
			"certificate-authority-data": "ca-cert",
			"server":                     "https://cluster.example.com",
		},
	}

	ctx.RedactCredentials()

	if len(ctx.Auth) != 0 {
		t.Errorf("expected Auth to be empty, got: %v", ctx.Auth)
	}
	if len(ctx.Cluster) != 0 {
		t.Errorf("expected Cluster to be empty, got: %v", ctx.Cluster)
	}
	if ctx.ID != "test-id" {
		t.Errorf("expected ID to be preserved, got %s", ctx.ID)
	}
	if ctx.Name != "test-name" {
		t.Errorf("expected Name to be preserved, got %s", ctx.Name)
	}
	if ctx.Server != "https://cluster.example.com" {
		t.Errorf("expected Server to be preserved, got %s", ctx.Server)
	}
	if ctx.ConnectionID != "conn-123" {
		t.Errorf("expected ConnectionID to be preserved, got %s", ctx.ConnectionID)
	}
	if ctx.DeploymentType != "out_cluster" {
		t.Errorf("expected DeploymentType to be preserved, got %s", ctx.DeploymentType)
	}
	if ctx.Version != "v1.30.0" {
		t.Errorf("expected Version to be preserved, got %s", ctx.Version)
	}

	// Nil receiver safety check
	var nilCtx *K8sContext
	nilCtx.RedactCredentials() // must not panic
}

func TestK8sContext_Redacted(t *testing.T) {
	orig := &K8sContext{
		ID:             "ctx-copy-id",
		Name:           "ctx-copy-name",
		Server:         "https://cluster-copy.example.com",
		ConnectionID:   "conn-copy-123",
		DeploymentType: "out_cluster",
		Version:        "v1.30.0",
		Auth: sql.Map{
			"client-key-data": "secret-key",
			"token":           "secret-token",
		},
		Cluster: sql.Map{
			"certificate-authority-data": "secret-ca",
		},
	}

	redacted := orig.Redacted()

	// Redacted copy must have credentials wiped
	if len(redacted.Auth) != 0 {
		t.Errorf("expected redacted Auth to be empty, got: %v", redacted.Auth)
	}
	if len(redacted.Cluster) != 0 {
		t.Errorf("expected redacted Cluster to be empty, got: %v", redacted.Cluster)
	}
	// Redacted copy must preserve metadata
	if redacted.ID != "ctx-copy-id" || redacted.Name != "ctx-copy-name" || redacted.Server != "https://cluster-copy.example.com" {
		t.Errorf("expected metadata to be preserved in copy, got %+v", redacted)
	}
	if redacted.ConnectionID != "conn-copy-123" || redacted.DeploymentType != "out_cluster" || redacted.Version != "v1.30.0" {
		t.Errorf("expected metadata to be preserved in copy, got %+v", redacted)
	}

	// Original must remain untouched (copy semantics)
	if len(orig.Auth) == 0 {
		t.Errorf("expected original Auth to remain untouched, got empty")
	}
	if len(orig.Cluster) == 0 {
		t.Errorf("expected original Cluster to remain untouched, got empty")
	}

	// Nil receiver safety
	var nilCtx *K8sContext
	nilRes := nilCtx.Redacted() // must not panic
	if len(nilRes.Auth) != 0 || len(nilRes.Cluster) != 0 {
		t.Errorf("expected empty maps from nil receiver, got %+v", nilRes)
	}
}

func TestRedactCredentialsForContext_NilSafety(t *testing.T) {
	// Calling with nil must not panic
	res := RedactCredentialsForContext(nil)
	if len(res.Auth) != 0 || len(res.Cluster) != 0 {
		t.Errorf("expected empty maps for nil input, got %+v", res)
	}

	// Calling with populated context
	ctx := &K8sContext{
		ID:     "legacy-id",
		Name:   "legacy-name",
		Server: "https://legacy.example.com",
		Auth: sql.Map{
			"token": "secret",
		},
	}
	resPop := RedactCredentialsForContext(ctx)
	if len(resPop.Auth) != 0 {
		t.Errorf("expected empty Auth, got %+v", resPop.Auth)
	}
	if resPop.MesheryInstanceID != nil {
		t.Errorf("expected nil MesheryInstanceID, got %+v", resPop.MesheryInstanceID)
	}
}

func testKubeLogger(t *testing.T) logger.Handler {
	t.Helper()
	log, err := logger.New("test", logger.Options{Format: logger.JsonLogFormat})
	if err != nil {
		t.Fatalf("failed to build logger: %v", err)
	}
	return log
}

func discoverForTest(t *testing.T, kubeconfig []byte) (KubeconfigDiscovery, error) {
	t.Helper()
	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))
	return DiscoverKubeconfigContexts(context.Background(), nil, uuid.Must(uuid.NewV4()).String(), nil, kubeconfig, &instanceID, map[string]interface{}{}, testKubeLogger(t), true)
}

func TestDiscoverKubeconfigContextsDistinguishesUnreadableEmptyAndSkipped(t *testing.T) {
	t.Run("unparseable input is an error", func(t *testing.T) {
		_, err := discoverForTest(t, []byte("this is not a kubeconfig {{{"))
		if !errors.Is(err, ErrKubeconfigUnreadable) {
			t.Fatalf("err = %v, want ErrKubeconfigUnreadable", err)
		}
	})

	t.Run("yaml that is not a kubeconfig is an error", func(t *testing.T) {
		_, err := discoverForTest(t, []byte("foo: bar\n"))
		if !errors.Is(err, ErrKubeconfigUnreadable) {
			t.Fatalf("err = %v, want ErrKubeconfigUnreadable", err)
		}
	})

	t.Run("a manifest that only has an apiVersion is an error", func(t *testing.T) {
		for _, doc := range []string{
			"apiVersion: v1\n",
			"apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: nginx\n",
		} {
			_, err := discoverForTest(t, []byte(doc))
			if !errors.Is(err, ErrKubeconfigUnreadable) {
				t.Fatalf("doc %q: err = %v, want ErrKubeconfigUnreadable", doc, err)
			}
		}
	})

	t.Run("a kubeconfig with no contexts is empty, not an error", func(t *testing.T) {
		got, err := discoverForTest(t, []byte("apiVersion: v1\nkind: Config\ncontexts: []\nclusters: []\nusers: []\n"))
		if err != nil {
			t.Fatalf("empty kubeconfig: %v", err)
		}
		if len(got.Contexts) != 0 || len(got.SkippedContexts) != 0 {
			t.Fatalf("contexts=%d skipped=%d, want both empty", len(got.Contexts), len(got.SkippedContexts))
		}
	})

	t.Run("exec auth that cannot be loaded is skipped with a reason", func(t *testing.T) {
		got, err := discoverForTest(t, []byte(`apiVersion: v1
kind: Config
current-context: eks
clusters:
- cluster:
    certificate-authority: /no/such/ca.crt
    server: https://127.0.0.1:1
  name: eks
contexts:
- context:
    cluster: eks
    user: eks-user
  name: eks
users:
- name: eks-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: aws
      args: ["eks", "get-token"]
`))
		if err != nil {
			t.Fatalf("discover: %v", err)
		}
		if len(got.Contexts) != 0 {
			t.Fatalf("imported %d contexts, want the exec context skipped", len(got.Contexts))
		}
		if len(got.SkippedContexts) != 1 {
			t.Fatalf("skipped %d, want 1", len(got.SkippedContexts))
		}
		skipped := got.SkippedContexts[0]
		if skipped.Name != "eks" {
			t.Fatalf("skipped name = %q, want eks", skipped.Name)
		}
		if !strings.Contains(skipped.Reason, "exec plugin") {
			t.Fatalf("reason = %q, want it to name the exec plugin", skipped.Reason)
		}
	})

	t.Run("auth-provider contexts are skipped with a reason", func(t *testing.T) {
		got, err := discoverForTest(t, []byte(`apiVersion: v1
kind: Config
current-context: gke
clusters:
- cluster:
    server: https://127.0.0.1:1
  name: gke
contexts:
- context:
    cluster: gke
    user: gke
  name: gke
users:
- name: gke
  user:
    auth-provider:
      name: gcp
`))
		if err != nil {
			t.Fatalf("discover: %v", err)
		}
		if len(got.Contexts) != 0 || len(got.SkippedContexts) != 1 {
			t.Fatalf("contexts=%d skipped=%d, want 0 and 1", len(got.Contexts), len(got.SkippedContexts))
		}
		if !strings.Contains(got.SkippedContexts[0].Reason, "auth-provider") {
			t.Fatalf("reason = %q, want it to name auth-provider", got.SkippedContexts[0].Reason)
		}
	})

	t.Run("a readable token context is returned and a skipped exec context is reported beside it", func(t *testing.T) {
		got, err := discoverForTest(t, []byte(`apiVersion: v1
kind: Config
current-context: good
clusters:
- cluster:
    server: https://127.0.0.1:1
  name: good
- cluster:
    certificate-authority: /no/such/ca.crt
    server: https://127.0.0.1:1
  name: bad
contexts:
- context:
    cluster: good
    user: good
  name: good
- context:
    cluster: bad
    user: bad
  name: bad
users:
- name: good
  user:
    token: token-good
- name: bad
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: aws
      args: ["eks", "get-token"]
`))
		if err != nil {
			t.Fatalf("discover: %v", err)
		}
		if len(got.Contexts) != 1 || got.Contexts[0].Name != "good" {
			t.Fatalf("contexts = %#v, want only good", got.Contexts)
		}
		if len(got.SkippedContexts) != 1 || got.SkippedContexts[0].Name != "bad" {
			t.Fatalf("skipped = %#v, want bad", got.SkippedContexts)
		}
	})

	t.Run("a well-formed exec kubeconfig is skipped before the API probe", func(t *testing.T) {
		got, err := discoverForTest(t, []byte(`apiVersion: v1
kind: Config
current-context: eks
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCg==
    server: https://127.0.0.1:1
  name: eks
contexts:
- context:
    cluster: eks
    user: eks-user
  name: eks
users:
- name: eks-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: aws
      args: ["eks", "get-token"]
`))
		if err != nil {
			t.Fatalf("discover: %v", err)
		}
		if len(got.Contexts) != 0 {
			t.Fatalf("imported %d contexts, want the exec context skipped even when a client can be built", len(got.Contexts))
		}
		if len(got.SkippedContexts) != 1 {
			t.Fatalf("skipped %d, want 1", len(got.SkippedContexts))
		}
		skipped := got.SkippedContexts[0]
		if skipped.Name != "eks" {
			t.Fatalf("skipped name = %q, want eks", skipped.Name)
		}
		if !strings.Contains(skipped.Reason, "exec plugin") {
			t.Fatalf("reason = %q, want it to name the exec plugin", skipped.Reason)
		}
		if !strings.Contains(skipped.Reason, "Meshery Server host") {
			t.Fatalf("reason = %q, want it to say the plugin would run on the Meshery Server host", skipped.Reason)
		}
	})
}

// Existing callers (import, component registration) keep receiving a context
// slice. An unreadable file stays an empty slice for them; the discovery
// endpoint is what turns that into an error.
func TestK8sContextsFromKubeconfigWithOptionsKeepsUnreadableAsEmpty(t *testing.T) {
	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))
	got := K8sContextsFromKubeconfigWithOptions(nil, uuid.Must(uuid.NewV4()).String(), nil, []byte("not a kubeconfig {{{"), &instanceID, map[string]interface{}{}, testKubeLogger(t), true)
	if len(got) != 0 {
		t.Fatalf("got %d contexts, want 0", len(got))
	}
}

func TestK8sContextImportErrorIsOptionalJSON(t *testing.T) {
	withReason, err := json.Marshal(K8sContext{Name: "bad", ImportError: "save failed"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(withReason), `"error":"save failed"`) {
		t.Fatalf("body = %s, want an error reason", withReason)
	}

	without, err := json.Marshal(K8sContext{Name: "ok"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(without), `"error"`) {
		t.Fatalf("body = %s, want error omitted when empty", without)
	}
}
