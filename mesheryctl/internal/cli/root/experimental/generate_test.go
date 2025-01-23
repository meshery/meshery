package experimental

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestExperimentalGenerate(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	//initialize mock server for handling requests
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

	// test scenarios for fetching data
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
			Name:             "Generate registered relationships",
			Args:             []string{"relationship", "generate", "$CRED", "", "-s", "1"},
			URL:              testContext.BaseURL + "/api/meshmodels/relationships",
			Fixture:          "generate.exp.relationship.api.response.golden",
			ExpectedResponse: "generate.exp.relationship.output.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

			utils.TokenFlag = tt.Token

			httpmock.RegisterResponder("GET", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Grab console prints
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			ExpCmd.SetArgs(tt.Args)
			ExpCmd.SetOutput(rescueStdout)
			err := ExpCmd.Execute()
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

			actualResponse := string(out)

			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			cleanedActualResponse := utils.CleanStringFromHandlePagination(actualResponse)
			cleanedExceptedResponse := utils.CleanStringFromHandlePagination(expectedResponse)

			utils.Equals(t, cleanedExceptedResponse, cleanedActualResponse)
		})
		t.Log("Generate experimental relationship test passed")
	}

	utils.StopMockery(t)
}

func TestExperimentalGenerate_MissingArguments(t *testing.T) {

	const errMsg = "Usage: mesheryctl exp relationship generate $CRED [google-sheets-credential] --sheetId [sheet-id]\nRun 'mesheryctl exp relationship generate --help' to see detailed help message"

	// test scenarios for fetching data
	tests := []struct {
		Name             string
		Args             []string
		ExpectedResponse string
		ExpectError      bool
	}{
		{
			Name:             "Missing Credentials",
			Args:             []string{"relationship", "generate"},
			ExpectedResponse: "Google Sheet Credentials is required\n" + errMsg,
			ExpectError:      false,
		},
		{
			Name:             "Missing Sheet ID",
			Args:             []string{"relationship", "generate", "$CRED", "--sheetId", ""},
			ExpectedResponse: "Sheet ID is required\n" + errMsg,
			ExpectError:      false,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			cmd := ExpCmd
			cmd.SetArgs(tt.Args)
			err := cmd.Execute()

			utils.Equals(t, errors.New(utils.RelationshipsError(tt.ExpectedResponse, "generate")), err)
		})
		t.Log("Generate relationships test passed")
	}
}
