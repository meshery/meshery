package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestGetCurrentContext(t *testing.T) {
	tests := []struct {
		name        string
		cfg         MesheryCtlConfig
		wantErr     bool
		errContains string
		wantName    string
	}{
		{
			name: "Valid current-context exists",
			cfg: MesheryCtlConfig{
				CurrentContext: "local",
				Contexts: map[string]Context{
					"local": {
						Endpoint: "http://localhost:9081",
						Platform: "docker",
						Channel:  "stable",
						Version:  "latest",
						Provider: "Local",
					},
				},
			},
			wantErr:  false,
			wantName: "local",
		},
		{
			name: "Empty current-context returns error without exit",
			cfg: MesheryCtlConfig{
				CurrentContext: "",
				Contexts: map[string]Context{
					"local": {Endpoint: "http://localhost:9081"},
				},
			},
			wantErr:     true,
			errContains: "Valid context is not available in meshconfig",
		},
		{
			name: "Current context does not exist in contexts map",
			cfg: MesheryCtlConfig{
				CurrentContext: "nonexistent",
				Contexts: map[string]Context{
					"local": {Endpoint: "http://localhost:9081"},
				},
			},
			wantErr:     true,
			errContains: "current context nonexistent does not exist",
		},
		{
			name: "Nil contexts map with current context set",
			cfg: MesheryCtlConfig{
				CurrentContext: "local",
				Contexts:       nil,
			},
			wantErr:     true,
			errContains: "current context local does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := tt.cfg.GetCurrentContext()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("GetCurrentContext() expected error, got nil")
				}
				if ctx != nil {
					t.Errorf("GetCurrentContext() expected nil context on error, got %v", ctx)
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetCurrentContext() error = %q, want substring %q", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Fatalf("GetCurrentContext() unexpected error: %v", err)
				}
				if ctx == nil {
					t.Fatalf("GetCurrentContext() expected valid context, got nil")
				}
				if ctx.Endpoint != tt.cfg.Contexts[tt.wantName].Endpoint {
					t.Errorf("GetCurrentContext() endpoint = %q, want %q", ctx.Endpoint, tt.cfg.Contexts[tt.wantName].Endpoint)
				}
			}
		})
	}
}

func TestGetContext(t *testing.T) {
	tests := []struct {
		name         string
		cfg          MesheryCtlConfig
		queryName    string
		wantErr      bool
		errContains  string
		wantEndpoint string
	}{
		{
			name: "Valid named context exists",
			cfg: MesheryCtlConfig{
				Contexts: map[string]Context{
					"k8s-ctx": {
						Endpoint: "http://localhost:9081",
						Platform: "kubernetes",
					},
				},
			},
			queryName:    "k8s-ctx",
			wantErr:      false,
			wantEndpoint: "http://localhost:9081",
		},
		{
			name: "Named context does not exist",
			cfg: MesheryCtlConfig{
				Contexts: map[string]Context{
					"k8s-ctx": {Endpoint: "http://localhost:9081"},
				},
			},
			queryName:   "missing-ctx",
			wantErr:     true,
			errContains: "context missing-ctx does not exist",
		},
		{
			name: "Empty context name returns error",
			cfg: MesheryCtlConfig{
				Contexts: map[string]Context{
					"k8s-ctx": {Endpoint: "http://localhost:9081"},
				},
			},
			queryName:   "",
			wantErr:     true,
			errContains: "context  does not exist",
		},
		{
			name: "Nil contexts map returns error",
			cfg: MesheryCtlConfig{
				Contexts: nil,
			},
			queryName:   "any-ctx",
			wantErr:     true,
			errContains: "context any-ctx does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := tt.cfg.GetContext(tt.queryName)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("GetContext(%q) expected error, got nil", tt.queryName)
				}
				if ctx != nil {
					t.Errorf("GetContext(%q) expected nil context on error, got %v", tt.queryName, ctx)
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetContext(%q) error = %q, want substring %q", tt.queryName, err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Fatalf("GetContext(%q) unexpected error: %v", tt.queryName, err)
				}
				if ctx == nil {
					t.Fatalf("GetContext(%q) expected valid context, got nil", tt.queryName)
				}
				if ctx.Endpoint != tt.wantEndpoint {
					t.Errorf("GetContext(%q) endpoint = %q, want %q", tt.queryName, ctx.Endpoint, tt.wantEndpoint)
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
