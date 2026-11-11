package machines

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// stubAction is a configurable machines.Action used to exercise SendEvent's
// transition loop without any real connection side effects. ExecuteOnEntry and
// ExecuteOnExit always succeed with a nil event (mirroring the real k8s actions,
// which only emit events on failure); Execute returns whatever the test wires
// up so we can simulate both success and a failing action that hands back a
// non-NoOp recovery event.
type stubAction struct {
	execNext  EventType
	execEvent *events.Event
	execErr   error
	execCount int
}

func (a *stubAction) ExecuteOnEntry(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

func (a *stubAction) Execute(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
	a.execCount++
	return a.execNext, a.execEvent, a.execErr
}

func (a *stubAction) ExecuteOnExit(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

// fakeProvider is a partial models.Provider implementation: it embeds the
// interface (so the type satisfies models.Provider) but only implements the two
// methods SendEvent calls on the status-update path. Any other method would
// panic on the nil embedded interface, which keeps the test honest about what
// SendEvent actually touches.
type fakeProvider struct {
	models.Provider

	// status is the currently "persisted" connection status. UpdateConnectionById
	// mutates it so a second SendEvent observes the state left by the first,
	// exercising the de-duplication gate.
	status connections.ConnectionStatus

	getErr    error
	updateErr error

	updateCalls  int
	lastUpdateTo connections.ConnectionStatus
}

func (f *fakeProvider) GetConnectionByID(_ string, connectionID core.Uuid) (*connections.Connection, int, error) {
	if f.getErr != nil {
		return nil, http.StatusInternalServerError, f.getErr
	}
	return &connections.Connection{
		ID:     connectionID,
		Kind:   "kubernetes",
		Status: f.status,
	}, http.StatusOK, nil
}

func (f *fakeProvider) UpdateConnectionById(_ string, conn *connections.ConnectionPayload, _ string) (*connections.Connection, error) {
	f.updateCalls++
	f.lastUpdateTo = conn.Status
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.status = conn.Status
	return &connections.Connection{
		ID:     conn.ID,
		Kind:   conn.Kind,
		Status: conn.Status,
	}, nil
}

// newTestMachine builds a minimal three-state machine:
//
//	initialized --Discovery--> discovered --NotFound--> not found
//
// The discovered state's action is supplied by the test so it can succeed or
// fail. The "not found" recovery state uses a no-op action that returns a nil
// event, which is exactly the overwrite that used to mask the original failure.
func newTestMachine(t *testing.T, provider models.Provider, discoveredAction Action) *StateMachine {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build test logger: %v", err)
	}
	connectionID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate connection UUID: %v", err)
	}
	return &StateMachine{
		ID:            core.Uuid(connectionID),
		Name:          "kubernetes",
		InitialState:  InitialState,
		CurrentState:  InitialState,
		PreviousState: DefaultState,
		Log:           log,
		Provider:      provider,
		States: States{
			InitialState: State{
				Events: Events{Discovery: DISCOVERED},
				Action: nil,
			},
			DISCOVERED: State{
				Events: Events{NotFound: NOTFOUND, Register: REGISTERED},
				Action: discoveredAction,
			},
			NOTFOUND: State{
				Events: Events{Discovery: DISCOVERED, Delete: DELETED},
				Action: &stubAction{execNext: NoOp},
			},
		},
	}
}

func newTestContext(t *testing.T) context.Context {
	t.Helper()
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user UUID: %v", err)
	}
	sysID := core.Uuid(uuid.Nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: core.Uuid(userID)})
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")
	return ctx
}

// TestSendEvent_SuccessReturnsInformationalEvent covers case (a): a clean
// transition returns a nil error and an Informational confirmation event, and
// the connection status is persisted to the state the machine settled in.
func TestSendEvent_SuccessReturnsInformationalEvent(t *testing.T) {
	provider := &fakeProvider{status: connections.CONNECTED}
	// A discovered action that returns NoOp with no error is a success: the loop
	// settles in the DISCOVERED state.
	sm := newTestMachine(t, provider, &stubAction{execNext: NoOp})
	ctx := newTestContext(t)

	event, err := sm.SendEvent(ctx, Discovery, nil)
	if err != nil {
		t.Fatalf("expected nil error on a successful transition, got %v", err)
	}
	if event == nil {
		t.Fatal("expected a confirmation event on success, got nil")
	}
	if event.Severity != events.Informational {
		t.Fatalf("expected Informational severity on success, got %q", event.Severity)
	}
	if sm.CurrentState != DISCOVERED {
		t.Fatalf("expected machine to settle in %q, got %q", DISCOVERED, sm.CurrentState)
	}
	if provider.lastUpdateTo != connections.ConnectionStatus(DISCOVERED) {
		t.Fatalf("expected connection status persisted as %q, got %q", DISCOVERED, provider.lastUpdateTo)
	}
}

