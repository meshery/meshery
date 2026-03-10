package filter

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestViewCmd(t *testing.T) {
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
