package filter

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestViewCmd(t *testing.T) {
	// setup current context
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

	// test scenrios for fetching data
	tests := []struct {
		Name             string
		Args             []string
		URL              string
		Fixture          string
		Token            string
		ExpectedResponse string
		ExpectError      bool
	}{

		{
			Name:             "Fetch Filter View",
			Args:             []string{"view", "KumaTest"},
			ExpectedResponse: "view.filter.output.golden",
			Fixture:          "view.filter.api.response.golden",
			URL:              testContext.BaseURL + "/api/filter",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "Fetch Kuma Filter View with ID",
			Args:             []string{"view", "957fbc9b-a655-4892-823d-375102a9587c"},
			ExpectedResponse: "view.id.filter.output.golden",
			Fixture:          "view.id.filter.api.response.golden",
			URL:              testContext.BaseURL + "/api/filter/957fbc9b-a655-4892-823d-375102a9587c",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "Fetch Filter View for non existing filter",
			Args:             []string{"view", "xyz"},
			ExpectedResponse: "view.nonexisting.filter.output.golden",
			Fixture:          "view.nonexisting.filter.api.response.golden",
			URL:              testContext.BaseURL + "/api/filter",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// View api response from golden files
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

			// set token
			utils.TokenFlag = tt.Token

			// mock response
			httpmock.RegisterResponder("GET", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			b := utils.SetupMeshkitLoggerTesting(t, false)
			FilterCmd.SetOutput(b)

			FilterCmd.SetArgs(tt.Args)
			err := FilterCmd.Execute()
			if err != nil {
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
				t.Fatal(err)
			}

			// response being printed in console
			output := b.String()
			actualResponse := output

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
	t.Log("test Passed")
}
