package v1beta1

import (
	"regexp"
	"strings"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	"github.com/meshery/schemas/models/v1beta1/component"

	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
)

// ComponentFilterWrapper is a wrapper around ComponentFilter that adds exclusion filtering.
type ComponentFilterWrapper struct {
	ComponentFilter *regv1beta1.ComponentFilter
	Exclude      string `json:"exclude,omitempty"`
	ExcludeRegex string `json:"exclude_regex,omitempty"`
}

// Create implements entity.Filter
func (f *ComponentFilterWrapper) Create(data map[string]interface{}) {
	f.ComponentFilter.Create(data)
}

// GetById implements entity.Filter
func (f *ComponentFilterWrapper) GetById(db *database.Handler) (entity.Entity, error) {
	return f.ComponentFilter.GetById(db)
}

// Get returns a list of components after applying exclusion filtering.
func (f *ComponentFilterWrapper) Get(db *database.Handler) ([]entity.Entity, int64, int, error) {
	entities, count, totalCount, err := f.ComponentFilter.Get(db)
	if err != nil {
		return nil, 0, 0, err
	}

	var comps []component.ComponentDefinition
	for _, entity := range entities {
		if comp, ok := entity.(*component.ComponentDefinition); ok {
			comps = append(comps, *comp)
		}
	}

	comps = filterComponentsByExclusion(comps, f.Exclude, f.ExcludeRegex)

	entities = make([]entity.Entity, len(comps))
	for i, comp := range comps {
		entities[i] = &comp
	}

	return entities, count, totalCount, nil
}

// filterComponentsByExclusion filters components based on exclude and excludeRegex patterns.
func filterComponentsByExclusion(comps []component.ComponentDefinition, exclude, excludeRegex string) []component.ComponentDefinition {
	if exclude == "" && excludeRegex == "" {
		return comps
	}

	var filtered []component.ComponentDefinition
	var regex *regexp.Regexp
	var err error

	if excludeRegex != "" {
		regex, err = regexp.Compile(excludeRegex)
		if err != nil {
			return comps
		}
	}

	excludeLower := strings.ToLower(exclude)

	for _, comp := range comps {
		shouldExclude := false

		kind := comp.Component.Kind
		kindLower := strings.ToLower(kind)

		if exclude != "" && strings.Contains(kindLower, excludeLower) {
			shouldExclude = true
		}

		if !shouldExclude && regex != nil && regex.MatchString(kind) {
			shouldExclude = true
		}

		if !shouldExclude {
			displayName := comp.DisplayName
			displayNameLower := strings.ToLower(displayName)

			if exclude != "" && strings.Contains(displayNameLower, excludeLower) {
				shouldExclude = true
			}

			if !shouldExclude && regex != nil && regex.MatchString(displayName) {
				shouldExclude = true
			}
		}

		if !shouldExclude {
			filtered = append(filtered, comp)
		}
	}

	return filtered
}