// TestSendEvent_FailureReturnsOriginalErrorEvent covers cases (b) and (c): an
// action that fails with a NotFound recovery transition must return the
// ORIGINAL Error event and the non-nil error (not the Informational fallback
// nor a nil error), and the connection status must still be updated to the
// recovery state (NOTFOUND).
func TestSendEvent_FailureReturnsOriginalErrorEvent(t *testing.T) {
	provider := &fakeProvider{status: connections.DISCOVERED}

	failErr := errors.New("unable to ping kubernetes context")
	failEvent := events.NewEvent().
		WithSeverity(events.Error).
		WithCategory("connection").
		WithAction("register").
		WithDescription("Unable to ping kubernetes context test at https://127.0.0.1:1").
		Build()

	// The discovered action fails and asks the machine to recover into NOTFOUND,
	// exactly like DiscoverAction.Execute / RegisterAction.Execute do on a
	// reachability failure.
	sm := newTestMachine(t, provider, &stubAction{
		execNext:  NotFound,
		execEvent: failEvent,
		execErr:   failErr,
	})
	ctx := newTestContext(t)

	event, err := sm.SendEvent(ctx, Discovery, nil)

	// (b) the non-nil error must propagate.
	if err == nil {
		t.Fatal("expected a non-nil error for a failing action, got nil — the failure was silently swallowed")
	}
	if !errors.Is(err, failErr) {
		t.Fatalf("expected the original action error to propagate, got %v", err)
	}

	// (b) the ORIGINAL Error event must be returned, not the Informational
	// fallback and not the nil event produced by the NotFound recovery action.
	if event == nil {
		t.Fatal("expected the original Error event, got nil")
	}
	if event != failEvent {
		t.Fatalf("expected the original Error event to be returned; got a different event %#v", event)
	}
	if event.Severity != events.Error {
		t.Fatalf("expected Error severity on the failure event, got %q", event.Severity)
	}

	// (c) the connection status must still be moved to the recovery state.
	if sm.CurrentState != NOTFOUND {
		t.Fatalf("expected machine to settle in %q, got %q", NOTFOUND, sm.CurrentState)
	}
	if provider.updateCalls == 0 {
		t.Fatal("expected the connection status to be updated even on failure")
	}
	if provider.lastUpdateTo != connections.NOTFOUND {
		t.Fatalf("expected connection status updated to %q, got %q", connections.NOTFOUND, provider.lastUpdateTo)
	}
}

// TestSendEvent_DeDuplicatesRepeatFailure guards the middleware anti-spam gate:
// K8sFSMMiddleware re-runs discovery on every request, so a cluster that stays
// down must not emit a fresh Error event each time. The first failure surfaces
// (status transitions into NOTFOUND); the second, with no status change, is
// suppressed (nil event, nil error) so callers persist/publish nothing.
func TestSendEvent_DeDuplicatesRepeatFailure(t *testing.T) {
	provider := &fakeProvider{status: connections.DISCOVERED}

	failErr := errors.New("unable to ping kubernetes context")
	failEvent := events.NewEvent().WithSeverity(events.Error).WithDescription("down").Build()

	sm := newTestMachine(t, provider, &stubAction{
		execNext:  NotFound,
		execEvent: failEvent,
		execErr:   failErr,
	})
	ctx := newTestContext(t)

	// First re-discovery: status changes DISCOVERED -> NOTFOUND, failure surfaces.
	sm.ResetState()
	event, err := sm.SendEvent(ctx, Discovery, nil)
	if err == nil || event == nil {
		t.Fatalf("expected the first failure to surface (event + error), got event=%v err=%v", event, err)
	}

	// Second re-discovery of the still-down cluster: status is already NOTFOUND,
	// so nothing new should be emitted.
	sm.ResetState()
	event, err = sm.SendEvent(ctx, Discovery, nil)
	if err != nil {
		t.Fatalf("expected repeat failure with no status change to be suppressed, got error %v", err)
	}
	if event != nil {
		t.Fatalf("expected no event on a suppressed repeat failure, got %#v", event)
	}
}

