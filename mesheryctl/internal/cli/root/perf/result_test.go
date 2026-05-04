package perf

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

var tempProfileID = "a2a555cf-ae16-479c-b5d2-a35656ba741e"

// PerformanceResultsAPIResponse is a local struct for testing unmarshal errors.
//
// The JSON tag mirrors the real models.PerformanceResultsAPIResponse.PageSize
// tag (`pageSize`) so the error message produced by json.Unmarshal references
// the same field name the production code would surface.
type PerformanceResultsAPIResponse struct {
	PageSize uint `json:"pageSize"`
}

func TestResultCmd(t *testing.T) {
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	testToken := filepath.Join(currDir, "fixtures", "auth.json")

	profileURL := testContext.BaseURL + "/api/user/performance/profiles"
	resultURL := testContext.BaseURL + "/api/user/performance/profiles/" + tempProfileID + "/results"

	listTests := []utils.MesheryMultiURLCommamdTest{
		{
			Name: "standard results output",
			Args: []string{"result", "abhishek"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "result.profile.response.golden", ResponseCode: 200},
				{Method: "GET", URL: resultURL, Response: "result.list.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "result.list.output.golden",
			ExpectError:      false,
		},
		{
			Name: "standard results in json output",
			Args: []string{"result", "abhishek", "-o", "json"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "result.profile.response.golden", ResponseCode: 200},
				{Method: "GET", URL: resultURL, Response: "result.list.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "result.json.output.golden",
			ExpectError:      false,
		},
		{
			Name: "standard results in yaml output",
			Args: []string{"result", "abhishek", "-o", "yaml"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "result.profile.response.golden", ResponseCode: 200},
				{Method: "GET", URL: resultURL, Response: "result.list.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "result.yaml.output.golden",
			ExpectError:      false,
		},
	}

	loggerTests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:             "invalid output format",
			Args:             []string{"result", "abhishek", "-o", "invalid"},
			URLs:             []utils.MockURL{},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    display.ErrInvalidOutputFormat("invalid"),
		},
		{
			Name: "Unmarshal error",
			Args: []string{"result", "abhishek"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "result.profile.response.golden", ResponseCode: 200},
				{Method: "GET", URL: resultURL, Response: "result.invalidJSON.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: func() error {
				cmdUsed = "result"

				// Replicate the exact JSON unmarshal error. Body matches the
				// canonical camelCase wire form in the
				// `result.invalidJSON.response.golden` fixture so the inner
				// json.Unmarshal error references `pageSize`, not `page_size`.
				var response PerformanceResultsAPIResponse
				innerErr := json.Unmarshal([]byte(`{"pageSize": "25"}`), &response)

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
				{Method: "GET", URL: profileURL, Response: "result.profile.response.golden", ResponseCode: 200},
				{Method: "GET", URL: resultURL, Response: "result.error.response.golden", ResponseCode: 400},
			},
			ExpectedResponse: "",
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
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    ErrNoProfileName(),
		},
	}

	// Run tests in list format
	utils.RunMesheryctlMultipleURLsListTests(t, update, PerfCmd, listTests, currDir, "perf", resetVariables)

	// Run tests in logger format
	utils.RunMesheryctlMultiURLTests(t, update, PerfCmd, loggerTests, currDir, "perf", resetVariables)
}
