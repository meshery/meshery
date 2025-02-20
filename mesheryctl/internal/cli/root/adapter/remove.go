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

package adapter

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "remove cloud and cloud native infrastructure",
		Long:  `remove cloud and cloud native infrastructure`,
		Example: `
// Remove Linkerd deployment
mesheryctl adapter remove linkerd

// Remove a Linkerd control plane found under a specific namespace (linkerd-ns)
mesheryctl adapter remove linkerd --namespace linkerd-ns
		`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			utils.Log.Info("Verifying prerequisites...")
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				utils.Log.Error(err)
				return nil
			}

			if err = validateAdapter(mctlCfg, meshName); err != nil {
				utils.Log.Error(err)
				return nil
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := utils.CreateDefaultSpinner(fmt.Sprintf("Removing %s", meshName), fmt.Sprintf("\n%s infrastructure removed", meshName))
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				utils.Log.Error(err)
				return nil
			}

			s.Start()
			_, err = sendOperationRequest(mctlCfg, strings.ToLower(meshName), true, "null")
			if err != nil {
				utils.Log.Error(ErrSendOperation(errors.Wrap(err, "error removing infrastructure")))
				return nil
			}
			s.Stop()

			return nil
		},
	}
)

func init() {
	removeCmd.Flags().StringVarP(
		&namespace, "namespace", "n", "default",
		"Kubernetes namespace where the mesh is deployed",
	)
}
