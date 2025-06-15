package kubernetes

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
)

type ConnectAction struct{}

// Execute On Entry and Exit should not return next eventtype i suppose, look again.
func (ca *ConnectAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (ca *ConnectAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	userUUID := uuid.FromStringOrNil(user.ID)
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.").WithSeverity(events.Error)

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}
	token, ok := ctx.Value(models.TokenCtxKey).(string)
	if !ok {
		errToken := ErrConnectActionFromString("failed to retrieve user token")
		eventBuilder.WithMetadata(map[string]interface{}{"error": errToken})
		return machines.NoOp, eventBuilder.Build(), errToken

	}
	connectionID := uuid.FromStringOrNil(machinectx.K8sContext.ID)
	connection, _, err := provider.GetConnectionByIDAndKind(
		token,
		connectionID,
		"kubernetes",
	)
	if err != nil {
		errConnection := ErrConnectAction(err)
		eventBuilder.WithMetadata(map[string]interface{}{"error": errConnection})
		return machines.NoOp, eventBuilder.Build(), errConnection

	}

	go func() {
		ctrlHelper := machinectx.MesheryCtrlsHelper.
			AddCtxControllerHandlers(machinectx.K8sContext).
			SetMeshsyncDeploymentMode(models.MeshsyncDeploymentModeFromString(connection.MeshsyncDeploymentMode)).
			UpdateOperatorsStatusMap(machinectx.OperatorTracker).
			DeployUndeployedOperators(machinectx.OperatorTracker)
		ctrlHelper.AddMeshsynDataHandlers(ctx, machinectx.K8sContext, userUUID, *sysID, provider)
	}()
	return machines.NoOp, nil, nil
}

func (ca *ConnectAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
