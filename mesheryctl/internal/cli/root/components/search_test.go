package components

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestSearchComponent(t *testing.T) {

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "given no valid-component is provided when running mesheryctl component search no-valid-component then an error message is displays",
			Args:             []string{"search"},
			URL:              "",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("[search term] isn't specified. Please enter component name to search\n\n%v", usageErrorMessage)),
		},
		{
			Name:             "given a valid component is provided when running mesheryctl component search valid-name then it displays every matching results in output",
			Args:             []string{"search", "Test"},
			URL:              fmt.Sprintf("/%s?pagesize=all&search=Test", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.search.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ComponentCmd, tests, currentDirectory, "component")

}
