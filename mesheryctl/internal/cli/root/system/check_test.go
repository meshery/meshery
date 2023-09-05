package system

import (
	"bytes"
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update golden files")

// This is an Integration test
func TestPreflightCmdIntegration(t *testing.T) {
	// skipping this integration test with --short flag
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	SetupContextEnv(t)
	tests := []utils.CmdTestInput{
		{
			Name:             "Run preflight check",
			Args:             []string{"check", "--pre"},
			ExpectedResponse: "check.output.golden",
		},
		{
			Name:             "Run preflight check",
			Args:             []string{"check", "--preflight"},
			ExpectedResponse: "check.output.golden",
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
			// testdataDir := filepath.Join(filepath.Dir(filename), tf.dir, tf.name)
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, filepath.Join(testdataDir, "check"))

			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			assert.Equal(t, expectedResponse, actualResponse)
		})
		t.Log("PreflightCmdIntegration Test Passed")
	}
}
