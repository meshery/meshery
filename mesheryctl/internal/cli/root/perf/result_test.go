package perf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var tempProfileID = "18358dfc-a009-4c76-9ab9-dfef33224b3b"

// golden file responses
var (
	// standard api response of 25 performance result
	result1001 = "1001.golden"
	// api response of performance result with searched term "kuma"
	result1002 = "1002.golden"
	// api response of performance result with searched term "no mesh"
	result1003 = "1003.golden"
	// api response of zero performance result
	result1004 = "1004.golden"
	// unable marshal response
	result1005 = "1005.golden"
	// empty response
	result1006 = "1006.golden"
)

// golden file mesheryctl outputs
var (
	// mesheryctl response of 25 performance result
	result1001output = "1001.golden"
	// mesheryctl response of performance result with searched term "kuma"
	result1002output = "1002.golden"
	// mesheryctl response of performance result with searched term "test 3"
	result1003output = "1003.golden"
	// mesheryctl response of 25 performance result in expanded output
	result1004output = "1004.golden"
	// mesheryctl response when no results found
	result1005output = "1005.golden"
	// mesheryctl response of 25 performance result in json output
	result1006output = "1006.golden"
	// mesheryctl response of 25 performance result in yaml output
	result1007output = "1007.golden"
	// mesheryctl response for invalid output format
	result1008output = "1008.golden"
	// mesheryctl response for server response code error
	result1009output = "1009.golden"
	// mesheryctl response for unmarshal error
	result1010output = "1010.golden"
	// mesheryctl response for failing attach authentication
	result1011output = "1011.golden"
	// mesheryctl response for no profile-id passed
	result1012output = "1012.golden"
)

func TestResultCmd(t *testing.T) {
	utils.SetupContextEnv(t)
	// utils.StartMockery(t)
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures", "result")
	testToken := filepath.Join(currDir, "fixtures", "auth.json")
	testdataDir := filepath.Join(currDir, "testdata", "result")
	resultURL := testContext.BaseURL + "/api/user/performance/profile/" + tempProfileID + "/results"

	tests := []tempTestStruct{
		{"standard results output", []string{"result", tempProfileID}, []utils.MockURL{{"GET", resultURL, result1001, 200}}, result1001output, testToken, false},
		{"results searching kuma", []string{"result", tempProfileID, "kuma"}, []utils.MockURL{{"GET", resultURL, result1002, 200}}, result1002output, testToken, false},
		{"results searching no mesh", []string{"result", tempProfileID, "no", "mesh"}, []utils.MockURL{{"GET", resultURL, result1003, 200}}, result1003output, testToken, false},
		{"standard results in expand output", []string{"result", tempProfileID, "--expand"}, []utils.MockURL{{"GET", resultURL, result1001, 200}}, result1004output, testToken, false},
		{"Unmarshal error", []string{"result", tempProfileID}, []utils.MockURL{{"GET", resultURL, result1005, 200}}, result1010output, testToken, true},
		{"failing add authentication test", []string{"result", tempProfileID}, []utils.MockURL{}, result1011output, testToken + "invalid-path", true},
		{"Server Error 500", []string{"result", tempProfileID}, []utils.MockURL{{"GET", resultURL, result1006, 500}}, result1009output, testToken, true},
		{"No profile passed", []string{"result"}, []utils.MockURL{}, result1012output, testToken, true},
	}

	testsforLogrusOutputs := []tempTestStruct{
		{"No results found", []string{"result", tempProfileID, "--expand"}, []utils.MockURL{{"GET", resultURL, result1004, 200}}, result1005output, testToken, false},
		{"standard results in json output", []string{"result", tempProfileID, "-o", "json"}, []utils.MockURL{{"GET", resultURL, result1001, 200}}, result1006output, testToken, false},
		{"standard results in yaml output", []string{"result", tempProfileID, "-o", "yaml"}, []utils.MockURL{{"GET", resultURL, result1001, 200}}, result1007output, testToken, false},
		{"invalid output format", []string{"result", "-o", tempProfileID, "invalid"}, []utils.MockURL{{"GET", resultURL, result1001, 200}}, result1008output, testToken, true},
	}

	// Run tests in list format
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tokenPath = testToken

			for _, mock := range tt.URLs {
				apiResponse := utils.NewGoldenFile(t, mock.Response, fixturesDir).Load()
				httpmock.RegisterResponder(mock.Method, mock.URL,
					httpmock.NewStringResponder(mock.ResponseCode, apiResponse))
			}

			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Grab console prints
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PerfCmd.SetArgs(tt.Args)
			PerfCmd.SetOutput(rescueStdout)
			err := PerfCmd.Execute()
			if err != nil {
				if tt.ExpectError {
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()

					utils.Equals(t, expectedResponse, err.Error())
					return
				}
				t.Error(err)
			}

			w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = rescueStdout

			// response being printed in console
			actualResponse := string(out)
			// write it in file
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			utils.Equals(t, expectedResponse, actualResponse)
		})
	}

	// Run tests in list format
	for _, tt := range testsforLogrusOutputs {
		t.Run(tt.Name, func(t *testing.T) {
			tokenPath = tt.Token

			for _, mock := range tt.URLs {
				apiResponse := utils.NewGoldenFile(t, mock.Response, fixturesDir).Load()
				httpmock.RegisterResponder(mock.Method, mock.URL,
					httpmock.NewStringResponder(mock.ResponseCode, apiResponse))
			}

			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			b := utils.SetupLogrusGrabTesting(t)

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
					return
				}
				t.Error(err)
			}

			// response being printed in console
			actualResponse := b.String()
			// write it in file
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			utils.Equals(t, expectedResponse, actualResponse)
		})
	}

	// stop mock server
	utils.StopMockery(t)
}
