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

// Operation is the common body type to be passed for Mesh Ops
type Operation struct {
	Adapter    string `json:"adapter"`
	CustomBody string `json:"customBody"`
	DeleteOp   string `json:"deleteOp"`
	Namespace  string `json:"namespace"`
	Query      string `json:"query"`
}

var spec string

var linkDocMeshValidate = map[string]string{
	"link":    "![mesh-validate-usage](/assets/img/mesheryctl/mesh-validate.png)",
	"caption": "Usage of mesheryctl mesh validate",
}

// validateCmd represents the service mesh validation command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate conformance to service mesh standards",
	Example: `
// Validate conformance to service mesh standards
mesheryctl mesh validate [mesh name] --adapter [name of the adapter] --tokenPath [path to token for authentication] --spec [specification to be used for conformance test] --namespace [namespace to be used]

// Validate Istio to service mesh standards
mesheryctl mesh validate istio --adapter meshery-istio --spec smi
	`,
	Annotations: linkDocMeshValidate,
	Long:        `Validate service mesh conformance to different standard specifications`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("Verifying prerequisites...")

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err)
		}

		prefs, err := utils.GetSessionData(mctlCfg)
		if err != nil {
			log.Fatalln(err)
		}
		//resolve adapterUrl to adapter Location
		for _, adapter := range prefs.MeshAdapters {
			adapterName := strings.Split(adapter.Location, ":")
			if adapterName[0] == adapterURL {
				adapterURL = adapter.Location
				meshName = adapter.Location
			}
		}
		//sync with available adapters
		if err = validateAdapter(mctlCfg, meshName); err != nil {
			log.Fatalln(err)
		}
		log.Info("verified prerequisites")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("Starting service mesh validation...")

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err)
		}
		s := utils.CreateDefaultSpinner(fmt.Sprintf("Validating %s", meshName), fmt.Sprintf("\n%s validation successful", meshName))
		s.Start()
		_, err = sendOperationRequest(mctlCfg, meshName, false, spec)
		if err != nil {
			log.Fatalln(err)
		}
		s.Stop()

		if watch {
			log.Infof("Verifying Operation")
			_, err = waitForValidateResponse(mctlCfg, "Smi conformance test")
			if err != nil {
				log.Fatalln(err)
			}
		}

		return nil
	},
}

func init() {
	validateCmd.Flags().StringVarP(&spec, "spec", "s", "smi", "(Required) specification to be used for conformance test (smi/istio-vet)")
	_ = validateCmd.MarkFlagRequired("spec")
	validateCmd.Flags().StringVarP(&adapterURL, "adapter", "a", "meshery-nsm", "(Required) Adapter to use for validation")
	_ = validateCmd.MarkFlagRequired("adapter")
	validateCmd.Flags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token for authenticating to Meshery API")
	validateCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch for events and verify operation (in beta testing)")
}
