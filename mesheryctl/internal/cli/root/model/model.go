// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by a, filepath.Dir(${1:}modelDefPathpplicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	modelsApiPath = "api/meshmodels/models"
	// Available model subcommads
	availableSubcommands = []*cobra.Command{listModelCmd, viewModelCmd, searchModelCmd, importModelCmd, exportModelCmd, generateModelCmd}
)

// ModelCmd represents the mesheryctl model command
var ModelCmd = &cobra.Command{
	Use:   "model",
	Short: "Manage models",
	Long: `Export, generate, import, list, search and view model(s) and detailed informations
Documentation for models can be found at https://docs.meshery.io/reference/mesheryctl/model`,
	Example: `
// Display number of available models in Meshery
mesheryctl model --count

// Export registred models
mesheryctl model export [model-name]

// Generate model(s)
mesheryctl model export [model-name]

// Import model(s)
mesheryctl model import -f [Uri]

// List available model(s)
mesheryctl model list

// Search for a specific model
mesheryctl model search [model-name]

// View a specific model
mesheryctl model view [model-name]
`,
	Args: func(cmd *cobra.Command, args []string) error {
		count, _ := cmd.Flags().GetBool("count")
		if len(args) == 0 && !count {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return utils.ErrInvalidArgument(errors.New("please provide a subcommand"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		countFlag, _ := cmd.Flags().GetBool("count")
		if countFlag {
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				log.Fatalln(err, "error processing config")
			}

			baseUrl := mctlCfg.GetBaseMesheryURL()
			url := fmt.Sprintf("%s/api/meshmodels/models?page=1", baseUrl)
			models, err := api.Fetch[models.MeshmodelsAPIResponse](url)

			if err != nil {
				return err
			}

			utils.DisplayCount("models", models.Count)

			return nil
		}

		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemModelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl model --help' to display usage guide.\n", args[0]), "model"))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		err = cmd.Usage()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	ModelCmd.AddCommand(availableSubcommands...)
	ModelCmd.Flags().BoolP("count", "", false, "(optional) Get the number of models in total")
}

func displayModels(modelsResponse *models.MeshmodelsAPIResponse, cmd *cobra.Command) error {
	header := []string{"Model", "Category", "Version"}
	rows := [][]string{}

	for _, model := range modelsResponse.Models {
		if len(model.DisplayName) > 0 {
			rows = append(rows, []string{model.Name, string(model.Category.Name), model.Version})
		}
	}

	count, _ := cmd.Flags().GetBool("count")

	dataToDisplay := display.DisplayedData{
		DataType:         "model",
		Header:           header,
		Rows:             rows,
		Count:            int64(modelsResponse.Count),
		DisplayCountOnly: count,
		IsPage:           cmd.Flags().Changed("page"),
	}

	err := display.List(dataToDisplay)
	if err != nil {
		return err
	}

	return nil
}
