package handlers

// Integration coverage verifying that a Kubernetes cluster that rejects
// Meshery's credential is probed once, parked in DISCONNECTED, and never
// probed again on idle requests until an explicit user action occurs.
//
// The test drives the real middleware pair - KubernetesMiddleware (greedy
// kubeconfig import) followed by K8sFSMMiddleware (FSM re-drive) - against a
// fake API server, for both provider shapes:
//
//	local  : LoadAllK8sContext ignores withStatus, as DefaultLocalProvider does,
//	         so a stored context reaches K8sFSMMiddleware whatever its status.
//	remote : LoadAllK8sContext honours withStatus=connected, as RemoteProvider
//	         does, so a parked context makes it return nothing and the greedy
//	         kubeconfig import runs on every request.
//
// Both shapes must end up probing the cluster zero times while the user is
// idle, which is the invariant the issue is about.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// k8sProviderStub implements the handful of Provider methods the two middlewares
// and the Kubernetes state machine reach. honorStatus switches between the two
// real LoadAllK8sContext behaviours.
type k8sProviderStub struct {
	models.Provider

	mu          sync.RWMutex
	known       bool // a kubernetes context (and its connection) is persisted
	conn        connections.Connection
	server      string
	honorStatus bool

	db      *database.Handler
	saves   atomic.Int32
	updated chan struct{}
}

func (p *k8sProviderStub) k8sContext() *models.K8sContext {
	return &models.K8sContext{
		ID:           "test-ctx-id",
		Name:         "test-context",
		Server:       p.server,
		ConnectionID: p.conn.ID.String(),
		Cluster: map[string]interface{}{"name": "test-cluster", "cluster": map[string]interface{}{
			"server": p.server, "insecure-skip-tls-verify": true}},
		Auth: map[string]interface{}{"name": "test-user", "user": map[string]interface{}{"token": "test-token"}},
	}
}

func (p *k8sProviderStub) LoadAllK8sContext(_ string) ([]*models.K8sContext, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if !p.known || (p.honorStatus && p.conn.Status != connections.CONNECTED) {
		return []*models.K8sContext{}, nil
	}
	return []*models.K8sContext{p.k8sContext()}, nil
}

// GetK8sContexts with an empty withStatus lists every status, for both providers.
func (p *k8sProviderStub) GetK8sContexts(_, _, _, _, _ string, withStatus string, _ bool) ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	page := models.MesheryK8sContextPage{Contexts: []*models.K8sContext{}}
	if p.known && (withStatus == "" || connections.ConnectionStatus(withStatus) == p.conn.Status) {
		page.TotalCount = 1
		page.Contexts = append(page.Contexts, p.k8sContext())
	}
	return json.Marshal(page)
}

func (p *k8sProviderStub) GetConnectionByID(_ string, _ core.Uuid) (*connections.Connection, int, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if !p.known {
		return nil, http.StatusNotFound, fmt.Errorf("connection not found")
	}
	conn := p.conn
	return &conn, http.StatusOK, nil
}

func (p *k8sProviderStub) UpdateConnectionById(_ string, payload *connections.ConnectionPayload, _ string) (*connections.Connection, error) {
	p.mu.Lock()
	p.conn.Status = payload.Status
	conn := p.conn
	p.mu.Unlock()
	select {
	case p.updated <- struct{}{}:
	default:
	}
	return &conn, nil
}

// SaveK8sContext mirrors ConnectionPersister.SaveConnection: the first save
// creates a DISCOVERED connection, a later save for the same cluster returns the
// existing row with its current status.
func (p *k8sProviderStub) SaveK8sContext(_ string, _ models.K8sContext, _ map[string]any) (connections.Connection, error) {
	p.saves.Add(1)
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.known {
		p.known = true
		p.conn.Status = connections.DISCOVERED
	}
	return p.conn, nil
}

func (p *k8sProviderStub) PersistEvent(_ events.Event, _ string) error { return nil }
func (p *k8sProviderStub) GetGenericPersister() *database.Handler      { return p.db }

func (p *k8sProviderStub) status() connections.ConnectionStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.conn.Status
}

// connect is what the UI's explicit "connect" action ends up persisting.
func (p *k8sProviderStub) connect() {
	p.mu.Lock()
	p.conn.Status = connections.CONNECTED
	p.mu.Unlock()
}

