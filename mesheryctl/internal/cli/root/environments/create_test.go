package environments

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestCreateEnvironment(t *testing.T) {
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

	// Test scenarios for environment creation
	tests := []struct {
		Name             string
		Args             []string
		URL              string
		Method           string
		Fixture          string
		RequestBody      string
		ExpectedResponse string
		ResponseCode     int
		Token            string
		ExpectError      bool
	}{
		{
			Name:             "Create environment without arguments",
			Args:             []string{"create"},
			URL:              testContext.BaseURL + "/api/environments",
			Method:           "POST",
			Fixture:          "",
			ExpectedResponse: "create.environment.without.name.golden",
			ResponseCode:     200,
			ExpectError:      true,
		},
		{
			Name:             "Create environment successfully",
			Args:             []string{"create", "--name", testConstants["environmentName"], "--description", "integration test", "--orgID", testConstants["orgID"]},
			URL:              testContext.BaseURL + "/api/environments",
			Method:           "POST",
			Fixture:          "create.environment.response.golden",
			ExpectedResponse: "create.environment.success.golden",
			ResponseCode:     200,
			ExpectError:      false,
		},
		{
			Name:             "Create environment failure with bad organization ID",
			Args:             []string{"create", "--name", testConstants["environmentName"], "--description", "integration test", "--orgID", testConstants["invalidOrgID"]},
			URL:              testContext.BaseURL + "/api/environments",
			Method:           "POST",
			Fixture:          "create.environment.failure.response.golden",
			ExpectedResponse: "create.environment.failure.golden",
			ResponseCode:     400,
			ExpectError:      true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Fixture != "" {
				apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

				utils.TokenFlag = utils.GetToken(t)

				httpmock.RegisterResponder(tt.Method, tt.URL,
					httpmock.NewStringResponder(tt.ResponseCode, apiResponse))
			}

			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			buf := utils.SetupMeshkitLoggerTesting(t, false)

			EnvironmentCmd.SetArgs(tt.Args)
			EnvironmentCmd.SetOut(buf)
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

			actualResponse := buf.String()

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
