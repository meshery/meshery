package relationships

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestList(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "List registered relationships empty",
			Args:             []string{"list"},
			URL:              "/api/meshmodels/relationships",
			Fixture:          "list.relationship.api.empty.response.golden",
			ExpectedResponse: "list.relationship.empty.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "List registered relationships",
			Args:             []string{"list"},
			URL:              "/api/meshmodels/relationships",
			Fixture:          "list.relationship.api.response.golden",
			ExpectedResponse: "list.relationship.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Display count of registered relationships empty result",
			Args:             []string{"list", "--count"},
			URL:              "/api/meshmodels/relationships",
			Fixture:          "list.relationship.api.empty.response.golden",
			ExpectedResponse: "list.count.relationship.empty.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Display count of registered relationships",
			Args:             []string{"list", "--count"},
			URL:              "/api/meshmodels/relationships",
			Fixture:          "list.relationship.api.response.golden",
			ExpectedResponse: "list.count.relationship.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, RelationshipCmd, tests, currDir, "relationships")
}
