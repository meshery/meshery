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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// represents the mesheryctl component list command
var listComponentCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered components",
	Long: `List all components registered in Meshery Server
Documentation for components can be found at https://docs.meshery.io/reference/mesheryctl/component/list`,
	Example: `
// View list of components
mesheryctl component list

// View list of components with specified page number (25 components per page)
mesheryctl component list --page [page-number]

// Display the number of components present in Meshery
mesheryctl component list --count
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := fmt.Sprintf("%s/%s?%s", baseUrl, componentApiPath, utils.GetPageQueryParameter(cmd, pageNumberFlag))

		return listComponents(cmd, url)
	},
}

func init() {
	// Add the new components commands to the ComponentsCmd
	listComponentCmd.Flags().IntP("page", "p", 1, "(optional) List next set of components with --page (default = 1)")
	listComponentCmd.Flags().BoolP("count", "c", false, "(optional) Display count only")
}
