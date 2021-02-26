package resolver

import (
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/database"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshsync/pkg/broker"
	"k8s.io/client-go/dynamic"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DBHandler        *database.Handler
	GetKubeClient    func() (*mesherykube.Client, error)
	GetDynamicClient func() (dynamic.Interface, error)

	operatorChannel chan *model.OperatorStatus
	meshsyncChannel chan *broker.Message
	addonChannel    chan []*model.AddonList
}
