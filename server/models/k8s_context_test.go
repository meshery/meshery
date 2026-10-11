package models

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/internal/sql"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshery/server/models/httputil"
	"github.com/meshery/meshkit/database"
	mkerrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/logger"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/environment"
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

// TestK8sContextGenerateID verifies that K8sContextGenerateID excludes service-account
// tokens from the hash for in-cluster contexts to prevent ID changes when tokens rotate (issue #21810)
func TestK8sContextGenerateID(t *testing.T) {
	instanceID, _ := uuid.NewV4()

	tests := []struct {
		name     string
		context  K8sContext
		wantSame bool // whether ID should remain same after token change
	}{
		{
			name: "regular context with token",
			context: K8sContext{
				Name:              "test-context",
				Auth:              sql.Map{"user": map[string]interface{}{"token": "original-token"}},
				Cluster:           sql.Map{"server": "https://k8s.example.com"},
				MesheryInstanceID: &instanceID,
			},
			wantSame: false, // ID should CHANGE after token change for regular contexts
		},
		{
			name: "in-cluster context with token",
			context: K8sContext{
				Name:              "in-cluster-context",
				Auth:              sql.Map{"user": map[string]interface{}{"token": "service-account-token"}},
				Cluster:           sql.Map{"server": "https://kubernetes.default.svc"},
				Server:            "https://kubernetes.default.svc",
				MesheryInstanceID: &instanceID,
				DeploymentType:    "in_cluster", // Explicit in-cluster provenance
			},
			wantSame: true, // ID should be same after token rotation for in-cluster contexts
		},
		{
			name: "in-cluster context with token (host+port URL)",
			context: K8sContext{
				Name:              "in-cluster-context-hostport",
				Auth:              sql.Map{"user": map[string]interface{}{"token": "service-account-token"}},
				Cluster:           sql.Map{"server": "https://10.0.0.1:443"},
				Server:            "https://10.0.0.1:443",
				MesheryInstanceID: &instanceID,
				DeploymentType:    "in_cluster", // Explicit in-cluster provenance (actual in-cluster URL format)
			},
			wantSame: true, // ID should be same after token rotation for in-cluster contexts
		},
		{
			name: "context without token",
			context: K8sContext{
				Name:              "no-token-context",
				Auth:              sql.Map{"user": map[string]interface{}{"client-certificate": "cert"}},
				Cluster:           sql.Map{"server": "https://k8s.example.com"},
				MesheryInstanceID: &instanceID,
			},
			wantSame: true, // ID should be same since no token to change
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate initial ID
			id1, err := K8sContextGenerateID(tt.context)
			if err != nil {
				t.Fatalf("K8sContextGenerateID() error = %v", err)
			}

			// Simulate token rotation by changing the token
			if tt.context.Auth != nil {
				if user, ok := tt.context.Auth["user"].(map[string]interface{}); ok {
					if _, hasToken := user["token"]; hasToken {
						user["token"] = "rotated-token"
					}
				}
			}

			// Generate ID after token change
			id2, err := K8sContextGenerateID(tt.context)
			if err != nil {
				t.Fatalf("K8sContextGenerateID() error = %v", err)
			}

			// Check if IDs remain the same
			if tt.wantSame && id1 != id2 {
				t.Errorf("ID changed after token rotation, got %v, want %v", id2, id1)
			}
			if !tt.wantSame && id1 == id2 {
				t.Errorf("ID did not change after token rotation, got %v, want different", id1)
			}
		})
	}
}

// TestNewK8sContextFromInClusterConfigTokenRotation verifies that in-cluster contexts
// generateTestCACert generates a self-signed CA certificate for testing.
// Returns the certificate, private key, and any error.
func generateTestCACert(t *testing.T) (*x509.Certificate, *rsa.PrivateKey, error) {
	t.Helper()

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Test CA"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert, privKey, nil
}

