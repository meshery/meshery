package model

import (
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
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "Search model without query",
			Args:             []string{"search"},
			URL:              "",
			Fixture:          "empty.golden",
			ExpectedResponse: "search.model.without.query.output.golden",
			ExpectError:      true,
		},
		{
			Name:             "Search model with query",
			Args:             []string{"search", querySearch},
			URL:              fmt.Sprintf("/%s?search=%s&pagesize=all", modelsApiPath, url.QueryEscape(querySearch)),
			Fixture:          "search.model.api.response.golden",
			ExpectedResponse: "search.model.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ModelCmd, tests, currDir, "model")
}
