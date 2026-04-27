package models

import (
	"encoding/json"

	models "github.com/meshery/meshkit/models/meshmodel/core/v1beta1"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta3/component"
)

// MeshmodelsAPIResponse is the Meshery-local response envelope for
// `GET /api/meshmodels/models`. Canonical wire form is camelCase per the
// identifier-naming migration; legacy snake_case keys (`page_size`,
// `total_count`) are emitted alongside for the deprecation window so
// external consumers still reading the old spellings keep working.
// Unmarshal likewise accepts either spelling so any inbound serialization
// (round-trip via clients) tolerates both.
//
// Once every known consumer has migrated off the snake_case keys, drop
// MarshalJSON/UnmarshalJSON and keep only the camelCase struct tags.
//
// TODO: Move to schemas
type MeshmodelsAPIResponse struct {
	Page       int                     `json:"page"`
	PageSize   int                     `json:"pageSize"`
	TotalCount int64                   `json:"totalCount"`
	Models     []model.ModelDefinition `json:"models"`
}

func (p MeshmodelsAPIResponse) MarshalJSON() ([]byte, error) {
	type alias MeshmodelsAPIResponse
	return json.Marshal(struct {
		alias
		PageSizeLegacy   int   `json:"page_size"`
		TotalCountLegacy int64 `json:"total_count"`
	}{
		alias:            alias(p),
		PageSizeLegacy:   p.PageSize,
		TotalCountLegacy: p.TotalCount,
	})
}

func (p *MeshmodelsAPIResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		Page             int                     `json:"page"`
		PageSize         *int                    `json:"pageSize"`
		PageSizeLegacy   *int                    `json:"page_size"`
		TotalCount       *int64                  `json:"totalCount"`
		TotalCountLegacy *int64                  `json:"total_count"`
		Models           []model.ModelDefinition `json:"models"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Page = raw.Page
	p.Models = raw.Models
	// Match stdlib json.Unmarshal semantics: fields absent from the input
	// reset to their zero value on the destination (important when callers
	// reuse the same MeshmodelsAPIResponse across decodes).
	p.PageSize = 0
	switch {
	case raw.PageSize != nil:
		p.PageSize = *raw.PageSize
	case raw.PageSizeLegacy != nil:
		p.PageSize = *raw.PageSizeLegacy
	}
	p.TotalCount = 0
	switch {
	case raw.TotalCount != nil:
		p.TotalCount = *raw.TotalCount
	case raw.TotalCountLegacy != nil:
		p.TotalCount = *raw.TotalCountLegacy
	}
	return nil
}

// MeshmodelsDuplicateAPIResponse is the response envelope for the
// duplicates variant of the meshmodels list endpoint. Wire-form rules
// mirror MeshmodelsAPIResponse above.
type MeshmodelsDuplicateAPIResponse struct {
	Page       int                       `json:"page"`
	PageSize   int                       `json:"pageSize"`
	TotalCount int64                     `json:"totalCount"`
	Models     []DuplicateResponseModels `json:"models"`
}

func (p MeshmodelsDuplicateAPIResponse) MarshalJSON() ([]byte, error) {
	type alias MeshmodelsDuplicateAPIResponse
	return json.Marshal(struct {
		alias
		PageSizeLegacy   int   `json:"page_size"`
		TotalCountLegacy int64 `json:"total_count"`
	}{
		alias:            alias(p),
		PageSizeLegacy:   p.PageSize,
		TotalCountLegacy: p.TotalCount,
	})
}

func (p *MeshmodelsDuplicateAPIResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		Page             int                       `json:"page"`
		PageSize         *int                      `json:"pageSize"`
		PageSizeLegacy   *int                      `json:"page_size"`
		TotalCount       *int64                    `json:"totalCount"`
		TotalCountLegacy *int64                    `json:"total_count"`
		Models           []DuplicateResponseModels `json:"models"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Page = raw.Page
	p.Models = raw.Models
	p.PageSize = 0
	switch {
	case raw.PageSize != nil:
		p.PageSize = *raw.PageSize
	case raw.PageSizeLegacy != nil:
		p.PageSize = *raw.PageSizeLegacy
	}
	p.TotalCount = 0
	switch {
	case raw.TotalCount != nil:
		p.TotalCount = *raw.TotalCount
	case raw.TotalCountLegacy != nil:
		p.TotalCount = *raw.TotalCountLegacy
	}
	return nil
}

