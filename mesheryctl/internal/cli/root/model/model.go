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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Maximum number of rows to be displayed in a page
	maxRowsPerPage = 25
	// Color for the whiteboard printer
	whiteBoardPrinter = color.New(color.FgHiBlack, color.BgWhite, color.Bold)
	// Available model subcommads
	availableSubcommands = []*cobra.Command{listModelCmd, viewModelCmd, searchModelCmd, importModelCmd, exportModelCmd, generateModelCmd}
)

// ModelCmd represents the mesheryctl model command
var ModelCmd = &cobra.Command{
	Use:   "model",
	Short: "Manage models",
	Long:  "Export, generate, import, list, search and view model(s) and detailed informations",
	Example: `
// Display number of available models in Meshery
mesheryctl model --count

// Export registred models
mesheeryctl model export [model-name]

// Generate model(s)
mesheeryctl model export [model-name]

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
			models, err := fetchModels(url)

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

func fetchModels(url string) (*models.MeshmodelsAPIResponse, error) {
	req, err := utils.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return nil, err
	}

	// defers the closing of the response body after its use, ensuring that the resources are properly released.
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	modelsResponse := &models.MeshmodelsAPIResponse{}
	err = json.Unmarshal(data, modelsResponse)
	if err != nil {
		return nil, err
	}

	return modelsResponse, nil
}
