package stages

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/pattern/core"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
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
	Provision(CompConfigPair) (string, error)
	GetRegistry() *meshmodel.RegistryManager
	Persist(string, core.Service, bool) error
	DryRun([]v1alpha1.Component) (map[string]map[string]core.DryRunResponse2, error)
	Mutate(*core.Pattern) //Uses pre-defined policies/configuration to mutate the pattern
}
