package resolver

import (
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/database"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshsync/pkg/broker"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DBHandler  *database.Handler
	KubeClient *mesherykube.Client

	operatorChannel chan *model.OperatorStatus
	meshsyncChannel chan *broker.Message
}