// TestSendEvent_InvalidTransitionPropagatesError covers the fatal
// transition-error path: a user-initiated event with no edge from the current
// state breaks out of the loop with an "Invalid status change requested" Error
// event, which must reach the caller with a non-nil error instead of falling
// through to the success path (which would return the event with a nil error
// and silently drop it).
func TestSendEvent_InvalidTransitionPropagatesError(t *testing.T) {
	provider := &fakeProvider{status: connections.CONNECTED}
	sm := newTestMachine(t, provider, &stubAction{execNext: NoOp})
	ctx := newTestContext(t)

	// InitialState only has a Discovery edge in the test machine, so a
	// user-initiated Register is an invalid transition on the first iteration.
	event, err := sm.SendEvent(ctx, Register, nil)
	if err == nil {
		t.Fatal("expected a non-nil error for an invalid transition, got nil — the transition error was silently swallowed")
	}
	if event == nil {
		t.Fatal("expected the invalid-transition Error event, got nil")
	}
	if event.Severity != events.Error {
		t.Fatalf("expected Error severity on the invalid-transition event, got %q", event.Severity)
	}
}

// TestSendEvent_TransitionErrorOnDiscoveryDeDuplicated guards the anti-spam
// gate for transition errors on the background Discovery path: a mid-loop
// transition error (a successful action emitting an event with no edge from the
// settled state) surfaces the first time the persisted status changes, and is
// suppressed on a repeat run with no status change — same contract as repeated
// action failures.
func TestSendEvent_TransitionErrorOnDiscoveryDeDuplicated(t *testing.T) {
	provider := &fakeProvider{status: connections.CONNECTED}

	// The discovered action succeeds but emits Connect, which has no edge from
	// DISCOVERED in the test machine, forcing a mid-loop transition error.
	sm := newTestMachine(t, provider, &stubAction{execNext: Connect})
	ctx := newTestContext(t)

	// First discovery: status changes CONNECTED -> DISCOVERED, so the transition
	// error surfaces.
	sm.ResetState()
	event, err := sm.SendEvent(ctx, Discovery, nil)
	if err == nil || event == nil {
		t.Fatalf("expected the first transition error to surface (event + error), got event=%v err=%v", event, err)
	}
	// The description must name the event that actually failed (the Connect the
	// action emitted mid-loop), not the Discovery the caller originally sent.
	wantDesc := "Invalid status change requested to connect for connection type kubernetes."
	if event.Description != wantDesc {
		t.Fatalf("expected invalid-transition description %q, got %q", wantDesc, event.Description)
	}
	if provider.lastUpdateTo != connections.DISCOVERED {
		t.Fatalf("expected connection status persisted as %q, got %q", connections.DISCOVERED, provider.lastUpdateTo)
	}

	// Repeat discovery settles in the same state with no status change: the
	// same transition error must be suppressed to avoid per-request spam.
	sm.ResetState()
	event, err = sm.SendEvent(ctx, Discovery, nil)
	if err != nil {
		t.Fatalf("expected repeat transition error with no status change to be suppressed, got error %v", err)
	}
	if event != nil {
		t.Fatalf("expected no event on a suppressed repeat transition error, got %#v", event)
	}
}

// TestSendEvent_UserInitiatedFailureAlwaysPropagates guards the boundary of the
// de-duplication gate: suppression is restricted to the background Discovery
// path. A user-initiated event (here Connect) whose action fails must surface
// its Error event + error even when the persisted status does not change, so a
// manual action never silently "succeeds" when it actually failed.
func TestSendEvent_UserInitiatedFailureAlwaysPropagates(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build test logger: %v", err)
	}
	connectionID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate connection UUID: %v", err)
	}

	// The connection is ALREADY persisted as NOTFOUND, so the failing Connect
	// below settles back in NOTFOUND with no status change.
	provider := &fakeProvider{status: connections.NOTFOUND}

	failErr := errors.New("connection unreachable")
	failEvent := events.NewEvent().WithSeverity(events.Error).WithDescription("connect failed").Build()

	sm := &StateMachine{
		ID:            core.Uuid(connectionID),
		Name:          "kubernetes",
		InitialState:  InitialState,
		CurrentState:  InitialState,
		PreviousState: DefaultState,
		Log:           log,
		Provider:      provider,
		States: States{
			InitialState: State{
				Events: Events{Connect: CONNECTED},
				Action: nil,
			},
			CONNECTED: State{
				Events: Events{NotFound: NOTFOUND, Disconnect: DISCONNECTED},
				Action: &stubAction{execNext: NotFound, execEvent: failEvent, execErr: failErr},
			},
			NOTFOUND: State{
				Events: Events{Discovery: DISCOVERED, Delete: DELETED},
				Action: &stubAction{execNext: NoOp},
			},
		},
	}
	ctx := newTestContext(t)

	// A user-initiated Connect that fails, with the connection already NOTFOUND.
	event, err := sm.SendEvent(ctx, Connect, nil)
	if err == nil {
		t.Fatal("expected a user-initiated failure to propagate even with no status change, got nil error")
	}
	if !errors.Is(err, failErr) {
		t.Fatalf("expected the original action error, got %v", err)
	}
	if event != failEvent {
		t.Fatalf("expected the original Error event, got %#v", event)
	}
}

