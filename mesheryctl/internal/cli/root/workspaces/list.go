// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package workspaces

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eiannone/keyboard"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listWorkspaceCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered workspaces",
	Long:  `List name of all registered workspaces`,
	Example: `
// List all registered workspace
mesheryctl exp workspace list --orgId [orgId]

// Documentation for workspace can be found at:
https://docs.layer5.io/cloud/spaces/workspaces/

`,

	Args: func(cmd *cobra.Command, args []string) error {
		// Check if the orgID is provided
		orgIdFlag, _ := cmd.Flags().GetString("orgId")
		if orgIdFlag == "" {
			if err := cmd.Usage(); err != nil {
				return nil
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

		url := fmt.Sprintf("%s/api/workspaces?orgID=%s", baseUrl, orgID)
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

		workspaceResponse := &models.WorkspacePage{}
		err = json.Unmarshal(data, workspaceResponse)
		if err != nil {
			return err
		}
		header := []string{"ID", "Name", "Organization ID", "Description", "Created At", "Updated At"}
		rows := [][]string{}
		for _, workspace := range workspaceResponse.Workspaces {
			rows = append(rows, []string{workspace.ID.String(), workspace.Name, workspace.OrganizationID.String(), workspace.Description, workspace.CreatedAt.String(), workspace.UpdatedAt.String()})
		}

		if len(rows) == 0 {
			utils.Log.Info("No workspace found")
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

				// Print number of workspaces and current page number
				whiteBoardPrinter.Println("Total number of workspaces: ", len(rows))
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
					utils.Log.Error(errors.New("error reading keyboard event"))
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
