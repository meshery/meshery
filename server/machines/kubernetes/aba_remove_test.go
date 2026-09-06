package kubernetes_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/core"
)

func TestDeleteABARemoval(t *testing.T) {
	log, _ := logger.New("meshery", logger.Options{Format: logger.SyslogLogFormat, LogLevel: 5})
	tracker := &machines.ConnectionToStateMachineInstanceTracker{}

	idStr := uuid.Must(uuid.NewV4()).String()
	id, _ := uuid.FromString(idStr)
	coreId := core.Uuid(id)

	// Create inst1
	inst1, _ := tracker.GetOrInitialize(coreId, func() (*machines.StateMachine, error) {
		sm, _ := kubernetes.New(idStr, coreId, log)
		return sm, nil
	})

	machineCtx := &kubernetes.MachineCtx{}
	inst1.Context = machineCtx

	_, err := inst1.Start(context.Background(), machineCtx, log, nil)
	if err != nil {
		t.Fatalf("failed to start: %v", err)
	}
	inst1.LifecycleCtx, inst1.CancelLifecycle = context.WithCancel(context.Background())

	done := make(chan struct{})
	deleteGenerationCtx := inst1.GetLifecycleCtx()

	// Simulate background cleanup logic
	go func() {
		<-done
		if inst1.GetLifecycleCtx() == deleteGenerationCtx {
			tracker.RemoveIfMatch(coreId, inst1)
		}
	}()

	// Simulate tracker removal (e.g., via CancelConnectionRegister)
	tracker.Remove(coreId)

	// Connect creates inst2
	inst2, _ := tracker.GetOrInitialize(coreId, func() (*machines.StateMachine, error) {
		sm, _ := kubernetes.New(idStr, coreId, log)
		return sm, nil
	})

	if inst1 == inst2 {
		t.Fatalf("inst2 should be different from inst1")
	}

	// Wait for cleanup to finish
	close(done)
	time.Sleep(100 * time.Millisecond)

	// Check if tracker still has inst2
	if current, ok := tracker.Get(coreId); !ok || current != inst2 {
		t.Fatalf("ABA bug reproduced: inst2 was removed by inst1's cleanup!")
	}
}

func TestConcurrentRemoveIfMatch(t *testing.T) {
	log, _ := logger.New("meshery", logger.Options{Format: logger.SyslogLogFormat, LogLevel: 5})
	tracker := &machines.ConnectionToStateMachineInstanceTracker{}

	idStr := uuid.Must(uuid.NewV4()).String()
	id, _ := uuid.FromString(idStr)
	coreId := core.Uuid(id)

	inst1, _ := tracker.GetOrInitialize(coreId, func() (*machines.StateMachine, error) {
		sm, _ := kubernetes.New(idStr, coreId, log)
		return sm, nil
	})

	// Concurrent RemoveIfMatch calls
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tracker.RemoveIfMatch(coreId, inst1)
		}()
	}

	wg.Wait()

	if _, ok := tracker.Get(coreId); ok {
		t.Fatalf("inst1 should have been removed")
	}
}