// TestSendEvent_NoOpHaltOnDiscoveryIsDeduplicated covers issue #20818: an action
// that fails while reporting NoOp (an error but no recovery event) — exactly what
// DiscoverAction.Execute does when GetMachineCtx fails for a broken/zombie
// connection — must NOT re-persist or re-broadcast an Error event on every
// background Discovery re-drive. The failure halts the machine in place (no
// transition, no status write), so on the background path it is suppressed every
// time rather than returning inline and spamming the notification center.
func TestSendEvent_NoOpHaltOnDiscoveryIsDeduplicated(t *testing.T) {
	provider := &fakeProvider{status: connections.CONNECTED}

	failErr := errors.New("nil interface cannot be type casted")
	// A NoOp failure hands back an Error event but no recovery event (execNext=NoOp).
	failEvent := events.NewEvent().
		WithSeverity(events.Error).
		WithDescription("Failed to interact with the connection.").
		Build()
	sm := newTestMachine(t, provider, &stubAction{execNext: NoOp, execEvent: failEvent, execErr: failErr})
	ctx := newTestContext(t)

	// The middleware re-runs ResetState()+SendEvent(Discovery) on every request;
	// none of them should surface an event/error for a permanently broken machine.
	for i := 0; i < 3; i++ {
		sm.ResetState()
		event, err := sm.SendEvent(ctx, Discovery, nil)
		if err != nil {
			t.Fatalf("run %d: expected a NoOp-halt Discovery failure to be suppressed, got error %v", i, err)
		}
		if event != nil {
			t.Fatalf("run %d: expected no event on a suppressed NoOp-halt Discovery failure, got %#v", i, event)
		}
	}

	// The machine must not have advanced past the reset InitialState, and the
	// persisted status must be left untouched — writing the un-advanced
	// CurrentState (InitialState) would corrupt CONNECTED into "initialized".
	if provider.updateCalls != 0 {
		t.Fatalf("expected no status write for a NoOp-halt failure, got %d UpdateConnectionById call(s)", provider.updateCalls)
	}
	if provider.status != connections.CONNECTED {
		t.Fatalf("expected persisted status to remain CONNECTED, got %q", provider.status)
	}
	if sm.CurrentState != InitialState {
		t.Fatalf("expected the machine to halt without transitioning (stay at %q), got %q", InitialState, sm.CurrentState)
	}
}

// TestSendEvent_NoOpHaltOnUserEventPropagates guards the boundary of the #20818
// fix: suppression is restricted to the background Discovery path. A
// user-initiated event whose action fails with NoOp must still surface its Error
// event + error (so a manual action never silently "succeeds"), and must not
// write a status.
func TestSendEvent_NoOpHaltOnUserEventPropagates(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build test logger: %v", err)
	}
	connectionID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate connection UUID: %v", err)
	}

	provider := &fakeProvider{status: connections.DISCOVERED}
	failErr := errors.New("registration failed")
	failEvent := events.NewEvent().WithSeverity(events.Error).WithDescription("register failed").Build()

	// InitialState --Register--> REGISTERED, whose action fails with NoOp.
	sm := &StateMachine{
		ID:            core.Uuid(connectionID),
		Name:          "kubernetes",
		InitialState:  InitialState,
		CurrentState:  InitialState,
		PreviousState: DefaultState,
		Log:           log,
		Provider:      provider,
		States: States{
			InitialState: State{
				Events: Events{Register: REGISTERED},
				Action: nil,
			},
			REGISTERED: State{
				Events: Events{Connect: CONNECTED, Ignore: IGNORED},
				Action: &stubAction{execNext: NoOp, execEvent: failEvent, execErr: failErr},
			},
		},
	}
	ctx := newTestContext(t)

	event, err := sm.SendEvent(ctx, Register, nil)
	if err == nil {
		t.Fatal("expected a user-initiated NoOp-halt failure to propagate, got nil error")
	}
	if !errors.Is(err, failErr) {
		t.Fatalf("expected the original action error to propagate, got %v", err)
	}
	if event == nil || event.Severity != events.Error {
		t.Fatalf("expected the original Error event to propagate, got %#v", event)
	}
	if provider.updateCalls != 0 {
		t.Fatalf("expected no status write on a NoOp-halt failure, got %d UpdateConnectionById call(s)", provider.updateCalls)
	}
}

// failingExitAction fails in ExecuteOnExit with a nil event — matching every
// real ExecuteOnExit in the tree, all of which currently return a nil event (and
// today, a nil error). Entry/Execute succeed so the test isolates the exit path.
type failingExitAction struct {
	exitErr error
}

