package kubernetes

import (
	"context"
	"errors"
	"fmt"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshkit/models/events"
)

func ResyncResources(ctx context.Context, sm *machines.StateMachine) error {
	// TODO what to put as event builder?
	mashineCtx, err := GetMachineCtx(sm.Context, events.NewEvent())
	if err != nil {
		return ErrResyncK8SResources(
			fmt.Errorf(
				"can not receive machine context for machine %v",
				sm.ID,
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
