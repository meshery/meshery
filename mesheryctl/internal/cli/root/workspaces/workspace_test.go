package workspaces

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
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
			Name:             "Display error when no args and flag provided",
			Args:             []string{},
			URL:              "",
			Fixture:          "list.workspace.empty.api.response.golden",
			ExpectedResponse: "workspace.no.args.no.flags.output.golden",
			ExpectError:      true,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, WorkSpaceCmd, tests, currentDirectory, "organization")
}
