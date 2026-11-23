package perf

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkiterrors "github.com/meshery/meshkit/errors"
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
		{
			Name: "given an exact profile id when delete with yes then delete it directly via GET endpoint",
			Args: []string{"delete", profileID, "--yes"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: deleteURL, Response: "profile.get.single.response.golden", ResponseCode: 200},
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
			Name: "given several profiles with the exact name when delete with yes then throw error",
			Args: []string{"delete", "istio-test", "--yes"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.delete.duplicate.response.golden", ResponseCode: 200},
			},
			ExpectError:   true,
			ExpectedError: utils.ErrFlagsInvalid(fmt.Errorf("2 profiles are named 'istio-test'; run without --yes to choose one")),
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
		{
			Name: "given nonexistent profile uuid when delete with yes then throw no profile found error",
			Args: []string{"delete", profileID, "--yes"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: deleteURL, Response: "profile.error.response.golden", ResponseCode: 404},
			},
			ExpectError:   true,
			ExpectedError: ErrNoProfileFound(),
		},
		{
			Name: "given server error on profile uuid lookup when delete with yes then throw error",
			Args: []string{"delete", profileID, "--yes"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: deleteURL, Response: "profile.error.response.golden", ResponseCode: 500},
			},
			ExpectError:   true,
			ExpectedError: utils.ErrMesheryServerInternalError(errors.New("")),
		},
		{
			Name: "given bad request on profile uuid lookup when delete with yes then throw error",
			Args: []string{"delete", profileID, "--yes"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: deleteURL, Response: "profile.error.response.golden", ResponseCode: 400},
			},
			ExpectError:   true,
			ExpectedError: utils.ErrFailReqStatus(400, ""),
		},
	}

	utils.RunMesheryctlMultiURLTests(t, update, PerfCmd, listTests, currDir, "perf", resetVariables)
	utils.RunMesheryctlMultiURLTests(t, update, PerfCmd, errorTests, currDir, "perf", resetVariables)
}

func TestDeleteCmdRejectsIncompleteProfileLookup(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	previousToken := utils.TokenFlag
	utils.TokenFlag = utils.GetToken(t)
	t.Cleanup(func() {
		utils.StopMockery(t)
		utils.ResetCommandFlags(PerfCmd, t)
		resetVariables()
		if previousToken != "" {
			utils.TokenFlag = previousToken
		}
	})

	profileURL := testContext.BaseURL + "/api/performance/profiles"
	profileID := "a947dff3-1415-4ca6-8922-e35b3232bd78"
	firstPage := fmt.Sprintf(`{"totalCount":2,"profiles":[{"id":"%s","name":"istio-test"}]}`, profileID)
	httpmock.RegisterResponderWithQuery("GET", profileURL, "pagesize=25&page=0&search=istio-test", httpmock.NewStringResponder(200, firstPage))
	httpmock.RegisterResponderWithQuery("GET", profileURL, "pagesize=25&page=1&search=istio-test", httpmock.NewStringResponder(200, `{"totalCount":1,"profiles":[]}`))
	httpmock.RegisterResponder("DELETE", profileURL+"/"+profileID, func(_ *http.Request) (*http.Response, error) {
		t.Error("delete requested after incomplete profile lookup")
		return httpmock.NewStringResponse(200, ""), nil
	})

	PerfCmd.SetArgs([]string{"delete", "istio-test", "--yes"})
	if err := PerfCmd.Execute(); meshkiterrors.GetCode(err) != ErrProfileTraversalCode {
		t.Fatalf("expected profile traversal error, got %v", err)
	}
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

func TestDeleteCmdByUUIDInteractive(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	previousToken := utils.TokenFlag
	utils.TokenFlag = utils.GetToken(t)
	t.Cleanup(func() {
		utils.StopMockery(t)
		utils.ResetCommandFlags(PerfCmd, t)
		resetVariables()
		if previousToken != "" {
			utils.TokenFlag = previousToken
		}
	})

	profileID := "a947dff3-1415-4ca6-8922-e35b3232bd78"
	singleProfileJSON := fmt.Sprintf(`{"id":"%s","name":"istio-test"}`, profileID)
	getURL := testContext.BaseURL + "/api/performance/profiles/" + profileID

	t.Run("when confirmed then delete profile", func(t *testing.T) {
		httpmock.Reset()
		deleted := false
		httpmock.RegisterResponder("GET", getURL, httpmock.NewStringResponder(200, singleProfileJSON))
		httpmock.RegisterResponder("DELETE", getURL, func(_ *http.Request) (*http.Response, error) {
			deleted = true
			return httpmock.NewStringResponse(200, ""), nil
		})

		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		_, _ = w.WriteString("yes\n")
		_ = w.Close()
		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }()
		os.Stdin = r

		PerfCmd.SetArgs([]string{"delete", profileID})
		err = PerfCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !deleted {
			t.Errorf("expected DELETE request to be executed on confirmation")
		}
	})

	t.Run("when aborted then do not delete profile", func(t *testing.T) {
		httpmock.Reset()
		deleted := false
		httpmock.RegisterResponder("GET", getURL, httpmock.NewStringResponder(200, singleProfileJSON))
		httpmock.RegisterResponder("DELETE", getURL, func(_ *http.Request) (*http.Response, error) {
			deleted = true
			return httpmock.NewStringResponse(200, ""), nil
		})

		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		_, _ = w.WriteString("no\n")
		_ = w.Close()
		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }()
		os.Stdin = r

		PerfCmd.SetArgs([]string{"delete", profileID})
		err = PerfCmd.Execute()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if deleted {
			t.Errorf("expected DELETE request NOT to be executed when aborted")
		}
	})
}

