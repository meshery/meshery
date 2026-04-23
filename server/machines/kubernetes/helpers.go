package kubernetes

import (
	"context"
	"fmt"

	"github.com/meshery/schemas/models/core"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/helpers"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/utils/kubernetes"
	"github.com/spf13/viper"
)

// extractActionContext pulls the user, system ID, and provider from the
// action context. Returns a descriptive error if any of the three is
// missing or of the wrong type, so lifecycle actions invoked without
// upstream HTTP middleware populating the context (scheduler events,
// internal invocations, timed-out contexts) can abort cleanly rather
// than panic later on user.ID / *sysID / provider.X().
func extractActionContext(ctx context.Context) (*models.User, *core.Uuid, models.Provider, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	if user == nil {
		return nil, nil, nil, fmt.Errorf("kubernetes state machine: missing user in context")
	}
	sysID, _ := ctx.Value(models.SystemIDKey).(*core.Uuid)
	if sysID == nil {
		return nil, nil, nil, fmt.Errorf("kubernetes state machine: missing system ID in context")
	}
	// Provider may legitimately be a typed-nil interface (tests and
	// internal invocations) and callers only use it after their own
	// error paths; surface it as-is instead of rejecting here.
	provider, _ := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return user, sysID, provider, nil
}

var (
	adapterURLs    = viper.GetStringSlice("ADAPTER_URLS")
	adapterTracker = helpers.NewAdaptersTracker(adapterURLs)
)

func GenerateClientSetAction(k8sContext *models.K8sContext, eventBuilder *events.EventBuilder, log logger.Handler) (*kubernetes.Client, error) {
	eventBuilder.ActedUpon(uuid.FromStringOrNil(k8sContext.ConnectionID))

	eventMetadata := map[string]interface{}{}

	handler, err := models.GenerateK8sClientSet(k8sContext, eventBuilder, eventMetadata, log)
	if handler == nil {
		return nil, err
	}

	return handler, nil
}

func GetMachineCtx(machinectx interface{}, eb *events.EventBuilder) (*MachineCtx, error) {
	machineCtx, ok := machinectx.(*MachineCtx)
	if !ok {
		err := machines.ErrAssertMachineCtx(fmt.Errorf("asserting of context %v failed", machinectx))
		if eb != nil {
			eb.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": err,
			})
		}
		return nil, err
	}
	return machineCtx, nil
}

func AssignClientSetToContext(machinectx *MachineCtx, eventBuilder *events.EventBuilder) error {
	if machinectx.clientset != nil {
		return nil
	}

	k8sContext := machinectx.K8sContext
	eventBuilder.ActedUpon(uuid.FromStringOrNil(k8sContext.ConnectionID))

	handler, err := GenerateClientSetAction(&k8sContext, eventBuilder, machinectx.log)
	if err != nil {
		// perofmr event publishinh and err handling
		return err
	}
	machinectx.clientset = handler
	return nil
}

func AssignControllerHandlers(machinectx *MachineCtx, systemID *core.Uuid, provider models.Provider) {
	machinectx.MesheryCtrlsHelper = models.NewMesheryControllersHelper(
		machinectx.log,
		models.NewOperatorDeploymentConfig(adapterTracker),
		models.GetDBInstance(),
		machinectx.EventBroadcaster,
		provider,
		systemID,
	)
}
