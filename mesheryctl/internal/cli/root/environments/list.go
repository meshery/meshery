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

package environments

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eiannone/keyboard"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models/environments"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listEnvironmentCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered environments",
	Long:  `List name of all registered environments`,
	Example: `
// List all registered environment
mesheryctl exp environment list --orgID [orgId]

// Documentation for environment can be found at:
https://docs.layer5.io/cloud/spaces/environments/

`,

	Args: func(cmd *cobra.Command, args []string) error {
		// Check if all flag is set
		orgIdFlag, _ := cmd.Flags().GetString("orgId")

		if orgIdFlag == "" {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return utils.ErrInvalidArgument(errors.New("Please provide a --orgId flag"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()

		url := fmt.Sprintf("%s/api/environments?orgID=%s", baseUrl, orgID)
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		environmentResponse := &environments.EnvironmentPage{}
		err = json.Unmarshal(data, environmentResponse)
		if err != nil {
			return err
		}
		header := []string{"ID", "Name", "Organization ID", "Description", "Created At", "Updated At"}
		rows := [][]string{}
		for _, environment := range environmentResponse.Environments {
			rows = append(rows, []string{environment.ID.String(), environment.Name, environment.OrganizationID.String(), environment.Description, environment.CreatedAt.String(), environment.UpdatedAt.String()})
		}

		if len(rows) == 0 {
			utils.Log.Info("No environment found")
			return nil
		}

		if cmd.Flags().Changed("page") {
			utils.PrintToTable(header, rows)
		} else {
			startIndex := 0
			endIndex := min(len(rows), startIndex+maxRowsPerPage)
			for {
				// Clear the entire terminal screen
				utils.ClearLine()

				// Print number of environments and current page number
				whiteBoardPrinter.Println("Total number of environments: ", len(rows))
				whiteBoardPrinter.Println("Page: ", startIndex/maxRowsPerPage+1)

				whiteBoardPrinter.Println("Press Enter or ↓ to continue. Press Esc or Ctrl+C to exit.")

				utils.PrintToTable(header, rows[startIndex:endIndex])
				keysEvents, err := keyboard.GetKeys(10)
				if err != nil {
					return err
				}

				defer func() {
					_ = keyboard.Close()
				}()

				event := <-keysEvents
				if event.Err != nil {
					utils.Log.Error(fmt.Errorf("unable to capture keyboard events"))
					break
				}

				if event.Key == keyboard.KeyEsc || event.Key == keyboard.KeyCtrlC {
					break
				}

				if event.Key == keyboard.KeyEnter || event.Key == keyboard.KeyArrowDown {
					startIndex += maxRowsPerPage
					endIndex = min(len(rows), startIndex+maxRowsPerPage)
				}

				if startIndex >= len(rows) {
					break
				}
			}
		}
		return nil
	},
}