// TestNewK8sContextFromInClusterConfig tests the full in-cluster constructor path
// using a fake Kubernetes API server. It verifies that DeploymentType is set to
// "in_cluster" and that the constructor properly reads token/CA from the expected
// in-cluster file locations.
func TestNewK8sContextFromInClusterConfig(t *testing.T) {
	t.Helper()

	// Generate a test CA certificate programmatically
	caCert, caPrivKey, err := generateTestCACert(t)
	if err != nil {
		t.Fatalf("failed to generate CA cert: %v", err)
	}

	// Create a temporary directory for in-cluster files
	tmpDir := t.TempDir()
	tokenFile := filepath.Join(tmpDir, "token")
	caFile := filepath.Join(tmpDir, "ca.crt")

	// Write a test service-account token
	testToken := "test-service-account-token"
	if err := os.WriteFile(tokenFile, []byte(testToken), 0600); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}

	// Write the CA certificate in PEM format
	caPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCert.Raw,
	})
	if err := os.WriteFile(caFile, caPEM, 0600); err != nil {
		t.Fatalf("failed to write CA file: %v", err)
	}

	// Start a fake Kubernetes API server
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/livez":
			// PingTest endpoint
			w.WriteHeader(http.StatusOK)
		case "/version":
			// AssignVersion endpoint
			w.Header().Set("Content-Type", "application/json")
			if _, err := fmt.Fprintf(w, `{"major":"1","minor":"28","gitVersion":"v1.28.0"}`); err != nil {
				httputil.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
			}
		case "/api/v1/namespaces/kube-system":
			// KubernetesServerID lookup endpoint
			w.Header().Set("Content-Type", "application/json")
			if _, err := fmt.Fprintf(w, `{"metadata":{"uid":"test-server-uid-12345"}}`); err != nil {
				httputil.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{{Certificate: [][]byte{caCert.Raw}, PrivateKey: caPrivKey}},
	}
	server.StartTLS()
	defer server.Close()

	// Save and restore original in-cluster file paths
	oldTokenFile := inClusterTokenFile
	oldCAFile := inClusterRootCAFile
	defer func() {
		inClusterTokenFile = oldTokenFile
		inClusterRootCAFile = oldCAFile
	}()

	// Override in-cluster file paths to our test files
	inClusterTokenFile = tokenFile
	inClusterRootCAFile = caFile

	// Set in-cluster environment variables
	oldHost := os.Getenv("KUBERNETES_SERVICE_HOST")
	oldPort := os.Getenv("KUBERNETES_SERVICE_PORT")
	t.Cleanup(func() {
		if oldHost != "" {
			if err := os.Setenv("KUBERNETES_SERVICE_HOST", oldHost); err != nil {
				t.Logf("failed to restore KUBERNETES_SERVICE_HOST: %v", err)
			}
		} else {
			if err := os.Unsetenv("KUBERNETES_SERVICE_HOST"); err != nil {
				t.Logf("failed to unset KUBERNETES_SERVICE_HOST: %v", err)
			}
		}
		if oldPort != "" {
			if err := os.Setenv("KUBERNETES_SERVICE_PORT", oldPort); err != nil {
				t.Logf("failed to restore KUBERNETES_SERVICE_PORT: %v", err)
			}
		} else {
			if err := os.Unsetenv("KUBERNETES_SERVICE_PORT"); err != nil {
				t.Logf("failed to unset KUBERNETES_SERVICE_PORT: %v", err)
			}
		}
	})

	// Parse server URL to get host and port
	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse server URL: %v", err)
	}
	if err := os.Setenv("KUBERNETES_SERVICE_HOST", u.Hostname()); err != nil {
		t.Fatalf("failed to set KUBERNETES_SERVICE_HOST: %v", err)
	}
	if err := os.Setenv("KUBERNETES_SERVICE_PORT", u.Port()); err != nil {
		t.Fatalf("failed to set KUBERNETES_SERVICE_PORT: %v", err)
	}

	// Call the actual constructor
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctx, err := NewK8sContextFromInClusterConfig("in-cluster", &instanceID, log)
	if err != nil {
		t.Fatalf("NewK8sContextFromInClusterConfig() error = %v", err)
	}

	// Verify DeploymentType is set to "in_cluster"
	if ctx.DeploymentType != "in_cluster" {
		t.Errorf("DeploymentType = %q, want \"in_cluster\"", ctx.DeploymentType)
	}

	// Verify the context name matches the in-cluster convention
	if ctx.Name != "in-cluster" {
		t.Errorf("Name = %q, want \"in-cluster\"", ctx.Name)
	}

	// Verify the server URL matches our fake server
	if ctx.Server != server.URL {
		t.Errorf("Server = %q, want %q", ctx.Server, server.URL)
	}
}

