package context

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
)

var b *bytes.Buffer
var mctlCfg *config.MesheryCtlConfig

func SetupFunc(t *testing.T) {
	b = bytes.NewBufferString("")
	logrus.SetOutput(b)
	utils.SetupLogrusFormatter()
	viewContextCmd.SetOut(b)
}
func BreakupFunc(t *testing.T) {
	viewContextCmd.Flags().VisitAll(setFlagValueAsUndefined)
}
func setFlagValueAsUndefined(flag *pflag.Flag) {
	_ = flag.Value.Set("")
}
func SetupContextEnv(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Error("unable to locate meshery directory")
	}
	viper.Reset()
	viper.SetConfigFile(path + "/../../../../../pkg/utils/TestContext1.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		t.Errorf("unable to read configuration from %v, %v", viper.ConfigFileUsed(), err.Error())
	}

	mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Error("error processing config", err)
	}

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
			Name: "view default context when context has token field present",
			Args: []string{"view"},
			ExpectedResponse: `
Current Context: local

endpoint: http://localhost:9081
token: Default
token-location: auth.json
platform: docker
adapters:
- meshery-istio
channel: stable
version: v0.5.10

`,
		},
		{
			Name: "view default context when context has token field absent",
			Args: []string{"view"},
			ExpectedResponse: `
Current Context: local

endpoint: http://localhost:9081
token: Default
token-location: auth.json
platform: docker
adapters:
- meshery-istio
channel: stable
version: v0.5.10

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			SetupFunc(t)
			ContextCmd.SetArgs(tt.Args)
			err := viewContextCmd.Execute()
			if err != nil {
				t.Error(err)
			}
			actualResponse := b.String()
			expectedResponse := tt.ExpectedResponse
			if expectedResponse != actualResponse {
				t.Errorf("expected response [%v] and actual response [%v] don't match", expectedResponse, actualResponse)
			}
			BreakupFunc(t)
		})
	}
}

func printContext(ctx *config.Context, name string) string {
	return fmt.Sprintf("")
}
