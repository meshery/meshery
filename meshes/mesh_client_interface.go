package meshes

import "context"

type MeshClient interface {
	// DeleteAllCreatedResources(ctx context.Context, namespace string)
	// DeleteResource(ctx context.Context, overallType, namespace, resName string) error
	MeshName() string
	ApplyRule(ctx context.Context, opName, username, namespace string) error
	Operations(ctx context.Context) (map[string]string, error)
}
