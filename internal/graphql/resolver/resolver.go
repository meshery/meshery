package resolver

import (
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Log        logger.Handler
	DBHandler  *database.Handler
	KubeClient *mesherykube.Client

	operatorChannel chan *model.OperatorStatus
	meshsyncChannel chan *broker.Message
	addonChannel    chan []*model.AddonList
}
