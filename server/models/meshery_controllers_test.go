package models

import (
	"context"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	controllersconfig "github.com/meshery/schemas/models/v1alpha1/controllers_config"
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

func TestGetControllerHandlersForEachContextReturnsCopy(t *testing.T) {
	mch := &MesheryControllersHelper{
		ctxControllerHandlers: map[MesheryController]controllers.IMesheryController{
			MesheryOperator: nil,
		},
	}

	got := mch.GetControllerHandlersForEachContext()
	if len(got) != 1 {
		t.Fatalf("expected 1 handler in the returned map, got %d", len(got))
	}

	// The returned map must be a copy: callers range it outside the lock, so
	// mutating it must not leak into the helper's internal state.
	got[MesheryBroker] = nil
	delete(got, MesheryOperator)

	if _, ok := mch.ctxControllerHandlers[MesheryOperator]; !ok {
		t.Fatal("deleting from the returned map removed an entry from the helper's internal map")
	}
	if _, ok := mch.ctxControllerHandlers[MesheryBroker]; ok {
		t.Fatal("adding to the returned map added an entry to the helper's internal map")
	}

	// A helper with no handlers attached returns nil, not a shared empty map.
	if (&MesheryControllersHelper{}).GetControllerHandlersForEachContext() != nil {
		t.Fatal("expected nil for a helper with no controller handlers")
	}
}

// TestMesheryControllersHelperConcurrentStateAccess drives the controller-state
// writers (the connect / reconcile / config-apply goroutines) against the
// readers (the controllers-status SSE stream and the status REST handlers) the
// way they run in production. It guards the data race fixed in #20807 and is
// meaningful under `go test -race`.
func TestMesheryControllersHelperConcurrentStateAccess(t *testing.T) {
	log, _ := logger.New("test", logger.Options{})

	// Seed a MeshSync data handler so RemoveMeshSyncDataHandler actually
	// exercises the clear-the-pointer-under-lock path. broker is nil and no
	// listeners are registered, so Stop() closes stopCh once and returns
	// immediately — safe to drive from the concurrent phase.
	seededHandler := &MeshsyncDataHandler{
		stopCh:     make(chan struct{}),
		stopOnce:   &sync.Once{},
		listenerWg: &sync.WaitGroup{},
	}

	mch := &MesheryControllersHelper{
		log:                    log,
		ctxControllerHandlers:  map[MesheryController]controllers.IMesheryController{MesheryBroker: nil},
		ctxOperatorStatus:      controllers.Unknown,
		ctxMeshsyncDataHandler: seededHandler,
		meshsyncDeploymentMode: connections.MeshsyncDeploymentModeOperator,
	}

	// Seed the tracker single-threaded, before any goroutine starts, so that
	// during the concurrent phase IsUndeployed only reads its own map — the
	// tracker's accessors are unsynchronized and out of scope for this test.
	ot := NewOperatorTracker(false)
	ot.Undeployed(mch.contextID, true)

	k8scontext := K8sContext{ID: "ctx", Name: "test"}

	ops := []func(){
		// writers
		func() { mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator) },
		func() { mch.SetControllersConfig(&controllersconfig.MesheryControllersConfig{}) },
		func() { mch.RemoveCtxControllerHandler(context.Background(), "ctx") },
		func() { mch.UpdateOperatorsStatusMap(ot) },
		func() {
			mch.AddMeshsyncDataHandlers(context.Background(), k8scontext, uuid.Nil, uuid.Nil, nil)
		},
		func() { mch.RemoveMeshSyncDataHandler(context.Background(), "ctx") },
		// readers
		func() { _ = mch.GetControllerHandlersForEachContext() },
		func() { _ = mch.GetOperatorsStatusMap() },
		func() { _ = mch.GetMeshsyncDeploymentMode() },
		func() { _ = mch.GetMeshSyncDataHandlersForEachContext() },
	}

	const iterations = 1000
	var wg sync.WaitGroup
	for _, op := range ops {
		op := op
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				op()
			}
		}()
	}
	wg.Wait()
}

