// Copyright 2023 Layer5, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"fmt"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configuration     *config.MesheryCtlConfig
	tempCntxt         = "local"
	set               = false
	components        = []string{}
	platform          = ""
	serverURL         = ""
	newContext        = ""
	currContext       string
	allContext        bool
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
	"link":    "![context-create-usage](/assets/img/mesheryctl/newcontext.png)",
	"caption": "Usage of mesheryctl context create",
}

// createContextCmd represents the create command
var createContextCmd = &cobra.Command{
	Use:   "create context-name",
	Short: "Create a new context (a named Meshery deployment)",
	Long:  `Add a new context to Meshery config.yaml file`,
	Example: `
// Create new context
mesheryctl system context create [context-name]

// Create new context and provide list of components, platform & URL
mesheryctl system context create context-name --components meshery-nsm --platform docker --url http://localhost:9081 --set --yes
	`,
	Annotations: linkDocContextCreate,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			const errMsg = `Please provide a context name.
Usage: mesheryctl system context create [context-name]`
			return errors.New(utils.SystemContextSubError(fmt.Sprintf("%s\n", errMsg), "create"))
		}

		tempCntxt := utils.TemplateContext

		if serverURL != "" {
			err := utils.ValidateURL(serverURL)
			if err != nil {
				return err
			}
			tempCntxt.Endpoint = serverURL
		}

		log.Debug("serverURL: `" + tempCntxt.Endpoint + "`")

		if platform != "" {
			tempCntxt.Platform = platform
		}

		if len(components) >= 1 {
			tempCntxt.Components = components
		}

		err := config.AddContextToConfig(args[0], tempCntxt, viper.ConfigFileUsed(), set, false)
		if err != nil {
			return err
		}

		log.Printf("Added `%s` context", args[0])
		return nil
	},
}

// deleteContextCmd represents the delete command
var deleteContextCmd = &cobra.Command{
	Use:   "delete context-name",
	Short: "delete context",
	Long:  `Delete an existing context (a named Meshery deployment) from Meshery config file`,
	Example: `
// Delete context
mesheryctl system context delete [context name]
	`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			const errMsg = `Please provide a context name to delete:
mesheryctl system context delete [context name]`
			return errors.New(utils.SystemContextSubError(fmt.Sprintf("%s\n", errMsg), "delete"))
		}
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return err
		}
		_, exists := configuration.Contexts[args[0]]
		if !exists {
			return fmt.Errorf("no context name found : %s", args[0])
		}

		if viper.GetString("current-context") == args[0] {
			var res bool
			if utils.SilentFlag {
				res = true
			} else {
				res = utils.AskForConfirmation("Are you sure you want to delete the current context")
			}

			if !res {
				log.Printf("Delete aborted")
				return nil
			}

			var result string

			if newContext != "" {
				_, exists := configuration.Contexts[newContext]
				if !exists {
					return errors.New("new context wrongly set")
				}

				if newContext == args[0] {
					return errors.New("choose a new context other than the context being deleted")
				}

				result = newContext
			} else {
				var listContexts []string
				for context := range configuration.Contexts {
					if context != args[0] {
						listContexts = append(listContexts, context)
					}
				}

				prompt := promptui.Select{
					Label: "Select context",
					Items: listContexts,
				}

				_, result, err = prompt.Run()

				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return err
				}
			}

			fmt.Printf("The current context is now %q\n", result)
			viper.Set("current-context", result)
		}
		delete(configuration.Contexts, args[0])
		viper.Set("contexts", configuration.Contexts)
		log.Printf("deleted context %s", args[0])
		err = viper.WriteConfig()

		return err
	},
}

// listContextCmd represents the list command
var listContextCmd = &cobra.Command{
	Use:   "list",
	Short: "list contexts",
	Long:  `List current context and available contexts`,
	Example: `
// List all contexts present
mesheryctl system context list
	`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemContextSubError("this command takes no arguments.\n", "list"))
		}
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return err
		}
		var contexts = configuration.Contexts
		if contexts == nil {
			log.Print("No contexts available. Use `mesheryctl system context create <name>` to create a new Meshery deployment context.\n")
			return nil
		}

		if currContext == "" {
			currContext = viper.GetString("current-context")
		}
		if currContext == "" {
			log.Print("Current context not set\n")
		} else {
			log.Printf("Current context: %s\n", currContext)
		}
		log.Print("Available contexts:\n")

		//sorting the contexts to get a consistent order on each subsequent run
		keys := make([]string, 0, len(contexts))
		for k := range contexts {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			log.Printf("- %s", k)
		}

		if currContext == "" {
			log.Print("\nRun `mesheryctl system context switch <context name>` to set the current context.")
		}
		return nil
	},
}

var linkDocContextView = map[string]string{
	"link":    "![context-view-usage](/assets/img/mesheryctl/context-view.png)",
	"caption": "Usage of mesheryctl context view",
}

