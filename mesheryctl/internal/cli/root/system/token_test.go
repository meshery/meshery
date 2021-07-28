package system

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestTokenCreateCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/fixtures/.meshery/meshery.yaml")
	tests := []utils.CmdTestInput{
		{
			Name:             "create the passed token with default location",
			Args:             []string{"token", "create", "Default2"},
			ExpectedResponse: "create.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			b := utils.SetupLogrusGrabTesting(t)
			SystemCmd.SetOut(b)
			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}
			//Check the stdout against the golden file
			actualResponse := b.String()
			// Expected response
			testdatatokenDir := filepath.Join(currDir, "testdata/token")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdatatokenDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			if expectedResponse != actualResponse {
				t.Errorf("expected response [%v] and actual response [%v] don't match", expectedResponse, actualResponse)
			}

			//Check the modified yaml against the golden file
			path, err := os.Getwd()
			if err != nil {
				t.Error("unable to locate meshery directory")
			}
			filepath := path + "/testdata/token/create.yaml"
			content, err := ioutil.ReadFile(filepath)
			if err != nil {
				t.Error(err)
			}
			actualResponse = string(content)
			golden = utils.NewGoldenFile(t, "createExpected.golden", testdatatokenDir)
			if *update {
				golden.Write(actualResponse)
			}
			createExpected := golden.Load()
			if actualResponse != createExpected {
				t.Errorf("expected response %v and actual response %v don't match", createExpected, actualResponse)
			}
			if err := copy(path+"/fixtures/.meshery/config.yaml", filepath); err != nil {
				t.Error(err, "Could not complete test. Unable to configure create test file")
			}
			BreakupFunc(t)
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
