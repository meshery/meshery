package system

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/manifoldco/promptui"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configuration     *config.MesheryCtlConfig
	tempCntxt         = "local"
	tokenNameLocation = map[string]string{} //maps each token name to its specified location
)

type contextWithLocation struct {
	Endpoint      string   `mapstructure:"endpoint,omitempty"`
	Token         string   `mapstructure:"token,omitempty"`
	Tokenlocation string   `mapstructure:"token,omitempty" yaml:"token-location,omitempty"`
	Platform      string   `mapstructure:"platform"`
	Components    []string `mapstructure:"components,omitempty"`
	Channel       string   `mapstructure:"channel,omitempty"`
	Version       string   `mapstructure:"version,omitempty"`
	Provider      string   `mapstructure:"provider,omitempty"`
}

var linkDocContextCreate = map[string]string{
	"link":    "![context-create-usage](/reference/images/newcontext.png)",
	"caption": "Usage of mesheryctl context create",
}

type cmdContextCreateFlags struct {
	URL        string   `json:"url" validate:"omitempty,url"`
	Components []string `json:"components" validate:"omitempty"`
	Platform   string   `json:"platform" validate:"omitempty,oneof=docker kubernetes"`
	Provider   string   `json:"provider" validate:"omitempty,oneof=Layer5 None"`
	Set        bool     `json:"set" validate:"boolean"`
}

type cmdContextDeleteFlags struct {
	Set string `json:"set" validate:"omitempty"`
}

var contextCreateFlags cmdContextCreateFlags
var contextDeleteFlags cmdContextDeleteFlags

// createContextCmd represents the create command
var createContextCmd = &cobra.Command{
	Use:   "create context-name",
	Short: "Create a new context (a named Meshery deployment)",
	Long: `Add a new context to Meshery config.yaml file.
Find more information at: https://docs.meshery.io/reference/mesheryctl/system/context/create`,
	Example: `
// Create new context
mesheryctl system context create [context-name]

// Create new context and provide list of components, platform & URL and set it as current context
mesheryctl system context create [context-name] --components [meshery-nsm] --platform [docker|kubernetes] --url [server-url] --set --yes
	`,
	Annotations: linkDocContextCreate,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &contextCreateFlags)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errArgMsg, contextCreateUsageMsg))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		tempCntxt := utils.TemplateContext

		if contextCreateFlags.URL != "" {
			err := utils.ValidateURL(contextCreateFlags.URL)
			if err != nil {
				return err
			}
			tempCntxt.Endpoint = contextCreateFlags.URL
		}

		utils.Log.Debug("serverURL: `" + tempCntxt.Endpoint + "`")

		if contextCreateFlags.Platform != "" {
			tempCntxt.Platform = contextCreateFlags.Platform
		}

		if contextCreateFlags.Provider != "" {
			tempCntxt.Provider = contextCreateFlags.Provider
		}

		if len(contextCreateFlags.Components) >= 1 {
			tempCntxt.Components = contextCreateFlags.Components
		}

		contextName := strings.ToLower(args[0])

		err := config.AddContextToConfig(contextName, tempCntxt, viper.ConfigFileUsed(), contextCreateFlags.Set, false)
		if err != nil {
			return err
		}

		utils.Log.Infof("Added `%s` context", contextName)
		return nil
	},
}

