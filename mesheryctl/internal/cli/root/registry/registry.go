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

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands = []*cobra.Command{generateCmd, publishCmd, updateCmd, purgeCmd}

	spreadsheeetID          string
	spreadsheeetCred        string
	modelName               string
	modelCSVFilePath        string
	componentCSVFilePath    string
	relationshipCSVFilePath string
	csvDirectory            string
)

// RegistryCmd represents the registry command, which manages the state and contents of Meshery's internal registry of capabilities.
var RegistryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage the capability registry",
	Long: `Generate, publish, update, and purge models and components in Meshery's capability registry.
Find more information at: https://docs.meshery.io/reference/mesheryctl/registry`,
	Example: `
// Generate models from a spreadsheet
mesheryctl registry generate --spreadsheet-id [spreadsheet-id] --spreadsheet-cred [base64-encoded-csv-content]

// Publish models to Meshery Cloud / Remote Provider
mesheryctl registry publish

// Update models from spreadsheets or website
mesheryctl registry update

// Dry run purging models directory, keeping only the latest version of each model
mesheryctl registry purge ./models --dry-run

// Purge older model versions, retaining the latest 2 versions per model
mesheryctl registry purge ./models --retain 2
`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		setupRegistryLogger()
		return nil
	},

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

func setupRegistryLogger() {
	err := meshkitRegistryUtils.SetLogger(true)
	if err != nil {
		utils.Log.Info("Error setting logger: ", err)
	}
}
