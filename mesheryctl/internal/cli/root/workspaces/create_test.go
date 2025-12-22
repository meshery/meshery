package workspaces

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestCreateWorkspace(t *testing.T) {
	// Get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	tests := []utils.MesheryCommamdTest{
		{
			Name:             "Given no flags provided trigger an error",
			Args:             []string{"create"},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf(createMissingArgumentsErrorMessage, strings.Join([]string{"orgId", "name", "description"}, " | "))),
		},
		{
			Name:             "Given missing flag orgId trigger an error",
			Args:             []string{"create", "-n", "workspace-test", "-d", "integration test"},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf(createMissingArgumentsErrorMessage, strings.Join([]string{"orgId"}, " | "))),
		},
		{
			Name:             "Given missing flag name trigger an error",
			Args:             []string{"create", "--orgId", testOrgId, "-d", "integration test"},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf(createMissingArgumentsErrorMessage, strings.Join([]string{"name"}, " | "))),
		},
		{
			Name:             "Given missing flag description trigger an error",
			Args:             []string{"create", "--orgId", testOrgId, "-n", "workspace-test"},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf(createMissingArgumentsErrorMessage, strings.Join([]string{"description"}, " | "))),
		},
		{
			Name:             "Given multiple missing flags trigger an error",
			Args:             []string{"create", "-d", "integration test"},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			HttpStatusCode:   404,
			Fixture:          "create.workspace.api.nil.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf(createMissingArgumentsErrorMessage, strings.Join([]string{"orgId", "name"}, " | "))),
		},
		{
			Name:             "Given an invalid organization Id trigger an error",
			Args:             []string{"create", "-n", "workspace-test-error", "-d", "integration test", "--orgId", testOrgId},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			HttpStatusCode:   404,
			Fixture:          "create.workspace.api.nil.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    returnFailedCreateWorkspaceError("workspace-test-error", testOrgId),
		},
		{
			Name:             "Given all requirements met, create workspace successfully",
			Args:             []string{"create", "-n", "workspace-test", "-d", "integration test", "--orgId", testOrgId},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			HttpStatusCode:   201,
			Fixture:          "create.workspace.api.response.golden",
			ExpectedResponse: "create.workspace.success.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, WorkSpaceCmd, tests, currDir, "workspaces")
}
