package credentials

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestCreateCredentialCmd(t *testing.T) {
	// Setup current context
	utils.SetupContextEnv(t)

	// Initialize mock server for handling requests
	utils.StartMockery(t)
	defer utils.StopMockery(t)

	// create a test helper
	testContext := utils.NewTestHelper(t)

	// Get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	// Generate a UUID for user-id
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("Failed to generate UUID: %v", err)
	}

	// Test scenarios for creating a credential
	tests := []struct {
		Name             string
		Args             []string
		URL              string
		ExpectedResponse string
		Response         string
		Token            string
		ExpectError      bool
	}{
		{
			Name:             "Create Credential",
			Args:             []string{"create", "--user-id", userID.String(), "--name", "test-name", "--type", "test-type", "--secret", "test-secret"},
			URL:              testContext.BaseURL + "/api/integrations/credentials",
			ExpectedResponse: "create.credential.output.golden",
			Response:         "create.credential.api.response.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		// Add more test cases for different scenarios (invalid input, server errors, etc.)
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			apiResponse := utils.NewGoldenFile(t, tt.Response, fixturesDir).Load()
			// set the token
			utils.TokenFlag = tt.Token

			// Mock the HTTP response
			httpmock.RegisterResponder("POST", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Redirect command output to buffer for testing
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			// No need to set args here, as we'll use prompts
			CredentialCmd.SetArgs(tt.Args)

			CredentialCmd.SetOutput(buf)
			// Execute the command
			err := CredentialCmd.Execute()

			// Verify the command output and error
			if err != nil {
				if tt.ExpectError {
					// write it in file
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()

					utils.Equals(t, expectedResponse, err.Error())
					return
				}
				t.Fatalf("Unexpected error: %v", err)
			}
			// Response being printed in console
			actualResponse := buf.String()
			// write it in file
			if *update {
				golden.Write(actualResponse)
			}
			// Compare with expected response
			expectedResponse := golden.Load()
			utils.Equals(t, expectedResponse, actualResponse)
		})
		t.Log("Create Credential test passed")
	}
}
