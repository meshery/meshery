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

package model

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var deleteModelCmd = &cobra.Command{
	Use:   "delete [model-id]",
	Short: "Delete a model",
	Long: `Delete a model by ID
Find more information at https://docs.meshery.io/reference/mesheryctl/model/delete`,
	Example: `
// Delete a model by ID
mesheryctl model delete [model-id]
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			const errMsg = "[ model-id ] is required\n\nUsage: mesheryctl model delete [model-id]\nRun 'mesheryctl model delete --help' to see detailed help message"
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}

		// if !utils.IsUUID(args[0]) {
		// 	return utils.ErrInvalidUUID(fmt.Errorf("invalid model ID: %q", args[0]))
		// }

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		modelArg := args[0]

		// Delete model by ID
		if utils.IsUUID(modelArg) {
			_, err := api.Delete(fmt.Sprintf("%s/%s", modelsApiPath, modelArg))
			if err != nil {
				return ErrDeleteModel(err, modelArg)
			}
			utils.Log.Infof("Model with ID %s has been deleted", modelArg)
			return nil
		}

		// Delete model by name

		return nil
	},
}
