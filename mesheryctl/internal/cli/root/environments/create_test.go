package environments

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateEnvironmentCmd(t *testing.T) {

	utils.SetupContextEnv(t)

	// Set up the mock server for handling requests
	utils.StartMockery(t)
	defer utils.StopMockery(t)

	testContext := utils.NewTestHelper(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")
	testdataDir := filepath.Join(currDir, "testdata")

	tests := []struct {
		Name             string
		Args             []string
		Token            string
		URL              string
		Fixture          string
		ExpectedResponse string
		ExpectedError    bool
	}{
		{
			Name:             "Missing Name Flag",
			Args:             []string{"create", "--orgID", "1234", "--description", "This is a test environment"},
			URL:              testContext.BaseURL + "/api/environments",
			ExpectedResponse: "create.environment.without.name.golden",
			ExpectedError:    true,
		},
		{
			Name:             "Invalid OrgId as Non-UUID",
			Args:             []string{"create", "--orgID", "invalid-org-id", "--name", "TestEnv", "--description", "This is a test environment"},
			URL:              testContext.BaseURL + "/api/environments",
			Fixture:          "create.environment.local.provider.invalid-orgid.response.golden",
			Token:            filepath.Join(utils.FixturesDir, "token.golden"),
			ExpectedResponse: "create.environment.invalid.orgId.golden",
			ExpectedError:    false,
		},
		{
			Name:             "Create Environment Successfully",
			Args:             []string{"create", "--orgID", "3f8319e0-33a9-4736-b248-12nm3kiuh3yu", "--name", "test-environment", "--description", "This is a test environment"},
			URL:              testContext.BaseURL + "/api/environments",
			Fixture:          "create.environment.response.golden",
			Token:            filepath.Join(utils.FixturesDir, "token.golden"),
			ExpectedResponse: "create.environment.success.golden",
			ExpectedError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Fixture != "" {
				apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()
				utils.TokenFlag = tt.Token

				// mock response
				httpmock.RegisterResponder("POST", tt.URL,
					httpmock.NewStringResponder(200, apiResponse))

			}
			// Expected response
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			b := utils.SetupMeshkitLoggerTesting(t, false)
			EnvironmentCmd.SetArgs(tt.Args)
			EnvironmentCmd.SetOut(b)
			err := EnvironmentCmd.Execute()

			if err != nil {
				if tt.ExpectedError {
					expectedResponse := golden.Load()
					utils.Equals(t, expectedResponse, err.Error())
					return
				}
				t.Fatal(err)
			}

			output := b.String()
			actualResponse := output
			expectedResponse := golden.Load()
			assert.Equal(t, expectedResponse, actualResponse)

		})
	}
	t.Log("Create test passed")
}