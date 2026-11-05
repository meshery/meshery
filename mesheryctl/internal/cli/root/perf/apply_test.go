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
	"os"
	"path/filepath"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
)

func clearPerfFlags() {
	profileName = ""
	testURL = ""
	testName = ""
	testMesh = ""
	qps = ""
	concurrentRequests = ""
	testDuration = ""
	loadGenerator = ""
	loadTestBody = ""
	filePath = ""
}

// TestPerformanceTestParams_FallbackWith verifies that unset fields are populated from fallback,
// while already populated fields remain intact.
func TestPerformanceTestParams_FallbackWith(t *testing.T) {
	primary := PerformanceTestParams{
		Name: "primary-name",
		URL:  "https://primary.example.com",
		QPS:  "100",
	}

	fallback := PerformanceTestParams{
		Name:               "fallback-name",
		URL:                "https://fallback.example.com",
		Mesh:               "istio",
		QPS:                "50",
		ConcurrentRequests: "10",
		Duration:           "45s",
		LoadGenerator:      "fortio",
		Body:               "payload",
	}

	primary.FallbackWith(fallback)

	if primary.Name != "primary-name" {
		t.Errorf("expected Name to remain 'primary-name', got %q", primary.Name)
	}
	if primary.URL != "https://primary.example.com" {
		t.Errorf("expected URL to remain 'https://primary.example.com', got %q", primary.URL)
	}
	if primary.QPS != "100" {
		t.Errorf("expected QPS to remain '100', got %q", primary.QPS)
	}
	if primary.Mesh != "istio" {
		t.Errorf("expected Mesh to fallback to 'istio', got %q", primary.Mesh)
	}
	if primary.ConcurrentRequests != "10" {
		t.Errorf("expected ConcurrentRequests to fallback to '10', got %q", primary.ConcurrentRequests)
	}
	if primary.Duration != "45s" {
		t.Errorf("expected Duration to fallback to '45s', got %q", primary.Duration)
	}
	if primary.LoadGenerator != "fortio" {
		t.Errorf("expected LoadGenerator to fallback to 'fortio', got %q", primary.LoadGenerator)
	}
	if primary.Body != "payload" {
		t.Errorf("expected Body to fallback to 'payload', got %q", primary.Body)
	}
}

// TestPerformanceTestParams_Validate tests validation rules for URL, QPS, concurrency, and test body.
func TestPerformanceTestParams_Validate(t *testing.T) {
	// Valid params pass validation
	valid := PerformanceTestParams{
		URL:                "https://valid.example.com",
		QPS:                "10",
		ConcurrentRequests: "2",
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid params to pass validation, got: %v", err)
	}

	// Missing URL
	noURL := PerformanceTestParams{URL: ""}
	err := noURL.Validate()
	if err == nil || meshkiterrors.GetCode(err) != ErrNoTestURLCode {
		t.Errorf("expected ErrNoTestURLCode, got %v", err)
	}

	// Invalid URL
	invalidURL := PerformanceTestParams{URL: "invalid-url"}
	err = invalidURL.Validate()
	if err == nil || meshkiterrors.GetCode(err) != ErrNotValidURLCode {
		t.Errorf("expected ErrNotValidURLCode, got %v", err)
	}

	// Invalid QPS
	invalidQPS := PerformanceTestParams{
		URL: "https://valid.example.com",
		QPS: "not-an-int",
	}
	err = invalidQPS.Validate()
	if err == nil || meshkiterrors.GetCode(err) != ErrConvertQPSCode {
		t.Errorf("expected ErrConvertQPSCode, got %v", err)
	}

	// Invalid ConcurrentRequests
	invalidConcurrency := PerformanceTestParams{
		URL:                "https://valid.example.com",
		ConcurrentRequests: "invalid-concurrency",
	}
	err = invalidConcurrency.Validate()
	if err == nil || meshkiterrors.GetCode(err) != ErrConvertConcurrentRequestCode {
		t.Errorf("expected ErrConvertConcurrentRequestCode, got %v", err)
	}

	// Body pointing to existing file on disk
	tempDir := t.TempDir()
	bodyFilePath := filepath.Join(tempDir, "request-body.json")
	expectedContent := `{"test":"content"}`
	if err := os.WriteFile(bodyFilePath, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("failed to write body file: %v", err)
	}
	bodyFileParams := PerformanceTestParams{
		URL:  "https://valid.example.com",
		Body: bodyFilePath,
	}
	if err := bodyFileParams.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if bodyFileParams.Body != expectedContent {
		t.Errorf("expected body to be loaded from file, got %q", bodyFileParams.Body)
	}

	// Body as string payload (not a file)
	rawBodyParams := PerformanceTestParams{
		URL:  "https://valid.example.com",
		Body: `{"raw":"json"}`,
	}
	if err := rawBodyParams.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
	if rawBodyParams.Body != `{"raw":"json"}` {
		t.Errorf("expected raw string body to be preserved, got %q", rawBodyParams.Body)
	}
}

