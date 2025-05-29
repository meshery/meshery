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

	tests := []struct {
		name               string
		args               []string
		profilesResponse   string
		profilesStatusCode int
		deleteResponse     string
		deleteStatusCode   int
		silent             bool
		expectError        bool
		tokenPath          string
		profileURL         string
		deleteURL          string
		profileID          string
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
			tokenPath:          filepath.Join(currDir, "fixtures", "auth.json"),
			profileURL:         testContext.BaseURL + "/api/user/performance/profiles",
			deleteURL:          testContext.BaseURL + "/api/user/performance/profiles/11111111-1111-1111-1111-111111111111",
			profileID:          "11111111-1111-1111-1111-111111111111",
		},
		{
			name:               "Profile not found",
			args:               []string{"delete", "nonexistent-profile"},
			profilesResponse:   `{"page": 0, "page_size": 10, "total_count": 0, "profiles": []}`,
			profilesStatusCode: 200,
			expectError:        false,
			tokenPath:          filepath.Join(currDir, "fixtures", "auth.json"),
			profileURL:         testContext.BaseURL + "/api/user/performance/profiles",
			deleteURL:          "",
			profileID:          "",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.TokenFlag = tt.tokenPath

			httpmock.RegisterResponder("GET", tt.profileURL,
				httpmock.NewStringResponder(tt.profilesStatusCode, tt.profilesResponse))

			if tt.deleteResponse != "" && tt.deleteURL != "" {
				httpmock.RegisterResponder("DELETE", tt.deleteURL,
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
