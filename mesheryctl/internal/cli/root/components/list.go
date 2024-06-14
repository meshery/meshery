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
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// represents the mesheryctl components list command
var listComponentCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered components",
	Long:  "List all components registered in Meshery Server",
	Example: `
// View list of components
mesheryctl components list

// View list of components with specified page number (25 components per page)
mesheryctl components list --page 2

// To view the number of components present in Meshery
mesheryctl components list --count
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
		url := fmt.Sprintf("%s/api/meshmodels/components?%s", baseUrl, utils.GetPageQueryParameter(cmd, pageNumberFlag))

		return listComponents(cmd, url, false)
	},
}

func init() {
	// Add the new components commands to the ComponentsCmd
	listComponentCmd.Flags().IntVarP(&pageNumberFlag, "page", "p", 1, "(optional) List next set of components with --page (default = 1)")
}
