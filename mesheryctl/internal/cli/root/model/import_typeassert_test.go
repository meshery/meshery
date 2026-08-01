package model

import (
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
)

// buildEntityTypeLine receives names/entityTypes as []interface{} decoded from a
// registry API JSON response. Its sibling builders (buildDescription,
// buildDescriptionList) guard every element with the comma-ok form, but
// buildEntityTypeLine asserted .(string) directly, so a non-string element
// panicked the mesheryctl model import error display. This test asserts it
// degrades gracefully instead of panicking.
func TestBuildEntityTypeLineHandlesNonStringEntries(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	tests := []struct {
		name        string
		names       []interface{}
		entityTypes []interface{}
		modelName   string
	}{
		{"non-string name with model filter", []interface{}{123}, []interface{}{"unknown"}, "mymodel"},
		{"non-string entityType", []interface{}{"mymodel"}, []interface{}{42}, "mymodel"},
		{"mixed non-string entries", []interface{}{true, "mymodel"}, []interface{}{nil, "unknown"}, "mymodel"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Must not panic on malformed (non-string) interface elements.
			_ = buildEntityTypeLine(tt.names, tt.entityTypes, "desc", "", "", tt.modelName)
		})
	}
}

// displaySuccessfulRelationships walks SuccessfulRelationships ([]map[string]interface{})
// straight from the registry API JSON. It asserted .(string)/.(map)/.([]interface{})
// directly and indexed from[0]/to[0] without a length check, so a malformed or
// sparse relationship entry panicked mesheryctl. This asserts it skips bad entries
// instead of crashing, matching displaySuccessfulComponents in the same file.
func TestDisplaySuccessfulRelationshipsHandlesMalformed(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	tests := []struct {
		name string
		rel  map[string]interface{}
	}{
		{
			name: "empty from/to lists (index out of range)",
			rel: map[string]interface{}{
				"Model": "mymodel", "Kind": "a", "Subtype": "b", "RelationshipType": "c",
				"Selectors": []interface{}{
					map[string]interface{}{"allow": map[string]interface{}{
						"from": []interface{}{}, "to": []interface{}{},
					}},
				},
			},
		},
		{
			name: "non-string kind with a valid selector",
			rel: map[string]interface{}{
				"Model": "mymodel", "Kind": 123, "Subtype": "b", "RelationshipType": "c",
				"Selectors": []interface{}{
					map[string]interface{}{"allow": map[string]interface{}{
						"from": []interface{}{map[string]interface{}{"kind": "X"}},
						"to":   []interface{}{map[string]interface{}{"kind": "Y"}},
					}},
				},
			},
		},
		{
			name: "selector is not a map",
			rel: map[string]interface{}{
				"Model": "mymodel", "Kind": "a", "Subtype": "b", "RelationshipType": "c",
				"Selectors": []interface{}{"not-a-map"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &models.RegistryAPIResponse{
				EntityTypeSummary: models.EntityTypeSummary{
					SuccessfulRelationships: []map[string]interface{}{tt.rel},
				},
			}
			// Must not panic on malformed relationship data.
			displaySuccessfulRelationships(resp, "mymodel")
		})
	}
}
