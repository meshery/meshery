package relationships

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

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

	mesheryDir := utils.IsolateMesheryHome(t)

	run := utils.MesheryCommandRun{
		Cmd:        RelationshipCmd,
		Args:       []string{"view", "kubernetes", "--output-format", "json", "--save"},
		CommandDir: utils.CallerDir(t),
		Mocks: []utils.MockURL{
			{
				URL:      "/api/registry/models/kubernetes/relationships",
				Response: "view.relationship.save.api.response.golden",
			},
		},
	}
	if _, err := utils.RunMesheryctlCommand(t, run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	utils.AssertFileExists(t, filepath.Join(mesheryDir, "relationship_kubernetes_aaaabbbb.json"))
}