// apiStats counts the reachability probe, GET /api/v1/namespaces/kube-system,
// and holds the status the probe is answered with so a test can change it.
type apiStats struct {
	probes atomic.Int32
	code   atomic.Int32
	hits   chan struct{}
}

func newFakeAPIServer(t *testing.T, code int) (*httptest.Server, *apiStats) {
	t.Helper()
	st := &apiStats{hits: make(chan struct{}, 64)}
	st.code.Store(int32(code))
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case st.hits <- struct{}{}:
		default:
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v1/namespaces/kube-system":
			// The reachability probe shared by kubeconfig discovery and the FSM.
			st.probes.Add(1)
			code := int(st.code.Load())
			reason := map[int]string{
				http.StatusUnauthorized:       "Unauthorized",
				http.StatusForbidden:          "Forbidden",
				http.StatusServiceUnavailable: "ServiceUnavailable",
			}[code]
			w.WriteHeader(code)
			_, _ = fmt.Fprintf(w, `{"kind":"Status","apiVersion":"v1","metadata":{},"status":"Failure","message":%q,"reason":%q,"code":%d}`, reason, reason, code)
		case "/version":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"major":"1","minor":"30","gitVersion":"v1.30.0"}`))
		default:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		}
	}))
	t.Cleanup(srv.Close)
	return srv, st
}

func writeKubeconfig(t *testing.T, dir, server string) {
	t.Helper()
	kubeconfig := fmt.Sprintf(`apiVersion: v1
kind: Config
current-context: test-context
clusters:
- name: test-cluster
  cluster:
    server: %s
    insecure-skip-tls-verify: true
contexts:
- name: test-context
  context:
    cluster: test-cluster
    user: test-user
users:
- name: test-user
  user:
    token: test-token
