package system

import (
	"testing"
)

func TestViewProviderCmd(t *testing.T) {
	SetupContextEnv(t)
	tests := []CmdTestInput{
		{
			Name:             "view without any parameter",
			Args:             []string{"provider", "view"},
			ExpectedResponse: PrintProviderToStdout(mctlCfg.Contexts["local"], "local") + "\n\n",
		},
		{
			Name:             "view with context override",
			Args:             []string{"provider", "view", "-c", "gke"},
			ExpectedResponse: PrintProviderToStdout(mctlCfg.Contexts["gke"], "gke") + "\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			SetupFunc()
			SystemCmd.SetArgs(tt.Args)
			err = SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := b.String()
			expectedResponse := tt.ExpectedResponse

			if expectedResponse != actualResponse {
				t.Errorf("\nExpected response\n\n%v\n\nActual response\n\n%v\n\nMismatch between expected and actual response\n", expectedResponse, actualResponse)
			}
			BreakupFunc()
		})
	}
}
func TestRunProviderWithNoCmdOrFlag(t *testing.T) {
	SetupFunc()
	SystemCmd.SetArgs([]string{"provider"})
	err = SystemCmd.Execute()

	actualResponse := err.Error()

	expectedResponse := "please specify a flag or subcommand. Use 'mesheryctl system provider --help' to display user guide." + "\n\n"
	expectedResponse += "See https://docs.meshery.io/reference/mesheryctl/system/provider for usage details\n"

	if expectedResponse != actualResponse {
		t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
	}
	BreakupFunc()
}
