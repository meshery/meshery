package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/layer5io/meshery/graphql/graph"
	"github.com/layer5io/meshery/graphql/graph/generated"
)

// GraphQLHandler returns handler to be used with graphql
func GraphQLHandler() http.Handler {

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// http.Handle("/query", srv)
	return srv
}
