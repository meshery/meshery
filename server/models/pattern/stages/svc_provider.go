package stages

import (
	"github.com/meshery/meshery/server/models/pattern/core"
	"github.com/meshery/meshery/server/models/pattern/patterns"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	schemacore "github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

type ServiceInfoProvider interface {
	GetMesheryPatternResource(
		name string,
		namespace string,
		typ string,
		oamType string,
	) (ID *schemacore.Uuid, err error)
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
