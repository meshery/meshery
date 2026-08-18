package kubernetes

import (
	"context"
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
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*core.Uuid)
	provider, _ := ctx.Value(models.ProviderCtxKey).(models.Provider)
	userUUID := user.ID

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromOwner(userUUID).WithDescription("Failed to interact with the connection.")

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	contextID := machinectx.K8sContext.ID

	machinectx.ActionMutex.Lock()
	go func() {
		defer machinectx.ActionMutex.Unlock()

		if ctx.Err() != nil {
			machinectx.log.Info("Delete side-effects aborted due to lifecycle cancellation")
			return
		}

		machinectx.MesheryCtrlsHelper.UpdateOperatorsStatusMap(machinectx.OperatorTracker).
			UndeployDeployedOperators(machinectx.OperatorTracker, contextID).
			RemoveCtxControllerHandler(ctx, contextID)

		machinectx.MesheryCtrlsHelper.RemoveMeshSyncDataHandler(ctx, contextID)

		models.FlushMeshSyncData(ctx, machinectx.K8sContext, provider, machinectx.EventBroadcaster, user.ID.String(), sysID, log)
	}()

	_ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	context.AfterFunc(_ctx, func() {
		// machinectx.MesheryCtrlsHelper.UpdateOperatorsStatusMap(machinectx.OperatorTracker)
	})

	return machines.NoOp, nil, nil
}

func (da *DeleteAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
