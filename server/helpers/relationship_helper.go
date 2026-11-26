package helpers

import (
	"strings"

	"github.com/meshery/schemas/models/v1alpha3/relationship"
)

// NormalizeRelationshipDefinition canonicalizes relationship classification fields
// (Kind, RelationshipType, SubType) to lowercase for consistent lookups and queries.
func NormalizeRelationshipDefinition(r *relationship.RelationshipDefinition) {
	if r == nil {
		return
	}
	r.Kind = relationship.RelationshipDefinitionKind(strings.ToLower(string(r.Kind)))
	r.RelationshipType = strings.ToLower(r.RelationshipType)
	r.SubType = strings.ToLower(r.SubType)
}
