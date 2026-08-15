package filter

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
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

	mesheryctlflags.InitValidators(FilterCmd)
	// Run tests using the pre-existing helper
	utils.InvokeMesheryctlTestCommand(t, update, FilterCmd, tests, currDir, "filter")
}

func TestFilterCmdCount(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	tests := []utils.MesheryListCommandTest{
		{
			Name:             "Fetch Filter Count on root command",
			Args:             []string{"--count"},
			ExpectedResponse: "list.filter.count.output.golden",
			Fixture:          "filter.list.api.response.golden",
			URL:              "/api/filter",
			ExpectError:      false,
		},
	}

	mesheryctlflags.InitValidators(FilterCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, FilterCmd, tests, currDir, "filter")
}
