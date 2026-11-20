package helpers

import (
	"testing"

	"github.com/meshery/schemas/models/v1alpha3/relationship"
)

func TestNormalizeRelationshipDefinition(t *testing.T) {
	tests := []struct {
		name            string
		input           *relationship.RelationshipDefinition
		expectedKind    relationship.RelationshipDefinitionKind
		expectedType    string
		expectedSubType string
	}{
		{
			name: "mixed case edge and type",
			input: &relationship.RelationshipDefinition{
				Kind:             relationship.RelationshipDefinitionKind("Edge"),
				RelationshipType: "Non-Binding",
				SubType:          "Inventory",
			},
			expectedKind:    relationship.RelationshipDefinitionKind("edge"),
			expectedType:    "non-binding",
			expectedSubType: "inventory",
		},
		{
			name: "uppercase hierarchical parent alias",
			input: &relationship.RelationshipDefinition{
				Kind:             relationship.RelationshipDefinitionKind("HIERARCHICAL"),
				RelationshipType: "PARENT",
				SubType:          "ALIAS",
			},
			expectedKind:    relationship.RelationshipDefinitionKind("hierarchical"),
			expectedType:    "parent",
			expectedSubType: "alias",
		},
		{
			name: "already lowercased",
			input: &relationship.RelationshipDefinition{
				Kind:             relationship.RelationshipDefinitionKind("edge"),
				RelationshipType: "sibling",
				SubType:          "match-labels",
			},
			expectedKind:    relationship.RelationshipDefinitionKind("edge"),
			expectedType:    "sibling",
			expectedSubType: "match-labels",
		},
		{
			name:  "nil relationship definition does not panic",
			input: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NormalizeRelationshipDefinition(tt.input)
			if tt.input == nil {
				return
			}
			if tt.input.Kind != tt.expectedKind {
				t.Errorf("Kind = %q, want %q", tt.input.Kind, tt.expectedKind)
			}
			if tt.input.RelationshipType != tt.expectedType {
				t.Errorf("RelationshipType = %q, want %q", tt.input.RelationshipType, tt.expectedType)
			}
			if tt.input.SubType != tt.expectedSubType {
				t.Errorf("SubType = %q, want %q", tt.input.SubType, tt.expectedSubType)
			}
		})
	}
}
