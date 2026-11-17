package kubernetes

import (
	"context"
	"testing"

	"github.com/meshery/meshery/server/machines"
)

func TestFlushAndResyncResourcesDoesNotFlushWithoutMesheryHelper(t *testing.T) {
	sm := &machines.StateMachine{
		Name:    "kubernetes",
		Context: &MachineCtx{},
	}

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
		t.Fatal("flush callback must not run when MesheryCtrlsHelper is nil")
	}
	if flushErr != nil {
		t.Fatalf("flush error = %v, want nil", flushErr)
	}
	if resyncErr == nil {
		t.Fatal("expected an error when MesheryCtrlsHelper is nil")
	}
}
