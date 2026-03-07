package handlers

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
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
							To:   []relationship.SelectorItem{{Id: &toID}},
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
							From: []relationship.SelectorItem{{Id: &fromID}},
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
							From: []relationship.SelectorItem{{Id: &fromID, Patch: nil}},
							To:   []relationship.SelectorItem{{Id: &toID}},
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
									Id: &fromID,
									Patch: &relationship.RelationshipDefinition_Selectors_Patch{
										MutatedRef: nil,
									},
								},
							},
							To: []relationship.SelectorItem{{Id: &toID}},
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
									Id: &fromID,
									Patch: &relationship.RelationshipDefinition_Selectors_Patch{
										MutatedRef: &emptyRefs,
									},
								},
							},
							To: []relationship.SelectorItem{{Id: &toID}},
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
			name: "nil to.Id returns false",
			input: func() relationship.RelationshipDefinition {
				refs := [][]string{{"configuration", "spec", "containers"}}
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{
								{
									Id: &fromID,
									Patch: &relationship.RelationshipDefinition_Selectors_Patch{
										MutatedRef: &refs,
									},
								},
							},
							To: []relationship.SelectorItem{{Id: nil}},
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
			name: "nil from.Id returns false",
			input: func() relationship.RelationshipDefinition {
				refs := [][]string{{"configuration", "spec", "containers"}}
				ss := relationship.SelectorSet{
					{
						Allow: relationship.Selector{
							From: []relationship.SelectorItem{
								{
									Id: nil,
									Patch: &relationship.RelationshipDefinition_Selectors_Patch{
										MutatedRef: &refs,
									},
								},
							},
							To: []relationship.SelectorItem{{Id: &toID}},
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
									Id: &fromID,
									Patch: &relationship.RelationshipDefinition_Selectors_Patch{
										MutatedRef: &refs,
									},
								},
							},
							To: []relationship.SelectorItem{{Id: &toID}},
						},
					},
				}
				rd := relationship.RelationshipDefinition{
					SubType:   RELATIONSHIP_SUBTYPE_ALIAS,
					Selectors: &ss,
				}
				rd.Id = relID
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
				assert.Equal(t, toID, alias.ImmediateParentId, "ImmediateParentId should match to.Id")
				assert.Equal(t, fromID, alias.AliasComponentId, "AliasComponentId should match from.Id")
				assert.Equal(t, relID, alias.RelationshipId, "RelationshipId should match the relationship's Id")
				assert.Equal(t, []string{"configuration", "spec", "containers"}, alias.ImmediateRefFieldPath, "ImmediateRefFieldPath should match first mutatedRef entry")
			}
		})
	}
}
