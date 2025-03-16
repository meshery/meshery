package model

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var viewModelCmd = &cobra.Command{
	Use:   "view",
	Short: "View model",
	Long:  "View a model queried by its name",
	Example: `
// View a specific model from current provider
mesheryctl model view [model-name]
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("output-format")
		if !slices.Contains(getValidOutputFormat(), format) {
			const errMsg = "[ yaml, json ] are the only format supported\n\nUsage: mesheryctl model view --output-format [yaml|json]\nRun 'mesheryctl model view --help' to see detailed help message"
			return ErrModelUnsupportedOutputFormat(errMsg)
		}
		return nil
	},
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

		modelsResponse, err := fetchModels(url)

		if err != nil {
			return err
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

		// user may pass flag in lower or upper case but we have to keep it lower
		// in order to make it consistent while checking output format
		outputFormat, _ := cmd.Flags().GetString("output-format")

		return viewModel(selectedModel, strings.ToLower(outputFormat))
	},
}

func getValidOutputFormat() []string {
	return []string{"yaml", "json"}
}

func viewModel(model model.ModelDefinition, format string) error {

	if format == "yaml" {
		output, err := yaml.Marshal(model)
		if err != nil {
			return errors.Wrap(err, "failed to format output in YAML")
		}
		fmt.Print(string(output))
	}

	if format == "json" {
		return outputJson(model)
	}

	return nil
}

func outputJson(model model.ModelDefinition) error {
	err := prettifyJson(model)

	if err != nil {
		// if prettifyJson return error, marshal output in conventional way using json.MarshalIndent
		// but it doesn't convert unicode to its corresponding HTML string (it is default behavior)
		// e.g unicode representation of '&' will be printed as '\u0026'
		output, err := json.MarshalIndent(model, "", "  ")

		if err != nil {
			return errors.Wrap(err, "failed to format output in JSON")
		} else {
			fmt.Print(string(output))
		}
	}

	return nil
}

// prettifyJson takes a model.ModelDefinition struct as input, marshals it into a nicely formatted JSON representation,
// and prints it to standard output with proper indentation and without escaping HTML entities.
func prettifyJson(model model.ModelDefinition) error {
	// Create a new JSON encoder that writes to the standard output (os.Stdout).
	enc := json.NewEncoder(os.Stdout)
	// Configure the JSON encoder settings.
	// SetEscapeHTML(false) prevents special characters like '<', '>', and '&' from being escaped to their HTML entities.
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	// Any errors during the encoding process will be returned as an error.
	return enc.Encode(model)
}

// selectModelPrompt lets user to select a model if models are more than one
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
	viewModelCmd.Flags().StringP("output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
