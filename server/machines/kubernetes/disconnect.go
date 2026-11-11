package kubernetes

import (
	"context"
	"time"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshkit/models/events"
)

type DisconnectAction struct{}

func (da *DisconnectAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil

}

// Execute cleanly tears down MeshSync data handlers and disconnects from the Kubernetes cluster.
// It leaves the cluster-scoped Meshery Operator installed.
func (da *DisconnectAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	userUUID, systemID := machines.ResolveUserAndSystemID(ctx)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(systemID).FromOwner(userUUID).WithDescription("Failed to interact with the connection.")

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	machinectx.setupLock.Lock()
	machinectx.disconnected = true
	machinectx.retryPending = false
	machinectx.setupLock.Unlock()

	go func() {
		machinectx.setupLock.Lock()
		defer machinectx.setupLock.Unlock()
		if machinectx.MesheryCtrlsHelper != nil {
			// Releases this connection's claim on the cluster's Meshery Operator
			// without uninstalling it: the operator is cluster-scoped and other
			// connections may still depend on it.
			machinectx.MesheryCtrlsHelper.DetachFromCluster(ctx, machinectx.OperatorTracker, machinectx.K8sContext)
		}
	}()

	_ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	context.AfterFunc(_ctx, func() {
		// machinectx.MesheryCtrlsHelper.UpdateOperatorsStatusMap(machinectx.OperatorTracker)
	})

	return machines.NoOp, nil, nil
}

func (da *DisconnectAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
