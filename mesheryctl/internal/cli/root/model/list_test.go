package model

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestListModel(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "List model",
			Args:             []string{"list"},
			URL:              fmt.Sprintf("/%s", modelsApiPath),
			Fixture:          "list.model.api.response.golden",
			ExpectedResponse: "list.model.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "List model no data",
			Args:             []string{"list"},
			URL:              fmt.Sprintf("/%s", modelsApiPath),
			Fixture:          "list.model.empty.api.response.golden",
			ExpectedResponse: "list.model.empty.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ModelCmd, tests, currDir, "model")
}
