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
			Name:             "given no arguments and no flags provided when component search then throw error",
			Args:             []string{"search"},
			URL:              "",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("provide a search query or use --model flag\n\n%s", searchUsageMsg)),
		},
		{
			Name:             "given multiple argument provided when component search then throw error",
			Args:             []string{"search", "args1", "args2"},
			URL:              "",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("at most one argument (search query) can be provided\n\n%s", searchUsageMsg)),
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
		{
			Name:             "given model flag provided when component search then display matching results",
			Args:             []string{"search", "--model", "component-test-0"},
			URL:              "/api/meshmodels/models/component-test-0/components?page=0&pagesize=10",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.search.success.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given model and query provided when component search then display matching results",
			Args:             []string{"search", "Test", "--model", "component-test-0"},
			URL:              "/api/meshmodels/models/component-test-0/components?search=Test&page=0&pagesize=10",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.search.success.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given valid name provided when component search with json output then display json",
			Args:             []string{"search", "Test", "-o", "json"},
			URL:              fmt.Sprintf("/%s?search=Test&pagesize=all", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.search.json.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given model flag provided when component search with yaml output then display yaml",
			Args:             []string{"search", "--model", "component-test-0", "-o", "yaml"},
			URL:              "/api/meshmodels/models/component-test-0/components?pagesize=all",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.search.yaml.output.golden",
			ExpectError:      false,
		},
	}

	mesheryctlflags.InitValidators(ComponentCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, ComponentCmd, tests, currentDirectory, "component")

}
