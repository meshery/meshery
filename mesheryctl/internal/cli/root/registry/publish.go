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

package system

import (
	// "fmt"

	// "github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	// "github.com/layer5io/meshery/mesheryctl/pkg/utils"
	// "github.com/pkg/errors"

	// log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

var (
	silentFlagSet bool
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish Meshery Models to Websites, Remote Provider, Meshery",
	Long:  `Publishes metadata about Meshery Models to Websites, Remote Provider, Meshery by reading from a Google Spreadsheet.`,
	Example: `
// Publish To Meshery
mesheryctl registry publish --system=meshery --credential=$GoogleCredential --sheetId=$GoogleSheetID --output=../../server/meshmodel

// Publish To Remote Provider
mesheryctl registry publish --system=remote-provider --credential=$GoogleCredential --sheetId=$GoogleSheetID --modeloutput=<remote-provider>/meshmodels/models --imgoutput=<remote-provider>/ui/public/img/meshmodels

// Publish To Website
mesheryctl registry publish --system=website --credential=$GoogleCredential --sheetId=$GoogleSheetID --output-path= --img-output=<remote-provider>/ui/public/img/meshmodels
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// main
	},
}