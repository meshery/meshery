// package perf

// import (
// 	"path/filepath"
// 	"runtime"
// 	"testing"

// 	"github.com/jarcoal/httpmock"
// 	"github.com/meshery/meshery/mesheryctl/pkg/utils"
// )

// var existingProfileID = "8f3daf25-e58e-4c59-8bf8-f474b76463ec"
// var newProfileID = "906f8876-33b5-4a97-906e-7a409d3b8ae9"

// // golden file responses
// var (
// 	// server returning existing profile "new"
// 	apply1001 = "1001.golden"
// 	// server response for no protocol added to url
// 	apply1003 = "1003.golden"
// 	// server response for creating new profile "test"
// 	apply1004 = "1004.golden"
// 	// server response for no profiles found
// 	apply1006 = "1006.golden"
// )

// var (
// 	// mesheryctl response for url not having protocol added
// 	apply1002output = "1002.golden"
// 	// mesheryctl response for invalid url
// 	apply1004output = "1004.golden"
// 	// mesheryctl response for no existing profile provided neither new profile-name
// 	apply1005output = "1005.golden"
// 	// mesheryctl response for no profiles found
// 	apply1006output = "1006.golden"
// )

// func TestApplyCmd(t *testing.T) {
// 	utils.SetupContextEnv(t)
// 	utils.StartMockery(t)
// 	testContext := utils.NewTestHelper(t)
// 	// get current directory
// 	_, filename, _, ok := runtime.Caller(0)
// 	if !ok {
// 		t.Fatal("Not able to get current working directory")
// 	}
// 	currDir := filepath.Dir(filename)
// 	fixturesDir := filepath.Join(currDir, "fixtures", "apply")
// 	testToken := filepath.Join(currDir, "fixtures", "auth.json")
// 	profileURL := testContext.BaseURL + "/api/performance/profiles"
// 	existingProfileRunTest := testContext.BaseURL + "/api/user/performance/profiles/" + existingProfileID + "/run"
// 	newProfileRunTest := testContext.BaseURL + "/api/user/performance/profiles/" + newProfileID + "/run"
// 	testdataDir := filepath.Join(currDir, "testdata", "apply")

// 	// test scenrios for fetching data
// 	tests := []tempTestStruct{
// 		{"Run Test with Existing profile", []string{"apply", "new"},
// 			[]utils.MockURL{
// 				{Method: "GET", URL: profileURL, Response: apply1001, ResponseCode: 200},
// 				{Method: "GET", URL: existingProfileRunTest, Response: apply1001, ResponseCode: 400},
// 			},
// 			apply1002output,
// 			testToken, false,
// 		},
// 		{"Run Test with Existing profile with --url", []string{"apply", "new", "--url", "https://www.google.com"},
// 			[]utils.MockURL{
// 				{Method: "GET", URL: profileURL, Response: apply1001, ResponseCode: 200},
// 				{Method: "GET", URL: existingProfileRunTest, Response: apply1001, ResponseCode: 400},
// 			},
// 			apply1002output,
// 			testToken, false,
// 		},
// 		{"Run Test with Existing profile with --url without protocol", []string{"apply", "new", "--url", "www.google.com"},
// 			[]utils.MockURL{
// 				{Method: "GET", URL: profileURL, Response: apply1001, ResponseCode: 200},
// 				{Method: "GET", URL: existingProfileRunTest, Response: apply1003, ResponseCode: 400},
// 			},
// 			apply1002output,
// 			testToken, true,
// 		},
// 		{"Run Test with new profile with --url without protocol", []string{"apply", "test", "--url", "www.google.com"},
// 			[]utils.MockURL{
// 				{Method: "POST", URL: profileURL, Response: apply1004, ResponseCode: 200},
// 				{Method: "GET", URL: newProfileRunTest, Response: apply1003, ResponseCode: 400},
// 			},
// 			apply1002output,
// 			testToken, true,
// 		},
// 		{"Run Test without profile-name and id", []string{"apply"},
// 			[]utils.MockURL{},
// 			apply1005output,
// 			testToken, true,
// 		},
// 		{"No profiles found with given name", []string{"apply", "new", "--yes"},
// 			[]utils.MockURL{
// 				{Method: "GET", URL: profileURL, Response: apply1006, ResponseCode: 200},
// 			},
// 			apply1006output,
// 			testToken, true,
// 		},
// 		{"Run Test with new profile with Invalid URL", []string{"apply", "test", "--url", "invalid-url"},
// 			[]utils.MockURL{},
// 			apply1004output,
// 			testToken, true,
// 		},
// 		{"Run Test with Existing profile with Invalid URL", []string{"apply", "new", "--url", "invalid-url"},
// 			[]utils.MockURL{{Method: "GET", URL: profileURL, Response: apply1001, ResponseCode: 200}},
// 			apply1004output,
// 			testToken, true,
// 		},
// 	}

