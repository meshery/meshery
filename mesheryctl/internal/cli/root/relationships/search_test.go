package relationships

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func expectedSearchMissingFlagsError() error {
	return utils.ErrFlagsInvalid(fmt.Errorf(
		"at least one of [--kind, --subtype, --type, --model] is required\n\n" +
			"Usage: mesheryctl exp relationship search [--kind <kind>] [--type <type>] [--subtype <subtype>] [--model <model>]\n" +
			"Run 'mesheryctl exp relationship search --help'",
	))
}

func TestSearch_WithoutFlags(t *testing.T) {
	mesheryctlflags.InitValidators(RelationshipCmd)

	// setup current context
	utils.SetupContextEnv(t)

	// test scenarios for fetching data
	tests := []struct {
		Name          string
		Args          []string
		ExpectError   bool
		ExpectedError error
	}{
		{
			Name:          "Search with missing arguments",
			Args:          []string{"search"},
			ExpectError:   true,
			ExpectedError: expectedSearchMissingFlagsError(),
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mesheryctlflags.InitValidators(RelationshipCmd)
			RelationshipCmd.SetArgs(tt.Args)
			err := RelationshipCmd.Execute()

			if err != nil {
				if tt.ExpectError {
					utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					return
				}
				t.Fatal(err)
			}

			if tt.ExpectError {
				t.Fatalf("expected an error but command succeeded")
			}
		})
		t.Log("Search experimental relationship test passed")
	}

	utils.StopMockery(t)
}

func TestSearch_WithFlags(t *testing.T) {
	mesheryctlflags.InitValidators(RelationshipCmd)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "Search registered relationships matching result(s) found",
			Args:             []string{"search", "--model", "kubernetes"},
			URL:              "/api/meshmodels/models/kubernetes/relationships",
			Fixture:          "search.relationship.api.response.matching.result.golden",
			ExpectedResponse: "search.relationship.output.matching.result.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "Search registered relationships no matching result(s) found",
			Args:             []string{"search", "--model", "kubernetes"},
			URL:              "/api/meshmodels/models/kubernetes/relationships",
			Fixture:          "search.relationship.api.response.no.matching.result.golden",
			ExpectedResponse: "search.relationship.output.no.matching.result.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
	}

	mesheryctlflags.InitValidators(RelationshipCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, RelationshipCmd, tests, currDir, "relationships")
}
