package resolver

import (
	"github.com/layer5io/meshkit/database"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DBHandler  *database.Handler
	KubeClient *mesherykube.Client
}
