package config

import (
	"fmt"
	"os"
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
	for _, test := range tests {
		token.SetLocation(test)
		got := token.GetLocation()
		want, err := os.UserHomeDir()
		want = want + "/.meshery/" + test
		if err != nil {
			t.Errorf("Fail")
		}
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetVersion(t *testing.T) {
	context := Context{}
	for _, test := range tests {
		context.SetVersion(test)
		got := context.GetVersion()
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
	context := Context{}
	for _, test := range tests {
		context.SetChannel(test)
		got := context.GetChannel()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetComponents(t *testing.T) {
	dummy := []string{"abc", "def", "ghi", "jkl", "mno", "pqr"}
	context := Context{"", "", "", dummy, "", "", ""}
	got := context.GetComponents()
	want := dummy
	for i, j := range got {
		if j != want[i] {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetPlatform(t *testing.T) {
	context := Context{}
	for _, test := range tests {
		context.SetPlatform(test)
		got := context.GetPlatform()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetToken(t *testing.T) {
	context := Context{}
	for _, test := range tests {
		context.SetToken(test)
		got := context.GetToken()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetEndpoint(t *testing.T) {
	context := Context{}
	for _, test := range tests {
		context.SetEndpoint(test)
		got := context.GetEndpoint()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
func TestGetCurrentContextName(t *testing.T) {
	for _, test := range tests {
		mesherycltconfig := MesheryCtlConfig{nil, test, nil}
		got := mesherycltconfig.GetCurrentContextName()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
func TestSetContext(t *testing.T) {
	for _, test := range tests {
		mesherycltconfig := MesheryCtlConfig{nil, test, nil}
		err := UpdateContextInConfig(nil, test)
		if err != nil {
			fmt.Print("Fail") //Internal:need to be fixed
		}
		got := mesherycltconfig.GetCurrentContextName()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
func TestSetEndpoint(t *testing.T) {
	context := Context{}
	for _, test := range tests {
		context.SetEndpoint(test)
		got := context.GetEndpoint()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetToken(t *testing.T) {
	context := Context{}
	for _, test := range tests {
		context.SetToken(test)
		got := context.GetToken()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetPlatform(t *testing.T) {
	context := Context{}
	for _, test := range tests {
		context.SetPlatform(test)
		got := context.GetPlatform()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetChannel(t *testing.T) {
	context := Context{}
	for _, test := range tests {
		context.SetChannel(test)
		got := context.GetChannel()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}
func TestSetVersion(t *testing.T) {
	context := Context{}
	for _, test := range tests {
		context.SetVersion(test)
		got := context.GetVersion()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetComponents(t *testing.T) {
	dummy := []string{"abc", "def", "ghi", "jkl", "mno", "pqr"}
	context := Context{"", "", "", dummy, "", "", ""}
	got := context.GetComponents()
	want := dummy
	for i, j := range got {
		if j != want[i] {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetProvider(t *testing.T) {
	for _, test := range tests {
		context := Context{"", "", "", nil, "", "", test}
		got := context.GetProvider()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestSetProvider(t *testing.T) {
	context := Context{}
	for _, test := range tests {
		context.SetProvider(test)
		got := context.GetProvider()
		want := test

		if got != want {
			t.Errorf("got %q want %q", got, want)
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
