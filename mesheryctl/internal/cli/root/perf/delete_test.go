package perf

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestDeleteCmd(t *testing.T) {
	testContext := utils.NewTestHelper(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	profileURL := testContext.BaseURL + "/api/performance/profiles"
	profileID := "a947dff3-1415-4ca6-8922-e35b3232bd78"
	deleteURL := profileURL + "/" + profileID

	listTests := []utils.MesheryMultiURLCommamdTest{
		{
			Name: "given a single matching profile when delete with yes then delete it",
			Args: []string{"delete", "istio-test", "--yes"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.delete.single.response.golden", ResponseCode: 200},
				{Method: "DELETE", URL: deleteURL, Response: "profile.delete.success.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "delete.success.output.golden",
		},
		{
			Name: "given several matches and one exact name when delete with yes then delete only the exact one",
			Args: []string{"delete", "istio-test", "--yes"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.delete.exact.response.golden", ResponseCode: 200},
				{Method: "DELETE", URL: deleteURL, Response: "profile.delete.success.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "delete.success.output.golden",
		},
	}

	errorTests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:          "given no argument when delete then throw no profile name error",
			Args:          []string{"delete"},
			ExpectError:   true,
			ExpectedError: ErrNoProfileName(),
		},
		{
			Name:          "given two arguments when delete then throw argument overflow error",
			Args:          []string{"delete", "one", "two"},
			ExpectError:   true,
			ExpectedError: ErrorArgumentOverflow(),
		},
		{
			Name: "given no matching profile when delete then throw no profile found error",
			Args: []string{"delete", "missing", "--yes"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.empty.response.golden", ResponseCode: 200},
			},
			ExpectError:   true,
			ExpectedError: ErrNoProfileFound(),
		},
		{
			Name: "given only a substring match when delete with yes then throw no profile found error",
			Args: []string{"delete", "istio", "--yes"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.delete.single.response.golden", ResponseCode: 200},
			},
			ExpectError:   true,
			ExpectedError: ErrNoProfileFound(),
		},
		{
			Name: "given server error when delete then throw error",
			Args: []string{"delete", "istio-test", "--yes"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.delete.single.response.golden", ResponseCode: 200},
				{Method: "DELETE", URL: deleteURL, Response: "profile.error.response.golden", ResponseCode: 400},
			},
			ExpectError:   true,
			ExpectedError: utils.ErrFailReqStatus(400, ""),
		},
	}

	utils.RunMesheryctlMultiURLTests(t, update, PerfCmd, listTests, currDir, "perf", resetVariables)
	utils.RunMesheryctlMultiURLTests(t, update, PerfCmd, errorTests, currDir, "perf", resetVariables)
}

func TestFetchAllPerformanceProfilesPagesThroughResults(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	utils.TokenFlag = utils.GetToken(t)

	defer func(old int) { pageSize = old }(pageSize)
	pageSize = 1

	profileURL := testContext.BaseURL + "/api/performance/profiles"
	ids := []string{"a947dff3-1415-4ca6-8922-e35b3232bd78", "f331c784-aba0-4944-8d7d-ae264221bcbf"}
	for page, id := range ids {
		body := fmt.Sprintf(`{"page":%d,"pageSize":1,"totalCount":2,"profiles":[{"id":"%s","name":"istio-%d"}]}`, page, id, page)
		httpmock.RegisterResponderWithQuery("GET", profileURL, fmt.Sprintf("pagesize=1&page=%d&search=istio", page),
			httpmock.NewStringResponder(200, body))
	}

	got, err := fetchAllPerformanceProfiles(testContext.BaseURL, "istio")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].Name != "istio-0" || got[1].Name != "istio-1" {
		t.Fatalf("expected both pages to be returned, got %+v", got)
	}
}
