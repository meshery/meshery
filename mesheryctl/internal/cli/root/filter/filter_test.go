package filter

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")

// TestFilterCmd tests the parent filter command behavior
// (not the subcommands - those have dedicated test files)
func TestFilterCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// Test cases for the parent FilterCmd
	tests := []utils.MesheryCommandTest{
		{
			Name:             "Filter command with no args",
			Args:             []string{},
			ExpectedResponse: "filter.help.output.golden",
			ExpectError:      false,
		},
	}

	// Run tests using the pre-existing helper
	utils.InvokeMesheryctlTestCommand(t, update, FilterCmd, tests, currDir, "filter")
}
