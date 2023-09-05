// Copyright 2023 Layer5, Inc.
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

package mesh

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
		Short: "remove a service mesh in the kubernetes cluster",
		Long:  `remove service mesh in the connected kubernetes cluster`,
		Example: `
// Remove a service mesh(linkerd)
mesheryctl mesh remove linkerd

// Remove a service mesh(linkerd) under a specific namespace(linkerd-ns)
mesheryctl mesh remove linkerd --namespace linkerd-ns
		`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			utils.Log.Info("Verifying prerequisites...")
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				utils.Log.Error(utils.ErrLoadConfig(err))
				return nil
			}

			if err = validateAdapter(mctlCfg, meshName); err != nil {
				utils.Log.Error(ErrValidatingAdapters(err))
				return nil
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) > 1 {
				return errors.New(utils.MeshError("'mesheryctl mesh remove' should not have more than one argument, it can remove only one adapter at a time.\n"))
			}

			s := utils.CreateDefaultSpinner(fmt.Sprintf("Removing %s", meshName), fmt.Sprintf("\n%s service mesh removed successfully", meshName))
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				utils.Log.Error(err)
				return nil
			}

			s.Start()
			_, err = sendOperationRequest(mctlCfg, strings.ToLower(meshName), true, "null")
			if err != nil {
				utils.Log.Error(errors.Wrap(ErrSendOperation(err), "error removing service mesh"))
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
