package connections

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
)

const ConnectionId = "11111111-1111-1111-1111-111111111111"

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
			Name:           "given no connection-id provided when running mesheryctl connection view then an error message is displayed",
			Args:           []string{"view"},
			Fixture:        "view.connection.api.empty.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("connection name or ID isn't specified\n")),
			IsOutputGolden: false,
		},
		{
			Name:           "given multiple connection-id provided when running mesheryctl connection view foo bar then an error message is displayed",
			Args:           []string{"view", "foo", "bar"},
			Fixture:        "view.connection.api.empty.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("too many arguments\n")),
			IsOutputGolden: false,
		},
		{
			Name:           "given an invalid argument for --output-format flag provided when running mesheryctl connection view connection-id --output-format foo then an error message is displayed",
			Args:           []string{"view", ConnectionID, "--output-format", "foo"},
			URL:            "/api/integrations/connections/" + ConnectionID,
			Fixture:        "view.connection.api.empty.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(errors.New(invalidOutputFormatMsg)),
			IsOutputGolden: false,
		},
		{
			Name:             "given a valid connection-id provided when running mesheryctl connection view connection-id then the detailed output is displayed",
			Args:             []string{"view", ConnectionID},
			URL:              "/api/integrations/connections/" + ConnectionID,
			Fixture:          "view.connection.api.response.golden",
			ExpectedResponse: "view.connection.yaml.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a valid --output-format argument provided when running mesheryctl connection view connection-id --output-format yaml then the output is displayed as specified format",
			Args:             []string{"view", ConnectionID, "--output-format", "yaml"},
			URL:              "/api/integrations/connections/" + ConnectionID,
			Fixture:          "view.connection.api.response.golden",
			ExpectedResponse: "view.connection.yaml.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a valid --output-format argument provided when running mesheryctl connection view connection-id --output-format json then the output is displayed as specified format",
			Args:             []string{"view", ConnectionID, "--output-format", "json"},
			URL:              "/api/integrations/connections/" + ConnectionID,
			Fixture:          "view.connection.api.response.golden",
			ExpectedResponse: "view.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
	}

	// Run tests
	utils.InvokeMesheryctlTestListCommand(t, update, ConnectionsCmd, tests, currDir, "connection")
}
