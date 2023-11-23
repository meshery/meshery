// Copyright 2023 Layer5, Inc.
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
	"fmt"
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/app"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/experimental"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/filter"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/mesh"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/pattern"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/perf"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose = false
)

var (
	availableSubcommands = []*cobra.Command{}
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mesheryctl",
	Short: "Meshery Command Line tool",
	Long:  `As a self-service engineering platform, Meshery enables collaborative design and operation of cloud native infrastructure.`,
	Example: `
// Base command
mesheryctl

// Display help about command/subcommand
mesheryctl --help
mesheryctl system start --help

// For viewing verbose output
mesheryctl -v [or] --verbose
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.RootError(fmt.Sprintf("'%s' is an invalid command. Use 'mesheryctl --help' to display usage guide.\n", args[0])))
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() error {
	//log formatter for improved UX
	utils.SetupLogrusFormatter()
	// Removing printing command usage on error
	RootCmd.SilenceUsage = true
	err := RootCmd.Execute()
	return err
}

func init() {
	err := utils.SetFileLocation()
	if err != nil {
		log.Fatal(err)
	}

	cobra.OnInitialize(setVerbose)
	cobra.OnInitialize(setupLogger)
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", utils.DefaultConfigPath, "path to config file")

	// Preparing for an "edge" channel
	// RootCmd.PersistentFlags().StringVar(&cfgFile, "edge", "", "flag to run Meshery as edge (one-time)")

	// global verbose flag for verbose logs
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	availableSubcommands = []*cobra.Command{
		completionCmd,
		versionCmd,
		system.SystemCmd,
		pattern.PatternCmd,
		perf.PerfCmd,
		mesh.MeshCmd,
		app.AppCmd,
		experimental.ExpCmd,
		filter.FilterCmd,
	}

	RootCmd.AddCommand(availableSubcommands...)
}

func TreePath() *cobra.Command {
	return RootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	utils.CfgFile = cfgFile
	// initialize the path to the kubeconfig file
	utils.SetKubeConfig()
	// Allow user to override config file with use of --config global flag
	if cfgFile != utils.DefaultConfigPath {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		// Otherwise, use the default `config.yaml` config file
	} else {
		stat, err := os.Stat(utils.DefaultConfigPath)
		if !os.IsNotExist(err) && stat.Size() == 0 {
			log.Println("Empty meshconfig. Please populate it before running a command")
		}
		if os.IsNotExist(err) {
			log.Printf("Missing Meshery config file.")
		}

		// Create a default meshconfig in each of the above two scenarios.
		if os.IsNotExist(err) || (!os.IsNotExist(err) && stat.Size() == 0) {
			// Check for Meshery existence and permission of application folder
			if _, err := os.Stat(utils.MesheryFolder); err != nil {
				if os.IsNotExist(err) {
					err = os.MkdirAll(utils.MesheryFolder, 0775)
					if err != nil {
						log.Fatal(err)
					}
				}
			}

			// Create config file if not present in meshery folder
			err = utils.CreateConfigFile()
			if err != nil {
				log.Fatal(err)
			}

			// Add Token to context file
			err = config.AddTokenToConfig(utils.TemplateToken, utils.DefaultConfigPath)
			if err != nil {
				log.Fatal(err)
			}

			// Add Context to context file
			err = config.AddContextToConfig("local", utils.TemplateContext, utils.DefaultConfigPath, true, false)
			if err != nil {
				log.Fatal(err)
			}

			log.Println(
				fmt.Sprintf("Default config file created at %s",
					utils.DefaultConfigPath,
				))
		}
		viper.SetConfigFile(utils.DefaultConfigPath)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	}
}

// setVerbose sets the log level to debug if the -v flag is set
func setVerbose() {
	log.SetLevel(log.InfoLevel)

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func setupLogger() {
	utils.SetupMeshkitLogger(verbose, nil)
}
