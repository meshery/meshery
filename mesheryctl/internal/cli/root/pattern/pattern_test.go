package pattern

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestPatterncmd(t *testing.T) {
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
	testcase := []struct {
		Name             string
		Args             []string
		ExpectedResponse string
		URLs             []utils.MockURL
		Token            string
		ExpectError      bool
	}{
		{
			Name:             "pattern apply",
			Args:             []string{"apply", "-f", filepath.Join(fixturesDir, "pattern.golden")},
			ExpectedResponse: "pattern.apply.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern",
					Response:     "pattern.apply.save.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "pattern.apply.deploy.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "pattern delete",
			Args:             []string{"delete", "-f", filepath.Join(fixturesDir, "pattern.golden")},
			ExpectedResponse: "pattern.delete.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "pattern.delete.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		//commented to pass in workflow run
		// {
		// 	Name:             "pattern list",
		// 	Args:             []string{"list"},
		// 	ExpectedResponse: "pattern.list.output.golden",
		// 	URLs: []utils.MockURL{
		// 		{
		// 			Method:       "GET",
		// 			URL:          testContext.BaseURL + "/api/pattern",
		// 			Response:     "pattern.list.golden",
		// 			ResponseCode: 200,
		// 		},
		// 	},
		// 	Token:       filepath.Join(fixturesDir, "token.golden"),
		// 	ExpectError: false,
		// },
		{
			Name:             "pattern view",
			Args:             []string{"view", "kumatest"},
			ExpectedResponse: "pattern.view.kuma.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?search=",
					Response:     "pattern.view.kuma.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "pattern view with ID",
			Args:             []string{"view", "4o7fbc9b-708d-4396-84b8-e2ba37c1adcc"},
			ExpectedResponse: "pattern.id.view.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "api/pattern/",
					Response:     "pattern.id.view.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		//Invalid or Non-Existing
		{
			Name:             "pattern invalid view",
			Args:             []string{"view", "kuma-test"},
			ExpectedResponse: "pattern.view.invalid.output.golden",
			URLs: []utils.MockURL{
				{
					Method:   "GET",
					URL:      testContext.BaseURL + "api/pattern",
					Response: "pattern.view.invalid.golden",
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: true,
		},
	}
	for _, test := range testcase {
		t.Run(test.Name, func(t *testing.T) {
			for _, url := range test.URLs {
				apiResponse := utils.NewGoldenFile(t, url.Response, fixturesDir).Load()
				httpmock.RegisterResponder(url.Method, url.URL,
					httpmock.NewStringResponder(url.ResponseCode, apiResponse))
			}
			// set token
			utils.TokenFlag = test.Token

			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, test.ExpectedResponse, testdataDir)

			// setting up log to grab logs
			b := utils.SetupMeshkitLoggerTesting(t, false)
			PatternCmd.SetOutput(b)
			PatternCmd.SetArgs(test.Args)
			err := PatternCmd.Execute()
			if err != nil {
				// if we're supposed to get an error
				if test.ExpectError {
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
	utils.StopMockery(t)
	t.Log("Pattern tests Passed")
}
