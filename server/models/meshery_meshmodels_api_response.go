package models

import "github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"

// API response model for meshmodel models API
type MeshmodelsAPIResponse struct {
	Count  int64            `json:"total_count"`
	Models []v1alpha1.Model `json:"models"`
}

// API response model for meshmodel components API
type MeshmodelComponentsAPIResponse struct {
	Count      int64                          `json:"total_count"`
	Components []v1alpha1.ComponentDefinition `json:"components"`
}

// API response model for meshmodel relationships API
type MeshmodelRelationshipsAPIResponse struct {
	Count         int64                             `json:"total_count"`
	Relationships []v1alpha1.RelationshipDefinition `json:"relationships"`
}
