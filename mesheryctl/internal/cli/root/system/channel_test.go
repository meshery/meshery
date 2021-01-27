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

	//fmt.Println(viper.AllKeys())
	b = bytes.NewBufferString("")
	logrus.SetOutput(b)
	logrus.SetFormatter(&utils.OnlyStringFormatterForLogrus{})
}

func TestViewWithoutAnyParameter(t *testing.T) {
	SetupFunc(t)
	viewCmd.SetOut(b)
	viewCmd.SetArgs([]string{})
	err = viewCmd.Execute()
	if err != nil {
		t.Error(err)
	}

	actualResponse := b.String()
	expectedResponse := PrintChannelAndVersionToStdout(mctlCfg.GetContextContent(), "local") + "\n"

	if expectedResponse != actualResponse {
		t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
	}
}

func TestViewWithCorrectContextOverride(t *testing.T) {
	SetupFunc(t)
	viewCmd.SetOut(b)
	viewCmd.SetArgs([]string{"-c", "gke"})
	err = viewCmd.Execute()
	if err != nil {
		t.Error(err)
	}

	actualResponse := b.String()
	expectedResponse := PrintChannelAndVersionToStdout(mctlCfg.Contexts["gke"], "gke") + "\n"

	if expectedResponse != actualResponse {
		t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
	}
}

func TestViewWithAllFlag(t *testing.T) {
	SetupFunc(t)
	viewCmd.SetOut(b)
	viewCmd.SetArgs([]string{"--all"})
	err = viewCmd.Execute()
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

func TestViewWithoutAnyParameterTemp(t *testing.T) {
	SetupFunc(t)
	viewCmd.SetOut(b)
	viewCmd.SetArgs([]string{"channel", "view"})
	err = viewCmd.Execute()
	if err != nil {
		t.Error(err)
	}

	actualResponse := b.String()
	expectedResponse := PrintChannelAndVersionToStdout(mctlCfg.GetContextContent(), "local") + "\n"

	if expectedResponse != actualResponse {
		t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
	}
}
