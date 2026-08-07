package filter

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestListCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenrios for fetching data
	tests := []utils.MesheryCommandTest{
		{
			Name:             "Fetch Filter List",
			Args:             []string{"list"},
			ExpectedResponse: "list.filter.output.golden",
			Fixture:          "list.filter.api.response.golden",
			URL:              "/api/filter",
			ExpectError:      false,
			HttpMethod:       "GET",
			HttpStatusCode:   200,
		},
	}
	// Run tests
	utils.InvokeMesheryctlTestCommand(t, update, FilterCmd, tests, currDir, "filter")
}

// TestListCmdKeepsFilterNamesOutsideTheExtension pins the NAME column of the
// non-verbose local-provider listing against `strings.Trim`, whose second
// argument is a cutset and not a suffix: every character of the extension was
// also stripped off the front of the name, so `sample.wasm` listed as `ple`.
// The shared golden fixtures authenticate as a remote provider, which takes the
// other branch, so this drives the command with a local-provider token.
func TestListCmdKeepsFilterNamesOutsideTheExtension(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	testContext := utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(FilterCmd, t)

	apiResponse := utils.NewGoldenFile(t, "list.filter.local.api.response.golden", filepath.Join(currDir, "fixtures")).Load()
	httpmock.RegisterResponder("GET", testContext.BaseURL+"/api/filter?pagesize=25&page=0",
		httpmock.NewStringResponder(200, apiResponse))

	previousToken := utils.TokenFlag
	utils.TokenFlag = filepath.Join(currDir, "fixtures", "local.token.golden")
	t.Cleanup(func() { utils.TokenFlag = previousToken })

	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	// The table is printed to os.Stdout directly by utils.PrintToTable, not to
	// the command's own writer, so it can only be captured by swapping stdout.
	os.Stdout = w
	t.Cleanup(func() {
		os.Stdout = originalStdout
		_ = r.Close()
	})

	// --page skips the interactive pagination prompt and prints the table once.
	FilterCmd.SetArgs([]string{"list", "--page", "1"})
	execErr := FilterCmd.Execute()

	// io.ReadAll below returns only once the write end is closed.
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close pipe: %v", err)
	}

	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read command output: %v", err)
	}

	assert.NoError(t, execErr)
	assert.Contains(t, string(output), "mesh-filter")
	assert.Contains(t, string(output), "sample")
}
