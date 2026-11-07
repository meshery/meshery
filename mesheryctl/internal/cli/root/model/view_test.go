package model

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestViewModel(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}

	const modelName = "model-test-0"
	const nonExistentModelName = "non-existent-model"
	const modelId = "00000000-0000-0000-0000-000000000000"

	currDir := filepath.Dir(filename)
	modelsApiPath = "api/registry/models"

	tests := []utils.MesheryCommandTest{
		{
			Name:           "given no argument provided when running mesheryctl model view then an error message is displayed",
			Args:           []string{"view"},
			URL:            "/api/registry/models",
			HttpMethod:     "GET",
			HttpStatusCode: 200,
			Fixture:        "list.model.empty.api.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(errors.New(errInvalidArg)),
			IsOutputGolden: false,
		},
		{
			Name:             "given a valid model-id when model view then display detailed information",
			Args:             []string{"view", modelId},
			URL:              fmt.Sprintf("/api/registry/models?id=%s&page=0&pagesize=10", url.QueryEscape(modelId)),
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			Fixture:          "list.model.api.response.golden",
			ExpectedResponse: "view.model.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a valid model-name when model view then display detailed information",
			Args:             []string{"view", modelName},
			URL:              fmt.Sprintf("/api/registry/models?page=0&pagesize=10&search=%s", url.QueryEscape(modelName)),
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			Fixture:          "list.model.api.response.golden",
			ExpectedResponse: "view.model.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a non-existent model-name when model view then display not found message",
			Args:             []string{"view", nonExistentModelName},
			URL:              fmt.Sprintf("/api/registry/models?page=0&pagesize=10&search=%s", url.QueryEscape(nonExistentModelName)),
			HttpMethod:       "GET",
			HttpStatusCode:   404,
			Fixture:          "list.model.api.response.golden",
			ExpectedResponse: "view.model.empty.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given multiple model-names when model view then throw error",
			Args:             []string{"view", "name1", "name2"},
			URL:              "",
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			Fixture:          "list.model.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    utils.ErrInvalidArgument(errors.New(errInvalidArg)),
			IsOutputGolden:   false,
		},
		{
			Name:           "given an invalid format when model view then throw error",
			Args:           []string{"view", modelName, "--output-format", "invalid-format"},
			URL:            fmt.Sprintf("/%s", getModelViewUrlPath(url.QueryEscape(modelName))),
			HttpMethod:     "GET",
			HttpStatusCode: 200,
			Fixture:        "list.model.empty.api.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --output-format 'invalid-format': valid values are json yaml")),
			IsOutputGolden: false,
		},
		{
			Name:             "given a valid format when model view then display detailed information",
			Args:             []string{"view", modelName, "--output-format", "json"},
			URL:              fmt.Sprintf("/%s?page=0&pagesize=10&search=%s", modelsApiPath, url.QueryEscape(modelName)),
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			Fixture:          "list.model.api.response.golden",
			ExpectedResponse: "view.json.api.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
	}

	mesheryctlflags.InitValidators(ModelCmd)
	utils.InvokeMesheryctlTestCommand(t, update, ModelCmd, tests, currDir, "model")
}

// TestViewModelSaveFlattensSeparatorsInName verifies that --save writes one
// file for a model whose display name has path separators, and that the real
// ~/.meshery is never touched.
func TestViewModelSaveFlattensSeparatorsInName(t *testing.T) {
	const modelId = "00000000-0000-0000-0000-000000000000"

	modelsApiPath = "api/registry/models"
	mesheryctlflags.InitValidators(ModelCmd)
	mesheryDir := utils.IsolateMesheryHome(t)

	run := utils.MesheryCommandRun{
		Cmd:        ModelCmd,
		Args:       []string{"view", modelId, "--save"},
		CommandDir: utils.CallerDir(t),
		Mocks: []utils.MockURL{
			{
				URL:      fmt.Sprintf("/api/registry/models?id=%s&page=0&pagesize=10", url.QueryEscape(modelId)),
				Response: "view.model.save.name.api.response.golden",
			},
		},
	}
	if _, err := utils.RunMesheryctlCommand(t, run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	utils.AssertFileExists(t, filepath.Join(mesheryDir, "model_AMD_GPU_Pro_v1.yaml"))
}
