package system

import (
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestStart(t *testing.T) {
	utils.SetupContextEnv(t)
	tests := []utils.CmdTestInput{
		{
			Name:        "Run start with --reset",
			Args:        []string{"start", "--reset", "-y"},
			ExpectError: true,
		},
		{
			Name:        "Run start with --skip-update",
			Args:        []string{"start", "--skip-update"},
			ExpectError: true,
		},
		{
			Name:        "Run start with --skip-browser",
			Args:        []string{"start", "--skip-browser"},
			ExpectError: true,
		},
	}

	for _, tt := range tests {
		utils.SetupContextEnv(t)
		t.Run(tt.Name, func(t *testing.T) {
			utils.SetupLogrusFormatter()
			cmd := SystemCmd
			cmd.SetArgs(tt.Args)
			err = cmd.Execute()
			_, errPo := utils.AreAllPodsRunning()
			if errPo != nil {
				t.Log("Pods are not running!")
			}
			if err != nil {
				t.Log("Error obtained: ", err)
			}
		})
	}
}
