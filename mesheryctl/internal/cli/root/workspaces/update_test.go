package workspaces

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestUpdateWorkspace(t *testing.T) {
	// Get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	workspaceID := "d56fb25b-f92c-4cd6-821b-2cfd6bb87259"

	tests := []utils.MesheryCommandTest{
		{
			Name:             "given no workspace ID when workspace update then throw error",
			Args:             []string{"update", "--orgId", testOrgId, "-n", validWorkspaceName},
			URL:              fmt.Sprintf("/api/workspaces/%s", workspaceID),
			HttpMethod:       "PUT",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("please provide exactly one workspace ID\n\n%v",
				"Usage: mesheryctl workspace update [workspace-id] --orgId [orgId] [--name NAME] [--description DESCRIPTION]\nRun 'mesheryctl workspace update --help' to see detailed help message")),
		},
		{
			Name:             "given an invalid workspace ID when workspace update then throw error",
			Args:             []string{"update", "not-a-uuid", "--orgId", testOrgId, "-n", validWorkspaceName},
			URL:              fmt.Sprintf("/api/workspaces/%s", workspaceID),
			HttpMethod:       "PUT",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: utils.ErrInvalidUUID(fmt.Errorf("invalid workspace ID: %s\n\n%v", "not-a-uuid",
				"Usage: mesheryctl workspace update [workspace-id] --orgId [orgId] [--name NAME] [--description DESCRIPTION]\nRun 'mesheryctl workspace update --help' to see detailed help message")),
		},
		{
			Name:             "given no --name or --description when workspace update then throw error",
			Args:             []string{"update", workspaceID, "--orgId", testOrgId},
			URL:              fmt.Sprintf("/api/workspaces/%s", workspaceID),
			HttpMethod:       "PUT",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("at least one of --name or --description must be provided\n\n%v",
				"Usage: mesheryctl workspace update [workspace-id] --orgId [orgId] [--name NAME] [--description DESCRIPTION]\nRun 'mesheryctl workspace update --help' to see detailed help message")),
		},
		{
			Name:             "given missing orgId when workspace update then throw error",
			Args:             []string{"update", workspaceID, "-n", validWorkspaceName},
			URL:              fmt.Sprintf("/api/workspaces/%s", workspaceID),
			HttpMethod:       "PUT",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --orgId ''")),
		},
		{
			Name:             "given a valid workspace ID, orgId and name when workspace update then workspace is updated",
			Args:             []string{"update", workspaceID, "--orgId", testOrgId, "-n", validWorkspaceName},
			URL:              fmt.Sprintf("/api/workspaces/%s", workspaceID),
			HttpMethod:       "PUT",
			HttpStatusCode:   200,
			Fixture:          "update.workspace.api.response.golden",
			ExpectedResponse: "update.workspace.success.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given a workspace ID that no longer exists when workspace update then it is reported as not found",
			Args:             []string{"update", workspaceID, "--orgId", testOrgId, "-d", workspaceDescription},
			URL:              fmt.Sprintf("/api/workspaces/%s", workspaceID),
			HttpMethod:       "PUT",
			HttpStatusCode:   404,
			Fixture:          "update.workspace.api.nil.response.golden",
			ExpectedResponse: "update.workspace.notfound.output.golden",
			ExpectError:      false,
		},
	}

	mesheryctlflags.InitValidators(WorkSpaceCmd)
	utils.InvokeMesheryctlTestCommand(t, update, WorkSpaceCmd, tests, currDir, "workspaces")
}