// MeshmodelComponentsAPIResponse is the response envelope for the
// meshmodel components list endpoint. Wire-form rules mirror
// MeshmodelsAPIResponse above.
type MeshmodelComponentsAPIResponse struct {
	Page       int                             `json:"page"`
	PageSize   int                             `json:"pageSize"`
	TotalCount int64                           `json:"totalCount"`
	Components []component.ComponentDefinition `json:"components"`
}

func (p MeshmodelComponentsAPIResponse) MarshalJSON() ([]byte, error) {
	type alias MeshmodelComponentsAPIResponse
	return json.Marshal(struct {
		alias
		PageSizeLegacy   int   `json:"page_size"`
		TotalCountLegacy int64 `json:"total_count"`
	}{
		alias:            alias(p),
		PageSizeLegacy:   p.PageSize,
		TotalCountLegacy: p.TotalCount,
	})
}

func (p *MeshmodelComponentsAPIResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		Page             int                             `json:"page"`
		PageSize         *int                            `json:"pageSize"`
		PageSizeLegacy   *int                            `json:"page_size"`
		TotalCount       *int64                          `json:"totalCount"`
		TotalCountLegacy *int64                          `json:"total_count"`
		Components       []component.ComponentDefinition `json:"components"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Page = raw.Page
	p.Components = raw.Components
	p.PageSize = 0
	switch {
	case raw.PageSize != nil:
		p.PageSize = *raw.PageSize
	case raw.PageSizeLegacy != nil:
		p.PageSize = *raw.PageSizeLegacy
	}
	p.TotalCount = 0
	switch {
	case raw.TotalCount != nil:
		p.TotalCount = *raw.TotalCount
	case raw.TotalCountLegacy != nil:
		p.TotalCount = *raw.TotalCountLegacy
	}
	return nil
}

// MeshmodelComponentsDuplicateAPIResponse is the response envelope for
// the duplicates variant of the components list endpoint. Wire-form rules
// mirror MeshmodelsAPIResponse above.
type MeshmodelComponentsDuplicateAPIResponse struct {
	Page       int                          `json:"page"`
	PageSize   int                          `json:"pageSize"`
	TotalCount int64                        `json:"totalCount"`
	Components []DuplicateResponseComponent `json:"components"`
}

func (p MeshmodelComponentsDuplicateAPIResponse) MarshalJSON() ([]byte, error) {
	type alias MeshmodelComponentsDuplicateAPIResponse
	return json.Marshal(struct {
		alias
		PageSizeLegacy   int   `json:"page_size"`
		TotalCountLegacy int64 `json:"total_count"`
	}{
		alias:            alias(p),
		PageSizeLegacy:   p.PageSize,
		TotalCountLegacy: p.TotalCount,
	})
}

func (p *MeshmodelComponentsDuplicateAPIResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		Page             int                          `json:"page"`
		PageSize         *int                         `json:"pageSize"`
		PageSizeLegacy   *int                         `json:"page_size"`
		TotalCount       *int64                       `json:"totalCount"`
		TotalCountLegacy *int64                       `json:"total_count"`
		Components       []DuplicateResponseComponent `json:"components"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Page = raw.Page
	p.Components = raw.Components
	p.PageSize = 0
	switch {
	case raw.PageSize != nil:
		p.PageSize = *raw.PageSize
	case raw.PageSizeLegacy != nil:
		p.PageSize = *raw.PageSizeLegacy
	}
	p.TotalCount = 0
	switch {
	case raw.TotalCount != nil:
		p.TotalCount = *raw.TotalCount
	case raw.TotalCountLegacy != nil:
		p.TotalCount = *raw.TotalCountLegacy
	}
	return nil
}

// MeshmodelRelationshipsAPIResponse is the response envelope for the
// relationships list endpoint. Wire-form rules mirror MeshmodelsAPIResponse.
type MeshmodelRelationshipsAPIResponse struct {
	Page          int             `json:"page"`
	PageSize      int             `json:"pageSize"`
	TotalCount    int64           `json:"totalCount"`
	Relationships []entity.Entity `json:"relationships"`
}