func (a *failingExitAction) ExecuteOnEntry(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

func (a *failingExitAction) Execute(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

func (a *failingExitAction) ExecuteOnExit(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, a.exitErr
}

// TestSendEvent_ExitActionFailureHaltsWithoutTransition covers the sibling of
// the #20818 bypass: an exit action failing aborts the transition before the
// machine advances, so like a NoOp halt it must route through the de-dup rather
// than returning inline. Two guarantees are asserted:
//
//  1. No status is written — the machine never advanced, and persisting the
//     un-advanced CurrentState would corrupt the connection's status.
//  2. A user-initiated event still surfaces a NON-NIL Error event with its
//     error. The old inline return handed back the action's event verbatim,
//     which every ExecuteOnExit in the tree builds as nil — and callers such as
//     K8sFSMMiddleware dereference (*event) whenever the error is non-nil, so
//     that path was a latent nil-pointer panic.
func TestSendEvent_ExitActionFailureHaltsWithoutTransition(t *testing.T) {
	provider := &fakeProvider{status: connections.CONNECTED}
	exitErr := errors.New("exit action failed")

	sm := newTestMachine(t, provider, &stubAction{execNext: NoOp})
	// The exit action that runs is the one on the CURRENT state, so attach the
	// failing action to InitialState and drive a transition out of it.
	sm.States[InitialState] = State{
		Events: Events{Discovery: DISCOVERED, Register: REGISTERED},
		Action: &failingExitAction{exitErr: exitErr},
	}
	sm.States[REGISTERED] = State{
		Events: Events{Connect: CONNECTED},
		Action: &stubAction{execNext: NoOp},
	}
	ctx := newTestContext(t)

	// User-initiated: the failure must propagate, with a usable event.
	event, err := sm.SendEvent(ctx, Register, nil)
	if err == nil {
		t.Fatal("expected an exit-action failure on a user-initiated event to propagate, got nil error")
	}
	if !errors.Is(err, exitErr) {
		t.Fatalf("expected the original exit-action error to propagate, got %v", err)
	}
	if event == nil {
		t.Fatal("expected a non-nil Error event alongside the error (callers dereference it), got nil")
	}
	if event.Severity != events.Error {
		t.Fatalf("expected an Error-severity event, got %q", event.Severity)
	}
	if provider.updateCalls != 0 {
		t.Fatalf("expected no status write when the machine never advanced, got %d UpdateConnectionById call(s)", provider.updateCalls)
	}
	if sm.CurrentState != InitialState {
		t.Fatalf("expected the machine to halt at %q, got %q", InitialState, sm.CurrentState)
	}

	// Background discovery: the same failure repeated must be suppressed rather
	// than re-broadcast on every request.
	for i := 0; i < 3; i++ {
		sm.ResetState()
		event, err := sm.SendEvent(ctx, Discovery, nil)
		if err != nil {
			t.Fatalf("run %d: expected a repeated background exit-action failure to be suppressed, got error %v", i, err)
		}
		if event != nil {
			t.Fatalf("run %d: expected no event on a suppressed background failure, got %#v", i, event)
		}
	}
	if provider.updateCalls != 0 {
		t.Fatalf("expected no status write across background re-drives, got %d UpdateConnectionById call(s)", provider.updateCalls)
	}
	if provider.status != connections.CONNECTED {
		t.Fatalf("expected persisted status to remain CONNECTED, got %q", provider.status)
	}
}

func TestResolveUserAndSystemID(t *testing.T) {
	t.Run("empty context returns nil uuids", func(t *testing.T) {
		userUUID, sysUUID := ResolveUserAndSystemID(context.Background())
		if userUUID != core.Uuid(uuid.Nil) {
			t.Errorf("expected userUUID to be Nil, got %v", userUUID)
		}
		if sysUUID != core.Uuid(uuid.Nil) {
			t.Errorf("expected sysUUID to be Nil, got %v", sysUUID)
		}
	})

	t.Run("populated context returns correct uuids", func(t *testing.T) {
		expectedUser := core.Uuid(uuid.Must(uuid.NewV4()))
		expectedSys := core.Uuid(uuid.Must(uuid.NewV4()))
		ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: expectedUser})
		ctx = context.WithValue(ctx, models.SystemIDKey, &expectedSys)

		userUUID, sysUUID := ResolveUserAndSystemID(ctx)
		if userUUID != expectedUser {
			t.Errorf("expected userUUID %v, got %v", expectedUser, userUUID)
		}
		if sysUUID != expectedSys {
			t.Errorf("expected sysUUID %v, got %v", expectedSys, sysUUID)
		}
	})

	t.Run("typed nil pointers in context return nil uuids without panic", func(t *testing.T) {
		var nilUser *models.User
		var nilSys *core.Uuid
		ctx := context.WithValue(context.Background(), models.UserCtxKey, nilUser)
		ctx = context.WithValue(ctx, models.SystemIDKey, nilSys)

		userUUID, sysUUID := ResolveUserAndSystemID(ctx)
		if userUUID != core.Uuid(uuid.Nil) {
			t.Errorf("expected userUUID to be Nil, got %v", userUUID)
		}
		if sysUUID != core.Uuid(uuid.Nil) {
			t.Errorf("expected sysUUID to be Nil, got %v", sysUUID)
		}
	})

	t.Run("invalid types in context return nil uuids without panic", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), models.UserCtxKey, "not-a-user")
		ctx = context.WithValue(ctx, models.SystemIDKey, 12345)

		userUUID, sysUUID := ResolveUserAndSystemID(ctx)
		if userUUID != core.Uuid(uuid.Nil) {
			t.Errorf("expected userUUID to be Nil, got %v", userUUID)
		}
		if sysUUID != core.Uuid(uuid.Nil) {
			t.Errorf("expected sysUUID to be Nil, got %v", sysUUID)
		}
	})
}

