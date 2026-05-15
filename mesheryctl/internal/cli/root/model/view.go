package model

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/spf13/cobra"
)

type cmdModelViewFlags struct {
	OutputFormat string `json:"output-format" validate:"oneof=json yaml"`
	Save         bool   `json:"save" validate:"boolean"`
}

var modelViewFlags cmdModelViewFlags

var viewModelCmd = &cobra.Command{
	Use:   "view",
	Short: "View model",
	Long: `View a model queried by its name or ID
Find more information at: https://docs.meshery.io/reference/mesheryctl/model/view`,
	Example: `
// View a specific model from current provider by using [model-name] or [model-id] in default format yaml
mesheryctl model view [model-name]

// View a specific model in specifed format
mesheryctl model view [model-name] --output-format [json|yaml]

// View a specific model in specified format and save it as a file
mesheryctl model view [model-name] --output-format [json|yaml] --save
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &modelViewFlags)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return utils.ErrInvalidArgument(errors.New(errInvalidArg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		modelNameOrId := args[0]

		var selectedModel *model.ModelDefinition
		var err error

		if !utils.IsUUID(modelNameOrId) {
			selectedModel, err = promptModelSelection(modelNameOrId, modelsApiPath)
		} else {
			queryParams := url.Values{}
			queryParams.Add("id", modelNameOrId)
			modelViewApiPath := fmt.Sprintf("%s?%s", modelsApiPath, queryParams.Encode())
			selectedModel, err = promptModelSelection("", modelViewApiPath)
		}

		if err != nil {
			return err
		}

		if selectedModel == nil {
			utils.Log.Infof("No model found with name or ID: %s", modelNameOrId)
			return nil
		}

		outputFormatterFactory := display.OutputFormatterFactory[*model.ModelDefinition]{}
		outputFormatter, err := outputFormatterFactory.New(strings.ToLower(modelViewFlags.OutputFormat), selectedModel)
		if err != nil {
			return err
		}

		err = outputFormatter.WithOutput(cmd.OutOrStdout()).Display()
		if err != nil {
			return err
		}

		if modelViewFlags.Save {
			outputFormatterSaverFactory := display.OutputFormatterSaverFactory[*model.ModelDefinition]{}
			outputFormatterSaver, err := outputFormatterSaverFactory.New(modelViewFlags.OutputFormat, outputFormatter)
			if err != nil {
				return err
			}

			modelString := strings.ReplaceAll(fmt.Sprintf("%v", selectedModel.DisplayName), " ", "_")
			fileName := filepath.Join(utils.MesheryFolder, fmt.Sprintf("model_%s.%s", modelString, modelViewFlags.OutputFormat))

			outputFormatterSaver = outputFormatterSaver.WithFilePath(fileName)
			err = outputFormatterSaver.Save()
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func getModelViewUrlPath(modelNameOrId string) string {
	queryParams := url.Values{}

	if !utils.IsUUID(modelNameOrId) {
		queryParams.Add("search", modelNameOrId)
	} else {
		queryParams.Add("id", modelNameOrId)
	}

	return fmt.Sprintf("%s?%s", modelsApiPath, queryParams.Encode())
}

func init() {
	viewModelCmd.Flags().StringVarP(&modelViewFlags.OutputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewModelCmd.Flags().BoolVarP(&modelViewFlags.Save, "save", "s", false, "(optional) save output as a JSON/YAML file")
}
