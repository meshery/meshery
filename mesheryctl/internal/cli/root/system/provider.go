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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	showProviderForAllContext bool
	forceSetProvider          bool
)

// PrintProviderToStdout to return provider details for a context
func PrintProviderToStdout(ctx config.Context, contextName string) string {
	return fmt.Sprintf("Context: %v\nProvider: %v", contextName, ctx.Provider)
}

var viewProviderCmd = &cobra.Command{
	Use:   "view",
	Short: "view provider",
	Long:  "View provider of context in focus",
	Example: `
// View current provider
mesheryctl system provider view
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemProviderSubError("this command takes no arguments\n", "view"))
		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		focusedContext := tempContext
		if focusedContext == "" {
			focusedContext = mctlCfg.CurrentContext
		}

		if showProviderForAllContext {
			for k, v := range mctlCfg.Contexts {
				log.Println(PrintProviderToStdout(v, k))
				log.Println()
			}
			log.Printf("Current Context: %v", focusedContext)
			return nil
		}

		err = mctlCfg.SetCurrentContext(focusedContext)
		if err != nil {
			utils.Log.Error(ErrSetCurrentContext(err))
			return nil
		}

		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(ErrGetCurrentContext(err))
			return nil
		}
		log.Print(PrintProviderToStdout(*currCtx, focusedContext))
		log.Println()
		return nil
	},
}

// listProviderCmd represents the list command
var listProviderCmd = &cobra.Command{
	Use:   "list",
	Short: "list available providers",
	Long:  "List current provider and available providers",
	Example: `
// List all available providers
mesheryctl system provider list
	`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemProviderSubError("this command takes no arguments\n", "list"))
		}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		focusedContext := tempContext
		if focusedContext == "" {
			focusedContext = mctlCfg.CurrentContext
		}

		err = mctlCfg.SetCurrentContext(focusedContext)
		if err != nil {
			utils.Log.Error(ErrSetCurrentContext(err))
			return nil
		}

		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(ErrGetCurrentContext(err))
			return nil
		}

		log.Printf("Current provider: %s\n", currCtx.Provider)

		providers, err := utils.GetProviderInfo(mctlCfg)
		if err != nil {
			log.Fatalln("could not fetch providers as Meshery server was unreachable\nStart Meshery to list available providers")
		}

		log.Print("Available providers:\n")

		//sorting the contexts to get a consistent order on each subsequent run
		keys := make([]string, 0, len(providers))
		for k := range providers {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			log.Printf("- %s", k)
		}

		if currCtx.Provider == "" {
			if tempContext == "" {
				log.Print("\nRun `mesheryctl system provider set [provider]` to set the provider")
			} else {
				log.Printf("\nRun `mesheryctl system provider set [provider] --context %s` to set the provider", tempContext)
			}
		}
		log.Println()
		return nil
	},
}

var setProviderCmd = &cobra.Command{
	Use:   "set [provider]",
	Short: "set provider",
	Long:  "Set provider of context in focus. Run `mesheryctl system provider list` to see the available providers.",
	Example: `
// Set provider
mesheryctl system provider set [provider]
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl system provider set [provider]\nRun 'mesheryctl system provider list' to see the available providers"
		if len(args) == 0 {
			return fmt.Errorf("provider not specified\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
		}

		focusedContext := tempContext
		if focusedContext == "" {
			focusedContext = mctlCfg.CurrentContext
		}

		err = mctlCfg.SetCurrentContext(focusedContext)
		if err != nil {
			utils.Log.Error(ErrSetCurrentContext(err))
			return nil
		}

		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(ErrGetCurrentContext(err))
			return nil
		}

		provider := args[0]

		if !forceSetProvider {
			// Verify provider
			availableProviders, err := utils.GetProviderInfo(mctlCfg)
			if err != nil {
				utils.Log.Error(ErrProviderInfo(err))
			}

			keys := make([]string, 0, len(availableProviders))
			isValidProvider := false

			for k := range availableProviders {
				if provider == k {
					isValidProvider = true
				}
				keys = append(keys, k)
			}

			if !isValidProvider {
				utils.Log.Error(ErrValidProvider())

				log.Print("Available providers:\n")
				//sorting the contexts to get a consistent order on each subsequent run
				sort.Strings(keys)
				for _, k := range keys {
					log.Printf("- %s", k)
				}
				os.Exit(1)
			}
		}

		currCtx.Provider = provider

		mctlCfg.Contexts[focusedContext] = *currCtx
		viper.Set("contexts", mctlCfg.Contexts)
		err = viper.WriteConfig()
		if err != nil {
			utils.Log.Error(ErrWriteConfig(err))
			return nil
		}

		log.Infof("Provider set to %s", currCtx.Provider)
		return nil
	},
}

