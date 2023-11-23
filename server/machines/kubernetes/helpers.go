package kubernetes

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/utils/kubernetes"
)

func GenerateClientSetAction(k8sContext *models.K8sContext, eventBuilder *events.EventBuilder) (*kubernetes.Client, error) {
	eventBuilder.ActedUpon(uuid.FromStringOrNil(k8sContext.ConnectionID))

	eventMetadata := map[string]interface{}{}

	handler, err := models.GenerateK8sClientSet(k8sContext, eventBuilder, eventMetadata)
	if handler == nil {
		return nil, err
	}

	return handler, nil
}

func GetMachineCtx(machinectx interface{}, eb *events.EventBuilder) (*MachineCtx, error) {
	machineCtx, ok := machinectx.(*MachineCtx)
	if !ok {
		err := machines.ErrAssertMachineCtx(fmt.Errorf("asserting of context %v failed", machinectx))
		eb.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": err,
		})
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

	handler, err := GenerateClientSetAction(&k8sContext, eventBuilder)
	if err != nil {
		// perofmr event publishinh and err handling
		return err
	}
	machinectx.clientset = handler
	return nil
}
