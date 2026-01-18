package system

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
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
				if tt.Name == "Error for viewing a non-existing context" {
					assert.Contains(t, err.Error(), "doesn't exists")
					return
				}
				t.Error(err)
			}

			actualResponse := b.String()
			switch tt.Name {
			case "view for default context":
				assert.Contains(t, actualResponse, "Current Context: local")
			case "view with specified context through context flag", "view with specified context as argument":
				assert.Contains(t, actualResponse, "Current Context: local2")
			case "view with all flag set":
				assert.Contains(t, actualResponse, "local:")
				assert.Contains(t, actualResponse, "local2:")
			}
			t.Log("ViewContextCmd test passed")
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
			assert.Contains(t, actualResponse, "Current context: local")
			assert.Contains(t, actualResponse, "- local")
			assert.Contains(t, actualResponse, "- local2")
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
			b := utils.SetupLogrusGrabTesting(t, false)
			SystemCmd.SetOut(b)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			assert.Contains(t, actualResponse, "deleted context local2")
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
			assert.Contains(t, actualResponse, "current-context: local")
			assert.NotContains(t, actualResponse, "local2:")

			//Repopulating Expected yaml
			if err := utils.Populate(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
		t.Log("DeleteContextCmd test Passed")
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
			utils.SetupMeshkitLoggerTesting(t, false)
			b := utils.SetupLogrusGrabTesting(t, false)
			SystemCmd.SetOut(b)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			assert.Contains(t, actualResponse, "Added `local3` context")
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
			assert.Contains(t, actualResponse, "local3:")
			assert.Contains(t, actualResponse, "provider: Layer5")

			//Repopulating Expected yaml
			if err := utils.Populate(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
		t.Log("CreateContextCmd test Passed")
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
			b := utils.SetupLogrusGrabTesting(t, false)
			SystemCmd.SetOut(b)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			assert.Contains(t, actualResponse, "switched to context 'local2'")

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
			assert.Contains(t, actualResponse, "current-context: local2")

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
