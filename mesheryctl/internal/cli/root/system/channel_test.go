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

func TestViewWithoutAnyParameter(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Error("unable to locate meshery directory")
	}
	viper.SetConfigFile(path + "/../../../../pkg/utils/TestConfig.yaml")
	//fmt.Println(viper.ConfigFileUsed())
	err = viper.ReadInConfig()
	if err != nil {
		t.Error("unable to read configuration from TestConfig.yaml", err)
	}

	mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Error("error processing config", err)
	}

	//fmt.Println(viper.AllKeys())
	b := bytes.NewBufferString("")
	logrus.SetOutput(b)
	logrus.SetFormatter(&utils.OnlyStringFormatterForLogrus{})
	SystemCmd.SetOut(b)
	SystemCmd.SetArgs([]string{"channel", "view"})
	err = SystemCmd.Execute()
	if err != nil {
		t.Error(err)
	}

	actualResponse := b.String()
	expectedResponse := fmt.Sprintf("Context: %v\nChannel: %v\nVersion: %v\n", "local", mctlCfg.GetContextContent().Channel, mctlCfg.GetContextContent().Version)

	if expectedResponse != actualResponse {
		t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
	}
}
