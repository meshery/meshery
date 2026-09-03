package kubernetes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/models/events"
	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/meshery/schemas/models/core"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	clientgo "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// TestAssignInitialCtx_AttachesLoggerBeforeClientSetAssignment guards against
// a nil-pointer panic on login when a persisted K8s context can't be reached:
// GenerateKubeHandler errors → GenerateK8sClientSet hits its log.Warn path →
// interface-method-on-nil panic.
//
// The bug was an ordering mistake in AssignInitialCtx: machinectx.log was
// assigned AFTER AssignClientSetToContext, which threaded the still-nil
// log through GenerateClientSetAction → GenerateK8sClientSet. Any persisted
// context whose API server wasn't reachable (common: stale contexts from a
// remote provider pointing at clusters this host can't route to) produced
// the panic on every /api request that went through K8sFSMMiddleware.
//
// This test exercises AssignInitialCtx with a K8sContext whose API server
// is unreachable, which forces the error path previously responsible for
// the nil-deref panic. The assertions:
//   - no panic on the error path
//   - AssignInitialCtx surfaces the underlying AssignClientSetToContext error
//     (so the caller can handle it) and does not return a populated context
//   - machinectx.log is the exact logger we passed (proving the attach
//     happened before any action could consume it)
func TestAssignInitialCtx_AttachesLoggerBeforeClientSetAssignment(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build test logger: %v", err)
	}

	// Fail fast on UUID generation so the event builder always sees a valid
	// user ID and the test setup is deterministic; silently leaving user.ID
	// unset would change the code path we're exercising.
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user UUID: %v", err)
	}
	user := &models.User{ID: core.Uuid(userID)}

	sysID := core.Uuid(uuid.Nil)

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	// ProviderCtxKey: a typed-nil is fine — AssignControllerHandlers is only
	// reached after AssignClientSetToContext, and that's the point we want
	// to defend. If AssignClientSetToContext returns an error we never reach
	// controller setup, which matches the production scenario.
	var provider models.Provider
	ctx = context.WithValue(ctx, models.ProviderCtxKey, provider)

	machinectx := &MachineCtx{
		K8sContext: models.K8sContext{
			// Deliberately empty: any unreachable/invalid kubeconfig is fine.
			// The point is to force GenerateKubeHandler to fail so the
			// previously panicking log.Warn path runs.
			Name:         "unreachable-test-context",
			Server:       "https://127.0.0.1:1", // RFC-reserved, refused instantly
			ConnectionID: uuid.Must(uuid.NewV4()).String(),
		},
		// clientset left nil to force AssignClientSetToContext to attempt
		// GenerateClientSetAction (the panicking path).
	}

	result, _, err := AssignInitialCtx(ctx, machinectx, log)

	// AssignClientSetToContext must fail for an unreachable/invalid context —
	// that's the exact production regression we're guarding. If this ever
	// returns nil here, either the test lost its repro or GenerateK8sClientSet
	// started tolerating unreachable servers, both of which invalidate this
	// guard.
	if err == nil {
		t.Fatal("expected AssignInitialCtx to return an error for an unreachable K8s context, got nil — the regression guard is no longer exercising the panicking path")
	}
	if result != nil {
		t.Fatalf("expected nil machine context on AssignClientSetToContext error, got %#v", result)
	}

	// Logger must be the exact instance we passed in: equality (not just
	// non-nil) proves the attach happened before AssignClientSetToContext
	// ran, which is the invariant the ordering fix establishes.
	if machinectx.log != log {
		t.Fatal("expected machinectx.log to be the logger passed into AssignInitialCtx and assigned before AssignClientSetToContext; a different or nil value reintroduces the login-panic ordering bug")
	}
}

func TestRegisteredState_HasNotFoundTransition(t *testing.T) {
	state := Registered()

	next, ok := state.Events[machines.NotFound]
	if !ok {
		t.Fatal("expected Registered state to define a NotFound transition")
	}

	if next != machines.NOTFOUND {
		t.Fatalf("expected NotFound transition to go to %q, got %q", machines.NOTFOUND, next)
	}
}

