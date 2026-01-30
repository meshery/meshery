package perf

import (
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// PerformanceProfilesAPIResponse is a local struct for testing unmarshal errors.
type PerformanceProfilesAPIResponse struct {
	PageSize uint `json:"page_size"`
}

var update = flag.Bool("update", false, "update golden files")

func TestProfileCmd(t *testing.T) {
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	testToken := filepath.Join(currDir, "fixtures", "auth.json")
	profileURL := testContext.BaseURL + "/api/user/performance/profiles"

	listTests := []utils.MesheryMultiURLCommamdTest{
		{
			Name: "standard profiles output",
			Args: []string{"profile"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.list.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "profile.list.output.golden",
			ExpectError:      false,
		},
		{
			Name: "profiles searching istio",
			Args: []string{"profile", "istio"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.searchIstio.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "profile.searchIstio.output.golden",
			ExpectError:      false,
		},
		{
			Name: "profiles searching test 3",
			Args: []string{"profile", "test", "3"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.searchTest3.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "profile.searchTest3.output.golden",
			ExpectError:      false,
		},
		{
			Name: "standard profiles in json output",
			Args: []string{"profile", "-o", "json"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.list.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "profile.json.output.golden",
			ExpectError:      false,
		},
		{
			Name: "standard profiles in yaml output",
			Args: []string{"profile", "-o", "yaml"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.list.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "profile.yaml.output.golden",
			ExpectError:      false,
		},
	}

	loggerTests := []utils.MesheryMultiURLCommamdTest{
		{
			Name: "No profiles found",
			Args: []string{"profile", "--view"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.empty.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "profile.noProfiles.output.golden",
			ExpectError:      false,
		},
		{
			Name: "invalid output format",
			Args: []string{"profile", "-o", "invalid"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.list.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf(invalidOutputFormatMsg, "invalid")),
		},
		{
			Name: "Unmarshal error",
			Args: []string{"profile"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.invalidJSON.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: func() error {
				cmdUsed = "profile"

				// Replicate the exact JSON unmarshal error using local struct
				var response PerformanceProfilesAPIResponse
				innerErr := json.Unmarshal([]byte(`{"page_size": "25"}`), &response)

				return utils.ErrLoadConfig(ErrFailUnmarshal(innerErr))
			}(),
		},
		{
			Name: "Server Error 400",
			Args: []string{"profile"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.error.response.golden", ResponseCode: 400},
			},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: func() error {
				cmdUsed = "profile"

				body := ""
				return utils.ErrLoadConfig(utils.ErrFailReqStatus(400, body))
			}(),
		},
		{
			Name:             "failing add authentication test",
			Args:             []string{"profile"},
			ExpectedResponse: "",
			Token:            testToken + "invalid-path",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: func() error {
				tokenPath := testToken + "invalid-path"
				innerErr := fmt.Errorf("%s does not exist", tokenPath)
				return utils.ErrLoadConfig(utils.ErrAttachAuthToken(innerErr))
			}(),
		},
	}

	// Run tests in list format
	utils.RunMesheryctlMultipleURLsListTests(t, update, PerfCmd, listTests, currDir, "perf", resetVariables)

	// Run tests in logger format
	utils.RunMesheryctlMultiURLTests(t, update, PerfCmd, loggerTests, currDir, "perf", resetVariables)
}
