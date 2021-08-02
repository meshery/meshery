package resolver

import (
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/logger"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Log             logger.Handler
	BrokerConn      broker.Handler
	MeshSyncChannel chan struct{}
	Config          *models.HandlerConfig

	operatorSyncChannel     chan struct{}
	meshsyncLivenessChannel chan struct{}
	operatorChannel         chan *model.OperatorStatus
	brokerChannel           chan *broker.Message
	addonChannel            chan []*model.AddonList
	controlPlaneChannel     chan []*model.ControlPlane
}