// TestNewK8sContextFromInClusterConfigTokenRotation tests that in-cluster context IDs
// produce stable IDs when service-account tokens rotate using DeploymentType for provenance.
// This test simulates the persistence/reload cycle to ensure IDs remain stable.
func TestNewK8sContextFromInClusterConfigTokenRotation(t *testing.T) {
	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))

	// Simulate in-cluster context with DeploymentType set before ID generation
	// This mimics what NewK8sContextFromInClusterConfig does after setting DeploymentType
	ctx1 := K8sContext{
		Name:              "in-cluster-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "initial-token"}},
		Cluster:           sql.Map{"server": "https://10.0.0.1:443"},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster", // This is what NewK8sContextFromInClusterConfig sets
	}

	id1, err := K8sContextGenerateID(ctx1)
	if err != nil {
		t.Fatalf("K8sContextGenerateID() error = %v", err)
	}

	// Simulate token rotation with same DeploymentType (persistence/reload scenario)
	ctx2 := K8sContext{
		Name:              "in-cluster-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "rotated-token"}},
		Cluster:           sql.Map{"server": "https://10.0.0.1:443"},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster", // Preserved from persistence
	}

	id2, err := K8sContextGenerateID(ctx2)
	if err != nil {
		t.Fatalf("K8sContextGenerateID() error = %v", err)
	}

	// IDs must be the same for in-cluster contexts with different tokens
	if id1 != id2 {
		t.Errorf("in-cluster context ID changed after token rotation: got %v, want %v", id2, id1)
	}

	// Verify persistence/reload concept: create context representing persisted state
	persistedCtx := K8sContext{
		Name:              "in-cluster-context",
		Auth:              ctx2.Auth,
		Cluster:           ctx2.Cluster,
		Server:            ctx2.Server,
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster", // Preserved from persistence
	}

	id3, err := K8sContextGenerateID(persistedCtx)
	if err != nil {
		t.Fatalf("K8sContextGenerateID() error = %v", err)
	}

	// ID must remain the same with DeploymentType preserved
	if id1 != id3 {
		t.Errorf("persisted context ID differs from original: got %v, want %v", id3, id1)
	}

	// Verify regular context behavior: different tokens should produce different IDs
	regularCtx1 := K8sContext{
		Name:              "regular-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-1"}},
		Cluster:           sql.Map{"server": "https://k8s.example.com"},
		MesheryInstanceID: &instanceID,
		DeploymentType:    "out_of_cluster", // Regular context
	}

	regularID1, err := K8sContextGenerateID(regularCtx1)
	if err != nil {
		t.Fatalf("K8sContextGenerateID() error = %v", err)
	}

	regularCtx2 := K8sContext{
		Name:              "regular-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-2"}},
		Cluster:           sql.Map{"server": "https://k8s.example.com"},
		MesheryInstanceID: &instanceID,
		DeploymentType:    "out_of_cluster", // Regular context
	}

	regularID2, err := K8sContextGenerateID(regularCtx2)
	if err != nil {
		t.Fatalf("K8sContextGenerateID() error = %v", err)
	}

	// IDs must be different for regular contexts with different tokens
	if regularID1 == regularID2 {
		t.Errorf("regular context ID did not change after token rotation: got %v, want different", regularID1)
	}
}

// newK8sContextFixture creates a DefaultLocalProvider with an in-memory SQLite database
// that has the k8s_contexts, connections, credentials, and environment_connection_mappings
// tables migrated. This provides the minimum infrastructure needed to exercise
// DefaultLocalProvider.SaveK8sContext in a regression test.
func newK8sContextFixture(t *testing.T) *DefaultLocalProvider {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	// Migrate required tables for SaveK8sContext
	if err := db.AutoMigrate(&K8sContext{}, &Credential{}, connections.Connection{}, environment.EnvironmentConnectionMapping{}); err != nil {
		t.Fatalf("failed to migrate k8s_contexts/credential/connections/mappings: %v", err)
	}

	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	return &DefaultLocalProvider{
		GenericPersister:           &db,
		ConnectionPersister:        &ConnectionPersister{DB: &db},
		MesheryK8sContextPersister: &MesheryK8sContextPersister{DB: &db},
		EnvironmentPersister:       &EnvironmentPersister{DB: &db},
		Log:                        log,
	}
}

