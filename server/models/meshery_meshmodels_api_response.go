package models

import (
	models "github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta1/component"
)

// API response model for meshmodel models API
type MeshmodelsAPIResponse struct {
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
	Count    int64                   `json:"total_count"`
	Models   []model.ModelDefinition `json:"models"`
}

// API response model for meshmodel models API that contains the number of duplicates for each model
type MeshmodelsDuplicateAPIResponse struct {
	Page     int                       `json:"page"`
	PageSize int                       `json:"page_size"`
	Count    int64                     `json:"total_count"`
	Models   []DuplicateResponseModels `json:"models"`
}

// API response model for meshmodel components API
type MeshmodelComponentsAPIResponse struct {
	Page       int                         `json:"page"`
	PageSize   int                         `json:"page_size"`
	Count      int64                       `json:"total_count"`
	Components []component.ComponentDefinition `json:"components"`
}

// API response model for meshmodel components API that contains the number of duplicates for each component
type MeshmodelComponentsDuplicateAPIResponse struct {
	Page       int                          `json:"page"`
	PageSize   int                          `json:"page_size"`
	Count      int64                        `json:"total_count"`
	Components []DuplicateResponseComponent `json:"components"`
}

// API response model for meshmodel relationships API
type MeshmodelRelationshipsAPIResponse struct {
	Page          int             `json:"page"`
	PageSize      int             `json:"page_size"`
	Count         int64           `json:"total_count"`
	Relationships []entity.Entity `json:"relationships"`
}

// API response model for meshmodel categories API
type MeshmodelCategoriesAPIResponse struct {
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	Count      int64           `json:"total_count"`
	Categories []entity.Entity `json:"categories"`
}

// API response model for meshmodel relationships API
type MeshmodelPoliciesAPIResponse struct {
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Count    int64           `json:"total_count"`
	Policies []entity.Entity `json:"policies"`
}

type DuplicateResponseComponent struct {
	component.ComponentDefinition
	Duplicates int `json:"duplicates"`
}

type DuplicateResponseModels struct {
	model.ModelDefinition
	Duplicates int `json:"duplicates"`
}

type MeshmodelRegistrantsAPIResponse struct {
	Page        int                                      `json:"page"`
	PageSize    int                                      `json:"page_size"`
	Count       int64                                    `json:"total_count"`
	Registrants []models.MeshModelHostsWithEntitySummary `json:"registrants"`
}

func FindDuplicateComponents(components []component.ComponentDefinition) []DuplicateResponseComponent {
	set := make(map[string]int)

	for _, comp := range components {
		key := comp.Component.Kind + "@" + comp.Component.Version + "@" + comp.Model.Name
		set[key]++
	}

	var comps []DuplicateResponseComponent

	for _, comp := range components {
		key := comp.Component.Kind + "@" + comp.Component.Version + "@" + comp.Model.Name

		comps = append(comps, DuplicateResponseComponent{
			ComponentDefinition: comp,
			Duplicates:          set[key] - 1,
		})
	}

	return comps
}

func FindDuplicateModels(models []model.ModelDefinition) []DuplicateResponseModels {
	set := make(map[string]int)

	for _, model := range models {
		key := model.Name + "@" + model.Version
		set[key]++
	}

	var mods []DuplicateResponseModels

	for _, model := range models {
		key := model.Name + "@" + model.Version

		mods = append(mods, DuplicateResponseModels{
			ModelDefinition: model,
			Duplicates:      set[key] - 1,
		})
	}

	return mods
}
