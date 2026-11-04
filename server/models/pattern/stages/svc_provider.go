package stages

import (
	patterncore "github.com/meshery/meshery/server/models/pattern/core"
	"github.com/meshery/meshery/server/models/pattern/patterns"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta2/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

type ServiceInfoProvider interface {
	GetMesheryPatternResource(
		name string,
		namespace string,
		typ string,
		oamType string,
	) (ID *core.Uuid, err error)
	IsDelete() bool
}

type ServiceActionProvider interface {
	Terminate(error)
	Log(msg string)
	Provision(CompConfigPair) ([]patterns.DeploymentMessagePerContext, error)
	GetRegistry() *registry.RegistryManager
	DryRun([]*component.ComponentDefinition) (map[string]map[string]patterncore.DryRunResponseWrapper, error)
	Mutate(*pattern.PatternFile) //Uses pre-defined policies/configuration to mutate the pattern
}
