package perf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkiterrors "github.com/meshery/meshkit/errors"
	meshkitlogger "github.com/meshery/meshkit/logger"
	core "github.com/meshery/schemas/models/core"
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
	"github.com/sirupsen/logrus"
)

var tempProfileID = "a2a555cf-ae16-479c-b5d2-a35656ba741e"

// PerformanceResultsAPIResponse is a local struct for testing unmarshal errors.
//
// The JSON tag mirrors the real perfprofile.PerformanceResultPage.PageSize
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

	profileURL := testContext.BaseURL + "/api/performance/profiles"
	resultURL := testContext.BaseURL + "/api/user/performance/profiles/" + tempProfileID + "/results"

	listTests := []utils.MesheryMultiURLCommamdTest{
		{
			Name: "results page does not paginate profile lookup",
			Args: []string{"result", "abhishek", "--page", "2"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "profile.empty.response.golden", ResponseCode: 200},
				{Method: "GET", URL: profileURL + "?pagesize=25&page=0&search=abhishek", Response: "result.profile.response.golden", ResponseCode: 200},
				{Method: "GET", URL: resultURL + "?pagesize=25&page=1", Response: "result.list.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "result.list.output.golden",
			ExpectError:      false,
		},
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
		{
			Name: "result output tolerates empty load generators",
			Args: []string{"result", "John Doe"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: profileURL, Response: "result.profile.emptyLoadGenerators.response.golden", ResponseCode: 200},
				{Method: "GET", URL: resultURL, Response: "result.list.response.golden", ResponseCode: 200},
			},
			ExpectedResponse: "result.list.output.golden",
			ExpectError:      false,
		},
	}

	loggerTests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:             "negative page",
			Args:             []string{"result", "abhishek", "--page", "-1"},
			URLs:             []utils.MockURL{},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("invalid value for --page '-1'")),
		},
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

func TestFetchAllPerformanceProfiles(t *testing.T) {
	for _, tt := range []struct {
		name      string
		total     int
		pageSize  int
		failPage  int
		wantPages int
	}{
		{name: "more than 25 matching profiles", total: 26, pageSize: 25, failPage: -1, wantPages: 2},
		{name: "exactly one full page", total: 25, pageSize: 25, failPage: -1, wantPages: 1},
		{name: "exactly two full pages", total: 50, pageSize: 25, failPage: -1, wantPages: 2},
		{name: "three pages", total: 51, pageSize: 25, failPage: -1, wantPages: 3},
		{name: "server caps page size", total: 26, pageSize: 10, failPage: -1, wantPages: 3},
		{name: "no matching profiles", pageSize: 25, failPage: -1, wantPages: 1},
		{name: "later page fails", total: 26, pageSize: 25, failPage: 1, wantPages: 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			testContext := utils.InitTestEnvironment(t)
			t.Cleanup(func() { utils.StopMockery(t) })
			previousToken := utils.TokenFlag
			utils.TokenFlag = utils.GetToken(t)
			t.Cleanup(func() { utils.TokenFlag = previousToken })

			profiles := make([]perfprofile.PerformanceProfile, tt.total)
			for i := range profiles {
				profiles[i] = perfprofile.PerformanceProfile{
					ID:   core.Uuid(uuid.FromStringOrNil(fmt.Sprintf("00000000-0000-0000-0000-%012d", i+1))),
					Name: "matching profile",
				}
			}

			requests := 0
			httpmock.RegisterResponder("GET", testContext.BaseURL+"/api/performance/profiles", func(req *http.Request) (*http.Response, error) {
				query := req.URL.Query()
				page, err := strconv.Atoi(query.Get("page"))
				if err != nil || page != requests || requests >= tt.wantPages {
					t.Errorf("unexpected profile page %q on request %d", query.Get("page"), requests)
					return httpmock.NewStringResponse(500, "{}"), nil
				}
				requests++
				if query.Get("search") != "matching profile" || query.Get("pagesize") != "25" {
					t.Errorf("unexpected profile query: %v", query)
				}
				if page == tt.failPage {
					return httpmock.NewStringResponse(500, "{}"), nil
				}

				start := page * tt.pageSize
				end := min(start+tt.pageSize, len(profiles))
				return httpmock.NewJsonResponse(200, perfprofile.PerformanceProfilePage{
					Page:       page,
					PageSize:   tt.pageSize,
					TotalCount: tt.total,
					Profiles:   profiles[start:end],
				})
			})

			got, err := fetchAllPerformanceProfiles(testContext.BaseURL, "matching%20profile")
			if tt.failPage >= 0 {
				if err == nil || len(got) != 0 {
					t.Fatalf("expected an error without partial profiles, got %d profiles and %v", len(got), err)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if len(got) != len(profiles) {
					t.Fatalf("expected all %d matching profiles, got %d", len(profiles), len(got))
				}
				for i := range profiles {
					if got[i].ID != profiles[i].ID {
						t.Errorf("profile %d: expected ID %s, got %s", i, profiles[i].ID, got[i].ID)
					}
				}
			}
			if requests != tt.wantPages {
				t.Errorf("expected %d profile page requests, got %d", tt.wantPages, requests)
			}
		})
	}
}

