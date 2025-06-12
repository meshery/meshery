package environments

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestCreateEnvironment(t *testing.T) {
	// Get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// Test scenarios for environment creation
	tests := []utils.MesheryCommamdTest{
		{
			Name:             "Create environment without arguments",
			Args:             []string{"create"},
			URL:              "/api/environments",
			HttpMethod:       "POST",
			Fixture:          "",
			ExpectedResponse: "create.environment.without.name.golden",
			ExpectError:      true,
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
