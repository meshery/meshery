package filter

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestViewCmd(t *testing.T) {
	mesheryctlflags.InitValidators(FilterCmd)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	ListTests := []utils.MesheryListCommandTest{
		{
			Name:             "Fetch Filter View",
			Args:             []string{"view", "KumaTest"},
			ExpectedResponse: "view.filter.output.golden",
			Fixture:          "view.filter.api.response.golden",
			URL:              "/api/filter",
			ExpectError:      false,
		},
		{
			Name:             "Fetch Kuma Filter View with ID",
			Args:             []string{"view", "957fbc9b-a655-4892-823d-375102a9587c"},
			ExpectedResponse: "view.id.filter.output.golden",
			Fixture:          "view.id.filter.api.response.golden",
			URL:              "/api/filter/957fbc9b-a655-4892-823d-375102a9587c",
			ExpectError:      false,
		},
	}

	loggerTests := []utils.MesheryCommandTest{
		{
			Name:             "Fetch Filter View for non existing filter",
			Args:             []string{"view", "xyz"},
			ExpectedResponse: "view.nonexisting.filter.output.golden",
			Fixture:          "view.nonexisting.filter.api.response.golden",
			URL:              "/api/filter",
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, FilterCmd, ListTests, currDir, "filter")
	utils.InvokeMesheryctlTestCommand(t, update, FilterCmd, loggerTests, currDir, "filter")
}

func TestFilterViewSaveCreatesFileWithExtension(t *testing.T) {
	mesheryctlflags.InitValidators(FilterCmd)

	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(FilterCmd, t)

	utils.TokenFlag = utils.GetToken(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	const filterID = "957fbc9b-a655-4892-823d-375102a9587c"
	apiResponse := utils.NewGoldenFile(t, "view.id.filter.api.response.golden", fixturesDir).Load()
	listResponse := utils.NewGoldenFile(t, "view.filter.api.response.golden", fixturesDir).Load()

	httpmock.RegisterResponder("GET",
		testContext.BaseURL+"/api/filter",
		httpmock.NewStringResponder(200, listResponse))
	httpmock.RegisterResponder("GET",
		testContext.BaseURL+"/api/filter/"+filterID,
		httpmock.NewStringResponder(200, apiResponse))

	tmpDir := t.TempDir()
	origMesheryFolder := utils.MesheryFolder
	utils.MesheryFolder = tmpDir
	t.Cleanup(func() { utils.MesheryFolder = origMesheryFolder })

	expectedFile := filepath.Join(tmpDir, "filter_KumaTest_957fbc9b.json")

	origStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		_ = w.Close()
		os.Stdout = origStdout
	}()

	_ = utils.SetupMeshkitLoggerTesting(t, false)
	FilterCmd.SetArgs([]string{"view", filterID, "--output-format", "json", "--save"})
	if err := FilterCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_ = w.Close()
	os.Stdout = origStdout

	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		entries, _ := os.ReadDir(tmpDir)
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("--save: expected file %q to exist, got: %v", expectedFile, names)
	}
}
