package filter

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestDeleteCmd(t *testing.T) {
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

	testcase := []struct {
		Name             string
		Args             []string
		URL              string
		Fixture          string
		Token            string
		ExpectedResponse string
		ExpectedError    bool
		Method           string
		ResponseCode     int
	}{
		{
			Name:             "Delete Kuma-Test",
			Args:             []string{"delete", "c0c6035a-b1b9-412d-aab2-4ed1f1d51f84", "Kuma-Test"},
			URL:              testContext.BaseURL + "/api/filter/c0c6035a-b1b9-412d-aab2-4ed1f1d51f84",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			Fixture:          "delete.kuma.api.response.golden",
			ExpectedResponse: "delete.kuma.output.golden",
			ExpectedError:    false,
			Method:           "DELETE",
			ResponseCode:     200,
		}, {
			Name:             "Delete RolloutAndIstio",
			Args:             []string{"delete", "d0e09134-acb6-4c71-b051-3d5611653f70", "RolloutAndIstio"},
			URL:              testContext.BaseURL + "/api/filter/d0e09134-acb6-4c71-b051-3d5611653f70",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			Fixture:          "delete.rollout.api.response.golden",
			ExpectedResponse: "delete.rollout.output.golden",
			ExpectedError:    false,
			Method:           "DELETE",
			ResponseCode:     200,
		},
	}
	for _, tt := range testcase {
		t.Run(tt.Name, func(t *testing.T) {
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()
			// set token
			utils.TokenFlag = tt.Token
			// mock response
			httpmock.RegisterResponder(tt.Method, tt.URL,
				httpmock.NewStringResponder(tt.ResponseCode, apiResponse))
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			b := utils.SetupMeshkitLoggerTesting(t, false)
			FilterCmd.SetOutput(b)

			FilterCmd.SetArgs(tt.Args)
			err := FilterCmd.Execute()
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
				t.Fatal(err)
			}
			// response being printed in console
			output := b.String()
			actualResponse := output

			// write it in file
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			utils.Equals(t, expectedResponse, actualResponse)
		})
	}
	utils.StopMockery(t)
	t.Log("test passed")
}