func TestFetchAllPerformanceProfilesRejectsUnstableTraversal(t *testing.T) {
	profiles := make([]perfprofile.PerformanceProfile, 51)
	for i := range profiles {
		profiles[i] = perfprofile.PerformanceProfile{
			ID:   core.Uuid(uuid.FromStringOrNil(fmt.Sprintf("00000000-0000-0000-0000-%012d", i+1))),
			Name: "matching profile",
		}
	}
	page := func(total, start, end int) perfprofile.PerformanceProfilePage {
		return perfprofile.PerformanceProfilePage{PageSize: 25, TotalCount: total, Profiles: profiles[start:end]}
	}
	for _, tt := range []struct {
		name  string
		pages []perfprofile.PerformanceProfilePage
	}{
		{name: "deletion leaves an empty page", pages: []perfprofile.PerformanceProfilePage{page(26, 0, 25), page(25, 26, 26)}},
		{name: "deletion leaves a nonempty page", pages: []perfprofile.PerformanceProfilePage{page(51, 0, 25), page(50, 26, 51)}},
		{name: "new match changes the count", pages: []perfprofile.PerformanceProfilePage{page(26, 0, 25), page(27, 25, 26)}},
		{name: "reordering repeats a profile", pages: []perfprofile.PerformanceProfilePage{page(26, 0, 25), page(26, 24, 25)}},
		{name: "duplicate within one page", pages: []perfprofile.PerformanceProfilePage{{TotalCount: 2, Profiles: []perfprofile.PerformanceProfile{profiles[0], profiles[0]}}}},
		{name: "empty page before initial total", pages: []perfprofile.PerformanceProfilePage{page(26, 0, 25), page(26, 26, 26)}},
		{name: "too many profiles for total", pages: []perfprofile.PerformanceProfilePage{page(1, 0, 2)}},
		{name: "negative total", pages: []perfprofile.PerformanceProfilePage{page(-1, 0, 0)}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			testContext := utils.InitTestEnvironment(t)
			t.Cleanup(func() { utils.StopMockery(t) })
			previousToken := utils.TokenFlag
			utils.TokenFlag = utils.GetToken(t)
			t.Cleanup(func() { utils.TokenFlag = previousToken })

			requests := 0
			httpmock.RegisterResponder("GET", testContext.BaseURL+"/api/performance/profiles", func(req *http.Request) (*http.Response, error) {
				pageNumber, err := strconv.Atoi(req.URL.Query().Get("page"))
				if err != nil || pageNumber != requests || requests >= len(tt.pages) {
					t.Errorf("unexpected profile page %q on request %d", req.URL.Query().Get("page"), requests)
					return httpmock.NewStringResponse(500, "{}"), nil
				}
				response := tt.pages[requests]
				response.Page = pageNumber
				requests++
				return httpmock.NewJsonResponse(200, response)
			})

			got, err := fetchAllPerformanceProfiles(testContext.BaseURL, "matching%20profile")
			if err == nil || len(got) != 0 {
				t.Fatalf("expected traversal error without partial profiles, got %d profiles and %v", len(got), err)
			}
			if meshkiterrors.GetCode(err) != ErrProfileTraversalCode {
				t.Errorf("expected profile traversal error, got %v", err)
			}
			if requests != len(tt.pages) {
				t.Errorf("expected %d profile requests, got %d", len(tt.pages), requests)
			}
		})
	}
}