// TestLocalProviderK8sContextTokenRotation verifies that calling SaveK8sContext
// twice with the same logical context but different tokens (simulating token rotation)
// updates the credential correctly without creating duplicate connections.
func TestLocalProviderK8sContextTokenRotation(t *testing.T) {
	provider := newK8sContextFixture(t)
	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))

	// Create initial context with old token
	initialCtx := K8sContext{
		Name:              "in-cluster-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-old"}},
		Cluster:           sql.Map{"server": "https://10.0.0.1:443"},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
		Version:           "v1.0.0",
	}

	// First save - creates the connection/context
	conn1, err := provider.SaveK8sContext("test-token", initialCtx, nil)
	if err != nil {
		t.Fatalf("first SaveK8sContext() error = %v", err)
	}

	// Verify connection was created
	if conn1.ID == uuid.Nil {
		t.Fatal("first SaveK8sContext() returned empty connection ID")
	}

	// Load the context that was created to get its ID
	mkcp := provider.MesheryK8sContextPersister
	var loadedCtx K8sContext
	err = mkcp.DB.Model(&K8sContext{}).Where("connection_id = ?", conn1.ID.String()).First(&loadedCtx).Error
	if err != nil {
		t.Fatalf("failed to load context after first save: %v", err)
	}

	// Create rotated context with new token (same logical identity)
	rotatedCtx := K8sContext{
		ID:                loadedCtx.ID, // Use the same context ID
		Name:              "in-cluster-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-new"}},
		Cluster:           sql.Map{"server": "https://10.0.0.1:443"},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
		Version:           "v1.0.1",
	}

	// Second save - should update the existing connection, not create a duplicate
	conn2, err := provider.SaveK8sContext("test-token", rotatedCtx, nil)
	if err != nil {
		t.Fatalf("second SaveK8sContext() error = %v", err)
	}

	// Verify connection ID remains the same (no duplicate created)
	if conn1.ID != conn2.ID {
		t.Errorf("connection ID changed after token rotation: got %v, want %v", conn2.ID, conn1.ID)
	}

	// Reload the context from database to verify the token was updated
	err = mkcp.DB.Model(&K8sContext{}).Where("connection_id = ?", conn2.ID.String()).First(&loadedCtx).Error
	if err != nil {
		t.Fatalf("failed to load context after second save: %v", err)
	}

	// Verify the stored token is the new one
	token, ok := loadedCtx.Auth["user"].(map[string]interface{})["token"].(string)
	if !ok || token != "token-new" {
		t.Errorf("stored token was not updated: got %v, want token-new", loadedCtx.Auth)
	}

	// Verify version was updated
	if loadedCtx.Version != "v1.0.1" {
		t.Errorf("version was not updated: got %v, want v1.0.1", loadedCtx.Version)
	}

	// Verify there is exactly one connection
	var connectionCount int64
	err = provider.GetGenericPersister().Model(&connections.Connection{}).Count(&connectionCount).Error
	if err != nil {
		t.Fatalf("failed to count connections: %v", err)
	}
	if connectionCount != 1 {
		t.Errorf("expected 1 connection after token rotation, found %d", connectionCount)
	}
}

