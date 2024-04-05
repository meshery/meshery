// Copyright 2024 Layer5, Inc.
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

package models

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	availableSubcommands = []*cobra.Command{pushModelCmd, pullModelCmd}

	username    string
	password    string
	registry    string
	repository  string
	tag         string
	pathToModel string
)

// ModelCmd represents the mesheryctl exp model command
var ModelCmd = &cobra.Command{
	Use:   "model",
	Short: "View list of models and detail of models",
	Long:  "View list of models and detailed information of a specific model",
	Example: `
// To view list of components
mesheryctl exp model list

// To view a specific model
mesheryctl exp model view [model-name]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return system.ErrGetCurrentContext(err)
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return errors.New("Subcommand is required")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemModelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl exp model --help' to display usage guide.\n", args[0]), "model"))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		err = cmd.Usage()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	pushModelCmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication")
	pushModelCmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	pushModelCmd.Flags().StringVarP(&registry, "registry", "r", "", "Registry to push the model to")
	pushModelCmd.Flags().StringVarP(&repository, "repository", "n", "", "Repository name to push the model to")
	pushModelCmd.Flags().StringVarP(&tag, "tag", "t", "", "Tag for the model")
	pushModelCmd.Flags().StringVarP(&pathToModel, "path", "f", "", "Path to the model")
	pullModelCmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication")
	pullModelCmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	pullModelCmd.Flags().StringVarP(&registry, "registry", "r", "", "Registry to push the model to")
	pullModelCmd.Flags().StringVarP(&repository, "repository", "n", "", "Repository name to push the model to")
	pullModelCmd.Flags().StringVarP(&tag, "tag", "t", "", "Tag for the model")

	ModelCmd.AddCommand(availableSubcommands...)
}
