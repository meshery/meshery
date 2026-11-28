package adapter

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
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

// adapterLocationFromFixture returns the location of the named adapter in a
// sync fixture.
func adapterLocationFromFixture(t *testing.T, path, name string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	var prefs struct {
		MeshAdapters []struct {
			Location string `json:"adapterLocation"`
			Name     string `json:"name"`
		} `json:"meshAdapters"`
	}
	require.NoError(t, json.Unmarshal(data, &prefs))
	for _, a := range prefs.MeshAdapters {
		if a.Name == name {
			require.NotEmpty(t, a.Location, "adapter %s has no location in %s", name, path)
			return a.Location
		}
	}
	t.Fatalf("adapter %s not found in %s", name, path)
	return ""
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

	// The validate request must carry the location of the ISTIO adapter from
	// the sync fixture. Read it from the fixture so the test does not depend
	// on the fixture's order or contents.
	expectedAdapter := adapterLocationFromFixture(t, filepath.Join(fixturesDir, "sync.golden"), "ISTIO")
	// sync.compose.golden lists the adapters under their Docker Compose
	// hostnames, so --adapter meshery-istio resolves without a positional mesh.
	composeIstioAdapter := adapterLocationFromFixture(t, filepath.Join(fixturesDir, "sync.compose.golden"), "ISTIO")

	syncURL := testContext.BaseURL + "/api/system/sync"
	operationURL := testContext.BaseURL + "/api/system/adapter/operation"

	tests := []struct {
		Name             string
		Args             []string
		URLs             []utils.MockURL
		ExpectedResponse string
		ExpectError      bool
		// query expected in the adapter operation POST body
		ExpectedQuery string
		// adapter expected in the adapter operation POST body; defaults to the
		// ISTIO location in sync.golden
		ExpectedAdapter string
		// number of GET requests expected on the sync endpoint; defaults to 1
		ExpectedSyncCalls int
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
			ExpectedQuery:          "smi_conformance",
			ExpectedOperationCalls: 1,
		},
		{
			Name: "Validate without the required spec flag",
			Args: []string{"validate", "istio", "--adapter", "meshery-istio"},
			// Cobra checks required flags after PreRunE, so the adapter sync
			// still runs; the validation request itself must never be sent.
			// The operation responder is registered anyway so that an
			// unexpected POST would be counted instead of going unmatched.
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
			ExpectedQuery:          "smi_conformance",
			ExpectedOperationCalls: 1,
		},
		{
			// The Istio adapter in sync.golden is at localhost:10000, not the
			// Docker Compose hostname meshery-istio:10000, so this guards against
			// istio-vet being tied to one deployment's adapter location.
			Name: "Validate Istio with the istio-vet spec at a non-default location",
			Args: []string{"validate", "istio", "--adapter", "meshery-istio", "--spec", "istio-vet"},
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
			ExpectedQuery:          "istio-vet",
			ExpectedOperationCalls: 1,
		},
		{
			Name: "Reject the istio-vet spec for a mesh other than Istio",
			Args: []string{"validate", "linkerd", "--adapter", "meshery-linkerd", "--spec", "istio-vet"},
			// The operation responder is registered so that a request that
			// should have been rejected would be counted.
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
			ExpectedResponse:       "validate.istio-vet-non-istio.output.golden",
			ExpectError:            true,
			ExpectedOperationCalls: 0,
		},
		{
			// With no positional mesh, validate's PreRunE derives meshName from
			// the adapter whose host matches --adapter, so istio-vet must work on
			// that path too. This path reads the session data twice: once to
			// resolve --adapter and once more in validateAdapter.
			Name: "Validate with the istio-vet spec when the mesh comes from --adapter",
			Args: []string{"validate", "--adapter", "meshery-istio", "--spec", "istio-vet"},
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          syncURL,
					Response:     "sync.compose.golden",
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
			ExpectedQuery:          "istio-vet",
			ExpectedAdapter:        composeIstioAdapter,
			ExpectedSyncCalls:      2,
			ExpectedOperationCalls: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			resetValidateFlags(t)
			// meshName is package state that the no-positional-mesh path only
			// fills when empty, as it is in a fresh mesheryctl process; clear it
			// so one case's mesh cannot leak into the next.
			meshName = ""
			wantAdapter := expectedAdapter
			if tc.ExpectedAdapter != "" {
				wantAdapter = tc.ExpectedAdapter
			}
			wantSyncCalls := 1
			if tc.ExpectedSyncCalls != 0 {
				wantSyncCalls = tc.ExpectedSyncCalls
			}
			httpmock.Reset()

			for _, mock := range tc.URLs {
				apiResponse := utils.NewGoldenFile(t, mock.Response, fixturesDir).Load()
				// mock response; for the operation POST, also check what was sent.
				// The request is form-encoded, so parse the raw body and compare
				// complete values rather than substrings.
				httpmock.RegisterResponder(mock.Method, mock.URL,
					func(req *http.Request) (*http.Response, error) {
						if mock.Method == http.MethodPost {
							// This runs inside the HTTP transport, not the test
							// goroutine, so use assert (not require, which calls
							// t.FailNow) and return the error to the client.
							body, err := io.ReadAll(req.Body)
							if !assert.NoError(t, err) {
								return nil, err
							}
							form, err := url.ParseQuery(string(body))
							if !assert.NoError(t, err) {
								return nil, err
							}
							assert.Equal(t, tc.ExpectedQuery, form.Get("query"))
							assert.Equal(t, wantAdapter, form.Get("adapter"))
						}
						return httpmock.NewStringResponse(mock.ResponseCode, apiResponse), nil
					})
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
			// Cobra checks required flags after PreRunE, so the sync request
			// runs in every case, including the missing-spec one.
			assert.Equal(t, wantSyncCalls, httpmock.GetCallCountInfo()["GET "+syncURL],
				"unexpected number of adapter sync requests")

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
