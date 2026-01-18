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

	tests := []utils.MesheryCommandTest{
		{
			Name:             "Fetch Filter List",
			Args:             []string{"list"},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              "/api/filter?pagesize=25&page=0",
			Fixture:          "list.filter.api.response.golden",
			ExpectedContains: []string{"No WASM Filter to display"},
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, FilterCmd, tests, currDir, "filter")
}
