package perf

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var existingProfileID = "8f3daf25-e58e-4c59-8bf8-f474b76463ec"
var newProfileID = "906f8876-33b5-4a97-906e-7a409d3b8ae9"

// golden file responses
var (
	// server returning existing profile "new"
	apply1001 = "1001.golden"
	// server response for no protocol added to url
	apply1003 = "1003.golden"
	// server response for creating new profile "test"
	apply1004 = "1004.golden"
	// server response for no profiles found
	apply1006 = "1006.golden"
)

var (
	// mesheryctl response for url not having protocol added
	apply1002output = "1002.golden"
	// mesheryctl response for invalid url
	apply1004output = "1004.golden"
	// mesheryctl response for no existing profile provided neither new profile-name
	apply1005output = "1005.golden"
	// mesheryctl response for no profiles found
	apply1006output = "1006.golden"
)

func TestApplyCmd(t *testing.T) {
	utils.SetupContextEnv(t)
	utils.StartMockery(t)
	testContext := utils.NewTestHelper(t)
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures", "apply")
	testToken := filepath.Join(currDir, "fixtures", "auth.json")
	profileURL := testContext.BaseURL + "/api/user/performance/profiles"
	existingProfileRunTest := testContext.BaseURL + "/api/user/performance/profiles/" + existingProfileID + "/run"
	newProfileRunTest := testContext.BaseURL + "/api/user/performance/profiles/" + newProfileID + "/run"
	testdataDir := filepath.Join(currDir, "testdata", "apply")

	// test scenrios for fetching data
	tests := []tempTestStruct{
		{"Run Test with Existing profile", []string{"apply", "new"},
			[]utils.MockURL{
				{Method: "GET", URL: profileURL, Response: apply1001, ResponseCode: 200},
				{Method: "GET", URL: existingProfileRunTest, Response: apply1001, ResponseCode: 400},
			},
			apply1002output,
			testToken, false,
		},
		{"Run Test with Existing profile with --url", []string{"apply", "new", "--url", "https://www.google.com"},
			[]utils.MockURL{
				{Method: "GET", URL: profileURL, Response: apply1001, ResponseCode: 200},
				{Method: "GET", URL: existingProfileRunTest, Response: apply1001, ResponseCode: 400},
			},
			apply1002output,
			testToken, false,
		},
		{"Run Test with Existing profile with --url without protocol", []string{"apply", "new", "--url", "www.google.com"},
			[]utils.MockURL{
				{Method: "GET", URL: profileURL, Response: apply1001, ResponseCode: 200},
				{Method: "GET", URL: existingProfileRunTest, Response: apply1003, ResponseCode: 400},
			},
			apply1002output,
			testToken, true,
		},
		{"Run Test with new profile with --url without protocol", []string{"apply", "test", "--url", "www.google.com"},
			[]utils.MockURL{
				{Method: "POST", URL: profileURL, Response: apply1004, ResponseCode: 200},
				{Method: "GET", URL: newProfileRunTest, Response: apply1003, ResponseCode: 400},
			},
			apply1002output,
			testToken, true,
		},
		{"Run Test without profile-name and id", []string{"apply"},
			[]utils.MockURL{},
			apply1005output,
			testToken, true,
		},
		{"No profiles found with given name", []string{"apply", "new", "--yes"},
			[]utils.MockURL{
				{Method: "GET", URL: profileURL, Response: apply1006, ResponseCode: 200},
			},
			apply1006output,
			testToken, true,
		},
		{"Run Test with new profile with Invalid URL", []string{"apply", "test", "--url", "invalid-url"},
			[]utils.MockURL{},
			apply1004output,
			testToken, true,
		},
		{"Run Test with Existing profile with Invalid URL", []string{"apply", "new", "--url", "invalid-url"},
			[]utils.MockURL{{Method: "GET", URL: profileURL, Response: apply1001, ResponseCode: 200}},
			apply1004output,
			testToken, true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			for _, url := range tt.URLs {
				apiResponse := utils.NewGoldenFile(t, url.Response, fixturesDir).Load()
				httpmock.RegisterResponder(url.Method, url.URL,
					httpmock.NewStringResponder(url.ResponseCode, apiResponse))
			}

			utils.TokenFlag = tt.Token
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			b := utils.SetupMeshkitLoggerTesting(t, false)

			PerfCmd.SetArgs(tt.Args)
			PerfCmd.SetOutput(b)
			err := PerfCmd.Execute()
			if err != nil {
				if tt.ExpectError {
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()
					utils.Equals(t, expectedResponse, err.Error())
					resetVariables()
					return
				}
				t.Error(err)
			}

			actualResponse := b.String()
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			utils.Equals(t, expectedResponse, actualResponse)
			resetVariables()
		})
		t.Log("Apply Perf test Passed")
	}
	utils.StopMockery(t)
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
