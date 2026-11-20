package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
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

// TestGetEndpointTrimsTrailingSlashes guards issue #23389: an endpoint stored with a trailing slash
// (which `system context create --url http://host:9081/` accepts) made every request path start with
// "//api/...". Meshery Server redirects that, and the client re-sent a redirected DELETE as a GET.
func TestGetEndpointTrimsTrailingSlashes(t *testing.T) {
	cases := map[string]string{
		"http://localhost:9081/":  "http://localhost:9081",
		"http://localhost:9081//": "http://localhost:9081",
		"http://localhost:9081":   "http://localhost:9081",
		"":                        "",
	}
	for stored, want := range cases {
		ctx := Context{Endpoint: stored}
		if got := ctx.GetEndpoint(); got != want {
			t.Errorf("GetEndpoint() with endpoint %q = %q, want %q", stored, got, want)
		}
	}
}

func TestGetBaseMesheryURLTrimsTrailingSlash(t *testing.T) {
	mctl := MesheryCtlConfig{
		Contexts:       map[string]Context{"local": {Endpoint: "http://localhost:9081/"}},
		CurrentContext: "local",
	}
	if got, want := mctl.GetBaseMesheryURL(), "http://localhost:9081"; got != want {
		t.Errorf("GetBaseMesheryURL() = %q, want %q", got, want)
	}
}

func TestAddContextRejectsDotBeforeWriting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	original := []byte("contexts: {}\ncurrent-context: local\n")
	if err := os.WriteFile(path, original, 0600); err != nil {
		t.Fatal(err)
	}
	if err := AddContextToConfig("local.dev", Context{}, path, true, false); err == nil || !strings.Contains(err.Error(), "cannot contain") {
		t.Fatalf("expected invalid name error, got %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(original) {
		t.Fatalf("config changed: %s", got)
	}
}

const tokenRefsConfig = `contexts:
  local:
    endpoint: http://localhost:9081
    token: default
    platform: docker
    provider: Meshery
  prod:
    endpoint: http://localhost:9082
    token: prod
    platform: kubernetes
    provider: Meshery
  staging:
    endpoint: http://localhost:9083
    token: prod
    platform: kubernetes
    provider: Meshery
current-context: local
tokens:
  - name: default
    location: auth.json
  - name: prod
    location: prod-auth.json
  - name: spare
    location: spare-auth.json
`

// writeTokenRefsConfig writes tokenRefsConfig to a private meshconfig and
// resets the global viper the token helpers read through.
func writeTokenRefsConfig(t *testing.T) string {
	t.Helper()
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(configPath, []byte(tokenRefsConfig), 0o600); err != nil {
		t.Fatal(err)
	}
	viper.Reset()
	t.Cleanup(viper.Reset)
	return configPath
}

func readTokenRefsConfig(t *testing.T, configPath string) *MesheryCtlConfig {
	t.Helper()
	v := viper.New()
	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	cfg, err := GetMesheryCtl(v)
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}

func TestSetTokenToConfigRequiresExistingToken(t *testing.T) {
	t.Run("unknown token is rejected and the context keeps its token", func(t *testing.T) {
		configPath := writeTokenRefsConfig(t)

		err := SetTokenToConfig("prodd", configPath, "prod")
		if err == nil {
			t.Fatal("SetTokenToConfig() accepted a token that is not defined in meshconfig")
		}
		if !strings.Contains(err.Error(), `token "prodd" does not exist`) {
			t.Errorf("SetTokenToConfig() error = %q, want it to name the missing token", err)
		}
		if got := readTokenRefsConfig(t, configPath).Contexts["prod"].Token; got != "prod" {
			t.Errorf("context prod token = %q after a rejected set, want %q", got, "prod")
		}
	})

	t.Run("existing token is set on the given context only", func(t *testing.T) {
		configPath := writeTokenRefsConfig(t)

		if err := SetTokenToConfig("spare", configPath, "prod"); err != nil {
			t.Fatalf("SetTokenToConfig() error = %v", err)
		}
		cfg := readTokenRefsConfig(t, configPath)
		if got := cfg.Contexts["prod"].Token; got != "spare" {
			t.Errorf("context prod token = %q, want %q", got, "spare")
		}
		if got := cfg.Contexts["local"].Token; got != "default" {
			t.Errorf("context local token = %q, want %q", got, "default")
		}
	})
}

func TestDeleteTokenFromConfigRefusesTokenInUse(t *testing.T) {
	t.Run("token used by contexts is kept and the contexts are named", func(t *testing.T) {
		configPath := writeTokenRefsConfig(t)

		err := DeleteTokenFromConfig("prod", configPath)
		if err == nil {
			t.Fatal("DeleteTokenFromConfig() deleted a token that contexts still use")
		}
		if !strings.Contains(err.Error(), `token "prod" is still used by context(s) prod, staging`) {
			t.Errorf("DeleteTokenFromConfig() error = %q, want it to name both contexts", err)
		}
		if got := len(readTokenRefsConfig(t, configPath).Tokens); got != 3 {
			t.Errorf("meshconfig has %d tokens after a refused delete, want 3", got)
		}
	})

	t.Run("unused token is deleted", func(t *testing.T) {
		configPath := writeTokenRefsConfig(t)

		if err := DeleteTokenFromConfig("spare", configPath); err != nil {
			t.Fatalf("DeleteTokenFromConfig() error = %v", err)
		}
		cfg := readTokenRefsConfig(t, configPath)
		for _, tok := range cfg.Tokens {
			if tok.Name == "spare" {
				t.Fatal("token spare is still in meshconfig after delete")
			}
		}
		if got := len(cfg.Tokens); got != 2 {
			t.Errorf("meshconfig has %d tokens after delete, want 2", got)
		}
	})

	t.Run("unknown token still reports that it does not exist", func(t *testing.T) {
		configPath := writeTokenRefsConfig(t)

		err := DeleteTokenFromConfig("missing", configPath)
		if err == nil || !strings.Contains(err.Error(), "no such token exists") {
			t.Errorf("DeleteTokenFromConfig() error = %v, want \"no such token exists\"", err)
		}
	})
}
