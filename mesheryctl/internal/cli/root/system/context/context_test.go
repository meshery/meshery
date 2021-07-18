package context

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var b *bytes.Buffer
var update = flag.Bool("update", false, "update golden files")

type CmdTestInput struct {
	Name             string
	Args             []string
	ExpectedResponse string
}

func SetupContextEnv(t *testing.T, testfile string) {
	path, err := os.Getwd()
	if err != nil {
		t.Error("unable to locate meshery directory")
	}
	viper.Reset()
	viper.SetConfigFile(path + testfile)
	err = viper.ReadInConfig()
	if err != nil {
		t.Errorf("unable to read configuration from %v, %v", viper.ConfigFileUsed(), err.Error())
	}

	configuration, err = config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Error("error processing config", err)
	}
}
func SetupFunc(t *testing.T) {
	b = bytes.NewBufferString("")
	logrus.SetOutput(b)
	utils.SetupLogrusFormatter()
	ContextCmd.SetOut(b)
}
func TestViewContextCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	SetupContextEnv(t, "/fixtures/.meshery/TestContext.yaml")

	tests := []CmdTestInput{
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
			SetupFunc(t)
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
	SetupContextEnv(t, "/fixtures/.meshery/TestContext.yaml")
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	tests := []CmdTestInput{
		{
			Name:             "list all contexts",
			Args:             []string{"list"},
			ExpectedResponse: "listExpected.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			SetupFunc(t)
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
	SetupContextEnv(t, "/testdata/ExpectedDelete.yaml")
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	tests := []CmdTestInput{
		{
			Name:             "delete given context",
			Args:             []string{"delete", "local2"},
			ExpectedResponse: "delete.context.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			SetupFunc(t)
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
			if err := copy(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
	}
}
func TestAddContextCmd(t *testing.T) {
	SetupContextEnv(t, "/testdata/ExpectedAdd.yaml")
	// get current directory
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	tests := []CmdTestInput{
		{
			Name:             "add given context",
			Args:             []string{"create", "local3"},
			ExpectedResponse: "Added `local3` context\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			SetupFunc(t)
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
			if err := copy(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
	}
}

func TestSwitchContextCmd(t *testing.T) {
	SetupContextEnv(t, "/testdata/ExpectedSwitch.yaml")
	// get current directory
	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	tests := []CmdTestInput{
		{
			Name:             "switch to a different context",
			Args:             []string{"switch", "local2"},
			ExpectedResponse: "switch.context.golden",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			SetupFunc(t)
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
			if err := copy(path+"/fixtures/.meshery/TestContext.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure delete test file")
			}
		})
	}
}

//utility function to repopulate config
func copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
