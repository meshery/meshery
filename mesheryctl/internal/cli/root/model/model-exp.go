// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by a, filepath.Dir(${1:}modelDefPathpplicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// TODO this is a temporal container, will be delete as soon as model init will be moved from exp to model.
// Had to create this container, otherwise if directly add model init to exp, a cyclical dependency in the tests appears.
// plus with such arrangement the access to model init in exp is through model, like exp model init.

var (
	// Available model subcommads
	modelExpAvailableSubcommands = []*cobra.Command{initModelCmd, buildModelCmd}
)

// ModelCmd represents the mesheryctl model command
var ModelExpCmd = &cobra.Command{
	Use:   "model",
	Short: "Experimental commands for mesheryctl model",
	Long:  `Temporal container to deliver model init and model build subcommands to exp command`,
	Example: `
// Scaffold a folder structure for model creation
mesheryctl exp model init [model-name]

// Create an OCI-compliant package from the model files
mesheryctl exp model build [path/to/model/version/folder]
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return utils.ErrInvalidArgument(errors.New("please provide a subcommand"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(modelExpAvailableSubcommands, args[0]); !ok {
			return errors.New(utils.ExpError(fmt.Sprintf("'%s' is an invalid command. Use 'mesheryctl exp model --help' to display usage guide.'\n", args[0])))
		}
		return nil
	},
}

func init() {
	ModelExpCmd.AddCommand(modelExpAvailableSubcommands...)
}
