package model

import (
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

	currDir := filepath.Dir(filename)

	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "View model without query",
			Args:             []string{"view"},
			URL:              "",
			Fixture:          "empty.golden",
			ExpectedResponse: "view.model.without.query.output.golden",
			ExpectError:      true,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ModelCmd, tests, currDir, "model")

}
