// Copyright 2024 Meshery Authors
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

package components

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// represents the mesheryctl exp components list command
var listComponentCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered components",
	Long:  "List all components registered in Meshery Server",
	Example: `
	// View list of components
mesheryctl exp components list
// View list of components with specified page number (25 components per page)
mesheryctl exp components list --page 2
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
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemModelSubError("this command takes no arguments\n", "list"))
		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		var url string
		if cmd.Flags().Changed("page") {
			url = fmt.Sprintf("%s/api/meshmodels/components?page=%d", baseUrl, pageNumberFlag)
		} else {
			url = fmt.Sprintf("%s/api/meshmodels/components?pagesize=all", baseUrl)
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

		componentsResponse := &models.MeshmodelComponentsAPIResponse{}
		err = json.Unmarshal(data, componentsResponse)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		header := []string{"Model", "kind", "Version"}
		rows := [][]string{}

		for _, component := range componentsResponse.Components {
			if len(component.DisplayName) > 0 {
				rows = append(rows, []string{component.Model.Name, component.Component.Kind, component.Component.Version})
			}
		}

		if len(rows) == 0 {
			// if no component is found
			fmt.Println("No components(s) found")
			return nil
		}

		if cmd.Flags().Changed("page") {
			utils.PrintToTable(header, rows)
		} else {
			maxRowsPerPage := 25
			err := utils.HandlePagination(maxRowsPerPage, "components", rows, header)
			if err != nil {
				utils.Log.Error(err)
				return err
			}
		}
		return nil
	},
}

func init() {
	// Add the new exp components commands to the ComponentsCmd
	listComponentCmd.Flags().IntVarP(&pageNumberFlag, "page", "p", 1, "(optional) List next set of models with --page (default = 1)")
}
