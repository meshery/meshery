// Copyright Meshery Authors
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
	"log"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	mesheryctllogger "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/logger"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/adapter"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/components"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/connections"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/design"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/environments"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/experimental"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/filter"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/model"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/organizations"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/perf"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/registry"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/relationships"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/system"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/workspaces"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	logrus "github.com/sirupsen/logrus"

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
	Long: `As a self-service engineering platform, Meshery enables collaborative design and operation of cloud native infrastructure.
Find more information at: https://docs.meshery.io/reference/mesheryctl#command-reference`,
	Example: `
// Base command:
mesheryctl

// Display help about command/subcommand:
mesheryctl --help
mesheryctl system start --help
// For viewing verbose output:
mesheryctl -v [or] --verbose
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize a validator and add it to the command context
		// This allows us to use the same validator instance across all subcommands and avoid initializing multiple instances of the validator
		mesheryctlflags.InitValidators(cmd)
		return nil
	},
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
// Execute adds all child commands
func Execute() error {
	//log formatter for improved UX
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
		design.DesignCmd,
		perf.PerfCmd,
		adapter.AdapterCmd,
		experimental.ExpCmd,
		filter.FilterCmd,
		registry.RegistryCmd,
		components.ComponentCmd,
		model.ModelCmd,
		environments.EnvironmentCmd,
		connections.ConnectionsCmd,
		organizations.OrgCmd,
		relationships.RelationshipCmd,
		workspaces.WorkSpaceCmd,
	}

	RootCmd.AddCommand(availableSubcommands...)
	RootCmd.SetHelpCommand(newHelpCommand())
}

func newHelpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "help [command]",
		Short: "Show help for any command",
		Long:  "Show help for any command.",
		Run: func(c *cobra.Command, args []string) {
			cmd, _, err := c.Root().Find(args)
			if cmd == nil || err != nil {
				c.Println(c.UsageString())
				return
			}
			cmd.InitDefaultHelpFlag()
			cmd.HelpFunc()(cmd, args)
		},
	}
}

func TreePath() *cobra.Command {
	return RootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	utils.CfgFile = cfgFile

	utils.SetKubeConfig()

	if cfgFile != utils.DefaultConfigPath {

		viper.SetConfigFile(cfgFile)

	} else {

		needsMutation, err := config.NeedsMutation(utils.DefaultConfigPath)
		if err != nil {
			utils.Log.Fatal(err)
		}

		if needsMutation {
			err = config.InitDefaultConfig(
				utils.DefaultConfigPath,
				utils.MesheryFolder,
				utils.TemplateToken,
				utils.TemplateContext,
				utils.CreateConfigFile,
			)

			if err != nil {
				utils.Log.Fatal(err)
			}
		}
	}

	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		utils.Log.Debugf("unable to read config file: %v", err)
	} else {
		utils.Log.Debug("Using config file:", viper.ConfigFileUsed())
	}
}

func setupLogger() {
	// default log level is info
	logLevel := logrus.InfoLevel

	// log level to debug if the -v flag is set
	if verbose {
		logLevel = logrus.DebugLevel
	}

func setupLogger() {

	utils.Log = utils.SetupMeshkitLogger("mesheryctl", verbose, RootCmd.OutOrStdout())

}
