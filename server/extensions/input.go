package extensions

import (
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
)

type ExtensionInput struct {
	DBHandler            *database.Handler
	MeshSyncChannel      chan struct{}
	Logger               logger.Handler
	BrokerConn           broker.Handler
	K8sConnectionTracker *machines.ConnectionToStateMachineInstanceTracker
}
