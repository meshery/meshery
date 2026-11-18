package kubernetes

import (
	"context"
	"testing"
)

func TestFlushAndResyncResourcesDoesNotFlushWithoutMesheryHelper(t *testing.T) {
	machineCtx := &MachineCtx{}

	flushCalled := false
	flushErr, resyncErr := FlushAndResyncResources(
		context.Background(),
		machineCtx,
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