// deleteContextCmd represents the delete command
var deleteContextCmd = &cobra.Command{
	Use:   "delete [context-name]",
	Short: "Delete context",
	Long: `Delete an existing context (a named Meshery deployment) from Meshery config file.
Find more information at: https://docs.meshery.io/reference/mesheryctl/system/context/delete`,
	Example: `
// ### Delete context
mesheryctl system context delete [context name]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &contextDeleteFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errArgMsg, contextDeleteUsageMsg))
		}
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return ErrUnmarshallConfig(err)
		}

		contextName := strings.ToLower(args[0])

		_, exists := configuration.Contexts[contextName]
		if !exists {
			return ErrContextNotExists(fmt.Errorf("no context name found : %s", contextName))
		}

		if viper.GetString("current-context") == contextName {
			var res bool
			if utils.SilentFlag {
				res = true
			} else {
				res = utils.AskForConfirmation("Are you sure you want to delete the current context")
			}

			if !res {
				utils.Log.Infof("Delete aborted")
				return nil
			}

			var newContext string

			if contextDeleteFlags.Set != "" {
				_, exists := configuration.Contexts[contextDeleteFlags.Set]
				if !exists {
					return ErrSetCurrentContext(fmt.Errorf("new context wrongly set"))
				}

				if contextDeleteFlags.Set == contextName {
					return ErrSetCurrentContext(fmt.Errorf("choose a new context other than the context being deleted"))
				}

				newContext = contextDeleteFlags.Set
			} else {
				var listContexts []string
				for context := range configuration.Contexts {
					if context != contextName {
						listContexts = append(listContexts, context)
					}
				}

				prompt := promptui.Select{
					Label: "Select context",
					Items: listContexts,
				}

				_, newContext, err = prompt.Run()
				if err != nil {
					return utils.ErrPromptCancelled()
				}
			}

			viper.Set("current-context", newContext)
			utils.Log.Infof("The current context is now %q", newContext)
		}
		delete(configuration.Contexts, contextName)
		viper.Set("contexts", configuration.Contexts)
		utils.Log.Infof("deleted context %s", contextName)
		err = viper.WriteConfig()
		if err != nil {
			return config.ErrWriteMeshConfig(err)
		}
		return nil
	},
}

// listContextCmd represents the list command
var listContextCmd = &cobra.Command{
	Use:   "list",
	Short: "list contexts",
	Long: `List current context and available contexts.
Find more information at: https://docs.meshery.io/reference/mesheryctl/system/context/list`,
	Example: `
// List all contexts present
mesheryctl system context list
	`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return utils.ErrInvalidArgument(errors.New(utils.SystemContextSubError("this command takes no arguments.\n", "list")))
		}
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return utils.ErrUnmarshal(err)
		}
		var currContext string
		var contexts = configuration.Contexts
		if contexts == nil {
			utils.Log.Info("No contexts available. Use `mesheryctl system context create <name>` to create a new Meshery deployment context.\n")
			return nil
		}

		if currContext == "" {
			currContext = viper.GetString("current-context")
		}
		if currContext == "" {
			utils.Log.Info("Current context not set\n")
		} else {
			utils.Log.Infof("Current context: %s", currContext)
		}
		utils.Log.Info("Available contexts:\n")

		//sorting the contexts to get a consistent order on each subsequent run
		keys := make([]string, 0, len(contexts))
		for k := range contexts {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			utils.Log.Infof("- %s", k)
		}

		if currContext == "" {
			utils.Log.Info("\nRun `mesheryctl system context switch <context name>` to set the current context.")
		}
		return nil
	},
}

var linkDocContextView = map[string]string{
	"link":    "![context-view-usage](/reference/images/context-view.png)",
	"caption": "Usage of mesheryctl context view",
}

type cmdContextViewFlags struct {
	Context string `json:"context" validate:"omitempty"`
	All     bool   `json:"all" validate:"boolean"`
}

var contextViewFlags cmdContextViewFlags

// viewContextCmd represents the view command
var viewContextCmd = &cobra.Command{
	Use:   "view [context-name | --context context-name | --all] --flags",
	Short: "Display the current Meshery CLI configuration context",
	Long: `Display the current Meshery CLI context configuration.
This command shows which Kubernetes cluster, platform, and provider Meshery is configured to communicate with.
Use this to verify or debug your current CLI settings.
Find more information at: https://docs.meshery.io/reference/mesheryctl/system/context/view`,
	Example: `
// View the default context
mesheryctl system context view

// View a specified context
mesheryctl system context view context-name

// View a specified context using the --context flag
mesheryctl system context view --context context-name

// View configuration of all contexts
mesheryctl system context view --all
    `,
	Annotations:  linkDocContextView,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &contextViewFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return ErrUnmarshallConfig(err)
		}
		//Storing all the tokens separately in a map, to get tokenlocation by token name.
		for _, tok := range configuration.Tokens {
			tokenNameLocation[tok.Name] = tok.Location
		}

		if contextViewFlags.All {
			tempcontexts := make(map[string]contextWithLocation)

			//Populating auxiliary struct with token-locations
			for k, v := range configuration.Contexts {
				if v.Token == "" {
					utils.Log.Warnf("[Warning]: Token not specified/empty for context \"%s\"", k)
					temp, _ := getContextWithTokenLocation(&v)
					tempcontexts[k] = *temp
				} else {
					temp, ok := getContextWithTokenLocation(&v)
					tempcontexts[k] = *temp
					if !ok {
						utils.Log.Warnf("[Warning]: Token \"%s\" could not be found! for context \"%s\"", tempcontexts[k].Token, k)
					}
				}

			}

			utils.Log.Info(getYAML(tempcontexts))

			return nil
		}
		if len(args) != 0 {
			contextViewFlags.Context = strings.ToLower(args[0])
		}
		if contextViewFlags.Context == "" {
			contextViewFlags.Context = viper.GetString("current-context")
		}
		if contextViewFlags.Context == "" {
			return ErrContextNotExists(fmt.Errorf("current context not set"))
		}

		contextData, ok := configuration.Contexts[contextViewFlags.Context]
		if !ok {
			return ErrContextNotExists(
				fmt.Errorf(
					"context `%s` does not exist",
					contextViewFlags.Context,
				),
			)
		}

		if contextData.Token == "" {
			utils.Log.Warnf("[Warning]: Token not specified/empty for context \"%s\"", contextViewFlags.Context)
			utils.Log.Infof("\nCurrent Context: %s\n", contextViewFlags.Context)
			utils.Log.Info(getYAML(contextData))
		} else {
			temp, ok := getContextWithTokenLocation(&contextData)
			utils.Log.Infof("\nCurrent Context: %s\n", contextViewFlags.Context)
			if !ok {
				utils.Log.Warnf("[Warning]: Token \"%s\" could not be found! for context \"%s\"", temp.Token, contextViewFlags.Context)
			}
			utils.Log.Info(getYAML(temp))
		}

		return nil
	},
}

