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

// Execute undeploys deployed operators, removes controller handlers, and flushes MeshSync cluster data.
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
	userUUID, systemID := machines.ResolveUserAndSystemID(ctx)
	provider, _ := ctx.Value(models.ProviderCtxKey).(models.Provider)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(systemID).FromOwner(userUUID).WithDescription("Failed to interact with the connection.")

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	// If no active user was found in context (e.g. background deletion), resolve the owner
	// from the persisted context so event attribution and MeshSync cleanup remain accurate.
	if userUUID == core.Uuid(uuid.Nil) {
		if machinectx.K8sContext.Owner != nil && *machinectx.K8sContext.Owner != core.Uuid(uuid.Nil) {
			userUUID = *machinectx.K8sContext.Owner
		} else if machinectx.K8sContext.CreatedBy != nil && *machinectx.K8sContext.CreatedBy != core.Uuid(uuid.Nil) {
			userUUID = *machinectx.K8sContext.CreatedBy
		}
	}

	contextID := machinectx.K8sContext.ID
	if machinectx.K8sCompRegHelper != nil {
		machinectx.K8sCompRegHelper.RemoveContext(contextID)
	}

	go func() {
		if machinectx.MesheryCtrlsHelper != nil {
			machinectx.MesheryCtrlsHelper.UpdateOperatorsStatusMap(machinectx.OperatorTracker).
				UndeployDeployedOperators(machinectx.OperatorTracker, contextID).
				RemoveCtxControllerHandler(ctx, contextID)

			machinectx.MesheryCtrlsHelper.RemoveMeshSyncDataHandler(ctx, contextID)
		}
	}()

	_ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	context.AfterFunc(_ctx, func() {
		// machinectx.MesheryCtrlsHelper.UpdateOperatorsStatusMap(machinectx.OperatorTracker)
	})

	go models.FlushMeshSyncData(ctx, machinectx.K8sContext, provider, machinectx.EventBroadcaster, userUUID.String(), &systemID, log)

	return machines.NoOp, nil, nil
}

func (da *DeleteAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
