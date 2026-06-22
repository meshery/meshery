package helpers

import (
	"context"
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

func TestInitializeMachineWithContext_ZombieProcess(t *testing.T) {
	log, err := logger.New("meshery", logger.Options{Format: logger.TerminalLogFormat})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}

	id := uuid.Must(uuid.NewV4())
	userID := uuid.Must(uuid.NewV4())

	// An initialization function that fails simulates a bad .Start() outcome.
	failingInitFunc := func(ctx context.Context, machineCtx interface{}, _ logger.Handler) (interface{}, *events.Event, error) {
		return nil, nil, errors.New("simulated manual initialization failure")
	}

	// Use a concrete available provider default instead of an uninitialized embedded interface mock
	provider := &models.DefaultLocalProvider{}

	inst, err := InitializeMachineWithContext(
		nil,
		context.Background(),
		id,
		userID,
		tracker,
		log,
		provider,
		machines.InitialState,
		"kubernetes",
		failingInitFunc,
	)

	// We expect InitializeMachineWithContext to fail.
	if err == nil {
		t.Fatalf("expected error from InitializeMachineWithContext due to failingInitFunc, got nil. Check: %v", inst)
	}

	// The BUG was that the failed instance was added to the tracker before checking the error.
	_, ok := tracker.Get(id)
	if ok {
		t.Fatalf("ZOMBIE STATE MACHINE CAUGHT: The failing state machine instance was stored in the tracker before returning!")
	}
}

func TestInitializeMachineWithContext_SuccessPath(t *testing.T) {
	log, err := logger.New("meshery", logger.Options{Format: logger.TerminalLogFormat})
	if err != nil {
		t.Fatalf("Failed to initialize logger: %v", err)
	}

	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}

	id := uuid.Must(uuid.NewV4())
	userID := uuid.Must(uuid.NewV4())

	// An initialization function that succeeds simulates a good .Start() outcome.
	successInitFunc := func(ctx context.Context, machineCtx interface{}, _ logger.Handler) (interface{}, *events.Event, error) {
		return nil, nil, nil
	}

	provider := &models.DefaultLocalProvider{}

	_, err = InitializeMachineWithContext(
		nil,
		context.Background(),
		id,
		userID,
		tracker,
		log,
		provider,
		machines.InitialState,
		"kubernetes",
		successInitFunc,
	)

	if err != nil {
		t.Fatalf("expected nil from InitializeMachineWithContext, got error: %v", err)
	}

	// Verify the successful instance was successfully added to the tracker.
	_, ok := tracker.Get(id)
	if !ok {
		t.Fatalf("MISSING STATE MACHINE: The successful state machine instance was not stored in the tracker cache!")
	}
}