// 	// Run tests
// 	for _, tt := range tests {
// 		t.Run(tt.Name, func(t *testing.T) {
// 			for _, url := range tt.URLs {
// 				apiResponse := utils.NewGoldenFile(t, url.Response, fixturesDir).Load()
// 				httpmock.RegisterResponder(url.Method, url.URL,
// 					httpmock.NewStringResponder(url.ResponseCode, apiResponse))
// 			}

// 			utils.TokenFlag = tt.Token
// 			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
// 			b := utils.SetupMeshkitLoggerTesting(t, false)

// 			PerfCmd.SetArgs(tt.Args)
// 			PerfCmd.SetOut(b)
// 			err := PerfCmd.Execute()
// 			if err != nil {
// 				if tt.ExpectError {
// 					if *update {
// 						golden.Write(err.Error())
// 					}
// 					expectedResponse := golden.Load()
// 					utils.Equals(t, expectedResponse, err.Error())
// 					resetVariables()
// 					return
// 				}
// 				t.Error(err)
// 			}

//				actualResponse := b.String()
//				if *update {
//					golden.Write(actualResponse)
//				}
//				expectedResponse := golden.Load()
//				utils.Equals(t, expectedResponse, actualResponse)
//				resetVariables()
//			})
//			t.Log("Apply Perf test Passed")
//		}
//		utils.StopMockery(t)
//	}
package perf

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
	"github.com/spf13/viper"
)

func TestSplitTestDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration string
		value    string
	}{
		{name: "empty duration defaults to 30s", duration: "", value: "30"},
		{name: "number without a unit means seconds", duration: "30", value: "30"},
		{name: "seconds", duration: "10s", value: "10"},
		{name: "surrounding spaces are trimmed", duration: " 5m ", value: "300"},
		{name: "hours", duration: "2h", value: "7200"},
		{name: "uppercase unit", duration: "5M", value: "300"},
		{name: "combined units", duration: "1h30m", value: "5400"},
		{name: "fractional value that is a whole number of seconds", duration: "1.5m", value: "90"},
		{name: "milliseconds that are a whole number of seconds", duration: "2000ms", value: "2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, unit, err := splitTestDuration(tt.duration)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if value != tt.value || unit != "s" {
				t.Fatalf("expected %q to be sent as %ss, got value=%q unit=%q", tt.duration, tt.value, value, unit)
			}
		})
	}
}

func TestSplitTestDurationRejectsInvalidDuration(t *testing.T) {
	for _, duration := range []string{"500ms", "1500ms", "1.5", "0s", "0", "-5s", "abc", "s", "5x"} {
		t.Run(duration, func(t *testing.T) {
			if value, unit, err := splitTestDuration(duration); err == nil {
				t.Fatalf("expected %q to be rejected, got value=%q unit=%q", duration, value, unit)
			}
		})
	}
}

func resetVariables() {
	// reset the variables after each test
	profileName = ""
	testURL = ""
	testName = ""
	testMesh = "none"
	qps = "0"
	concurrentRequests = "1"
	testDuration = "30s"
	loadGenerator = "fortio"
	filePath = ""
	outputFormatFlag = ""
	viewSingleProfile = false
	viewSingleResult = false
	additionalOptions = ""
	certPath = ""
}

func TestCreatePerformanceProfile_ServerErrors(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	origToken := utils.TokenFlag
	utils.TokenFlag = filepath.Join(currDir, "fixtures", "auth.json")
	defer func() { utils.TokenFlag = origToken }()

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Fatalf("failed to get mesheryctl config: %v", err)
	}

	profileURL := testContext.BaseURL + "/api/performance/profiles"

	tests := []struct {
		name          string
		statusCode    int
		responseBody  string
		expectedInErr string
	}{
		{
			name:          "server returns 500 internal server error",
			statusCode:    http.StatusInternalServerError,
			responseBody:  `{"error":"failed to save profile in database"}`,
			expectedInErr: "failed to save profile in database",
		},
		{
			name:          "server returns 400 with empty body uses status code",
			statusCode:    http.StatusBadRequest,
			responseBody:  "",
			expectedInErr: "400",
		},
		{
			name:          "server returns 401 unauthorized",
			statusCode:    http.StatusUnauthorized,
			responseBody:  `{"error":"unauthorized access"}`,
			expectedInErr: "unauthorized access",
		},
		{
			name:          "server returns 204 no content",
			statusCode:    http.StatusNoContent,
			responseBody:  "",
			expectedInErr: "204",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetVariables()
			profileName = "test-profile"
			testURL = "https://meshery.io"

			httpmock.RegisterResponder("POST", profileURL,
				httpmock.NewStringResponder(tt.statusCode, tt.responseBody))

			id, name, err := createPerformanceProfile(mctlCfg)
			if err == nil {
				t.Fatalf("expected error for status code %d, got nil", tt.statusCode)
			}
			if !strings.Contains(err.Error(), tt.expectedInErr) {
				t.Errorf("expected error to contain %q, got %q", tt.expectedInErr, err.Error())
			}
			if id != "" || name != "" {
				t.Errorf("expected empty profile ID and name on error, got id=%q, name=%q", id, name)
			}
		})
	}
}

