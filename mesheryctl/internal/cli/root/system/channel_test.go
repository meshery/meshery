package system

import (
	"slices"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestViewCmd(t *testing.T) {
	setupContextTestEnv(t)
	tests := []CmdTestInput{
		{
			Name:             "given no argument when view then display current context channel and version",
			Args:             []string{"channel", "view"},
			ExpectedResponse: PrintChannelAndVersionToStdout(mctlCfg.Contexts["local"], "local") + "\n\n",
		},
		{
			Name:             "given context override when view then display specified context channel and version",
			Args:             []string{"channel", "view", "-c", "gke"},
			ExpectedResponse: PrintChannelAndVersionToStdout(mctlCfg.Contexts["gke"], "gke") + "\n\n",
		},
		{
			Name: "given --all flag when view then display all contexts channel and version",
			Args: []string{"channel", "view", "--all"},
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
					output.WriteString(PrintChannelAndVersionToStdout(context, contextName) + "\n\n")
				}

				output.WriteString("Current Context: local\n")
				return output.String()
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			buf := utils.SetupMeshkitLoggerTesting(t, false)
			defer buf.Reset()
			SystemCmd.SetOut(buf)
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

func TestSetCmd(t *testing.T) {
	setupContextTestEnv(t)
	tests := []CmdTestInput{
		{
			Name:             "Set Docker Platform  Channel",
			Args:             []string{"channel", "set", "stable", "-c", "local"},
			ExpectedResponse: "Channel set to stable-latest" + "\n",
		},
		{
			Name:             "Set GKE Platform Channel",
			Args:             []string{"channel", "set", "stable", "-c", "gke"},
			ExpectedResponse: "Channel set to stable-latest" + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			buf := setupSystemOutCmdTest(t)
			defer func() {
				buf.Reset()
				// resetCmdFlags(SystemCmd, t)
			}()
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
	t.Log("Set ChannelCmd Passed")
}

//UNDER REVIEW
// func TestSwitchCmd(t *testing.T) {
// 	_, filename, _, ok := runtime.Caller(0)

// 	if !ok {
// 		t.Fatal("Not able to get current working directory")
// 	}
// 	currDir := filepath.Dir(filename)
// 	SetupContextEnv(t)
// 	tests := []utils.CmdTestInput{
// 		{
// 			Name:             "Switch docker Channel",
// 			Args:             []string{"channel", "switch", "stable", "-y", "-c", "local"},
// 			ExpectedResponse: "switch.docker.output.golden",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.Name, func(t *testing.T) {
// 			buff := utils.SetupLogrusGrabTesting(t, false)
// 			SystemCmd.SetOut(buff)
// 			SystemCmd.SetArgs(tt.Args)
// 			err := SystemCmd.Execute()
// 			if err != nil {
// 				if errSubstrs := tt.ErrorStringContains; len(errSubstrs) > 0 && checkErrorContains(err, errSubstrs) {
// 					return
// 				}
// 				t.Error(err)
// 			}
// 			actualResponse := buff.String()
// 			testdataDir := filepath.Join(currDir, "testdata/channel/")
// 			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
// 			if *update {
// 				golden.Write(actualResponse)
// 			}
// 			expectedResponse := golden.Load()
// 			assert.Equal(t, expectedResponse, actualResponse)

// 			path, err := os.Getwd()
// 			if err != nil {
// 				t.Error("unable to locate meshery directory")
// 			}
// 			filepath := path + "/testdata/channel/switch.output.golden"
// 			content, err := os.ReadFile(filepath)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			actualResponse = string(content)
// 			golden = utils.NewGoldenFile(t, "switchExpected.golden", testdataDir)
// 			if *update {
// 				golden.Write(actualResponse)
// 			}
// 			switchExpected := golden.Load()
// 			assert.Equal(t, switchExpected, actualResponse)
// 			BreakupFunc()
// 		})
// 	}
// 	t.Log("Switch ChannelCmd Passed")
// }
