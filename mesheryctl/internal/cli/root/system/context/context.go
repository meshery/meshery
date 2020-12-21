// Copyright 2020 Layer5, Inc.
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

package context

import (
	"fmt"
	"os"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configuration models.MesheryCtlConfig

var (
	availableSubcommands []*cobra.Command
)

// ContextCmd represents the update command
var ContextCmd = &cobra.Command{
	Use:   "context",
	Short: "Operate on multiple meshery contexts",
	Long:  `Operate between different meshery server version hosted on different URLs.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\"", args[0])))
		}
		return nil
	},
}

func init() {
	homeDir, err := os.UserHomeDir()
	viper.SetConfigName("config")
	viper.AddConfigPath(homeDir + string(os.PathSeparator) + utils.MesheryFolder)
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	if err != nil {
		log.Fatalf("Unable to Access home directory, %s", err)
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file %s", err)
	}
	err = viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode struct, %s", err)
	}
	availableSubcommands = []*cobra.Command{
		createContextCmd,
		deleteContextCmd,
		switchContextCmd,
		viewContextCmd,
	}
	ContextCmd.AddCommand(availableSubcommands...)
}
