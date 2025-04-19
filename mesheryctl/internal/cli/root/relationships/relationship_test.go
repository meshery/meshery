package relationships

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")

func TestRelationship(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "Display error without any flags or args",
			Args:             []string{},
			URL:              "",
			Fixture:          "list.relationship.api.empty.response.golden",
			ExpectedResponse: "relationship.no.args.no.flag.output.golden",
			ExpectError:      true,
		},
		{
			Name:             "Display error given invalid command",
			Args:             []string{"invalidCommand"},
			URL:              "",
			Fixture:          "list.relationship.api.empty.response.golden",
			ExpectedResponse: "relationship.invalid.command.output.golden",
			ExpectError:      true,
		},
		{
			Name:             "Display count of registered relationships empty",
			Args:             []string{"--count"},
			URL:              "/api/meshmodels/relationships",
			Fixture:          "list.relationship.api.empty.response.golden",
			ExpectedResponse: "list.count.relationship.empty.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Display count of registered relationships",
			Args:             []string{"--count"},
			URL:              "/api/meshmodels/relationships",
			Fixture:          "list.relationship.api.response.golden",
			ExpectedResponse: "list.count.relationship.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, RelationshipCmd, tests, currDir, "relationship")
}
