package graphql

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/internal/graphql/generated"
	"github.com/layer5io/meshery/server/internal/graphql/resolver"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils/broadcast"
)

type Options struct {
	Logger      logger.Handler
	BrokerConn  broker.Handler
	Config      *models.HandlerConfig
	URL         string
	Broadcaster broadcast.Broadcaster
}

// New returns a graphql handler instance
func New(opts Options) http.Handler {
	res := &resolver.Resolver{
		Log:                          opts.Logger,
		MeshSyncChannelPerK8sContext: make(map[string]chan struct{}),
		BrokerConn:                   opts.BrokerConn,
		Config:                       opts.Config,

		Broadcast: opts.Broadcaster,
	}

	config := generated.Config{
		Resolvers: res,
	}
	config.Directives.KubernetesMiddleware = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {

		handler, _ := ctx.Value(models.HandlerKey).(*handlers.Handler)
		user, _ := ctx.Value(models.UserCtxKey).(*models.User)
		provider, _ := ctx.Value(models.ProviderCtxKey).(models.Provider)
		ctx, err = handlers.KubernetesMiddleware(ctx, handler, provider, user, []string{})
		if err != nil {
			return nil, err
		}

		ctx, err = handlers.MesheryControllersMiddleware(ctx, handler)
		if err != nil {
			return nil, err
		}
		return next(ctx)
	}

	srv := handler.New(generated.NewExecutableSchema(config))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow any origin to establish websocket connection
				return true
			},
		},
	})

	return srv
}

// NewPlayground returns a graphql playground instance
func NewPlayground(opts Options) http.Handler {
	return playground.Handler("GraphQL playground", opts.URL)
}
