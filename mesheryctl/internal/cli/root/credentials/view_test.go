package credentials

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestViewCredentialCmd(t *testing.T) {
	// Setup current context
	utils.SetupContextEnv(t)

	// Initialize mock server for handling requests
	utils.StartMockery(t)
	defer utils.StopMockery(t)

	// Create a test helper
	testContext := utils.NewTestHelper(t)

	// Get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	// Create a mock credential file
	mockCredentialFile := filepath.Join(fixturesDir, "credential.json")
	err := os.WriteFile(mockCredentialFile, []byte(`{"token": "mockToken"}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create mock credential file: %v", err)
	}
	// Ensure the mock credential file is removed after the test
	defer func() {
		err := os.Remove(mockCredentialFile)
		if err != nil {
			t.Errorf("Failed to remove mock credential file: %v", err)
		}
	}()
	// Point to the mock credential file
	os.Setenv("MESHERY_CREDENTIAL_FILE", mockCredentialFile)

	// Generate a UUID for credential ID
	credentialID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("Failed to generate UUID: %v", err)
	}

	// Test scenarios for viewing a credential
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
			Name:             "View Credential",
			Args:             []string{"view", credentialID.String()},
			URL:              testContext.BaseURL + "/api/integrations/credentials/" + credentialID.String(),
			ExpectedResponse: "view.credential.output.golden",
			Response:         "view.credential.api.response.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		// Add more test cases for different scenarios (invalid input, server errors, etc.)
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			apiResponse := utils.NewGoldenFile(t, tt.Response, fixturesDir).Load()
			utils.TokenFlag = tt.Token
			// Setup http mock
			httpmock.RegisterResponder("GET", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Redirect command output to buffer for testing
			buf := utils.SetupMeshkitLoggerTesting(t, false)

			// Set command arguments
			CredentialCmd.SetArgs(tt.Args)
			CredentialCmd.SetOutput(buf)

			// Execute the command
			err := CredentialCmd.Execute()

			// Verify the command output and error
			if err != nil {
				if tt.ExpectError {
					// Write it in file if *update
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

			// Write it in file if *update
			if *update {
				golden.Write(actualResponse)
			}

			// Compare with expected response
			expectedResponse := golden.Load()
			utils.Equals(t, expectedResponse, actualResponse)
		})
		t.Log("View Credential test passed")
	}
}
