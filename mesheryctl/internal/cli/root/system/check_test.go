package system

import (
	"path/filepath"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var baseTestDataPath = "mesheryctl/internal/cli/root/system/testdata"

func TestPreflightCmd(t *testing.T) {
	tests := []utils.CmdTestInput{
		{
			Name:             "Run preflight check",
			Args:             []string{"check", "--pre"},
			ExpectedResponse: "default-check-output.golden",
			MockResponse:     "default-check-input.golden",
		},
		{
			Name:             "Run preflight check",
			Args:             []string{"check", "--preflight"},
			ExpectedResponse: "default-check-output.golden",
			MockResponse:     "default-check-input.golden",
		},
		{
			Name:             "Docker not available",
			Args:             []string{"check", "--pre"},
			ExpectedResponse: "docker-not-check-output.golden",
			MockResponse:     "docker-not-check-input.golden",
		},
		{
			Name:             "Kubernetes not available",
			Args:             []string{"check", "--pre"},
			ExpectedResponse: "k8s-not-check-output.golden",
			MockResponse:     "k8s-not-check-input.golden",
		},
		{
			Name:             "No requirments met",
			Args:             []string{"check", "--pre"},
			ExpectedResponse: "no-req-met-check-output.golden",
			MockResponse:     "no-req-met-check-input.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			checkCmd.SetArgs(tt.Args)
			err = checkCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, filepath.Join(baseTestDataPath, "check"))

			expectedResponse := golden.Load()

			if expectedResponse != actualResponse {
				t.Fatalf("diff: %v", utils.Diff(expectedResponse, actualResponse))
				// t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
			}
		})
	}
}
