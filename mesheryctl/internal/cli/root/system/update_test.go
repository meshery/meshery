package system

import (
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestUpdate(t *testing.T) {
	utils.SetupContextEnv(t)
	tests := []utils.CmdTestInput{
		{
			Name:        "Run update",
			Args:        []string{"update", "-y"},
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