func TestSendEvent_NilContextKeysDoNotPanic(t *testing.T) {
	provider := &fakeProvider{status: connections.CONNECTED}
	sm := newTestMachine(t, provider, &stubAction{execNext: NoOp})

	// 1. Invalid transition with context.Background() (no user or sysID)
	// Triggers invalidTransitionEvent which previously dereferenced *sysID and user.ID
	event, err := sm.SendEvent(context.Background(), Delete, nil)
	if err == nil {
		t.Fatal("expected error on invalid transition Delete from InitialState, got nil")
	}
	if event == nil {
		t.Fatal("expected non-nil event on invalid transition")
	}
	if event.Severity != events.Error {
		t.Errorf("expected Error severity, got %v", event.Severity)
	}

	// 2. Valid transition with context.Background()
	// Advances InitialState -> DISCOVERED, triggering status update event
	_, err = sm.SendEvent(context.Background(), Discovery, nil)
	if err != nil {
		t.Fatalf("expected nil error on valid transition, got %v", err)
	}
	if sm.CurrentState != DISCOVERED {
		t.Errorf("expected CurrentState DISCOVERED, got %v", sm.CurrentState)
	}
}

// The unreachable condition (meshery/schemas event correlation contract): a
// background Discovery failure that moves the connection to NOTFOUND reports
// it, the de-duplication of repeats is unchanged, and the Discovery that later
// moves the connection out of NOTFOUND resolves it.
func TestSendEvent_DiscoveryFailureReportsUnreachableCondition(t *testing.T) {
	provider := &fakeProvider{status: connections.CONNECTED}
	failing := &stubAction{
		execNext:  NotFound,
		execEvent: events.NewEvent().WithSeverity(events.Error).WithDescription("cluster down").Build(),
		execErr:   errors.New("dial tcp: connection refused"),
	}
	sm := newTestMachine(t, provider, failing)
	sm.ConditionComponent = "k8s-connection"
	ctx := newTestContext(t)

	event, err := sm.SendEvent(ctx, Discovery, nil)
	if err == nil || event == nil {
		t.Fatalf("got event %v, err %v; want the failure", event, err)
	}
	want := models.ConditionKey("k8s-connection", sm.ID.String(), ConditionUnreachable)
	if event.CorrelationKey == nil || *event.CorrelationKey != want {
		t.Fatalf("correlationKey = %v, want %q", event.CorrelationKey, want)
	}
	if event.ResolvedAt != nil {
		t.Fatalf("a failure carries resolvedAt %v", event.ResolvedAt)
	}

	// The connection is NOTFOUND now; the same failure again changes nothing,
	// so it is not reported to the user again (issue #20818), but it comes back
	// as a Recurrence, without an error, to count on the open condition.
	sm.ResetState()
	failing.execEvent = events.NewEvent().WithSeverity(events.Error).WithDescription("cluster down").Build()
	event, err = sm.SendEvent(ctx, Discovery, nil)
	if !IsRecurrence(event, err) {
		t.Fatalf("repeat failure = event %v, err %v; want a Recurrence", event, err)
	}
	if *event.CorrelationKey != want {
		t.Fatalf("recurrence key = %q, want %q", *event.CorrelationKey, want)
	}

	// The cluster answers again: the Discovery that leaves NOTFOUND resolves the
	// condition, under its key, with a success severity.
	sm.ResetState()
	failing.execNext, failing.execEvent, failing.execErr = NoOp, nil, nil
	event, err = sm.SendEvent(ctx, Discovery, nil)
	if err != nil || event == nil {
		t.Fatalf("got event %v, err %v; want the resolution", event, err)
	}
	if event.CorrelationKey == nil || *event.CorrelationKey != want || event.ResolvedAt == nil {
		t.Fatalf("got key %v resolvedAt %v; want a resolution of %q", event.CorrelationKey, event.ResolvedAt, want)
	}
	if event.Severity != events.Success || event.Category != "connection" {
		t.Fatalf("resolution severity %q category %q, want success and connection", event.Severity, event.Category)
	}

	// A healthy Discovery afterwards is routine: no second resolution.
	sm.ResetState()
	event, err = sm.SendEvent(ctx, Discovery, nil)
	if err != nil || event == nil {
		t.Fatalf("got event %v, err %v", event, err)
	}
	if event.ResolvedAt != nil || event.CorrelationKey != nil {
		t.Fatalf("a routine Discovery carries key %v resolvedAt %v", event.CorrelationKey, event.ResolvedAt)
	}
}