func TestFetchPerformanceProfile(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	previousToken := utils.TokenFlag
	utils.TokenFlag = utils.GetToken(t)
	t.Cleanup(func() {
		utils.StopMockery(t)
		if previousToken != "" {
			utils.TokenFlag = previousToken
		}
	})

	profileID := "a947dff3-1415-4ca6-8922-e35b3232bd78"
	url := testContext.BaseURL + "/api/performance/profiles/" + profileID

	t.Run("successful profile lookup", func(t *testing.T) {
		httpmock.Reset()
		body := fmt.Sprintf(`{"id":"%s","name":"istio-test"}`, profileID)
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, body))

		profile, err := fetchPerformanceProfile(testContext.BaseURL, profileID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if profile.ID.String() != profileID || profile.Name != "istio-test" {
			t.Errorf("unexpected profile data: %+v", profile)
		}
	})

	t.Run("profile not found returns ErrNoProfileFound", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(404, `{"message":"not found"}`))

		_, err := fetchPerformanceProfile(testContext.BaseURL, profileID)
		if meshkiterrors.GetCode(err) != ErrNoProfileFoundCode {
			t.Fatalf("expected ErrNoProfileFoundCode, got %v", err)
		}
	})

	t.Run("server error returns ErrMesheryServerInternalError", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(500, "internal error"))

		_, err := fetchPerformanceProfile(testContext.BaseURL, profileID)
		if meshkiterrors.GetCode(err) != utils.ErrMesheryServerInternalErrorCode {
			t.Fatalf("expected ErrMesheryServerInternalErrorCode, got %v", err)
		}
	})

	t.Run("bad request error returns ErrFailReqStatus", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(400, "bad request"))

		_, err := fetchPerformanceProfile(testContext.BaseURL, profileID)
		if meshkiterrors.GetCode(err) != utils.ErrFailReqStatusCode {
			t.Fatalf("expected ErrFailReqStatusCode, got %v", err)
		}
	})

	t.Run("invalid json response returns ErrFailUnmarshal", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, "not-valid-json"))

		_, err := fetchPerformanceProfile(testContext.BaseURL, profileID)
		if meshkiterrors.GetCode(err) != ErrFailUnmarshalCode {
			t.Fatalf("expected ErrFailUnmarshalCode, got %v", err)
		}
	})
}
