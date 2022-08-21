package graphql

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/layer5io/meshery/internal/graphql/generated"
	"github.com/layer5io/meshery/internal/graphql/resolver"
	"github.com/layer5io/meshery/models"
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

	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: res,
	}))

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
