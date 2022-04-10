package system

import (
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/micmonay/keybd_event"
)

func TestLogin(t *testing.T) {
	utils.SetupContextEnv(t)
	utils.StartMockery(t)
	tests := []utils.CmdTestInput{
		{
			Name:             "Run login",
			Args:             []string{"login"},
			ExpectedResponse: "login.output.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.SetupLogrusFormatter()

			kb, err := keybd_event.NewKeyBonding()
			if err != nil {
				t.Error(err)
			}
			cmd := SystemCmd
			cmd.SetArgs(tt.Args)
			err = cmd.Execute()
			kb.SetKeys(keybd_event.VK_ENTER) // To press and simulate Enter key
			if err != nil {
				t.Error(err)
			}
		})
	}
}
