package system

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// This is a Unit test
func TestResetCmd(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/reset/TestResetContexts.yaml")

	tests := []utils.CmdTestInput{
		//Test for platform docker
		{
			Name:             "Reset the meshery config file with docker platform",
			Args:             []string{"reset", "-y", "-c", "local2"},
			ExpectedResponse: "reset_docker.output.golden",
		},
		// This test is failing with error : "failed to make GET request: could not find path: install/deployment_yamls/k8s in the manifest tree"
		//Test for platform kubernetes
		// {
		// 	Name:             "Reset the meshery config file with kubernetes platform",
		// 	Args:             []string{"reset", "-y", "-c", "local"},
		// 	ExpectedResponse: "reset_kubernetes.output.golden",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// setting up log to grab logs
			b := utils.SetupLogrusGrabTesting(t, false)
			SystemCmd.SetOut(b)
			SystemCmd.SetArgs(tt.Args)
			err = SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}
			//Collecting the actual response
			actualResponse := b.String()

			testdataDir := filepath.Join(currDir, "testdata/reset/")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			if *update {
				golden.Write(actualResponse)
			}

			//Collecting the expected response
			expectedResponse := golden.Load()

			//Comparing the expected response with the actual response
			if expectedResponse != actualResponse {
				utils.Equals(t, expectedResponse, actualResponse)
			}
		})
	}
	//Removing the files and directory that were created during the test
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
