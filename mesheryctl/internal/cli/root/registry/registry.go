// # Copyright Meshery Authors
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

package registry

import (
	"fmt"

	"errors"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands = []*cobra.Command{generateCmd, publishCmd, updateCmd}

	spreadsheeetID          string
	spreadsheeetCred        string
	modelName               string
	modelCSVFilePath        string
	componentCSVFilePath    string
	relationshipCSVFilePath string
	csvDirectory            string
)

// PublishCmd represents the publish command to publish Meshery Models to Websites, Remote Provider, Meshery
var RegistryCmd = &cobra.Command{
	Use:     "registry",
	Short:   "Model Database",
	Long:    `Manage the state and contents of Meshery’s internal registry of capabilities.`,
	Example: `mesheryctl registry [subcommand]`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.RegistryError(fmt.Sprintf("'%s' is an invalid command.  Use 'mesheryctl registry --help' to display usage guide.\n", args[0]), "registry"))
		}
		return nil
	},
}

func init() {
	RegistryCmd.AddCommand(availableSubcommands...)

}
