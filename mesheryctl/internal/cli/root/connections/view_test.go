package connections

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
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

// TestConnectionViewSaveCreatesFile verifies that "connection view --save"
// writes a file whose name includes the connection name and the format
// extension (e.g. connection_minikube.yaml).
//
// This is a regression guard for the bug where os.UserHomeDir() was called
// unconditionally — moving it inside the --save block means the non-save
// path (exercised by TestConnectionViewCmd above) cannot be broken by a
// missing HOME directory.
func TestConnectionViewSaveCreatesFile(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(ConnectionsCmd, t)

	utils.TokenFlag = utils.GetToken(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	apiResponse := utils.NewGoldenFile(t, "view.connection.api.response.golden", fixturesDir).Load()
	httpmock.RegisterResponder("GET",
		testContext.BaseURL+"/api/integrations/connections/"+connectionId,
		httpmock.NewStringResponder(200, apiResponse))

	// Determine the save path that the command will use: ~/.meshery/connection_minikube.yaml
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("skipping: cannot resolve home directory: %v", err)
	}
	mesheryDir := filepath.Join(homeDir, ".meshery")
	if mkErr := os.MkdirAll(mesheryDir, 0755); mkErr != nil {
		t.Fatalf("cannot create %s: %v", mesheryDir, mkErr)
	}
	expectedFile := filepath.Join(mesheryDir, "connection_minikube.yaml")
	t.Cleanup(func() { _ = os.Remove(expectedFile) })

	origStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		_ = w.Close()
		os.Stdout = origStdout
	}()

	_ = utils.SetupMeshkitLoggerTesting(t, false)
	ConnectionsCmd.SetArgs([]string{"view", connectionId, "--save"})
	if execErr := ConnectionsCmd.Execute(); execErr != nil {
		t.Fatalf("unexpected error: %v", execErr)
	}

	_ = w.Close()
	os.Stdout = origStdout

	if _, statErr := os.Stat(expectedFile); os.IsNotExist(statErr) {
		t.Errorf("--save flag: expected file %q to exist after 'connection view --save', but it was not created", expectedFile)
	}
}
