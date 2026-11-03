package kubernetes

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
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
		return machines.NoOp, nil, models.ErrInitLogger(err)
	}
	userUUID, sysID := machines.ExtractActorIDs(ctx)
	provider, _ := machines.ExtractProvider(ctx)

	eventBuilder := events.NewEvent().WithCategory("connection").WithAction("update").WithDescription("Failed to interact with the connection.")
	if userUUID != core.Uuid(uuid.Nil) {
		eventBuilder.ActedUpon(userUUID)
	}
	machines.ApplyEventAttribution(eventBuilder, userUUID, sysID)

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	contextID := machinectx.K8sContext.ID

	if machinectx.MesheryCtrlsHelper != nil {
		go func() {
			machinectx.MesheryCtrlsHelper.UpdateOperatorsStatusMap(machinectx.OperatorTracker).
				UndeployDeployedOperators(machinectx.OperatorTracker, contextID).
				RemoveCtxControllerHandler(ctx, contextID)

			machinectx.MesheryCtrlsHelper.RemoveMeshSyncDataHandler(ctx, contextID)
		}()
	}

	_ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	context.AfterFunc(_ctx, func() {
		// machinectx.MesheryCtrlsHelper.UpdateOperatorsStatusMap(machinectx.OperatorTracker)
	})

	var userIDStr string
	if userUUID != core.Uuid(uuid.Nil) {
		userIDStr = userUUID.String()
	}
	if provider != nil {
		go models.FlushMeshSyncData(ctx, machinectx.K8sContext, provider, machinectx.EventBroadcaster, userIDStr, sysID, log)
	}

	return machines.NoOp, nil, nil
}

func (da *DeleteAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
