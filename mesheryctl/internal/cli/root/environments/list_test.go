package environments

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
)

func TestListEnvironment(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "List environments organization ID not provided",
			Args:             []string{"list"},
			URL:              "/api/environments",
			Fixture:          "list.environment.without.orgID.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(errors.New("[ orgID ] isn't specified\n\nUsage: mesheryctl environment list --orgID [orgID]\nRun 'mesheryctl environment list --help' to see detailed help message")),
		},
		{
			Name:             "List environments non available",
			Args:             []string{"list", "--orgID", testConstants["orgID"]},
			URL:              fmt.Sprintf("/api/environments?orgID=%s", testConstants["orgID"]),
			Fixture:          "list.environment.empty.response.golden",
			ExpectedContains: []string{"No environments found"},
			ExpectError:      false,
		},
		{
			Name:             "List environments available",
			Args:             []string{"list", "--orgID", testConstants["orgID"]},
			URL:              fmt.Sprintf("/api/environments?orgID=%s", testConstants["orgID"]),
			Fixture:          "list.environment.response.golden",
			ExpectedContains: []string{"Total number of environments: 1", "NAME", "test-environment"},
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, EnvironmentCmd, tests, currDir, "environments")
}
