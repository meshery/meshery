package helpers

import (
	"context"
	"fmt"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
)

// FlushAndResyncResources flushes cached resources and requests a resync
// for the specified machine type.
func FlushAndResyncResources(
	ctx context.Context,
	sm *machines.StateMachine,
	flush func() error,
) (flushErr error, resyncErr error) {
	if sm.Name != "kubernetes" {
		return nil, ErrResyncResources(
			fmt.Errorf(
				"FlushAndResyncResources operation is not implemented for machine type %s",
				sm.Name,
			),
		)
	}

	return kubernetes.FlushAndResyncResources(ctx, sm, flush)
}

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
