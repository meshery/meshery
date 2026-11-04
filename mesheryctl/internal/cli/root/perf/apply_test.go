// package perf

// import (
// 	"path/filepath"
// 	"runtime"
// 	"testing"

// 	"github.com/jarcoal/httpmock"
// 	"github.com/meshery/meshery/mesheryctl/pkg/utils"
// )

// var existingProfileID = "8f3daf25-e58e-4c59-8bf8-f474b76463ec"
// var newProfileID = "906f8876-33b5-4a97-906e-7a409d3b8ae9"

// // golden file responses
// var (
// 	// server returning existing profile "new"
// 	apply1001 = "1001.golden"
// 	// server response for no protocol added to url
// 	apply1003 = "1003.golden"
// 	// server response for creating new profile "test"
// 	apply1004 = "1004.golden"
// 	// server response for no profiles found
// 	apply1006 = "1006.golden"
// )

// var (
// 	// mesheryctl response for url not having protocol added
// 	apply1002output = "1002.golden"
// 	// mesheryctl response for invalid url
// 	apply1004output = "1004.golden"
// 	// mesheryctl response for no existing profile provided neither new profile-name
// 	apply1005output = "1005.golden"
// 	// mesheryctl response for no profiles found
// 	apply1006output = "1006.golden"
// )

// func TestApplyCmd(t *testing.T) {
// 	utils.SetupContextEnv(t)
// 	utils.StartMockery(t)
// 	testContext := utils.NewTestHelper(t)
// 	// get current directory
// 	_, filename, _, ok := runtime.Caller(0)
// 	if !ok {
// 		t.Fatal("Not able to get current working directory")
// 	}
// 	currDir := filepath.Dir(filename)
// 	fixturesDir := filepath.Join(currDir, "fixtures", "apply")
// 	testToken := filepath.Join(currDir, "fixtures", "auth.json")
// 	profileURL := testContext.BaseURL + "/api/performance/profiles"
// 	existingProfileRunTest := testContext.BaseURL + "/api/user/performance/profiles/" + existingProfileID + "/run"
// 	newProfileRunTest := testContext.BaseURL + "/api/user/performance/profiles/" + newProfileID + "/run"
// 	testdataDir := filepath.Join(currDir, "testdata", "apply")

// 	// test scenrios for fetching data
// 	tests := []tempTestStruct{
// 		{"Run Test with Existing profile", []string{"apply", "new"},
// 			[]utils.MockURL{
// 				{Method: "GET", URL: profileURL, Response: apply1001, ResponseCode: 200},
// 				{Method: "GET", URL: existingProfileRunTest, Response: apply1001, ResponseCode: 400},
// 			},
// 			apply1002output,
// 			testToken, false,
// 		},
// 		{"Run Test with Existing profile with --url", []string{"apply", "new", "--url", "https://www.google.com"},
// 			[]utils.MockURL{
// 				{Method: "GET", URL: profileURL, Response: apply1001, ResponseCode: 200},
// 				{Method: "GET", URL: existingProfileRunTest, Response: apply1001, ResponseCode: 400},
// 			},
// 			apply1002output,
// 			testToken, false,
// 		},
// 		{"Run Test with Existing profile with --url without protocol", []string{"apply", "new", "--url", "www.google.com"},
// 			[]utils.MockURL{
// 				{Method: "GET", URL: profileURL, Response: apply1001, ResponseCode: 200},
// 				{Method: "GET", URL: existingProfileRunTest, Response: apply1003, ResponseCode: 400},
// 			},
// 			apply1002output,
// 			testToken, true,
// 		},
// 		{"Run Test with new profile with --url without protocol", []string{"apply", "test", "--url", "www.google.com"},
// 			[]utils.MockURL{
// 				{Method: "POST", URL: profileURL, Response: apply1004, ResponseCode: 200},
// 				{Method: "GET", URL: newProfileRunTest, Response: apply1003, ResponseCode: 400},
// 			},
// 			apply1002output,
// 			testToken, true,
// 		},
// 		{"Run Test without profile-name and id", []string{"apply"},
// 			[]utils.MockURL{},
// 			apply1005output,
// 			testToken, true,
// 		},
// 		{"No profiles found with given name", []string{"apply", "new", "--yes"},
// 			[]utils.MockURL{
// 				{Method: "GET", URL: profileURL, Response: apply1006, ResponseCode: 200},
// 			},
// 			apply1006output,
// 			testToken, true,
// 		},
// 		{"Run Test with new profile with Invalid URL", []string{"apply", "test", "--url", "invalid-url"},
// 			[]utils.MockURL{},
// 			apply1004output,
// 			testToken, true,
// 		},
// 		{"Run Test with Existing profile with Invalid URL", []string{"apply", "new", "--url", "invalid-url"},
// 			[]utils.MockURL{{Method: "GET", URL: profileURL, Response: apply1001, ResponseCode: 200}},
// 			apply1004output,
// 			testToken, true,
// 		},
// 	}

// 	// Run tests
// 	for _, tt := range tests {
// 		t.Run(tt.Name, func(t *testing.T) {
// 			for _, url := range tt.URLs {
// 				apiResponse := utils.NewGoldenFile(t, url.Response, fixturesDir).Load()
// 				httpmock.RegisterResponder(url.Method, url.URL,
// 					httpmock.NewStringResponder(url.ResponseCode, apiResponse))
// 			}

// 			utils.TokenFlag = tt.Token
// 			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
// 			b := utils.SetupMeshkitLoggerTesting(t, false)

