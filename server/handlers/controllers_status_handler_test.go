package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/core"
)

// trackerWith builds an instance tracker holding a single connection->machine
// mapping, for exercising machineCtxForConnection without standing up real FSMs.
func trackerWith(id core.Uuid, inst *machines.StateMachine) *machines.ConnectionToStateMachineInstanceTracker {
	return &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: map[core.Uuid]*machines.StateMachine{id: inst},
	}
}

// machineCtxForConnection must treat an unassigned Context on a tracked machine
// - a non-kubernetes connection, or a kubernetes one whose cluster was
// unreachable when the machine was created - as an expected "not ready" state
// (nil,false) and, crucially, must NOT type-assert that nil Context, which
// previously logged meshkit-11180 on every controller-status poll. Missing,
// nil, typed-nil, wrongly-typed and half-built contexts are all not-ready;
// only a fully-initialized kubernetes machine (a *kubernetes.MachineCtx
// carrying a controllers helper) resolves as ready.
func TestMachineCtxForConnection(t *testing.T) {
	connID := uuid.Must(uuid.NewV4())
	mctx := &kubernetes.MachineCtx{MesheryCtrlsHelper: &models.MesheryControllersHelper{}}

	tests := []struct {
		name     string
		tracker  *machines.ConnectionToStateMachineInstanceTracker
		lookupID uuid.UUID
		wantOK   bool
		wantCtx  *kubernetes.MachineCtx
	}{
		{
			name:     "nil Context is not-ready",
			tracker:  trackerWith(connID, &machines.StateMachine{ID: connID, Context: nil}),
			lookupID: connID,
			wantOK:   false,
		},
		{
			name:     "untracked connection is not-ready",
			tracker:  trackerWith(uuid.Must(uuid.NewV4()), &machines.StateMachine{Context: nil}),
			lookupID: uuid.Must(uuid.NewV4()),
			wantOK:   false,
		},
		{
			name:     "nil tracked instance is not-ready",
			tracker:  trackerWith(connID, nil),
			lookupID: connID,
			wantOK:   false,
		},
		{
			// A boxed typed-nil *MachineCtx is non-nil as an interface, so a bare
			// `Context == nil` check would let it through and the cast would hand
			// back a nil pointer with a nil error - dereferenced one line later.
			name:     "typed-nil Context is not-ready",
			tracker:  trackerWith(connID, &machines.StateMachine{ID: connID, Context: (*kubernetes.MachineCtx)(nil)}),
			lookupID: connID,
			wantOK:   false,
		},
		{
			// A non-kubernetes Context is a genuine type error, distinct from
			// "not ready": it takes the cast's err != nil branch and is logged.
			name:     "wrong Context type is not-ready",
			tracker:  trackerWith(connID, &machines.StateMachine{ID: connID, Context: &struct{ notAMachineCtx bool }{}}),
			lookupID: connID,
			wantOK:   false,
		},
		{
			// A kubernetes machine that initialized but has no controllers helper
			// cannot answer a status query either.
			name:     "missing controllers helper is not-ready",
			tracker:  trackerWith(connID, &machines.StateMachine{ID: connID, Context: &kubernetes.MachineCtx{}}),
			lookupID: connID,
			wantOK:   false,
		},
		{
			name:     "valid Context is ready",
			tracker:  trackerWith(connID, &machines.StateMachine{ID: connID, Context: mctx}),
			lookupID: connID,
			wantOK:   true,
			wantCtx:  mctx,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				config:                                  &models.HandlerConfig{},
				log:                                     newTestLogger(t),
				ConnectionToStateMachineInstanceTracker: tt.tracker,
			}

			ctx, ok := h.machineCtxForConnection(tt.lookupID.String())
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if tt.wantCtx == nil && ctx != nil {
				t.Fatalf("expected nil machine context, got %#v", ctx)
			}
			if tt.wantCtx != nil && ctx != tt.wantCtx {
				t.Fatalf("expected the tracked machine context to be returned, got %#v", ctx)
			}
		})
	}
}

