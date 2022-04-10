package system

import (
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestStop(t *testing.T) {
	tests := []utils.CmdTestInput{
		{
			Name:        "Run stop",
			Args:        []string{"stop"},
			ExpectError: true,
		},
		{
			Name:        "Run stop with --force",
			Args:        []string{"stop", "--force"},
			ExpectError: false,
		},
	}

	for _, tt := range tests {
		utils.SetupContextEnv(t)
		t.Run(tt.Name, func(t *testing.T) {
			utils.SetupLogrusFormatter()
			cmd := SystemCmd
			cmd.SetArgs(tt.Args)
			err = cmd.Execute()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
