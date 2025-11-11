package v1beta1

import (
	"testing"

	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/stretchr/testify/assert"
)

func TestFilterComponentsByExclusion(t *testing.T) {
	comps := []component.ComponentDefinition{
		{
			Component: component.Component{
				Kind: "Pod",
			},
			DisplayName: "Pod",
		},
		{
			Component: component.Component{
				Kind: "Service",
			},
			DisplayName: "Service",
		},
		{
			Component: component.Component{
				Kind: "PodList",
			},
			DisplayName: "Pod List",
		},
	}

	tests := []struct {
		name           string
		exclude        string
		excludeRegex   string
		expectedCount  int
		expectedTitles []string
	}{
		{
			name:           "Exclude with string",
			exclude:        "List",
			expectedCount:  2,
			expectedTitles: []string{"Pod", "Service"},
		},
		{
			name:           "Exclude with regex",
			excludeRegex:   ".*List$",
			expectedCount:  2,
			expectedTitles: []string{"Pod", "Service"},
		},
		{
			name:           "No exclusion",
			expectedCount:  3,
			expectedTitles: []string{"Pod", "Service", "PodList"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filteredComps := filterComponentsByExclusion(comps, tt.exclude, tt.excludeRegex)
			assert.Equal(t, tt.expectedCount, len(filteredComps))

			var titles []string
			for _, c := range filteredComps {
				titles = append(titles, c.Component.Kind)
			}
			assert.ElementsMatch(t, tt.expectedTitles, titles)
		})
	}
}
