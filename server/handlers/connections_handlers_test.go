package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/schemas/models/core"
)

func newRegistrationTracker() *machines.ConnectionToStateMachineInstanceTracker {
	return &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}
}

func cancelRegistrationRequest(registrationID string) *http.Request {
	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/integrations/connections/register/"+registrationID,
		nil,
	)
	// The handler reads the path parameter via mux.Vars; SetURLVars stands in
	// for the router's extraction in unit tests.
	return mux.SetURLVars(req, map[string]string{"registrationId": registrationID})
}

func TestCancelConnectionRegistration_RemovesTrackedInstance(t *testing.T) {
	tracker := newRegistrationTracker()
	registrationID := uuid.Must(uuid.NewV4())
	tracker.Add(registrationID, &machines.StateMachine{})

	h := &Handler{ConnectionToStateMachineInstanceTracker: tracker}
	rec := httptest.NewRecorder()

	h.CancelConnectionRegistration(rec, cancelRegistrationRequest(registrationID.String()), nil, nil, nil)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d", http.StatusNoContent, rec.Code)
	}
	if _, ok := tracker.Get(registrationID); ok {
		t.Fatal("expected the tracked state machine to be removed")
	}
}

func TestCancelConnectionRegistration_UnknownIdIsIdempotent(t *testing.T) {
	tracker := newRegistrationTracker()
	h := &Handler{ConnectionToStateMachineInstanceTracker: tracker}
	rec := httptest.NewRecorder()

	h.CancelConnectionRegistration(rec, cancelRegistrationRequest(uuid.Must(uuid.NewV4()).String()), nil, nil, nil)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestCancelConnectionRegistration_MalformedIdIsRejected(t *testing.T) {
	tracker := newRegistrationTracker()
	// A malformed id must not be coerced to uuid.Nil: anything tracked under
	// the nil UUID would be silently discarded. Seed a nil-keyed instance and
	// assert it survives the rejected request.
	tracker.Add(uuid.Nil, &machines.StateMachine{})

	h := &Handler{
		log:                                     newTestLogger(t),
		ConnectionToStateMachineInstanceTracker: tracker,
	}
	rec := httptest.NewRecorder()

	h.CancelConnectionRegistration(rec, cancelRegistrationRequest("not-a-uuid"), nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if _, ok := tracker.Get(uuid.Nil); !ok {
		t.Fatal("expected the nil-keyed state machine to remain tracked")
	}
}
