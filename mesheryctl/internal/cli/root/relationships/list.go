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

package relationships

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// represents the mesheryctl exp relationships list command
var listRelationshipsCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered relationships",
	Long:  "List all relationships registered in Meshery Server",
	Example: `
	View list of relationship
    mesheryctl exp relationship list
    View list of relationship with specified page number (25 relationships per page)
    mesheryctl exp relationship list --page 2
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check prerequisites for the command here

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp relationship list \nRun 'mesheryctl exp relationship list --help' to see detailed help message"
		if len(args) != 0 {
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.RelationshipsError("this command takes no arguments\n", "list"))
		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		var url string
		if cmd.Flags().Changed("page") {
			url = fmt.Sprintf("%s/api/meshmodels/relationships?page=%d", baseUrl, pageNumberFlag)
		} else {
			url = fmt.Sprintf("%s/api/meshmodels/relationships?pagesize=all", baseUrl)
		}
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		relationshipsResponse := &models.MeshmodelRelationshipsAPIResponse{}
		err = json.Unmarshal(data, relationshipsResponse)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		header := []string{"Entity Details", "Type"}
		rows := [][]string{}

		for _, relationship := range relationshipsResponse.Relationships {
			if len(relationship.GetEntityDetail()) > 0 {
				rows = append(rows, []string{relationship.GetEntityDetail(), string(relationship.Type())})
			}
		}

		if len(rows) == 0 {
			// if no relationship is found
			fmt.Println("No relationship(s) found")
			return nil
		}

		if cmd.Flags().Changed("page") {
			utils.PrintToTable(header, rows)
		} else {
			maxRowsPerPage := 25
			startIndex := 0
			endIndex := min(len(rows), startIndex+maxRowsPerPage)
			whiteBoardPrinter := color.New(color.FgHiBlack, color.BgWhite, color.Bold)
			for {
				// Clear the entire terminal screen
				utils.ClearLine()

				// Print number of models and current page number
				whiteBoardPrinter.Print("Total number of relationships: ", len(rows))
				fmt.Println()
				whiteBoardPrinter.Print("Page: ", startIndex/maxRowsPerPage+1)
				fmt.Println()
				whiteBoardPrinter.Print("Total pages: ", len(rows)/maxRowsPerPage+1)
				fmt.Println()
				whiteBoardPrinter.Println("Press ↑/← or ↓/→ to navigate, Esc or Ctrl+C to exit")
				fmt.Println()

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

				// Navigate through the list of relationships
				if event.Key == keyboard.KeyArrowRight || event.Key == keyboard.KeyArrowDown {
					startIndex += maxRowsPerPage
					endIndex = min(len(rows), startIndex+maxRowsPerPage)
				}

				if event.Key == keyboard.KeyArrowUp || event.Key == keyboard.KeyArrowLeft {
					startIndex -= maxRowsPerPage
					if startIndex < 0 {
						startIndex = 0
					}
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

func init() {
	// Add the new exp relationship commands to the listRelationshipsCmd
	listRelationshipsCmd.Flags().IntVarP(&pageNumberFlag, "page", "p", 1, "(optional) List next set of models with --page (default = 1)")
}
