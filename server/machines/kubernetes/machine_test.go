package kubernetes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	mkkube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/meshery/schemas/models/core"
	k8sclient "k8s.io/client-go/kubernetes"
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

	sysID := core.Uuid(uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000"))

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

func TestAssignInitialCtx_NilContextDoesNotPanic(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build test logger: %v", err)
	}

	machinectx := &MachineCtx{
		K8sContext: models.K8sContext{
			Name:         "unreachable-test-context",
			Server:       "https://127.0.0.1:1",
			ConnectionID: uuid.Must(uuid.NewV4()).String(),
		},
	}

	// Calling AssignInitialCtx with context.Background() where UserCtxKey and SystemIDKey are absent.
	// Previously this caused a nil-pointer dereference panic immediately when extracting user.ID or *sysID.
	result, _, err := AssignInitialCtx(context.Background(), machinectx, log)
	if err == nil {
		t.Fatal("expected AssignInitialCtx to return an error for an unreachable K8s context, got nil")
	}
	if result != nil {
		t.Fatalf("expected nil machine context on AssignClientSetToContext error, got %#v", result)
	}
}

// TestDiscoverAction_ContextCancellationReleasesStateMachineLock asserts that canceling
// the context during an active in-flight discovery request aborts promptly and releases
// the state machine mutex lock so subsequent operations like ResetState succeed.
func TestDiscoverAction_ContextCancellationReleasesStateMachineLock(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build test logger: %v", err)
	}

	connectionUUID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate connection UUID: %v", err)
	}
	userUUID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user UUID: %v", err)
	}

	sm, err := New(connectionUUID.String(), userUUID, log)
	if err != nil {
		t.Fatalf("failed to initialize K8s state machine: %v", err)
	}

	requestReceived := make(chan struct{}, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case requestReceived <- struct{}{}:
		default:
		}
		<-r.Context().Done()
	}))
	defer srv.CloseClientConnections()
	defer srv.Close()

	clientset, err := k8sclient.NewForConfig(&rest.Config{Host: srv.URL, Timeout: 30 * time.Second})
	if err != nil {
		t.Fatalf("build clientset: %v", err)
	}

	sm.Context = &MachineCtx{
		K8sContext: models.K8sContext{
			Name:         "test-cluster",
			Server:       srv.URL,
			ConnectionID: connectionUUID.String(),
		},
		clientset: &mkkube.Client{KubeClient: clientset},
		log:       log,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sendDone := make(chan error, 1)
	start := time.Now()
	go func() {
		_, sendErr := sm.SendEvent(ctx, machines.Discovery, nil)
		sendDone <- sendErr
	}()

	select {
	case <-requestReceived:
		// Active in-flight request reached mock server; cancel context during request
		cancel()
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for request to reach mock server")
	}

	select {
	case sendErr := <-sendDone:
		if sendErr == nil {
			t.Fatal("expected SendEvent to return error on context cancellation, got nil")
		}
		if time.Since(start) > 2*time.Second {
			t.Fatalf("SendEvent took %v; expected prompt abort upon cancellation", time.Since(start))
		}
	case <-time.After(3 * time.Second):
		t.Fatal("SendEvent did not return within 3s after cancellation")
	}

	// Verify that the state machine lock was released by ensuring ResetState succeeds immediately
	lockAcquired := make(chan struct{})
	go func() {
		sm.ResetState()
		close(lockAcquired)
	}()

	select {
	case <-lockAcquired:
		// success: lock was released and ResetState succeeded
	case <-time.After(1 * time.Second):
		t.Fatal("state machine mutex lock was not released after context cancellation")
	}
}

