package models

import (
	"github.com/layer5io/meshsync/pkg/model"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

type MeshSyncResourcesAPIResponse struct {
	Page       int                        `json:"page"`
	PageSize   int                        `json:"page_size"`
	TotalCount int64                      `json:"total_count"`
	Resources  []model.KubernetesResource `json:"resources"`
	Design     pattern.PatternFile        `json:"design"`
}

type MeshSyncResourcesSummaryAPIResponse struct {
	Kinds []struct {
		Kind  string
		Model string
		Count int64
	} `json:"kinds"`
	Namespaces []string                   `json:"namespaces"`
	Labels     []model.KubernetesKeyValue `json:"labels"`
}
