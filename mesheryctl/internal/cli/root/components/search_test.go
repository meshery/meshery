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
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("%v\n\n%v", errNoArg, usageErrorMessage)),
		},
		{
			Name:             "given a multiple argument is provided when running mesheryctl component search arg1 arg2 then an error message is displayed",
			Args:             []string{"search", "args1", "args2"},
			URL:              "",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("%v\n\n%v", errMultiArg, usageErrorMessage)),
		},
		{
			Name:             "given a valid component is provided when running mesheryctl component search valid-name then it displays every matching results in output",
			Args:             []string{"search", "Test"},
			URL:              fmt.Sprintf("/%s?search=Test&page=0&pagesize=10", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.search.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given a valid component is provided when running mesheryctl component search valid-name --page int then it displays every matching results in output",
			Args:             []string{"search", "Test", "--page", "1"},
			URL:              fmt.Sprintf("/%s?search=Test&page=0&pagesize=10", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.search.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ComponentCmd, tests, currentDirectory, "component")

}
