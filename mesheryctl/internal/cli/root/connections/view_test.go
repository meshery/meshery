package connections

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestConnectionViewCmd(t *testing.T) {

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:           "given no argument provided when connection view then throw error",
			Args:           []string{"view"},
			Fixture:        "view.connection.api.empty.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errNoArgMsg, viewUsageMsg)),
			IsOutputGolden: false,
		},
		{
			Name:           "given multiple arguments provided when connection view then throw error",
			Args:           []string{"view", "foo", "bar"},
			Fixture:        "view.connection.api.empty.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errMultiArgMsg, viewUsageMsg)),
			IsOutputGolden: false,
		},
		{
			Name:           "given an invalid argument for --output-format flag provided when connection view then throw error",
			Args:           []string{"view", connectionId, "--output-format", "foo"},
			URL:            "/api/integrations/connections/" + connectionId,
			Fixture:        "view.connection.api.empty.response.golden",
			ExpectError:    true,
			ExpectedError:  display.ErrInvalidOutputFormat("foo"),
			IsOutputGolden: false,
		},
		{
			Name:             "given a valid connection-id provided when connection view then display detailed information",
			Args:             []string{"view", connectionId},
			URL:              "/api/integrations/connections/" + connectionId,
			Fixture:          "view.connection.api.response.golden",
			ExpectedResponse: "view.connection.yaml.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a valid --output-format argument provided when connection view then display detailed information in the specified format",
			Args:             []string{"view", connectionId, "--output-format", "yaml"},
			URL:              "/api/integrations/connections/" + connectionId,
			Fixture:          "view.connection.api.response.golden",
			ExpectedResponse: "view.connection.yaml.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
	}

	// Run tests
	utils.InvokeMesheryctlTestListCommand(t, update, ConnectionsCmd, tests, currDir, "connection")
}
