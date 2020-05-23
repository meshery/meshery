// Copyright 2019 The Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package root

import (
	"errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/perf"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//TerminalFormatter is exported
type TerminalFormatter struct{}

var (
	cfgFile     string
	mctlCfgFile string
)

var cmdDetails = `
Meshery is the service mesh management plane, providing lifecycle, performance, and configuration management of service meshes and their workloads.

Usage:
  mesheryctl [command]

Available Commands:
  help        Help about any command
  perf        Performance Management 
  system      Meshery Lifecyle Management
  version     Print mesheryctl version


Flags:
      --config string      config file (default location is: $HOME/.meshery/` + utils.DockerComposeFile + `)
      --mesheryctl-config  mesheryctl config file (default location is: <unset>. Uses default config.)
  -h, --help               help for mesheryctl
  -v, --version            version of mesheryctl
  -d, --debug              enable debug logging

Use "mesheryctl [command] --help" for more information about a command.
`

//Format is exported
func (f *TerminalFormatter) Format(entry *log.Entry) ([]byte, error) {
	return append([]byte(entry.Message), '\n'), nil
}

var (
	availableSubcommands = []*cobra.Command{}
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mesheryctl",
	Short: "Meshery Command Line tool",
	Long:  `Meshery is the service mesh management plane, providing lifecycle, performance, and configuration management of service meshes and their workloads.`,
	//Args:  cobra.MinimumNArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		b, _ := cmd.Flags().GetBool("version")
		if b {
			versionCmd.Run(nil, nil)
			return nil
		}

		if len(args) == 0 {
			log.Print(cmdDetails)
			return nil
		}

		for _, subcommand := range availableSubcommands {
			if args[0] == subcommand.Name() {
				return nil
			}
		}

		return errors.New("sub-command not found : " + "\"" + args[0] + "\"")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	log.SetLevel(log.InfoLevel)

	if debug, err := RootCmd.Flags().GetBool("debug"); err == nil && debug {
		log.SetLevel(log.DebugLevel)
	}

	//log formatter for improved UX
	log.SetFormatter(new(TerminalFormatter))
	_ = RootCmd.Execute()
}

func init() {
	utils.SetFileLocation() //from vars.go
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default location is: "+utils.DockerComposeFile+")")
	RootCmd.PersistentFlags().StringVar(&mctlCfgFile, "mesheryctl-config", "", "mesheryctl config file to override defaults (default file: <unset>")

	// Preparing for an "edge" channel
	// RootCmd.PersistentFlags().StringVar(&cfgFile, "edge", "", "flag to run Meshery as edge (one-time)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("version", "v", false, "Version flag")
	RootCmd.Flags().BoolP("debug", "d", false, "Debug flag")

	availableSubcommands = []*cobra.Command{
		versionCmd,
		system.SystemCmd,
		perf.PerfCmd,
	}

	RootCmd.AddCommand(availableSubcommands...)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Use default ".meshery" folder location.
		viper.AddConfigPath(utils.MesheryFolder)
		log.Debug("initConfig: ", utils.MesheryFolder)
		viper.SetConfigName("meshery")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	} else {
		log.WithError(err).Errorf("failed to read in meshery configv")
	}

	// Read in mesheryctl config or use defaults
	// Default config:
	//  baseMesheryURL: "http://localhost:9081/api",
	//  perf:
	//	  authTokenURI:    "/gettoken",
	//	  loadTestSmpsURI: "/load-test-smps",
	if mctlCfgFile != "" {
		viper.SetConfigFile(mctlCfgFile)
		if err := viper.ReadInConfig(); err == nil {
			log.Debugf("Using mesheryctl config file: %s", viper.ConfigFileUsed())
		} else {
			log.Errorf("failed to read in mesheryctl config - %v", err)
		}
	} else {
		setMesheryctlConfigDefaults(map[string]interface{}{
			"baseMesheryURL": "http://localhost:9081/api",
			"perf": map[string]interface{}{
				"authTokenURI":    "/gettoken",
				"loadTestSmpsURI": "/load-test-smps",
			},
		})
	}
}

// setMesheryctlConfigDefaults loads the hardcoded defaults in to viper kv store
func setMesheryctlConfigDefaults(defaults map[string]interface{}) {
	v := viper.GetViper()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
}
