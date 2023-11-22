package kubernetes

import (
	"context"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshkit/models/events"
)

type NotFoundAction struct{}

func (ia *NotFoundAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {

	return machines.NoOp, nil, nil
}
func (ia *NotFoundAction) Execute(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	// Just pass along, the status is update as we exit from the event.
	return machines.NoOp, nil, nil
}

func (ia *NotFoundAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
