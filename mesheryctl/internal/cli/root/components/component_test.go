package components

import (
	"flag"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")

func TestComponent(t *testing.T) {

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "launch component with invalid subcommand name",
			Args:             []string{"invalidCommand"},
			URL:              "",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.invalid.subcommand.output.golden",
			ExpectError:      true,
		},
		{
			Name:             "display components count",
			Args:             []string{"--count"},
			URL:              fmt.Sprintf("/%s", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.list.count.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ComponentCmd, tests, currentDirectory, "component")

}
