package relationships

import (
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitschema "github.com/meshery/meshkit/schema"
)

func TestValidate(t *testing.T) {
	mesheryctlflags.InitValidators(RelationshipCmd)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	validRelationshipPath := filepath.Join(currDir, "testdata", "validate.relationship.valid.yaml")
	invalidRelationshipPath := filepath.Join(currDir, "testdata", "validate.relationship.invalid.yaml")

	tests := []utils.MesheryCommandTest{
		{
			Name:             "given a valid relationship file when running relationship validate then validation succeeds",
			Args:             []string{"validate", "--file", validRelationshipPath},
			ExpectedResponse: "validate.relationship.valid.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:        "given an invalid relationship file when running relationship validate then validation fails",
			Args:        []string{"validate", "--file", invalidRelationshipPath},
			ExpectError: true,
			ExpectedError: ErrRelationshipValidationFailed(meshkitschema.ValidationDetails{
				Violations: []meshkitschema.Violation{{InstancePath: "/kind", Message: "invalid value"}},
			}),
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, RelationshipCmd, tests, currDir, "relationships")
}