var linkDocContextSwitch = map[string]string{
	"link":    "![context-switch-usage](/reference/images/contextswitch.png)",
	"caption": "Usage of mesheryctl context switch",
}

// switchContextCmd represents the switch command
var switchContextCmd = &cobra.Command{
	Use:   "switch context-name",
	Short: "switch context",
	Long: `Configure mesheryctl to actively use one one context vs. another context.
Find more information at: https://docs.meshery.io/reference/mesheryctl/system/context/switch`,
	Example: `
// Switch to context named "sample"
mesheryctl system context switch sample
	`,
	Annotations: linkDocContextSwitch,
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("%s", errArgMsg))
		}
		return nil
	},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return ErrUnmarshallConfig(err)
		}

		contextName := strings.ToLower(args[0])

		_, exists := configuration.Contexts[contextName]
		if !exists {
			const errMsg = `Try running the following to create the context:
mesheryctl system context create `

			return ErrContextNotExists(fmt.Errorf("requested context does not exist \n\n%v%s", errMsg, contextName))
		}
		if viper.GetString("current-context") == contextName {
			return ErrSetCurrentContext(fmt.Errorf("already using context '%s'", contextName))
		}
		//check if meshery is running
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return ErrGetCurrentContext(err)
		}
		isRunning, _ := utils.AreMesheryComponentsRunning(currCtx.GetPlatform())
		// if meshery running stop meshery before context switch
		if isRunning {
			utils.Log.Info("Meshery is running... switching context without stopping Meshery deployments.")
		}

		configuration.CurrentContext = contextName
		viper.Set("current-context", configuration.CurrentContext)
		utils.Log.Infof("switched to context '%s'", contextName)
		err = viper.WriteConfig()
		if err != nil {
			return config.ErrWriteMeshConfig(err)
		}

		return nil
	},
}

// ContextCmd represents the context command
var ContextCmd = &cobra.Command{
	Use:   "context [command]",
	Short: "Configure your Meshery deployment(s)",
	Long: `Configure and switch between different named Meshery server and component versions and deployments.
Find more information at: https://docs.meshery.io/reference/mesheryctl/system/context`,
	Example: `
// Base command
mesheryctl system context
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Display the help message
			_ = cmd.Help()
			os.Exit(0)
		}

		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemContextSubError(fmt.Sprintf("'%s' is an invalid command. Include one of these arguments: [ create | delete | list | switch | view ]. Use 'mesheryctl system context --help' to display sample usage.\n", args[0]), "context"))
		}

		return nil
	},
}

func init() {
	availableSubcommands = []*cobra.Command{
		createContextCmd,
		deleteContextCmd,
		switchContextCmd,
		viewContextCmd,
		listContextCmd,
	}
	createContextCmd.Flags().StringVarP(&contextCreateFlags.URL, "url", "u", "", "Meshery Server URL with Port (default: http://localhost:9081)")
	createContextCmd.Flags().BoolVarP(&contextCreateFlags.Set, "set", "s", false, "Set as current context")
	createContextCmd.Flags().StringArrayVarP(&contextCreateFlags.Components, "components", "a", nil, "List of components")
	createContextCmd.Flags().StringVarP(&contextCreateFlags.Platform, "platform", "p", "", "Platform to deploy Meshery (docker or kubernetes)")
	createContextCmd.Flags().StringVar(&contextCreateFlags.Provider, "provider", "", "Provider to use with the Meshery server (Layer5 or None)")
	deleteContextCmd.Flags().StringVarP(&contextDeleteFlags.Set, "set", "s", "", "New context to deploy Meshery")
	viewContextCmd.Flags().StringVarP(&contextViewFlags.Context, "context", "c", "", "Show config for the context")
	viewContextCmd.Flags().BoolVar(&contextViewFlags.All, "all", false, "Show configs for all of the context")
	ContextCmd.PersistentFlags().StringVarP(&tempCntxt, "context", "c", "", "(optional) temporarily change the current context.")
	ContextCmd.AddCommand(availableSubcommands...)
}

// getYAML takes in a struct and converts it into yaml
func getYAML(strct interface{}) string {
	out, _ := yaml.Marshal(strct)
	return string(out)
}

func getContextWithTokenLocation(c *config.Context) (*contextWithLocation, bool) {
	temp := contextWithLocation{
		Endpoint:      c.Endpoint,
		Token:         c.Token,
		Tokenlocation: tokenNameLocation[c.Token],
		Platform:      c.Platform,
		Components:    c.Components,
		Channel:       c.Channel,
		Version:       c.Version,
		Provider:      c.Provider,
	}
	if temp.Tokenlocation == "" {
		return &temp, false
	}
	return &temp, true
}
