package components

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestListComponent(t *testing.T) {

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "list components with page number",
			Args:             []string{"list"},
			URL:              fmt.Sprintf("/%s", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.list.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "list components non empty count",
			Args:             []string{"list", "--count"},
			URL:              fmt.Sprintf("/%s?pagesize=all", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.list.count.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "list components empty count",
			Args:             []string{"list", "--count"},
			URL:              fmt.Sprintf("/%s?pagesize=all", componentApiPath),
			Fixture:          "components.empty.api.response.golden",
			ExpectedResponse: "components.list.count.empty.ouput.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ComponentCmd, tests, currentDirectory, "component")

}
