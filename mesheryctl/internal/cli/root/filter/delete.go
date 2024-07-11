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

package filter

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [filter-name | ID]",
	Short: "Delete a filter file",
	Long:  `Delete a filter file using the name or ID of a filter.  Multi-word filter names must be enclosed in quotes.`,
	Example: `
// Delete the specified WASM filter file using name or ID
// A unique prefix of the name or ID can also be provided. If the prefix is not unique, the first match will be deleted.
mesheryctl filter delete [filter-name | ID]
// For multi-word filter names, use quotes:
mesheryctl filter delete "My Filter Name"
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		if len(args) == 0 {
			return errors.New(utils.FilterDeleteError("filter name or ID not provided\nUse 'mesheryctl filter delete --help' to display usage guide\n"))
		}

		var filterNameOrID string
		if len(args) == 1 {
			// Single argument case (single word or quoted multi-word)
			filterNameOrID = args[0]
		} else {
			// Multiple arguments case (unquoted multi-word)
			return errors.New(utils.FilterDeleteError("multi-word filter names must be enclosed in quotes. Example: mesheryctl filter delete \"My Filter Name\""))
		}

		// Remove surrounding quotes if present
		filterNameOrID = strings.Trim(filterNameOrID, "\"")

		var filterID, filterName string
		var isValidID, isValidName bool

		filterID, isValidID, err = utils.ValidId(mctlCfg.GetBaseMesheryURL(), filterNameOrID, "filter")
		if err != nil {
			utils.Log.Error(ErrFilterNameOrID(err))
			return nil
		}

		if !isValidID {
			filterName, filterID, isValidName, err = utils.ValidName(mctlCfg.GetBaseMesheryURL(), filterNameOrID, "filter")
			if err != nil {
				utils.Log.Error(ErrFilterNameOrID(err))
				return nil
			}
		}

		// Delete the filter using the id
		if isValidID || isValidName {
			err := utils.DeleteConfiguration(mctlCfg.GetBaseMesheryURL(), filterID, "filter")

			var filter string
			if isValidID {
				filter = filterID
			} else {
				filter = filterName
			}
			if err != nil {
				return errors.Wrap(err, utils.FilterDeleteError(fmt.Sprintf("failed to delete filter %s", filter)))
			}
			utils.Log.Info("Filter ", filter, " deleted successfully")
			return nil
		}

		return errors.New(utils.FilterDeleteError(fmt.Sprintf("filter with name or ID having prefix %s does not exist", filterNameOrID)))
	},
}
