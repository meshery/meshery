package system

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// represents the `mesheryctl exp components list` command
var listComponentCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered components",
	Long:  `List all components registered in Meshery Server`,
	Example: `
	// View list of components
mesheryctl exp components list

// View list of components with specified page number (25 components per page)
mesheryctl exp components list --page 2
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check prerequisites for the command here

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
		if len(args) != 0 {
			return errors.New(utils.SystemModelSubError("this command takes no arguments\n", "list"))
		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		var url string
		if cmd.Flags().Changed("page") {
			url = fmt.Sprintf("%s/api/meshmodels/components?page=%d", baseUrl, pageNumberFlag)
		} else {
			url = fmt.Sprintf("%s/api/meshmodels/components?pagesize=all", baseUrl)
		}
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		componentsResponse := &models.MeshmodelComponentsAPIResponse{}
		err = json.Unmarshal(data, componentsResponse)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		header := []string{"TypeMeta", "Model", "Schema"}
		rows := [][]string{}

		for _, component := range componentsResponse.Components {
			if len(component.DisplayName) > 0 {
				rows = append(rows, []string{component.TypeMeta.Kind, component.TypeMeta.APIVersion, component.Model.Name, component.Schema})
			}
		}

		if len(rows) == 0 {
			// if no component is found
			utils.Log.Info("No components(s) found")
		} else {
			utils.PrintToTable(header, rows)
		}

		return nil
	},
}

// represents the `mesheryctl exp components view [component-name]` subcommand.
var viewComponentCmd = &cobra.Command{
	Use:   "view",
	Short: "view registered components",
	Long:  "view a component registered in Meshery Server",
	Example: `
// View details of a specific component
mesheryctl exp components view [component-name]
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
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp component view [component-name]\nRun 'mesheryctl exp component view --help' to see detailed help message"
		if len(args) == 0 {
			return fmt.Errorf("component name isn't specified\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		component := args[0]

		url := fmt.Sprintf("%s/api/meshmodels/components/%s?pagesize=all", baseUrl, component)
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		componentResponse := &models.MeshmodelComponentsAPIResponse{}
		err = json.Unmarshal(data, componentResponse)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		var selectedComponent v1alpha1.ComponentDefinition

		if componentResponse.Count == 0 {
			utils.Log.Info("No component(s) found for the given name ", component)
			return nil
		} else if componentResponse.Count == 1 {
			selectedComponent = componentResponse.Components[0] // Update the type of selectedModel
		} else {
			selectedComponent = selectComponentPrompt(componentResponse.Components)
		}

		var output []byte

		// user may pass flag in lower or upper case but we have to keep it lower
		// in order to make it consistent while checking output format
		outFormatFlag = strings.ToLower(outFormatFlag)
		if outFormatFlag == "yaml" {
			if output, err = yaml.Marshal(selectedComponent); err != nil {
				return errors.Wrap(err, "failed to format output in YAML")
			}
			fmt.Print(string(output))
		} else if outFormatFlag == "json" {
			return outputComponentJson(selectedComponent)
		} else {
			return errors.New("output-format choice invalid, use [json|yaml]")
		}

		return nil
	},
}

// represents the `mesheryctl exp components search [query-text]` subcommand.
var searchComponentsCmd = &cobra.Command{
	Use:   "search",
	Short: "search registered components",
	Long:  "search components registered in Meshery Server",
	Example: `
// Search for components using a query
mesheryctl exp components search [query-text]
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
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp component search [query-text]\nRun 'mesheryctl exp component search --help' to see detailed help message"
		if len(args) == 0 {
			return fmt.Errorf("Search term is missing\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		queryText := args[0]

		url := fmt.Sprintf("%s/api/meshmodels/components?search=%s&pagesize=all", baseUrl, queryText)
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		componentsResponse := &models.MeshmodelComponentsAPIResponse{}
		err = json.Unmarshal(data, componentsResponse)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		header := []string{"TypeMeta", "Model", "Schema"}
		rows := [][]string{}

		for _, component := range componentsResponse.Components {
			if len(component.DisplayName) > 0 {
				rows = append(rows, []string{component.TypeMeta.Kind, component.TypeMeta.APIVersion, component.Model.Name, component.Schema})
			}
		}

		if len(rows) == 0 {
			// if no model is found
			utils.Log.Info("No components(s) found")
		} else {
			utils.PrintToTable(header, rows)
		}

		return nil
	},
}

// ComponentsCmd represents the `mesheryctl exp components` command
var ComponentsCmd = &cobra.Command{
	Use:   "components",
	Short: "View list of components and detail of components",
	Long:  `View list of components and detailed information of a specific component`,
	Example: `
// To view list of components
mesheryctl exp components list

// To view a specific component
mesheryctl exp components view [component-name]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			_ = cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemModelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl exp component --help' to display usage guide.\n", args[0]), "component"))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}
		err = viewProviderCmd.RunE(cmd, args)
		if err != nil {
			return err
		}
		err = cmd.Usage()
		if err != nil {
			return err
		}
		return nil
	},
}

// `selectComponentPrompt` lets user to select a model if models are more than one

func selectComponentPrompt(components []v1alpha1.ComponentDefinition) v1alpha1.ComponentDefinition {
	componentNames := []string{}
	componentArray := []v1alpha1.ComponentDefinition{}

	componentArray = append(componentArray, components...)

	for _, component := range componentArray {
		componentName := fmt.Sprintf("%s, version: %s", component.DisplayName, component.TypeMeta.APIVersion)
		componentNames = append(componentNames, componentName)
	}

	prompt := promptui.Select{
		Label: "Select component",
		Items: componentNames,
	}

	for {
		i, _, err := prompt.Run()
		if err != nil {
			continue
		}

		return componentArray[i]
	}
}

func outputComponentJson(component v1alpha1.ComponentDefinition) error {
	if err := prettifyComponentJson(component); err != nil {
		// if prettifyJson return error, marshal output in conventional way using json.MarshalIndent
		// but it doesn't convert unicode to its corresponding HTML string (it is default behavior)
		// e.g unicode representation of '&' will be printed as '\u0026'
		if output, err := json.MarshalIndent(component, "", "  "); err != nil {
			return errors.Wrap(err, "failed to format output in JSON")
		} else {
			fmt.Print(string(output))
		}
	}
	return nil
}


// prettifyJson takes a v1alpha1.Model struct as input, marshals it into a nicely formatted JSON representation,
// and prints it to standard output with proper indentation and without escaping HTML entities.
func prettifyComponentJson(component v1alpha1.ComponentDefinition) error {
	// Create a new JSON encoder that writes to the standard output (os.Stdout).
	enc := json.NewEncoder(os.Stdout)
	// Configure the JSON encoder settings.
	// SetEscapeHTML(false) prevents special characters like '<', '>', and '&' from being escaped to their HTML entities.
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	// Any errors during the encoding process will be returned as an error.
	return enc.Encode(component)
}

func init() {
	// Add the new exp components commands to the ComponentsCmd
	availableSubcommands := []*cobra.Command{listComponentCmd, viewComponentCmd, searchComponentsCmd}
	ComponentsCmd.AddCommand(availableSubcommands...)
}