// A machine that does not opt in keeps its events uncorrelated.
func TestSendEvent_WithoutConditionComponentReportsNoCondition(t *testing.T) {
	provider := &fakeProvider{status: connections.CONNECTED}
	sm := newTestMachine(t, provider, &stubAction{
		execNext:  NotFound,
		execEvent: events.NewEvent().WithSeverity(events.Error).Build(),
		execErr:   errors.New("down"),
	})
	event, err := sm.SendEvent(newTestContext(t), Discovery, nil)
	if err == nil || event == nil {
		t.Fatalf("got event %v, err %v; want the failure", event, err)
	}
	if event.CorrelationKey != nil {
		t.Fatalf("correlationKey = %q on a machine without ConditionComponent", *event.CorrelationKey)
	}
}

// Deleting a connection that is NOTFOUND resolves its unreachable condition.
func TestSendEvent_DeleteFromNotFoundResolvesTheCondition(t *testing.T) {
	provider := &fakeProvider{status: connections.ConnectionStatus(NOTFOUND)}
	sm := newTestMachine(t, provider, &stubAction{execNext: NoOp})
	sm.ConditionComponent = "k8s-connection"
	sm.States[DELETED] = State{Events: Events{}, Action: &stubAction{execNext: NoOp}}
	sm.CurrentState = NOTFOUND

	event, err := sm.SendEvent(newTestContext(t), Delete, nil)
	if err != nil || event == nil {
		t.Fatalf("got event %v, err %v; want the resolution", event, err)
	}
	want := models.ConditionKey("k8s-connection", sm.ID.String(), ConditionUnreachable)
	if event.CorrelationKey == nil || *event.CorrelationKey != want || event.ResolvedAt == nil {
		t.Fatalf("got key %v resolvedAt %v; want a resolution of %q", event.CorrelationKey, event.ResolvedAt, want)
	}
	if IsRecurrence(event, err) {
		t.Fatal("a resolution is not a Recurrence")
	}
}

// A request that saved the new status before driving the machine names the
// old one with WithPriorStatus: a connection that left NOTFOUND resolves its
// condition although the status the machine reads is already the new one, and
// one that was not NOTFOUND resolves nothing.
func TestSendEvent_PriorStatusResolvesAManualReconnect(t *testing.T) {
	for _, c := range []struct {
		prior       connections.ConnectionStatus
		wantResolve bool
	}{
		{connections.ConnectionStatus(NOTFOUND), true},
		{connections.CONNECTED, false},
		{"", false},
	} {
		provider := &fakeProvider{status: connections.ConnectionStatus(DISCOVERED)}
		sm := newTestMachine(t, provider, &stubAction{execNext: NoOp})
		sm.ConditionComponent = "k8s-connection"
		ctx := newTestContext(t)
		if c.prior != "" {
			ctx = WithPriorStatus(ctx, c.prior)
		}
		event, err := sm.SendEvent(ctx, Discovery, nil)
		if err != nil || event == nil {
			t.Fatalf("prior %q: got event %v, err %v", c.prior, event, err)
		}
		if got := event.ResolvedAt != nil; got != c.wantResolve {
			t.Errorf("prior %q: resolved = %v, want %v", c.prior, got, c.wantResolve)
		}
	}
}

