package model

import (
	"fmt"
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
	Long: `View a model queried by its name
Documentation for models view can be found at https://docs.meshery.io/reference/mesheryctl/model/view`,
	Example: `
// View a specific model from current provider
mesheryctl model view [model-name]
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !slices.Contains(getValidOutputFormat(), strings.ToLower(modelViewOutputFormat)) {
			const errMsg = "[ yaml, json ] are the only format supported\n\nUsage: mesheryctl model view --output-format [yaml|json]\nRun 'mesheryctl model view --help' to see detailed help message"
			return ErrModelUnsupportedOutputFormat(errMsg)
		}
		return nil
	},
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl model view [model-name]\nRun 'mesheryctl model view --help' to see detailed help message"
		if len(args) == 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("model name isn't specified\n\n%v", errMsg))
		} else if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments\n\n%v", errMsg))
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

func init() {
	viewModelCmd.Flags().StringVarP(&modelViewOutputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
