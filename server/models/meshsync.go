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

// MeshSyncLabel is one distinct label found on the resources in scope of
// GET /api/system/meshsync/resources/summary.
//
// gorm scans the `DISTINCT kubernetes_key_values.key, kubernetes_key_values.value`
// projection into it, so it carries exactly those two columns. It replaces
// meshsync's model.KubernetesKeyValue on the wire: that type's id, unique_id and
// kind were always empty here because the query never selects them, and its
// value was dropped (omitempty) for labels whose value is empty. The wire shape
// is MeshSyncLabel in meshery/schemas v1beta1/meshsync; switch to the generated
// type once a schemas release carries it.
type MeshSyncLabel struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MeshSyncResourcesSummaryAPIResponse struct {
	Kinds      []MeshSyncResourceKindSummary `json:"kinds"`
	Namespaces []string                      `json:"namespaces"`
	Labels     []MeshSyncLabel               `json:"labels"`
}
