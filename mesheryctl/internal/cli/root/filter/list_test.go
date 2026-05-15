package filter

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestListCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenrios for fetching data
	tests := []utils.MesheryCommandTest{
		{
			Name:             "Fetch Filter List",
			Args:             []string{"list"},
			ExpectedResponse: "list.filter.output.golden",
			Fixture:          "list.filter.api.response.golden",
			URL:              "/api/filter",
			ExpectError:      false,
			HttpMethod:       "GET",
			HttpStatusCode:   200,
		},
	}
	// Run tests
	utils.InvokeMesheryctlTestCommand(t, update, FilterCmd, tests, currDir, "filter")
}
