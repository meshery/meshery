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

package credentials

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// Maximum number of rows to be displayed in a page
	maxRowsPerPage = 25

	// Color for the whiteboard printer
	whiteBoardPrinter = color.New(color.FgHiBlack, color.BgWhite, color.Bold)

	availableSubcommands = []*cobra.Command{listCredentialCmd, createCredentialCmd}
)
var CredentialCmd = &cobra.Command{
	Use:   "credential",
	Short: "Manage credentials",
	Long: `View and manage list of credentials for Meshery Environment.
Find more information at: https://docs.meshery.io/reference/mesheryctl#command-reference`,
	Example: `
// Base command for credentials:
mesheryctl exp credential [subcommands]
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.CredentialsError(fmt.Sprintf("invalid subcommand: %s. Please provide options from [view]. Use 'mesheryctl exp credential --help' to display usage guide.\n", args[0]), "credential"))
		}
		return nil
	},
}

func init() {
	createCredentialCmd.Flags().StringP("name", "n", "", "Name of the credential")
	createCredentialCmd.Flags().StringP("user-id", "u", "", "User ID of the credential")
	createCredentialCmd.Flags().StringP("type", "t", "", "Type of the credential")
	createCredentialCmd.Flags().StringP("secret", "s", "", "Secret of the credential")

	CredentialCmd.AddCommand(availableSubcommands...)
}
