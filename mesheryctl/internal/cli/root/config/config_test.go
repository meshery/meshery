package config

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

// var update = flag.Bool("update", false, "update golden files")

var tests = []string{"ab2q3^$serhj", "klg434%$^", "m234@#$n", "op$#%G", "q$@#$r", "st4@#$23", "uv$@#$", "$@#TGGB#wx", "y@#$FG@$z", "1234"}

func TestGetCommitSHA(t *testing.T) {
	for _, test := range tests {
		version := Version{"", test, ""}
		got := version.GetCommitSHA()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
func TestGetBuild(t *testing.T) {
	for _, test := range tests {
		version := Version{test, "", ""}
		got := version.GetBuild()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetLocation(t *testing.T) {
	token := Token{}
	mesheryFolder := ".meshery"

	for _, test := range tests {
		token.SetLocation(test)
		got := token.GetLocation()
		want, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("os.UserHomeDir() returned unexpected error: %v", err)
		}
		want = filepath.Join(want, mesheryFolder, test)
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetVersion(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetVersion(test)
		got := meshConfigCtx.GetVersion()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetName(t *testing.T) {
	token := Token{}
	for _, test := range tests {
		token.SetName(test)
		got := token.Name
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetChannel(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetChannel(test)
		got := meshConfigCtx.GetChannel()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetComponents(t *testing.T) {
	dummy := []string{"abc", "def", "ghi", "jkl", "mno", "pqr"}
	meshConfigCtx := Context{"", "", "", dummy, "", "", "", "", nil}
	got := meshConfigCtx.GetComponents()
	want := dummy
	for i, j := range got {
		if j != want[i] {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetPlatform(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetPlatform(test)
		got := meshConfigCtx.GetPlatform()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetToken(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetToken(test)
		got := meshConfigCtx.GetToken()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetEndpoint(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetEndpoint(test)
		got := meshConfigCtx.GetEndpoint()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
func TestGetCurrentContextName(t *testing.T) {
	for _, test := range tests {
		mesheryctlconfig := MesheryCtlConfig{nil, test, nil}
		got := mesheryctlconfig.GetCurrentContextName()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
func TestSetContext(t *testing.T) {
	// SetCurrentContext sets CurrentContext and validates it against known contexts.
	// We test only the setting behavior here using a context name that exists in the map.
	contextName := "local"
	ctx := Context{
		Endpoint: "http://localhost:9081",
		Platform: "docker",
		Token:    "Default",
		Channel:  "stable",
		Version:  "latest",
	}
	mesheryctlconfig := MesheryCtlConfig{
		Contexts:       map[string]Context{contextName: ctx},
		CurrentContext: "",
		Tokens:         nil,
	}

	err := mesheryctlconfig.SetCurrentContext(contextName)
	if err != nil {
		t.Fatalf("SetCurrentContext(%q) returned unexpected error: %v", contextName, err)
	}

	got := mesheryctlconfig.GetCurrentContextName()
	if got != contextName {
		t.Errorf("GetCurrentContextName() = %q, want %q", got, contextName)
	}
}
func TestSetEndpoint(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetEndpoint(test)
		got := meshConfigCtx.GetEndpoint()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetToken(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetToken(test)
		got := meshConfigCtx.GetToken()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetPlatform(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetPlatform(test)
		got := meshConfigCtx.GetPlatform()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetChannel(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetChannel(test)
		got := meshConfigCtx.GetChannel()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
func TestSetVersion(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetVersion(test)
		got := meshConfigCtx.GetVersion()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetComponents(t *testing.T) {
	dummy := []string{"abc", "def", "ghi", "jkl", "mno", "pqr"}
	meshConfigCtx := Context{"", "", "", dummy, "", "", "", "", nil}
	got := meshConfigCtx.GetComponents()
	want := dummy
	for i, j := range got {
		if j != want[i] {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetProvider(t *testing.T) {
	for _, test := range tests {
		meshConfigCtx := Context{"", "", "", nil, "", "", test, "", nil}
		got := meshConfigCtx.GetProvider()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetProvider(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetProvider(test)
		got := meshConfigCtx.GetProvider()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetOperatorStatus(t *testing.T) {
	for _, test := range tests {
		meshConfigCtx := Context{"", "", "", nil, "", "", "", test, nil}
		got := meshConfigCtx.GetOperatorStatus()
		want := test

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	}
}

func TestSetOperatorStatus(t *testing.T) {
	meshConfigCtx := Context{}
	for _, test := range tests {
		meshConfigCtx.SetOperatorStatus(test)
		got := meshConfigCtx.GetOperatorStatus()
		want := test

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	}
}

// roundTripFunc is an adapter to allow the use of an ordinary function as an http.RoundTripper.
type roundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip executes a single HTTP transaction, returning a Response for the provided Request.
func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// TestValidateVersion tests the version validation logic for contexts.
func TestValidateVersion(t *testing.T) {
	// 1. Test empty version defaults to "latest" and returns nil
	emptyCtx := Context{Version: ""}
	if err := emptyCtx.ValidateVersion(); err != nil {
		t.Fatalf("expected nil error for empty version, got: %v", err)
	}
	if emptyCtx.Version != "latest" {
		t.Fatalf("expected version to be set to 'latest', got: %s", emptyCtx.Version)
	}

	// 2. Test "latest" version returns nil without making HTTP call
	latestCtx := Context{Version: "latest"}
	if err := latestCtx.ValidateVersion(); err != nil {
		t.Fatalf("expected nil error for 'latest' version, got: %v", err)
	}

	// Preserve original client transport and restore after test
	origTransport := validateVersionClient.Transport
	defer func() {
		validateVersionClient.Transport = origTransport
	}()

	testCases := []struct {
		name          string
		version       string
		statusCode    int
		respErr       error
		wantErr       bool
		expectedCode  string
		expectedError string
	}{
		{
			name:       "valid release version 200 OK",
			version:    "v0.7.0",
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:          "invalid release version 404 Not Found",
			version:       "v999.999.999",
			statusCode:    http.StatusNotFound,
			wantErr:       true,
			expectedCode:  ErrInvalidVersionCode,
			expectedError: "version 'v999.999.999' is not a valid Meshery release",
		},
		{
			name:          "server error 500",
			version:       "v0.7.1",
			statusCode:    http.StatusInternalServerError,
			wantErr:       true,
			expectedCode:  ErrInvalidVersionCode,
			expectedError: "failed to validate Meshery release version v0.7.1 (server returned HTTP 500)",
		},
		{
			name:          "network error",
			version:       "v0.7.2",
			respErr:       errors.New("connection refused"),
			wantErr:       true,
			expectedError: "failed to make GET request",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validateVersionClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
				if tc.respErr != nil {
					return nil, tc.respErr
				}
				return &http.Response{
					StatusCode: tc.statusCode,
					Body:       io.NopCloser(strings.NewReader("")),
					Header:     make(http.Header),
				}, nil
			})

			ctx := Context{Version: tc.version}
			err := ctx.ValidateVersion()

			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for case %q, got nil", tc.name)
				}
				if tc.expectedCode != "" && meshkiterrors.GetCode(err) != tc.expectedCode {
					t.Fatalf("expected error code %q, got %q", tc.expectedCode, meshkiterrors.GetCode(err))
				}
				if !strings.Contains(err.Error(), tc.expectedError) {
					t.Fatalf("expected error containing %q, got %q", tc.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error for case %q, got: %v", tc.name, err)
				}
			}
		})
	}
}

// TODO: Shift Testing utility functions to meshkit so import cycle problems can be eliminated in future

// func TestChangePlatform(t *testing.T) {
// 	type args struct {
// 		contextName string
// 		platform    string
// 	}

// 	currDir := utils.GetBasePath(t)
// 	fixtureDir := currDir + "/fixtures"
// 	testConfigPath := fixtureDir + "/TestConfig.yaml"

// 	// Read and write to the test config file
// 	utils.SetupCustomContextEnv(t, testConfigPath)

// 	mctlCfg, _ := GetMesheryCtl(viper.GetViper())

// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 		golden  string
// 	}{
// 		{
// 			name:    "Update platform in gke context (valid context)",
// 			args:    args{contextName: "gke", platform: "testplatform"},
// 			wantErr: false,
// 			golden:  "changeplatform.valid.golden",
// 		},
// 		{
// 			name:    "Update platform in kubernetes context (invalid context)",
// 			args:    args{contextName: "kubernetes", platform: "testplatform"},
// 			wantErr: true,
// 			golden:  "changeplatform.invalid.golden",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := mctlCfg.SetCurrentContext(tt.args.contextName)
// 			if err != nil {
// 				if !tt.wantErr {
// 					t.Fatal("Error setting context", err)
// 				} else {
// 					// handles the case when an invalid context was intentionally supplied
// 					return
// 				}
// 			}

// 			currCtx, err := mctlCfg.GetCurrentContext()
// 			if err != nil {
// 				t.Fatal("Error processing context from config: ", err)
// 			}

// 			currCtx.SetPlatform(tt.args.platform)

// 			if err := UpdateContextInConfig(tt.args.contextName, currCtx, testConfigPath); err != nil && !tt.wantErr {
// 				t.Errorf("UpdateContextInConfig() error = %v, wantErr %v", err, tt.wantErr)
// 			}

// 			// Actual file contents
// 			actualContent, err := os.ReadFile(testConfigPath)
// 			if err != nil {
// 				t.Error("Error reading actual file contents: ", err)
// 			}

// 			actualFileContent := string(actualContent)

// 			// Expected file contents
// 			testdataDir := currDir + "/testdata"

// 			golden := utils.NewGoldenFile(t, tt.golden, testdataDir)
// 			if *update {
// 				golden.Write(actualFileContent)
// 			}
// 			expectedFileContent := golden.Load()

// 			if expectedFileContent != actualFileContent {
// 				t.Errorf("Expected file content \n[%v]\n and actual file content \n[%v]\n don't match", expectedFileContent, actualFileContent)
// 			}

// 			// Repopulating Expected yaml
// 			if err := utils.Populate(currDir+"/fixtures/original/TestConfig.yaml", testConfigPath); err != nil {
// 				t.Fatal(err, "Could not complete test. Unable to repopulate fixture")
// 			}
// 		})
// 	}
// }

// func TestChangeConfigEndpoint(t *testing.T) {
// 	// Setup path to test config file
// 	currDir := utils.GetBasePath(t)
// 	testConfigPath := currDir + "/fixtures/TestChangeEndpointConfig.yaml"

// 	utils.SetupCustomContextEnv(t, testConfigPath)

// 	mctlCfg, _ := GetMesheryCtl(viper.GetViper())

// 	tests := []struct {
// 		name            string
// 		ctxName         string
// 		endpointAddress string
// 		golden          string
// 		wantErr         bool
// 	}{
// 		{
// 			name:            "ChangeConfigEndpoint with platform docker",
// 			ctxName:         "local",
// 			endpointAddress: "http://localhost:55555",
// 			golden:          "changeconfigendpoint.expect.docker.golden",
// 			wantErr:         false,
// 		},
// 		{
// 			name:            "ChangeConfigEndpoint with platform kubernetes",
// 			ctxName:         "local2",
// 			endpointAddress: "http://localhost:44444",
// 			golden:          "changeconfigendpoint.expect.kubernetes.golden",
// 			wantErr:         false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := mctlCfg.SetCurrentContext(tt.ctxName)
// 			if err != nil {
// 				t.Fatal("error setting context", err)
// 			}

// 			currCtx, err := mctlCfg.GetCurrentContext()
// 			if err != nil {
// 				t.Fatal("error processing context from config", err)
// 			}

// 			currCtx.SetEndpoint(tt.endpointAddress)

// 			if err := UpdateContextInConfig(tt.ctxName, currCtx, testConfigPath); (err != nil) != tt.wantErr {
// 				t.Errorf("UpdateContextInConfig() error = %v, wantErr %v", err, tt.wantErr)
// 			}

// 			// Actual file contents
// 			actualContent, err := os.ReadFile(testConfigPath)
// 			if err != nil {
// 				t.Error("Error reading actual file contents: ", err)
// 			}

// 			actualFileContent := string(actualContent)

// 			// Expected file contents
// 			testdataDir := currDir + "/testdata"

// 			golden := utils.NewGoldenFile(t, tt.golden, testdataDir)
// 			if *update {
// 				golden.Write(actualFileContent)
// 			}
// 			expectedFileContent := golden.Load()

// 			if expectedFileContent != actualFileContent {
// 				t.Errorf("Expected file content \n[%v]\n and actual file content \n[%v]\n don't match", expectedFileContent, actualFileContent)
// 			}

// 			// Repopulating Expected yaml
// 			if err := utils.Populate(currDir+"/fixtures/platform/original/TestChangeEndpointConfig.yaml", testConfigPath); err != nil {
// 				t.Fatal(err, "Could not complete test. Unable to repopulate fixture")
// 			}
// 		})
// 	}
// }
