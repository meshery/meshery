package perf

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
)

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
	profileURL := testContext.BaseURL + "/api/performance/profiles"

	listTests := []utils.MesheryMultiURLCommamdTest{
		{
			Name: "given no argument when profile then display profiles",
			Args: []string{"profile"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.list.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "profile.list.output.golden",
			ExpectError:      false,
		},
		{
			Name: "given argument istio when profile then display istio profiles",
			Args: []string{"profile", "istio"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.searchIstio.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "profile.searchIstio.output.golden",
			ExpectError:      false,
		},
		{
			Name: "given argument test 3 when profile then display test 3 profiles",
			Args: []string{"profile", "test", "3"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.searchTest3.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "profile.searchTest3.output.golden",
			ExpectError:      false,
		},
		{
			Name: "given json output flag when profile then display profiles in json format",
			Args: []string{"profile", "-o", "json"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.list.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "profile.json.output.golden",
			ExpectError:      false,
		},
		{
			Name: "given yaml output flag when profile then display profiles in yaml format",
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
			Name: "given no profiles found when profile then display no profiles message",
			Args: []string{"profile", "--view"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.empty.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "profile.noProfiles.output.golden",
			ExpectError:      false,
		},
		{
			Name: "given invalid output format when profile then throw error",
			Args: []string{"profile", "-o", "invalid"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.list.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    display.ErrInvalidOutputFormat("invalid"),
		},
		{
			Name: "given invalid API response when profile then throw error",
			Args: []string{"profile"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.invalidJSON.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: func() error {
				cmdUsed = "profile"

				// Replicate the exact JSON unmarshal error using the schema type.
				// Body matches the canonical camelCase wire form in the
				// `profile.invalidJSON.response.golden` fixture so the inner
				// json.Unmarshal error references `pageSize`, not `page_size`.
				var response perfprofile.PerformanceProfilePage
				innerErr := json.Unmarshal([]byte(`{"pageSize": "25"}`), &response)

				return utils.ErrLoadConfig(ErrFailUnmarshal(innerErr))
			}(),
		},
		{
			Name: "given server error 400 when profile then throw error",
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
			Name:             "given invalid authentication token when profile then throw error",
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

// TestFetchPerformanceProfilesEscapesSearch pins that the profile name reaches
// the server unchanged. It used to be sent with only spaces encoded, so "a&b"
// arrived as search=a and "50%" was dropped by the server's query parser.
func TestFetchPerformanceProfilesEscapesSearch(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	utils.TokenFlag = utils.GetToken(t)

	var gotSearch string
	httpmock.RegisterResponder("GET", testContext.BaseURL+"/api/performance/profiles",
		func(req *http.Request) (*http.Response, error) {
			gotSearch = req.URL.Query().Get("search")
			return httpmock.NewStringResponse(200, `{"profiles":[]}`), nil
		})

	for _, search := range []string{"a&b", "50%", "test 3", "a+b"} {
		gotSearch = ""
		if _, _, err := fetchPerformanceProfiles(testContext.BaseURL, search, 25, 0); err != nil {
			t.Fatalf("fetchPerformanceProfiles(%q): %v", search, err)
		}
		if gotSearch != search {
			t.Errorf("server received search=%q, want %q", gotSearch, search)
		}
	}
}
