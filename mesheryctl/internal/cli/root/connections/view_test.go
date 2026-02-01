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
			Name:           "View connection without connection id",
			Args:           []string{"view"},
			URL:            "/api/integrations/connections",
			Fixture:        "view.connection.api.empty.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("connection name or ID isn't specified\n")),
			IsOutputGolden: false,
		},
		{
			Name:           "View connection with multiple connection id",
			Args:           []string{"view", "foo", "bar"},
			Fixture:        "view.connection.api.empty.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("too many arguments\n")),
			IsOutputGolden: false,
		},
		{
			Name:             "View connection with invalid connection id",
			Args:             []string{"view", "foo"},
			URL:              "/api/integrations/connections/foo",
			Fixture:          "view.connection.api.empty.response.golden",
			ExpectedResponse: "view.connection.invalid.output.golden",
			ExpectError:      true,
			IsOutputGolden:   true,
		},
		{
			Name:             "View connection with --output-format",
			Args:             []string{"view", ConnectionID, "--output-format"},
			URL:              "/api/integrations/connections/" + ConnectionID,
			Fixture:          "view.connection.api.empty.response.golden",
			ExpectedResponse: "view.connection.invalid.outputflag.output.golden",
			ExpectError:      true,
			IsOutputGolden:   true,
		},
		{
			Name:           "View connection with --output-format foo",
			Args:           []string{"view", ConnectionID, "--output-format", "foo"},
			URL:            "/api/integrations/connections/" + ConnectionID,
			Fixture:        "view.connection.api.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(errors.New(invalidOutputFormatMsg)),
			IsOutputGolden: false,
		},
		{
			Name:             "View connection with --output-format yaml",
			Args:             []string{"view", ConnectionID, "--output-format", "yaml"},
			URL:              "/api/integrations/connections/" + ConnectionID,
			Fixture:          "view.connection.api.response.golden",
			ExpectedResponse: "view.connection.yaml.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "View connection with --output-format json",
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
