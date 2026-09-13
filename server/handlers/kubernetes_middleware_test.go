package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
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
)

// localProviderMock reproduces the one DefaultLocalProvider behaviour this test
// depends on: LoadAllK8sContext returns EVERY persisted Kubernetes context,
// whatever its connection status.
//
// DefaultLocalProvider.LoadAllK8sContext passes string(connections.CONNECTED)
// into DefaultLocalProvider.GetK8sContexts, but that method never reads its
// withStatus parameter - it forwards to
// MesheryK8sContextPersister.GetMesheryK8sContexts(search, order, page, pageSize),
// which has no status parameter and no join to the connections table. A mock
// that filtered DISCONNECTED contexts out would delete the production failure
// mode this test exists to catch, so it must not filter.
type localProviderMock struct {
	models.Provider

	mu     sync.RWMutex
	conn   connections.Connection
	server string
	db     *database.Handler

	updateCalls atomic.Int32

	// updated receives once per UpdateConnectionById call. StateMachine.SendEvent
	// always writes the connection status when it advances, so this is an exact,
	// non-timing signal that the state machine actually ran.
	updated chan struct{}
}

func (m *localProviderMock) LoadAllK8sContext(_ string) ([]*models.K8sContext, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Cluster and Auth are real persisted columns on k8s_contexts (sql.Map, no
	// gorm:"-"), so a context loaded from the local provider can always build a
	// working kube client. Reproduce that here or AssignInitialCtx fails and the
	// machine is skipped for the wrong reason.
	return []*models.K8sContext{
		{
			ID:           "test-ctx-id",
			Name:         "test-context",
			Server:       m.server,
			ConnectionID: m.conn.ID.String(),
			Cluster: map[string]interface{}{
				"name": "test-cluster",
				"cluster": map[string]interface{}{
					"server":                   m.server,
					"insecure-skip-tls-verify": true,
				},
			},
			Auth: map[string]interface{}{
				"name": "test-user",
				"user": map[string]interface{}{"token": "test-token"},
			},
		},
	}, nil
}

func (m *localProviderMock) GetConnectionByID(_ string, _ core.Uuid) (*connections.Connection, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Hand back a copy: StateMachine.SendEvent reads Status and Metadata off this
	// value outside our lock while UpdateConnectionById mutates the stored one.
	conn := m.conn
	return &conn, http.StatusOK, nil
}

func (m *localProviderMock) UpdateConnectionById(_ string, payload *connections.ConnectionPayload, _ string) (*connections.Connection, error) {
	m.mu.Lock()
	m.conn.Status = payload.Status
	conn := m.conn
	m.mu.Unlock()

	m.updateCalls.Add(1)
	select {
	case m.updated <- struct{}{}:
	default:
	}
	return &conn, nil
}

func (m *localProviderMock) SaveK8sContext(_ string, _ models.K8sContext, _ map[string]any) (connections.Connection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.conn, nil
}

func (m *localProviderMock) PersistEvent(_ events.Event, _ string) error { return nil }

func (m *localProviderMock) GetGenericPersister() *database.Handler { return m.db }

func (m *localProviderMock) status() connections.ConnectionStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.conn.Status
}

