package model

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	availableSubcommands = []*cobra.Command{listModelCmd, viewModelCmd, searchModelCmd, exportModelCmd}
	// flag used to specify the page number in list command
	pageNumberFlag int
	// flag used to specify format of output of view {model-name} command
	outFormatFlag string

	// Maximum number of rows to be displayed in a page
	maxRowsPerPage = 25

	// Color for the whiteboard printer
	whiteBoardPrinter = color.New(color.FgHiBlack, color.BgWhite, color.Bold)
)

// ModelCmd represents the mesheryctl exp model command
var ModelCmd = &cobra.Command{
	Use:   "model",
	Short: "View list of models and detail of models",
	Long:  "View list of models and detailed information of a specific model",
	Example: `
// To view list of components
mesheryctl exp model list

// To view a specific model
mesheryctl exp model view [model-name]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemModelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl exp model --help' to display usage guide.\n", args[0]), "model"))
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
	listModelCmd.Flags().IntVarP(&pageNumberFlag, "page", "p", 1, "(optional) List next set of models with --page (default = 1)")
	listModelCmd.Flags().BoolP("count", "c", false, "(optional) Get the number of models in total")
	viewModelCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	wd,err := os.Getwd()
	if err != nil { 
        fmt.Println(err) 
    }
	exportModelCmd.Flags().StringVarP(&outputPath, "output", "o", wd, "location for exported oci file")
	ModelCmd.AddCommand(availableSubcommands...)
}

// selectModelPrompt lets user to select a model if models are more than one
func selectModelPrompt(models []v1alpha1.Model) v1alpha1.Model {
	modelArray := []v1alpha1.Model{}
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

func outputJson(model v1alpha1.Model) error {
	if err := prettifyJson(model); err != nil {
		// if prettifyJson return error, marshal output in conventional way using json.MarshalIndent
		// but it doesn't convert unicode to its corresponding HTML string (it is default behavior)
		// e.g unicode representation of '&' will be printed as '\u0026'
		if output, err := json.MarshalIndent(model, "", "  "); err != nil {
			return errors.Wrap(err, "failed to format output in JSON")
		} else {
			fmt.Print(string(output))
		}
	}
	return nil
}

// prettifyJson takes a v1alpha1.Model struct as input, marshals it into a nicely formatted JSON representation,
// and prints it to standard output with proper indentation and without escaping HTML entities.
func prettifyJson(model v1alpha1.Model) error {
	// Create a new JSON encoder that writes to the standard output (os.Stdout).
	enc := json.NewEncoder(os.Stdout)
	// Configure the JSON encoder settings.
	// SetEscapeHTML(false) prevents special characters like '<', '>', and '&' from being escaped to their HTML entities.
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	// Any errors during the encoding process will be returned as an error.
	return enc.Encode(model)
}
