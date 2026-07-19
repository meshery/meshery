package models

import (
	"context"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
)

func TestControllerEventActedUponPrefersConnectionID(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())
	connectionID := uuid.Must(uuid.NewV4())

	actedUpon := controllerEventActedUpon(userID, map[string]any{
		"connectionID": connectionID.String(),
	})

	if actedUpon != connectionID {
		t.Fatalf("expected actedUpon to use connectionID, got %s", actedUpon)
	}
}

func TestControllerEventActedUponFallsBackToUserID(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())

	actedUpon := controllerEventActedUpon(userID, map[string]any{
		"connectionID": "not-a-uuid",
	})

	if actedUpon != userID {
		t.Fatalf("expected actedUpon to fall back to userID, got %s", actedUpon)
	}
}

func TestControllerEventActedUponReturnsNilWithoutValidIDs(t *testing.T) {
	actedUpon := controllerEventActedUpon(uuid.Nil, map[string]any{
		"connectionID": "",
	})

	if actedUpon != uuid.Nil {
		t.Fatalf("expected actedUpon to be nil UUID, got %s", actedUpon)
	}
}

func TestShouldPersistControllerEvent(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())
	resourceID := uuid.Must(uuid.NewV4())

	if !shouldPersistControllerEvent(userID, resourceID) {
		t.Fatal("expected controller event to be persisted when user and resource IDs are valid")
	}

	if shouldPersistControllerEvent(uuid.Nil, resourceID) {
		t.Fatal("expected controller event persistence to be skipped when user ID is nil")
	}

	if shouldPersistControllerEvent(userID, uuid.Nil) {
		t.Fatal("expected controller event persistence to be skipped when actedUpon is nil")
	}
}

func TestAddCtxControllerHandlersReturnsEarlyOnInvalidConfig(t *testing.T) {
	// Create an empty K8sContext which will result in an invalid kubeconfig
	// and eventually a failure in mesherykube.New()
	ctx := K8sContext{
		ID:   "test-context",
		Name: "test-cluster",
	}

	log, _ := logger.New("test", logger.Options{})

	mch := NewMesheryControllersHelper(
		log,
		controllers.OperatorDeploymentConfig{},
		nil,
		nil,
		nil,
		nil,
	)
	
	// A successful execution of this function without panicking indicates the fix is working.
	// We wrap in a defer-recover just to explicitly fail the test if a panic occurs.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("AddCtxControllerHandlers panicked, indicating the nil pointer bug is still present: %v", r)
		}
	}()

	mch.AddCtxControllerHandlers(ctx)

	// Additionally, verify that ctxControllerHandlers is still empty/nil because it should have returned early
	if len(mch.ctxControllerHandlers) != 0 {
		t.Fatalf("expected ctxControllerHandlers to be empty, got %d", len(mch.ctxControllerHandlers))
	}
}

// newMeshsyncLifecycleTestHelper builds a helper with just the dependencies the
// MeshSync data-handler lifecycle methods touch (a logger); the broker, DB,
// provider and event dependencies are not reached by the paths under test.
func newMeshsyncLifecycleTestHelper(t *testing.T) *MesheryControllersHelper {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build test logger: %v", err)
	}
	return NewMesheryControllersHelper(log, controllers.OperatorDeploymentConfig{}, nil, nil, nil, nil)
}

// TestRemoveMeshSyncDataHandlerIsIdempotent verifies the data-handler teardown
// clears the handler and is safe to call again (a second connect/disconnect cycle
// or a fan-out config-apply can call it after the handler is already gone).
func TestRemoveMeshSyncDataHandlerIsIdempotent(t *testing.T) {
	mch := newMeshsyncLifecycleTestHelper(t)

	// Attach a handler the way AddMeshsyncDataHandlers does (under msMu). A
	// zero-value handler is safe to Stop: a nil broker/StopFunc are skipped and
	// listenerWg.Wait returns immediately.
	mch.msMu.Lock()
	mch.ctxMeshsyncDataHandler = &MeshsyncDataHandler{}
	mch.msMu.Unlock()

	mch.RemoveMeshSyncDataHandler(context.Background(), "ctx-id")
	if mch.GetMeshSyncDataHandlersForEachContext() != nil {
		t.Fatal("expected the MeshSync data handler to be nil after RemoveMeshSyncDataHandler")
	}

	// Repeated teardown must be a no-op, not a panic.
	mch.RemoveMeshSyncDataHandler(context.Background(), "ctx-id")
	if mch.GetMeshSyncDataHandlersForEachContext() != nil {
		t.Fatal("expected the MeshSync data handler to remain nil after a repeated teardown")
	}
}

// TestMeshsyncDataHandlerLifecycleConcurrent exercises the MeshSync data-handler
// lifecycle from many goroutines at once — the connect action, the deployment-mode
// reconcile handler, and the controllers-config apply fan-out all mutate the same
// helper concurrently in production. Run with -race: before msMu guarded
// ctxMeshsyncDataHandler / brokerPortForward, the assign in AddMeshsyncDataHandlers
// raced the clear in RemoveMeshSyncDataHandler and the diagnostics accessors, and
// two callers could each attach a handler (duplicate broker subscriptions + a
// leaked, orphaned handler). It must now be race-free and settle to a single owner.
func TestMeshsyncDataHandlerLifecycleConcurrent(t *testing.T) {
	mch := newMeshsyncLifecycleTestHelper(t)

	const (
		writers    = 4
		readers    = 4
		iterations = 200
	)
	ctx := context.Background()

	var wg sync.WaitGroup

	// Writers attach a handler (mirroring the field assign in AddMeshsyncDataHandlers,
	// under msMu) and tear it down, over and over.
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				mch.msMu.Lock()
				mch.ctxMeshsyncDataHandler = &MeshsyncDataHandler{}
				mch.msMu.Unlock()
				mch.RemoveMeshSyncDataHandler(ctx, "ctx-id")
			}
		}()
	}

	// Readers hammer the accessors the status / diagnostics HTTP handlers call.
	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				_ = mch.GetMeshSyncDataHandlersForEachContext()
				_ = mch.GetBrokerPortForwardAddr()
			}
		}()
	}

	wg.Wait()

	// A final deterministic teardown must leave exactly zero live handlers — no
	// orphan survived the concurrent churn.
	mch.RemoveMeshSyncDataHandler(ctx, "ctx-id")
	if mch.GetMeshSyncDataHandlersForEachContext() != nil {
		t.Fatal("expected a single owner: no live MeshSync data handler should remain after teardown")
	}
}
