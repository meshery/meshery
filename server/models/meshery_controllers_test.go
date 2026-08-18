package models

import (
	"context"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/schemas/models/core"
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

// TestMesheryControllersHelperConcurrency verifies safe R/W to map under RWMutex.
// Added as part of #21265 lifecycle serialization.
func TestMesheryControllersHelperConcurrency(t *testing.T) {
	mch := &MesheryControllersHelper{
		ctxControllerHandlers: make(map[MesheryController]controllers.IMesheryController),
	}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			mch.GetControllerHandlersForEachContext()
			mch.GetMeshsyncDataHandler()
		}()
		go func() {
			defer wg.Done()
			mch.mu.Lock()
			mch.ctxControllerHandlers = make(map[MesheryController]controllers.IMesheryController)
			mch.mu.Unlock()
		}()
	}
	wg.Wait()
}

// TestMesheryControllersHelperOperatorStatusConcurrency verifies concurrent calls to
// UpdateOperatorsStatusMap, GetOperatorsStatusMap, and operatorStatusObserved
// are safe under the race detector.
func TestMesheryControllersHelperOperatorStatusConcurrency(t *testing.T) {
	stub := &stubController{status: controllers.Deployed}
	mch := &MesheryControllersHelper{
		ctxControllerHandlers: map[MesheryController]controllers.IMesheryController{
			MesheryOperator: stub,
		},
		meshsyncDeploymentMode: connections.MeshsyncDeploymentModeOperator,
	}
	ot := NewOperatorTracker(false)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(4)
		go func() {
			defer wg.Done()
			mch.UpdateOperatorsStatusMap(ot)
		}()
		go func() {
			defer wg.Done()
			_ = mch.GetOperatorsStatusMap()
		}()
		go func() {
			defer wg.Done()
			_ = mch.operatorStatusObserved()
		}()
		go func() {
			defer wg.Done()
			mch.mu.Lock()
			mch.ctxOperatorStatus = controllers.Undeployed
			mch.mu.Unlock()
		}()
	}
	wg.Wait()
}

// TestMesheryControllersHelper_ConcurrentMeshsyncInitialization verifies that
// concurrent calls to AddMeshsyncDataHandlers and RemoveMeshSyncDataHandler
// are properly serialized by meshsyncInitMu without race conditions or multiple handlers.
func TestMesheryControllersHelper_ConcurrentMeshsyncInitialization(t *testing.T) {
	log, _ := logger.New("test", logger.Options{})
	mch := &MesheryControllersHelper{
		ctxControllerHandlers:  make(map[MesheryController]controllers.IMesheryController),
		meshsyncDeploymentMode: connections.MeshsyncDeploymentMode("unknown-mode"),
		log:                    log,
	}

	ctx := context.Background()
	k8sCtx := K8sContext{ID: "test-ctx", ConnectionID: uuid.Must(uuid.NewV4()).String()}
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))

	startGate := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(4)
		go func() {
			defer wg.Done()
			<-startGate
			mch.AddMeshsyncDataHandlers(ctx, k8sCtx, userID, sysID, nil)
		}()
		go func() {
			defer wg.Done()
			<-startGate
			mch.RemoveMeshSyncDataHandler(ctx, k8sCtx.ID)
		}()
		go func() {
			defer wg.Done()
			<-startGate
			_ = mch.GetMeshsyncDataHandler()
		}()
		go func() {
			defer wg.Done()
			<-startGate
			_ = mch.ResyncMeshsync(ctx)
		}()
	}
	close(startGate)
	wg.Wait()
}

// TestMesheryControllersHelper_TeardownReleasesSharedLockBeforeStop verifies that
// RemoveMeshSyncDataHandler releases mch.mu before calling handler.Stop(), while
// still holding meshsyncInitMu to serialize against concurrent initialization.
func TestMesheryControllersHelper_TeardownReleasesSharedLockBeforeStop(t *testing.T) {
	log, _ := logger.New("test", logger.Options{})
	mch := &MesheryControllersHelper{
		ctxControllerHandlers:  make(map[MesheryController]controllers.IMesheryController),
		meshsyncDeploymentMode: connections.MeshsyncDeploymentMode("unknown-mode"),
		log:                    log,
	}

	stopCalled := make(chan struct{})
	stopUnblock := make(chan struct{})

	handler := &MeshsyncDataHandler{
		StopFunc: func() {
			close(stopCalled)
			<-stopUnblock
		},
	}
	mch.ctxMeshsyncDataHandler = handler

	ctx := context.Background()
	done := make(chan struct{})
	go func() {
		mch.RemoveMeshSyncDataHandler(ctx, "test-ctx")
		close(done)
	}()

	// Wait until handler.Stop() is executing
	<-stopCalled

	// 1. Verify mch.mu is NOT held: GetMeshsyncDataHandler can acquire RLock and returns nil immediately.
	got := mch.GetMeshsyncDataHandler()
	if got != nil {
		t.Fatalf("expected ctxMeshsyncDataHandler to be nil while Stop() runs, got: %v", got)
	}

	// 2. Verify meshsyncInitMu IS held during Stop(): TryLock must fail.
	if mch.meshsyncInitMu.TryLock() {
		mch.meshsyncInitMu.Unlock()
		t.Fatal("expected meshsyncInitMu to remain locked while Stop() is executing")
	}

	// 3. Verify initialization cannot begin while Stop() is running:
	// A concurrent AddMeshsyncDataHandlers caller will block on meshsyncInitMu.
	initStarted := make(chan struct{})
	initDone := make(chan struct{})
	k8sCtx := K8sContext{ID: "test-ctx", ConnectionID: uuid.Must(uuid.NewV4()).String()}
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	go func() {
		close(initStarted)
		mch.AddMeshsyncDataHandlers(ctx, k8sCtx, userID, sysID, nil)
		close(initDone)
	}()

	<-initStarted
	select {
	case <-initDone:
		t.Fatal("expected AddMeshsyncDataHandlers to block while Stop() is holding meshsyncInitMu")
	default:
	}

	// Unblock Stop() and wait for RemoveMeshSyncDataHandler to finish.
	close(stopUnblock)
	<-done
	<-initDone

	// Verify meshsyncInitMu is released after teardown and init complete.
	if !mch.meshsyncInitMu.TryLock() {
		t.Fatal("expected meshsyncInitMu to be unlocked after RemoveMeshSyncDataHandler returns")
	}
	mch.meshsyncInitMu.Unlock()
}
