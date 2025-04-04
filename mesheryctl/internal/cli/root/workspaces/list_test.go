package workspaces

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestListWorkspaces(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "List workspaces whithout providing organization ID",
			Args:             []string{"list"},
			URL:              "",
			Fixture:          "list.workspace.api.response.golden",
			ExpectedResponse: "list.workspace.no.orgID.output.golden",
			ExpectError:      true,
		},
		{
			Name:             "List workspaces providing organization ID",
			Args:             []string{"list", "--orgId", testOrgId},
			URL:              fmt.Sprintf("/%s?orgID=%s", workspacesApiPath, testOrgId),
			Fixture:          "list.workspace.api.response.golden",
			ExpectedResponse: "list.workspace.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, WorkSpaceCmd, tests, currentDirectory, "organization")
}
