package system

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestContextViewCmd(t *testing.T) {
	resetVariables()
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/fixtures/.meshery/TestContext.yaml")

	tests := []utils.CmdTestInput{
		{
			Name:             "view for default context",
			Args:             []string{"context", "view"},
			ExpectedResponse: "view.golden",
		},
		{
			Name:             "view with specified context through context flag",
			Args:             []string{"context", "view", "--context", "local2"},
			ExpectedResponse: "viewWithContextExpected.golden",
		},
		{
			Name:             "Error for viewing a non-existing context",
			Args:             []string{"context", "view", "local3"},
			IsOutputGolden:   false,
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    ErrContextNotExists(fmt.Errorf("context `local3` does not exist")),
		},
		{
			Name:             "view with specified context as argument",
			Args:             []string{"context", "view", "local2"},
			ExpectedResponse: "viewWithContextExpected.golden",
		},
		{
			Name:             "view with all flag set",
			Args:             []string{"context", "view", "--all"},
			ExpectedResponse: "viewAllExpected.golden",
		},
	}

	mesheryctlflags.InitValidators(SystemCmd)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			defer buf.Reset()

			SystemCmd.SetOut(buf)
			SystemCmd.SetErr(buf)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()

			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					return
				}
				t.Fatal(err)
			}

			if tt.ExpectError {
				t.Fatalf("expected error, got nil")
			}

			testdataDir := filepath.Join(currDir, "testdata/context")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			actualResponse := buf.String()

			if *update {
				golden.Write(actualResponse)
			}

			expectedResponse := golden.Load()
			cleanedActualResponse := utils.CleanStringFromHandlePagination(actualResponse)
			cleanedExceptedResponse := utils.CleanStringFromHandlePagination(expectedResponse)

			utils.Equals(t, cleanedExceptedResponse, cleanedActualResponse)

		})
		t.Logf("List %s test", "context")
	}
}

func TestContextListCmd(t *testing.T) {
	resetVariables()
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/fixtures/.meshery/TestContext.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:             "list all contexts",
			Args:             []string{"context", "list"},
			ExpectedResponse: "listExpected.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			defer buf.Reset()

			SystemCmd.SetOut(buf)
			SystemCmd.SetErr(buf)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := buf.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata/context")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			//t.Logf("Expected response:\n%s", expectedResponse)
			//t.Logf("Actual response:\n%s", actualResponse)

			assert.Equal(t, expectedResponse, actualResponse)
		})
		t.Log("ListContextCmd test passed")
	}
}

