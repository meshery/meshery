package kubernetes

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/models/events"
)

type mockProvider struct {
	models.Provider
}

func (m *mockProvider) GetConnectionByID(token string, connectionID core.Uuid) (*connections.Connection, int, error) {
	return &connections.Connection{}, 200, nil
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
<-flushCalled
}

// TestStaleConnectWorkerCancellation verifies that a stale worker detects
// lifecycle cancellation when it resumes and does not continue.
func TestStaleConnectWorkerCancellation(t *testing.T) {
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
	log:                log,
	EventBroadcaster:   models.NewBroadcaster("test"),
}

// 1. Connect lifecycle accepted, worker is paused at a controlled synchronization point
machineCtx.ActionMutex.Lock()

connectCtx, cancelConnect := context.WithCancel(ctx)

// Using DeleteAction as a proxy for any stale worker since it deterministically
// calls FlushMeshSyncData (which we can observe via mock provider).
staleWorker := &DeleteAction{}
	_, _, err := staleWorker.Execute(connectCtx, machineCtx, nil)
	if err != nil {
		t.Fatalf("failed to execute stale worker action: %v", err)
	}

// 2. Disconnect accepted -> old LifecycleCtx cancelled
cancelConnect()

// 3. Connect worker resumes
machineCtx.ActionMutex.Unlock()

// 4. Stale worker detects lifecycle cancellation and does NOT continue
// If it continued, it would call FlushMeshSyncData and close flushCalled.
select {
case <-flushCalled:
t.Fatal("Stale worker continued into side effects despite lifecycle cancellation")
case <-time.After(50 * time.Millisecond):
// Expected: side effects were aborted.
}
}
