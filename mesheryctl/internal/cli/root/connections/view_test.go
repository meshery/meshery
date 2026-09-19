package connections

import (
	"bytes"
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
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

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

	utils.InvokeMesheryctlTestListCommand(t, update, ConnectionsCmd, tests, currDir, "connection")
}

// runConnectionViewTest handles all shared scaffolding — mock setup, token,
// cobra output capture, and command execution. Tests only contain what is
// unique to their scenario.
func runConnectionViewTest(t *testing.T, args []string) error {
	t.Helper()
	testContext := utils.InitTestEnvironment(t)
	t.Cleanup(func() { utils.StopMockery(t) })
	t.Cleanup(func() { utils.ResetCommandFlags(ConnectionsCmd, t) })

	utils.TokenFlag = utils.GetToken(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine current working directory")
	}
	fixturesDir := filepath.Join(filepath.Dir(filename), "fixtures")

	apiResponse := utils.NewGoldenFile(t, "view.connection.api.response.golden", fixturesDir).Load()
	httpmock.RegisterResponder("GET",
		testContext.BaseURL+"/api/integrations/connections/"+connectionId,
		httpmock.NewStringResponder(200, apiResponse))

	buf := &bytes.Buffer{}
	ConnectionsCmd.SetOut(buf)
	ConnectionsCmd.SetErr(buf)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	ConnectionsCmd.SetArgs(args)
	return ConnectionsCmd.Execute()
}

// TestConnectionViewSaveCreatesFile verifies --save writes a file with the
// correct name and extension. Uses a temp dir so ~/.meshery is never touched.
func TestConnectionViewSaveCreatesFile(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	mesheryDir := filepath.Join(tmpHome, ".meshery")
	if err := os.MkdirAll(mesheryDir, 0755); err != nil {
		t.Fatalf("cannot create %s: %v", mesheryDir, err)
	}

	if err := runConnectionViewTest(t, []string{"view", connectionId, "--save"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFile := filepath.Join(mesheryDir, "connection_minikube.yaml")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		entries, _ := os.ReadDir(mesheryDir)
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("--save: expected file %q to exist, got: %v", expectedFile, names)
	}
}

// TestConnectionViewNoSaveWithBrokenHome verifies that view without --save
// succeeds even when HOME is unset — proving os.UserHomeDir() is not called
// on the non-save path.
func TestConnectionViewNoSaveWithBrokenHome(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")

	if err := runConnectionViewTest(t, []string{"view", connectionId}); err != nil {
		t.Fatalf("view without --save should succeed even with no HOME: %v", err)
	}
}
