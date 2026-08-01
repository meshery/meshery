package model

import (
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
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
