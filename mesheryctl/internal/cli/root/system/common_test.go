package system

import (
	"bytes"
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
	providerViewFlags.All = false
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

// setupContextTestEnv points viper at a private copy of the shared meshconfig
// fixture rather than the fixture itself, because commands exercised from this
// package persist the active meshconfig through viper.WriteConfig. See
// utils.CopyMeshconfigFixture for what writing the shared file does to the
// mesheryctl packages running alongside this one.
func setupContextTestEnv(t *testing.T) {
	viper.Reset()

	configPath := utils.CopyMeshconfigFixture(t, utils.SharedTestConfigPath(t))
	viper.SetConfigFile(configPath)

	// `system login` and `system reset` persist the meshconfig through
	// utils.DefaultConfigPath rather than through viper. Pointing only viper at
	// the copy would leave those two writing wherever that variable was last
	// set - the shared fixture among the possibilities, which is precisely what
	// the copy exists to keep them away from.
	defaultConfigPath := utils.DefaultConfigPath
	utils.DefaultConfigPath = configPath
	t.Cleanup(func() { utils.DefaultConfigPath = defaultConfigPath })

	err := viper.ReadInConfig()
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
