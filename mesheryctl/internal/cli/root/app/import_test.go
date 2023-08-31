package app

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestImportCmd(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	// initialize mock server for handling requests
	utils.StartMockery(t)

	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Unable to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")
	tests := []struct {
		Name             string
		Args             []string
		ExpectedResponse string
		URL              string
		Token            string
		Method           string
		Fixture          string
		ExpectedError    bool
	}{
		{
			Name:             "Test Import file",
			Args:             []string{"import", "-f", filepath.Join(fixturesDir, "importapp.golden"), "-s", "Docker Compose"},
			URL:              testContext.BaseURL + "/api/application",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			Method:           "GET",
			Fixture:          "apply.patternSave.response.golden",
			ExpectedError:    false,
			ExpectedResponse: "import.output.golden",
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
			httpmock.RegisterResponder(tt.Method, tt.URL,
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
				if tt.ExpectedError {
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
}
