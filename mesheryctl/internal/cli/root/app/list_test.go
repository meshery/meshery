package app

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestAppList(t *testing.T) {
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
			Name:             "Fetch Applications List",
			Args:             []string{"list"},
			ExpectedResponse: "list.application.output.golden",
			Fixture:          "list.application.api.response.golden",
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
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			// Grab console prints
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			AppCmd.SetArgs(tt.Args)
			AppCmd.SetOutput(rescueStdout)
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
				t.Fatal(err)
			}

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			// response being printed in console
			actualResponse := string(out)

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
