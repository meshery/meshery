package system

import (
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestLogout(t *testing.T) {
	tests := []utils.CmdTestInput{
		{
			Name:             "Run logout",
			Args:             []string{"logout"},
			ExpectedResponse: "logout.output.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.SetupLogrusFormatter()

			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
