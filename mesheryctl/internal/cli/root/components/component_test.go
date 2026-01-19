package components

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
)

func TestComponent(t *testing.T) {

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "Given invalid subcommand name provided trigger an error",
			Args:             []string{"invalidCommand"},
			URL:              "",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("'%s' is an invalid subcommand. Please provide required options from [list, search, view]. Use 'mesheryctl component --help' to display usage guide", "invalidCommand")),
		},
		{
			Name:             "Given no args provided and flag count is false trigger an error",
			Args:             []string{},
			URL:              "",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(errors.New("please provide a subcommand")),
		},
		{
			Name:             "Given an empty args provided and flags count is false trigger an error",
			Args:             []string{""},
			URL:              "",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(errors.New("please provide a subcommand")),
		},
		{
			Name:             "display components count",
			Args:             []string{"--count"},
			URL:              fmt.Sprintf("/%s", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedContains: []string{"Total number of components: 1"},
			ExpectError:      false,
			ExpectedError:    nil,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ComponentCmd, tests, currentDirectory, "component")

}
