package kubernetes

import (
	"context"
	"errors"
	"fmt"

	"github.com/meshery/meshery/server/machines"
)

// FlushAndResyncResources flushes cached resources and requests a MeshSync
// resync while serialized with MeshSync handler lifecycle changes.
func FlushAndResyncResources(
	ctx context.Context,
	machineCtx *MachineCtx,
	flush func() error,
) (flushErr error, resyncErr error) {
	if machineCtx == nil {
		return nil, ErrResyncK8SResources(fmt.Errorf("machine context is nil"))
	}
	mesheryCtrlsHelper := machineCtx.MesheryCtrlsHelper
	if mesheryCtrlsHelper == nil {
		return nil, ErrResyncK8SResources(fmt.Errorf("machine context does not contain reference to MesheryCtrlsHelper"))
	}
	flushErr, resyncErr = mesheryCtrlsHelper.FlushAndResyncMeshsync(flush)
	if resyncErr != nil {
		resyncErr = ErrResyncK8SResources(errors.Join(fmt.Errorf("error calling FlushAndResyncMeshsync"), resyncErr))
	}
	return flushErr, resyncErr
}

func ResyncResources(ctx context.Context, sm *machines.StateMachine) error {
	mashineCtx, err := GetMachineCtx(sm.GetContext(), nil)
	if err != nil {
		return ErrResyncK8SResources(
			fmt.Errorf(
				"failed to get machine context for machine %v: %w",
				sm.ID,
				err,
			),
		)
	}
	mesheryCtrlsHelper := mashineCtx.MesheryCtrlsHelper
	if mesheryCtrlsHelper == nil {
		return ErrResyncK8SResources(
			fmt.Errorf(
				"machine context does not contain reference to MesheryCtrlsHelper for machine %v",
				sm.ID,
			),
		)
	}

	if err := mesheryCtrlsHelper.ResyncMeshsync(ctx); err != nil {
		return ErrResyncK8SResources(
			errors.Join(
				fmt.Errorf("error calling ResyncMeshsync for machine %v", sm.ID),
				err,
			),
		)
	}

	return nil
}
