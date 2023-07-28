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
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	query    string
	applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "Apply a service operation",
		Long:  "Apply a service mesh operation for configuration and management of a service mesh",
		Example: `
//Apply a Specific mesh operation 
mesheryctl mesh apply [operation-name] service mesh

//Apply name of operation like mTLS to secure communication within a service mesh
mesheryctl mesh apply [mtls] istio`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			utils.Log.Info("Verifying prerequisites...")
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return err
			}
			meshName, err = validateMesh(mctlCfg, meshName)
			if err != nil {
				return err
			}
			err = validateAdapter(mctlCfg, meshName)
			if err != nil {
				return ErrValidateAdapter(err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return err
			}
			err = ApplyOperation(mctlCfg, query)
			if err != nil {
				return err
			}
			return nil
		},
	}
)
