package perf

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResultViewPrintsLatenciesInMilliseconds runs `perf result --view` against a
// result as fortio reports it, in seconds, and checks the line the user reads.
func TestResultViewPrintsLatenciesInMilliseconds(t *testing.T) {
	output := runResultViewWithPercentiles(t, `[
		{"Percentile": 50, "Value": 0.0015},
		{"Percentile": 75, "Value": 0.002},
		{"Percentile": 90, "Value": 0.003},
		{"Percentile": 99, "Value": 0.008},
		{"Percentile": 99.9, "Value": 0.0095}
	]`)
	assert.Contains(t, output, "Latencies _ms: Avg: 2, Max: 10, Min: 1, P50: 1.5, P90: 3, P99: 8")
}

// TestResultViewMatchesPercentilesByField runs `perf result --view` against a
// result whose percentile list is out of order and carries an extra 99.99
// entry, and checks the percentiles are still matched by their Percentile
// field rather than by slice position.
func TestResultViewMatchesPercentilesByField(t *testing.T) {
	output := runResultViewWithPercentiles(t, `[
		{"Percentile": 99, "Value": 0.008},
		{"Percentile": 99.99, "Value": 0.0099},
		{"Percentile": 50, "Value": 0.0015},
		{"Percentile": 90, "Value": 0.003},
		{"Percentile": 99.9, "Value": 0.0095}
	]`)
	assert.Contains(t, output, "Latencies _ms: Avg: 2, Max: 10, Min: 1, P50: 1.5, P90: 3, P99: 8")
}

// runResultViewWithPercentiles runs `perf result --view` against a single
// mocked result whose DurationHistogram carries the given percentile list,
// and returns the captured stdout.
func runResultViewWithPercentiles(t *testing.T, percentilesJSON string) string {
	t.Helper()
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok, "not able to get current working directory")
	origToken := utils.TokenFlag
	utils.TokenFlag = filepath.Join(filepath.Dir(filename), "fixtures", "auth.json")
	defer func() { utils.TokenFlag = origToken }()

	const profileID = "906f8876-33b5-4a97-906e-7a409d3b8ae9"
	httpmock.RegisterResponder(http.MethodGet, testContext.BaseURL+"/api/performance/profiles",
		httpmock.NewStringResponder(http.StatusOK,
			`{"profiles":[{"id":"`+profileID+`","name":"my-profile","endpoints":["https://meshery.io"]}],"totalCount":1}`))
	httpmock.RegisterResponder(http.MethodGet, testContext.BaseURL+"/api/user/performance/profiles/"+profileID+"/results",
		httpmock.NewStringResponder(http.StatusOK, `{"results":[{
			"name": "my-result",
			"testStartTime": "2026-11-06T10:00:00Z",
			"runnerResults": {
				"DurationHistogram": {
					"Avg": 0.002, "Max": 0.01, "Min": 0.001,
					"Percentiles": `+percentilesJSON+`
				}
			}
		}]}`))

	resetVariables()
	defer resetVariables()
	viewSingleResult = true
	defer func() { viewSingleResult = false }()

	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	stdout := os.Stdout
	os.Stdout = writer
	defer func() { os.Stdout = stdout }()

	runErr := resultCmd.RunE(resultCmd, []string{"my-profile"})

	require.NoError(t, writer.Close())
	os.Stdout = stdout
	output, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.NoError(t, runErr)

	return string(output)
}
