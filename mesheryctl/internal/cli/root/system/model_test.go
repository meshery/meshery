package system

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

// This is an Integration test
func TestListModelCmd(t *testing.T) {
	SetupContextEnv(t)
	// initialize mock server for handling requests
	utils.StartMockery(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	// create a test helper
	testContext := utils.NewTestHelper(t)
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

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
			Name:             "view list of models",
			Args:             []string{"model", "list"},
			ExpectedResponse: "list.model.output.golden",
			Fixture:          "list.model.api.response.golden",
			URL:              testContext.BaseURL + "/api/meshmodels/models?pagesize=all",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// View api response from golden files
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

			// set token
			utils.TokenFlag = tt.Token

			// mock response
			httpmock.RegisterResponder("GET", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			testdataDir := filepath.Join(currDir, "testdata/model/")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Grab console prints
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			SystemCmd.SetArgs(tt.Args)
			SystemCmd.SetOutput(rescueStdout)
			err = SystemCmd.Execute()
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

			utils.Equals(t, expectedResponse, actualResponse)
		})

		// we're done with testing stop mock server
		utils.StopMockery(t)
	}
}

func TestModelViewCmd(t *testing.T) {
	SetupContextEnv(t)
	// initialize mock server for handling requests
	utils.StartMockery(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	// create a test helper
	testContext := utils.NewTestHelper(t)
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

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
			Name:             "view a requested model in yaml format",
			Args:             []string{"model", "view", "spire"},
			ExpectedResponse: "view.model.yaml.output.golden",
			Fixture:          "view.model.api.response.golden",
			URL:              testContext.BaseURL + "/api/meshmodels/models/spire?pagesize=all",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "view a requested model in json format",
			Args:             []string{"model", "view", "spire", "-o", "json"},
			ExpectedResponse: "view.model.json.output.golden",
			Fixture:          "view.model.api.response.golden",
			URL:              testContext.BaseURL + "/api/meshmodels/models/spire?pagesize=all",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// View api response from golden files
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

			// set token
			utils.TokenFlag = tt.Token

			// mock response
			httpmock.RegisterResponder("GET", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			testdataDir := filepath.Join(currDir, "testdata/model/")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Grab console prints
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			SystemCmd.SetArgs(tt.Args)
			SystemCmd.SetOutput(rescueStdout)
			err = SystemCmd.Execute()
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

			utils.Equals(t, expectedResponse, actualResponse)
		})

		// we're done with testing stop mock server
		utils.StopMockery(t)
	}
}
