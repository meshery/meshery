package handlers

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

func TestActionMutexBypassDuringModeChange(t *testing.T) {
	connID := uuid.Must(uuid.NewV4())
	log := newTestLogger(t)

	provider := &lifecycleTestMockProvider{
		k8sContext: models.K8sContext{
			ID:           connID.String(),
			Name:         "test-k8s",
			ConnectionID: connID.String(),
		},
	}
	h, tracker := newLifecycleTestHandler(t, provider)

	// Create and register a kubernetes machine in tracker
	sm, err := kubernetes.New(connID.String(), core.Uuid(connID), log)
	if err != nil {
		t.Fatalf("failed to create machine: %v", err)
	}
	sm.Provider = provider
	sm.CurrentState = machines.CONNECTED

	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))})
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")

	machineCtx := &kubernetes.MachineCtx{
		K8sContext:         provider.k8sContext,
		ActionMutex:        &sync.Mutex{},
		MesheryCtrlsHelper: models.NewMesheryControllersHelper(log, controllers.OperatorDeploymentConfig{}, nil, nil, provider, nil),
		OperatorTracker:    models.NewOperatorTracker(false),
	}
	_, err = sm.Start(ctx, machineCtx, log, func(c context.Context, mc interface{}, l logger.Handler) (interface{}, *events.Event, error) {
		return mc, nil, nil
	})
	if err != nil {
		t.Fatalf("start machine: %v", err)
	}

	tracker.Add(core.Uuid(connID), sm)

	// 1. Lock ActionMutex manually to simulate an ongoing Connect side-effect
	machineCtx.ActionMutex.Lock()

	// 2. Start MeshSync mode reconciliation concurrently
	modeChangeFinished := make(chan struct{})
	go func() {
		// This should block on ActionMutex
		_ = h.reconcileMeshsyncDeploymentMode(ctx, core.Uuid(connID), connections.MeshsyncDeploymentModeEmbedded, nil, core.Uuid(uuid.Must(uuid.NewV4())), provider)
		close(modeChangeFinished)
	}()

	// 3. Ensure mode reconciliation blocks and does not race
	select {
	case <-modeChangeFinished:
		t.Fatalf("reconcileMeshsyncDeploymentMode completed while ActionMutex was held by another operation. It bypassed the mutex!")
	case <-time.After(100 * time.Millisecond):
		// This is the expected behavior, it should block
	}

	// 4. Release ActionMutex and allow mode change to complete
	machineCtx.ActionMutex.Unlock()

	select {
	case <-modeChangeFinished:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatalf("reconcileMeshsyncDeploymentMode failed to complete after ActionMutex was released")
	}
}

func TestMeshSyncModeChangeConcurrentLifecycleTransition(t *testing.T) {
	connID := uuid.Must(uuid.NewV4())
	log := newTestLogger(t)

	provider := &lifecycleTestMockProvider{
		k8sContext: models.K8sContext{
			ID:           connID.String(),
			Name:         "test-k8s",
			ConnectionID: connID.String(),
		},
	}
	h, tracker := newLifecycleTestHandler(t, provider)

	sm, err := kubernetes.New(connID.String(), core.Uuid(connID), log)
	if err != nil {
		t.Fatalf("failed to create machine: %v", err)
	}
	sm.Provider = provider
	sm.CurrentState = machines.CONNECTED

	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))})
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")

	machineCtx := &kubernetes.MachineCtx{
		K8sContext:         provider.k8sContext,
		ActionMutex:        &sync.Mutex{},
		MesheryCtrlsHelper: models.NewMesheryControllersHelper(log, controllers.OperatorDeploymentConfig{}, nil, nil, provider, nil),
		OperatorTracker:    models.NewOperatorTracker(false),
	}

	_, err = sm.Start(ctx, machineCtx, log, func(c context.Context, mc interface{}, l logger.Handler) (interface{}, *events.Event, error) {
		return mc, nil, nil
	})
	if err != nil {
		t.Fatalf("start machine: %v", err)
	}

	tracker.Add(core.Uuid(connID), sm)

	// We want to run SendEvent and reconcileMeshsyncDeploymentMode concurrently
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// Reconcile mode
		_ = h.reconcileMeshsyncDeploymentMode(ctx, core.Uuid(connID), connections.MeshsyncDeploymentModeEmbedded, nil, core.Uuid(uuid.Must(uuid.NewV4())), provider)
	}()

	go func() {
		defer wg.Done()
		// SendEvent that modifies CurrentState
		_, _ = sm.SendEvent(ctx, machines.Disconnect, nil)
	}()

	wg.Wait()
}
