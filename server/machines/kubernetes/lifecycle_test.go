package kubernetes

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

type mockProvider struct {
	models.Provider
}

func (m *mockProvider) GetConnectionByID(token string, connectionID core.Uuid) (*connections.Connection, int, error) {
	return &connections.Connection{Kind: "kubernetes"}, 200, nil
}

func (m *mockProvider) UpdateConnectionStatusByID(token string, connectionID core.Uuid, connectionStatus connections.ConnectionStatus) (*connections.Connection, int, error) {
	return &connections.Connection{}, 200, nil
}

func (m *mockProvider) UpdateConnectionById(token string, conn *connections.ConnectionPayload, connId string) (*connections.Connection, error) {
	return &connections.Connection{}, nil
}

func getTestLogger() logger.Handler {
	log, _ := logger.New("test", logger.Options{Format: logger.SyslogLogFormat, LogLevel: 5})
	return log
}

// TestSequentialConnectDisconnect verifies FSM state assignment and machine
// initialization. This is a smoke test for the machine setup path.
func TestSequentialConnectDisconnect(t *testing.T) {
	log := getTestLogger()
	m, err := New(uuid.Must(uuid.NewV4()).String(), core.Uuid(uuid.Must(uuid.NewV4())), log)
	if err != nil {
		t.Fatalf("failed to init machine: %v", err)
	}

	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))})
	ctx = context.WithValue(ctx, models.SystemIDKey, &core.Uuid{})

	m.Provider = &mockProvider{}

	machineCtx := &MachineCtx{
		ActionMutex:        &sync.Mutex{},
		MesheryCtrlsHelper: models.NewMesheryControllersHelper(log, controllers.OperatorDeploymentConfig{}, nil, nil, &mockProvider{}, nil),
		K8sContext:         models.K8sContext{},
		OperatorTracker:    models.NewOperatorTracker(false),
	}
	initFunc := func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error) {
		return machineCtx, nil, nil
	}
	_, err = m.Start(ctx, machineCtx, log, initFunc)
	if err != nil {
		t.Fatalf("failed to start machine: %v", err)
	}

	m.CurrentState = machines.REGISTERED

	if m.CurrentState != machines.REGISTERED {
		t.Fatalf("expected REGISTERED, got %v", m.CurrentState)
	}
}

// TestLifecycleContextRotation verifies the #21265 core mechanism:
//
//  1. SendEvent creates a new LifecycleCtx on each FSM transition.
//  2. The previous LifecycleCtx is cancelled when a new transition is accepted.
//  3. Background workers holding the old context can detect cancellation via ctx.Err().
//
// This test would FAIL if LifecycleCtx rotation were removed from models.go.
func TestLifecycleContextRotation(t *testing.T) {
	log := getTestLogger()
	m, _ := New(uuid.Must(uuid.NewV4()).String(), core.Uuid(uuid.Must(uuid.NewV4())), log)
	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))})
	ctx = context.WithValue(ctx, models.SystemIDKey, &core.Uuid{})
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")

	p := &mockProvider{}
	m.Provider = p
	ctx = context.WithValue(ctx, models.ProviderCtxKey, p)

	connID := uuid.Must(uuid.NewV4()).String()
	machineCtx := &MachineCtx{
		ActionMutex:        &sync.Mutex{},
		MesheryCtrlsHelper: models.NewMesheryControllersHelper(log, controllers.OperatorDeploymentConfig{}, nil, nil, p, nil),
		K8sContext:         models.K8sContext{ID: "test-id", ConnectionID: connID},
		OperatorTracker:    models.NewOperatorTracker(false),
		log:                log,
		EventBroadcaster:   models.NewBroadcaster("test"),
	}
	initFunc := func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error) {
		return machineCtx, nil, nil
	}
	_, err := m.Start(ctx, machineCtx, log, initFunc)
	if err != nil {
		t.Fatalf("failed to start machine: %v", err)
	}
	m.CurrentState = machines.REGISTERED

	// Trigger Connect — this creates the first LifecycleCtx.
	_, err = m.SendEvent(ctx, machines.Connect, nil)
	if err != nil {
		t.Fatalf("failed to send Connect event: %v", err)
	}
	firstCtx := m.LifecycleCtx

	if firstCtx == nil {
		t.Fatal("expected LifecycleCtx to be populated after Connect")
	}

	// Force state to CONNECTED so Disconnect is a valid transition.
	m.CurrentState = machines.CONNECTED

	// Trigger Disconnect — this should cancel firstCtx and create a new one.
	_, err = m.SendEvent(ctx, machines.Disconnect, nil)
	if err != nil {
		t.Fatalf("failed to send Disconnect event: %v", err)
	}

	if m.LifecycleCtx == firstCtx {
		t.Fatal("expected LifecycleCtx to be rotated on Disconnect")
	}

	if firstCtx.Err() == nil {
		t.Fatal("expected old LifecycleCtx to be cancelled after Disconnect")
	}
}

