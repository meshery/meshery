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
	"encoding/json"
	"fmt"
	"io"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Maximum number of rows to be displayed in a page
	maxRowsPerPage = 25

	// Color for the whiteboard printer
	whiteBoardPrinter = color.New(color.FgHiBlack, color.BgWhite, color.Bold)
)

var listCredentialCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the credentials",
	Long:  `Display list of all the available credentials`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// List all the credentials
mesheryctl exp credential list
`,

	RunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		req, err := utils.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/integrations/credentials", nil)
		if err != nil {
			return err
		}
		res, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			utils.Log.Error(utils.ErrReadResponseBody(err))
			return nil
		}
		credentialResponse := models.CredentialsPage{}
		err = json.Unmarshal(body, &credentialResponse)
		if err != nil {
			return utils.ErrUnmarshal(err)
			
		}
		header := []string{"ID", "User-Id", "Name", "Type", "Secrets", "Created At", "Updated At"}
		data := [][]string{}
		for _, credential := range credentialResponse.Credentials {
			secret := fmt.Sprintf("%v", credential.Secret)
			data = append(data, []string{credential.ID.String(), credential.Name, credential.Type, secret, credential.CreatedAt.String(), credential.UpdatedAt.String()})
		}
		if len(data) == 0 {
			utils.Log.Info("No credentials found")
			return nil
		}
		if cmd.Flags().Changed("page") {
			utils.PrintToTable(header, data)
		} else {
			startIndex := 0
			endIndex := min(len(data), startIndex+maxRowsPerPage)
			for {
				// Clear the entire terminal screen
				utils.ClearLine()
				// Print the number of credentials and current page number
				whiteBoardPrinter.Println("Total Credentials: ", len(data))
				whiteBoardPrinter.Println("Page: ", startIndex/maxRowsPerPage+1)
				whiteBoardPrinter.Println("Press Enter or ↓ to continue, Esc or Ctrl+C (Ctrl+Cmd for OS user) to exit")

				// Print the table
				utils.PrintToTable(header, data[startIndex:endIndex])
				// Wait for the user to press a key
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
					endIndex = min(len(data), startIndex+maxRowsPerPage)
				}

				if startIndex >= len(data) {
					break
				}
			}
		}
		return nil
	},
}
