package model

import (
	goerrors "errors"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestSearchModel(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	querySearch := "model-test"
	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "Search model without query",
			Args:             []string{"search"},
			URL:              "",
			Fixture:          "empty.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(goerrors.New(errSearchModelName)),
		},
		{
			Name:     "Search model with query",
			Args:     []string{"search", querySearch, "--page", "1"},
			URL:      fmt.Sprintf("/%s?search=%s&pagesize=all", modelsApiPath, url.QueryEscape(querySearch)),
			Fixture:  "search.model.api.response.golden",
			ExpectedContains: []string{
				"Total number of models: 41",
				"MODEL",
				"CATEGORY",
				"VERSION",
				"model-test-0",
				"category-test-0",
				"v1.0.0",
			},
			ExpectError: false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ModelCmd, tests, currDir, "model")
}
