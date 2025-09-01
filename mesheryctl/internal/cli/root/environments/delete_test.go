package environments

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestDeleteEnvironment(t *testing.T) {
	// Setup current context
	utils.SetupContextEnv(t)

	// Initialize mock server for handling requests
	utils.StartMockery(t)

	// Create a test helper
	testContext := utils.NewTestHelper(t)

	// Get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	environmentId := "d56fb25b-f92c-4cd6-821b-2cfd6bb87259"

	// Test scenarios for environment creation
	tests := []struct {
		Name             string
		Args             []string
		URL              string
		Method           string
		Fixture          string
		RequestBody      string
		ExpectedResponse string
		Token            string
		ExpectError      bool
	}{
		{
			Name:             "Delete environment without arguments",
			Args:             []string{"delete"},
			URL:              testContext.BaseURL + "/api/environments",
			Method:           "DELETE",
			Fixture:          "",
			ExpectedResponse: "delete.environment.without.name.golden",
			ExpectError:      true,
		},
		{
			Name:             "Delete environment successfully",
			Args:             []string{"delete", environmentId},
			URL:              fmt.Sprintf("%s/api/environments/%s", testContext.BaseURL, environmentId),
			Method:           "DELETE",
			Fixture:          "delete.environment.response.golden",
			ExpectedResponse: "delete.environment.success.golden",
			ExpectError:      false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Fixture != "" {
				apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

				utils.TokenFlag = utils.GetToken(t)

				httpmock.RegisterResponder(tt.Method, tt.URL,
					httpmock.NewStringResponder(200, apiResponse))
			}

			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			b := utils.SetupMeshkitLoggerTesting(t, false)
			EnvironmentCmd.SetArgs(tt.Args)
			EnvironmentCmd.SetOut(b)
			err := EnvironmentCmd.Execute()

			if err != nil {

				if tt.ExpectError {

					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()

					utils.Equals(t, expectedResponse, err.Error())
					return
				}
				t.Error(err)
			}

			actualResponse := b.String()

			if *update {
				golden.Write(actualResponse)
			}

			expectedResponse := golden.Load()

			cleanedActualResponse := utils.CleanStringFromHandlePagination(actualResponse)
			cleanedExpectedResponse := utils.CleanStringFromHandlePagination(expectedResponse)

			utils.Equals(t, cleanedExpectedResponse, cleanedActualResponse)
		})
		t.Logf("Create environment test '%s' passed", tt.Name)
	}

	utils.StopMockery(t)
}
