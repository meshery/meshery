package system

import (
	"bytes"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestStartCmd(t *testing.T) {
	SetupContextEnv(t)
	tests := []utils.CmdTestInput{
		{
			Name:             "start without configurating",
			Args:             []string{"start"},
			ExpectedResponse: "start.output.golden",
		},
		{
			Name:             "start with configurating",
			Args:             []string{"start", "--config-cluster", "aks"},
			ExpectedResponse: "start.output.golden",
		},
		{
			Name:             "start with configurating",
			Args:             []string{"start", "--config-cluster", "eks"},
			ExpectedResponse: "start.output.golden",
		},
		{
			Name:             "start with configurating",
			Args:             []string{"start", "-p", "docker", "--config-cluster", "gke"},
			ExpectedResponse: "start.output.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var buf bytes.Buffer
			startCmd.SetArgs(tt.Args)
			err := startCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			output := buf.String()
			actualResponse := output

			_, filename, _, ok := runtime.Caller(0)
			if !ok {
				t.Fatal("Not able to get current working directory")
			}

			currDir := filepath.Dir(filename)
			testDataDir := filepath.Join(currDir, "testdata")

			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, filepath.Join(testDataDir, "start"))

			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			assert.Equal(t, expectedResponse, actualResponse)
		})
		t.Log("start test passed")
	}
}
