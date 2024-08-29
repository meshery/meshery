package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var viewModelCmd = &cobra.Command{
	Use:   "view",
	Short: "view model",
	Long:  "view a model queried by its name",
	Example: `
// View current provider
mesheryctl model view [model-name]
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl model view [model-name]\nRun 'mesheryctl model view --help' to see detailed help message"
		if len(args) == 0 {
			return fmt.Errorf("model name isn't specified\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		modelDefinition := args[0]

		url := fmt.Sprintf("%s/api/meshmodels/models/%s?pagesize=all", baseUrl, modelDefinition)
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

		modelsResponse := &models.MeshmodelsAPIResponse{}
		err = json.Unmarshal(data, modelsResponse)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		var selectedModel model.ModelDefinition

		if modelsResponse.Count == 0 {
			fmt.Println("No model(s) found for the given name ", modelDefinition)
			return nil
		} else if modelsResponse.Count == 1 {
			selectedModel = modelsResponse.Models[0]
		} else {
			selectedModel = selectModelPrompt(modelsResponse.Models)
		}

		var output []byte

		// user may pass flag in lower or upper case but we have to keep it lower
		// in order to make it consistent while checking output format
		outFormatFlag = strings.ToLower(outFormatFlag)
		if outFormatFlag == "yaml" {
			if output, err = yaml.Marshal(selectedModel); err != nil {
				return errors.Wrap(err, "failed to format output in YAML")
			}
			fmt.Print(string(output))
		} else if outFormatFlag == "json" {
			return outputJson(selectedModel)
		} else {
			return errors.New("output-format choice invalid, use [json|yaml]")
		}

		return nil
	},
}
