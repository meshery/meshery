package perf

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
)

var update = flag.Bool("update", false, "update golden files")

func TestFetchList(t *testing.T) {
	utils.SetupContextEnv(t)
	// Get viper instance used for context
	mctlCfg, _ := config.GetMesheryCtl(viper.GetViper())
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("problems recovering caller information")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	tests := []struct {
		Name             string
		ExpectedResponse string
		Fixture          string
		URL              string
		Token            string
	}{
		{
			Name:             "Fetch Profiles",
			ExpectedResponse: "profile.output.golden",
			Fixture:          "profile.api.response.golden",
			URL:              mctlCfg.GetBaseMesheryURL() + "/api/user/performance/profiles",
			Token:            filepath.Join(fixturesDir, "token.golden"),
		},
		{
			Name:             "Fetch Results",
			ExpectedResponse: "result.output.golden",
			Fixture:          "result.api.response.golden",
			URL:              mctlCfg.GetBaseMesheryURL() + "/api/user/performance/profiles/ecddef09-7411-4b9e-b06c-fd55ff5debbc/results",
			Token:            filepath.Join(fixturesDir, "token.golden"),
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Fetch api response from golden files
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()
			// Start a local HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test request parameters
				utils.Equals(t, req.URL.String(), tt.URL)
				// Send response to be tested
				_, err := rw.Write([]byte(apiResponse))
				if err != nil {
					t.Fatal(err)
				}
			}))
			// t.Log(string(apiResponse))

			// Close the server when test finishes
			defer server.Close()
			// set token
			tokenPath = tt.Token

			var actualResponse []byte
			var err error
			if tt.Name == "Fetch Profiles" {
				_, actualResponse, err = fetchPerformanceProfiles(server.Client(), tt.URL)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				t.Log(tt.URL)
				_, actualResponse, err = fetchPerformanceProfileResults(server.Client(), tt.URL, "ecddef09-7411-4b9e-b06c-fd55ff5debbc")
				if err != nil {
					t.Fatal(err)
				}
				// t.Log(string(apiResponse))
			}

			testdataDir := filepath.Join(currDir, "testdata")
			// testdataDir := filepath.Join(filepath.Dir(filename), tf.dir, tf.name)
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			if *update {
				golden.Write(string(actualResponse))
			}
			// t.Log(string(actualResponse))
			expectedResponse := golden.LoadByte()

			utils.Equals(t, string(expectedResponse), string(actualResponse))
		})
	}
}
