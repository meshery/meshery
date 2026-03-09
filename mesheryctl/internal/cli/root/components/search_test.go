package components

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestSearchComponent(t *testing.T) {

	mesheryctlflags.InitValidators(ComponentCmd)
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "given no arguments provided when component search  then throw error",
			Args:             []string{"search"},
			URL:              "",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("%v\n%v", errInvalidArg, searchUsageMsg)),
		},
		{
			Name:             "given multiple argument provided when component search then throw error",
			Args:             []string{"search", "args1", "args2"},
			URL:              "",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("%v\n%v", errInvalidArg, searchUsageMsg)),
		},
		{
			Name:             "given valid name provided when component search then display matching results",
			Args:             []string{"search", "Test"},
			URL:              fmt.Sprintf("/%s?search=Test&page=0&pagesize=10", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.search.success.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given valid name provided with page flag when component search then display matching results",
			Args:             []string{"search", "Test", "--page", "2"},
			URL:              fmt.Sprintf("/%s?search=Test&page=1&pagesize=10", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.search.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ComponentCmd, tests, currentDirectory, "component")

}
