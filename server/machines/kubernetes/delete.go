package kubernetes

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DeleteAction struct{}

func (da *DeleteAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (da *DeleteAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	logLevel := viper.GetInt("LOG_LEVEL")
	if viper.GetBool("DEBUG") {
		logLevel = int(logrus.DebugLevel)
	}
	// Initialize Logger instance
	log, err := logger.New("meshery", logger.Options{
		Format:   logger.SyslogLogFormat,
		LogLevel: logLevel,
	})
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	user, ok := ctx.Value(models.UserCtxKey).(*models.User)
	if !ok || user == nil {
		err := fmt.Errorf("user missing from context")
		return machines.NoOp, events.NewEvent().WithCategory("connection").WithAction("update").WithDescription(err.Error()).WithSeverity(events.Error).WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}
	sysID, ok := ctx.Value(models.SystemIDKey).(*core.Uuid)
	if !ok || sysID == nil {
		err := fmt.Errorf("system ID missing from context")
		return machines.NoOp, events.NewEvent().WithCategory("connection").WithAction("update").WithDescription(err.Error()).WithSeverity(events.Error).WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}
	provider, ok := ctx.Value(models.ProviderCtxKey).(models.Provider)
	if !ok || machines.IsProviderNil(provider) {
		err := fmt.Errorf("provider missing from context")
		return machines.NoOp, events.NewEvent().WithCategory("connection").WithAction("update").WithDescription(err.Error()).WithSeverity(events.Error).WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}
	userUUID := user.ID

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromOwner(userUUID).WithDescription("Failed to interact with the connection.")

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	contextID := machinectx.K8sContext.ID

	go func() {

		machinectx.MesheryCtrlsHelper.UpdateOperatorsStatusMap(machinectx.OperatorTracker).
			UndeployDeployedOperators(machinectx.OperatorTracker).
			RemoveCtxControllerHandler(ctx, contextID)

		machinectx.MesheryCtrlsHelper.RemoveMeshSyncDataHandler(ctx, contextID)
	}()

	_ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	context.AfterFunc(_ctx, func() {
		// machinectx.MesheryCtrlsHelper.UpdateOperatorsStatusMap(machinectx.OperatorTracker)
	})

	go models.FlushMeshSyncData(ctx, machinectx.K8sContext, provider, machinectx.EventBroadcaster, user.ID.String(), sysID, log)

	return machines.NoOp, nil, nil
}

func (da *DeleteAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