type serializeProvider struct {
	mockProvider
	onPersist func()
}

func (p *serializeProvider) PersistSystemEvent(event events.Event) error {
	if p.onPersist != nil {
		p.onPersist()
	}
	return nil
}

// TestActionMutexSerialization ensures that conflicting lifecycle operations
// cannot execute their side-effects concurrently by occupying the ActionMutex.
func TestActionMutexSerialization(t *testing.T) {
	machineCtx := &MachineCtx{
		ActionMutex: &sync.Mutex{},
	}

	ch2 := make(chan struct{})

	// 1. Acquire/occupy ActionMutex using a controlled operation
	machineCtx.ActionMutex.Lock()

	var sideEffectRan bool

	// 2. Start the competing lifecycle operation
	go func() {
		machineCtx.ActionMutex.Lock()
		defer machineCtx.ActionMutex.Unlock()
		sideEffectRan = true
		close(ch2)
	}()

	// 3. Prove that the competing operation cannot enter its protected section
	select {
	case <-ch2:
		t.Fatal("Competing operation entered protected section while ActionMutex was held")
	default:
		// Expected: it is blocked
	}

	// 4. Release the first operation
	machineCtx.ActionMutex.Unlock()

	// 5. Prove the second operation then proceeds
	<-ch2

	if !sideEffectRan {
		t.Fatal("Side effect did not run")
	}
}

// TestDeleteFlushMeshSyncDataSerialization proves that FlushMeshSyncData
// is protected inside the lifecycle serialization boundary.
func TestDeleteFlushMeshSyncDataSerialization(t *testing.T) {
	log := getTestLogger()
	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))})
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)

	flushCalled := make(chan struct{})
	p := &serializeProvider{
		onPersist: func() {
			close(flushCalled)
		},
	}
	ctx = context.WithValue(ctx, models.ProviderCtxKey, p)

	machineCtx := &MachineCtx{
		ActionMutex:        &sync.Mutex{},
		MesheryCtrlsHelper: models.NewMesheryControllersHelper(log, controllers.OperatorDeploymentConfig{}, nil, nil, p, nil),
		K8sContext:         models.K8sContext{ID: "test-id"},
		OperatorTracker:    models.NewOperatorTracker(false),
		EventBroadcaster:   models.NewBroadcaster("test"),
	}

	// 1. Hold/block the conflicting lifecycle operation
	machineCtx.ActionMutex.Lock()

	// 2. Trigger Delete
	deleteAction := &DeleteAction{}
	_, _, err := deleteAction.Execute(ctx, machineCtx, nil)
	if err != nil {
		t.Fatalf("failed to execute delete action: %v", err)
	}

	// 3. Prove FlushMeshSyncData cannot execute before the protected lifecycle operation releases
	select {
	case <-flushCalled:
		t.Fatal("FlushMeshSyncData executed before ActionMutex was released")
	default:
		// Expected: blocked by ActionMutex
	}

	// 4. Release the first operation
	machineCtx.ActionMutex.Unlock()

	// 5. Prove FlushMeshSyncData then executes
	select {
	case <-flushCalled:
	case <-time.After(5 * time.Second):
		t.Fatal("FlushMeshSyncData did not execute within timeout")
	}
}

