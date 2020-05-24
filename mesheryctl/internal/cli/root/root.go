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
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/perf"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//TerminalFormatter is exported
type TerminalFormatter struct{}

var cfgFile string

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
	Args:  cobra.MinimumNArgs(1),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := cmd.Flags().GetBool("version")
		if err != nil {
			return err
		}

		if b {
			versionCmd.Run(nil, nil)
			return nil
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	//log formatter for improved UX
	log.SetFormatter(new(TerminalFormatter))
	_ = RootCmd.Execute()
}

func init() {
	err := utils.SetFileLocation()
	if err != nil {
		log.Fatal(err)
	}
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default location is: "+utils.DockerComposeFile+")")

	// Preparing for an "edge" channel
	// RootCmd.PersistentFlags().StringVar(&cfgFile, "edge", "", "flag to run Meshery as edge (one-time)")

	RootCmd.Flags().BoolP("version", "v", false, "Version flag")

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
		viper.SetConfigFile(utils.DockerComposeFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Fatal(err)
	}
}