// TestParseSMPConfigFile verifies SMP configuration file loading and error handling.
func TestParseSMPConfigFile(t *testing.T) {
	// Nonexistent file
	_, err := parseSMPConfigFile("nonexistent.yaml")
	if err == nil || meshkiterrors.GetCode(err) != ErrReadFilepathCode {
		t.Errorf("expected ErrReadFilepathCode for missing file, got %v", err)
	}

	tempDir := t.TempDir()

	// Invalid YAML
	badYamlPath := filepath.Join(tempDir, "bad.yaml")
	if err := os.WriteFile(badYamlPath, []byte("test: [invalid: yaml"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	_, err = parseSMPConfigFile(badYamlPath)
	if err == nil || meshkiterrors.GetCode(err) != ErrFailUnmarshalFileCode {
		t.Errorf("expected ErrFailUnmarshalFileCode for bad YAML, got %v", err)
	}

	// Empty config with no clients
	emptyClientsPath := filepath.Join(tempDir, "empty-clients.yaml")
	if err := os.WriteFile(emptyClientsPath, []byte("mesh: istio\ntest:\n  name: test\n"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	_, err = parseSMPConfigFile(emptyClientsPath)
	if err == nil || meshkiterrors.GetCode(err) != ErrInvalidTestConfigFileCode {
		t.Errorf("expected ErrInvalidTestConfigFileCode for empty clients, got %v", err)
	}

	// Valid SMP file
	validPath := filepath.Join(tempDir, "valid-smp.yaml")
	validYAML := `
mesh: linkerd
test:
  name: smp-test
  duration: 2m
  clients:
    - endpointUrls:
        - https://smp.example.com
      rps: 300
      connections: 15
      loadGenerator: fortio
      body: "sample-body"
`
	if err := os.WriteFile(validPath, []byte(validYAML), 0644); err != nil {
		t.Fatalf("failed to write valid SMP file: %v", err)
	}
	params, err := parseSMPConfigFile(validPath)
	if err != nil {
		t.Fatalf("unexpected error parsing valid SMP file: %v", err)
	}
	if params.Name != "smp-test" {
		t.Errorf("expected Name 'smp-test', got %q", params.Name)
	}
	if params.URL != "https://smp.example.com" {
		t.Errorf("expected URL 'https://smp.example.com', got %q", params.URL)
	}
	if params.Mesh != "linkerd" {
		t.Errorf("expected Mesh 'linkerd', got %q", params.Mesh)
	}
	if params.QPS != "300" {
		t.Errorf("expected QPS '300', got %q", params.QPS)
	}
	if params.ConcurrentRequests != "15" {
		t.Errorf("expected ConcurrentRequests '15', got %q", params.ConcurrentRequests)
	}
	if params.Duration != "2m" {
		t.Errorf("expected Duration '2m', got %q", params.Duration)
	}
	if params.LoadGenerator != "fortio" {
		t.Errorf("expected LoadGenerator 'fortio', got %q", params.LoadGenerator)
	}
	if params.Body != "sample-body" {
		t.Errorf("expected Body 'sample-body', got %q", params.Body)
	}
}

// TestExtractProfileParams verifies mapping from PerformanceProfile into PerformanceTestParams.
func TestExtractProfileParams(t *testing.T) {
	profile := perfprofile.PerformanceProfile{
		Name:              "my-profile",
		Endpoints:         []string{"https://profile.example.com"},
		LoadGenerators:    []string{"wrk2"},
		ConcurrentRequest: 8,
		QPS:               75,
		Duration:          "15s",
		ServiceMesh:       "osm",
		RequestBody:       `{"profile":true}`,
	}

	params := extractProfileParams(profile)

	if params.URL != "https://profile.example.com" {
		t.Errorf("expected URL 'https://profile.example.com', got %q", params.URL)
	}
	if params.LoadGenerator != "wrk2" {
		t.Errorf("expected LoadGenerator 'wrk2', got %q", params.LoadGenerator)
	}
	if params.ConcurrentRequests != "8" {
		t.Errorf("expected ConcurrentRequests '8', got %q", params.ConcurrentRequests)
	}
	if params.QPS != "75" {
		t.Errorf("expected QPS '75', got %q", params.QPS)
	}
	if params.Duration != "15s" {
		t.Errorf("expected Duration '15s', got %q", params.Duration)
	}
	if params.Mesh != "osm" {
		t.Errorf("expected Mesh 'osm', got %q", params.Mesh)
	}
	if params.Body != `{"profile":true}` {
		t.Errorf("expected Body '{\"profile\":true}', got %q", params.Body)
	}
}

// TestResolvePerformanceParams_Precedence verifies the strict precedence: Flags > File > Profile > Defaults.
func TestResolvePerformanceParams_Precedence(t *testing.T) {
	clearPerfFlags()
	defer resetVariables()

	tempDir := t.TempDir()
	smpFile := filepath.Join(tempDir, "smp.yaml")
	smpContent := `
mesh: istio
test:
  name: file-name
  duration: 1m
  clients:
    - endpointUrls:
        - https://file.example.com
      rps: 200
      connections: 5
      loadGenerator: fortio
      body: "file-body"
`
	if err := os.WriteFile(smpFile, []byte(smpContent), 0644); err != nil {
		t.Fatalf("failed to write smp file: %v", err)
	}

	profile := perfprofile.PerformanceProfile{
		Name:              "profile-name",
		Endpoints:         []string{"https://profile.example.com"},
		LoadGenerators:    []string{"wrk2"},
		ConcurrentRequest: 2,
		QPS:               100,
		Duration:          "30s",
		ServiceMesh:       "linkerd",
		RequestBody:       "profile-body",
	}

	// Flag overrides URL and QPS
	flags := PerformanceTestParams{
		URL: "https://flag.example.com",
		QPS: "500",
	}

	resolved, err := resolvePerformanceParams(flags, smpFile, &profile, true)
	if err != nil {
		t.Fatalf("unexpected error resolving params: %v", err)
	}

	// Flags win over File and Profile
	if resolved.URL != "https://flag.example.com" {
		t.Errorf("expected URL from flag, got %q", resolved.URL)
	}
	if resolved.QPS != "500" {
		t.Errorf("expected QPS from flag, got %q", resolved.QPS)
	}

	// File wins over Profile when flag is unset
	if resolved.Name != "file-name" {
		t.Errorf("expected Name from file, got %q", resolved.Name)
	}
	if resolved.Mesh != "istio" {
		t.Errorf("expected Mesh from file, got %q", resolved.Mesh)
	}
	if resolved.Duration != "1m" {
		t.Errorf("expected Duration from file, got %q", resolved.Duration)
	}
	if resolved.ConcurrentRequests != "5" {
		t.Errorf("expected ConcurrentRequests from file, got %q", resolved.ConcurrentRequests)
	}
	if resolved.Body != "file-body" {
		t.Errorf("expected Body from file, got %q", resolved.Body)
	}
}

// TestResolvePerformanceParams_ProfileFallback verifies fallback to profile when no file is provided.
func TestResolvePerformanceParams_ProfileFallback(t *testing.T) {
	clearPerfFlags()
	defer resetVariables()

	profile := perfprofile.PerformanceProfile{
		Name:              "profile-target",
		Endpoints:         []string{"https://profile.example.com"},
		LoadGenerators:    []string{"wrk2"},
		ConcurrentRequest: 4,
		QPS:               80,
		Duration:          "45s",
		ServiceMesh:       "consul",
		RequestBody:       "profile-payload",
	}

	flags := PerformanceTestParams{
		Name: "custom-test-name",
	}

	resolved, err := resolvePerformanceParams(flags, "", &profile, true)
	if err != nil {
		t.Fatalf("unexpected error resolving profile params: %v", err)
	}

	if resolved.Name != "custom-test-name" {
		t.Errorf("expected Name from flags, got %q", resolved.Name)
	}
	if resolved.URL != "https://profile.example.com" {
		t.Errorf("expected URL from profile, got %q", resolved.URL)
	}
	if resolved.LoadGenerator != "wrk2" {
		t.Errorf("expected LoadGenerator from profile, got %q", resolved.LoadGenerator)
	}
	if resolved.ConcurrentRequests != "4" {
		t.Errorf("expected ConcurrentRequests from profile, got %q", resolved.ConcurrentRequests)
	}
	if resolved.QPS != "80" {
		t.Errorf("expected QPS from profile, got %q", resolved.QPS)
	}
	if resolved.Duration != "45s" {
		t.Errorf("expected Duration from profile, got %q", resolved.Duration)
	}
	if resolved.Mesh != "consul" {
		t.Errorf("expected Mesh from profile, got %q", resolved.Mesh)
	}
	if resolved.Body != "profile-payload" {
		t.Errorf("expected Body from profile, got %q", resolved.Body)
	}
}

// TestResolvePerformanceParams_ApplyDefaults verifies standard defaults are applied when unset across all sources.
func TestResolvePerformanceParams_ApplyDefaults(t *testing.T) {
	clearPerfFlags()
	defer resetVariables()

	flags := PerformanceTestParams{
		URL: "https://example.com",
	}

	resolved, err := resolvePerformanceParams(flags, "", nil, true)
	if err != nil {
		t.Fatalf("unexpected error resolving params with defaults: %v", err)
	}

	if resolved.Mesh != "None" {
		t.Errorf("expected default Mesh 'None', got %q", resolved.Mesh)
	}
	if resolved.QPS != "0" {
		t.Errorf("expected default QPS '0', got %q", resolved.QPS)
	}
	if resolved.ConcurrentRequests != "1" {
		t.Errorf("expected default ConcurrentRequests '1', got %q", resolved.ConcurrentRequests)
	}
	if resolved.Duration != "30s" {
		t.Errorf("expected default Duration '30s', got %q", resolved.Duration)
	}
	if resolved.LoadGenerator != "fortio" {
		t.Errorf("expected default LoadGenerator 'fortio', got %q", resolved.LoadGenerator)
	}
	if resolved.Name == "" {
		t.Error("expected random Name to be generated, got empty")
	}
}

// TestPerformanceTestParams_ApplyToGlobals verifies syncing parameters back to package global variables.
func TestPerformanceTestParams_ApplyToGlobals(t *testing.T) {
	clearPerfFlags()
	defer resetVariables()

	params := PerformanceTestParams{
		Name:               "global-name",
		URL:                "https://global.example.com",
		Mesh:               "osm",
		QPS:                "350",
		ConcurrentRequests: "12",
		Duration:           "2m",
		LoadGenerator:      "fortio",
		Body:               "global-body",
	}

	params.ApplyToGlobals()

	if testName != "global-name" {
		t.Errorf("expected testName 'global-name', got %q", testName)
	}
	if testURL != "https://global.example.com" {
		t.Errorf("expected testURL 'https://global.example.com', got %q", testURL)
	}
	if testMesh != "osm" {
		t.Errorf("expected testMesh 'osm', got %q", testMesh)
	}
	if qps != "350" {
		t.Errorf("expected qps '350', got %q", qps)
	}
	if concurrentRequests != "12" {
		t.Errorf("expected concurrentRequests '12', got %q", concurrentRequests)
	}
	if testDuration != "2m" {
		t.Errorf("expected testDuration '2m', got %q", testDuration)
	}
	if loadGenerator != "fortio" {
		t.Errorf("expected loadGenerator 'fortio', got %q", loadGenerator)
	}
	if loadTestBody != "global-body" {
		t.Errorf("expected loadTestBody 'global-body', got %q", loadTestBody)
	}
}

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

// TestProfileMetadataFromFlags verifies the metadata built from --options and
// --cert-path, and that invalid values are rejected.
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

// TestCreatePerformanceProfile_RejectsMissingCertificate verifies that no profile
// is created when the --cert-path file can't be read.
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

// existingProfileSchedule is the schedule stored on the mocked existing profile.
const existingProfileSchedule = "2c1f0a3e-5d4b-4c6e-9a7f-3b8d2e1f4a60"

// existingProfileRequests records the requests perf apply sends for an existing
// profile, and the profile it saved.
type existingProfileRequests struct {
	order []string
	saved perfprofile.PerformanceProfile
}

// mockExistingProfile serves an existing profile. As on the server, the profile
// list leaves out the schedule, which only the single-profile read returns.
func mockExistingProfile(t *testing.T, baseURL string, profileID string) *existingProfileRequests {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	origToken := utils.TokenFlag
	utils.TokenFlag = filepath.Join(filepath.Dir(filename), "fixtures", "auth.json")
	t.Cleanup(func() { utils.TokenFlag = origToken })

	fields := `"id":"` + profileID + `","name":"existing","qps":7,"duration":"30s",` +
		`"endpoints":["https://meshery.io"],"loadGenerators":["fortio"],` +
		`"metadata":{"additional_options":"{\"old\":true}","kept":"value"}`
	listedProfile := `{` + fields + `}`
	fullProfile := `{` + fields + `,"schedule":"` + existingProfileSchedule + `"}`

	requests := &existingProfileRequests{}
	profilesURL := baseURL + "/api/performance/profiles"
	httpmock.RegisterResponder("GET", profilesURL,
		httpmock.NewStringResponder(http.StatusOK, `{"profiles":[`+listedProfile+`],"totalCount":1}`))
	httpmock.RegisterResponder("GET", profilesURL+"/"+profileID, func(*http.Request) (*http.Response, error) {
		requests.order = append(requests.order, "fetch")
		return httpmock.NewStringResponse(http.StatusOK, fullProfile), nil
	})
	httpmock.RegisterResponder("POST", profilesURL, func(req *http.Request) (*http.Response, error) {
		requests.order = append(requests.order, "save")
		if err := json.NewDecoder(req.Body).Decode(&requests.saved); err != nil {
			return nil, err
		}
		return httpmock.NewStringResponse(http.StatusOK, fullProfile), nil
	})
	httpmock.RegisterResponder("GET", baseURL+"/api/user/performance/profiles/"+profileID+"/run",
		func(*http.Request) (*http.Response, error) {
			requests.order = append(requests.order, "run")
			return httpmock.NewStringResponse(http.StatusOK, `{}`), nil
		})

	return requests
}

// TestApplyUpdatesExistingProfileBeforeRun verifies that --options reaches an
// existing profile before the test starts, because the server reads it from the
// saved profile, and that the rest of the stored profile is kept.
func TestApplyUpdatesExistingProfileBeforeRun(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	const profileID = "906f8876-33b5-4a97-906e-7a409d3b8ae9"
	requests := mockExistingProfile(t, testContext.BaseURL, profileID)

	resetVariables()
	defer resetVariables()
	testURL = "https://meshery.io"
	additionalOptions = `{"MethodOverride":"POST"}`

	if err := applyCmd.RunE(applyCmd, []string{"existing"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(requests.order, []string{"fetch", "save", "run"}) {
		t.Fatalf("expected the full profile to be fetched and saved before the run, got requests %v", requests.order)
	}
	saved := requests.saved
	if saved.ID.String() != profileID || saved.QPS != 7 {
		t.Errorf("expected the existing profile to be updated, got id=%s qps=%d", saved.ID, saved.QPS)
	}
	if saved.Schedule == nil || saved.Schedule.String() != existingProfileSchedule {
		t.Errorf("expected the stored schedule to be kept, got %v", saved.Schedule)
	}
	if got := saved.Metadata["additional_options"]; got != additionalOptions {
		t.Errorf("expected --options to replace the saved options, got %v", got)
	}
	if got := saved.Metadata["kept"]; got != "value" {
		t.Errorf("expected other metadata to be kept, got %v", got)
	}
}

// TestApplyDoesNotUpdateProfileWhenValidationFails verifies that a command that
// fails on its own inputs leaves the saved profile unchanged.
func TestApplyDoesNotUpdateProfileWhenValidationFails(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	requests := mockExistingProfile(t, testContext.BaseURL, "906f8876-33b5-4a97-906e-7a409d3b8ae9")

	resetVariables()
	defer resetVariables()
	testURL = "https://meshery.io"
	additionalOptions = `{"MethodOverride":"POST"}`
	testDuration = "500ms"

	if err := applyCmd.RunE(applyCmd, []string{"existing"}); err == nil {
		t.Fatal("expected an invalid duration to fail the command")
	}
	if len(requests.order) != 0 {
		t.Errorf("expected no profile update or run, got requests %v", requests.order)
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
