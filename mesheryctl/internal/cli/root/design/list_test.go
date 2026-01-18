package design

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestDesignListCmd(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenrios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "Fetch Design List",
			Args:             []string{"list", "--page", "1"},
			Fixture:          "list.design.api.response.golden",
			URL:              "/api/pattern?page=0&pagesize=10",
			ExpectedContains: []string{"Total number of designs: 343", "Page: 1", "DESIGN ID", "Untitled Design"},
			ExpectError:      false,
		},
		{
			Name:             "Fetch Design List with Local provider",
			Args:             []string{"list", "--page", "1"},
			Fixture:          "list.design.local.api.response.golden",
			URL:              "/api/pattern?page=0&pagesize=10",
			ExpectedContains: []string{"Total number of designs: 1", "Page: 1", "DESIGN ID", "IstioSM"},
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, DesignCmd, tests, currDir, "design")
}
