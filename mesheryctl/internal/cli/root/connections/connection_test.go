package connections

import (
	"flag"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
)

var update = flag.Bool("update", false, "update golden files")

func TestConnection(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	//test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:           "given no subcommands provided when running mesheryctl connection then an error message is displayed",
			Args:           []string{},
			URL:            "",
			Fixture:        "list.connection.api.empty.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(errors.New("no subcommand provided for connection")),
			IsOutputGolden: false,
		},
		{
			Name:           "given an invalid subcommands provided when running mesheryctl connection then an error message is displayed",
			Args:           []string{"invalidCommand"},
			URL:            "",
			Fixture:        "list.connection.api.empty.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("'%s' is an invalid subcommand. Use 'mesheryctl connection --help' to display usage", "invalidCommand")),
			IsOutputGolden: false,
		},
		{
			Name:             "given a --count flag provided when running mesheryctl connection --count then the total count of connections is displayed",
			Args:             []string{"--count"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.api.empty.response.golden",
			ExpectedResponse: "list.count.connection.empty.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ConnectionsCmd, tests, currDir, "connection")
}
