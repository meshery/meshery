package context

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

// var b *bytes.Buffer
var update = flag.Bool("update", false, "update golden files")

func TestViewContextCmd(t *testing.T) {
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
			Args:             []string{"view"},
			ExpectedResponse: "view.golden",
		},
		{
			Name:             "view with specified context through context flag",
			Args:             []string{"view", "--context", "local2"},
			ExpectedResponse: "viewWithContextExpected.golden",
		},
		{
			Name:             "Error for viewing a non-existing context",
			Args:             []string{"view", "local3"},
			ExpectedResponse: "view.notexist.golden",
		},
		{
			Name:             "view with specified context as argument",
			Args:             []string{"view", "local2"},
			ExpectedResponse: "viewWithContextExpected.golden",
		},
		{
			Name:             "view with all flag set",
			Args:             []string{"view", "--all"},
			ExpectedResponse: "viewAllExpected.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b := utils.SetupLogrusGrabTesting(t)
			ContextCmd.SetOut(b)
			ContextCmd.SetArgs(tt.Args)
			err := ContextCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			if expectedResponse != actualResponse {
				t.Errorf("expected response [%v] and actual response [%v] don't match", expectedResponse, actualResponse)
			}
		})
	}
}
func TestListContextCmd(t *testing.T) {
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
			Args:             []string{"list"},
			ExpectedResponse: "listExpected.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b := utils.SetupLogrusGrabTesting(t)
			ContextCmd.SetOut(b)
			ContextCmd.SetArgs(tt.Args)
			err := ContextCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			if expectedResponse != actualResponse {
				t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
			}
		})
	}
}

func TestDeleteContextCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/ExpectedDelete.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:             "delete given context",
			Args:             []string{"delete", "local2"},
			ExpectedResponse: "delete.context.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b := utils.SetupLogrusGrabTesting(t)
			ContextCmd.SetOut(b)
			ContextCmd.SetArgs(tt.Args)
			err := ContextCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			if expectedResponse != actualResponse {
				t.Errorf("expected response [%v] and actual response [%v] don't match", expectedResponse, actualResponse)
			}
			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/ExpectedDelete.yaml"

			content, err := ioutil.ReadFile(filepath)
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
				t.Errorf("expected response [%v] and actual response [%v] don't match", deleteExpected, actualResponse)
			}

			//Repopulating Expected yaml
			if err := utils.Copy(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
	}
}
func TestAddContextCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/ExpectedAdd.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:             "add given context",
			Args:             []string{"create", "local3"},
			ExpectedResponse: "createContext.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b := utils.SetupLogrusGrabTesting(t)
			ContextCmd.SetOut(b)
			ContextCmd.SetArgs(tt.Args)
			err := ContextCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			if expectedResponse != actualResponse {
				t.Errorf("expected response [%v] and actual response [%v] don't match", expectedResponse, actualResponse)
			}
			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/ExpectedAdd.yaml"

			content, err := ioutil.ReadFile(filepath)
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
				t.Errorf("expected response [%v] and actual response [%v] don't match", addExpected, actualResponse)
			}

			//Repopulating Expected yaml
			if err := utils.Copy(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
	}
}

func TestSwitchContextCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/ExpectedSwitch.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:             "switch to a different context",
			Args:             []string{"switch", "local2"},
			ExpectedResponse: "switch.context.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b := utils.SetupLogrusGrabTesting(t)
			ContextCmd.SetOut(b)
			ContextCmd.SetArgs(tt.Args)
			err := ContextCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
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
			filepath := path + "/testdata/ExpectedSwitch.yaml"
			content, err := ioutil.ReadFile(filepath)
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
			if err := utils.Copy(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
	}
}
