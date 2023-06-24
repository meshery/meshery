package models

import "github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"

// API response model for meshmodel models API
type MeshmodelsAPIResponse struct {
	Page        int 			 `json:"page"`
	PageSize 	int 			 `json:"page_size"`	
	Count       int64            `json:"total_count"`
	Models      []v1alpha1.Model `json:"models"`
}

// API response model for meshmodel components API
type MeshmodelComponentsAPIResponse struct {
	Page        int 			 			   `json:"page"`
	PageSize 	int 			 			   `json:"page_size"`	
	Count       int64                          `json:"total_count"`
	Components  []v1alpha1.ComponentDefinition `json:"components"`
}

// API response model for meshmodel relationships API
type MeshmodelRelationshipsAPIResponse struct {
	Page          int 			 			   		`json:"page"`
	PageSize 	  int 			 			   		`json:"page_size"`	
	Count         int64                             `json:"total_count"`
	Relationships []v1alpha1.RelationshipDefinition `json:"relationships"`
}

// API response model for meshmodel categories API
type MeshmodelCategoriesAPIResponse struct {
	Page          int 			 		`json:"page"`
	PageSize 	  int 			 		`json:"page_size"`	
	Count         int64                 `json:"total_count"`
	Categories 	  []v1alpha1.Category	`json:"categories"`
}
