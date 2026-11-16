package handlers

import (
	"testing"
	"time"

	"github.com/meshery/meshery/server/meshes"
)

// TestAdapterOperationOwners verifies that an operation belongs only to the
// user who started it, and only until it expires or is forgotten.
func TestAdapterOperationOwners(t *testing.T) {
	owners := newAdapterOperationOwners()
	now := time.Now()
	owners.now = func() time.Time { return now }

	owners.record("op-1", "alice")

	if !owners.ownedBy("op-1", "alice") {
		t.Error("expected op-1 to belong to alice")
	}
	if owners.ownedBy("op-1", "bob") {
		t.Error("expected op-1 not to belong to bob")
	}
	if owners.ownedBy("op-unknown", "alice") {
		t.Error("expected an unknown operation to belong to nobody")
	}

	owners.forget("op-1")
	if owners.ownedBy("op-1", "alice") {
		t.Error("expected a forgotten operation to belong to nobody")
	}

	owners.record("op-2", "alice")
	now = now.Add(adapterOperationOwnerTTL + time.Minute)
	if owners.ownedBy("op-2", "alice") {
		t.Error("expected an expired operation to belong to nobody")
	}

	owners.record("op-3", "bob")
	if _, ok := owners.owners["op-2"]; ok {
		t.Error("expected recording an operation to drop expired ones")
	}
}

// TestAdapterOperationOwnersNil verifies that a handler built without the
// tracker drops operation events instead of panicking.
func TestAdapterOperationOwnersNil(t *testing.T) {
	var owners *adapterOperationOwners
	owners.record("op-1", "alice")
	owners.forget("op-1")
	if owners.ownedBy("op-1", "alice") {
		t.Error("expected no operation to belong to anyone without a tracker")
	}
}

// TestForwardAdapterEvent verifies that a user only receives the events of
// their own operations, plus events that aren't tied to an operation.
func TestForwardAdapterEvent(t *testing.T) {
	owners := newAdapterOperationOwners()
	owners.record("alice-op", "alice")

	tests := []struct {
		name        string
		operationID string
		userID      string
		want        bool
	}{
		{name: "own operation", operationID: "alice-op", userID: "alice", want: true},
		{name: "another user's operation", operationID: "alice-op", userID: "bob", want: false},
		{name: "operation started elsewhere", operationID: "unknown-op", userID: "alice", want: false},
		{name: "event without an operation", operationID: "", userID: "bob", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &meshes.EventsResponse{OperationId: tt.operationID, Summary: "summary"}
			if got := forwardAdapterEvent(owners, event, tt.userID); got != tt.want {
				t.Errorf("forwardAdapterEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
