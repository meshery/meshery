package relationships

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestView(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "given no model name provided when running relationship view then throw error",
			Args:             []string{"view"},
			URL:              "/api/meshmodels/models/kubernetes/relationships?pagesize=all",
			Fixture:          "view.relationship.empty.response.golden",
			ExpectedResponse: "",
			IsOutputGolden:   false,
			ExpectError:      true,
			ExpectedError:    utils.ErrInvalidArgument(errNoModelNameProvided),
		},
		{
			Name:             "given model name provided when running relationship view then display registered relationship",
			Args:             []string{"view", "kubernetes"},
			URL:              "/api/meshmodels/models/kubernetes/relationships?pagesize=all",
			Fixture:          "view.relationship.api.response.golden",
			ExpectedResponse: "view.relationship.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given non existing model name provided when running relationship view then display no relationship found",
			Args:             []string{"view", "nonexistent"},
			URL:              "/api/meshmodels/models/nonexistent/relationships?pagesize=all",
			Fixture:          "view.relationship.empty.response.golden",
			ExpectedResponse: "view.relationship.non.existent.output.golden",
			ExpectError:      false,
		},
	}

	// Run tests
	utils.InvokeMesheryctlTestListCommand(t, update, RelationshipCmd, tests, currDir, "relationships")
}
