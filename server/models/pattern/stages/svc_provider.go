package stages

import (
	"github.com/gofrs/uuid"
    meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/patterns"
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
	Persist(string, patterns.Service, bool) error
	DryRun([]v1alpha1.Component) (map[string]map[string]patterns.DryRunResponseWrapper, error)
	Mutate(*patterns.Pattern) //Uses pre-defined policies/configuration to mutate the pattern
}
