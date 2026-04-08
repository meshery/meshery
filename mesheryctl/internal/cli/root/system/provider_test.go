package system

import (
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewProviderCmd(t *testing.T) {
	setupContextTestEnv(t)
	tests := []CmdTestInput{
		{
			Name:             "given provider context when provider view then display provided provider",
			Args:             []string{"provider", "view", "-c", "gke"},
			ExpectedResponse: PrintProviderToStdout(mctlCfg.Contexts["gke"], "gke") + "\n\n",
		},
		{
			Name:             "given no provider context when provider view then display default provider",
			Args:             []string{"provider", "view"},
			ExpectedResponse: PrintProviderToStdout(mctlCfg.Contexts["local"], "local") + "\n\n",
		},
		{
			Name: "given --all flag provided when provider view then display all providers",
			Args: []string{"provider", "view", "--all"},
			ExpectedResponse: func() string {
				output := strings.Builder{}
				keys := make([]string, 0, len(mctlCfg.Contexts))
				for k := range mctlCfg.Contexts {
					keys = append(keys, k)
				}

				// 2. Sort the slice alphabetically
				slices.Sort(keys)

				for _, contextName := range keys {
					context := mctlCfg.Contexts[contextName]
					output.WriteString(PrintProviderToStdout(context, contextName) + "\n\n")
				}

				output.WriteString("Current Context: local\n")
				return output.String()
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			buf := setupSystemOutCmdTest(t)
			defer buf.Reset()
			SystemCmd.SetArgs(tt.Args)
			err = SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			actualResponse := buf.String()
			expectedResponse := tt.ExpectedResponse

			assert.Equal(t, expectedResponse, actualResponse)

			BreakupFunc()
		})
	}
}
func TestRunProviderWithNoCmdOrFlag(t *testing.T) {
	setupSystemOutCmdTest(t)
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
