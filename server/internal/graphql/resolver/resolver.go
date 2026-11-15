package resolver

import (
	"context"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/utils/broadcast"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// getProvider safely extracts the Provider from the request context using the comma-ok pattern.
// If the provider key is absent, nil, or is not of type models.Provider,
// it returns ErrInvalidRequest. Middlewares and internal callers must never
// store a typed nil in context: a typed nil passes the interface nil check,
// so it is the caller's responsibility to store only a live provider.
func getProvider(ctx context.Context) (models.Provider, error) {
	if ctx == nil {
		return nil, ErrInvalidRequest
	}
	provider, ok := ctx.Value(models.ProviderCtxKey).(models.Provider)
	if !ok || provider == nil {
		return nil, ErrInvalidRequest
	}
	return provider, nil
}

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
