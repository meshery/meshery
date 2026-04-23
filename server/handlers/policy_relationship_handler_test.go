package handlers

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/stretchr/testify/assert"
)

func TestParseRelationshipToAlias(t *testing.T) {
	fromID := uuid.Must(uuid.NewV4())
	toID := uuid.Must(uuid.NewV4())
	relID := uuid.Must(uuid.NewV4())

	tests := []struct {
		name   string
		input  relationship.RelationshipDefinition
		wantOk bool
	}{
		{
			name: "wrong subtype returns false",
			input: relationship.RelationshipDefinition{
				SubType: "not-alias",
			},
			wantOk: false,
		},
		{
			name: "nil Selectors returns false",
			input: relationship.RelationshipDefinition{
				SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
				Selectors: nil,
			},
			wantOk: false,
		},
		{
			name: "empty Selectors returns false",
			input: relationship.RelationshipDefinition{
				SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
				Selectors: &relationship.SelectorSet{},
			},
			wantOk: false,
		},
		{
			name: "empty From set returns false",
			input: func() relationship.RelationshipDefinition {
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{},
							To:   []relationship.SelectorItem{{ID: &toID}},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "empty To set returns false",
			input: func() relationship.RelationshipDefinition {
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{{ID: &fromID}},
							To:   []relationship.SelectorItem{},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "nil Patch returns false",
			input: func() relationship.RelationshipDefinition {
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{{ID: &fromID, RelationshipDefinitionSelectorsPatch: nil}},
							To:   []relationship.SelectorItem{{ID: &toID}},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "nil MutatedRef returns false",
			input: func() relationship.RelationshipDefinition {
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{
								{
									ID: &fromID,
									RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
										MutatedRef: nil,
									},
								},
							},
							To: []relationship.SelectorItem{{ID: &toID}},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "empty MutatedRef returns false",
			input: func() relationship.RelationshipDefinition {
				emptyRefs := [][]string{}
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{
								{
									ID: &fromID,
									RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
										MutatedRef: &emptyRefs,
									},
								},
							},
							To: []relationship.SelectorItem{{ID: &toID}},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "nil to.ID returns false",
			input: func() relationship.RelationshipDefinition {
				refs := [][]string{{"configuration", "spec", "containers"}}
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{
								{
									ID: &fromID,
									RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
										MutatedRef: &refs,
									},
								},
							},
							To: []relationship.SelectorItem{{ID: nil}},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "nil from.ID returns false",
			input: func() relationship.RelationshipDefinition {
				refs := [][]string{{"configuration", "spec", "containers"}}
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{
								{
									ID: nil,
									RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
										MutatedRef: &refs,
									},
								},
							},
							To: []relationship.SelectorItem{{ID: &toID}},
						},
					},
				}
				return relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
			}(),
			wantOk: false,
		},
		{
			name: "valid alias relationship returns true",
			input: func() relationship.RelationshipDefinition {
				refs := [][]string{{"configuration", "spec", "containers"}}
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{
								{
									ID: &fromID,
									RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
										MutatedRef: &refs,
									},
								},
							},
							To: []relationship.SelectorItem{{ID: &toID}},
						},
					},
				}
				rd := relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
				rd.ID = relID
				return rd
			}(),
			wantOk: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alias, ok := parseRelationshipToAlias(tt.input)
			assert.Equal(t, tt.wantOk, ok, "parseRelationshipToAlias() ok mismatch")

			if tt.wantOk {
				assert.Equal(t, toID, alias.ImmediateParentId, "ImmediateParentId should match to.ID")
				assert.Equal(t, fromID, alias.AliasComponentId, "AliasComponentId should match from.ID")
				assert.Equal(t, relID, alias.RelationshipId, "RelationshipId should match the relationship's Id")
				assert.Equal(t, []string{"configuration", "spec", "containers"}, alias.ImmediateRefFieldPath, "ImmediateRefFieldPath should match first mutatedRef entry")
			}
		})
	}
}
