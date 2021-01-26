package system

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var b *bytes.Buffer

func SetupFunc(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Error("unable to locate meshery directory")
	}
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

	//fmt.Println(viper.AllKeys())
	b = bytes.NewBufferString("")
	logrus.SetOutput(b)
	logrus.SetFormatter(&utils.OnlyStringFormatterForLogrus{})
	SystemCmd.SetOut(b)
}

func TestViewWithoutAnyParameter(t *testing.T) {
	SetupFunc(t)
	SystemCmd.SetArgs([]string{"channel", "view"})
	err = SystemCmd.Execute()
	if err != nil {
		t.Error(err)
	}

	actualResponse := b.String()
	expectedResponse := PrintChannelAndVersionToStdout(mctlCfg.GetContextContent(), "local") + "\n"

	if expectedResponse != actualResponse {
		t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
	}
}

func TestViewWithAllFlag(t *testing.T) {
	SetupFunc(t)
	SystemCmd.SetArgs([]string{"channel", "view", "--all"})
	err = SystemCmd.Execute()
	if err != nil {
		t.Error(err)
	}

	actualResponse := b.String()
	expectedResponse := ""
	for k, v := range mctlCfg.Contexts {
		expectedResponse += PrintChannelAndVersionToStdout(v, k) + "\n"
	}
	expectedResponse += fmt.Sprintf("Current Context: %v\n", mctlCfg.CurrentContext)

	if expectedResponse != actualResponse {
		t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
	}
}
