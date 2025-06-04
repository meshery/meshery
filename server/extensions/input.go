package extensions

import (
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
)

type ExtensionInput struct {
	DBHandler            *database.Handler
	MeshSyncChannel      chan struct{}
	Logger               logger.Handler
	BrokerConn           broker.Handler
	K8sConnectionTracker *machines.ConnectionToStateMachineInstanceTracker
}
