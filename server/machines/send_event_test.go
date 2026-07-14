package machines

import (
	"context"
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// redirectAction mimics the real kubernetes RegisterAction: its Execute returns
// a redirect event (and optionally an error), asking the engine to move the
// machine on to another state.
type redirectAction struct {
	redirect EventType
	err      error
}

func (a *redirectAction) ExecuteOnEntry(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}
func (a *redirectAction) Execute(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
	return a.redirect, nil, a.err
}
func (a *redirectAction) ExecuteOnExit(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

// terminalAction ends a transition (always NoOp).
type terminalAction struct{}

func (a *terminalAction) ExecuteOnEntry(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}
func (a *terminalAction) Execute(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}
func (a *terminalAction) ExecuteOnExit(context.Context, interface{}, interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

// fakeProvider records the status persisted via UpdateConnectionById. Only the
// two methods SendEvent calls are implemented; the rest are inherited from the
// embedded (nil) interface and are never invoked on these paths.
type fakeProvider struct {
	models.Provider
	updateCalled  bool
	updatedStatus connections.ConnectionStatus
}

func (f *fakeProvider) GetConnectionByID(string, core.Uuid) (*connections.Connection, int, error) {
	return &connections.Connection{Kind: "test"}, 200, nil
}
func (f *fakeProvider) UpdateConnectionById(_ string, conn *connections.ConnectionPayload, _ string) (*connections.Connection, error) {
	f.updateCalled = true
	f.updatedStatus = conn.Status
	return &connections.Connection{ID: conn.ID, Kind: conn.Kind}, nil
}

func testCtx() context.Context {
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: core.Uuid(uuid.Must(uuid.NewV4()))})
	return context.WithValue(ctx, models.SystemIDKey, &sysID)
}

func newTestLogger(t *testing.T) logger.Handler {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	return log
}

const (
	stateStart StateType = "start"
	stateMid   StateType = "mid"
	stateFinal StateType = "final"

	eventGo       EventType = "go"
	eventFallback EventType = "fallback"
)

// TestSendEvent_FollowsRedirectToDefinedEdge covers the fix for #20642: when an
// action redirects (as RegisterAction does to NOTFOUND on ping failure) and the
// current state declares that edge, the machine must follow it and persist the
// resolved final state — not stall in the intermediate state.
func TestSendEvent_FollowsRedirectToDefinedEdge(t *testing.T) {
	provider := &fakeProvider{}
	sm := &StateMachine{
		ID:           core.Uuid(uuid.Must(uuid.NewV4())),
		Name:         "test",
		CurrentState: stateStart,
		Log:          newTestLogger(t),
		Provider:     provider,
		States: States{
			stateStart: {Events: Events{eventGo: stateMid}},
			// mid redirects onward to final, and declares the edge for it.
			stateMid:   {Events: Events{eventFallback: stateFinal}, Action: &redirectAction{redirect: eventFallback, err: errors.New("ping failed")}},
			stateFinal: {Events: Events{}, Action: &terminalAction{}},
		},
	}

	if _, err := sm.SendEvent(testCtx(), eventGo, nil); err != nil {
		t.Fatalf("SendEvent returned error: %v", err)
	}
	if sm.CurrentState != stateFinal {
		t.Fatalf("machine state = %q, want %q (redirect not followed)", sm.CurrentState, stateFinal)
	}
	if !provider.updateCalled || provider.updatedStatus != connections.ConnectionStatus(stateFinal) {
		t.Fatalf("persisted status = %q (called=%v), want %q", provider.updatedStatus, provider.updateCalled, stateFinal)
	}
}

// TestSendEvent_DoesNotPersistPartialTransition covers the second half of
// #20642: when an action redirects to an event the current state has no edge
// for, the transition is invalid and the machine must NOT persist the
// intermediate state it was passing through.
func TestSendEvent_DoesNotPersistPartialTransition(t *testing.T) {
	provider := &fakeProvider{}
	sm := &StateMachine{
		ID:           core.Uuid(uuid.Must(uuid.NewV4())),
		Name:         "test",
		CurrentState: stateStart,
		Log:          newTestLogger(t),
		Provider:     provider,
		States: States{
			stateStart: {Events: Events{eventGo: stateMid}},
			// mid redirects to fallback but declares NO edge for it -> invalid.
			stateMid: {Events: Events{}, Action: &redirectAction{redirect: eventFallback}},
		},
	}

	event, _ := sm.SendEvent(testCtx(), eventGo, nil)
	if provider.updateCalled {
		t.Fatalf("status was persisted (%q) on an invalid transition; expected no persistence", provider.updatedStatus)
	}
	if event == nil || event.Severity != events.Error {
		t.Fatalf("expected an error-severity event on an invalid transition, got %#v", event)
	}
}
