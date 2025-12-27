package perf

import (
	"flag"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// sanitizePaths replaces absolute paths with a placeholder to make golden file comparisons portable
func sanitizePaths(output, basePath string) string {
	return strings.ReplaceAll(output, basePath, "<TEST_PATH>")
}

var update = flag.Bool("update", false, "update golden files")

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

type tempTestStruct struct {
	Name             string
	Args             []string
	URLs             []utils.MockURL
	ExpectedResponse string
	Token            string
	ExpectError      bool
	SanitizePath     bool  // set to true if output contains dynamic paths
	IsOutputGolden   bool  `default:"true"`
	ExpectedError    error `default:"nil"`
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

	tests := []tempTestStruct{
		{
			Name: "standard profiles output",
			Args: []string{"profile"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: profile1001, ResponseCode: 200},
			},
			ExpectedResponse: profile1001output,
			Token:            testToken,
			ExpectError:      false,
		},
		{
			Name: "profiles searching istio",
			Args: []string{"profile", "istio"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: profile1002, ResponseCode: 200},
			},
			ExpectedResponse: profile1002output,
			Token:            testToken,
			ExpectError:      false,
		},
		{
			Name: "profiles searching test 3",
			Args: []string{"profile", "test", "3"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: profile1003, ResponseCode: 200},
			},
			ExpectedResponse: profile1003output,
			Token:            testToken,
			ExpectError:      false,
		},
	}

	testsforLogrusOutputs := []tempTestStruct{
		{
			Name: "No profiles found",
			Args: []string{"profile", "--view"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: profile1004, ResponseCode: 200},
			},
			ExpectedResponse: profile1005output,
			Token:            testToken,
		},
		{
			Name: "standard profiles in json output",
			Args: []string{"profile", "-o", "json"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: profile1001, ResponseCode: 200},
			},
			ExpectedResponse: profile1006output,
			Token:            testToken,
		},
		{
			Name: "standard profiles in yaml output",
			Args: []string{"profile", "-o", "yaml"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: profile1001, ResponseCode: 200},
			},
			ExpectedResponse: profile1007output,
			Token:            testToken,
		},
		{
			Name: "invalid output format",
			Args: []string{"profile", "-o", "invalid"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: profile1001, ResponseCode: 200},
			},
			ExpectedResponse: profile1008output,
			Token:            testToken,
			ExpectError:      true,
		},
		{
			Name: "Unmarshal error",
			Args: []string{"profile"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: profile1005, ResponseCode: 200},
			},
			ExpectedResponse: profile1010output,
			Token:            testToken,
			ExpectError:      true,
		},
		{
			Name: "Server Error 400",
			Args: []string{"profile"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: profile1006, ResponseCode: 400},
			},
			ExpectedResponse: profile1009output,
			Token:            testToken,
			ExpectError:      true,
		},
		{
			Name:             "failing add authentication test",
			Args:             []string{"profile"},
			ExpectedResponse: profile1011output,
			Token:            testToken + "invalid-path",
			ExpectError:      true,
			SanitizePath:     true,
		},
	}

	// Run tests in list format
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.TokenFlag = tt.Token

			for _, mock := range tt.URLs {
				apiResponse := utils.NewGoldenFile(t, mock.Response, fixturesDir).Load()
				httpmock.RegisterResponder(mock.Method, mock.URL,
					httpmock.NewStringResponder(mock.ResponseCode, apiResponse))
			}

			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			_ = utils.SetupMeshkitLoggerTesting(t, false)

			// Grab console prints with proper cleanup
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Ensure stdout is always restored
			defer func() {
				os.Stdout = originalStdout
			}()

			PerfCmd.SetArgs(tt.Args)
			PerfCmd.SetOut(originalStdout)
			err := PerfCmd.Execute()

			// Close write end before reading
			w.Close()

			if err != nil {
				if tt.ExpectError {
					errStr := err.Error()
					if tt.SanitizePath {
						errStr = sanitizePaths(errStr, currDir)
					}
					if *update {
						golden.Write(errStr)
					}
					expectedResponse := golden.Load()
					utils.Equals(t, expectedResponse, errStr)
					resetVariables()
					return
				}
				t.Error(err)
			}

			out, _ := io.ReadAll(r)

			// response being printed in console
			actualResponse := string(out)
			// write it in file
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			cleanedActualResponse := utils.CleanStringFromHandlePagination(actualResponse)
			cleanedexpectedResponse := utils.CleanStringFromHandlePagination(expectedResponse)

			utils.Equals(t, cleanedexpectedResponse, cleanedActualResponse)
			resetVariables()
		})
	}

	// Run tests in list format
	for _, tt := range testsforLogrusOutputs {
		t.Run(tt.Name, func(t *testing.T) {
			utils.TokenFlag = tt.Token

			for _, mock := range tt.URLs {
				apiResponse := utils.NewGoldenFile(t, mock.Response, fixturesDir).Load()
				httpmock.RegisterResponder(mock.Method, mock.URL,
					httpmock.NewStringResponder(mock.ResponseCode, apiResponse))
			}

			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			b := utils.SetupMeshkitLoggerTesting(t, false)

			PerfCmd.SetArgs(tt.Args)
			PerfCmd.SetOut(b)
			err := PerfCmd.Execute()
			if err != nil {
				if tt.ExpectError {
					errStr := err.Error()
					if tt.SanitizePath {
						errStr = sanitizePaths(errStr, currDir)
					}
					if *update {
						golden.Write(errStr)
					}
					expectedResponse := golden.Load()
					utils.Equals(t, expectedResponse, errStr)
					resetVariables()
					return
				}
				t.Error(err)
			}

			// response being printed in console
			actualResponse := b.String()
			if tt.SanitizePath {
				actualResponse = sanitizePaths(actualResponse, currDir)
			}
			// write it in file
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			utils.Equals(t, expectedResponse, actualResponse)
			resetVariables()
		})
		t.Log("Profile Perf test Passed")
	}

	// stop mock server
	utils.StopMockery(t)
}
