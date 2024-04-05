// Copyright 2024 Layer5, Inc.
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

package models

import (
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/models/oci"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var pushModelCmd = &cobra.Command{
	Use:   "push",
	Short: "push model",
	Long:  "push a model to the registry",
	Example: `
// Push a model
mesheryctl exp model push [path-to-model] --username [username] --password [password] --registry [registry] --tag [tag] --repository [repository]
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Check if all flags are set
		usernameFlag, _ := cmd.Flags().GetString("username")
		passwordFlag, _ := cmd.Flags().GetString("password")
		registryFlag, _ := cmd.Flags().GetString("registry")
		tagFlag, _ := cmd.Flags().GetString("tag")
		repositoryFlag, _ := cmd.Flags().GetString("repository")

		if usernameFlag == "" || passwordFlag == "" || registryFlag == "" || tagFlag == "" || repositoryFlag == "" {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return errors.New("all flags are required")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// path to model
		modelPath := pathToModel
		err := oci.PushToOCIRegistry(modelPath, registry, repository, tag, username, password)
		if err != nil {
			return err
		}

		utils.Log.Info("Model pushed successfully")
		return nil
	},
}
