package system

import (
	"bytes"
	"fmt"
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

func PrintSetChannel(stable, Edge string) string {
	var ctx config.Context
	return fmt.Sprintf("Channel set to %s-%s", ctx.Channel, ctx.Version)
}

type CmdTestInput struct {
	Name             string
	Args             []string
	ExpectedResponse string
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
}

// func TestSetCmd(t *testing.T) {
// 	tests := []CmdTestInput{
// 		{
// 			Name:             "Set Channel",
// 			Args:             []string{"channel", "set", "-c", "local", "latest"},
// 			ExpectedResponse: PrintSetChannel("stable", "latest"),
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.Name, func(t *testing.T) {
// 			SetupFunc()
// 			SystemCmd.SetArgs(tt.Args)
// 			err = SystemCmd.Execute()
// 			if err != nil {
// 				t.Error(err)
// 			}

// 			actualResponse := b.String()
// 			expectedResponse := tt.ExpectedResponse
// 			assert.Equal(t, expectedResponse, actualResponse)
// 			BreakupFunc()
// 		})
// 	}
// }