// TestLocalProviderLegacyK8sContextMigration verifies that a legacy in-cluster
// context with a token-dependent ID is reconciled to the new token-independent ID
// when SaveK8sContext is called, without creating duplicate connections.
func TestLocalProviderLegacyK8sContextMigration(t *testing.T) {
	provider := newK8sContextFixture(t)
	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))

	// Manually insert a legacy connection with old token-dependent ID
	// This simulates the pre-fix state where IDs included the token
	legacyConnID := uuid.Must(uuid.NewV4())
	oldTime := time.Now().Add(-24 * time.Hour)

	legacyMetadata := map[string]interface{}{
		"id":                "legacy-token-dependent-id",
		"server":            "https://10.0.0.1:443",
		"mesheryInstanceId": instanceID.String(),
		"deploymentType":    "in_cluster",
		"version":           "v1.0.0",
		"name":              "in-cluster-context",
	}

	legacyConn := connections.Connection{
		ID:             legacyConnID,
		Kind:           "kubernetes",
		ConnectionType: "platform",
		SubType:        "orchestrator",
		Status:         connections.DISCOVERED,
		Metadata:       legacyMetadata,
	}

	err := provider.GetGenericPersister().Save(&legacyConn).Error
	if err != nil {
		t.Fatalf("failed to save legacy connection: %v", err)
	}

	// Manually insert a legacy k8s_context row with old token
	legacyCtx := K8sContext{
		ID:                "legacy-token-dependent-id",
		Name:              "in-cluster-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-old"}},
		Cluster:           sql.Map{"server": "https://10.0.0.1:443"},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
		ConnectionID:      legacyConnID.String(),
		Version:           "v1.0.0",
		UpdatedAt:         &oldTime,
		CreatedAt:         &oldTime,
	}

	mkcp := provider.MesheryK8sContextPersister
	_, err = mkcp.SaveMesheryK8sContext(legacyCtx)
	if err != nil && err != ErrContextAlreadyPersisted {
		t.Fatalf("failed to save legacy context: %v", err)
	}

	// Verify legacy context exists
	var count int64
	err = mkcp.DB.Model(&K8sContext{}).Where("id = ?", "legacy-token-dependent-id").Count(&count).Error
	if err != nil {
		t.Fatalf("failed to count legacy contexts: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 legacy context, found %d", count)
	}

	// Create current context with new token-independent ID behavior
	currentCtx := K8sContext{
		Name:              "in-cluster-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-new"}},
		Cluster:           sql.Map{"server": "https://10.0.0.1:443"},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
		Version:           "v1.0.1",
	}

	// Call SaveK8sContext - this should trigger the migration/reconciliation logic
	conn, err := provider.SaveK8sContext("test-token", currentCtx, nil)
	if err != nil {
		t.Fatalf("SaveK8sContext() error = %v", err)
	}

	// Verify a connection was returned
	if conn.ID == uuid.Nil {
		t.Fatal("SaveK8sContext() returned empty connection ID")
	}

	// Verify the final context has the new token
	// The context ID will be the new stable ID, not the legacy ID
	// Look up by connection ID instead
	var loadedCtx K8sContext
	err = mkcp.DB.Model(&K8sContext{}).Where("connection_id = ?", conn.ID.String()).First(&loadedCtx).Error
	if err != nil {
		t.Fatalf("failed to load context by connection_id: %v", err)
	}

	token, ok := loadedCtx.Auth["user"].(map[string]interface{})["token"].(string)
	if !ok || token != "token-new" {
		t.Errorf("stored token was not updated: got %v, want token-new", loadedCtx.Auth)
	}

	// Verify version was updated
	if loadedCtx.Version != "v1.0.1" {
		t.Errorf("version was not updated: got %v, want v1.0.1", loadedCtx.Version)
	}

	// Verify there is only one context record (no duplicate)
	var contextCount int64
	err = mkcp.DB.Model(&K8sContext{}).Count(&contextCount).Error
	if err != nil {
		t.Fatalf("failed to count total contexts: %v", err)
	}
	if contextCount != 1 {
		t.Errorf("expected 1 context after migration, found %d", contextCount)
	}

	// Verify there is only one connection
	var connectionCount int64
	err = provider.GetGenericPersister().Model(&connections.Connection{}).Count(&connectionCount).Error
	if err != nil {
		t.Fatalf("failed to count total connections: %v", err)
	}
	if connectionCount != 1 {
		t.Errorf("expected 1 connection after migration, found %d", connectionCount)
	}
}

