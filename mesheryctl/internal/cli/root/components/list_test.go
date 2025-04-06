package components

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")

func TestComponentList(t *testing.T) {

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "list components no results",
			Args:             []string{"list", "--count"},
			URL:              "/api/meshmodels/components",
			Fixture:          "components.list.count.only.empty.golden",
			ExpectedResponse: "components.list.count.only.empty.ouput.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ComponentCmd, tests, currDir, "component")

}
