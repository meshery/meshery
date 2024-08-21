// Copyright 2024 Meshery Authors
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

package relationships

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	searchModelName string
	searchType      string
	searchSubType   string
	searchKind      string
)

// represents the mesheryctl exp relationship search [query-text] subcommand.
var SearchComponentsCmd = &cobra.Command{
	Use:   "search",
	Short: "Searches registered relationships",
	Long:  "Searches and finds the realtionship used by different models based on the query-text.",
	Example: `
// Search for relationship using a query
mesheryctl exp relationship search --[flag] [query-text]`,
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp relationship search --[flag] [query-text] \nRun 'mesheryctl exp relationship search --help' to see detailed help message"
		if len(args) == 1 {
			return fmt.Errorf("flag is missing. Please provide flag \n\n%v", errMsg)
		}
		kind, _ := cmd.Flags().GetString("kind")
		subType, _ := cmd.Flags().GetString("subtype")
		relType, _ := cmd.Flags().GetString("type")
		modelname, _ := cmd.Flags().GetString("model")
		if kind == "" && subType == "" && relType == "" && modelname == "" {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return utils.ErrInvalidArgument(errors.New("please provide a --kind, --subtype or --type or --model"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := ""
		kind, _ := cmd.Flags().GetString("kind")
		subType, _ := cmd.Flags().GetString("subtype")
		relType, _ := cmd.Flags().GetString("type")
		modelname, _ := cmd.Flags().GetString("model")
		if modelname == "" {
			url = fmt.Sprintf("%s/api/meshmodels/relationships?type=%s&kind=%s&subType=%s&pagesize=all", baseUrl, relType, kind, subType)

		} else {
			url = fmt.Sprintf("%s/api/meshmodels/models/%s/relationships?type=%s&kind=%s&subType=%s&pagesize=all", baseUrl, modelname, relType, kind, subType)

		}
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		relationshipResponse := &MeshmodelRelationshipsAPIResponse{}
		err = json.Unmarshal(data, relationshipResponse)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		header := []string{"kind", "apiVersion", "model-name", "subType", "regoQuery"}
		rows := [][]string{}

		for _, relationship := range relationshipResponse.Relationships {
			if len(relationship.Type()) > 0 {
				evaluationQuery := ""
				if relationship.EvaluationQuery != nil {
					evaluationQuery = *relationship.EvaluationQuery
				}
				rows = append(rows, []string{string(relationship.Kind), relationship.SchemaVersion, relationship.Model.DisplayName, relationship.SubType, evaluationQuery})
			}
		}

		if len(rows) == 0 {
			// if no relationship is found
			fmt.Println("No relationship(s) found with the search term")
			return nil
		} else {
			// Print the result in tabular format
			utils.PrintToTable(header, rows)
		}
		startIndex := 0
		endIndex := min(len(rows), startIndex+maxRowsPerPage)
		for {
			// Clear the entire terminal screen
			utils.ClearLine()

			// Print number of relationships and current page number
			fmt.Print("Total number of relationships: ", len(rows))
			fmt.Println()
			fmt.Print("Page: ", startIndex/maxRowsPerPage+1)
			fmt.Println()

			fmt.Println("Press Enter or ↓ to continue. Press Esc or Ctrl+C to exit.")

			utils.PrintToTable(header, rows[startIndex:endIndex])
			keysEvents, err := keyboard.GetKeys(10)
			if err != nil {
				return err
			}

			defer func() {
				_ = keyboard.Close()
			}()

			event := <-keysEvents
			if event.Err != nil {
				utils.Log.Error(fmt.Errorf("unable to capture keyboard events"))
				break
			}

			if event.Key == keyboard.KeyEsc || event.Key == keyboard.KeyCtrlC {
				break
			}

			if event.Key == keyboard.KeyEnter || event.Key == keyboard.KeyArrowDown {
				startIndex += maxRowsPerPage
				endIndex = min(len(rows), startIndex+maxRowsPerPage)
			}

			if startIndex >= len(rows) {
				break
			}
		}
		return nil
	},
}

func init() {
	SearchComponentsCmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})
	SearchComponentsCmd.Flags().StringVarP(&searchKind, "kind", "k", "", "search particular kind of relationships")
	SearchComponentsCmd.Flags().StringVarP(&searchSubType, "subtype", "s", "", "search particular subtype of relationships")
	SearchComponentsCmd.Flags().StringVarP(&searchModelName, "model", "m", "", "search relationships of particular model name")
	SearchComponentsCmd.Flags().StringVarP(&searchType, "type", "t", "", "search particular type of relationships")
}