// 			PerfCmd.SetArgs(tt.Args)
// 			PerfCmd.SetOut(b)
// 			err := PerfCmd.Execute()
// 			if err != nil {
// 				if tt.ExpectError {
// 					if *update {
// 						golden.Write(err.Error())
// 					}
// 					expectedResponse := golden.Load()
// 					utils.Equals(t, expectedResponse, err.Error())
// 					resetVariables()
// 					return
// 				}
// 				t.Error(err)
// 			}

//				actualResponse := b.String()
//				if *update {
//					golden.Write(actualResponse)
//				}
//				expectedResponse := golden.Load()
//				utils.Equals(t, expectedResponse, actualResponse)
//				resetVariables()
//			})
//			t.Log("Apply Perf test Passed")
//		}
//		utils.StopMockery(t)
//	}
package perf

import (
	"os"
	"path/filepath"
	"testing"

	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
)

func clearPerfFlags() {
	profileName = ""
	testURL = ""
	testName = ""
	testMesh = ""
	qps = ""
	concurrentRequests = ""
	testDuration = ""
	loadGenerator = ""
	loadTestBody = ""
	filePath = ""
}

func TestApplyProfileParams(t *testing.T) {
	clearPerfFlags()
	defer resetVariables()

	// Pre-set flags: URL and QPS set via CLI flags
	testURL = "https://flag-override.example.com"
	qps = "250"

	profile := perfprofile.PerformanceProfile{
		Name:              "test-profile",
		Endpoints:         []string{"https://profile.example.com"},
		LoadGenerators:    []string{"wrk2"},
		ConcurrentRequest: 8,
		QPS:               50,
		Duration:          "45s",
		ServiceMesh:       "linkerd",
		RequestBody:       "{\"test\":true}",
	}

	applyProfileParams(profile)

	if profileName != "test-profile" {
		t.Errorf("expected profileName 'test-profile', got %q", profileName)
	}
	if testURL != "https://flag-override.example.com" {
		t.Errorf("expected testURL from flag, got %q", testURL)
	}
	if qps != "250" {
		t.Errorf("expected qps from flag '250', got %q", qps)
	}
	if loadGenerator != "wrk2" {
		t.Errorf("expected loadGenerator 'wrk2', got %q", loadGenerator)
	}
	if concurrentRequests != "8" {
		t.Errorf("expected concurrentRequests '8', got %q", concurrentRequests)
	}
	if testDuration != "45s" {
		t.Errorf("expected testDuration '45s', got %q", testDuration)
	}
	if testMesh != "linkerd" {
		t.Errorf("expected testMesh 'linkerd', got %q", testMesh)
	}
	if loadTestBody != "{\"test\":true}" {
		t.Errorf("expected loadTestBody from profile, got %q", loadTestBody)
	}
}

func TestApplyFileParams(t *testing.T) {
	clearPerfFlags()
	defer resetVariables()

	// Error case: nonexistent file
	err := applyFileParams("nonexistent-file.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}

	// Valid SMP config
	tempDir := t.TempDir()
	confPath := filepath.Join(tempDir, "smp.yaml")
	content := `
mesh: istio
test:
  name: smp-test
  duration: 1m
  clients:
    - endpointUrls:
        - https://smp.example.com
      rps: 500
      connections: 20
      loadGenerator: fortio
      body: "payload"
`
	if err := os.WriteFile(confPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Pre-set testName via flag
	testName = "flag-test-name"

	if err := applyFileParams(confPath); err != nil {
		t.Fatalf("unexpected error applying file params: %v", err)
	}

	if testName != "flag-test-name" {
		t.Errorf("expected testName from flag, got %q", testName)
	}
	if testURL != "https://smp.example.com" {
		t.Errorf("expected testURL from file, got %q", testURL)
	}
	if testMesh != "istio" {
		t.Errorf("expected testMesh 'istio', got %q", testMesh)
	}
	if qps != "500" {
		t.Errorf("expected qps '500', got %q", qps)
	}
	if concurrentRequests != "20" {
		t.Errorf("expected concurrentRequests '20', got %q", concurrentRequests)
	}
	if testDuration != "1m" {
		t.Errorf("expected testDuration '1m', got %q", testDuration)
	}
	if loadGenerator != "fortio" {
		t.Errorf("expected loadGenerator 'fortio', got %q", loadGenerator)
	}
	if loadTestBody != "payload" {
		t.Errorf("expected loadTestBody 'payload', got %q", loadTestBody)
	}
}

func TestSplitTestDurationDefaultsEmptyDuration(t *testing.T) {
	value, unit := splitTestDuration("")

	if value != "30" || unit != "s" {
		t.Fatalf("expected empty duration to default to 30s, got value=%q unit=%q", value, unit)
	}
}

func TestSplitTestDurationTrimsDuration(t *testing.T) {
	value, unit := splitTestDuration(" 5m ")

	if value != "5" || unit != "m" {
		t.Fatalf("expected trimmed duration to split into 5 and m, got value=%q unit=%q", value, unit)
	}
}

func TestSplitTestDurationDefaultsUnitlessDuration(t *testing.T) {
	value, unit := splitTestDuration("30")

	if value != "30" || unit != "s" {
		t.Fatalf("expected unitless duration to default to seconds, got value=%q unit=%q", value, unit)
	}
}

func TestSplitTestDurationSupportsMultiCharacterUnit(t *testing.T) {
	value, unit := splitTestDuration("500ms")

	if value != "500" || unit != "ms" {
		t.Fatalf("expected multi-character unit to be preserved, got value=%q unit=%q", value, unit)
	}
}

func resetVariables() {
	// reset the variables after each test
	profileName = ""
	testURL = ""
	testName = ""
	testMesh = "none"
	qps = "0"
	concurrentRequests = "1"
	testDuration = "30s"
	loadGenerator = "fortio"
	filePath = ""
	outputFormatFlag = ""
	viewSingleProfile = false
	viewSingleResult = false
}
