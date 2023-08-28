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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var linkDocMeshDeploy = map[string]string{
	"link":    "![mesh-deploy-usage](/assets/img/mesheryctl/deploy-mesh.png)",
	"caption": "Usage of mesheryctl mesh deploy",
}

var (
	meshName  string
	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a service mesh to the Kubernetes cluster",
		Long:  `Deploy a service mesh to the connected Kubernetes cluster`,
		Example: `
// Deploy a service mesh from an interactive on the default namespace
mesheryctl mesh deploy

// Deploy a service mesh
mesheryctl mesh deploy linkerd

// Deploy Linkerd mesh on a specific namespace
mesheryctl mesh deploy linkerd --namespace linkerd-ns

// Deploy Linkerd mesh and wait for it to be deployed
mesheryctl mesh deploy linkerd --watch
		`,
		Annotations: linkDocMeshDeploy,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			log.Infof("Verifying prerequisites...")
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return utils.ErrLoadConfig(err)
			}

			if err = validateAdapter(mctlCfg, meshName); err != nil {
				// ErrValidatingAdapter
				return ErrValidatingAdapters(err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return utils.ErrLoadConfig(err)
			}
			s := utils.CreateDefaultSpinner(fmt.Sprintf("Deploying %s", meshName), fmt.Sprintf("\n%s service mesh deployed successfully", meshName))
			s.Start()
			_, err = sendOperationRequest(mctlCfg, strings.ToLower(meshName), false, "null")
			if err != nil {
				return ErrSendOperation(err)
			}
			s.Stop()

			if watch {
				log.Infof("Verifying Operation")
				_, err = waitForDeployResponse(mctlCfg, "mesh is now installed")
				if err != nil {
					return ErrWaitValidateResponse(err)
				}
			}

			return nil
		},
	}
)

func init() {
	deployCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace to be used for deploying the validation tests and sample workload")
	deployCmd.Flags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token for authenticating to Meshery API")
	deployCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch for events and verify operation (in beta testing)")
}
