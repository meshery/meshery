package relationships

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestRelationship(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "Display error without any flags or args",
			Args:             []string{},
			URL:              "",
			Fixture:          "list.relationship.api.empty.response.golden",
			ExpectedContains: []string{"no command specified."},
			ExpectError:      true,
		},
		{
			Name:             "Display error given invalid command",
			Args:             []string{"invalidCommand"},
			URL:              "",
			Fixture:          "list.relationship.api.empty.response.golden",
			ExpectedContains: []string{"is an invalid subcommand"},
			ExpectError:      true,
		},
		{
			Name:             "Display count of registered relationships empty",
			Args:             []string{"--count"},
			URL:              "/api/meshmodels/relationships",
			Fixture:          "list.relationship.api.empty.response.golden",
			ExpectedContains: []string{"No relationships found"},
			ExpectError:      false,
		},
		{
			Name:             "Display count of registered relationships",
			Args:             []string{"--count"},
			URL:              "/api/meshmodels/relationships",
			Fixture:          "list.relationship.api.response.golden",
			ExpectedContains: []string{"Total number of relationships: 1"},
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, RelationshipCmd, tests, currDir, "relationship")
}
