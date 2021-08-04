package stages

import (
	"github.com/gofrs/uuid"
)

type ServiceInfoProvider interface {
	GetMesheryPatternResource(
		name string,
		namespace string,
		typ string,
		oamType string,
	) (ID *uuid.UUID, err error)
	GetServiceMesh() (name string, version string)
	GetAPIVersionForKind() string
}
