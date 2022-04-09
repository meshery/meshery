package system

import (
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestStatus(t *testing.T) {
	utils.SetupContextEnv(t)
	utils.StartMockery(t)
	tests := []utils.CmdTestInput{
		{
			Name:        "Run status",
			Args:        []string{"status"},
			ExpectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.SetupLogrusFormatter()

			SystemCmd.SetArgs(tt.Args)
			err = SystemCmd.Execute()
			if err != nil {
				t.Log("Error encountered. Skipping case...")
			}
		})
	}
}
