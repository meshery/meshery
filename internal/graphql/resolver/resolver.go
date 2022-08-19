package resolver

import (
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils/broadcast"
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
	performanceChannel  chan *model.PerfPageResult
	brokerChannel       chan *broker.Message
	addonChannel        chan []*model.AddonList
	controlPlaneChannel chan []*model.ControlPlane
	dataPlaneChannel    chan []*model.DataPlane
}
