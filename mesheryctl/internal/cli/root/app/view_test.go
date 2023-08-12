package app

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")
var tempAppID = "e39df138-bd73-47f1-8db4-edfc4027f178"
var tempAppName = "OApplicationPattern"

func TestAppView(t *testing.T) {
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
		View             string
		ExpectedResponse string
		Fixture          string
		URL              string
		Token            string
		ExpectError      bool
	}{
		{
			Name:             "View Applications with ID",
			Args:             []string{"view", tempAppID},
			View:             "ApplicationId",
			ExpectedResponse: "view.applicationid.output.golden",
			Fixture:          "view.applicationid.api.response.golden",
			URL:              testContext.BaseURL + "/api/application/" + tempAppID,
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "View Applications with Name",
			Args:             []string{"view", tempAppName},
			View:             "ApplicationName",
			ExpectedResponse: "view.applicationName.output.golden",
			Fixture:          "view.applicationName.api.response.golden",
			URL:              testContext.BaseURL + "/api/application?search=" + tempAppName,
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "View Invalid Application Name",
			Args:             []string{"view", "invalid-name"},
			View:             "ApplicationsName",
			ExpectedResponse: "view.invalid.application.output.golden",
			Fixture:          "view.invalid.application.api.response.golden",
			URL:              testContext.BaseURL + "/api/application?search=invalid-name",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      true,
		},
		{
			Name:             "View All Applications",
			Args:             []string{"view", "--all"},
			View:             "Applications",
			ExpectedResponse: "view.application.output.golden",
			Fixture:          "view.application.api.response.golden",
			URL:              testContext.BaseURL + "/api/application",
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

			//Grab Logs
			b := utils.SetupMeshkitLoggerTesting(t, false)

			AppCmd.SetOut(b)
			AppCmd.SetArgs(tt.Args)
			// AppCmd.SetOutput(rescueStdout)
			err := AppCmd.Execute()
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
				t.Error(err)
			}

			// response being printed in console
			actualResponse := b.String()

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
