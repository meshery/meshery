package system

import (
	"bytes"
	"os"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var b *bytes.Buffer

func SetupContextEnv(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Error("unable to locate meshery directory")
	}
	viper.Reset()
	viper.SetConfigFile(path + "/../../../../pkg/utils/TestConfig.yaml")
	//fmt.Println(viper.ConfigFileUsed())
	err = viper.ReadInConfig()
	if err != nil {
		t.Errorf("unable to read configuration from %v, %v", viper.ConfigFileUsed(), err.Error())
	}

	mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Error("error processing config", err)
	}
}

func SetupFunc() {
	//fmt.Println(viper.AllKeys())
	b = bytes.NewBufferString("")
	logrus.SetOutput(b)
	utils.SetupLogrusFormatter()
	SystemCmd.SetOut(b)
}

func BreakupFunc() {
	viewCmd.Flags().VisitAll(setFlagValueAsUndefined)
}

func setFlagValueAsUndefined(flag *pflag.Flag) {
	_ = flag.Value.Set("")
}

type CmdTestInput struct {
	Name             string
	Args             []string
	ExpectedResponse string
	Token string
}

func TestViewCmd(t *testing.T) {
	SetupContextEnv(t)
	tests := []CmdTestInput{
		{
			Name:             "view without any parameter",
			Args:             []string{"channel", "view"},
			ExpectedResponse: PrintChannelAndVersionToStdout(mctlCfg.Contexts["local"], "local") + "\n\n",
		},
		{
			Name:             "view with context override",
			Args:             []string{"channel", "view", "-c", "gke"},
			ExpectedResponse: PrintChannelAndVersionToStdout(mctlCfg.Contexts["gke"], "gke") + "\n\n",
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
			assert.Equal(t, expectedResponse, actualResponse)
			BreakupFunc()
		})
	}
	t.Log("View channelCmd Passed")
}

func TestSetCmd(t *testing.T) {
	SetupContextEnv(t)
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
			SetupFunc()
			SystemCmd.SetArgs(tt.Args)
			err = SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}
			actualResponse := b.String()
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