// viewContextCmd represents the view command
var viewContextCmd = &cobra.Command{
	Use:   "view [context-name | --context context-name| --all] --flags",
	Short: "view current context",
	Long:  `Display active Meshery context`,
	Example: `
// View default context
mesheryctl system context view

// View specified context
mesheryctl system context view context-name

// View specified context with context flag
mesheryctl system context view --context context-name

// View config of all contexts
mesheryctl system context view --all
	`,
	Annotations:  linkDocContextView,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return err
		}
		//Storing all the tokens separately in a map, to get tokenlocation by token name.
		for _, tok := range configuration.Tokens {
			tokenNameLocation[tok.Name] = tok.Location
		}

		if allContext {
			tempcontexts := make(map[string]contextWithLocation)

			//Populating auxiliary struct with token-locations
			for k, v := range configuration.Contexts {
				if v.Token == "" {
					log.Warnf("[Warning]: Token not specified/empty for context \"%s\"", k)
					temp, _ := getContextWithTokenLocation(&v)
					tempcontexts[k] = *temp
				} else {
					temp, ok := getContextWithTokenLocation(&v)
					tempcontexts[k] = *temp
					if !ok {
						log.Warnf("[Warning]: Token \"%s\" could not be found! for context \"%s\"", tempcontexts[k].Token, k)
					}
				}

			}

			log.Print(getYAML(tempcontexts))

			return nil
		}
		if len(args) != 0 {
			currContext = args[0]
		}
		if currContext == "" {
			currContext = viper.GetString("current-context")

		}
		if currContext == "" {
			return errors.New("current context not set")
		}

		contextData, ok := configuration.Contexts[currContext]
		if !ok {
			log.Printf("context \"%s\" doesn't exists, run the following to create:\n\nmesheryctl system context create %s", currContext, currContext)
			return nil
		}

		if contextData.Token == "" {
			log.Warnf("[Warning]: Token not specified/empty for context \"%s\"", currContext)
			log.Printf("\nCurrent Context: %s\n", currContext)
			log.Print(getYAML(contextData))
		} else {
			temp, ok := getContextWithTokenLocation(&contextData)
			log.Printf("\nCurrent Context: %s\n", currContext)
			if !ok {
				log.Warnf("[Warning]: Token \"%s\" could not be found! for context \"%s\"", temp.Token, currContext)
			}
			log.Print(getYAML(temp))
		}

		return nil
	},
}

var linkDocContextSwitch = map[string]string{
	"link":    "![context-switch-usage](/assets/img/mesheryctl/contextswitch.png)",
	"caption": "Usage of mesheryctl context switch",
}

// switchContextCmd represents the switch command
var switchContextCmd = &cobra.Command{
	Use:   "switch context-name",
	Short: "switch context",
	Long:  `Configure mesheryctl to actively use one one context vs. another context`,
	Example: `
// Switch to context named "sample"
mesheryctl system context switch sample
	`,
	Annotations: linkDocContextSwitch,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = `Usage: mesheryctl system context switch [context name]
Example: mesheryctl system context switch k8s-sample
Description: Configures mesheryctl to actively use one one context vs. the another context`

		if len(args) != 1 {
			return fmt.Errorf("please provide exactly one context name\n\n%v", errMsg)
		}
		return nil
	},
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return err
		}
		_, exists := configuration.Contexts[args[0]]
		if !exists {
			const errMsg = `Try running the following to create the context:
mesheryctl system context create `

			return fmt.Errorf("requested context does not exist \n\n%v%s", errMsg, args[0])
		}
		if viper.GetString("current-context") == args[0] {
			return errors.New("already using context '" + args[0] + "'")
		}
		//check if meshery is running
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		isRunning, _ := utils.AreMesheryComponentsRunning(currCtx.GetPlatform())
		//if meshery running stop meshery before context switch
		if isRunning {
			if err := stop(); err != nil {
				return errors.Wrap(err, utils.SystemError("Failed to stop Meshery before switching context"))
			} else if !userResponse && err == nil {
				return nil
			}
		}

		configuration.CurrentContext = args[0]
		viper.Set("current-context", configuration.CurrentContext)
		log.Printf("switched to context '%s'", args[0])
		err = viper.WriteConfig()
		if isRunning {
			if Starterr := start(); Starterr != nil {
				return errors.Wrap(Starterr, utils.SystemError("Failed to start Meshery while switching context"))
			}
		}
		return err
	},
}

// ContextCmd represents the context command
var ContextCmd = &cobra.Command{
	Use:   "context [command]",
	Short: "Configure your Meshery deployment(s)",
	Long:  `Configure and switch between different named Meshery server and component versions and deployments.`,
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
			return errors.New(utils.SystemContextSubError(fmt.Sprintf("'%s' is a invalid command. Include one of these arguments: [ create | delete | list | switch | view ]. Use 'mesheryctl system context --help' to display sample usage.\n", args[0]), "context"))
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
	createContextCmd.Flags().StringVarP(&serverURL, "url", "u", "", "Meshery Server URL with Port")
	createContextCmd.Flags().BoolVarP(&set, "set", "s", false, "Set as current context")
	createContextCmd.Flags().StringArrayVarP(&components, "components", "a", []string{}, "List of components")
	createContextCmd.Flags().StringVarP(&platform, "platform", "p", "", "Platform to deploy Meshery")
	deleteContextCmd.Flags().StringVarP(&newContext, "set", "s", "", "New context to deploy Meshery")
	viewContextCmd.Flags().StringVarP(&currContext, "context", "c", "", "Show config for the context")
	viewContextCmd.Flags().BoolVar(&allContext, "all", false, "Show configs for all of the context")
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
