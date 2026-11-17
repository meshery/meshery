package helpers

import (
	"context"
	"testing"

	"github.com/meshery/meshery/server/machines"
)

func TestFlushAndResyncResourcesRejectsUnsupportedMachine(t *testing.T) {
	sm := &machines.StateMachine{Name: "prometheus"}

	flushCalled := false
	flushErr, resyncErr := FlushAndResyncResources(
		context.Background(),
		sm,
		func() error {
			flushCalled = true
			return nil
		},
	)

	if flushCalled {
		t.Fatal("flush callback must not run for an unsupported machine type")
	}
	if flushErr != nil {
		t.Fatalf("flush error = %v, want nil", flushErr)
	}
	if resyncErr == nil {
		t.Fatal("expected an error for an unsupported machine type")
	}
}