var switchProviderCmd = &cobra.Command{
	Use:   "switch [provider]",
	Short: "switch provider and redeploy",
	Long:  "Switch provider of context in focus and redeploy Meshery. Run `mesheryctl system provider list` to see the available providers.",
	Example: `
// Switch provider and redeploy Meshery
mesheryctl system provider switch [provider]
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl system provider switch [provider]\nRun 'mesheryctl system provider list' to see the available providers\n"
		if len(args) == 0 {
			return fmt.Errorf("provider not specified\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments.\n\n%v", errMsg)
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
		hcOptions := &HealthCheckOptions{
			IsPreRunE:  true,
			PrintLogs:  false,
			Subcommand: cmd.Use,
		}
		hc, err := NewHealthChecker(hcOptions)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		return hc.RunPreflightHealthChecks()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		userResponse := false

		mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
		}

		focusedContext := tempContext
		if focusedContext == "" {
			focusedContext = mctlCfg.CurrentContext
		}

		//skip asking confirmation if -y flag used
		if utils.SilentFlag {
			userResponse = true
		} else {
			// ask user for confirmation
			userResponse = utils.AskForConfirmation("The Meshery deployment in context '" + focusedContext + "' will be replaced with a new Meshery deployment with provider set to '" + args[0] + "'. Are you sure you want to continue")
		}

		if !userResponse {
			return errors.New("provider switch aborted")
		}

		err = setProviderCmd.RunE(cmd, args)
		if err != nil {
			return err
		}
		err = restartCmd.RunE(cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

// resetProviderCmd represents the reset command
var resetProviderCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset provider to default",
	Long:  "Reset provider for current context to default (Meshery)",
	Example: `
// Reset provider to default
mesheryctl system provider reset
	`,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemProviderSubError("this command takes no arguments.\n", "reset"))
		}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
		}

		focusedContext := tempContext
		if focusedContext == "" {
			focusedContext = mctlCfg.CurrentContext
		}

		err = mctlCfg.SetCurrentContext(focusedContext)
		if err != nil {
			utils.Log.Error(ErrSetCurrentContext(err))
			return nil
		}

		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(ErrGetCurrentContext(err))
			return nil
		}

		currCtx.Provider = "Meshery"

		mctlCfg.Contexts[focusedContext] = *currCtx
		viper.Set("contexts", mctlCfg.Contexts)
		err = viper.WriteConfig()
		if err != nil {
			utils.Log.Error(ErrWriteConfig(err))
			return nil
		}

		log.Infof("Provider reset to Meshery")
		return nil
	},
}

// providerCmd represents the provider command
var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Switch between providers",
	Long:  `Enforce a provider. Choose between available Meshery providers`,
	Example: `
// To view provider
mesheryctl system provider view
// To list all available providers
mesheryctl system provider list
// To set a provider
mesheryctl system provider set [provider]
// To switch provider and redeploy Meshery
mesheryctl system provider switch [provider]
// To reset provider to default
mesheryctl system provider reset
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New(utils.SystemProviderSubError("please specify a flag or subcommand. Use 'mesheryctl system provider --help' to display user guide.\n", "provider"))
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemProviderSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl system provider --help' to display usage guide.\n", args[0]), "provider"))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
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

func init() {
	viewProviderCmd.Flags().BoolVarP(&showProviderForAllContext, "all", "a", false, "Show provider for all contexts")
	setProviderCmd.Flags().BoolVarP(&forceSetProvider, "force", "f", false, "Force set provider")
	providerCmd.AddCommand(viewProviderCmd, listProviderCmd, setProviderCmd, switchProviderCmd, resetProviderCmd)
}
