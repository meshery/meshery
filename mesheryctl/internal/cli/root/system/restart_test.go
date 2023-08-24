package system

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestRestartCmd(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Can't get working Directory")
	}
	currDir := filepath.Dir(filename)
	restartDir := "/testdata/restart/TestRestartContext.yaml"
	utils.SetupCustomContextEnv(t, currDir+restartDir)
	testcases := []utils.CmdTestInput{
		{
			Name:             "Restart meshery on docker",
			Args:             []string{"restart","-y", "-c", "docker"},
			ExpectedResponse: "restart_k8.output.golden",
		},
		{
			Name:             "Restart Meshery on k8",
			Args:             []string{"restart","-y", "-c", "k8s"},
			ExpectedResponse: "restart_docker.output.golden",
		},
	}
	for _, test := range testcases {
		t.Run(test.Name, func(t *testing.T) {
			buff := utils.SetupLogrusGrabTesting(t, false)
			cmd := SystemCmd
			cmd.SetOut(buff)
			cmd.SetArgs(test.Args)
			err = cmd.Execute()
			if err != nil {
				t.Error(err)
			}
			actualResponse := buff.String()
			testdataDir := filepath.Join(currDir, "testdata/restart/")
			golden := utils.NewGoldenFile(t, test.ExpectedResponse, testdataDir)

			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			assert.Equal(t, expectedResponse, actualResponse)

		})
	}
}
