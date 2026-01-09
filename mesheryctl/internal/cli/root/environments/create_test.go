package environments

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
)

func TestCreateEnvironment(t *testing.T) {
	// Get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// Test scenarios for environment creation
	tests := []utils.MesheryCommandTest{
		{
			Name:             "Create environment without arguments",
			Args:             []string{"create"},
			URL:              "/api/environments",
			HttpMethod:       "POST",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(errors.New("[ Organization ID | Name | Description ] aren't specified\n\nUsage: mesheryctl environment create --orgID [orgID] --name [name] --description [description]\nRun 'mesheryctl environment create --help' to see detailed help message")),
		},
		{
			Name:             "Create environment successfully",
			Args:             []string{"create", "--name", testConstants["environmentName"], "--description", "integration test", "--orgID", testConstants["orgID"]},
			URL:              "/api/environments",
			HttpMethod:       "POST",
			HttpStatusCode:   200,
			Fixture:          "create.environment.response.golden",
			ExpectedResponse: "create.environment.success.golden",
			ExpectError:      false,
		},
	}

	// Run tests
	utils.InvokeMesheryctlTestCommand(t, update, EnvironmentCmd, tests, currDir, "environments")
}
