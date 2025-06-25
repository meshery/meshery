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
	// Get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// Example model name to use
	querySearch := "model-test"

	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "View model without query",
			Args:             []string{"view"},
			URL:              "",
			Fixture:          "empty.golden",
			ExpectedResponse: "view.model.without.query.output.golden",
			ExpectError:      true,
		},
		{
			Name:             "View model with too many arguments",
			Args:             []string{"view", "model1", "model2"},
			URL:              "",
			Fixture:          "empty.golden",
			ExpectedResponse: "view.model.too.many.arguments.output.golden",
			ExpectError:      true,
		},
		{
			Name:             "View model with invalid output format",
			Args:             []string{"view", querySearch, "--output-format", "xml"},
			URL:              "",
			Fixture:          "empty.golden",
			ExpectedResponse: "view.model.invalid.output.format.golden",
			ExpectError:      true,
		},
		{
			Name:             "View model no models found",
			Args:             []string{"view", querySearch, "--output-format", "yaml"},
			URL:              fmt.Sprintf("/%s/%s?pagesize=all", modelsApiPath, url.PathEscape(querySearch)),
			Fixture:          "view.model.no.models.api.response.golden",
			ExpectedResponse: "view.model.no.models.found.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "View model exactly one found - YAML output",
			Args:             []string{"view", querySearch, "--output-format", "yaml"},
			URL:              fmt.Sprintf("/%s/%s?pagesize=all", modelsApiPath, url.PathEscape(querySearch)),
			Fixture:          "view.model.one.found.api.response.golden",
			ExpectedResponse: "view.model.one.found.yaml.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "View model exactly one found - JSON output",
			Args:             []string{"view", querySearch, "--output-format", "json"},
			URL:              fmt.Sprintf("/%s/%s?pagesize=all", modelsApiPath, url.PathEscape(querySearch)),
			Fixture:          "view.model.one.found.api.response.golden",
			ExpectedResponse: "view.model.one.found.json.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ModelCmd, tests, currDir, "model")
}