func TestContextDeleteCmd(t *testing.T) {
	resetVariables()
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	configPath := filepath.Join(currDir, "testdata/context/ExpectedDelete.yaml")
	fixturePath := filepath.Join(currDir, "fixtures/.meshery/TestContext.yaml")

	t.Cleanup(func() {
		_ = utils.Populate(fixturePath, configPath)
	})

	tests := []struct {
		name               string
		args               []string
		setupConfig        func(t *testing.T, path string)
		expectError        bool
		expectedError      error
		expectedResponse   string
		expectedSubstr     string
		expectedCurrentCtx string
		expectedDeletedCtx string
	}{
		{
			name:               "delete given context",
			args:               []string{"context", "delete", "local2"},
			expectedResponse:   "delete.context.golden",
			expectedCurrentCtx: "local",
			expectedDeletedCtx: "local2",
		},
		{
			name:               "delete current context with --set flag",
			args:               []string{"context", "delete", "local", "-y", "--set", "local2"},
			expectedSubstr:     "deleted context local",
			expectedCurrentCtx: "local2",
			expectedDeletedCtx: "local",
		},
		{
			name:               "delete current context with mixed-case --set flag",
			args:               []string{"context", "delete", "local", "-y", "--set", "Local2"},
			expectedSubstr:     "deleted context local",
			expectedCurrentCtx: "local2",
			expectedDeletedCtx: "local",
		},
		{
			name: "error when deleting the only context",
			args: []string{"context", "delete", "local", "-y"},
			setupConfig: func(t *testing.T, path string) {
				singleCtx := `contexts:
  local:
    endpoint: http://localhost:9081
    token: default
    platform: kubernetes
current-context: local
`
				if err := os.WriteFile(path, []byte(singleCtx), 0644); err != nil {
					t.Fatal(err)
				}
			},
			expectError:   true,
			expectedError: utils.ErrInvalidArgument(errors.New("cannot delete the only context; create another with `mesheryctl system context create` and try again")),
		},
	}

	mesheryctlflags.InitValidators(SystemCmd)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Always start with fresh fixture
			if err := utils.Populate(fixturePath, configPath); err != nil {
				t.Fatal(err)
			}
			if tt.setupConfig != nil {
				tt.setupConfig(t, configPath)
			}
			utils.SetupCustomContextEnv(t, configPath)

			// Reset flags to prevent bleed-over between tests
			contextDeleteFlags = cmdContextDeleteFlags{}
			utils.SilentFlag = false

			buf := utils.SetupMeshkitLoggerTesting(t, false)
			defer buf.Reset()

			SystemCmd.SetOut(buf)
			SystemCmd.SetErr(buf)
			SystemCmd.SetArgs(tt.args)
			err := SystemCmd.Execute()

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				utils.AssertMeshkitErrorsEqual(t, err, tt.expectedError)
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			actualResponse := buf.String()
			if tt.expectedResponse != "" {
				testdataDir := filepath.Join(currDir, "testdata/context")
				golden := utils.NewGoldenFile(t, tt.expectedResponse, testdataDir)
				if *update {
					golden.Write(actualResponse)
				}
				expectedResponse := golden.Load()
				assert.Equal(t, expectedResponse, actualResponse)
			}

			if tt.expectedSubstr != "" {
				assert.Contains(t, actualResponse, tt.expectedSubstr)
			}

			content, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatal(err)
			}

			var data map[string]interface{}
			if err := yaml.Unmarshal(content, &data); err != nil {
				t.Fatal(err)
			}

			if tt.expectedCurrentCtx != "" {
				assert.Equal(t, tt.expectedCurrentCtx, data["current-context"])
			}

			if tt.expectedDeletedCtx != "" {
				if contexts, ok := data["contexts"].(map[string]interface{}); ok {
					_, exists := contexts[tt.expectedDeletedCtx]
					assert.False(t, exists, "deleted context should not exist in contexts")
				}
			}

			// Repopulating Expected yaml for cleanup
			if err := utils.Populate(fixturePath, configPath); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestContextCreateCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/context/ExpectedAdd.yaml")
	tests := []struct {
		utils.CmdTestInput
		Validation func(*testing.T)
	}{
		{
			CmdTestInput: utils.CmdTestInput{
				Name:                 "given context name provided when system context create then context is created",
				Args:                 []string{"context", "create", "local3"},
				ExpectedResponse:     "createContext.golden",
				ExpectedResponseYaml: "addExpected.golden",
			},
		},
		{
			CmdTestInput: utils.CmdTestInput{
				Name:           "given no context provided when system context create  then thorw error",
				Args:           []string{"context", "create"},
				ExpectError:    true,
				ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errArgMsg, contextCreateUsageMsg)),
				IsOutputGolden: false,
			},
		},
		{
			CmdTestInput: utils.CmdTestInput{
				Name:           "given multiple context name provided when system context create then throw error",
				Args:           []string{"context", "create", "local1", "local2"},
				ExpectError:    true,
				ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errArgMsg, contextCreateUsageMsg)),
				IsOutputGolden: false,
			},
		},
		{
			CmdTestInput: utils.CmdTestInput{
				Name:                 "Create a context with a valid component succeeds",
				Args:                 []string{"context", "create", "local-valid", "--components", "meshery-istio"},
				ExpectedResponse:     "createContextValid.golden",
				ExpectedResponseYaml: "addExpectedValid.golden",
			},
		},
		{
			CmdTestInput: utils.CmdTestInput{
				Name:           "Create a context with a non-existent component fails",
				Args:           []string{"context", "create", "local4", "--components", "invalid-component"},
				ExpectError:    true,
				ExpectedError:  ErrInvalidComponent("invalid-component"),
				IsOutputGolden: false,
			},
			Validation: func(t *testing.T) {
				mctl, err := config.GetMesheryCtl(viper.GetViper())
				if err != nil {
					t.Fatalf("Failed to load configuration: %v", err)
				}
				if _, ok := mctl.Contexts["local4"]; ok {
					t.Errorf("Expected context 'local4' to not be saved, but it was found")
				}
			},
		},
	}

	mesheryctlflags.InitValidators(SystemCmd)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Reset flags to prevent bleed-over between tests
			contextCreateFlags = cmdContextCreateFlags{}
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			defer buf.Reset()

			SystemCmd.SetOut(buf)
			SystemCmd.SetErr(buf)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()

			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					if tt.Validation != nil {
						tt.Validation(t)
					}
					return
				}
				t.Fatal(err)
			}

			if tt.ExpectError {
				t.Fatalf("expected error, got nil")
			}

			actualResponse := buf.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata/context")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			assert.Equal(t, expectedResponse, actualResponse)
			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/context/ExpectedAdd.yaml"

			content, err := os.ReadFile(filepath)
			if err != nil {
				t.Error(err)
			}
			actualResponse = string(content)
			golden = utils.NewGoldenFile(t, tt.ExpectedResponseYaml, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			addExpected := golden.Load()
			assert.Equal(t, addExpected, actualResponse)

			//Repopulating Expected yaml
			if err := utils.Populate(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}

			if tt.Validation != nil {
				tt.Validation(t)
			}
		})
		t.Log("CreateContextCmd test Passed")
	}

}

func TestContextAddUppercaseCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/context/ExpectedAdd.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:                 "given context name which contains uppercase provided when system context create then context is created in lowercase",
			Args:                 []string{"context", "create", "Local3"},
			ExpectedResponse:     "createContext.golden",
			ExpectedResponseYaml: "addExpected.golden",
		},
	}

	mesheryctlflags.InitValidators(SystemCmd)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Reset flags to prevent bleed-over between tests
			contextCreateFlags = cmdContextCreateFlags{}
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			defer buf.Reset()

			SystemCmd.SetOut(buf)
			SystemCmd.SetErr(buf)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()

			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					return
				}
				t.Fatal(err)
			}

			if tt.ExpectError {
				t.Fatalf("expected error, got nil")
			}

			actualResponse := buf.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata/context")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			assert.Equal(t, expectedResponse, actualResponse)
			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/context/ExpectedAdd.yaml"

			content, err := os.ReadFile(filepath)
			if err != nil {
				t.Error(err)
			}
			actualResponse = string(content)
			golden = utils.NewGoldenFile(t, tt.ExpectedResponseYaml, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			addExpected := golden.Load()
			assert.Equal(t, addExpected, actualResponse)

			// To check context is lowercase
			contentFile, err := os.ReadFile(filepath)
			if err != nil {
				t.Fatal(err)
			}

			var data map[string]interface{}
			err = yaml.Unmarshal(contentFile, &data)
			if err != nil {
				t.Fatal(err)
			}

			contexts := data["contexts"].(map[string]interface{})

			if len(tt.Args) > 2 {
				_, existsUpper := contexts[tt.Args[2]]
				if existsUpper {
					t.Fatalf("uppercase context should not exist")
				}
			}

			//Repopulating Expected yaml
			if err := utils.Populate(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
		t.Log("CreateUppercaseContextCmd test Passed")
	}

}

func TestContextSwitchCmd(t *testing.T) {
	resetVariables()
	// get current directory
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/context/ExpectedSwitch.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:             "switch to a different context",
			Args:             []string{"context", "switch", "local2"},
			ExpectedResponse: "switch.context.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			defer buf.Reset()

			SystemCmd.SetOut(buf)
			SystemCmd.SetErr(buf)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := buf.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata/context")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			assert.Equal(t, expectedResponse, actualResponse)

			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/context/ExpectedSwitch.yaml"
			content, err := os.ReadFile(filepath)
			if err != nil {
				t.Error(err)
			}
			actualResponse = string(content)
			golden = utils.NewGoldenFile(t, "switchExpected.golden", testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			switchExpected := golden.Load()
			assert.Equal(t, switchExpected, actualResponse)

			if err := utils.Populate(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
		t.Log("SwitchContextCmd test passed")
	}
}

func resetVariables() {
	//reset context before tests
	tempCntxt = ""
}
