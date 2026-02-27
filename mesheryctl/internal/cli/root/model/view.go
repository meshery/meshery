package model

import (
	"fmt"
	"slices"
	"strings"

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
	Long: `View a model queried by its name
Find more information at: https://docs.meshery.io/reference/mesheryctl/model/view`,
	Example: `
// View a specific model from current provider
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
		modelDefinition := args[0]

		modelsResponse, err := api.Fetch[models.MeshmodelsAPIResponse](fmt.Sprintf("%s/%s?pagesize=all", modelsApiPath, modelDefinition))

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
			selectedModel, err = selectModelPrompt(modelsResponse.Models)
			if err != nil {
				return err
			}
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

func selectModelPrompt(models []model.ModelDefinition) (model.ModelDefinition, error) {
	modelNames := make([]string, len(models))

	for i, model := range models {
		modelNames[i] = fmt.Sprintf("%s, version: %s", model.DisplayName, model.Version)
	}

	i, err := utils.RunSelectPrompt("Select a model", modelNames)
	if err != nil {
		return model.ModelDefinition{}, err
	}

	return models[i], nil
}

func init() {
	viewModelCmd.Flags().StringVarP(&modelViewOutputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
