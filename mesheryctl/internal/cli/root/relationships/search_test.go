package relationships

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestSearch(t *testing.T) {
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
			Name:             "Search with missing arguments",
			Args:             []string{"search"},
			URL:              "",
			Fixture:          "",
			ExpectedResponse: "search.missing.args.output.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      true,
		},
		{
			Name:             "Search registered relationships",
			Args:             []string{"search", "--model", "kubernetes"},
			URL:              testContext.BaseURL + "/api/meshmodels/models/kubernetes",
			Fixture:          "search.relationship.api.response.golden",
			ExpectedResponse: "search.relationship.output.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Fixture != "" {
				apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()
				httpmock.RegisterResponder("GET", tt.URL,
					httpmock.NewStringResponder(200, apiResponse))
			}

			utils.TokenFlag = tt.Token

			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Grab console prints
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			RelationshipCmd.SetArgs(tt.Args)
			RelationshipCmd.SetOutput(rescueStdout)
			err := RelationshipCmd.Execute()
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
		t.Log("Search experimental relationship test passed")
	}

	utils.StopMockery(t)
}
