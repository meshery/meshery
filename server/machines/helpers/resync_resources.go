package helpers

import (
	"context"
	"fmt"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
)

func ResyncResources(ctx context.Context, sm *machines.StateMachine) error {
	if sm.Name != "kubernetes" {
		return ErrResyncResources(
			fmt.Errorf(
				"ResyncResources operation is not implemented for machine type %s",
				sm.Name,
			),
		)
	}
	return kubernetes.ResyncResources(ctx, sm)
}
