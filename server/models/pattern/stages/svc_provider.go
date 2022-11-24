package stages

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/pattern/core"
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
	Provision(CompConfigPair) (string, error)
	Persist(string, core.Service, bool) error
	Mutate(*core.Pattern) //Uses pre-defined policies/configuration to mutate the pattern
}
