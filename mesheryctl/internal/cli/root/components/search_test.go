package components

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestSearchComponent(t *testing.T) {

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "Search components with query parameter",
			Args:             []string{"search"},
			URL:              "",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.search.no.agrs.output.golden",
			ExpectError:      true,
		},
		{
			Name:             "Search components with query parameter",
			Args:             []string{"search", "Test"},
			URL:              fmt.Sprintf("/%s?search=Test&pagesize=all", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.list.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ComponentCmd, tests, currentDirectory, "component")

}
