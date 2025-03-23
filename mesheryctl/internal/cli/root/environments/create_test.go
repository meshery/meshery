package environments

import (
	"flag"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update golden files")

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
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      true,
		},
		{
			Name:             "Create environment successfully",
			Args:             []string{"create", "--name", "test-environment", "--description", "integration test", "--orgID", "3f8319e0-33a9-4736-b248-12nm3kiuh3yu"},
			URL:              testContext.BaseURL + "/api/environments",
			Method:           "POST",
			Fixture:          "create.environment.response.golden",
			RequestBody:      `{"name":"test-environment"}`,
			ExpectedResponse: "create.environment.success.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},

		// {
		// 	Name:             "Create environment with server error",
		// 	Args:             []string{"create", "error-environment"},
		// 	URL:              testContext.BaseURL + "/api/environments",
		// 	Method:           "POST",
		// 	Fixture:          "create.environment.error.golden",
		// 	RequestBody:      `{"name":"error-environment"}`,
		// 	ExpectedResponse: "create.environment.server.error.golden",
		// 	Token:            filepath.Join(fixturesDir, "token.golden"),
		// 	ExpectError:      true,
		// },
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Setup token
			utils.TokenFlag = tt.Token

			// Register mock responder if fixture is provided
			if tt.Fixture != "" {
				apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

				// Register the responder with request body verification
				httpmock.RegisterResponder(tt.Method, tt.URL,
					func(req *http.Request) (*http.Response, error) {
						// Read request body if needed for verification
						if tt.RequestBody != "" {
							body, _ := io.ReadAll(req.Body)
							assert.JSONEq(t, tt.RequestBody, string(body), "Request body doesn't match expected")
						}
						return httpmock.NewStringResponse(200, apiResponse), nil
					})
			}

			// Setup golden file
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Capture stdout
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Setup logger and execute command
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			EnvironmentCmd.SetArgs(tt.Args)
			EnvironmentCmd.SetOut(rescueStdout)
			err := EnvironmentCmd.Execute()

			// Handle error cases
			if err != nil {
				if tt.ExpectError {
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()
					actualResponse := err.Error()
					utils.Equals(t, expectedResponse, actualResponse)
					return
				}
				t.Error(err)
			}

			// Capture output
			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout
			actualResponse := string(out)

			// Update golden file if requested
			if *update {
				golden.Write(actualResponse)
			}

			// Compare with expected response
			expectedResponse := golden.Load()
			cleanedActualResponse := utils.CleanStringFromHandlePagination(actualResponse)
			cleanedExpectedResponse := utils.CleanStringFromHandlePagination(expectedResponse)
			utils.Equals(t, cleanedExpectedResponse, cleanedActualResponse)
		})
		t.Logf("Create environment test '%s' passed", tt.Name)
	}

	utils.StopMockery(t)
}