`, server)
	if err := os.WriteFile(filepath.Join(dir, "config"), []byte(kubeconfig), 0o600); err != nil {
		t.Fatal(err)
	}
}

func newK8sMiddlewareHandler(t *testing.T, kubeConfigFolder string) *Handler {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatal(err)
	}
	return &Handler{
		log: log,
		config: &models.HandlerConfig{
			KubeConfigFolder:  kubeConfigFolder,
			EventBroadcaster:  models.NewBroadcaster("test"),
			OperatorTracker:   models.NewOperatorTracker(false),
			K8scontextChannel: models.NewContextHelper(),
		},
		ConnectionToStateMachineInstanceTracker: &machines.ConnectionToStateMachineInstanceTracker{
			ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
		},
		MesheryCtrlsHelper: &models.MesheryControllersHelper{},
	}
}

// syncRequest runs one /api/system/sync-shaped request through both middlewares
// and returns the number of reachability probes it caused. Both middlewares
// dispatch FSM work in goroutines, so it waits for the request to go quiet
// rather than for a fixed duration.
func syncRequest(ctx context.Context, t *testing.T, h *Handler, p *k8sProviderStub, user *models.User, st *apiStats) int32 {
	t.Helper()
	before := st.probes.Load()
	drain(st, p)

	reqCtx, err := KubernetesMiddleware(ctx, h, p, user, []string{"all"})
	if err != nil {
		t.Fatalf("KubernetesMiddleware: %v", err)
	}
	K8sFSMMiddleware(reqCtx, h, p, user)

	const quiet = 400 * time.Millisecond
	deadline := time.After(10 * time.Second)
	for {
		select {
		case <-st.hits:
		case <-p.updated:
		case <-time.After(quiet):
			return st.probes.Load() - before
		case <-deadline:
			t.Log("request did not go quiet within 10s")
			return st.probes.Load() - before
		}
	}
}

func drain(st *apiStats, p *k8sProviderStub) {
	for {
		select {
		case <-st.hits:
		case <-p.updated:
		default:
			return
		}
	}
}

// TestKubernetesMiddleware_DoesNotReprobeParkedClusters verifies that parked clusters
// are not re-probed on idle requests. The 503 rows are the control: a transient failure must stay
// retryable, so only a rejected credential may park a connection.
func TestKubernetesMiddleware_DoesNotReprobeParkedClusters(t *testing.T) {
	for _, tc := range []struct {
		name string
		// honorStatus picks the provider shape; code is what the API server
		// answers to the reachability probe.
		honorStatus bool
		code        int
		// persisted starts with a CONNECTED connection already stored. 401/403
		// start from nothing, which is the issue's shape: the cluster is only
		// known through the kubeconfig on disk.
		persisted bool
		// wantStatus is the status the connection must end in after the first
		// request, and wantRetry whether an idle request may probe again.
		wantStatus connections.ConnectionStatus
		wantRetry  bool
	}{
		{name: "local/401", code: http.StatusUnauthorized, wantStatus: connections.DISCONNECTED},
		{name: "local/403", code: http.StatusForbidden, wantStatus: connections.DISCONNECTED},
		{name: "local/503", code: http.StatusServiceUnavailable, persisted: true, wantStatus: connections.NOTFOUND, wantRetry: true},
		{name: "remote/401", honorStatus: true, code: http.StatusUnauthorized, wantStatus: connections.DISCONNECTED},
		{name: "remote/403", honorStatus: true, code: http.StatusForbidden, wantStatus: connections.DISCONNECTED},
		{name: "remote/503", honorStatus: true, code: http.StatusServiceUnavailable, persisted: true, wantStatus: connections.NOTFOUND, wantRetry: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv, st := newFakeAPIServer(t, tc.code)
			db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
			if err != nil {
				t.Fatal(err)
			}
			dir := t.TempDir()
			writeKubeconfig(t, dir, srv.URL)

			systemID := core.Uuid(uuid.Must(uuid.NewV4()))
			// INSTANCE_ID is package-global; DiscoverK8SContextFromKubeConfig reads it.
			previousInstanceID := viper.Get("INSTANCE_ID")
			viper.Set("INSTANCE_ID", &systemID)
			t.Cleanup(func() { viper.Set("INSTANCE_ID", previousInstanceID) })

			provider := &k8sProviderStub{
				known:       tc.persisted,
				conn:        connections.Connection{ID: core.Uuid(uuid.Must(uuid.NewV4())), Kind: "kubernetes"},
				server:      srv.URL,
				honorStatus: tc.honorStatus,
				db:          &database.Handler{DB: db},
				updated:     make(chan struct{}, 8),
			}
			if tc.persisted {
				provider.conn.Status = connections.CONNECTED
			}

			h := newK8sMiddlewareHandler(t, dir)
			user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}
			ctx := context.WithValue(context.Background(), models.TokenCtxKey, "test-token")
			ctx = context.WithValue(ctx, models.UserCtxKey, user)
			ctx = context.WithValue(ctx, models.ProviderCtxKey, models.Provider(provider))
			ctx = context.WithValue(ctx, models.SystemIDKey, &systemID)

			// Request 1: the cluster must actually be probed and the failure classified.
			if probes := syncRequest(ctx, t, h, provider, user, st); probes == 0 {
				t.Fatal("first request did not probe the cluster, so the failure path was never exercised")
			}
			if got := provider.status(); got != tc.wantStatus {
				t.Fatalf("after a %d from the API server: status = %q, want %q", tc.code, got, tc.wantStatus)
			}

			// Requests 2 and 3: the user has done nothing.
			savesBefore := provider.saves.Load()
			idle := syncRequest(ctx, t, h, provider, user, st) + syncRequest(ctx, t, h, provider, user, st)
			if tc.wantRetry {
				if idle == 0 {
					t.Error("a transient failure must stay retryable, but idle requests never probed again")
				}
				return
			}
			if idle != 0 {
				t.Errorf("idle requests probed the cluster %d times; want 0 probes for parked connection", idle)
			}
			if reimports := provider.saves.Load() - savesBefore; reimports != 0 {
				t.Errorf("idle requests re-imported the context %d times", reimports)
			}

			// An explicit reconnect must probe again: the park is not permanent, it
			// only withholds automatic rediscovery. The cluster answers the fresh
			// probe with a transient failure, which keeps this to one DisconnectAction
			// for the whole subtest - driving a second disconnect would trigger
			// unsynchronized concurrent access inside MesheryControllersHelper.
			st.code.Store(http.StatusServiceUnavailable)
			provider.connect()
			if probes := syncRequest(ctx, t, h, provider, user, st); probes != 1 {
				t.Errorf("an explicit reconnect must probe exactly once, probed %d times", probes)
			}
			if got := provider.status(); got != connections.NOTFOUND {
				t.Errorf("after reconnecting to a cluster that is now merely unreachable: status = %q, want %q", got, connections.NOTFOUND)
			}
		})
	}
}
