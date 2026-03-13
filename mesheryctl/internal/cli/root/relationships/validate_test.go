package relationships

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func expectedValidateMissingFileError() error {
	return utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --file ''"))
}

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
			Name:             "given an invalid relationship file when running relationship validate then validation failures are surfaced",
			Args:             []string{"validate", "--file", invalidRelationshipPath},
			ExpectedResponse: "validate.relationship.invalid.output.golden",
			ExpectError:      true,
			IsOutputGolden:   true,
		},
		{
			Name:          "given no file flag when running relationship validate then shared flag validation fails",
			Args:          []string{"validate"},
			ExpectError:   true,
			ExpectedError: expectedValidateMissingFileError(),
		},
		{
			Name:        "given an unexpected positional argument when running relationship validate then argument validation fails",
			Args:        []string{"validate", "unexpected", "--file", validRelationshipPath},
			ExpectError: true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf(
				"too many arguments specified\n\n%s",
				relationshipValidateUsage,
			)),
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, RelationshipCmd, tests, currDir, "relationships")
}
