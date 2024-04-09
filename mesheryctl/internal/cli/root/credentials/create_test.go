package credentials

import (
	"net/http"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestCreateCredentialCmd(t *testing.T) {
	// Setup current context
	utils.SetupContextEnv(t)

	// Initialize mock server for handling requests
	utils.StartMockery(t)

	// create a test helper
	testContext := utils.NewTestHelper(t)

	// Get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	// Test scenarios for creating a credential
	tests := []struct {
		Name             string
		Args             []string
		URL              string
		ExpectedResponse string
		ResponseCode     int
		Fixture          string
		Method           string
		Token            string
		ExpectError      bool
	}{
		{
			Name:             "Create Credential",
			Args:             []string{"create", "--name", "testCredential", "--user-id", "12345", "--type", "kubernetes", "--secret", "abc123"},
			URL:              testContext.BaseURL + "/api/integrations/credentials",
			ExpectedResponse: "create.credential.output.golden",
			ResponseCode:     http.StatusOK,
			Fixture:          "create.credential.api.response.golden",
			Method:           "POST",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		// Add more test cases for different scenarios (invalid input, server errors, etc.)
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			// Mock the HTTP response
			httpmock.RegisterResponder(tt.Method, tt.URL,
				httpmock.NewStringResponder(tt.ResponseCode, tt.Fixture))

			// set the token
			utils.TokenFlag = tt.Token

			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Redirect command output to buffer for testing
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			createCredentialCmd.SetArgs(tt.Args)
			createCredentialCmd.SetOutput(buf)

			// Execute the command
			err := createCredentialCmd.Execute()

			// Verify the command output and error
			if err != nil {
				if tt.ExpectError {
					// If we're supposed to get an error
					if *update {
						golden.Write(err.Error())
					}

					expectedResponse := golden.Load()

					utils.Equals(t, expectedResponse, err.Error())
					return
				}
				t.Fatal(err)
			}
			// Response being printed in console
			actualResponse := buf.String()

			// Write it in file
			if *update {
				golden.Write(actualResponse)
			}
			// Compare with expected response
			expectedResponse := golden.Load()
			utils.Equals(t, expectedResponse, actualResponse)
		})
		t.Log("Create Credential test passed")
	}
	// stop mock server
	utils.StopMockery(t)
}
