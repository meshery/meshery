package perf

import (
	"bytes"
	"path/filepath"
	"runtime"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestApplyCmd(t *testing.T) {
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
		ExpectedResponse string
		URLs             []utils.MockURL
		Token            string
		ExpectError      bool
	}{
		{
			Name:             "Run Test with Existing profile",
			Args:             []string{"apply", "new", "--token", filepath.Join(fixturesDir, "token.golden")},
			ExpectedResponse: "apply.success.output.golden",
			URLs: []utils.MockURL{
				{
					Method:   "GET",
					URL:      testContext.BaseURL + "/api/user/performance/profiles?search=new",
					Response: "apply.fetch.profile.response.golden",
				},
				{
					Method:   "GET",
					URL:      testContext.BaseURL + "/api/user/performance/profiles/" + "8f3daf25-e58e-4c59-8bf8-f474b76463ec" + "/run",
					Response: "apply.run.existing.perf.test.output.golden",
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "Run Test with Existing profile with new URL",
			Args:             []string{"apply", "new", "--url", "https://www.google.com", "--token", filepath.Join(fixturesDir, "token.golden")},
			ExpectedResponse: "apply.success.output.golden",
			URLs: []utils.MockURL{
				{
					Method:   "GET",
					URL:      testContext.BaseURL + "/api/user/performance/profiles?search=new",
					Response: "apply.fetch.profile.response.golden",
				},
				{
					Method:   "GET",
					URL:      testContext.BaseURL + "/api/user/performance/profiles/" + "8f3daf25-e58e-4c59-8bf8-f474b76463ec" + "/run",
					Response: "apply.run.new.perf.test.output.golden",
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "Run Test with new profile",
			Args:             []string{"apply", "--profile", "test", "--url", "https://www.google.com", "--token", filepath.Join(fixturesDir, "token.golden")},
			ExpectedResponse: "apply.success.output.golden",
			URLs: []utils.MockURL{
				{
					Method:   "POST",
					URL:      testContext.BaseURL + "/api/user/performance/profiles",
					Response: "apply.create.profile.response.golden",
				},
				{
					Method:   "GET",
					URL:      testContext.BaseURL + "/api/user/performance/profiles/" + "906f8876-33b5-4a97-906e-7a409d3b8ae9" + "/run",
					Response: "apply.run.new.profile.test.output.golden",
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "Run Test with new profile but without URL",
			Args:             []string{"apply", "--profile", "test", "--url", "", "--token", filepath.Join(fixturesDir, "token.golden")},
			ExpectedResponse: "apply.no.url.output.golden",
			URLs: []utils.MockURL{
				{
					Method:   "POST",
					URL:      testContext.BaseURL + "/api/user/performance/profiles",
					Response: "apply.create.profile.response.golden",
				},
				{
					Method:   "GET",
					URL:      testContext.BaseURL + "/api/user/performance/profiles/" + "906f8876-33b5-4a97-906e-7a409d3b8ae9" + "/run",
					Response: "apply.run.new.profile.test.output.golden",
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: true,
		},
		{
			Name:             "Run Test with new profile with Invalid URL",
			Args:             []string{"apply", "--profile", "test", "--url", "invalid-url", "--token", filepath.Join(fixturesDir, "token.golden")},
			ExpectedResponse: "apply.invalid.url.output.golden",
			URLs: []utils.MockURL{
				{
					Method:   "POST",
					URL:      testContext.BaseURL + "/api/user/performance/profiles",
					Response: "apply.create.profile.response.golden",
				},
				{
					Method:   "GET",
					URL:      testContext.BaseURL + "/api/user/performance/profiles/" + "906f8876-33b5-4a97-906e-7a409d3b8ae9" + "/run",
					Response: "apply.run.new.profile.test.output.golden",
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: true,
		},
		{
			Name:             "Run Test with Existing profile with Invalid URL",
			Args:             []string{"apply", "new", "--url", "invalid-url", "--token", filepath.Join(fixturesDir, "token.golden")},
			ExpectedResponse: "apply.invalid.url.output.golden",
			URLs: []utils.MockURL{
				{
					Method:   "GET",
					URL:      testContext.BaseURL + "/api/user/performance/profiles?search=new",
					Response: "apply.fetch.profile.response.golden",
				},
				{
					Method:   "GET",
					URL:      testContext.BaseURL + "/api/user/performance/profiles/" + "8f3daf25-e58e-4c59-8bf8-f474b76463ec" + "/run",
					Response: "apply.run.new.perf.test.output.golden",
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			for _, url := range tt.URLs {
				// View api response from golden files
				apiResponse := utils.NewGoldenFile(t, url.Response, fixturesDir).Load()

				// mock response
				httpmock.RegisterResponder(url.Method, url.URL,
					httpmock.NewStringResponder(200, apiResponse))
			}

			// set token
			tokenPath = tt.Token

			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// setting up log to grab logs
			var buf bytes.Buffer
			log.SetOutput(&buf)
			utils.SetupLogrusFormatter()

			PerfCmd.SetArgs(tt.Args)
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
				} else {
					t.Error(err)
				}
			}

			// response being printed in console
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
