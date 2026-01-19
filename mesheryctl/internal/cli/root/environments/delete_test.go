package environments

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
)

func TestDeleteEnvironment(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	environmentId := "d56fb25b-f92c-4cd6-821b-2cfd6bb87259"

	// Test scenarios for environment deletion
	tests := []utils.MesheryCommandTest{
		{
			Name:             "Delete environment without arguments",
			Args:             []string{"delete"},
			HttpMethod:       "DELETE",
			HttpStatusCode:   200,
			URL:              "/api/environments",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(errors.New("[ Environment ID ] isn't specified\n\nUsage: mesheryctl environment delete [environmentId]\nRun 'mesheryctl environment delete --help' to see detailed help message")),
		},
		{
			Name:             "Delete environment successfully",
			Args:             []string{"delete", environmentId},
			HttpMethod:       "DELETE",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/api/environments/%s", environmentId),
			Fixture:          "delete.environment.response.golden",
			ExpectedContains: []string{"has been deleted"},
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, EnvironmentCmd, tests, currDir, "environments")
}