// fakeController is a controllers.IMesheryController whose methods mutate an
// internal counter instead of reaching a cluster, so a data race on a shared
// handler surfaces under `go test -race`. It mirrors how the real meshkit
// handlers cache their status in place.
type fakeController struct {
	calls int
}

func (f *fakeController) GetName() string {
	f.calls++
	return "fake"
}

func (f *fakeController) GetStatus() controllers.MesheryControllerStatus {
	f.calls++
	return controllers.Deployed
}

func (f *fakeController) Deploy(force bool) error {
	f.calls++
	return nil
}

func (f *fakeController) Undeploy() error {
	f.calls++
	return nil
}

func (f *fakeController) GetPublicEndpoint() (string, error) {
	f.calls++
	return "", nil
}

func (f *fakeController) GetVersion() (string, error) {
	f.calls++
	return "v1", nil
}

func (f *fakeController) GetEndpointForPort(string) (string, error) {
	f.calls++
	return "", nil
}

// TestControllerHandlerCallsSerialized drives the shared meshkit controller
// handler concurrently from both paths that use it in production — the status
// handlers (through the wrappers GetControllerHandlersForEachContext hands out)
// and the FSM reconcile (UpdateOperatorsStatusMap) — against a handler that
// mutates internal state on every call. It guards the data race raised in review
// on #20887 (a copied map still leaks the underlying mutable handlers) and is
// meaningful under `go test -race`.
func TestControllerHandlerCallsSerialized(t *testing.T) {
	mch := &MesheryControllersHelper{
		ctxControllerHandlers: map[MesheryController]controllers.IMesheryController{
			MesheryOperator: &fakeController{},
		},
		meshsyncDeploymentMode: connections.MeshsyncDeploymentModeOperator,
	}

	// Seed the tracker single-threaded so IsUndeployed only reads during the
	// concurrent phase, and return false so UpdateOperatorsStatusMap actually
	// invokes the operator handler's GetStatus — the call being serialized.
	ot := NewOperatorTracker(false)
	ot.Undeployed(mch.contextID, false)

	ops := []func(){
		// reader: the status handlers only ever touch the handlers via these wrappers
		func() {
			for _, h := range mch.GetControllerHandlersForEachContext() {
				if h == nil {
					continue
				}
				_ = h.GetStatus()
				_, _ = h.GetVersion()
			}
		},
		// writer: the FSM reconcile invokes the same handler's GetStatus
		func() { mch.UpdateOperatorsStatusMap(ot) },
	}

	const (
		iterations  = 1000
		concurrency = 4
	)
	var wg sync.WaitGroup
	for _, op := range ops {
		op := op
		for g := 0; g < concurrency; g++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < iterations; i++ {
					op()
				}
			}()
		}
	}
	wg.Wait()
}

// TestOperatorFSMHandlesNilHandler ensures the operator FSM methods don't panic
// when ctxControllerHandlers holds a present-but-nil MesheryOperator entry — the
// map type allows it and the handler getter guards against it, so the FSM must
// gate on the handler value, not merely the map key's presence.
func TestOperatorFSMHandlesNilHandler(t *testing.T) {
	log, _ := logger.New("test", logger.Options{})
	mch := &MesheryControllersHelper{
		log:                    log,
		ctxControllerHandlers:  map[MesheryController]controllers.IMesheryController{MesheryOperator: nil},
		ctxOperatorStatus:      controllers.NotDeployed,
		meshsyncDeploymentMode: connections.MeshsyncDeploymentModeOperator,
	}

	ot := NewOperatorTracker(false)
	ot.Undeployed(mch.contextID, false)

	// A present-but-nil handler must be skipped, not dereferenced.
	mch.UpdateOperatorsStatusMap(ot)
	mch.DeployUndeployedOperators(ot)
	mch.UndeployDeployedOperators(ot)
}
