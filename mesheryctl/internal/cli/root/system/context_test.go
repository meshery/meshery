package system

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
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
			ExpectedResponse: "view.notexist.golden",
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
			b := utils.SetupLogrusGrabTesting(t, false)
			SystemCmd.SetOut(b)
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
				t.Errorf("expected response [%v] and actual response [%v] don't match", expectedResponse, actualResponse)
			} else {
				t.Log("test passed")
			}
		})
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
			b := utils.SetupLogrusGrabTesting(t, false)
			SystemCmd.SetOut(b)
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
		t.Log("test passed")
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
			b := utils.SetupLogrusGrabTesting(t, false)
			SystemCmd.SetOut(b)
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
	}
}
func TestAddContextCmd(t *testing.T) {
	resetVariables()
	// get current directory
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/context/ExpectedAdd.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:             "add given context",
			Args:             []string{"context", "create", "local3"},
			ExpectedResponse: "createContext.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b := utils.SetupLogrusGrabTesting(t, false)
			SystemCmd.SetOut(b)
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
				t.Errorf("Context: expected response [%v] and actual response [%v] don't match", expectedResponse, actualResponse)
			}
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
			golden = utils.NewGoldenFile(t, "addExpected.golden", testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			addExpected := golden.Load()
			if actualResponse != addExpected {
				t.Errorf("ExpectedAdd: expected response [%v] and actual response [%v] don't match", addExpected, actualResponse)
			}

			//Repopulating Expected yaml
			if err := utils.Populate(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
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
			Args:             []string{"context", "switch", "local2", "-y"},
			ExpectedResponse: "switch.context.golden",
		},
		{
			Name: "error switching to an empty context",
			Args: []string{"context", "switch"},
			ErrorStringContains: []string{
				"please provide exactly one context name",
				"Usage: mesheryctl system context switch [context name]",
				"Example: mesheryctl system context switch k8s-sample",
				"Description: Configures mesheryctl to actively use one one context vs. the another context",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b := utils.SetupLogrusGrabTesting(t, false)
			SystemCmd.SetOut(b)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				if errSubstrs := tt.ErrorStringContains; len(errSubstrs) > 0 && checkErrorContains(err, errSubstrs) {
					return
				}
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
				t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
			}
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
			if actualResponse != switchExpected {
				t.Errorf("expected response %v and actual response %v don't match", switchExpected, actualResponse)
			}
			if err := utils.Populate(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
	}
}

func resetVariables() {
	//reset context before tests
	newContext = ""
	currContext = ""
	tempCntxt = ""
}

func checkErrorContains(err error, substrs []string) bool {
	for _, substr := range substrs {
		if !strings.Contains(err.Error(), substr) {
			return false
		}
	}
	return true
}
