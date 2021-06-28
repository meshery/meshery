package perf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestPerfView(t *testing.T) {
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
		Args             []string
		View             string
		ExpectedResponse string
		Fixture          string
		URL              string
		Token            string
		ExpectError      bool
	}{
		{
			Name:             "View Profiles",
			Args:             []string{"view", "test", "--token", filepath.Join(fixturesDir, "token.golden")},
			View:             "Profiles",
			ExpectedResponse: "view.profile.output.golden",
			Fixture:          "view.profile.api.response.golden",
			URL:              testContext.BaseURL + "/api/user/performance/profiles?page_size=25&search=test",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "View Results",
			Args:             []string{"view", "ecddef09-7411-4b9e-b06c-fd55ff5debbc", "app%20mesh"},
			View:             "Results",
			ExpectedResponse: "view.result.output.golden",
			Fixture:          "view.result.api.response.golden",
			URL:              testContext.BaseURL + "/api/user/performance/profiles/ecddef09-7411-4b9e-b06c-fd55ff5debbc/results?pageSize=25&search=app%20mesh",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "View Profiles with no token",
			Args:             []string{"view", "test"},
			View:             "Profiles",
			ExpectedResponse: "no.token.golden",
			Fixture:          "view.profile.api.response.golden",
			URL:              testContext.BaseURL + "/api/user/performance/profiles?page_size=25&search=test",
			Token:            "",
			ExpectError:      true,
		},
		{
			Name:             "View Results with No token",
			Args:             []string{"view", "ecddef09-7411-4b9e-b06c-fd55ff5debbc", "app%20mesh"},
			View:             "Results",
			ExpectedResponse: "no.token.golden",
			Fixture:          "view.result.api.response.golden",
			URL:              testContext.BaseURL + "/api/user/performance/profiles/ecddef09-7411-4b9e-b06c-fd55ff5debbc/results?pageSize=25&search=app%20mesh",
			Token:            "",
			ExpectError:      true,
		},
		{
			Name:             "View Invalid Profile Name",
			Args:             []string{"view", "invalid-name", "-t", filepath.Join(fixturesDir, "token.golden")},
			View:             "Profiles",
			ExpectedResponse: "view.invalid.profile.output.golden",
			Fixture:          "view.invalid.profile.api.response.golden",
			URL:              testContext.BaseURL + "/api/user/performance/profiles?page_size=25&search=invalid-name",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      true,
		},
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

			// Grab console prints
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PerfCmd.SetArgs(tt.Args)
			PerfCmd.SetOutput(rescueStdout)
			err := PerfCmd.Execute()
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

			w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = rescueStdout

			// response being printed in console
			actualResponse := string(out)

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
