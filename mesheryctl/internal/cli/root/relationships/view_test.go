package relationships

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestView(t *testing.T) {
	mesheryctlflags.InitValidators(RelationshipCmd)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	tests := []utils.MesheryListCommandTest{
		{
			Name:             "given no model name provided when running relationship view then throw error",
			Args:             []string{"view"},
			URL:              "/api/registry/models/kubernetes/relationships",
			Fixture:          "view.relationship.empty.response.golden",
			ExpectedResponse: "",
			IsOutputGolden:   false,
			ExpectError:      true,
			ExpectedError:    utils.ErrInvalidArgument(errors.New(errInvalidArg)),
		},
		{
			Name:             "given model name provided when running relationship view then display registered relationship",
			Args:             []string{"view", "kubernetes"},
			URL:              "/api/registry/models/kubernetes/relationships?page=0&pagesize=10",
			Fixture:          "view.relationship.api.response.golden",
			ExpectedResponse: "view.relationship.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given non existing model name provided when running relationship view then display no relationship found",
			Args:             []string{"view", "nonexistent"},
			URL:              "/api/registry/models/nonexistent/relationships?page=0&pagesize=10",
			Fixture:          "view.relationship.empty.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrNotFound(fmt.Errorf("No relationship(s) found for the model with name: %s", "nonexistent")),
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, RelationshipCmd, tests, currDir, "relationships")
}

func TestRelationshipViewSaveCreatesFileWithExtension(t *testing.T) {
	mesheryctlflags.InitValidators(RelationshipCmd)

	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(RelationshipCmd, t)

	utils.TokenFlag = utils.GetToken(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	apiResponse := utils.NewGoldenFile(t, "view.relationship.save.api.response.golden", fixturesDir).Load()

	httpmock.RegisterResponder("GET",
		testContext.BaseURL+"/api/registry/models/kubernetes/relationships",
		httpmock.NewStringResponder(200, apiResponse))

	tmpDir := t.TempDir()
	origMesheryFolder := utils.MesheryFolder
	utils.MesheryFolder = tmpDir
	t.Cleanup(func() { utils.MesheryFolder = origMesheryFolder })

	expectedFile := filepath.Join(tmpDir, "relationship_kubernetes_aaaabbbb.json")

	origStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		_ = w.Close()
		os.Stdout = origStdout
	}()

	_ = utils.SetupMeshkitLoggerTesting(t, false)
	RelationshipCmd.SetArgs([]string{"view", "kubernetes", "--output-format", "json", "--save"})
	if err := RelationshipCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_ = w.Close()
	os.Stdout = origStdout

	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		entries, _ := os.ReadDir(tmpDir)
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("--save: expected file %q to exist, got: %v", expectedFile, names)
	}
}
