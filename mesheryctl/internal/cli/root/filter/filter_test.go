package filter

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestFilterCmd(t *testing.T) {
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
	currentDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currentDir, "fixtures")

	FilterTestCases := []struct {
		Name             string
		Args             []string
		URL              string
		Fixture          string
		Token            string
		ExpectedResponse string
		ExpectHelp       bool
		ExpectErr        bool
	}{
		{
			Name:             "filter viewcmd with name",
			Args:             []string{"view", "view-filter-name"},
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectedResponse: "filter.id.view.output.golden",
			Fixture:          "filter.name.view.api.response.golden",
			URL:              testContext.BaseURL + "/api/filter",
			ExpectHelp:       false,
			ExpectErr:        false,
		},
		{
			Name:             "filter viewcmd with ID",
			Args:             []string{"view", "c0c6035a-b1b9-412d-aab2-4ed1f1d51f84"},
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectedResponse: "filter.id.view.output.golden",
			Fixture:          "filter.id.view.api.response.golden",
			URL:              testContext.BaseURL + "/api/filter/c0c6035a-b1b9-412d-aab2-4ed1f1d51f84",
			ExpectHelp:       false,
			ExpectErr:        false,
		},
		{
			Name:             "filter deletecmd",
			Args:             []string{"delete", "c0c6035a-b1b9-412d-aab2-4ed1f1d51f84", "Kuma-test"},
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectedResponse: "filter.id.delete.output.golden",
			Fixture:          "filter.id.delete.api.response.golden",
			URL:              testContext.BaseURL + "/api/filter/deploy/c0c6035a-b1b9-412d-aab2-4ed1f1d51f84",
			ExpectHelp:       false,
			ExpectErr:        false,
		},
		{
			Name:             "filter listcmd",
			Args:             []string{"list"},
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectedResponse: "filter.list.output.golden",
			Fixture:          "filter.list.api.response.golden",
			URL:              testContext.BaseURL + "/api/filter",
			ExpectHelp:       false,
			ExpectErr:        false,
		},
		{
			Name:             "filter viewcmd For Non-Existing Name",
			Args:             []string{"view", "test-existence"},
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectedResponse: "filter.nonexisting-name.view.output.golden",
			Fixture:          "filter.nonexisting-name.view.response.golden",
			URL:              testContext.BaseURL + "/api/filter",
			ExpectHelp:       true,
			ExpectErr:        true,
		},
		{
			Name:             "filter viewcmd with Invalid ID",
			Args:             []string{"view", "957fbc9b-a655-4892-823d-375102a9587c"},
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectedResponse: "filter.invalidID.view.output.golden",
			Fixture:          "filter.invalidID.view.response.golden",
			URL:              testContext.BaseURL + "/api/filter/957fbc9b-a655-4892-823d-375102a9587c",
			ExpectHelp:       true,
			ExpectErr:        true,
		},
	}
	for _, tc := range FilterTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			// View api response from golden files
			apiResponse := utils.NewGoldenFile(t, tc.Fixture, fixturesDir).Load()

			//set token
			utils.TokenFlag = tc.Token
			// mock response
			httpmock.RegisterResponder("GET", tc.URL,
				httpmock.NewStringResponder(200, apiResponse))

			httpmock.RegisterResponder("DELETE", tc.URL,
				httpmock.NewStringResponder(200, apiResponse))

			//Expected Response
			testdataDir := filepath.Join(currentDir, "testdata")
			golden := utils.NewGoldenFile(t, tc.ExpectedResponse, testdataDir)

			//Console Prints
			StdOut := os.Stdout
			//return expected pair of files/output
			read, write, _ := os.Pipe()
			os.Stdout = write
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			cmd := FilterCmd
			cmd.SetArgs(tc.Args)
			cmd.SetOutput(StdOut)

			err := cmd.Execute()
			if err != nil {
				if tc.ExpectHelp && tc.ExpectErr {
					//write in file
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()
					utils.Equals(t, expectedResponse, err.Error())
					return
				}
				t.Fatal(err)
			}
			write.Close()
			output, _ := io.ReadAll(read)
			os.Stdout = StdOut
			//print response string to console
			Response := string(output)
			// write it in file
			if *update {
				golden.Write(Response)
			}
			expectedResponse := golden.Load()

			utils.Equals(t, expectedResponse, Response)

		})
	}
	// stop mock server
	utils.StopMockery(t)
	t.Log("Filter tests Passed")
}