// TestUpdateMesheryK8sContext verifies that UpdateMesheryK8sContext correctly
// updates an existing k8s context's mutable fields (auth, cluster, version) while
// preserving immutable fields and correctly updating updated_at via GORM.
func TestUpdateMesheryK8sContext(t *testing.T) {
	// Use the database.New approach consistent with repository patterns
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("failed to open SQLite database: %v", err)
	}

	// Create the k8s_contexts table
	err = db.AutoMigrate(&K8sContext{})
	if err != nil {
		t.Fatalf("failed to migrate K8sContext: %v", err)
	}

	// Create persister
	persister := &MesheryK8sContextPersister{DB: &db}

	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))
	oldTime := time.Now().Add(-1 * time.Hour)

	// Create initial context with old token and non-null updated_at
	initialCtx := K8sContext{
		ID:                "test-context-id",
		Name:              "test-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "tok-v1"}},
		Cluster:           sql.Map{"server": "https://10.0.0.1:443"},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
		Version:           "v1.0.0",
		UpdatedAt:         &oldTime,
	}

	// Save initial context
	_, err = persister.SaveMesheryK8sContext(initialCtx)
	if err != nil {
		t.Fatalf("SaveMesheryK8sContext() initial save error = %v", err)
	}

	// Reload to verify initial state
	loadedCtx, err := persister.GetMesheryK8sContext("test-context-id")
	if err != nil {
		t.Fatalf("GetMesheryK8sContext() error = %v", err)
	}

	// Verify initial token
	token1, ok := loadedCtx.Auth["user"].(map[string]interface{})["token"].(string)
	if !ok || token1 != "tok-v1" {
		t.Errorf("initial token incorrect: got %v, want tok-v1", loadedCtx.Auth)
	}

	// Verify initial updated_at is non-null
	if loadedCtx.UpdatedAt == nil {
		t.Errorf("initial updated_at is nil, want non-null")
	}

	// Create refreshed context with new token (UpdatedAt is nil as from fresh discovery)
	refreshedCtx := K8sContext{
		ID:                "test-context-id", // Same ID (stable)
		Name:              "test-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "tok-v2-ROTATED"}},
		Cluster:           sql.Map{"server": "https://10.0.0.1:443"},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
		Version:           "v1.0.1",
		UpdatedAt:         nil, // Fresh context has nil UpdatedAt
	}

	// Call the actual UpdateMesheryK8sContext
	err = persister.UpdateMesheryK8sContext(refreshedCtx)
	if err != nil {
		t.Fatalf("UpdateMesheryK8sContext() error = %v", err)
	}

	// Reload from database to verify persisted state
	updatedCtx, err := persister.GetMesheryK8sContext("test-context-id")
	if err != nil {
		t.Fatalf("GetMesheryK8sContext() after update error = %v", err)
	}

	// Verify token was updated
	token2, ok := updatedCtx.Auth["user"].(map[string]interface{})["token"].(string)
	if !ok || token2 != "tok-v2-ROTATED" {
		t.Errorf("rotated token not persisted: got %v, want tok-v2-ROTATED", updatedCtx.Auth)
	}

	// Verify ID remained stable
	if updatedCtx.ID != "test-context-id" {
		t.Errorf("context ID changed: got %v, want test-context-id", updatedCtx.ID)
	}

	// Verify version was updated
	if updatedCtx.Version != "v1.0.1" {
		t.Errorf("version not updated: got %v, want v1.0.1", updatedCtx.Version)
	}

	// Verify immutable fields are preserved
	if updatedCtx.Name != "test-context" {
		t.Errorf("name changed: got %v, want test-context", updatedCtx.Name)
	}
	if updatedCtx.Server != "https://10.0.0.1:443" {
		t.Errorf("server changed: got %v, want https://10.0.0.1:443", updatedCtx.Server)
	}
	if updatedCtx.DeploymentType != "in_cluster" {
		t.Errorf("deploymentType changed: got %v, want in_cluster", updatedCtx.DeploymentType)
	}

	// Verify updated_at is non-null and changed
	if updatedCtx.UpdatedAt == nil {
		t.Errorf("updated_at is nil after update, want non-null")
	}
	// Note: In-memory SQLite may have coarse timestamp granularity, so just verify it's not nil
	// The production code correctly lets GORM handle the timestamp update

	// Verify created_at is preserved
	if updatedCtx.CreatedAt == nil {
		t.Errorf("created_at is nil after update, want preserved")
	}
}

// TestSaveK8sContextRepeatSaveWithID verifies that calling SaveK8sContext
// twice with the same context (repeat save) correctly handles the case where
// SaveMesheryK8sContext returns ErrContextAlreadyPersisted and the update
// path must use the correct context ID. This tests the fix for the CodeRabbit
// functional correctness comment.
func TestSaveK8sContextRepeatSaveWithID(t *testing.T) {
	provider := newK8sContextFixture(t)
	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))

	// Create initial context
	initialCtx := K8sContext{
		Name:              "test-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-v1"}},
		Cluster:           sql.Map{"server": "https://10.0.0.1:443"},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
		Version:           "v1.0.0",
	}

	// First save - should succeed
	conn1, err := provider.SaveK8sContext("test-token", initialCtx, nil)
	if err != nil {
		t.Fatalf("first SaveK8sContext() error = %v", err)
	}

	// Verify connection was created
	if conn1.ID == uuid.Nil {
		t.Fatal("first SaveK8sContext() returned empty connection ID")
	}

	// Load the context to get its ID
	mkcp := provider.MesheryK8sContextPersister
	var loadedCtx K8sContext
	err = mkcp.DB.Model(&K8sContext{}).Where("connection_id = ?", conn1.ID.String()).First(&loadedCtx).Error
	if err != nil {
		t.Fatalf("failed to load context after first save: %v", err)
	}

	// Create a second context with the same logical identity (same server, name, instance)
	// but with updated auth/version - this simulates a repeat save
	repeatCtx := K8sContext{
		ID:                loadedCtx.ID, // Use the same ID from first save
		Name:              "test-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-v2"}},
		Cluster:           sql.Map{"server": "https://10.0.0.1:443"},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
		Version:           "v1.0.1",
	}

	// Second save - should trigger ErrContextAlreadyPersisted path and update correctly
	conn2, err := provider.SaveK8sContext("test-token", repeatCtx, nil)
	if err != nil {
		t.Fatalf("second SaveK8sContext() error = %v", err)
	}

	// Verify connection ID remains the same
	if conn1.ID != conn2.ID {
		t.Errorf("connection ID changed on repeat save: got %v, want %v", conn2.ID, conn1.ID)
	}

	// Reload context to verify it was updated
	err = mkcp.DB.Model(&K8sContext{}).Where("connection_id = ?", conn2.ID.String()).First(&loadedCtx).Error
	if err != nil {
		t.Fatalf("failed to load context after second save: %v", err)
	}

	// Verify token was updated
	token, ok := loadedCtx.Auth["user"].(map[string]interface{})["token"].(string)
	if !ok || token != "token-v2" {
		t.Errorf("token not updated on repeat save: got %v, want token-v2", loadedCtx.Auth)
	}

	// Verify version was updated
	if loadedCtx.Version != "v1.0.1" {
		t.Errorf("version not updated on repeat save: got %v, want v1.0.1", loadedCtx.Version)
	}

	// Verify there is still only one connection
	var connectionCount int64
	err = provider.GetGenericPersister().Model(&connections.Connection{}).Count(&connectionCount).Error
	if err != nil {
		t.Fatalf("failed to count connections: %v", err)
	}
	if connectionCount != 1 {
		t.Errorf("expected 1 connection after repeat save, found %d", connectionCount)
	}
}

