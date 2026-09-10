package models

import (
	"github.com/meshery/meshsync/pkg/model"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

type MeshSyncResourcesAPIResponse struct {
	Page       int                        `json:"page"`
	PageSize   int                        `json:"pageSize"`
	TotalCount int64                      `json:"totalCount"`
	Resources  []model.KubernetesResource `json:"resources"`
	Design     pattern.PatternFile        `json:"design"`
}

// MeshSyncResourceKindSummary is one row of the resource census returned by
// GET /api/system/meshsync/resources/summary.
//
// The field names are what gorm scans the `kind, model, count(*) as count`
// aggregate into, so they must not be renamed without updating that query; the
// json tags are what goes on the wire and are camelCase like the rest of the API.
type MeshSyncResourceKindSummary struct {
	Kind  string `json:"kind"`
	Model string `json:"model"`
	Count int64  `json:"count"`
}

type MeshSyncResourcesSummaryAPIResponse struct {
	Kinds      []MeshSyncResourceKindSummary `json:"kinds"`
	Namespaces []string                      `json:"namespaces"`
	Labels     []model.KubernetesKeyValue    `json:"labels"`
}
