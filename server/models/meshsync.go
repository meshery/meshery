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
		Count int64
	} `json:"kinds"`
	Namespaces   []string                   `json:"namespaces"`
	Labels       []model.KubernetesKeyValue `json:"labels"`
	PodSummaries []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	} `json:"podSummaries"`
	NodeSummaries []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	} `json:"nodeSummaries"`
	Usage []struct {
		Resource   string  `json:"resource"`
		Percentage float64 `json:"percentage"`
	} `json:"usage"`
}
