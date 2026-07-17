package machines

import (
	"context"
	"errors"
	"net/http"
	"testing"

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
}

func (a *stubAction) ExecuteOnEntry(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

func (a *stubAction) Execute(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
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
	sysID := core.Uuid(uuid.FromStringOrNil("00000000-0000-0000-0000-000000000000"))
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