func newControllersStatusTestHandler(t *testing.T) *Handler {
	t.Helper()
	// A nil ConnectionToStateMachineInstanceTracker is fine: unresolved
	// connections degrade to an empty snapshot, which is exactly the path we
	// exercise without standing up real FSM instances and cluster clients.
	return &Handler{
		config: &models.HandlerConfig{},
		log:    newTestLogger(t),
	}
}

func TestSubscribeMesheryControllersStatusHandler_EmitsInitialSnapshot(t *testing.T) {
	h := newControllersStatusTestHandler(t)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequest(http.MethodGet,
		"/api/system/controllers/status/subscribe?connectionIds="+uuid.Must(uuid.NewV4()).String(), nil).
		WithContext(ctx)

	writeCh := make(chan struct{}, 8)
	rec := newSSERecorder(writeCh)

	done := make(chan struct{})
	go func() {
		h.SubscribeMesheryControllersStatusHandler(rec, req, nil, user, nil)
		close(done)
	}()

	// The initial snapshot is written synchronously before the poll loop, so it
	// should land promptly. Unknown connections yield an empty array.
	deadline := time.After(testTimeout)
	for !strings.Contains(rec.body.String(), "data: [") {
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for initial snapshot; got: %q", rec.body.String())
		case <-writeCh:
		case <-time.After(20 * time.Millisecond):
		}
	}

	body := rec.body.String()
	if !strings.Contains(body, "data: []") {
		t.Fatalf("expected empty snapshot for unresolved connection, got: %q", body)
	}
	if !strings.Contains(body, "\n\n") {
		t.Fatalf("snapshot not framed as SSE data message: %q", body)
	}
	if got := rec.header.Get("Content-Type"); got != "text/event-stream" {
		t.Fatalf("Content-Type = %q, want text/event-stream", got)
	}

	cancel()
	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Fatal("handler did not return after context cancellation")
	}
}

func TestSubscribeMesheryControllersStatusHandler_StopsOnContextCancellation(t *testing.T) {
	h := newControllersStatusTestHandler(t)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/api/system/controllers/status/subscribe", nil).WithContext(ctx)
	rec := newSSERecorder(nil)

	done := make(chan struct{})
	go func() {
		h.SubscribeMesheryControllersStatusHandler(rec, req, nil, user, nil)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Fatal("handler did not return after context cancellation")
	}
}

func TestOperatorStatusHandler_UnknownForUnresolvedConnection(t *testing.T) {
	h := newControllersStatusTestHandler(t)

	req := httptest.NewRequest(http.MethodGet,
		"/api/system/controllers/operator/status?connectionId="+uuid.Must(uuid.NewV4()).String(), nil)
	rec := httptest.NewRecorder()

	h.OperatorStatusHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `"controller":"OPERATOR"`) || !strings.Contains(body, `"status":"`+string(controllersStatusUnknown)+`"`) {
		t.Fatalf("expected unknown operator status payload, got: %q", body)
	}
	if !strings.Contains(body, `"connectionId":`) {
		t.Fatalf("expected canonical connectionId key, got: %q", body)
	}
}

// A malformed connectionId must be rejected with 400 rather than falling back
// to uuid.Nil and echoing a zero UUID. Covers all three one-shot status
// handlers, which share the same validation.
func TestControllerStatusHandlers_RejectInvalidConnectionID(t *testing.T) {
	h := newControllersStatusTestHandler(t)

	handlers := map[string]func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider){
		"operator": h.OperatorStatusHandler,
		"meshsync": h.MeshsyncStatusHandler,
		"broker":   h.BrokerStatusHandler,
	}

	for name, handler := range handlers {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet,
				"/api/system/controllers/"+name+"/status?connectionId=not-a-uuid", nil)
			rec := httptest.NewRecorder()

			handler(rec, req, nil, nil, nil)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400 for a malformed connectionId", rec.Code)
			}
			if strings.Contains(rec.Body.String(), uuid.Nil.String()) {
				t.Fatalf("response must not echo the zero UUID, got: %q", rec.Body.String())
			}
		})
	}
}
