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
	availableSubcommands = []*cobra.Command{importCmd, publishCmd}

	spreadsheeetID   string
	spreadsheeetCred string
)

// PublishCmd represents the publish command to publish Meshery Models to Websites, Remote Provider, Meshery
var RegistryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Meshery Registry Management",
	Long:  `Manage the state and configuration of Meshery Registry.`,
	Example: `
	mesheryctl registry [subcommand]
	`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("'%s' is an invalid command.  Use 'mesheryctl registry --help' to display usage guide.\n", args[0])))
		}
		return nil
	},
}

func init() {
	RegistryCmd.AddCommand(availableSubcommands...)
	importCmd.PersistentFlags().StringVar(&pathToRegistrantConnDefinition, "registrant_def", "", "path pointing to the registrant connection definition")
	importCmd.PersistentFlags().StringVar(&pathToRegistrantCredDefinition, "registrant_cred", "", "path pointing to the registrant credetial definition")

	importCmd.MarkFlagsRequiredTogether("registrant_def", "registrant_cred")

	importCmd.MarkFlagsMutuallyExclusive("spreadsheet_id", "registrant_def")
	// importCmd.MarkFlagsMutuallyExclusive("spreadsheet_cred", "registrant_cred")
	importCmd.PersistentFlags().StringVarP(&outputLocation, "output", "o", "../server/meshmodels", "location to output generated models, defaults to ../server/meshmodels")

}
