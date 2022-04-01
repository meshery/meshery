package system

import (
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestStop(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping unit test")
	}
	utils.SetupContextEnv(t)
	utils.StartMockery(t)
	tests := []utils.CmdTestInput{
		{
			Name:             "Run stop",
			Args:             []string{"stop"},
			ExpectError:      true,
		},
		{
			Name:			 "Run stop with --force",
			Args:			 []string{"stop", "--force"},
			ExpectError:	 true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.SetupLogrusFormatter()

			SystemCmd.SetArgs(tt.Args)
			err = SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}
		})
	}
	utils.StopMockery(t)
}