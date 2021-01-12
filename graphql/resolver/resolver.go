package resolver

import "github.com/layer5io/meshery/graphql/meshes/istio"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Istio *istio.IstioHandler
}
