package model

import (
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/spf13/cobra"
)

var modelViewOutputFormat string

var viewModelCmd = &cobra.Command{
	Use:   "view",
	Short: "View model",
	Long: `View a model queried by its name or id
Find more information at: https://docs.meshery.io/reference/mesheryctl/model/view`,
	Example: `
// View a specific model from current provider by using [model-name] or [model-id] in default format yaml
mesheryctl model view [model-name]


`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !slices.Contains(getValidOutputFormat(), strings.ToLower(modelViewOutputFormat)) {
			return ErrModelUnsupportedOutputFormat(formaterrMsg)
		}
		return nil
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("%s%s", errNoArg, viewUsageMsg))
		} else if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("%s%s", errMultiArg, viewUsageMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err, urlPath := getModelViewUrlPath(args[0])
		if err != nil {
			return err
		}

		modelDefinition := args[0]

		modelsResponse, err := api.Fetch[models.MeshmodelsAPIResponse](urlPath)

		if err != nil {
			return err
		}

		var selectedModel model.ModelDefinition

		switch modelsResponse.Count {
		case 0:
			fmt.Println("No model(s) found for the given name ", modelDefinition)
			return nil
		case 1:
			selectedModel = modelsResponse.Models[0]
		default:
			selectedModel = selectModelPrompt(modelsResponse.Models)
		}

		outputFormatterFactory := display.OutputFormatterFactory[model.ModelDefinition]{}
		outputFormatter, err := outputFormatterFactory.New(strings.ToLower(modelViewOutputFormat), selectedModel)
		if err != nil {
			return err
		}

		err = outputFormatter.Display()
		if err != nil {
			return err
		}

		return nil
	},
}

func getValidOutputFormat() []string {
	return []string{"yaml", "json"}
}

func selectModelPrompt(models []model.ModelDefinition) model.ModelDefinition {
	modelArray := []model.ModelDefinition{}
	modelNames := []string{}

	modelArray = append(modelArray, models...)

	for _, model := range modelArray {
		modelName := fmt.Sprintf("%s, version: %s", model.DisplayName, model.Version)
		modelNames = append(modelNames, modelName)
	}

	prompt := promptui.Select{
		Label: "Select a model",
		Items: modelNames,
	}

	for {
		i, _, err := prompt.Run()
		if err != nil {
			continue
		}

		return modelArray[i]
	}
}

func getModelViewUrlPath(modelname string) (error, string) {
	queryParams := url.Values{}
	var modelsUrlPath string

	isID, err := regexp.MatchString("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$", modelname)
	if err != nil {
		return err, ""
	}

	if !isID {
		modelsUrlPath = fmt.Sprintf("%s/%s", modelsApiPath, url.PathEscape(modelname))
	} else {
		modelsUrlPath = modelsApiPath
		queryParams.Add("id", url.QueryEscape(modelname))
	}

	queryParams.Add("pagesize", "all")

	return nil, fmt.Sprintf("%s?%s", modelsUrlPath, queryParams.Encode())
}

func init() {
	viewModelCmd.Flags().StringVarP(&modelViewOutputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