func TestDiscoveredState_HasTransitions(t *testing.T) {
	state := Discovered()

	expected := map[machines.EventType]machines.StateType{
		machines.Register:   machines.REGISTERED,
		machines.NotFound:   machines.NOTFOUND,
		machines.Delete:     machines.DELETED,
		machines.Disconnect: machines.DISCONNECTED,
	}

	for event, wantState := range expected {
		gotState, ok := state.Events[event]
		if !ok {
			t.Fatalf("expected Discovered state to define transition for %q", event)
		}
		if gotState != wantState {
			t.Fatalf("expected event %q transition to be %q, got %q", event, wantState, gotState)
		}
	}
}

type recordingTestProvider struct {
	models.Provider
	mu       sync.Mutex
	conn     *connections.Connection
	lastConn *connections.ConnectionPayload
}

func (p *recordingTestProvider) GetConnectionByID(token string, connectionID core.Uuid) (*connections.Connection, int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.conn != nil {
		return p.conn, 200, nil
	}
	return &connections.Connection{
		ID:     connectionID,
		Kind:   "kubernetes",
		Status: connections.DISCOVERED,
	}, 200, nil
}

func (p *recordingTestProvider) UpdateConnectionById(token string, conn *connections.ConnectionPayload, connId string) (*connections.Connection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lastConn = conn
	if p.conn != nil {
		p.conn.Status = conn.Status
	}
	return &connections.Connection{
		ID:     conn.ID,
		Kind:   conn.Kind,
		Status: conn.Status,
	}, nil
}

func (p *recordingTestProvider) SaveK8sContext(token string, k8sCtx models.K8sContext, metadata map[string]any) (connections.Connection, error) {
	return connections.Connection{
		ID:     core.Uuid(uuid.FromStringOrNil(k8sCtx.ConnectionID)),
		Kind:   "kubernetes",
		Status: connections.DISCOVERED,
	}, nil
}

// TestDiscoverAction_ErrorClassification verifies that DiscoverAction distinguishes
// authorization failures (403/401) from transient failures (503/network) on AssignServerID.
func TestDiscoverAction_ErrorClassification(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build test logger: %v", err)
	}

	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	sysID := core.Uuid(uuid.Nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userID})
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	ctx = context.WithValue(ctx, models.ProviderCtxKey, models.Provider(&recordingTestProvider{}))
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")

	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		wantEventType machines.EventType
		wantForbidden bool
		wantUnauth    bool
	}{
		{
			name:          "403 Forbidden transitions to Disconnect",
			statusCode:    http.StatusForbidden,
			responseBody:  `{"kind":"Status","apiVersion":"v1","metadata":{},"status":"Failure","message":"namespaces \"kube-system\" is forbidden: User \"system:serviceaccount:meshery-extensions:default\" cannot get resource \"namespaces\" in API group \"\" in the namespace \"kube-system\"","reason":"Forbidden","code":403}`,
			wantEventType: machines.Disconnect,
			wantForbidden: true,
		},
		{
			name:          "401 Unauthorized transitions to Disconnect",
			statusCode:    http.StatusUnauthorized,
			responseBody:  `{"kind":"Status","apiVersion":"v1","metadata":{},"status":"Failure","message":"Unauthorized","reason":"Unauthorized","code":401}`,
			wantEventType: machines.Disconnect,
			wantUnauth:    true,
		},
		{
			name:          "503 Service Unavailable transitions to NotFound (transient retry)",
			statusCode:    http.StatusServiceUnavailable,
			responseBody:  `{"kind":"Status","apiVersion":"v1","metadata":{},"status":"Failure","message":"service unavailable","reason":"ServiceUnavailable","code":503}`,
			wantEventType: machines.NotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/api/v1/namespaces/kube-system" {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tc.statusCode)
					_, _ = w.Write([]byte(tc.responseBody))
					return
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			restConfig := &rest.Config{
				Host:            server.URL,
				TLSClientConfig: rest.TLSClientConfig{Insecure: true},
			}
			clientset, err := clientgo.NewForConfig(restConfig)
			if err != nil {
				t.Fatalf("failed to build clientset: %v", err)
			}

			machinectx := &MachineCtx{
				K8sContext: models.K8sContext{
					Name:         "test-context",
					Server:       server.URL,
					ConnectionID: uuid.Must(uuid.NewV4()).String(),
				},
				clientset: &meshkitkube.Client{KubeClient: clientset},
				log:       log,
			}

			action := &DiscoverAction{}
			eventType, event, err := action.Execute(ctx, machinectx, nil)
			if err == nil {
				t.Fatal("expected DiscoverAction.Execute to fail, got nil error")
			}
			if eventType != tc.wantEventType {
				t.Fatalf("expected event type %q, got %q", tc.wantEventType, eventType)
			}
			if event == nil {
				t.Fatal("expected DiscoverAction to emit an event on AssignServerID failure")
			}
			if tc.wantForbidden && !k8serrors.IsForbidden(err) {
				t.Fatalf("expected error to satisfy k8serrors.IsForbidden, got: %v", err)
			}
			if tc.wantUnauth && !k8serrors.IsUnauthorized(err) {
				t.Fatalf("expected error to satisfy k8serrors.IsUnauthorized, got: %v", err)
			}
		})
	}
}

