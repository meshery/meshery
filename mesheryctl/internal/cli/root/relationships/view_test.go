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
			Name:             "given empty model name provided when running relationship view then throw error",
			Args:             []string{"view", "   "},
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
			Name:             "given model name with spaces when running relationship view then escape URL path",
			Args:             []string{"view", "Kubernetes Manifest"},
			URL:              "/api/registry/models/Kubernetes%20Manifest/relationships?page=0&pagesize=10",
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

func TestRelationshipViewSaveSanitizesSpacesInFilename(t *testing.T) {
	mesheryctlflags.InitValidators(RelationshipCmd)

	mesheryDir := utils.IsolateMesheryHome(t)

	run := utils.MesheryCommandRun{
		Cmd:        RelationshipCmd,
		Args:       []string{"view", "Kubernetes Manifest", "--output-format", "json", "--save"},
		CommandDir: utils.CallerDir(t),
		Mocks: []utils.MockURL{
			{
				URL:      "/api/registry/models/Kubernetes%20Manifest/relationships",
				Response: "view.relationship.save.spaces.api.response.golden",
			},
		},
	}
	if _, err := utils.RunMesheryctlCommand(t, run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	utils.AssertFileExists(t, filepath.Join(mesheryDir, "relationship_Kubernetes_Manifest_aaaabbbb.json"))
}

func TestFormatSaveFileName(t *testing.T) {
	tests := []struct {
		name         string
		modelName    string
		id           string
		outputFormat string
		expected     string
	}{
		{
			name:         "model name with spaces",
			modelName:    "Kubernetes Manifest",
			id:           "aaaabbbb-cccc-dddd-eeee-111122223333",
			outputFormat: "json",
			expected:     "relationship_Kubernetes_Manifest_aaaabbbb.json",
		},
		{
			name:         "model name with forward slash and backslash",
			modelName:    "repo/model\\sub",
			id:           "aaaabbbb-cccc-dddd-eeee-111122223333",
			outputFormat: "yaml",
			expected:     "relationship_repo_model_sub_aaaabbbb.yaml",
		},
		{
			name:         "model name with mixed spaces and slashes",
			modelName:    "my model / test\\path",
			id:           "aaaabbbb",
			outputFormat: "JSON",
			expected:     "relationship_my_model___test_path_aaaabbbb.json",
		},
		{
			name:         "model name with leading and trailing spaces",
			modelName:    "  Kubernetes Manifest  ",
			id:           "aaaabbbb",
			outputFormat: "yaml",
			expected:     "relationship_Kubernetes_Manifest_aaaabbbb.yaml",
		},
		{
			name:         "model name with colon and question mark",
			modelName:    "Kubernetes:Manifest?test",
			id:           "aaaabbbb",
			outputFormat: "json",
			expected:     "relationship_Kubernetes_Manifest_test_aaaabbbb.json",
		},
		{
			name:         "model name with quotes, asterisk, pipe and angle brackets",
			modelName:    `model*"<pipe|test>"`,
			id:           "aaaabbbb",
			outputFormat: "yaml",
			expected:     "relationship_model___pipe_test___aaaabbbb.yaml",
		},
		{
			name:         "empty model name falls back to shortID",
			modelName:    "",
			id:           "aaaabbbb-cccc-dddd-eeee-111122223333",
			outputFormat: "json",
			expected:     "relationship_aaaabbbb.json",
		},
		{
			name:         "whitespace-only model name falls back to shortID",
			modelName:    "   ",
			id:           "aaaabbbb",
			outputFormat: "yaml",
			expected:     "relationship_aaaabbbb.yaml",
		},
		{
			name:         "empty model name and empty ID falls back to unknown",
			modelName:    "",
			id:           "",
			outputFormat: "yaml",
			expected:     "relationship_unknown.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := formatSaveFileName(tt.modelName, tt.id, tt.outputFormat)
			if actual != tt.expected {
				t.Errorf("formatSaveFileName(%q, %q, %q) = %q, want %q", tt.modelName, tt.id, tt.outputFormat, actual, tt.expected)
			}
		})
	}
}
