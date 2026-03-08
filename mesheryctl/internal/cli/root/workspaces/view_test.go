package workspaces

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func expectedViewFlagError(outputFormat string) error {
	fv := mesheryctlflags.NewFlagValidator()
	return fv.Validate(&workspaceViewFlags{OutputFormat: outputFormat})
}

func TestViewWorkspace(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	mesheryctlflags.InitValidators(WorkSpaceCmd)

	testWorkspaceID := "0dd47d1a-d1c9-47dc-897c-40bf4a71d96b"

	tests := []utils.MesheryCommandTest{
		{
			Name:             "given no arguments when running workspace view then return error",
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
			Name:             "given too many arguments when running workspace view then return error",
			Args:             []string{"view", testWorkspaceID, "extra-arg"},
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
			Name:             "given valid workspace ID when running workspace view then display workspace details",
			Args:             []string{"view", testWorkspaceID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s/%s", workspacesApiPath, testWorkspaceID),
			Fixture:          "view.workspace.api.response.golden",
			ExpectedResponse: "view.workspace.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given workspace name without orgId when running workspace view then return error",
			Args:             []string{"view", "my-workspace"},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s", workspacesApiPath),
			Fixture:          "view.workspace.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("--orgId is required when searching by name\n\nUsage: mesheryctl exp workspace view [workspace-name] --orgId [orgId]")),
		},
		{
			Name:             "given invalid output format flag when running workspace view then return error",
			Args:             []string{"view", testWorkspaceID, "--output-format", "invalid"},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s/%s", workspacesApiPath, testWorkspaceID),
			Fixture:          "view.workspace.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    expectedViewFlagError("invalid"),
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, WorkSpaceCmd, tests, currDir, "workspace")
}