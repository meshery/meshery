package stages

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/patterns"
	"github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

type ServiceInfoProvider interface {
	GetMesheryPatternResource(
		name string,
		namespace string,
		typ string,
		oamType string,
	) (ID *uuid.UUID, err error)
	IsDelete() bool
}

type ServiceActionProvider interface {
	Terminate(error)
	Log(msg string)
	Provision(CompConfigPair) ([]patterns.DeploymentMessagePerContext, error)
	GetRegistry() *registry.RegistryManager
	DryRun([]*component.ComponentDefinition) (map[string]map[string]core.DryRunResponseWrapper, error)
	Mutate(*pattern.PatternFile) //Uses pre-defined policies/configuration to mutate the pattern
}
