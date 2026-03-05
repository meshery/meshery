package workspaces

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestViewWorkspace(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// --- INJECT THE FLAG VALIDATOR CONTEXT FOR THE TEST ---
	ctx := context.WithValue(context.Background(), mesheryctlflags.FlagValidatorKey, mesheryctlflags.NewFlagValidator())
	WorkSpaceCmd.SetContext(ctx)
	// ------------------------------------------------------

	testWorkspaceID := "0dd47d1a-d1c9-47dc-897c-40bf4a71d96b"

	tests := []utils.MesheryCommandTest{
		{
			Name:             "View workspace without arguments",
			Args:             []string{"view"},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s", workspacesApiPath),
			Fixture:          "view.workspace.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("please provide exactly one workspace name or ID\n\nUsage: mesheryctl exp workspace view [workspace-name|workspace-id]\nRun 'mesheryctl exp workspace view --help' to see detailed help message")),
		},
		{
			Name:             "View workspace by ID",
			Args:             []string{"view", testWorkspaceID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s/%s", workspacesApiPath, testWorkspaceID),
			Fixture:          "view.workspace.api.response.golden",
			ExpectedResponse: "view.workspace.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, WorkSpaceCmd, tests, currDir, "workspace")
}
