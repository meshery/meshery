package environments

import (
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)


func environmentViewRun(commandDir string, args []string) utils.MesheryCommandRun {
	return utils.MesheryCommandRun{
		Cmd:        EnvironmentCmd,
		Args:       args,
		CommandDir: commandDir,
		Mocks: []utils.MockURL{
			{
				URL:      "/api/environments?orgId=" + testConstants["orgId"],
				Response: "view.environment.api.response.golden",
			},
		},
	}
}


func TestEnvironmentViewNoSaveWithBrokenHome(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")

	run := environmentViewRun(utils.CallerDir(t), []string{"view", "--orgId", testConstants["orgId"]})
	if _, err := utils.RunMesheryctlCommand(t, run); err != nil {
		t.Fatalf("view without --save should succeed even with no HOME: %v", err)
	}
}


func TestEnvironmentViewSaveCreatesFile(t *testing.T) {
	mesheryDir := utils.IsolateMesheryHome(t)

	run := environmentViewRun(utils.CallerDir(t), []string{"view", "--orgId", testConstants["orgId"], "--save"})
	if _, err := utils.RunMesheryctlCommand(t, run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	utils.AssertFileExists(t, filepath.Join(mesheryDir, "environment_test-environment.yaml"))
}