package system

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestViewContextCmd(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b := utils.SetupMeshkitLoggerTesting(t, false)
			SystemCmd.SetOut(b)
			SystemCmd.SetErr(b)
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
			actualResponse := b.String()

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
func TestListContextCmd(t *testing.T) {
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
			b := utils.SetupMeshkitLoggerTesting(t, false)
			SystemCmd.SetOut(b)
			SystemCmd.SetErr(b)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
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

func TestDeleteContextCmd(t *testing.T) {
	resetVariables()
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/context/ExpectedDelete.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:             "delete given context",
			Args:             []string{"context", "delete", "local2"},
			ExpectedResponse: "delete.context.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b := utils.SetupMeshkitLoggerTesting(t, false)
			SystemCmd.SetOut(b)
			SystemCmd.SetErr(b)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata/context")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			if expectedResponse != actualResponse {
				t.Error("Expected response not obtained")
				t.Errorf("Expected: %v", expectedResponse)
				t.Errorf("Actual: %v", actualResponse)
			}
			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/context/ExpectedDelete.yaml"

			content, err := os.ReadFile(filepath)
			if err != nil {
				t.Error(err)
			}
			actualResponse = string(content)
			golden = utils.NewGoldenFile(t, "deleteExpected.golden", testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			deleteExpected := golden.Load()
			if actualResponse != deleteExpected {
				t.Error("Contexts are mismatched")
				t.Error("Expected:")
				t.Errorf("%v", deleteExpected)
				t.Error("Actual:")
				t.Errorf("%v", actualResponse)
			}

			//Repopulating Expected yaml
			if err := utils.Populate(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
		t.Log("DeleteContextCmd test Passed")
	}
}
func TestAddContextCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/context/ExpectedAdd.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:                 "given context name provided when system context create [valid-name] then context is created",
			Args:                 []string{"context", "create", "local3"},
			ExpectedResponse:     "createContext.golden",
			ExpectedResponseYaml: "addExpected.golden",
		},
		{
			Name:           "given no context provided when system context create  then thorw error",
			Args:           []string{"context", "create"},
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errArgMsg, contextCreateUsageMsg)),
			IsOutputGolden: false,
		},
		{
			Name:           "given multiple context name provided when system context create then throw error",
			Args:           []string{"context", "create", "local1", "local2"},
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errArgMsg, contextCreateUsageMsg)),
			IsOutputGolden: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b := utils.SetupMeshkitLoggerTesting(t, false)
			utils.Log.SetLevel(logrus.InfoLevel)

			SystemCmd.SetOut(b)
			SystemCmd.SetErr(b)
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

			actualResponse := b.String()
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
		})
		t.Log("CreateContextCmd test Passed")
	}

}

func TestAddUppercaseContextCmd(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b := utils.SetupMeshkitLoggerTesting(t, false)
			utils.Log.SetLevel(logrus.InfoLevel)

			SystemCmd.SetOut(b)
			SystemCmd.SetErr(b)
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

			actualResponse := b.String()
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

func TestSwitchContextCmd(t *testing.T) {
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
			b := utils.SetupMeshkitLoggerTesting(t, false)

			SystemCmd.SetOut(b)
			SystemCmd.SetErr(b)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
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
	newContext = ""
	currContext = ""
	tempCntxt = ""
}
