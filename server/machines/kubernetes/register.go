package kubernetes

import (
	"context"
	"fmt"

	"github.com/meshery/schemas/models/core"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	meshmodelcore "github.com/meshery/meshery/server/models/meshmodel/core"
	"github.com/meshery/meshkit/models/events"
)

type RegisterAction struct{}

// ExecuteOnEntry is a no-op; RegisterAction performs all of its work in Execute.
//
// TODO: Execute On Entry and Exit should not return next eventtype i suppose, look again.
func (ra *RegisterAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (ra *RegisterAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	userUUID, sysID := machines.ExtractActorIDs(ctx)
	provider, ok := machines.ExtractProvider(ctx)
	eventBuilder := events.NewEvent().ActedUpon(uuid.Nil).WithCategory("connection").WithAction("register").WithDescription("Failed to interact with the connection.")
	machines.ApplyEventAttribution(eventBuilder, userUUID, sysID)

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	if !ok {
		errProvider := machines.ErrNilProvider(fmt.Errorf("provider is missing or invalid"))
		eventBuilder.WithMetadata(map[string]interface{}{"error": errProvider})
		return machines.NoOp, eventBuilder.Build(), errProvider
	}


	machinectx.log.Debug("executing ping test for connection", machinectx.K8sContext.ConnectionID)
	err = machinectx.K8sContext.PingTest()

	if err != nil {
		// Mark the ping/reachability failure as an Error so it is findable under
		// the notification center's Error filter, matching the sibling
		// ConnectAction and DiscoverAction failure events. Without this the event
		// lands with an empty severity and is effectively invisible to users.
		eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to ping kubernetes context %s at %s", machinectx.K8sContext.Name, machinectx.K8sContext.Server)).WithMetadata(map[string]interface{}{"error": err})
		machinectx.log.Error(err)
		return machines.NotFound, eventBuilder.Build(), err
	}

	if machinectx.K8sCompRegHelper != nil {
		context := []*models.K8sContext{&machinectx.K8sContext}
		var userIDStr string
		if userUUID != core.Uuid(uuid.Nil) {
			userIDStr = userUUID.String()
		}

		machinectx.K8sCompRegHelper.UpdateContexts(context).RegisterComponents(context, []models.K8sRegistrationFunction{meshmodelcore.RegisterK8sMeshModelComponents}, machinectx.RegistryManager, machinectx.EventBroadcaster, provider, userIDStr, true)
	}

	return machines.Connect, nil, nil
}

func (ra *RegisterAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
