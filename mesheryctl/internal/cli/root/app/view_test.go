package app

import (
	"bytes"
	"flag"

	// "io/ioutil"
	"log"
	// "os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")
var tempAppID = "e39df138-bd73-47f1-8db4-edfc4027f178"

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
			Name:             "View Applications",
			Args:             []string{"view", "--all"},
			View:             "Applications",
			ExpectedResponse: "view.application.output.golden",
			Fixture:          "view.application.api.response.golden",
			URL:              testContext.BaseURL + "/api/experimental/application?page_size=10000",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		// {
		// 	Name:             "View Applications with ID",
		// 	Args:             []string{"view", tempAppID},
		// 	View:             "Results",
		// 	ExpectedResponse: "view.result.output.golden",
		// 	Fixture:          "view.result.api.response.golden",
		// 	URL:              testContext.BaseURL + "/api/experimental/application/" + tempAppID,
		// 	Token:            filepath.Join(fixturesDir, "token.golden"),
		// 	ExpectError:      false,
		// },
		// {
		// 	Name:             "View Applications with no token",
		// 	Args:             []string{"view", "test"},
		// 	View:             "Applications",
		// 	ExpectedResponse: "no.token.golden",
		// 	Fixture:          "view.application.api.response.golden",
		// 	URL:              testContext.BaseURL + "/api/experimental/application?page_size=10000&search=test",
		// 	Token:            "",
		// 	ExpectError:      true,
		// },
		// {
		// 	Name:             "View Results with No token",
		// 	Args:             []string{"view", tempAppID, "app%20mesh"},
		// 	View:             "Results",
		// 	ExpectedResponse: "no.token.golden",
		// 	Fixture:          "view.result.api.response.golden",
		// 	URL:              testContext.BaseURL + "/api/experimental/application/" + tempAppID + "/results?pageSize=25&search=app%20mesh",
		// 	Token:            "",
		// 	ExpectError:      true,
		// },
		// {
		// 	Name:             "View Invalid Application Name",
		// 	Args:             []string{"view", "invalid-name", "-t", filepath.Join(fixturesDir, "token.golden")},
		// 	View:             "Applications",
		// 	ExpectedResponse: "view.invalid.application.output.golden",
		// 	Fixture:          "view.invalid.application.api.response.golden",
		// 	URL:              testContext.BaseURL + "/api/experimental/application?page_size=25&search=invalid-name",
		// 	Token:            filepath.Join(fixturesDir, "token.golden"),
		// 	ExpectError:      true,
		// },
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// View api response from golden files
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

			// set token
			tokenPath = tt.Token

			// mock response
			httpmock.RegisterResponder("GET", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			//Grab Logs
			var buf bytes.Buffer
			log.SetOutput(&buf)
			utils.SetupLogrusFormatter()

			// Grab console prints
			// rescueStdout := os.Stdout
			// r, w, _ := os.Pipe()
			// os.Stdout = w

			AppCmd.SetArgs(tt.Args)
			// AppCmd.SetOutput(rescueStdout)
			err := AppCmd.Execute()
			t.Log("Executed")
			if err != nil {
				log.Print("We got errorrrrrrrrrrr")
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

			// w.Close()
			// log.Print("We are safeeeeeeeeeeeeeeeeee")
			// out, _ := ioutil.ReadAll(r)
			// os.Stdout = rescueStdout

			// // response being printed in console
			// actualResponse := string(out)
			actualResponse := buf.String()

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
