package kubernetes

import (
	"testing"

	"github.com/meshery/meshery/server/machines"
)

// TestRegisteredStateRoutesNotFound guards the fix for #20642: RegisterAction
// redirects the machine to NOTFOUND when the cluster ping fails, so the
// REGISTERED state must declare that edge. Without it the redirect is dropped
// and an unreachable cluster is wrongly persisted as REGISTERED.
func TestRegisteredStateRoutesNotFound(t *testing.T) {
	st := Registered()

	if got, ok := st.Events[machines.NotFound]; !ok || got != machines.NOTFOUND {
		t.Fatalf("Registered must route NotFound -> NOTFOUND; got %q present=%v", got, ok)
	}

	// Existing edges must remain intact.
	if st.Events[machines.Connect] != machines.CONNECTED {
		t.Errorf("Registered lost its Connect -> CONNECTED edge")
	}
	if st.Events[machines.Ignore] != machines.IGNORED {
		t.Errorf("Registered lost its Ignore -> IGNORED edge")
	}
}