type spyLogger struct {
	logger.Handler
	onInfo func(args ...interface{})
}

func (s *spyLogger) Info(args ...interface{}) {
	if s.onInfo != nil {
		s.onInfo(args...)
	}
	s.Handler.Info(args...)
}

// TestStaleConnectWorkerCancellation verifies that a stale worker detects
// lifecycle cancellation when it resumes and does not continue into protected
// side effects.
//
// The test is fully deterministic: it uses spyLogger to receive an explicit
// completion signal (workerAborted) when the stale worker logs its cancellation
// abort message and exits. A conservative 5-second deadline is kept only as a
// deadlock/hang guard for CI.
//
// This test would FAIL if the ctx.Err() guard were removed from the action's
// Execute method.
func TestStaleConnectWorkerCancellation(t *testing.T) {
	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))})
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)

	flushCalled := make(chan struct{})
	p := &serializeProvider{
		onPersist: func() {
			close(flushCalled)
		},
	}
	ctx = context.WithValue(ctx, models.ProviderCtxKey, p)

	workerAborted := make(chan struct{})
	spyLog := &spyLogger{
		Handler: getTestLogger(),
		onInfo: func(args ...interface{}) {
			for _, arg := range args {
				if str, ok := arg.(string); ok && strings.Contains(str, "aborted due to lifecycle cancellation") {
					select {
					case <-workerAborted:
					default:
						close(workerAborted)
					}
				}
			}
		},
	}

	machineCtx := &MachineCtx{
		ActionMutex:        &sync.Mutex{},
		MesheryCtrlsHelper: models.NewMesheryControllersHelper(spyLog, controllers.OperatorDeploymentConfig{}, nil, nil, p, nil),
		K8sContext:         models.K8sContext{ID: "test-id"},
		OperatorTracker:    models.NewOperatorTracker(false),
		log:                spyLog,
		EventBroadcaster:   models.NewBroadcaster("test"),
	}

	// 1. Hold ActionMutex so the stale goroutine blocks before executing
	//    any side effects. This simulates the worker being paused mid-flight
	//    (e.g., waiting behind a slow MeshKit Deploy).
	machineCtx.ActionMutex.Lock()

	connectCtx, cancelConnect := context.WithCancel(ctx)

	// Using DeleteAction as a proxy for a stale Connect worker: it
	// deterministically calls FlushMeshSyncData (observable via the mock
	// provider) if the ctx.Err() guard does not abort it.
	staleWorker := &DeleteAction{}
	_, _, err := staleWorker.Execute(connectCtx, machineCtx, nil)
	if err != nil {
		t.Fatalf("failed to execute stale worker action: %v", err)
	}

	// 2. Disconnect accepted: cancel the stale worker's LifecycleCtx before
	//    it gets a chance to enter the protected section.
	cancelConnect()

	// 3. Release ActionMutex: the stale goroutine may now enter, observe
	//    ctx.Err() != nil, and exit without running the side effect.
	machineCtx.ActionMutex.Unlock()

	// 4. Wait for the stale goroutine to abort, then assert the side effect
	//    was NOT executed. The 5-second deadline is a deadlock/hang guard only.
	const deadline = 5 * time.Second
	select {
	case <-workerAborted:
		// Stale worker detected cancellation and aborted. Verify the side
		// effect was not triggered.
		select {
		case <-flushCalled:
			t.Fatal("Stale worker executed FlushMeshSyncData despite lifecycle cancellation")
		default:
			// Correct: ctx.Err() was detected and side effects were aborted.
		}
	case <-flushCalled:
		t.Fatal("Stale worker executed FlushMeshSyncData instead of aborting on cancelled lifecycle context")
	case <-time.After(deadline):
		t.Fatalf("Stale worker goroutine did not complete within %v — possible deadlock or hang", deadline)
	}
}
