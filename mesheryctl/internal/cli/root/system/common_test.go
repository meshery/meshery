package system

import (
	"bytes"
	"os"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func BreakupFunc() {
	viewCmd.Flags().VisitAll(setFlagValueAsUndefined)
	viewProviderCmd.Flags().VisitAll(setFlagValueAsUndefined)
	channelCmd.Flags().VisitAll(setFlagValueAsUndefined)
	SystemCmd.PersistentFlags().VisitAll(setFlagValueAsUndefined)
	showForAllContext = false
	showProviderForAllContext = false
	tempContext = ""
	utils.SilentFlag = false
}

func setFlagValueAsUndefined(flag *pflag.Flag) {
	_ = flag.Value.Set("")
}

type CmdTestInput struct {
	Name             string
	Args             []string
	ExpectedResponse string
	Token            string
}

func setupContextTestEnv(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Error("unable to locate meshery directory")
	}
	viper.Reset()
	viper.SetConfigFile(path + "/../../../../pkg/utils/TestConfig.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		t.Errorf("unable to read configuration from %v, %v", viper.ConfigFileUsed(), err.Error())
	}

	mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Error("error processing config", err)
	}

	// Initialize the flag validator on SystemCmd so child commands
	// can access it via their command context.
	mesheryctlflags.InitValidators(SystemCmd)
}

func setupSystemOutCmdTest(t *testing.T) *bytes.Buffer {
	buf := utils.SetupMeshkitLoggerTesting(t, true)
	SystemCmd.SetOut(buf)
	return buf
}