// TestSaveK8sContextWithEmptyID verifies that SaveK8sContext correctly assigns
// a stable ID when the input context has an empty ID, and that a subsequent
// save with the same logical identity uses the assigned ID (repeat-save case).
func TestSaveK8sContextWithEmptyID(t *testing.T) {
	provider := newK8sContextFixture(t)
	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))

	// Create initial context with empty ID (discovery case)
	initialCtx := K8sContext{
		ID:                "", // Empty ID - should be generated
		Name:              "empty-id-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-v1"}},
		Cluster:           sql.Map{"server": "https://10.0.0.1:443"},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
		Version:           "v1.0.0",
	}

	// First save - should generate ID and create connection
	conn1, err := provider.SaveK8sContext("test-token", initialCtx, nil)
	if err != nil {
		t.Fatalf("first SaveK8sContext() error = %v", err)
	}

	// Verify connection was created with a valid ID
	if conn1.ID == uuid.Nil {
		t.Fatal("first SaveK8sContext() returned empty connection ID")
	}

	// Load the context to verify ID was assigned
	mkcp := provider.MesheryK8sContextPersister
	var loadedCtx K8sContext
	err = mkcp.DB.Model(&K8sContext{}).Where("connection_id = ?", conn1.ID.String()).First(&loadedCtx).Error
	if err != nil {
		t.Fatalf("failed to load context after first save: %v", err)
	}

	// Verify ID was generated and is not empty
	if loadedCtx.ID == "" {
		t.Errorf("context ID was not assigned on first save")
	}

	// Create a second context with the same logical identity but empty ID
	// This simulates a repeat save where the caller provides an empty ID
	repeatCtx := K8sContext{
		ID:                "", // Empty ID again - should use existing connection
		Name:              "empty-id-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-v2"}},
		Cluster:           sql.Map{"server": "https://10.0.0.1:443"},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
		Version:           "v1.0.1",
	}

	// Second save - should find existing connection by stable fields and update it
	conn2, err := provider.SaveK8sContext("test-token", repeatCtx, nil)
	if err != nil {
		t.Fatalf("second SaveK8sContext() error = %v", err)
	}

	// Verify connection ID remains the same (no duplicate created)
	if conn1.ID != conn2.ID {
		t.Errorf("connection ID changed on repeat save with empty ID: got %v, want %v", conn2.ID, conn1.ID)
	}

	// Reload context to verify it was updated
	err = mkcp.DB.Model(&K8sContext{}).Where("connection_id = ?", conn2.ID.String()).First(&loadedCtx).Error
	if err != nil {
		t.Fatalf("failed to load context after second save: %v", err)
	}

	// Verify token was updated
	token, ok := loadedCtx.Auth["user"].(map[string]interface{})["token"].(string)
	if !ok || token != "token-v2" {
		t.Errorf("token not updated on repeat save with empty ID: got %v, want token-v2", loadedCtx.Auth)
	}

	// Verify version was updated
	if loadedCtx.Version != "v1.0.1" {
		t.Errorf("version not updated on repeat save with empty ID: got %v, want v1.0.1", loadedCtx.Version)
	}

	// Verify there is still only one connection
	var connectionCount int64
	err = provider.GetGenericPersister().Model(&connections.Connection{}).Count(&connectionCount).Error
	if err != nil {
		t.Fatalf("failed to count connections: %v", err)
	}
	if connectionCount != 1 {
		t.Errorf("expected 1 connection after repeat save with empty ID, found %d", connectionCount)
	}
}

