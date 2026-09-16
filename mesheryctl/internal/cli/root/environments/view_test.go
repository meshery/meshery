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

// runEnvironmentViewTest handles all shared scaffolding for environment view tests.
func runEnvironmentViewTest(t *testing.T, args []string) error {
	t.Helper()
	testContext := utils.InitTestEnvironment(t)
	t.Cleanup(func() { utils.StopMockery(t) })
	t.Cleanup(func() { utils.ResetCommandFlags(EnvironmentCmd, t) })

	utils.TokenFlag = utils.GetToken(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine current working directory")
	}
	fixturesDir := filepath.Join(filepath.Dir(filename), "fixtures")

	apiResponse := utils.NewGoldenFile(t, "view.environment.api.response.golden", fixturesDir).Load()
	httpmock.RegisterResponder("GET",
		testContext.BaseURL+"/api/environments?orgId="+testConstants["orgId"],
		httpmock.NewStringResponder(200, apiResponse))

	buf := &bytes.Buffer{}
	EnvironmentCmd.SetOut(buf)
	EnvironmentCmd.SetErr(buf)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	EnvironmentCmd.SetArgs(args)
	return EnvironmentCmd.Execute()
}

// TestEnvironmentViewNoSaveWithBrokenHome verifies that "environment view --orgId"
// without --save succeeds even when HOME/USERPROFILE is unset, proving
// os.UserHomeDir() is not called on the non-save path.
func TestEnvironmentViewNoSaveWithBrokenHome(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")
	if err := runEnvironmentViewTest(t, []string{"view", "--orgId", testConstants["orgId"]}); err != nil {
		t.Fatalf("view without --save should succeed even with no HOME: %v", err)
	}
}

// TestEnvironmentViewSaveCreatesFile verifies saved file uses correct
// extension and is written into an isolated temp home directory.
func TestEnvironmentViewSaveCreatesFile(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	mesheryDir := filepath.Join(tmpHome, ".meshery")
	if err := os.MkdirAll(mesheryDir, 0755); err != nil {
		t.Fatalf("cannot create %s: %v", mesheryDir, err)
	}

	if err := runEnvironmentViewTest(t, []string{"view", "--orgId", testConstants["orgId"], "--save"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFile := filepath.Join(mesheryDir, "environment_test-environment.yaml")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		entries, _ := os.ReadDir(mesheryDir)
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("--save: expected file %q to exist, got: %v", expectedFile, names)
	}
}