func TestCreatePerformanceProfile_RejectsInvalidDuration(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Fatalf("failed to get mesheryctl config: %v", err)
	}

	profileURL := testContext.BaseURL + "/api/performance/profiles"
	httpmock.RegisterResponder("POST", profileURL, httpmock.NewStringResponder(http.StatusOK, `{}`))

	resetVariables()
	defer resetVariables()
	profileName = "test-profile"
	testURL = "https://meshery.io"
	testDuration = "500ms"

	if _, _, err := createPerformanceProfile(mctlCfg); err == nil {
		t.Fatal("expected an invalid duration to be rejected")
	}
	if calls := httpmock.GetCallCountInfo()["POST "+profileURL]; calls != 0 {
		t.Errorf("expected no profile to be created, got %d request(s)", calls)
	}
}

func TestProfileMetadataFromFlags(t *testing.T) {
	dir := t.TempDir()
	optionsFile := filepath.Join(dir, "options.json")
	certFile := filepath.Join(dir, "ca.pem")
	if err := os.WriteFile(optionsFile, []byte(`{"MethodOverride":"POST"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(certFile, []byte("certificate data"), 0o600); err != nil {
		t.Fatal(err)
	}

	certificate := map[string]interface{}{"name": "ca.pem", "file": "certificate data"}

	tests := []struct {
		name     string
		options  string
		cert     string
		expected map[string]interface{}
		wantErr  bool
	}{
		{name: "no flags", expected: map[string]interface{}{}},
		{
			name:     "options as a json string",
			options:  `{"MethodOverride":"PUT"}`,
			expected: map[string]interface{}{"additional_options": `{"MethodOverride":"PUT"}`},
		},
		{
			name:     "options from a file",
			options:  optionsFile,
			expected: map[string]interface{}{"additional_options": `{"MethodOverride":"POST"}`},
		},
		{name: "options that are not json", options: "not-json", wantErr: true},
		{
			name:     "certificate stored under the keys the server reads",
			cert:     certFile,
			expected: map[string]interface{}{"ca_certificate": certificate},
		},
		{name: "certificate file that does not exist", cert: filepath.Join(dir, "missing.pem"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetVariables()
			defer resetVariables()
			additionalOptions = tt.options
			certPath = tt.cert

			metadata, err := profileMetadataFromFlags()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got metadata %v", metadata)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(metadata, tt.expected) {
				t.Errorf("expected metadata %v, got %v", tt.expected, metadata)
			}
		})
	}
}

func TestCreatePerformanceProfile_RejectsMissingCertificate(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Fatalf("failed to get mesheryctl config: %v", err)
	}

	profileURL := testContext.BaseURL + "/api/performance/profiles"
	httpmock.RegisterResponder("POST", profileURL, httpmock.NewStringResponder(http.StatusOK, `{}`))

	resetVariables()
	defer resetVariables()
	profileName = "test-profile"
	testURL = "https://meshery.io"
	certPath = filepath.Join(t.TempDir(), "missing.pem")

	if _, _, err := createPerformanceProfile(mctlCfg); err == nil {
		t.Fatal("expected a missing certificate file to be rejected")
	}
	if calls := httpmock.GetCallCountInfo()["POST "+profileURL]; calls != 0 {
		t.Errorf("expected no profile to be created, got %d request(s)", calls)
	}
}

// TestApplyUpdatesExistingProfileBeforeRun verifies that --options reaches an
// existing profile before the test starts, because the server reads it from the
// saved profile, and that the rest of the profile is kept.
func TestApplyUpdatesExistingProfileBeforeRun(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	origToken := utils.TokenFlag
	utils.TokenFlag = filepath.Join(filepath.Dir(filename), "fixtures", "auth.json")
	defer func() { utils.TokenFlag = origToken }()

	const profileID = "906f8876-33b5-4a97-906e-7a409d3b8ae9"
	existingProfile := `{"id":"` + profileID + `","name":"existing","qps":7,"duration":"30s",` +
		`"endpoints":["https://meshery.io"],"loadGenerators":["fortio"],` +
		`"metadata":{"additional_options":"{\"old\":true}","kept":"value"}}`

	profilesURL := testContext.BaseURL + "/api/performance/profiles"
	httpmock.RegisterResponder("GET", profilesURL,
		httpmock.NewStringResponder(http.StatusOK, `{"profiles":[`+existingProfile+`],"totalCount":1}`))

	var requests []string
	var saved perfprofile.PerformanceProfile
	httpmock.RegisterResponder("POST", profilesURL, func(req *http.Request) (*http.Response, error) {
		requests = append(requests, "save")
		if err := json.NewDecoder(req.Body).Decode(&saved); err != nil {
			return nil, err
		}
		return httpmock.NewStringResponse(http.StatusOK, existingProfile), nil
	})
	httpmock.RegisterResponder("GET", testContext.BaseURL+"/api/user/performance/profiles/"+profileID+"/run",
		func(*http.Request) (*http.Response, error) {
			requests = append(requests, "run")
			return httpmock.NewStringResponse(http.StatusOK, `{}`), nil
		})

	resetVariables()
	defer resetVariables()
	testURL = "https://meshery.io"
	additionalOptions = `{"MethodOverride":"POST"}`

	if err := applyCmd.RunE(applyCmd, []string{"existing"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(requests, []string{"save", "run"}) {
		t.Fatalf("expected the profile to be saved before the run, got requests %v", requests)
	}
	if saved.ID.String() != profileID || saved.QPS != 7 {
		t.Errorf("expected the existing profile to be updated, got id=%s qps=%d", saved.ID, saved.QPS)
	}
	if got := saved.Metadata["additional_options"]; got != additionalOptions {
		t.Errorf("expected --options to replace the saved options, got %v", got)
	}
	if got := saved.Metadata["kept"]; got != "value" {
		t.Errorf("expected other metadata to be kept, got %v", got)
	}
}

func TestCreatePerformanceProfile_NullResponseGuardsAgainstPanic(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	origToken := utils.TokenFlag
	utils.TokenFlag = filepath.Join(currDir, "fixtures", "auth.json")
	defer func() { utils.TokenFlag = origToken }()

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Fatalf("failed to get mesheryctl config: %v", err)
	}

	profileURL := testContext.BaseURL + "/api/performance/profiles"

	resetVariables()
	profileName = "test-profile-null"
	testURL = "https://meshery.io"

	httpmock.RegisterResponder("POST", profileURL,
		httpmock.NewStringResponder(http.StatusOK, "null"))

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("createPerformanceProfile panicked on null response: %v", r)
		}
	}()

	id, name, err := createPerformanceProfile(mctlCfg)
	if err == nil {
		t.Fatal("expected error on null response, got nil")
	}
	if !strings.Contains(err.Error(), "empty performance profile") {
		t.Errorf("expected error to contain 'empty performance profile', got: %v", err)
	}
	if id != "" || name != "" {
		t.Errorf("expected empty profile ID and name, got id=%q, name=%q", id, name)
	}
}

func TestCreatePerformanceProfile_Success(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	origToken := utils.TokenFlag
	utils.TokenFlag = filepath.Join(currDir, "fixtures", "auth.json")
	defer func() { utils.TokenFlag = origToken }()

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Fatalf("failed to get mesheryctl config: %v", err)
	}

	profileURL := testContext.BaseURL + "/api/performance/profiles"

	successJSON := `{"id":"906f8876-33b5-4a97-906e-7a409d3b8ae9","name":"test-success-profile"}`

	tests := []struct {
		name       string
		statusCode int
	}{
		{name: "success with 200 OK", statusCode: http.StatusOK},
		{name: "success with 201 Created", statusCode: http.StatusCreated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetVariables()
			profileName = "test-success-profile"
			testURL = "https://meshery.io"

			httpmock.RegisterResponder("POST", profileURL,
				httpmock.NewStringResponder(tt.statusCode, successJSON))

			id, name, err := createPerformanceProfile(mctlCfg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if id != "906f8876-33b5-4a97-906e-7a409d3b8ae9" {
				t.Errorf("expected profile ID '906f8876-33b5-4a97-906e-7a409d3b8ae9', got %q", id)
			}
			if name != "test-success-profile" {
				t.Errorf("expected profile name 'test-success-profile', got %q", name)
			}
		})
	}
}
