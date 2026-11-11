package organizations

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestListOrganizations(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "List organizations without providing organization ID",
			Args:             []string{"list"},
			URL:              fmt.Sprintf("/%s?page=0&pagesize=10", organizationsApiPath),
			Fixture:          "list.organization.response.golden",
			ExpectedResponse: "list.organization.golden",
			ExpectError:      false,
		},
		{
			Name:             "given extra positional arguments when organization list then throw error",
			Args:             []string{"list", "extra-arg"},
			URL:              "",
			Fixture:          "list.organization.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", listUsageMsg)),
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, OrgCmd, tests, currentDirectory, "organization")
}
