package model

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestViewModel(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "View model by name (success)",
			Args:             []string{"view", "test-model"},
			URL:              "/models/test-model?pagesize=all",
			Fixture:          "view.model.api.response.golden",
			ExpectedResponse: "view.model.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "View model by name (not found)",
			Args:             []string{"view", "notfound-model"},
			URL:              "/models/notfound-model?pagesize=all",
			Fixture:          "view.model.notfound.api.response.golden",
			ExpectedResponse: "view.model.notfound.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "View model with multiple versions (prompt)",
			Args:             []string{"view", "multi-version-model"},
			URL:              "/models/multi-version-model?pagesize=all",
			Fixture:          "view.model.multiversion.api.response.golden",
			ExpectedResponse: "view.model.multiversion.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Invalid output format",
			Args:             []string{"view", "test-model", "--output-format", "invalid"},
			URL:              "",
			Fixture:          "",
			ExpectedResponse: "view.model.invalid-output-format.output.golden",
			ExpectError:      true,
		},
		{
			Name:             "Missing model name argument",
			Args:             []string{"view"},
			URL:              "",
			Fixture:          "",
			ExpectedResponse: "view.model.missing-argument.output.golden",
			ExpectError:      true,
		},
		{
			Name:             "Too many arguments",
			Args:             []string{"view", "model1", "model2"},
			URL:              "",
			Fixture:          "",
			ExpectedResponse: "view.model.too-many-arguments.output.golden",
			ExpectError:      true,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ModelCmd, tests, currDir, "model")
} 