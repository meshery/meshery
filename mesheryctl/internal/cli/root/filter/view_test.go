package filter

import (
	"io"
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

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching filter data
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

	// Run tests that capture stdout output
	utils.InvokeMesheryctlTestListCommand(t, update, FilterCmd, ListTests, currDir, "filter")

	// Run tests that capture logger output
	utils.InvokeMesheryctlTestCommand(t, update, FilterCmd, loggerTests, currDir, "filter")
}

func TestGetFilterViewFilePath_AppendsExtension(t *testing.T) {
	originalMesheryFolder := utils.MesheryFolder
	originalDefaultConfigPath := utils.DefaultConfigPath
	originalTokenFlag := utils.TokenFlag
	originalArgs := append([]string(nil), FilterCmd.Flags().Args()...)
	originalOut := FilterCmd.OutOrStdout()
	t.Cleanup(func() {
		utils.MesheryFolder = originalMesheryFolder
		utils.DefaultConfigPath = originalDefaultConfigPath
		utils.TokenFlag = originalTokenFlag
		FilterCmd.SetArgs(originalArgs)
		FilterCmd.SetOut(originalOut)
		utils.ResetCommandFlags(FilterCmd, t)
	})

	utils.MesheryFolder = t.TempDir()
	utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	testContext := utils.InitTestEnvironment(t)
	t.Cleanup(func() {
		utils.StopMockery(t)
	})

	fixturesDir := filepath.Join(currDir, "fixtures")
	apiResponse := utils.NewGoldenFile(t, "view.filter.api.response.golden", fixturesDir).Load()
	httpmock.RegisterResponder("GET", testContext.BaseURL+"/api/filter", httpmock.NewStringResponder(200, apiResponse))
	utils.TokenFlag = utils.GetToken(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	mesheryctlflags.InitValidators(FilterCmd)
	FilterCmd.SetArgs([]string{"view", "KumaTest", "--output-format", "json", "--save"})
	FilterCmd.SetOut(io.Discard)
	if err := FilterCmd.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedBase := filepath.Join(utils.MesheryFolder, "filter_KumaTest_957fbc9b")
	jsonPath := expectedBase + ".json"
	if _, err := os.Stat(jsonPath); err != nil {
		t.Fatalf("expected saved file %q to exist, got %v", jsonPath, err)
	}
	if _, err := os.Stat(expectedBase); !os.IsNotExist(err) {
		t.Fatalf("expected unsuffixed save path %q to not exist, got %v", expectedBase, err)
	}
}