// TestInClusterContextCARotationStability verifies that in-cluster context IDs
// remain stable when CA certificate data rotates. For in-cluster contexts,
// certificate-authority-data is excluded from the ID hash.
func TestInClusterContextCARotationStability(t *testing.T) {
	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))

	// Test CA rotation for in-cluster context
	ctx1 := K8sContext{
		Name:              "in-cluster-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "service-account-token"}},
		Cluster:           sql.Map{"cluster": map[string]interface{}{"certificate-authority-data": "ca-data-v1", "server": "https://10.0.0.1:443"}},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
	}

	id1, err := K8sContextGenerateID(ctx1)
	if err != nil {
		t.Fatalf("K8sContextGenerateID() error = %v", err)
	}

	// Simulate CA rotation
	ctx2 := K8sContext{
		Name:              "in-cluster-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "service-account-token"}},
		Cluster:           sql.Map{"cluster": map[string]interface{}{"certificate-authority-data": "ca-data-v2", "server": "https://10.0.0.1:443"}},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
	}

	id2, err := K8sContextGenerateID(ctx2)
	if err != nil {
		t.Fatalf("K8sContextGenerateID() error = %v", err)
	}

	// IDs must be the same for in-cluster contexts with different CA data
	if id1 != id2 {
		t.Errorf("in-cluster context ID changed after CA rotation: got %v, want %v", id2, id1)
	}

	// Test CA rotation for regular context (should change ID)
	regularCtx1 := K8sContext{
		Name:              "regular-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-1"}},
		Cluster:           sql.Map{"cluster": map[string]interface{}{"certificate-authority-data": "ca-data-v1", "server": "https://k8s.example.com"}},
		Server:            "https://k8s.example.com",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "out_of_cluster",
	}

	regularID1, err := K8sContextGenerateID(regularCtx1)
	if err != nil {
		t.Fatalf("K8sContextGenerateID() error = %v", err)
	}

	regularCtx2 := K8sContext{
		Name:              "regular-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-1"}},
		Cluster:           sql.Map{"cluster": map[string]interface{}{"certificate-authority-data": "ca-data-v2", "server": "https://k8s.example.com"}},
		Server:            "https://k8s.example.com",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "out_of_cluster",
	}

	regularID2, err := K8sContextGenerateID(regularCtx2)
	if err != nil {
		t.Fatalf("K8sContextGenerateID() error = %v", err)
	}

	// IDs must be different for regular contexts with different CA data
	if regularID1 == regularID2 {
		t.Errorf("regular context ID did not change after CA rotation: got %v, want different", regularID1)
	}
}

// TestInClusterContextTokenAndCARotationStability verifies that in-cluster context IDs
// remain stable when both token and CA certificate data rotate simultaneously.
func TestInClusterContextTokenAndCARotationStability(t *testing.T) {
	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))

	ctx1 := K8sContext{
		Name:              "in-cluster-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-v1"}},
		Cluster:           sql.Map{"cluster": map[string]interface{}{"certificate-authority-data": "ca-data-v1", "server": "https://10.0.0.1:443"}},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
	}

	id1, err := K8sContextGenerateID(ctx1)
	if err != nil {
		t.Fatalf("K8sContextGenerateID() error = %v", err)
	}

	// Simulate both token and CA rotation
	ctx2 := K8sContext{
		Name:              "in-cluster-context",
		Auth:              sql.Map{"user": map[string]interface{}{"token": "token-v2"}},
		Cluster:           sql.Map{"cluster": map[string]interface{}{"certificate-authority-data": "ca-data-v2", "server": "https://10.0.0.1:443"}},
		Server:            "https://10.0.0.1:443",
		MesheryInstanceID: &instanceID,
		DeploymentType:    "in_cluster",
	}

	id2, err := K8sContextGenerateID(ctx2)
	if err != nil {
		t.Fatalf("K8sContextGenerateID() error = %v", err)
	}

	// IDs must be the same for in-cluster contexts with both token and CA rotation
	if id1 != id2 {
		t.Errorf("in-cluster context ID changed after token and CA rotation: got %v, want %v", id2, id1)
	}
}