// IsRecurrence holds only for an event without an error that carries an open
// condition's key.
func TestIsRecurrence(t *testing.T) {
	key := "k"
	now := time.Now()
	open := &events.Event{CorrelationKey: &key}
	resolved := &events.Event{CorrelationKey: &key, ResolvedAt: &now}
	for _, c := range []struct {
		name  string
		event *events.Event
		err   error
		want  bool
	}{
		{"open key, no error", open, nil, true},
		{"open key with an error", open, errors.New("down"), false},
		{"resolution", resolved, nil, false},
		{"no key", &events.Event{}, nil, false},
		{"no event", nil, nil, false},
	} {
		if got := IsRecurrence(c.event, c.err); got != c.want {
			t.Errorf("%s: IsRecurrence = %v, want %v", c.name, got, c.want)
		}
	}
}

// A status update that cannot read or write the connection is a failure, and
// says so: without a severity its event used to land in the Notification
// Center as informational.
func TestSendEvent_StatusUpdateFailureIsAnError(t *testing.T) {
	for name, provider := range map[string]*fakeProvider{
		"retrieve": {status: connections.CONNECTED, getErr: errors.New("record not found")},
		"update":   {status: connections.CONNECTED, updateErr: errors.New("database is locked")},
	} {
		sm := newTestMachine(t, provider, &stubAction{execNext: NoOp})
		initialState := sm.CurrentState
		initialPrevState := sm.PreviousState

		event, err := sm.SendEvent(newTestContext(t), Discovery, nil)
		if err == nil || event == nil {
			t.Fatalf("%s: SendEvent = (%v, %v), want a failure event and an error", name, event, err)
		}
		if event.Severity != events.Error {
			t.Errorf("%s: severity = %q, want %q", name, event.Severity, events.Error)
		}
		if sm.CurrentState != initialState {
			t.Errorf("%s: currentState = %q, want rolled back state %q", name, sm.CurrentState, initialState)
		}
		if sm.PreviousState != initialPrevState {
			t.Errorf("%s: previousState = %q, want rolled back state %q", name, sm.PreviousState, initialPrevState)
		}
	}
}

// TestSendEvent_StatusUpdateFailureRollbackEnablesRetry verifies that when a
// database status update fails, in-memory state rolls back so subsequent retry
// events do not fail with ErrInvalidTransitionEvent.
func TestSendEvent_StatusUpdateFailureRollbackEnablesRetry(t *testing.T) {
	provider := &fakeProvider{status: connections.CONNECTED, updateErr: errors.New("database is locked")}
	action := &stubAction{execNext: NoOp}
	sm := newTestMachine(t, provider, action)

	// First attempt fails during status persistence.
	_, err := sm.SendEvent(newTestContext(t), Discovery, nil)
	if err == nil {
		t.Fatal("expected SendEvent to fail when DB update errors")
	}
	if sm.CurrentState != InitialState {
		t.Fatalf("expected state to roll back to %q, got %q", InitialState, sm.CurrentState)
	}
	if action.execCount != 1 {
		t.Fatalf("expected action to have executed once on first attempt, got %d", action.execCount)
	}

	// Resolve the transient database error and retry the same event.
	provider.updateErr = nil
	event, err := sm.SendEvent(newTestContext(t), Discovery, nil)
	if err != nil {
		t.Fatalf("expected retry SendEvent to succeed after DB recovery, got error: %v", err)
	}
	if event == nil || event.Severity != events.Informational {
		t.Fatalf("expected informational success event on retry, got %v", event)
	}
	if sm.CurrentState != DISCOVERED {
		t.Fatalf("expected currentState = %q after successful retry, got %q", DISCOVERED, sm.CurrentState)
	}
	if action.execCount != 2 {
		t.Fatalf("expected action to have re-executed safely on retry, got %d", action.execCount)
	}
}

// TestSendEvent_AmbiguousUpdateSucceededDoesNotRollback verifies that when
// UpdateConnectionById returns an error (e.g. read-back or transport error) but
// the database write actually succeeded, SendEvent reconciles the persisted state
// and does not roll back in-memory state.
func TestSendEvent_AmbiguousUpdateSucceededDoesNotRollback(t *testing.T) {
	provider := &fakeProvider{
		status:    connections.DISCOVERED, // simulated state persisted in database
		updateErr: errors.New("read-back failed after write"),
	}
	action := &stubAction{execNext: NoOp}
	sm := newTestMachine(t, provider, action)

	event, err := sm.SendEvent(newTestContext(t), Discovery, nil)
	if err != nil {
		t.Fatalf("expected SendEvent to reconcile and succeed without error, got %v", err)
	}
	if sm.CurrentState != DISCOVERED {
		t.Fatalf("expected currentState to remain %q, got %q", DISCOVERED, sm.CurrentState)
	}
	if event == nil || event.Severity != events.Informational {
		t.Fatalf("expected informational success event, got %v", event)
	}
}
