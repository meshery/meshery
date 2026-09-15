package components

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestListComponent(t *testing.T) {

	mesheryctlflags.InitValidators(ComponentCmd)
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)
	testUrl := fmt.Sprintf("/%s?page=0&pagesize=10", componentApiPath)
	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "list components with page number",
			Args:             []string{"list", "--page", "1"},
			URL:              testUrl,
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.list.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "list components non empty count",
			Args:             []string{"list", "--count"},
			URL:              testUrl,
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.list.count.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "list components empty count",
			Args:             []string{"list", "--count"},
			URL:              testUrl,
			Fixture:          "components.empty.api.response.golden",
			ExpectedResponse: "components.list.count.empty.ouput.golden",
			ExpectError:      false,
		},
		{
			Name:             "list components with model",
			Args:             []string{"list", "--model", "kubernetes"},
			URL:              "/api/registry/models/kubernetes/components?page=0&pagesize=10",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.list.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "list components with model and count",
			Args:             []string{"list", "--model", "kubernetes", "--count"},
			URL:              "/api/registry/models/kubernetes/components?page=0&pagesize=10",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.list.count.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "list components with model and pagination",
			Args:             []string{"list", "--model", "kubernetes", "--page", "2", "--pagesize", "20"},
			URL:              "/api/registry/models/kubernetes/components?page=1&pagesize=20",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.search.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "list components with special character model name",
			Args:             []string{"list", "--model", "test model"},
			URL:              "/api/registry/models/test%20model/components?page=0&pagesize=10",
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.list.output.golden",
			ExpectError:      false,
		},
	}

	mesheryctlflags.InitValidators(ComponentCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, ComponentCmd, tests, currentDirectory, "component")

}
