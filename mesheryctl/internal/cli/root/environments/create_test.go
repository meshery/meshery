package environments

import (
	
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"

	// "github.com/layer5io/meshkit/utils/kubernetes/describe"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCreateEnvironmentCmdValidation(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("orgId", "", "")
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("description", "", "")

	tests := []struct {
		name          string
		args          []string
		expectedError bool
	}{
		{
			name:          "All flags present",
			args:          []string{"--orgId", "org123", "--name", "TestEnv", "--description", "Test Description"},
			expectedError: false,
		},
		{
			name:          "Missing orgId",
			args:          []string{"--name", "TestEnv", "--description", "Test Description"},
			expectedError: true,
		},
		{
			name:          "Missing name",
			args:          []string{"--orgId", "org123", "--description", "Test Description"},
			expectedError: true,
		},
		{
			name:          "Missing description",
			args:          []string{"--orgId", "org123", "--name", "TestEnv"},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.Flags().Set("orgId", "")
			cmd.Flags().Set("name", "")
			cmd.Flags().Set("description", "")

			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			// cmd.SetArgs(tt.args)
			cmd.ParseFlags(tt.args)

			err := createEnvironmentCmd.Args(cmd, []string{})
			if (err != nil) != tt.expectedError {
				t.Errorf("createEnvironmentCmd.Args() error = %v, wantErr %v", err, tt.expectedError)
			}else{
				t.Log("Flags test passed")
			}

		})
		// 
	}
}

func TestCreateEnvironmentCmd(t *testing.T) {

	utils.SetupContextEnv(t)

	// Set up the mock server for handling requests
	utils.StartMockery(t)
	defer utils.StopMockery(t)

	// Create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")


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
			Name:             "Create Environment Successfully",
			Args:             []string{"create", "--orgId", "422595a1-bbe3-4355-ac80-5efa0b35c9da", "--name", "TestEnv", "--description", "This is a test environment"},
			URL:              testContext.BaseURL + "/api/environments",
			Fixture:          "create.environment.response.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectedResponse: "create.environment.output.golden",
			ExpectedError:      false,
		},
		// {
		// 	Name:             "Missing Name Flag",
		// 	Args:             []string{"exp", "environment", "create", "--orgID", "1234", "--description", "This is a test environment"},
		// 	URL:               testContext.BaseURL+"/api/environments",
		// 	ExpectedResponse: "create.environment.missing.arg.golden",
		// 	ExpectedError:      true,
		// },
		// {
		// 	Name:             "Missing OrgID Flag",
		// 	Args:             []string{"exp", "environment", "create", "--name", "TestEnv", "--description", "This is a test environment"},
		// 	URL:               testContext.BaseURL+"/api/environments",
		// 	ExpectedResponse: "create.environment.missing.arg.golden",
		// 	ExpectedError:      true,
		// },
	}

	// Loop through each test case
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Mock API response for a successful environment creation
			// mockResponse := &utils.MockResponse{
			// 	StatusCode: http.StatusOK,
			// 	Body:       `{"message": "environment created successfully"}`,
			apiResponse:=  utils.NewGoldenFile(t,tt.Fixture,fixturesDir).Load()
			// set token
			utils.TokenFlag = tt.Token

			// mock response
			httpmock.RegisterResponder("POST", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)


			cmd := EnvironmentCmd
			cmd.SetArgs(tt.Args)
			b := utils.SetupMeshkitLoggerTesting(t, false)
			cmd.SetOutput(b)
			err:=cmd.Execute()
			
			if err != nil {
				// if we're supposed to get an error
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
