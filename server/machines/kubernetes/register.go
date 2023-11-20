package kubernetes

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshery/server/models/meshmodel/core"
	"github.com/layer5io/meshkit/models/events"
)

type RegisterAction struct {}

// Execute On Entry and Exit should not return next eventtype i suppose, look again.
func(ra *RegisterAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func(ra *RegisterAction) Execute(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	userUUID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().ActedUpon(uuid.Nil).WithCategory("connection").WithAction("register").FromSystem(*sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.")

	
	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}
	
	machinectx.log.Debug("executing ping test for connection", machinectx.K8sContext.ConnectionID)
	err = machinectx.K8sContext.PingTest()
	
	if err != nil {
		eventBuilder.WithDescription(fmt.Sprintf("Unable to ping kubernetes context %s at %s", machinectx.K8sContext.Name,machinectx.K8sContext.Server)).WithMetadata(map[string]interface{}{"error": err})
		machinectx.log.Error(err)	
		return machines.NotFound, eventBuilder.Build(), err
	}

	token, _ := ctx.Value(models.TokenCtxKey).(string)
	context := []*models.K8sContext{&machinectx.K8sContext}

	connection, statusCode, err := machinectx.Provider.UpdateConnectionStatusByID(token, uuid.FromStringOrNil(machinectx.K8sContext.ConnectionID), connections.REGISTERED)

	if err != nil {
		return machines.NoOp, eventBuilder.WithDescription(fmt.Sprintf("Failed to register the connection \"%s\" at %s", machinectx.K8sContext.Name, machinectx.K8sContext.Server)).WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}

	machinectx.K8sCompRegHelper.UpdateContexts(context).RegisterComponents(context, []models.K8sRegistrationFunction{core.RegisterK8sMeshModelComponents}, machinectx.RegistryManager, machinectx.EventBroadcaster, machinectx.Provider, user.ID, true)

	machinectx.log.Debug("HTTP status:", statusCode, "updated status for connection", connection.ID)
	machinectx.log.Debug("exiting execute func from registered state", connection)

	return machines.Connect, nil, nil
}

func(ra *RegisterAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
