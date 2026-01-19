package system

import (
	"bytes"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

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

			actualResponse := buf.String()
			assert.Contains(t, actualResponse, "Docker")
			assert.Contains(t, actualResponse, "Kubernetes API")
			assert.Contains(t, actualResponse, "Kubernetes Version")
			assert.Contains(t, actualResponse, "Meshery prerequisites")
		})
		t.Log("PreflightCmdIntegration Test Passed")
	}
}
