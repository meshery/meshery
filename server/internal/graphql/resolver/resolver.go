package resolver

import (
	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/utils/broadcast"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Log                          logger.Handler
	BrokerConn                   broker.Handler
	MeshSyncChannelPerK8sContext map[string]chan struct{}
	Config                       *models.HandlerConfig
	Broadcast                    broadcast.Broadcaster
	// operatorChannel         chan *model.OperatorStatus
	// performanceChannel  chan *model.PerfPageResult
	brokerChannel       chan *broker.Message
	addonChannel        chan []*model.AddonList        //nolint:unused // used by listenToAddonState
	controlPlaneChannel chan []*model.ControlPlane     //nolint:unused // used by listenToControlPlaneState
	dataPlaneChannel    chan []*model.DataPlane        //nolint:unused // used by listenToDataPlaneState
}
