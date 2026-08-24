package machines

// Tests in this file verify the lifecycle-context generation semantics of
// StateMachine.SendEvent.
//
// Lifecycle-context generation model
// -----------------------------------
// Each time a transition begins executing (after ExecuteOnExit succeeds),
// SendEvent cancels the previous lifecycle context and installs a fresh one.
// This "generation" context is the one passed to ExecuteOnEntry and Execute.
//
// When a transition halts early — because the entry or execute action returns
// NoOp with a non-nil error — the new generation context (C2) is left
// installed in LifecycleCtx even though CurrentState was not advanced. The
// previous generation (C1) is cancelled at that point.
//
// Invariant under test:
//   After SendEvent returns, GetLifecycleCtx() always yields a non-nil,
//   non-cancelled context, regardless of whether the transition succeeded,
//   advanced to a recovery state, or halted without changing CurrentState.

import (
	"context"
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// haltingEntryAction returns NoOp + a non-nil error from ExecuteOnEntry.
// This simulates any action whose guard/prerequisite check fails before
// the state is accepted as active (e.g. a missing machine context).
type haltingEntryAction struct{ err error }

func (a *haltingEntryAction) ExecuteOnEntry(_ context.Context, _ interface{}, _ interface{}) (EventType, *events.Event, error) {
	ev := events.NewEvent().WithSeverity(events.Error).WithDescription("entry action failed").Build()
	return NoOp, ev, a.err
}
func (a *haltingEntryAction) Execute(_ context.Context, _ interface{}, _ interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}
func (a *haltingEntryAction) ExecuteOnExit(_ context.Context, _ interface{}, _ interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

// haltingExecAction returns NoOp + a non-nil error from Execute.
// This simulates an action whose main body fails with no recovery transition
// (e.g. a connection that cannot be reached and has no fallback state).
type haltingExecAction struct{ err error }

func (a *haltingExecAction) ExecuteOnEntry(_ context.Context, _ interface{}, _ interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}
func (a *haltingExecAction) Execute(_ context.Context, _ interface{}, _ interface{}) (EventType, *events.Event, error) {
	ev := events.NewEvent().WithSeverity(events.Error).WithDescription("execute action failed").Build()
	return NoOp, ev, a.err
}
func (a *haltingExecAction) ExecuteOnExit(_ context.Context, _ interface{}, _ interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

// succeedingAction completes all three lifecycle methods without error.
type succeedingAction struct{}

func (a *succeedingAction) ExecuteOnEntry(_ context.Context, _ interface{}, _ interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}
func (a *succeedingAction) Execute(_ context.Context, _ interface{}, _ interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}
func (a *succeedingAction) ExecuteOnExit(_ context.Context, _ interface{}, _ interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

// lcSendCtx builds the minimum context values SendEvent reads.
func lcSendCtx(t *testing.T) context.Context {
	t.Helper()
	uid, _ := uuid.NewV4()
	sysID := core.Uuid(uuid.Nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: core.Uuid(uid)})
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")
	return ctx
}

// lcNewMachine builds a two-state machine (CONNECTED ↔ DISCONNECTED) with
// Provider nil so the status-persistence block is skipped entirely. Both
// CurrentState and the installed lifecycle context must be supplied by the
// test before calling SendEvent.
func lcNewMachine(t *testing.T, connectedAction, disconnectedAction Action) *StateMachine {
	t.Helper()
	log, _ := logger.New("test", logger.Options{})
	id, _ := uuid.NewV4()
	return &StateMachine{
		ID:            core.Uuid(id),
		Name:          "kubernetes",
		InitialState:  InitialState,
		CurrentState:  CONNECTED,
		PreviousState: REGISTERED,
		Log:           log,
		Provider:      nil,
		States: States{
			CONNECTED: State{
				Events: Events{Disconnect: DISCONNECTED},
				Action: connectedAction,
			},
			DISCONNECTED: State{
				Events: Events{Connect: CONNECTED},
				Action: disconnectedAction,
			},
		},
	}
}

// installGeneration attaches a fresh, independent lifecycle context to sm and
// returns it so the test can observe its cancellation state afterward.
func installGeneration(sm *StateMachine) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sm.LifecycleCtx = ctx
	sm.CancelLifecycle = cancel
	return ctx
}

// TestLifecycleCtxRemainsValidAfterHaltedTransition verifies that when a
// transition halts without advancing CurrentState — because the entry action
// or the execute action returns NoOp with a non-nil error — SendEvent still
// leaves a valid (non-cancelled) lifecycle context in the machine.
//
// The previous generation is cancelled as part of beginning the transition
// attempt. The newly-created generation is installed unconditionally and
// remains valid even when the transition does not complete.
func TestLifecycleCtxRemainsValidAfterHaltedTransition(t *testing.T) {
	haltErr := errors.New("action failed")

	tests := []struct {
		name             string
		disconnectedAction Action
	}{
		{
			name:             "entry action halts",
			disconnectedAction: &haltingEntryAction{err: haltErr},
		},
		{
			name:             "execute action halts",
			disconnectedAction: &haltingExecAction{err: haltErr},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctx := lcSendCtx(t)
			sm := lcNewMachine(t, &succeedingAction{}, tc.disconnectedAction)
			prev := installGeneration(sm)

			_, _ = sm.SendEvent(ctx, Disconnect, nil)

			// CurrentState must not have advanced: a halted transition must
			// leave the state machine where it found it.
			if sm.CurrentState != CONNECTED {
				t.Errorf("CurrentState = %q; want CONNECTED: a halted transition must not advance state", sm.CurrentState)
			}

			// The previous generation must be cancelled — SendEvent cancels it
			// before invoking entry/execute actions, regardless of whether
			// those actions succeed.
			if prev.Err() == nil {
				t.Error("previous lifecycle generation was not cancelled when the new transition began")
			}

			// The current generation must be a distinct, non-cancelled context.
			// These are two separate facts: the old one being dead does not
			// imply the new one is alive, and vice versa.
			curr := sm.LifecycleCtx
			if curr == nil {
				t.Fatal("LifecycleCtx is nil after SendEvent")
			}
			if curr == prev {
				t.Fatal("LifecycleCtx was not rotated: still holds the previous generation")
			}
			if curr.Err() != nil {
				t.Errorf("current lifecycle generation is cancelled after a halted transition: %v", curr.Err())
			}
		})
	}
}

// TestMachineUsableAfterHaltedTransition verifies that the machine can accept
// and complete a legitimate transition after a previous one halted without
// advancing CurrentState.
//
// This guards against a scenario where a halted transition leaves the machine
// in a state from which no further progress is possible (e.g. a cancelled
// lifecycle context that blocks subsequent background workers).
func TestMachineUsableAfterHaltedTransition(t *testing.T) {
	ctx := lcSendCtx(t)
	haltErr := errors.New("transient entry failure")

	sm := lcNewMachine(t, &succeedingAction{}, &haltingEntryAction{err: haltErr})
	installGeneration(sm)

	// First send: halts, CurrentState stays CONNECTED.
	_, _ = sm.SendEvent(ctx, Disconnect, nil)
	if sm.CurrentState != CONNECTED {
		t.Fatalf("precondition: expected CONNECTED after halted transition, got %q", sm.CurrentState)
	}
	afterHalt := sm.LifecycleCtx

	// Swap in a succeeding action and retry the same transition.
	sm.States[DISCONNECTED] = State{
		Events: Events{Connect: CONNECTED},
		Action: &succeedingAction{},
	}

	_, err := sm.SendEvent(ctx, Disconnect, nil)
	if err != nil {
		t.Fatalf("SendEvent failed after recovering from a halted transition: %v", err)
	}
	if sm.CurrentState != DISCONNECTED {
		t.Errorf("CurrentState = %q; want DISCONNECTED after a successful transition", sm.CurrentState)
	}

	// The generation that was current after the halted transition must now
	// be cancelled — the successful transition rotated it out.
	if afterHalt.Err() == nil {
		t.Error("generation left by the halted transition was not cancelled by the subsequent successful transition")
	}

	// The newly-installed generation must be distinct from the previous one
	// and must itself be valid — two independent assertions.
	curr := sm.LifecycleCtx
	if curr == nil {
		t.Fatal("LifecycleCtx is nil after the successful transition")
	}
	if curr == afterHalt {
		t.Error("LifecycleCtx was not rotated by the successful transition")
	}
	if curr.Err() != nil {
		t.Errorf("current lifecycle generation is cancelled after a successful transition: %v", curr.Err())
	}
}
