package workspaces

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestListWorkspaces(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "given no organization ID when workspace list then throw error",
			Args:             []string{"list"},
			URL:              "",
			Fixture:          "list.workspace.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --orgId ''")),
		},
		{
			Name:             "given organization ID when workspace list then return workspaces under that organization",
			Args:             []string{"list", "--orgId", testOrgId},
			URL:              fmt.Sprintf("/%s?orgId=%s&page=0&pagesize=10", workspacesApiPath, testOrgId),
			Fixture:          "list.workspace.api.response.golden",
			ExpectedResponse: "list.workspace.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given organization ID without workspaces when workspace list then return empty result",
			Args:             []string{"list", "--orgId", testOrgId},
			URL:              fmt.Sprintf("/%s?orgId=%s&page=0&pagesize=10", workspacesApiPath, testOrgId),
			Fixture:          "list.workspace.empty.api.response.golden",
			ExpectedResponse: "list.workspace.empty.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given extra positional arguments when workspace list then throw error",
			Args:             []string{"list", "extra-arg", "--orgId", testOrgId},
			URL:              "",
			Fixture:          "list.workspace.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\nUsage: mesheryctl workspace list --orgId [orgId]\nRun 'mesheryctl workspace list --help' to see detailed help message")),
		},
	}

	mesheryctlflags.InitValidators(WorkSpaceCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, WorkSpaceCmd, tests, currentDirectory, "organization")
}
