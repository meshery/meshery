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
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

type Options struct {
	DBHandler       *database.Handler
	Logger          logger.Handler
	KubeClient      *mesherykube.Client
	MeshSyncChannel chan struct{}
	URL             string
}

// New returns a graphql handler instance
func New(opts Options) http.Handler {
	resolver := &resolver.Resolver{
		Log:             opts.Logger,
		DBHandler:       opts.DBHandler,
		KubeClient:      opts.KubeClient,
		MeshSyncChannel: opts.MeshSyncChannel,
	}

	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
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