func TestResultCmdReportsIncompleteProfileLookup(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	t.Cleanup(func() { utils.StopMockery(t) })
	t.Cleanup(resetVariables)
	previousToken := utils.TokenFlag
	utils.TokenFlag = utils.GetToken(t)
	t.Cleanup(func() {
		utils.ResetCommandFlags(PerfCmd, t)
		utils.TokenFlag = previousToken
	})
	var buf bytes.Buffer
	log, err := meshkitlogger.New("mesheryctl-test", meshkitlogger.Options{
		Format:      meshkitlogger.JsonLogFormat,
		LogLevel:    int(logrus.InfoLevel),
		Output:      &buf,
		ErrorOutput: &buf,
	})
	if err != nil {
		t.Fatal(err)
	}
	previousLog := utils.Log
	utils.Log = log
	t.Cleanup(func() { utils.Log = previousLog })
	previousOut, previousErr := PerfCmd.OutOrStdout(), PerfCmd.ErrOrStderr()
	PerfCmd.SetOut(&buf)
	PerfCmd.SetErr(&buf)
	t.Cleanup(func() {
		PerfCmd.SetOut(previousOut)
		PerfCmd.SetErr(previousErr)
	})

	requests := 0
	httpmock.RegisterResponder("GET", testContext.BaseURL+"/api/performance/profiles", func(req *http.Request) (*http.Response, error) {
		page, err := strconv.Atoi(req.URL.Query().Get("page"))
		if err != nil || page != requests || page > 1 {
			t.Errorf("unexpected profile page %q", req.URL.Query().Get("page"))
			return httpmock.NewStringResponse(500, "{}"), nil
		}
		requests++
		response := perfprofile.PerformanceProfilePage{Page: page, PageSize: 25, TotalCount: 26}
		if page == 0 {
			for i := 0; i < 25; i++ {
				response.Profiles = append(response.Profiles, perfprofile.PerformanceProfile{
					ID: core.Uuid(uuid.FromStringOrNil(fmt.Sprintf("00000000-0000-0000-0000-%012d", i+1))),
				})
			}
		}
		return httpmock.NewJsonResponse(200, response)
	})

	PerfCmd.SetArgs([]string{"result", "matching profile", "--page", "2"})
	err = PerfCmd.Execute()
	if meshkiterrors.GetCode(err) != ErrProfileTraversalCode {
		t.Fatalf("expected profile traversal error, got %v", err)
	}
	if requests != 2 {
		t.Errorf("expected two profile requests, got %d", requests)
	}
	if !strings.Contains(buf.String(), ErrProfileTraversalCode) || !strings.Contains(buf.String(), "Rerun mesheryctl perf result") {
		t.Errorf("expected structured error code and retry guidance in output, got %q", buf.String())
	}
}

func TestProfilesToStringArraysHandlesMissingLoadGenerators(t *testing.T) {
	t.Parallel()

	id := core.Uuid(uuid.Nil)
	data := profilesToStringArrays([]perfprofile.PerformanceProfile{
		{
			ID:             id,
			Name:           "test-profile",
			TotalResults:   3,
			LoadGenerators: nil,
		},
	})

	if len(data) != 1 {
		t.Fatalf("expected 1 row, got %d", len(data))
	}

	if len(data[0]) != 5 {
		t.Fatalf("expected 5 columns, got %d", len(data[0]))
	}

	if data[0][3] != "" {
		t.Fatalf("expected empty load generator cell, got %q", data[0][3])
	}
}
