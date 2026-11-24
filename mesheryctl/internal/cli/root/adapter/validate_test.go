package adapter

import (
	"net/http"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resetValidateFlags restores every validate flag to its default and clears
// its "changed" state. AdapterCmd is shared across test cases, so without this
// a flag set by one case (for example --spec) would still count as set in the
// next one.
func resetValidateFlags(t *testing.T) {
	t.Helper()
	validateCmd.Flags().VisitAll(func(f *pflag.Flag) {
		require.NoError(t, f.Value.Set(f.DefValue))
		f.Changed = false
	})
}

func TestValidateMesh(t *testing.T) {
	utils.SetupContextEnv(t)

	// initialize mock server for handling requests
	utils.StartMockery(t)

	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")
	testdataDir := filepath.Join(currDir, "testdata")

	syncURL := testContext.BaseURL + "/api/system/sync"
	operationURL := testContext.BaseURL + "/api/system/adapter/operation"

	tests := []struct {
		Name             string
		Args             []string
		URLs             []utils.MockURL
		ExpectedResponse string
		ExpectError      bool
		// number of POST requests expected on the adapter operation endpoint
		ExpectedOperationCalls int
	}{
		{
			Name: "Validate Istio with the smi spec",
			Args: []string{"validate", "istio", "--adapter", "meshery-istio", "--spec", "smi"},
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          syncURL,
					Response:     "sync.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          operationURL,
					Response:     "deploy.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse:       "validate.istio.output.golden",
			ExpectedOperationCalls: 1,
		},
		{
			Name: "Validate without the required spec flag",
			Args: []string{"validate", "istio", "--adapter", "meshery-istio"},
			// Cobra checks required flags after PreRunE, so the adapter sync
			// still runs; the validation request itself must never be sent.
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          syncURL,
					Response:     "sync.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse:       "validate.missing-spec.output.golden",
			ExpectError:            true,
			ExpectedOperationCalls: 0,
		},
		{
			Name: "Validate when the adapter operation endpoint returns an error",
			Args: []string{"validate", "istio", "--adapter", "meshery-istio", "--spec", "smi"},
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          syncURL,
					Response:     "sync.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          operationURL,
					Response:     "deploy.golden",
					ResponseCode: http.StatusInternalServerError,
				},
			},
			ExpectedResponse:       "validate.operation-error.output.golden",
			ExpectError:            true,
			ExpectedOperationCalls: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			resetValidateFlags(t)
			httpmock.Reset()

			for _, url := range tc.URLs {
				apiResponse := utils.NewGoldenFile(t, url.Response, fixturesDir).Load()
				// mock response
				httpmock.RegisterResponder(url.Method, url.URL,
					httpmock.NewStringResponder(url.ResponseCode, apiResponse))
			}
			utils.TokenFlag = filepath.Join(fixturesDir, "token.golden")

			golden := utils.NewGoldenFile(t, tc.ExpectedResponse, testdataDir)
			buff := utils.SetupMeshkitLoggerTesting(t, false)
			cmd := AdapterCmd
			cmd.SetArgs(tc.Args)
			cmd.SetOut(buff)
			err := cmd.Execute()

			assert.Equal(t, tc.ExpectedOperationCalls,
				httpmock.GetCallCountInfo()["POST "+operationURL],
				"unexpected number of adapter operation requests")

			if tc.ExpectError {
				require.Error(t, err)
				if *update {
					golden.Write(err.Error())
				}
				assert.Equal(t, golden.Load(), err.Error())
				return
			}

			require.NoError(t, err)
			actualResponse := buff.String()
			if *update {
				golden.Write(actualResponse)
			}
			assert.Equal(t, golden.Load(), actualResponse)
		})
	}
	utils.StopMockery(t)
}
