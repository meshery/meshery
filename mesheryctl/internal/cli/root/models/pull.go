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
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/models/oci"
	"github.com/spf13/cobra"
)

var pullModelCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull model",
	Long:  "pull a model to the registry",
	Example: `
// Push a model
mesheryctl exp model push --username [username] --password [password] --registry [registry] --tag [tag] --repository [repository]
	`,
	// skip preRunE as it is not required for this command
	Args: func(cmd *cobra.Command, args []string) error {
		errMsg := "Usage: mesheryctl exp model pull --username [username] --password [password] --registry [registry] --tag [tag] --repository [repository]\nRun 'mesheryctl exp model pull --help' to see detailed help message"
		if len(args) < 5 {
			// also print 5 argument are required in the error message
			return utils.ErrInvalidArgument(fmt.Errorf("%s: expected 5 arguments (username, password, registryURL, repository, image_tag)", errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		dirPath := "./meshery/models"
		err := oci.PullFromOCIRegistry(dirPath, registry, repository, tag, username, password)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		fmt.Println("Model pulled successfully")
		return nil
	},
}
