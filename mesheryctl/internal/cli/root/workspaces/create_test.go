package workspaces

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestCreateWorkspace(t *testing.T) {
	// Get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	tests := []utils.MesheryCommandTest{
		{
			Name:             "given no flags provided when workspace create throw error",
			Args:             []string{"create"},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --orgId '', Invalid value for --name ''")),
		},
		{
			Name:             "given missing flag orgId when workspace create throw error",
			Args:             []string{"create", "-n", validWorkspaceName, "-d", workspaceDescription},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --orgId ''")),
		},
		{
			Name:             "given missing flag name when workspace create throw error",
			Args:             []string{"create", "--orgId", testOrgId, "-d", workspaceDescription},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --name ''")),
		},
		{
			Name:             "given missing flag description when workspace create then workspace is created with empty description",
			Args:             []string{"create", "--orgId", testOrgId, "-n", validWorkspaceName},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			HttpStatusCode:   201,
			Fixture:          "create.workspace.api.response.golden",
			ExpectedResponse: "create.workspace.success.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given missing name and orgId flags when workspace create then throw error",
			Args:             []string{"create", "-d", workspaceDescription},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			HttpStatusCode:   404,
			Fixture:          "create.workspace.api.nil.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --orgId '', Invalid value for --name ''")),
		},
		{
			Name:             "given an invalid organization Id when workspace create throw error",
			Args:             []string{"create", "-n", validWorkspaceName, "-d", workspaceDescription, "--orgId", testOrgId},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			HttpStatusCode:   404,
			Fixture:          "create.workspace.api.nil.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    returnFailedCreateWorkspaceError(validWorkspaceName, testOrgId),
		},
		{
			Name:             "given all flags provided when workspace create then workspace is created",
			Args:             []string{"create", "-n", validWorkspaceName, "-d", workspaceDescription, "--orgId", testOrgId},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			HttpStatusCode:   201,
			Fixture:          "create.workspace.api.response.golden",
			ExpectedResponse: "create.workspace.success.output.golden",
			ExpectError:      false,
		},
	}

	mesheryctlflags.InitValidators(WorkSpaceCmd)
	utils.InvokeMesheryctlTestCommand(t, update, WorkSpaceCmd, tests, currDir, "workspaces")
}
