package perf

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")

func TestFetchList(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	// initialize mock server for handling requests
	utils.StartMockery(t)

	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("problems recovering caller information")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	// test scenrios for fetching data
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
			URL:              testContext.BaseURL + "/api/user/performance/profiles",
			Token:            filepath.Join(fixturesDir, "token.golden"),
		},
		{
			Name:             "Fetch Results",
			ExpectedResponse: "result.output.golden",
			Fixture:          "result.api.response.golden",
			URL:              testContext.BaseURL + "/api/user/performance/profiles/ecddef09-7411-4b9e-b06c-fd55ff5debbc/results",
			Token:            filepath.Join(fixturesDir, "token.golden"),
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Fetch api response from golden files
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

			// set token
			tokenPath = tt.Token

			// mock response
			httpmock.RegisterResponder("GET", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			var actualResponse []byte
			var err error
			if tt.Name == "Fetch Profiles" {
				_, actualResponse, err = fetchPerformanceProfiles(tt.URL)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				_, actualResponse, err = fetchPerformanceProfileResults(tt.URL, "ecddef09-7411-4b9e-b06c-fd55ff5debbc")
				if err != nil {
					t.Fatal(err)
				}
			}

			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			if *update {
				golden.Write(string(actualResponse))
			}
			expectedResponse := golden.LoadByte()

			utils.Equals(t, string(expectedResponse), string(actualResponse))
		})
	}
	utils.StopMockery(t)
}
