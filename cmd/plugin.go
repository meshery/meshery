package main

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/layer5io/meshery/graphql/generated"
	"github.com/layer5io/meshery/graphql/meshes/istio"
	"github.com/layer5io/meshery/graphql/resolver"
	"github.com/layer5io/meshkit/database"
)

// Run will run the server for now
func Run(db *database.Handler) (*handler.Server, error) {
	fmt.Println("plugin running")
	resolver := &resolver.Resolver{
		Istio: &istio.IstioHandler{DB: db},
	}

	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver})), nil
}
