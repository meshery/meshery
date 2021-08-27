package perf

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")

var tempProfileID = "ecddef09-7411-4b9e-b06c-fd55ff5debbc"

// golden file responses
var (
	// standard api response of 25 performance profile
	profile1001 = "1001.golden"
	// api response of performance profile with searched term "istio"
	profile1002 = "1002.golden"
	// api response of performance profile with searched term "test 3"
	profile1003 = "1003.golden"
	// api response of zero performance profile
	profile1004 = "1004.golden"
	// unable marshal response
	profile1005 = "1005.golden"
	// empty response
	profile1006 = "1006.golden"
)

// golden file mesheryctl outputs
var (
	// mesheryctl response of 25 performance profile
	profile1001output = "1001.golden"
	// mesheryctl response of performance profile with searched term "istio"
	profile1002output = "1002.golden"
	// mesheryctl response of performance profile with searched term "test 3"
	profile1003output = "1003.golden"
	// mesheryctl response of 25 performance profile in expanded output
	profile1004output = "1004.golden"
	// mesheryctl response when no profiles found
	profile1005output = "1005.golden"
	// mesheryctl response of 25 performance profile in json output
	profile1006output = "1006.golden"
	// mesheryctl response of 25 performance profile in yaml output
	profile1007output = "1007.golden"
	// mesheryctl response for invalid output format
	profile1008output = "1008.golden"
	// mesheryctl response for server response code error
	profile1009output = "1009.golden"
	// mesheryctl response for unmarshal error
	profile1010output = "1010.golden"
	// mesheryctl response for failing attach authentication
	profile1011output = "1011.golden"
)

type profileTestStruct struct {
	Name             string
	Args             []string
	URLs             []utils.MockURL
	ExpectedResponse string
	Token            string
	ExpectError      bool
}

func TestProfileCmd(t *testing.T) {
	utils.SetupContextEnv(t)
	utils.StartMockery(t)
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures", "profile")
	testToken := filepath.Join(currDir, "fixtures", "auth.json")
	testdataDir := filepath.Join(currDir, "testdata", "profile")
	profileURL := testContext.BaseURL + "/api/user/performance/profiles"

	// test scenrios for fetching data
	tests := []profileTestStruct{
		{"standard profiles output", []string{"profile"}, []utils.MockURL{{"GET", profileURL, profile1001, 200}}, profile1001output, testToken, false},
		{"profiles searching istio", []string{"profile", "istio"}, []utils.MockURL{{"GET", profileURL, profile1002, 200}}, profile1002output, testToken, false},
		{"profiles searching test 3", []string{"profile", "test", "3"}, []utils.MockURL{{"GET", profileURL, profile1003, 200}}, profile1003output, testToken, false},
		{"standard profiles in expand output", []string{"profile", "--expand"}, []utils.MockURL{{"GET", profileURL, profile1001, 200}}, profile1004output, testToken, false},
		{"Unmarshal error", []string{"profile"}, []utils.MockURL{{"GET", profileURL, profile1005, 200}}, profile1010output, testToken, true},
		{"failing add authentication test", []string{"profile"}, []utils.MockURL{}, profile1011output, testToken + "invalid-path", true},
		{"Server Error 500", []string{"profile"}, []utils.MockURL{{"GET", profileURL, profile1006, 500}}, profile1009output, testToken, true},
	}

	testsforLogrusOutputs := []profileTestStruct{
		{"No profiles found", []string{"profile", "--expand"}, []utils.MockURL{{"GET", profileURL, profile1004, 200}}, profile1005output, testToken, false},
		{"standard profiles in json output", []string{"profile", "-o", "json"}, []utils.MockURL{{"GET", profileURL, profile1001, 200}}, profile1006output, testToken, false},
		{"standard profiles in yaml output", []string{"profile", "-o", "yaml"}, []utils.MockURL{{"GET", profileURL, profile1001, 200}}, profile1007output, testToken, false},
		{"invalid output format", []string{"profile", "-o", "invalid"}, []utils.MockURL{{"GET", profileURL, profile1001, 200}}, profile1008output, testToken, true},
	}

	// Run tests in list format
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			tokenPath = tt.Token

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
