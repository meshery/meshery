package handlers

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

func TestDeleteCleanupCannotRemoveReconnectedStateMachine(t *testing.T) {
	provider := &lifecycleTestMockProvider{}
	h, tracker := newLifecycleTestHandler(t, provider)

	connectionID := uuid.Must(uuid.NewV4())

	sysID := uuid.Must(uuid.NewV4())

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: uuid.Must(uuid.NewV4())})
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	_ = context.WithValue(ctx, models.ProviderCtxKey, provider)

	sm, err := kubernetes.New(connectionID.String(), core.Uuid(connectionID), h.log)
	if err != nil {
		t.Fatalf("failed to create state machine: %v", err)
	}
	sm.Provider = provider

	machineCtx := &kubernetes.MachineCtx{
		K8sContext: models.K8sContext{
			ID:           "test-k8s-context",
			ConnectionID: connectionID.String(),
		},
	}
	sm.Context = machineCtx

	tracker.Add(core.Uuid(connectionID), sm)
	inst := sm

	// Inject a custom DeleteAction that allows us to pause and resume the cleanup.
	pauseCleanup := make(chan struct{})
	resumeCleanup := make(chan struct{})

	delState := inst.States[machines.DELETED]
	delState.Events = machines.Events{
		machines.Connect: machines.CONNECTED,
	}
	delState.Action = &testPausedDeleteAction{
		pauseCleanup:  pauseCleanup,
		resumeCleanup: resumeCleanup,
	}
	inst.States[machines.DELETED] = delState

	// 1. Delete generation D1
	detachedCtx := context.WithoutCancel(context.Background())
	detachedCtx = context.WithValue(detachedCtx, models.UserCtxKey, &models.User{ID: uuid.Must(uuid.NewV4())})
	detachedCtx = context.WithValue(detachedCtx, models.SystemIDKey, &sysID)
	detachedCtx = context.WithValue(detachedCtx, models.ProviderCtxKey, provider)

	// Create the done channel for the Delete event
	done := make(chan struct{})

	// 2. Start Delete generation D1
	deleteEvent, _ := inst.SendEvent(detachedCtx, machines.Delete, done)
	if deleteEvent == nil {
		t.Log("deleteEvent is nil")
	}

	deleteGenerationCtx := inst.GetLifecycleCtx()

	// Mimic the cleanup goroutine logic from contexts_handler.go / connections_handlers.go
	cleanupFinished := make(chan struct{})
	go func() {
		<-done // wait for Action's goroutine to finish

		// Wait for test to pause us right BEFORE RemoveIfMatch
		<-pauseCleanup

		// New implementation:
		tracker.RemoveIfMatchAndGeneration(connectionID, inst, deleteGenerationCtx)
		close(cleanupFinished)
	}()

	// Wait for the action to reach its pause point (it closed `done`)
	// In testPausedDeleteAction, we close `done` then block on resumeCleanup.

	// Ensure the Delete transition has actually occurred
	if inst.GetCurrentState() != machines.DELETED {
		t.Fatalf("expected state DELETED, got %s", inst.GetCurrentState())
	}

	// 3. Reuse M1 for CONNECT generation C2
	type key string
	connectCtx := context.WithValue(detachedCtx, key("generation"), 2)
	_ = context.WithValue(connectCtx, models.TokenCtxKey, "test-token")

	// 3. Override CONNECTED state action to avoid panics from complex dependencies
	connState := inst.States[machines.CONNECTED]
	connState.Action = &testPausedDeleteAction{}
	inst.States[machines.CONNECTED] = connState

	// 4. Send Connect
	connectCtx = context.WithValue(detachedCtx, key("generation"), 2)
	_, err = inst.SendEvent(connectCtx, machines.Connect, nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	if inst.GetCurrentState() != machines.CONNECTED {
		t.Fatalf("expected state CONNECTED, got %s", inst.GetCurrentState())
	}

	// 5. Verify LifecycleCtx changed
	if inst.GetLifecycleCtx() == deleteGenerationCtx {
		t.Fatalf("expected LifecycleCtx to change")
	}

	// Signal the cleanup goroutine that it can evaluate generation (it is waiting on pauseCleanup)
	close(pauseCleanup)

	// 6. Resume old D1 cleanup
	close(resumeCleanup)

	// Wait for cleanup to finish
	<-cleanupFinished

	// 7. Assert tracker[id] STILL == M1
	if trackedInst, ok := tracker.Get(connectionID); !ok || trackedInst != inst {
		t.Fatalf("Tracker lost the active StateMachine due to TOCTOU bug")
	}
}

type testPausedDeleteAction struct {
	pauseCleanup  chan struct{}
	resumeCleanup chan struct{}
}

func (a *testPausedDeleteAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (a *testPausedDeleteAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	go func() {
		if done, ok := data.(chan struct{}); ok && done != nil {
			close(done) // This unblocks `<-done` in the cleanup goroutine
		}
		<-a.resumeCleanup
	}()
	return machines.NoOp, nil, nil
}

func (a *testPausedDeleteAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
