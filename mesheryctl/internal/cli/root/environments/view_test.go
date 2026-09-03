package environments

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// TestEnvironmentViewNoSaveWithBrokenHome verifies that "environment view --orgId"
// without --save succeeds even when HOME/USERPROFILE is unset, proving
// os.UserHomeDir() is not called on the non-save path.
func TestEnvironmentViewNoSaveWithBrokenHome(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(EnvironmentCmd, t)

	utils.TokenFlag = utils.GetToken(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	apiResponse := utils.NewGoldenFile(t, "view.environment.api.response.golden", fixturesDir).Load()
	httpmock.RegisterResponder("GET",
		testContext.BaseURL+"/api/environments?orgId="+testConstants["orgId"],
		httpmock.NewStringResponder(200, apiResponse))

	// Unset HOME and USERPROFILE to simulate unavailable home directory
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")

	buf := &bytes.Buffer{}
	EnvironmentCmd.SetOut(buf)
	EnvironmentCmd.SetErr(buf)

	_ = utils.SetupMeshkitLoggerTesting(t, false)
	EnvironmentCmd.SetArgs([]string{"view", "--orgId", testConstants["orgId"]})
	if execErr := EnvironmentCmd.Execute(); execErr != nil {
		t.Fatalf("view without --save should succeed even with no HOME: %v", execErr)
	}
}

// TestConnectionViewNoSaveWithBrokenHome equivalent for environments --save path:
// verifies saved file uses correct extension and isolated temp home.
func TestEnvironmentViewSaveCreatesFile(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(EnvironmentCmd, t)

	utils.TokenFlag = utils.GetToken(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	apiResponse := utils.NewGoldenFile(t, "view.environment.api.response.golden", fixturesDir).Load()
	httpmock.RegisterResponder("GET",
		testContext.BaseURL+"/api/environments?orgId="+testConstants["orgId"],
		httpmock.NewStringResponder(200, apiResponse))

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	mesheryDir := filepath.Join(tmpHome, ".meshery")
	if err := os.MkdirAll(mesheryDir, 0755); err != nil {
		t.Fatalf("cannot create %s: %v", mesheryDir, err)
	}
	expectedFile := filepath.Join(mesheryDir, "environment_test-environment.yaml")

	buf := &bytes.Buffer{}
	EnvironmentCmd.SetOut(buf)
	EnvironmentCmd.SetErr(buf)

	_ = utils.SetupMeshkitLoggerTesting(t, false)
	EnvironmentCmd.SetArgs([]string{"view", "--orgId", testConstants["orgId"], "--save"})
	if execErr := EnvironmentCmd.Execute(); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}

	if _, statErr := os.Stat(expectedFile); os.IsNotExist(statErr) {
		entries, _ := os.ReadDir(mesheryDir)
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("--save: expected file %q to exist, got: %v", expectedFile, names)
	}
}