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

package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//TerminalFormatter is exported
type TerminalFormatter struct{}

var cfgFile string

var cmdDetails = `
Meshery is the service mesh management plane, providing lifecycle, performance, and configuration management of service meshes and their workloads.

Usage:
  mesheryctl [command]

Available Commands:
  cleanup     Clean up Meshery
  help        Help about any command
  logs        Print logs
  perf        Performance testing and benchmarking
  start       Start Meshery
  status      Check Meshery status
  stop        Stop Meshery
  update      Pull new Meshery images from Docker Hub
  version     Version of mesheryctl

Flags:
      --config string   config file (default location is: $HOME/.meshery/config.yaml)
  -h, --help            help for mesheryctl
  -t, --toggle          Help message for toggle
  -v, --version         Version of mesheryctl

Use "mesheryctl [command] --help" for more information about a command.
`

//Format is exported
func (f *TerminalFormatter) Format(entry *log.Entry) ([]byte, error) {
	return append([]byte(entry.Message), '\n'), nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mesheryctl",
	Short: "Meshery Command Line tool",
	Long:  `Meshery is the service mesh management plane, providing lifecycle, performance, and configuration management of service meshes and their workloads.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		b, _ := cmd.Flags().GetBool("version")
		if b {
			versionCmd.Run(nil, nil)
		}
		if len(args) == 0 {
			log.Print(cmdDetails)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	fmt.Println("\n")
	//log formatter for improved UX
	log.SetFormatter(new(TerminalFormatter))
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("\n")
}

func init() {
	setFileLocation() //from vars.go
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default location is: "+dockerComposeFile+")")

	// Preparing for an "edge" channel
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "edge", "", "flag to run Meshery as edge (one-time)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "v", false, "Version flag")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Use default ".meshery" folder location.
		viper.AddConfigPath(mesheryFolder)
		log.Debug("initConfig: ", mesheryFolder)
		viper.SetConfigName("config.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
