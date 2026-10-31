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
	Use:   "delete [model-id | model-name]",
	Short: "Delete a model",
	Long: `Delete a model by ID or Name
Find more information at https://docs.meshery.io/reference/mesheryctl/model/delete`,
	Example: `
// Delete a model by ID
mesheryctl model delete [model-id]

// Delete a model by name
mesheryctl model delete [model-name]
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return utils.ErrInvalidArgument(errors.New(errDeleteInvalidArg))
		}
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

		// Delete model by name, for multiple matches use pagination selection prompt
		selectedModel, err := promptModelSelection(modelArg, modelsApiPath)
		if err != nil {
			return err
		}

		if selectedModel == nil {
			utils.Log.Infof("No model(s) found with the name: %s", modelArg)
			return nil
		}

		// Delete the selected model by its UUID
		_, err = api.Delete(fmt.Sprintf("%s/%s", modelsApiPath, selectedModel.ID.String()))
		if err != nil {
			return ErrDeleteModel(err, modelArg)
		}
		utils.Log.Infof("Model '%s' (ID: %s) has been deleted", selectedModel.DisplayName, selectedModel.ID.String())

		return nil
	},
}
