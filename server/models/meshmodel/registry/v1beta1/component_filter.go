package v1beta1

import (
	"fmt"
	"strings"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	"github.com/meshery/schemas/models/v1beta1/component"
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

func (f *ComponentFilter) Create(_ map[string]interface{}) {
	// This method intentionally left empty.
	// It's only needed to satisfy the entity.Filter interface.
}

func (f *ComponentFilter) Get(db *database.Handler) ([]entity.Entity, int64, int, error) {
	if db == nil {
		return nil, 0, 0, models.ErrDBConnection
	}
	var err error
	query := db.Model(&component.ComponentDefinition{}).
		Preload("Model").
		Preload("Model.Category")

	if f.Name != "" {
		if f.Greedy {
			query = query.Where("component ->> 'kind' ILIKE ?", "%"+f.Name+"%")
		} else {
			query = query.Where("component ->> 'kind' = ?", f.Name)
		}
	}
	if f.ModelName != "" {
		query = query.Joins("JOIN model_dbs on component_definitions.model_id = model_dbs.id").
			Where("model_dbs.name = ?", f.ModelName)
	}
	if f.CategoryName != "" {
		query = query.Joins("JOIN model_dbs on component_definitions.model_id = model_dbs.id").
			Joins("JOIN category_dbs on model_dbs.category_id = category_dbs.id").
			Where("category_dbs.name = ?", f.CategoryName)
	}

	if f.APIVersion != "" {
		query = query.Where("component ->> 'version' = ?", f.APIVersion)
	}

	if f.Version != "" {
		query = query.Joins("JOIN model_dbs on component_definitions.model_id = model_dbs.id").
			Where("model_dbs.version = ?", f.Version)
	}
	if f.Annotations != "" {
		switch f.Annotations {
		case "true":
			query = query.Where("component -> 'metadata' ->> 'isAnnotation' = 'true'")
		case "false":
			query = query.Where("component -> 'metadata' ->> 'isAnnotation' = 'false' OR component -> 'metadata' ->> 'isAnnotation' IS NULL")
		}
	}

	if f.Exclude != "" {
		excludeLower := strings.ToLower(f.Exclude)
		query = query.Where("LOWER(component ->> 'kind') NOT LIKE ?", "%"+excludeLower+"%").
			Where("LOWER(display_name) NOT LIKE ?", "%"+excludeLower+"%")
	}

	if f.ExcludeRegex != "" {
		switch db.Dialector.Name() {
		case "postgres":
			query = query.Where("component ->> 'kind' !~ ?", f.ExcludeRegex).
				Where("display_name !~ ?", f.ExcludeRegex)
		case "mysql":
			query = query.Where("component ->> 'kind' NOT REGEXP ?", f.ExcludeRegex).
				Where("display_name NOT REGEXP ?", f.ExcludeRegex)
		}
	}

	if f.Trim {
		query = query.Select("id", "component", "display_name", "model_id")
	}

	var count int64
	query.Count(&count)

	if f.OrderOn != "" {
		query = query.Order(fmt.Sprintf("%s %s", f.OrderOn, f.Sort))
	}

	if f.Limit != 0 {
		query = query.Limit(f.Limit)
	}

	if f.Offset != 0 {
		query = query.Offset(f.Offset)
	}

	var comps []component.ComponentDefinition
	err = query.Find(&comps).Error
	if err != nil {
		return nil, 0, 0, err
	}

	var entities []entity.Entity
	for i := range comps {
		entities = append(entities, &comps[i])
	}

	return entities, count, 0, nil
}



func (f *ComponentFilter) GetById(db *database.Handler) (entity.Entity, error) {
	if db == nil {
		return nil, models.ErrDBConnection
	}
	var comp component.ComponentDefinition
	err := db.First(&comp, f.Name).Error
	return &comp, err
}
