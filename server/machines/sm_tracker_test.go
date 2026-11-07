package machines

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"
)

// RemoveIfMatch is what lets a late cleanup for a retired machine run safely
// after a connection with the same deterministic id was re-imported.
func TestRemoveIfMatchKeepsReplacement(t *testing.T) {
	id := uuid.Must(uuid.NewV4())
	old, fresh := &StateMachine{ID: id}, &StateMachine{ID: id}
	tracker := &ConnectionToStateMachineInstanceTracker{ConnectToInstanceMap: map[core.Uuid]*StateMachine{}}

	tracker.Add(id, old)
	tracker.Add(id, fresh)
	if tracker.RemoveIfMatch(id, old) {
		t.Fatal("RemoveIfMatch removed an entry holding a different instance")
	}
	if got, ok := tracker.Get(id); !ok || got != fresh {
		t.Fatal("replacement machine was dropped by the old machine's cleanup")
	}

	if !tracker.RemoveIfMatch(id, fresh) {
		t.Fatal("RemoveIfMatch did not remove the matching instance")
	}
	if _, ok := tracker.Get(id); ok {
		t.Fatal("entry still present after a matching RemoveIfMatch")
	}
	if tracker.RemoveIfMatch(id, fresh) {
		t.Fatal("RemoveIfMatch reported a removal for an absent id")
	}
}
