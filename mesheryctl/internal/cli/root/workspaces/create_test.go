package workspaces

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestCreateWorkspace(t *testing.T) {
	// Get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// Test scenarios for workspace creation
	tests := []utils.MesheryCommamdTest{
		{
			Name:             "Create workspace without arguments",
			Args:             []string{"create"},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			Fixture:          "",
			ExpectedResponse: "create.workspace.missing.flag.output.golden",
			ExpectError:      true,
		},
		{
			Name:             "Create workspace successfully",
			Args:             []string{"create", "-n", "workspace-test", "-d", "integration test", "--orgId", testOrgId},
			URL:              "/api/workspaces",
			HttpMethod:       "POST",
			HttpStatusCode:   201,
			Fixture:          "create.workspace.api.response.golden",
			ExpectedResponse: "create.workspace.success.output.golden",
			ExpectError:      false,
		},
	}

	// Run tests
	utils.InvokeMesheryctlTestCommand(t, update, WorkSpaceCmd, tests, currDir, "workspaces")
}
