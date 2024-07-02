package stages

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/patterns"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/layer5io/meshkit/models/meshmodel/registry"
)

type ServiceInfoProvider interface {
	GetMesheryPatternResource(
		name string,
		namespace string,
		typ string,
		oamType string,
	) (ID *uuid.UUID, err error)
	GetServiceMesh() (name string, version string)
	GetAPIVersionForKind(kind string) string
	IsDelete() bool
}

type ServiceActionProvider interface {
	Terminate(error)
	Log(msg string)
	Provision(CompConfigPair) ([]patterns.DeploymentMessagePerContext, error)
	GetRegistry() *registry.RegistryManager
	Persist(string, core.Service, bool) error
	DryRun([]v1beta1.Component) (map[string]map[string]core.DryRunResponseWrapper, error)
	Mutate(*core.Pattern) //Uses pre-defined policies/configuration to mutate the pattern
}
