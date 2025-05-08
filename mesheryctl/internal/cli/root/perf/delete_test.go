package perf

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestDeleteCmd(t *testing.T) {
	// Setup test environment
	utils.SetupContextEnv(t)
	utils.StartMockery(t)
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	testToken := filepath.Join(currDir, "fixtures", "auth.json")

	// Mock responses directly
	profileURL := testContext.BaseURL + "/api/user/performance/profiles"

	tests := []struct {
		name               string
		args               []string
		profilesResponse   string
		profilesStatusCode int
		deleteResponse     string
		deleteStatusCode   int
		silent             bool
		expectError        bool
	}{
		{
			name: "Profile found and deleted",
			args: []string{"delete", "test-profile"},
			profilesResponse: `{
				"page": 0,
				"page_size": 10,
				"total_count": 1,
				"profiles": [
					{
						"id": "11111111-1111-1111-1111-111111111111",
						"name": "test-profile",
						"load_generators": ["fortio"],
						"endpoints": ["http://example.com"],
						"service_mesh": "None",
						"total_results": 1
					}
				]
			}`,
			profilesStatusCode: 200,
			deleteResponse:     `{"status": "ok"}`,
			deleteStatusCode:   200,
			silent:             true,
			expectError:        false,
		},
		{
			name:               "Profile not found",
			args:               []string{"delete", "nonexistent-profile"},
			profilesResponse:   `{"page": 0, "page_size": 10, "total_count": 0, "profiles": []}`,
			profilesStatusCode: 200,
			expectError:        false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.TokenFlag = testToken

			httpmock.RegisterResponder("GET", profileURL,
				httpmock.NewStringResponder(tt.profilesStatusCode, tt.profilesResponse))

			if tt.deleteResponse != "" {
				deleteURL := testContext.BaseURL + "/api/user/performance/profiles/11111111-1111-1111-1111-111111111111"
				httpmock.RegisterResponder("DELETE", deleteURL,
					httpmock.NewStringResponder(tt.deleteStatusCode, tt.deleteResponse))
			}

			utils.SilentFlag = tt.silent

			b := utils.SetupMeshkitLoggerTesting(t, false)

			PerfCmd.SetArgs(tt.args)
			PerfCmd.SetOut(os.Stdout)
			err := PerfCmd.Execute()

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error %v, got %v", tt.expectError, err)
			}

			_ = b.String() 

			// Reset for next test
			httpmock.Reset()
			resetVariables()
		})
	}

	utils.StopMockery(t)
}
