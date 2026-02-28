package model

import (
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestViewModel(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}

	const modelName = "model-test-0"
	const modelId = "00000000-0000-0000-0000-000000000000"

	currDir := filepath.Dir(filename)
	modelsApiPath = "api/meshmodels/models"

	// test scenarios for fetching data
	tests := []utils.MesheryCommandTest{
		{
			Name:           "given no argument provided when running mesheryctl model view then an error message is displayed",
			Args:           []string{"view"},
			URL:            "api/meshmodels/models",
			HttpMethod:     "GET",
			HttpStatusCode: 200,
			Fixture:        "list.model.empty.api.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("%s%s", errNoArg, viewUsageMsg)),
			IsOutputGolden: false,
		},
		{
			Name:             "given an invalid model-name provided when running mesheryctl model view invalid-name then an error message is displayed",
			Args:             []string{"view", "invalid"},
			URL:              fmt.Sprintf("/%s/%s?pagesize=all", modelsApiPath, url.QueryEscape("invalid")),
			Fixture:          "list.model.empty.api.response.golden",
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			ExpectedResponse: "view.model.empty.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a valid model-id provided when running mesheryctl model view valid-id then the detailed information of the model is displayed",
			Args:             []string{"view", modelId},
			URL:              fmt.Sprintf("/%s?id=%s", modelsApiPath, url.QueryEscape(modelId)),
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			Fixture:          "list.model.api.response.golden",
			ExpectedResponse: "view.model.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a valid model-name provided when running mesheryctl model view valid-name then the detailed information of the model is displayed",
			Args:             []string{"view", modelName},
			URL:              fmt.Sprintf("/%s/%s?pagesize=all", modelsApiPath, url.QueryEscape(modelName)),
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			Fixture:          "list.model.api.response.golden",
			ExpectedResponse: "view.model.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a multiple model-name provided when running mesheryctl model view name1 name2 then an error message is displayed",
			Args:             []string{"view", "name1", "name2"},
			URL:              "",
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			Fixture:          "list.model.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("%s%s", errMultiArg, viewUsageMsg)),
			IsOutputGolden:   false,
		},
		{
			Name:           "given an invalid format provided for --output-format flag when running mesheryctl model view valid-name --output-format invalid-format then an error message is displayed",
			Args:           []string{"view", modelName, "--output-format", "invalid-format"},
			URL:            fmt.Sprintf("/%s/%s?pagesize=all", modelsApiPath, url.QueryEscape(modelName)),
			HttpMethod:     "GET",
			HttpStatusCode: 200,
			Fixture:        "list.model.empty.api.response.golden",
			ExpectError:    true,
			ExpectedError:  ErrModelUnsupportedOutputFormat(formaterrMsg),
			IsOutputGolden: false,
		},
		{
			Name:             "given a valid format provided when running mesheryctl model view valid-name --output-format valid-format then a detailed information of the model is displayed",
			Args:             []string{"view", modelName, "--output-format", "json"},
			URL:              fmt.Sprintf("/%s/%s?pagesize=all", modelsApiPath, url.QueryEscape(modelName)),
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			Fixture:          "list.model.api.response.golden",
			ExpectedResponse: "view.json.api.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, ModelCmd, tests, currDir, "model")
}
