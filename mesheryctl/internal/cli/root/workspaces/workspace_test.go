package workspaces

import (
	"flag"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")
var testOrgId = "2d2c0b60-076a-4f0a-8a63-de538570a553"

func TestWorkspaces(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "Given no argument provided trigger an error",
			Args:             []string{},
			URL:              "",
			Fixture:          "list.workspace.empty.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    returnInvalidArgumentProvidedError(),
		},
		{
			Name:             "Given invalid subcommand name provided trigger an error",
			Args:             []string{"invalidCommand"},
			URL:              "",
			Fixture:          "list.workspace.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("'%s' is an invalid subcommand. Please provide required options from [create, list]. Use 'mesheryctl exp workspace --help' to display usage guide", "invalidCommand")),
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, WorkSpaceCmd, tests, currentDirectory, "workspace")
}
