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
		ExpectedError      bool
	}{
		{
			Name:             "Missing Name Flag",
			Args:             []string{"create", "--orgId", "1234", "--description", "This is a test environment"},
			URL:              testContext.BaseURL+"/api/environments",
			ExpectedResponse: "create.environment.missing.flag.golden",
			ExpectedError:    true,
		},
		{
			Name:             "Invalid OrgId as Non-UUID",
			Args:             []string{"create", "--orgId", "invalid-org-id", "--name", "TestEnv", "--description", "This is a test environment"},
			URL:              testContext.BaseURL + "/api/environments",
			Fixture:          "create.environment.local.provider.invalid-orgid.response.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectedResponse: "create.environment.output.golden",
			ExpectedError:    false,
		},
		{
			Name:             "Create Environment Successfully",
			Args:             []string{"create", "--orgId", "422595a1-bbe3-4355-ac80-5efa0b35c9da", "--name", "TestEnv", "--description", "This is a test environment"},
			URL:              testContext.BaseURL + "/api/environments",
			Fixture:          "create.environment.response.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectedResponse: "create.environment.output.golden",
			ExpectedError:    false,
		},
		}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if(tt.Fixture!=""){
				apiResponse:=  utils.NewGoldenFile(t,tt.Fixture,fixturesDir).Load()
				utils.TokenFlag = tt.Token

				// mock response
				httpmock.RegisterResponder("POST", tt.URL,
					httpmock.NewStringResponder(200, apiResponse))

			}
			// Expected response
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)


			cmd := EnvironmentCmd
			cmd.SetArgs(tt.Args)

			b := utils.SetupMeshkitLoggerTesting(t, false)
			cmd.SetOut(b)
			err:=cmd.Execute()

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