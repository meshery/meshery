package meshes

import "context"

// MeshClient is the spec for all mesh implementations to follow
type MeshClient interface {
	MeshName() string
	ApplyRule(ctx context.Context, opName, username, namespace string) error
	Operations(ctx context.Context) (map[string]string, error)
}
