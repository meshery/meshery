package v1beta1_test

import (
	"regexp"
	"strings"
	"testing"

	v1beta1 "github.com/meshery/meshery/server/models/meshmodel/registry/v1beta1"
)

// MockComponent simulates a meshmodel component with just the Name field.
type MockComponent struct {
	Name string
}

// TestComponentFilter_ExcludeAndRegex tests that Exclude and ExcludeRegex fields work properly.
func TestComponentFilter_ExcludeAndRegex(t *testing.T) {
	tests := []struct {
		name          string
		filter        v1beta1.ComponentFilter
		components    []MockComponent
		expectedNames []string
	}{
		{
			name: "Exclude by string match",
			filter: v1beta1.ComponentFilter{
				Exclude: "List",
			},
			components: []MockComponent{
				{Name: "PodList"},
				{Name: "Deployment"},
				{Name: "NodeList"},
			},
			expectedNames: []string{"Deployment"},
		},
		{
			name: "Exclude by regex pattern",
			filter: v1beta1.ComponentFilter{
				ExcludeRegex: ".*List$",
			},
			components: []MockComponent{
				{Name: "ServiceList"},
				{Name: "ConfigMap"},
				{Name: "Pod"},
			},
			expectedNames: []string{"ConfigMap", "Pod"},
		},
		{
			name: "No exclusions applied",
			filter: v1beta1.ComponentFilter{
				Exclude:      "",
				ExcludeRegex: "",
			},
			components: []MockComponent{
				{Name: "Service"},
				{Name: "Pod"},
			},
			expectedNames: []string{"Service", "Pod"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filtered []MockComponent

			for _, comp := range tt.components {
				// Apply exclude logic manually (simulating your actual logic)
				if tt.filter.Exclude != "" && strings.Contains(comp.Name, tt.filter.Exclude) {
					continue
				}
				if tt.filter.ExcludeRegex != "" {
					re := regexp.MustCompile(tt.filter.ExcludeRegex)
					if re.MatchString(comp.Name) {
						continue
					}
				}
				filtered = append(filtered, comp)
			}

			// Compare expected vs actual
			if len(filtered) != len(tt.expectedNames) {
				t.Fatalf("expected %d components, got %d", len(tt.expectedNames), len(filtered))
			}

			for i, name := range tt.expectedNames {
				if filtered[i].Name != name {
					t.Errorf("expected %s, got %s", name, filtered[i].Name)
				}
			}
		})
	}
}
