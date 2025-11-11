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

	// 🆕 Added for exclusion filtering
	Exclude      string `json:"exclude,omitempty"`
	ExcludeRegex string `json:"exclude_regex,omitempty"`
}

func (f *ComponentFilter) Create(_ map[string]interface{}) {}
func (f *ComponentFilter) Get(db *database.Handler) ([]entity.Entity, int64, int, error) {
	return []entity.Entity{}, 0, 0, nil
}
func (f *ComponentFilter) GetById(db *database.Handler) (entity.Entity, error) {
	return nil, nil
}
