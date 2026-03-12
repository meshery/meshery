package relationships

import (
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
	meshkitschema "github.com/meshery/meshkit/schema"
	"github.com/stretchr/testify/assert"
)

func TestRelationshipValidationErrorFormatting(t *testing.T) {
	err := ErrRelationshipValidationFailedForFile("/tmp/relationship.yaml", meshkitschema.ValidationDetails{
		Violations: []meshkitschema.Violation{
			{
				InstancePath: "/kind",
				Message:      `value is not one of the allowed values ["hierarchical","edge","sibling"]`,
			},
			{
				Message: "schema does not match any allowed shape",
			},
		},
	})

	assert.Equal(t, "Relationship validation failed: 2 schema violations found in relationship.yaml\n\nViolations:\n  1. /kind: value is not one of the allowed values [\"hierarchical\",\"edge\",\"sibling\"]\n  2. document: schema does not match any allowed shape", err.Error())
	assert.Equal(t, ErrRelationshipValidationFailedCode, meshkiterrors.GetCode(err))
}
