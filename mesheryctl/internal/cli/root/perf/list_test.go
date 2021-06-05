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
		Fetch            string
		ExpectedResponse string
		Fixture          string
		URL              string
		Token            string
		ExpectError      bool
	}{
		{
			Name:             "Fetch Profiles",
			Fetch:            "Profiles",
			ExpectedResponse: "profile.output.golden",
			Fixture:          "profile.api.response.golden",
			URL:              testContext.BaseURL + "/api/user/performance/profiles",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "Fetch Results",
			Fetch:            "Results",
			ExpectedResponse: "result.output.golden",
			Fixture:          "result.api.response.golden",
			URL:              testContext.BaseURL + "/api/user/performance/profiles/ecddef09-7411-4b9e-b06c-fd55ff5debbc/results",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "Fetch Profiles with no token",
			Fetch:            "Profiles",
			ExpectedResponse: "profile.no.token.golden",
			Fixture:          "profile.api.response.golden",
			URL:              testContext.BaseURL + "/api/user/performance/profiles",
			Token:            "",
			ExpectError:      true,
		},
		{
			Name:             "Fetch Results with No token",
			Fetch:            "Results",
			ExpectedResponse: "result.no.token.golden",
			Fixture:          "result.api.response.golden",
			URL:              testContext.BaseURL + "/api/user/performance/profiles/ecddef09-7411-4b9e-b06c-fd55ff5debbc/results",
			Token:            "",
			ExpectError:      true,
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

			var data [][]string
			var err error
			if tt.Fetch == "Profiles" {
				data, _, err = fetchPerformanceProfiles(tt.URL)
			} else {
				data, _, err = fetchPerformanceProfileResults(tt.URL, "ecddef09-7411-4b9e-b06c-fd55ff5debbc")
			}

			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// if we're supposed to get an error
			if tt.ExpectError {
				// write it in file
				if *update {
					golden.Write(err.Error())
				}
				expectedResponse := golden.Load()

				utils.Equals(t, expectedResponse, err.Error())
				return
			}
			// when we're supposed to fetch data
			var actualResponse string
			if tt.Fetch == "Profiles" {
				actualResponse = utils.PrintToTableInStringFormat([]string{"Name", "ID", "RESULTS", "LAST-RUN"}, data)
			} else {
				actualResponse = utils.PrintToTableInStringFormat([]string{"NAME", "MESH", "START-TIME", "QPS", "DURATION", "P50", "P99.9"}, data)
			}

			// write it in file
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			utils.Equals(t, expectedResponse, actualResponse)
		})
	}

	// stop mock server
	utils.StopMockery(t)
}