func (p MeshmodelRelationshipsAPIResponse) MarshalJSON() ([]byte, error) {
	type alias MeshmodelRelationshipsAPIResponse
	return json.Marshal(struct {
		alias
		PageSizeLegacy   int   `json:"page_size"`
		TotalCountLegacy int64 `json:"total_count"`
	}{
		alias:            alias(p),
		PageSizeLegacy:   p.PageSize,
		TotalCountLegacy: p.TotalCount,
	})
}

func (p *MeshmodelRelationshipsAPIResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		Page             int             `json:"page"`
		PageSize         *int            `json:"pageSize"`
		PageSizeLegacy   *int            `json:"page_size"`
		TotalCount       *int64          `json:"totalCount"`
		TotalCountLegacy *int64          `json:"total_count"`
		Relationships    []entity.Entity `json:"relationships"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Page = raw.Page
	p.Relationships = raw.Relationships
	p.PageSize = 0
	switch {
	case raw.PageSize != nil:
		p.PageSize = *raw.PageSize
	case raw.PageSizeLegacy != nil:
		p.PageSize = *raw.PageSizeLegacy
	}
	p.TotalCount = 0
	switch {
	case raw.TotalCount != nil:
		p.TotalCount = *raw.TotalCount
	case raw.TotalCountLegacy != nil:
		p.TotalCount = *raw.TotalCountLegacy
	}
	return nil
}

type EntityCount struct {
	CompCount     int `json:"comp_count"`
	RelCount      int `json:"rel_count"`
	ModelCount    int `json:"model_count"`
	ErrCompCount  int `json:"err_comp_count"`
	ErrRelCount   int `json:"err_rel_count"`
	ErrModelCount int `json:"err_model_count"`
	TotalErrCount int `json:"total_err_count"`
}

type EntityTypeSummary struct {
	SuccessfulComponents            []map[string]interface{} `json:"successful_components"`
	SuccessfulRelationships         []map[string]interface{} `json:"successful_relationships"`
	SuccessfulModels                []map[string]interface{} `json:"successful_models"`
	UnsuccessfulEntityNameWithError []interface{}            `json:"unsuccessful_component_names"`
}

type RegistryAPIResponse struct {
	EntityCount       EntityCount       `json:"entity_count"`
	ErrMsg            string            `json:"err_msg"`
	EntityTypeSummary EntityTypeSummary `json:"entity_type_summary"`
	ModelName         []string          `json:"model_name"`
}

// MeshmodelCategoriesAPIResponse is the response envelope for the
// categories list endpoint. Wire-form rules mirror MeshmodelsAPIResponse.
type MeshmodelCategoriesAPIResponse struct {
	Page       int             `json:"page"`
	PageSize   int             `json:"pageSize"`
	TotalCount int64           `json:"totalCount"`
	Categories []entity.Entity `json:"categories"`
}

func (p MeshmodelCategoriesAPIResponse) MarshalJSON() ([]byte, error) {
	type alias MeshmodelCategoriesAPIResponse
	return json.Marshal(struct {
		alias
		PageSizeLegacy   int   `json:"page_size"`
		TotalCountLegacy int64 `json:"total_count"`
	}{
		alias:            alias(p),
		PageSizeLegacy:   p.PageSize,
		TotalCountLegacy: p.TotalCount,
	})
}

func (p *MeshmodelCategoriesAPIResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		Page             int             `json:"page"`
		PageSize         *int            `json:"pageSize"`
		PageSizeLegacy   *int            `json:"page_size"`
		TotalCount       *int64          `json:"totalCount"`
		TotalCountLegacy *int64          `json:"total_count"`
		Categories       []entity.Entity `json:"categories"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Page = raw.Page
	p.Categories = raw.Categories
	p.PageSize = 0
	switch {
	case raw.PageSize != nil:
		p.PageSize = *raw.PageSize
	case raw.PageSizeLegacy != nil:
		p.PageSize = *raw.PageSizeLegacy
	}
	p.TotalCount = 0
	switch {
	case raw.TotalCount != nil:
		p.TotalCount = *raw.TotalCount
	case raw.TotalCountLegacy != nil:
		p.TotalCount = *raw.TotalCountLegacy
	}
	return nil
}

