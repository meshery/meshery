package perf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
)

var tempProfileID = "a2a555cf-ae16-479c-b5d2-a35656ba741e"

// golden file responses
var (
	// standard api response for abhishek profile
	result1000 = "1000.golden"
	// standard api response of 25 performance result
	result1001 = "1001.golden"
	// unable marshal response
	result1005 = "1005.golden"
	// empty response
	result1006 = "1006.golden"
)

// golden file mesheryctl outputs
var (
	// mesheryctl response of 25 performance result
	result1001output = "1001.golden"
	// mesheryctl response of 25 performance result in json output
	result1006output = "1006.golden"
	// mesheryctl response of 25 performance result in yaml output
	result1007output = "1007.golden"
)

func TestResultCmd(t *testing.T) {
	utils.SetupContextEnv(t)
	utils.StartMockery(t)
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

	profileURL := testContext.BaseURL + "/api/user/performance/profiles"
	resultURL := testContext.BaseURL + "/api/user/performance/profiles/" + tempProfileID + "/results"

	tests := []tempTestStruct{
		{
			Name: "standard results output",
			Args: []string{"result", "abhishek"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: result1000, ResponseCode: 200},
				{Method: "GET", URL: resultURL, Response: result1001, ResponseCode: 200},
			},
			ExpectedResponse: result1001output,
			Token:            testToken,
			ExpectError:      false,
		},
	}

	testsforLogrusOutputs := []tempTestStruct{
		{
			Name: "standard results in json output",
			Args: []string{"result", "abhishek", "-o", "json"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: result1000, ResponseCode: 200},
				{Method: "GET", URL: resultURL, Response: result1001, ResponseCode: 200},
			},
			ExpectedResponse: result1006output,
			Token:            testToken,
			ExpectError:      false,
		},
		{
			Name: "standard results in yaml output",
			Args: []string{"result", "abhishek", "-o", "yaml"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: result1000, ResponseCode: 200},
				{Method: "GET", URL: resultURL, Response: result1001, ResponseCode: 200},
			},
			ExpectedResponse: result1007output,
			Token:            testToken,
			ExpectError:      false,
		},
		{
			Name: "invalid output format",
			Args: []string{"result", "abhishek", "-o", "invalid"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: result1000, ResponseCode: 200},
				{Method: "GET", URL: resultURL, Response: result1001, ResponseCode: 200},
			},
			ExpectedResponse: "",
			Token:            testToken,
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    ErrInvalidOutputChoice(),
		},
		{
			Name: "Unmarshal error",
			Args: []string{"result", "abhishek"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: result1000, ResponseCode: 200},
				{Method: "GET", URL: resultURL, Response: result1005, ResponseCode: 200},
			},
			ExpectedResponse: "",
			Token:            testToken,
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: func() error {
				cmdUsed = "result"

				// Replicate the exact JSON unmarshal error
				var response models.PerformanceResultsAPIResponse
				innerErr := json.Unmarshal([]byte(`{"page_size": "25"}`), &response)

				return ErrPerformanceProfileResult(ErrFailUnmarshal(innerErr))
			}(),
		},
		{
			Name:             "failing add authentication test",
			Args:             []string{"result", "abhishek"},
			ExpectedResponse: "",
			Token:            testToken + "invalid-path",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: func() error {
				tokenPath := testToken + "invalid-path"
				innerErr := fmt.Errorf("%s does not exist", tokenPath)
				return utils.ErrAttachAuthToken(innerErr)
			}(),
		},
		{
			Name: "Server Error 400",
			Args: []string{"result", "abhishek"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: result1000, ResponseCode: 200},
				{Method: "GET", URL: resultURL, Response: result1006, ResponseCode: 400},
			},
			ExpectedResponse: "",
			Token:            testToken,
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: func() error {
				cmdUsed = "result"
				body := "{}"

				// The inner error from ErrFailReqStatus
				innerErr := utils.ErrFailReqStatus(400, body)

				return ErrPerformanceProfileResult(innerErr)
			}(),
		},
		{
			Name:             "No profile passed",
			Args:             []string{"result"},
			ExpectedResponse: "",
			Token:            testToken,
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    ErrNoProfileName(),
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

			// Skip golden file creation for error tests that use ExpectedError instead
			var golden *utils.GoldenFile
			if tt.ExpectedResponse != "" {
				golden = utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			}
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
					if tt.IsOutputGolden {

						if *update {
							golden.Write(err.Error())
						}
						expectedResponse := golden.Load()
						utils.Equals(t, expectedResponse, err.Error())
						resetVariables()
						return
					}

					utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
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

			// Skip golden file creation for error tests that use ExpectedError instead
			var golden *utils.GoldenFile
			if tt.ExpectedResponse != "" {
				golden = utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			}

			b := utils.SetupMeshkitLoggerTesting(t, false)

			PerfCmd.SetArgs(tt.Args)
			PerfCmd.SetOut(b)
			err := PerfCmd.Execute()
			if err != nil {
				if tt.ExpectError {
					if tt.IsOutputGolden {

						if *update {
							golden.Write(err.Error())
						}
						expectedResponse := golden.Load()
						utils.Equals(t, expectedResponse, err.Error())
						resetVariables()
						return
					}

					utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					resetVariables()
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
			resetVariables()
		})
		t.Log("Result perf test passed")
	}
	// stop mock server
	utils.StopMockery(t)
}
