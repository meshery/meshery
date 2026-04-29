package kubernetes

import (
	"context"
	"errors"
	"fmt"

	"github.com/meshery/meshery/server/machines"
)

func ResyncResources(ctx context.Context, sm *machines.StateMachine) error {
	mashineCtx, err := GetMachineCtx(sm.Context, nil)
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
