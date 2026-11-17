package perf

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// --page pages the profile's results. The profile is still looked up on the first
// page of the name search: page 2 of a search that matches one profile is empty.
func TestResultPageFlagPagesResultsNotTheProfileLookup(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	defer resetVariables()
	defer utils.ResetCommandFlags(PerfCmd, t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	fixturesDir := filepath.Join(filepath.Dir(filename), "fixtures")
	utils.TokenFlag = utils.GetToken(t)

	profileURL := testContext.BaseURL + "/api/performance/profiles"
	resultURL := testContext.BaseURL + "/api/user/performance/profiles/" + tempProfileID + "/results"
	profiles := utils.NewGoldenFile(t, "result.profile.response.golden", fixturesDir).Load()
	results := utils.NewGoldenFile(t, "result.list.response.golden", fixturesDir).Load()

	// the one profile named abhishek is on the first page of the search; any other
	// page of it is empty, as the server returns it
	httpmock.RegisterResponderWithQuery("GET", profileURL, "page=0&pagesize=25&search=abhishek",
		httpmock.NewStringResponder(200, profiles))
	httpmock.RegisterResponder("GET", profileURL,
		httpmock.NewStringResponder(200, `{"page":1,"pageSize":25,"totalCount":1,"profiles":[]}`))
	httpmock.RegisterResponderWithQuery("GET", resultURL, "page=1&pagesize=25",
		httpmock.NewStringResponder(200, results))

	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	defer func() { os.Stdout = stdout; _ = r.Close() }()

	PerfCmd.SetArgs([]string{"result", "abhishek", "--page", "2"})
	err := PerfCmd.Execute()
	_ = w.Close()
	if err != nil {
		t.Fatalf("perf result abhishek --page 2: %v", err)
	}

calls := httpmock.GetCallCountInfo()
if n := calls["GET "+profileURL+"?page=0&pagesize=25&search=abhishek"]; n != 1 {
	t.Errorf("expected the profile lookup to hit page 0 once, got %d calls: %v", n, calls)
}
if n := calls["GET "+resultURL+"?page=1&pagesize=25"]; n != 1 {
	t.Errorf("expected page 2 of the results (page=1) to be fetched once, got %d calls: %v", n, calls)
}
}
