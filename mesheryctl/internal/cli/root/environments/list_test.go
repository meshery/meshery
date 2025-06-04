package environments

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestListEnvironment(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "List environments organization ID not provided",
			Args:             []string{"list"},
			URL:              "/api/environments",
			Fixture:          "list.environment.without.orgID.golden",
			ExpectedResponse: "list.environment.without.orgID.golden",
			ExpectError:      true,
		},
		{
			Name:             "List environments non available",
			Args:             []string{"list", "--orgID", testConstants["orgID"]},
			URL:              fmt.Sprintf("/api/environments?orgID=%s", testConstants["orgID"]),
			Fixture:          "list.environment.empty.response.golden",
			ExpectedResponse: "list.environment.empty.golden",
			ExpectError:      true,
		},
		{
			Name:             "List environments available",
			Args:             []string{"list", "--orgID", testConstants["orgID"]},
			URL:              fmt.Sprintf("/api/environments?orgID=%s", testConstants["orgID"]),
			Fixture:          "list.environment.response.golden",
			ExpectedResponse: "list.environment.success.golden",
			ExpectError:      true,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, EnvironmentCmd, tests, currDir, "environments")
}