// TestStateMachine_DiscoveryAuthorizationFailure_TransitionsToDisconnected exercises the
// entire FSM lifecycle path for an authorization failure:
//
// When AssignServerID returns 403 Forbidden:
//  1. DiscoverAction returns machines.Disconnect
//  2. StateMachine transitions DISCOVERED → DISCONNECTED
//  3. Connection status in provider is persisted as DISCONNECTED
//  4. connections.ShouldConnectionBeManaged returns false (halting infinite discovery loops)
//  5. Error event is returned for the Notification Center
func TestStateMachine_DiscoveryAuthorizationFailure_TransitionsToDisconnected(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build test logger: %v", err)
	}
	connUUID := uuid.Must(uuid.NewV4())
	userUUID := core.Uuid(uuid.Must(uuid.NewV4()))

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/namespaces/kube-system" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"kind":"Status","apiVersion":"v1","metadata":{},"status":"Failure","message":"namespaces \"kube-system\" is forbidden: User \"system:serviceaccount:meshery-extensions:default\" cannot get resource \"namespaces\"","reason":"Forbidden","code":403}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	restConfig := &rest.Config{
		Host:            server.URL,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}
	clientset, err := clientgo.NewForConfig(restConfig)
	if err != nil {
		t.Fatalf("failed to build clientset: %v", err)
	}

	provider := &recordingTestProvider{
		conn: &connections.Connection{
			ID:     core.Uuid(connUUID),
			Kind:   "kubernetes",
			Status: connections.DISCOVERED,
		},
	}

	sm, err := New(connUUID.String(), userUUID, log)
	if err != nil {
		t.Fatalf("failed to create StateMachine: %v", err)
	}
	sm.AssignProvider(provider)

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userUUID})
	ctx = context.WithValue(ctx, models.SystemIDKey, &core.Uuid{})
	ctx = context.WithValue(ctx, models.ProviderCtxKey, models.Provider(provider))
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")

	machineCtx := &MachineCtx{
		MesheryCtrlsHelper: models.NewMesheryControllersHelper(log, controllers.OperatorDeploymentConfig{}, nil, nil, provider, nil),
		K8sContext: models.K8sContext{
			ID:           "test-ctx-id",
			Name:         "forbidden-cluster",
			Server:       server.URL,
			ConnectionID: connUUID.String(),
		},
		OperatorTracker: models.NewOperatorTracker(false),
		clientset:       &meshkitkube.Client{KubeClient: clientset},
		log:             log,
	}

	_, err = sm.Start(ctx, machineCtx, log, func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error) {
		return machineCtx, nil, nil
	})
	if err != nil {
		t.Fatalf("failed to start StateMachine: %v", err)
	}

	event, err := sm.SendEvent(ctx, machines.Discovery, nil)
	if err == nil {
		t.Fatal("expected SendEvent to return non-nil error on 403 Forbidden discovery failure")
	}

	// 1. FSM must settle in DISCONNECTED state
	if sm.CurrentState != machines.DISCONNECTED {
		t.Fatalf("expected FSM current state to be %q, got %q", machines.DISCONNECTED, sm.CurrentState)
	}

	// 2. Provider must have persisted the DISCONNECTED status
	if provider.lastConn == nil || provider.lastConn.Status != connections.DISCONNECTED {
		t.Fatalf("expected persisted connection status to be %q, got %#v", connections.DISCONNECTED, provider.lastConn)
	}

	// 3. Connection is no longer managed (halts infinite discovery loop)
	persistedConn := connections.Connection{Status: provider.lastConn.Status}
	if connections.ShouldConnectionBeManaged(persistedConn) {
		t.Fatalf("expected ShouldConnectionBeManaged to return false for %q connection, got true", connections.DISCONNECTED)
	}

	// 4. Notification event has Error severity
	if event == nil || event.Severity != events.Error {
		t.Fatalf("expected Error severity notification event, got %#v", event)
	}

	// 5. Server ID was not persisted in metadata on authorization failure
	if provider.lastConn != nil && provider.lastConn.MetaData != nil {
		if serverID, ok := provider.lastConn.MetaData["kubernetesServerId"]; ok && serverID != "" && serverID != nil {
			t.Fatalf("expected kubernetesServerId to not be persisted in metadata on authorization failure, got %v", serverID)
		}
	}
}

