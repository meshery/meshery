package system

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
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

// TestSetCmdWritesToTheActiveConfigNotTheSharedFixture pins the invariant that
// keeps `go test ./mesheryctl/...` deterministic: pkg/utils/TestConfig.yaml is
// read by every mesheryctl test package, those packages run as concurrent
// processes, and viper.WriteConfig truncates before it rewrites. A sibling
// package reading the fixture inside that truncate window parsed an empty
// document without error and exited its whole test binary from
// GetBaseMesheryURL's Log.Fatal, which surfaced as a bare package-level FAIL
// naming no test.
//
// `channel set edge` is used rather than `stable` precisely because it differs
// from the fixture, so a write that reaches the shared file changes its bytes
// instead of rewriting them identically.
func TestSetCmdWritesToTheActiveConfigNotTheSharedFixture(t *testing.T) {
	fixture := utils.SharedTestConfigPath(t)
	before, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatalf("unable to read shared meshconfig fixture: %v", err)
	}
	info, err := os.Stat(fixture)
	if err != nil {
		t.Fatalf("unable to stat shared meshconfig fixture: %v", err)
	}
	t.Cleanup(func() {
		// A regression here corrupts the fixture for every other package, so
		// put it back rather than letting one failure cascade.
		after, err := os.ReadFile(fixture)
		if err != nil || bytes.Equal(before, after) {
			return
		}

		// Put it back by rename rather than by writing it in place. Writing
		// truncates first, and a sibling package reading the fixture inside
		// that window is the very failure this test exists to prevent - the
		// repair must not reproduce it.
		restored, err := os.CreateTemp(filepath.Dir(fixture), ".TestConfig.yaml-restore-*")
		if err != nil {
			t.Errorf("unable to create the meshconfig fixture restore file: %v", err)
			return
		}
		defer func() { _ = os.Remove(restored.Name()) }()

		if _, err := restored.Write(before); err != nil {
			_ = restored.Close()
			t.Errorf("unable to write the meshconfig fixture restore file: %v", err)
			return
		}

		if err := restored.Close(); err != nil {
			t.Errorf("unable to close the meshconfig fixture restore file: %v", err)
			return
		}

		if err := os.Chmod(restored.Name(), info.Mode().Perm()); err != nil {
			t.Errorf("unable to set the mode of the meshconfig fixture restore file: %v", err)
			return
		}

		if err := os.Rename(restored.Name(), fixture); err != nil {
			t.Errorf("unable to restore shared meshconfig fixture: %v", err)
		}
	})

	setupContextTestEnv(t)
	t.Cleanup(BreakupFunc)
	active := viper.ConfigFileUsed()

	buf := setupSystemOutCmdTest(t)
	SystemCmd.SetArgs([]string{"channel", "set", "edge", "-c", "local"})
	if err := SystemCmd.Execute(); err != nil {
		t.Fatalf("channel set: %v", err)
	}
	assert.Equal(t, "Channel set to edge-latest\n", buf.String())

	after, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatalf("unable to re-read shared meshconfig fixture: %v", err)
	}
	assert.Equal(t, string(before), string(after),
		"`channel set` wrote through to the shared meshconfig fixture %s; every mesheryctl test package reads that file concurrently", fixture)

	// The write still has to land somewhere, or this passes for the wrong reason.
	assert.NotEqual(t, fixture, active)
	written := viper.New()
	written.SetConfigFile(active)
	if err := written.ReadInConfig(); err != nil {
		t.Fatalf("unable to read the active meshconfig %v: %v", active, err)
	}
	persisted, err := config.GetMesheryCtl(written)
	if err != nil {
		t.Fatalf("unable to decode the active meshconfig: %v", err)
	}
	assert.Equal(t, "edge", persisted.Contexts["local"].Channel)
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
