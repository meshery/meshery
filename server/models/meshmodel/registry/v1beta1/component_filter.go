package v1beta1

import (
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/meshmodel/entity"
)

// ComponentFilter is used to filter meshmodel components during retrieval.
type ComponentFilter struct {
	Name         string `json:"name,omitempty"`
	DisplayName  string `json:"display_name,omitempty"`
	CategoryName string `json:"category_name,omitempty"`
	ModelName    string `json:"model_name,omitempty"`
	Version      string `json:"version,omitempty"`
	APIVersion   string `json:"api_version,omitempty"`
	Limit        int    `json:"limit,omitempty"`
	Offset       int    `json:"offset,omitempty"`
	OrderOn      string `json:"order_on,omitempty"`
	Sort         string `json:"sort,omitempty"`
	Greedy       bool   `json:"greedy,omitempty"`
	Annotations  string `json:"annotations,omitempty"`
	Trim         bool   `json:"trim,omitempty"`
}

func (f *ComponentFilter) Create(_ map[string]interface{}) {
	// This method intentionally left empty.
	// It's only needed to satisfy the entity.Filter interface.
}

func (f *ComponentFilter) Get(db *database.Handler) ([]entity.Entity, int64, int, error) {
	// For now, just return an empty list. We'll hook the actual filtering logic later.
	return []entity.Entity{}, 0, 0, nil
}

func (f *ComponentFilter) GetById(db *database.Handler) (entity.Entity, error) {
	// For now, just return nil until we implement the actual lookup logic.
	return nil, nil
}