// TestStateMachine_DiscoveryTransientFailure_TransitionsToNotFound verifies that a temporary
// 503 / API server outage moves the connection to NOTFOUND (retaining retry eligibility)
// rather than prematurely disconnecting it.
func TestStateMachine_DiscoveryTransientFailure_TransitionsToNotFound(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build test logger: %v", err)
	}
	connUUID := uuid.Must(uuid.NewV4())
	userUUID := core.Uuid(uuid.Must(uuid.NewV4()))

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/namespaces/kube-system" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"kind":"Status","apiVersion":"v1","metadata":{},"status":"Failure","message":"service unavailable","reason":"ServiceUnavailable","code":503}`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	restConfig := &rest.Config{
		Host:            server.URL,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true},
	}
	clientset, err := clientgo.NewForConfig(restConfig)
	if err != nil {
		t.Fatalf("failed to build clientset: %v", err)
	}

	provider := &recordingTestProvider{
		conn: &connections.Connection{
			ID:     core.Uuid(connUUID),
			Kind:   "kubernetes",
			Status: connections.DISCOVERED,
		},
	}

	sm, err := New(connUUID.String(), userUUID, log)
	if err != nil {
		t.Fatalf("failed to create StateMachine: %v", err)
	}
	sm.AssignProvider(provider)

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userUUID})
	ctx = context.WithValue(ctx, models.SystemIDKey, &core.Uuid{})
	ctx = context.WithValue(ctx, models.ProviderCtxKey, models.Provider(provider))
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")

	machineCtx := &MachineCtx{
		MesheryCtrlsHelper: models.NewMesheryControllersHelper(log, controllers.OperatorDeploymentConfig{}, nil, nil, provider, nil),
		K8sContext: models.K8sContext{
			ID:           "test-ctx-id",
			Name:         "transient-cluster",
			Server:       server.URL,
			ConnectionID: connUUID.String(),
		},
		OperatorTracker: models.NewOperatorTracker(false),
		clientset:       &meshkitkube.Client{KubeClient: clientset},
		log:             log,
	}

	_, err = sm.Start(ctx, machineCtx, log, func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error) {
		return machineCtx, nil, nil
	})
	if err != nil {
		t.Fatalf("failed to start StateMachine: %v", err)
	}

	event, err := sm.SendEvent(ctx, machines.Discovery, nil)
	if err == nil {
		t.Fatal("expected SendEvent to return non-nil error on 503 Service Unavailable discovery failure")
	}

	// 1. FSM must settle in NOTFOUND state
	if sm.CurrentState != machines.NOTFOUND {
		t.Fatalf("expected FSM current state to be %q, got %q", machines.NOTFOUND, sm.CurrentState)
	}

	// 2. Provider must have persisted the NOTFOUND status
	if provider.lastConn == nil || provider.lastConn.Status != connections.NOTFOUND {
		t.Fatalf("expected persisted connection status to be %q, got %#v", connections.NOTFOUND, provider.lastConn)
	}

	// 3. Connection is STILL managed (eligible for retry on next discovery cycle)
	persistedConn := connections.Connection{Status: provider.lastConn.Status}
	if !connections.ShouldConnectionBeManaged(persistedConn) {
		t.Fatalf("expected ShouldConnectionBeManaged to return true for %q connection, got false", connections.NOTFOUND)
	}

	// 4. Notification event has Error severity
	if event == nil || event.Severity != events.Error {
		t.Fatalf("expected Error severity notification event, got %#v", event)
	}
}
