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

// The state a MesheryControllersHelper keeps for its Kubernetes context is
// written from the connect goroutine that machines/kubernetes spawns
// (AddCtxControllerHandlers -> SetMeshsyncDeploymentMode -> SetControllersConfig
// -> UpdateOperatorsStatusMap -> AddMeshsyncDataHandlers) while the
// controller-status and MeshSync HTTP handlers read the same fields through the
// exported getters. The tests below reproduce that shape so `go test -race`
// fails if the guarding on those fields regresses.

func newTestControllersHelper() *MesheryControllersHelper {
	log, _ := logger.New("test", logger.Options{})
	return NewMesheryControllersHelper(log, controllers.OperatorDeploymentConfig{}, nil, nil, nil, nil)
}

// raceHammer runs a writer and a reader concurrently over the same helper.
func raceHammer(write, read func()) {
	const iterations = 200
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			write()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			read()
		}
	}()
	wg.Wait()
}

func TestMeshsyncDeploymentModeConcurrentAccess(t *testing.T) {
	mch := newTestControllersHelper()
	modes := []connections.MeshsyncDeploymentMode{
		connections.MeshsyncDeploymentModeOperator,
		connections.MeshsyncDeploymentModeEmbedded,
	}
	i := 0
	raceHammer(
		func() { mch.SetMeshsyncDeploymentMode(modes[i%len(modes)]); i++ },
		func() { _ = mch.GetMeshsyncDeploymentMode() },
	)
}

func TestControllersConfigConcurrentAccess(t *testing.T) {
	mch := newTestControllersHelper()
	cfg := &controllersconfig.MesheryControllersConfig{}
	raceHammer(
		func() { mch.SetControllersConfig(cfg) },
		func() { _ = mch.getControllersConfig() },
	)
}

func TestControllerHandlersConcurrentAccess(t *testing.T) {
	mch := newTestControllersHelper()
	raceHammer(
		func() {
			mch.setCtxControllerHandlers(map[MesheryController]controllers.IMesheryController{})
			mch.RemoveCtxControllerHandler(context.Background(), "ctx")
		},
		func() { _ = mch.GetControllerHandlersForEachContext() },
	)
}

func TestOperatorStatusConcurrentAccess(t *testing.T) {
	mch := newTestControllersHelper()
	statuses := []controllers.MesheryControllerStatus{
		controllers.Deployed,
		controllers.NotDeployed,
		controllers.Undeployed,
	}
	i := 0
	raceHammer(
		func() { mch.setCtxOperatorStatus(statuses[i%len(statuses)]); i++ },
		func() { _ = mch.GetOperatorsStatusMap() },
	)
}

func TestMeshsyncDataHandlerConcurrentAccess(t *testing.T) {
	mch := newTestControllersHelper()
	raceHammer(
		func() { mch.setCtxMeshsyncDataHandler(nil) },
		func() { _ = mch.GetMeshSyncDataHandlersForEachContext() },
	)
}