// MeshmodelPoliciesAPIResponse is the response envelope for the
// policies list endpoint. Wire-form rules mirror MeshmodelsAPIResponse.
type MeshmodelPoliciesAPIResponse struct {
	Page       int             `json:"page"`
	PageSize   int             `json:"pageSize"`
	TotalCount int64           `json:"totalCount"`
	Policies   []entity.Entity `json:"policies"`
}

func (p MeshmodelPoliciesAPIResponse) MarshalJSON() ([]byte, error) {
	type alias MeshmodelPoliciesAPIResponse
	return json.Marshal(struct {
		alias
		PageSizeLegacy   int   `json:"page_size"`
		TotalCountLegacy int64 `json:"total_count"`
	}{
		alias:            alias(p),
		PageSizeLegacy:   p.PageSize,
		TotalCountLegacy: p.TotalCount,
	})
}

func (p *MeshmodelPoliciesAPIResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		Page             int             `json:"page"`
		PageSize         *int            `json:"pageSize"`
		PageSizeLegacy   *int            `json:"page_size"`
		TotalCount       *int64          `json:"totalCount"`
		TotalCountLegacy *int64          `json:"total_count"`
		Policies         []entity.Entity `json:"policies"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Page = raw.Page
	p.Policies = raw.Policies
	p.PageSize = 0
	switch {
	case raw.PageSize != nil:
		p.PageSize = *raw.PageSize
	case raw.PageSizeLegacy != nil:
		p.PageSize = *raw.PageSizeLegacy
	}
	p.TotalCount = 0
	switch {
	case raw.TotalCount != nil:
		p.TotalCount = *raw.TotalCount
	case raw.TotalCountLegacy != nil:
		p.TotalCount = *raw.TotalCountLegacy
	}
	return nil
}

type DuplicateResponseComponent struct {
	component.ComponentDefinition
	Duplicates int `json:"duplicates"`
}

type DuplicateResponseModels struct {
	model.ModelDefinition
	Duplicates int `json:"duplicates"`
}

// MeshmodelRegistrantsAPIResponse is the response envelope for the
// registrants list endpoint. Wire-form rules mirror MeshmodelsAPIResponse.
type MeshmodelRegistrantsAPIResponse struct {
	Page        int                                      `json:"page"`
	PageSize    int                                      `json:"pageSize"`
	TotalCount  int64                                    `json:"totalCount"`
	Registrants []models.MeshModelHostsWithEntitySummary `json:"registrants"`
}

func (p MeshmodelRegistrantsAPIResponse) MarshalJSON() ([]byte, error) {
	type alias MeshmodelRegistrantsAPIResponse
	return json.Marshal(struct {
		alias
		PageSizeLegacy   int   `json:"page_size"`
		TotalCountLegacy int64 `json:"total_count"`
	}{
		alias:            alias(p),
		PageSizeLegacy:   p.PageSize,
		TotalCountLegacy: p.TotalCount,
	})
}

func (p *MeshmodelRegistrantsAPIResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		Page             int                                      `json:"page"`
		PageSize         *int                                     `json:"pageSize"`
		PageSizeLegacy   *int                                     `json:"page_size"`
		TotalCount       *int64                                   `json:"totalCount"`
		TotalCountLegacy *int64                                   `json:"total_count"`
		Registrants      []models.MeshModelHostsWithEntitySummary `json:"registrants"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	p.Page = raw.Page
	p.Registrants = raw.Registrants
	p.PageSize = 0
	switch {
	case raw.PageSize != nil:
		p.PageSize = *raw.PageSize
	case raw.PageSizeLegacy != nil:
		p.PageSize = *raw.PageSizeLegacy
	}
	p.TotalCount = 0
	switch {
	case raw.TotalCount != nil:
		p.TotalCount = *raw.TotalCount
	case raw.TotalCountLegacy != nil:
		p.TotalCount = *raw.TotalCountLegacy
	}
	return nil
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
