package system

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"
)

//var update1 = flag.Bool("update", true, "update golden files")

// This is a Unit test
func TestResetCmd(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/reset/TestResetContexts.yaml")

	tests := []utils.CmdTestInput{
		{
			Name:             "Reset the meshery config file with docker platform",
			Args:             []string{"reset", "-y", "-c", "local2"},
			ExpectedResponse: "reset_docker.output.golden",
		},
		{
			Name:             "Reset the meshery config file with kubernetes platform",
			Args:             []string{"reset", "-y", "-c", "local"},
			ExpectedResponse: "reset_kubernetes.output.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// setting up log to grab logs
			b := utils.SetupLogrusGrabTesting(t)
			SystemCmd.SetOut(b)
			SystemCmd.SetArgs(tt.Args)
			err = SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()

			testdataDir := filepath.Join(currDir, "testdata/reset/")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			//if *update1 {
			//	golden.Write(actualResponse)
			//}
			expectedResponse := golden.Load()

			//Comparing the expected response with the actual response
			if expectedResponse != actualResponse {
				//utils.Equals(t, expectedResponse, actualResponse)
				t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
			}

		})
	}
	RemoveFileAndDirectory()
}

func RemoveFileAndDirectory() {
	// Removing the meshery.yaml file
	// from the directory
	removeMesheryYaml := os.Remove("meshery.yaml")
	if removeMesheryYaml != nil {
		log.Fatal(removeMesheryYaml)
	}

	// Removing the .meshery directory
	removeMesheryHiddenDirectory := os.RemoveAll(".meshery")
	if removeMesheryHiddenDirectory != nil {
		log.Fatal(removeMesheryHiddenDirectory)
	}
}
