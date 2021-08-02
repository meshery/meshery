package system

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"
)

//var update1 = flag.Bool("update", true, "update golden files")

// This is an Integration test
func TestResetCmd(t *testing.T) {
	SetupContextEnv(t)
	tests := []utils.CmdTestInput{
		{
			Name:             "Reset the meshery config file with docker platform",
			Args:             []string{"reset", "-y"},
			ExpectedResponse: "reset.output.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// setting up log to grab logs
			var buf bytes.Buffer
			log.SetOutput(&buf)
			utils.SetupLogrusFormatter()

			SystemCmd.SetArgs(tt.Args)
			err = SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			output := buf.String()
			actualResponse := output

			// get current directory
			_, filename, _, ok := runtime.Caller(0)
			if !ok {
				t.Fatal("Not able to get current working directory")
			}

			currDir := filepath.Dir(filename)
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, filepath.Join(testdataDir, "reset"))

			//if *update1 {
			//	golden.Write(actualResponse)
			//}
			expectedResponse := golden.Load()

			//Comparing the expected response with the actual response
			if expectedResponse != actualResponse {
				utils.Equals(t, expectedResponse, actualResponse)
				// t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
			}
			RemoveFileAndDirectory()
		})
	}
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
