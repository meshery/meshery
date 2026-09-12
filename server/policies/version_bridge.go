//go:build !js

package policies

import (
	patternutils "github.com/meshery/meshery/server/models/pattern/utils"
	relationshipv1alpha3 "github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

func bridgeRelationship(r any) *relationship.RelationshipDefinition {
	switch v := r.(type) {
	case *relationshipv1alpha3.RelationshipDefinition:
		return patternutils.RelationshipV1alpha3ToV1beta2(v)
	case relationshipv1alpha3.RelationshipDefinition:
		return patternutils.RelationshipV1alpha3ToV1beta2(&v)
	}
	return nil
}
