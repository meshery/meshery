package design

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")

func TestDesignListCmd(t *testing.T) {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenrios for fetching data
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "Fetch Pattern List",
			Args:             []string{"list", "--page", "2"},
			ExpectedResponse: "list.design.output.golden",
			Fixture:          "list.design.api.response.golden",
			URL:              "/api/pattern",
			ExpectError:      false,
		},
		{
			Name:             "Fetch Pattern List with Local provider",
			Args:             []string{"list", "--page", "1"},
			ExpectedResponse: "list.design.local.output.golden",
			Fixture:          "list.design.local.api.response.golden",
			URL:              "/api/pattern",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, DesignCmd, tests, currDir, "designs")
}
