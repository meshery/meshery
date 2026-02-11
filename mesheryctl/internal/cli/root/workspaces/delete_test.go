package workspaces

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestDeleteWorkspace(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	workspaceID := "d56fb25b-f92c-4cd6-821b-2cfd6bb87259"

	// Test scenarios for environment deletion
	tests := []utils.MesheryCommandTest{
		{
			Name:             "given no workspace-id is provided when running mesheryctl exp workspace delete then an error message is displayed",
			Args:             []string{"delete"},
			HttpMethod:       "DELETE",
			HttpStatusCode:   200,
			URL:              workspacesApiPath,
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("[ Workspace ID ] not specified")),
		},
		{
			Name:             "given an invalid workspace-id is provided when running mesheryctl exp workspace delete invalid-id then an error message is displayed",
			Args:             []string{"delete", "invalid"},
			HttpMethod:       "DELETE",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s/%s", workspacesApiPath, "invalid"),
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    utils.ErrInvalidUUID(fmt.Errorf("invalid workspace ID: %q", "invalid")),
		},
		{
			Name:             "given a valid workspace-id is provided when running mesheryctl exp workspace delete valid-id then the specified workspace gets deleted",
			Args:             []string{"delete", workspaceID},
			HttpMethod:       "DELETE",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s/%s", workspacesApiPath, workspaceID),
			Fixture:          "delete.workspace.response.golden",
			ExpectedResponse: "delete.workspace.success.golden",
			ExpectError:      false,
		},
	}
	fmt.Printf("%s/%s", workspacesApiPath, workspaceID)

	utils.InvokeMesheryctlTestCommand(t, update, WorkSpaceCmd, tests, currDir, "workspace")
}
