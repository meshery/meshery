package app

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestAppCmd(t *testing.T) {
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
		ExpectedResponse string
		URLs             []utils.MockURL
		Token            string
		ExpectError      bool
	}{
		{
			Name:             "App list Test",
			Args:             []string{"list"},
			ExpectedResponse: "applist.output.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/application",
					Response:     "app.list.api.response.golden",
					ResponseCode: 200,
				},
			},
		},
		{
			Name:             "App Offboard Test",
			Args:             []string{"offboard", "-f", filepath.Join(fixturesDir, "sampleApp.golden")},
			ExpectedResponse: "app.offborad.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern",
					Response:     "apply.patternSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/application/deploy",
					Response:     "app.offboard.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		// {
		// 	Name:             "App View All Test",
		// 	Args:             []string{"view", "--all"},
		// 	ExpectedResponse: "app.viewall.output.golden",
		// 	URLs: []utils.MockURL{
		// 		{
		// 			Method:       "GET",
		// 			URL:          testContext.BaseURL + "/api/application",
		// 			Response:     "app.viewall.api.response.golden",
		// 			ResponseCode: 200,
		// 		},
		// 	},
		// 	Token:       filepath.Join(fixturesDir, "token.golden"),
		// 	ExpectError: false,
		// },
		{
			Name:             "App view test with ID",
			Args:             []string{"view", "41f81b3d-64dc-42de-8b68-fa904ec46da9"},
			ExpectedResponse: "app.viewid.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/application/41f81b3d-64dc-42de-8b68-fa904ec46da9",
					Response:     "app.viewid.api.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "App Invalid view ID",
			Args:             []string{"view", "bbf81b3d-64dc-42de-8b68-uy904ec46da9"},
			ExpectedResponse: "app.invalid.view.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/application/bbf81b3d-64dc-42de-8b68-uy904ec46da9",
					Response:     "app.view.invalidid.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: true,
		},
	}
	// Run tests
	for _, tt := range tests {
		// View api response from golden files
		t.Run(tt.Name, func(t *testing.T) {
			for _, url := range tt.URLs {
				apiResponse := utils.NewGoldenFile(t, url.Response, fixturesDir).Load()

				// mock response
				httpmock.RegisterResponder(url.Method, url.URL,
					httpmock.NewStringResponder(url.ResponseCode, apiResponse))
			}

			// set token
			utils.TokenFlag = tt.Token

			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			b := utils.SetupMeshkitLoggerTesting(t, false)
			AppCmd.SetOutput(b)
			AppCmd.SetArgs(tt.Args)
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

			assert.Equal(t, expectedResponse, actualResponse)
		})
	}

	// stop mock server
	utils.StopMockery(t)
}