// TestKubernetesMiddleware_SecondInvocationAfter403 drives the middleware pair
// exactly as router/server.go wires /api/system/sync:
//
//	KubernetesMiddleware -> K8sFSMMiddleware
//
// The first request lets the real authorization-failure flow move the connection
// to DISCONNECTED. The second request must then produce no Kubernetes API
// traffic at all, because connections.ShouldConnectionBeManaged rejects a
// DISCONNECTED connection before K8sFSMMiddleware can call ResetState() and
// re-drive Discovery (issue #14083).
//
// The connection status is never set by hand: if the production transition
// breaks, phase 1 fails rather than silently handing phase 2 a fabricated
// precondition.
func TestKubernetesMiddleware_SecondInvocationAfter403(t *testing.T) {
	var kubeRequests atomic.Int32
	// requestPaths carries the path of every request the cluster receives, so
	// phase 2 can fail fast and name the request that should not have happened
	// instead of comparing a counter after a fixed sleep.
	requestPaths := make(chan string, 16)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		kubeRequests.Add(1)
		select {
		case requestPaths <- r.URL.Path:
		default:
		}

		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v1/namespaces/kube-system" {
			// Unconditional 403. This test provokes exactly one discovery, so the
			// mock never has to guess which request in a sequence it is answering.
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"kind":"Status","apiVersion":"v1","metadata":{},"status":"Failure","message":"namespaces \"kube-system\" is forbidden","reason":"Forbidden","code":403}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open in-memory database (this test needs CGO_ENABLED=1 and a C toolchain, like the other DB-backed tests in this package): %v", err)
	}

	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	provider := &localProviderMock{
		conn: connections.Connection{
			ID:     connID,
			Kind:   "kubernetes",
			Status: connections.DISCOVERED,
		},
		server:  server.URL,
		db:      &db,
		updated: make(chan struct{}, 4),
	}

	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("build test logger: %v", err)
	}

	h := &Handler{
		log: log,
		config: &models.HandlerConfig{
			KubeConfigFolder:  t.TempDir(),
			EventBroadcaster:  models.NewBroadcaster("test"),
			OperatorTracker:   models.NewOperatorTracker(false),
			K8scontextChannel: models.NewContextHelper(),
		},
		ConnectionToStateMachineInstanceTracker: &machines.ConnectionToStateMachineInstanceTracker{
			ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
		},
		MesheryCtrlsHelper: &models.MesheryControllersHelper{},
	}

	sysID := core.Uuid(uuid.Nil)
	user := &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))}
	baseCtx := context.WithValue(context.Background(), models.TokenCtxKey, "test-token")
	baseCtx = context.WithValue(baseCtx, models.UserCtxKey, user)
	baseCtx = context.WithValue(baseCtx, models.ProviderCtxKey, models.Provider(provider))
	baseCtx = context.WithValue(baseCtx, models.SystemIDKey, &sysID)

	// --- Phase 1: a DISCOVERED connection is discovered, gets 403, disconnects.
	ctx, err := KubernetesMiddleware(baseCtx, h, provider, user, []string{"all"})
	if err != nil {
		t.Fatalf("first KubernetesMiddleware invocation: %v", err)
	}
	K8sFSMMiddleware(ctx, h, provider, user)

	select {
	case <-provider.updated:
	case <-time.After(15 * time.Second):
		t.Fatal("timed out waiting for the state machine to persist a status change on the first invocation")
	}

	if got := provider.status(); got != connections.DISCONNECTED {
		t.Fatalf("after a 403 on discovery, want connection status %q, got %q", connections.DISCONNECTED, got)
	}
	requestsAfterFirst := kubeRequests.Load()
	if requestsAfterFirst == 0 {
		t.Fatal("the first invocation made no Kubernetes API request, so the 403 path was never exercised")
	}
	updatesAfterFirst := provider.updateCalls.Load()

	// --- Phase 2: the same connection, now DISCONNECTED, must be left alone.
	for drained := false; !drained; {
		select {
		case <-requestPaths:
		case <-provider.updated:
		default:
			drained = true
		}
	}

	ctx, err = KubernetesMiddleware(baseCtx, h, provider, user, []string{"all"})
	if err != nil {
		t.Fatalf("second KubernetesMiddleware invocation: %v", err)
	}
	K8sFSMMiddleware(ctx, h, provider, user)

	select {
	case path := <-requestPaths:
		t.Fatalf("the second invocation issued a Kubernetes API request to %q for a %q connection; "+
			"K8sFSMMiddleware did not gate rediscovery on ShouldConnectionBeManaged (issue #14083)",
			path, connections.DISCONNECTED)
	case <-provider.updated:
		t.Fatalf("the second invocation drove the state machine for a %q connection; "+
			"K8sFSMMiddleware did not gate rediscovery on ShouldConnectionBeManaged (issue #14083)",
			connections.DISCONNECTED)
	case <-time.After(2 * time.Second):
		// Nothing happened, which is the point.
	}

	if got := kubeRequests.Load(); got != requestsAfterFirst {
		t.Fatalf("want %d Kubernetes API requests after the second invocation, got %d", requestsAfterFirst, got)
	}
	if got := provider.updateCalls.Load(); got != updatesAfterFirst {
		t.Fatalf("want %d connection status writes after the second invocation, got %d", updatesAfterFirst, got)
	}
	if got := provider.status(); got != connections.DISCONNECTED {
		t.Fatalf("want the connection to remain %q, got %q", connections.DISCONNECTED, got)
	}

	inst, ok := h.ConnectionToStateMachineInstanceTracker.Get(connID)
	if !ok {
		t.Fatal("expected the Kubernetes state machine to still be tracked after the second invocation")
	}
	// Reading CurrentState without the machine's lock is safe here: control only
	// reaches this line when the select above timed out, which means no SendEvent
	// goroutine was ever started.
	if inst.CurrentState != machines.DISCONNECTED {
		t.Fatalf("want the machine to still be in %q, got %q - ResetState() ran, so discovery was re-driven",
			machines.DISCONNECTED, inst.CurrentState)
	}
}
