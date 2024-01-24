// # Copyright Meshery Authors
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

package registry

import (
	// "fmt"

	// "github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	// "github.com/layer5io/meshery/mesheryctl/pkg/utils"
	// "github.com/pkg/errors"

	// log "github.com/sirupsen/logrus"

	"fmt"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

// publishCmd represents the publish command to publish Meshery Models to Websites, Remote Provider, Meshery
var publishCmd = &cobra.Command{
	Use:   "publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path]",
	Short: "Publish Meshery Models to Websites, Remote Provider, Meshery",
	Long:  `Publishes metadata about Meshery Models to Websites, Remote Provider and Meshery by reading from a Google Spreadsheet.`,
	Example: `
// Publish To System
mesheryctl publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path]

// Publish To Meshery
mesheryctl registry publish meshery GoogleCredential GoogleSheetID <repo>/server/meshmodel

// Publish To Remote Provider
mesheryctl registry publish --system=remote-provider GoogleCredential GoogleSheetID <repo>/meshmodels/models <repo>/ui/public/img/meshmodels

// Publish To Website
mesheryctl registry publish --system=website GoogleCredential $GoogleSheetID <repo>/integrations <repo>/ui/public/img/meshmodels
	`,
	// PreRunE: func(cmd *cobra.Command, args []string) error {
	// 	//Check prerequisite
	// },
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		fmt.Println("args:", args)
		return nil
	},
}

func init() {
	publishCmd.Flags().StringVarP(&system, "system", "", "", "System to publish to. Available options: meshery, remote-provider, website")
	publishCmd.Flags().StringVarP(&credential, "credential", "", "", "Google Credential File")
	publishCmd.Flags().StringVarP(&sheetId, "sheetId", "", "", "Google Sheet ID to read from")
	publishCmd.Flags().StringVarP(&modelsOutput, "models-output", "", "", "Path to Meshery Models output")
	publishCmd.Flags().StringVarP(&imgsOutput, "imgs-output", "", "", "Path to Meshery Models Images output")
}
