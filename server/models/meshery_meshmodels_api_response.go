package models

import "github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"

// API response model for meshmodel models API
type MeshmodelsAPIResponse struct {
	Page        int 			 `json:"page"`
	PageSize 	int 			 `json:"page_size"`	
	Count       int64            `json:"total_count"`
	Models      []v1alpha1.Model `json:"models"`
}

// API response model for meshmodel models API that contains the number of duplicates for each model
type MeshmodelsDuplicateAPIResponse struct {
	Page        int 			 			`json:"page"`
	PageSize 	int 			 			`json:"page_size"`	
	Count       int64            			`json:"total_count"`
	Models      []DuplicateResponseModels 	`json:"models"`
}

// API response model for meshmodel components API
type MeshmodelComponentsAPIResponse struct {
	Page        int 			 			   `json:"page"`
	PageSize 	int 			 			   `json:"page_size"`	
	Count       int64                          `json:"total_count"`
	Components  []v1alpha1.ComponentDefinition `json:"components"`
}

// API response model for meshmodel components API that contains the number of duplicates for each component
type MeshmodelComponentsDuplicateAPIResponse struct {
	Page        int 			 			   `json:"page"`
	PageSize 	int 			 			   `json:"page_size"`
	Count       int64                          `json:"total_count"`
	Components  []DuplicateResponseComponent   `json:"components"`
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

type DuplicateResponseComponent struct {
	v1alpha1.ComponentDefinition
	Duplicates int `json:"duplicates"`
}

type DuplicateResponseModels struct {
	v1alpha1.Model
	Duplicates int `json:"duplicates"`
}

func FindDuplicateComponents(components []v1alpha1.ComponentDefinition) []DuplicateResponseComponent {
	set := make(map[string]int)

	for _, comp := range components {
		key := comp.Kind + "@" + comp.APIVersion + "@" + comp.Model.Name
		set[key]++
	}

	var comps []DuplicateResponseComponent

	for _, comp := range components {
		key := comp.Kind + "@" + comp.APIVersion + "@" + comp.Model.Name

		comps = append(comps, DuplicateResponseComponent{
			ComponentDefinition: comp,
			Duplicates: set[key],
		})
	}

	return comps
}

func FindDuplicateModels(models []v1alpha1.Model) []DuplicateResponseModels {
	set := make(map[string]int)

	for _, model := range models {
		key := model.Name + "@" + model.Version
		set[key]++
	}

	var mods []DuplicateResponseModels

	for _, model := range models {
		key := model.Name + "@" + model.Version

		mods = append(mods, DuplicateResponseModels{
			Model: model,
			Duplicates: set[key],
		})
	}

	return mods
}