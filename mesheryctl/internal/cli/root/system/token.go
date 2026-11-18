// Copyright Meshery Authors
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

	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	tokenPath     string
	viewAllTokens bool
	set           bool
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage Meshery user tokens",
	Long: `
	Manipulate user tokens and their context assignments in your meshconfig`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			if err := cmd.Help(); err != nil {
				return errors.Wrap(err, "displaying help menu")
			}
			return nil
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\"", args[0])))
		}
		return nil
	},
}

func checkTokenName(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n || args[0] == "" {
			return fmt.Errorf("token name is required in command, accepts %d arg(s), received %d or empty string", n, len(args))
		}
		return nil
	}
}

var createTokenCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a token in your meshconfig",
	Long: `Create the token with provided token name (optionally token path) to your meshconfig tokens.
	Find more information at: https://docs.meshery.io/reference/mesheryctl/system/token/create`,
	Example: `
mesheryctl system token create [token-name] -f [token-path]
mesheryctl system token create [token-name] (default path is auth.json)
mesheryctl system token create [token-name] -f [token-path] --set
	`,
	Args: checkTokenName(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenName := args[0]
		configPath := utils.GetActiveConfigPath()
		resolvedTokenPath := tokenPath
		if resolvedTokenPath == "" {
			resolvedTokenPath = "auth.json"
		}

		token := config.Token{
			Name:     tokenName,
			Location: resolvedTokenPath,
		}
		if err := config.AddTokenToConfig(token, configPath); err != nil {
			return errors.Wrap(err, "Could not create specified token to config")
		}
		utils.Log.Info(fmt.Sprintf("Token %s created.", tokenName))
		if set {
			targetContext := resolveTokenTargetContext(cmd, configPath)
			if err = assignTokenWithRollback(tokenName, configPath, targetContext, config.SetTokenToConfig, config.DeleteTokenFromConfig); err != nil {
				utils.Log.Error(err)
				return err
			}
			utils.Log.Info(fmt.Sprintf("Token: %s set on context %s.", tokenName, targetContext))
		}
		return nil
	},
}
var deleteTokenCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a token from your meshconfig",
	Long:  "Delete the token with provided token name from your meshconfig tokens. A token that a context still uses cannot be deleted; set another token on that context first.",
	Example: `
mesheryctl system token delete [token-name]
	`,
	Args: checkTokenName(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenName := args[0]
		configPath := utils.GetActiveConfigPath()

		if err = config.DeleteTokenFromConfig(tokenName, configPath); err != nil {
			return errors.Wrapf(err, "Could not delete token \"%s\" from config", tokenName)
		}
		utils.Log.Infof("Token %s deleted.", tokenName)
		return nil
	},
}
var setTokenCmd = &cobra.Command{
	Use:   "set",
	Short: "Set token for context",
	Long:  "Set token for current context or context specified with --context flag. The token must already exist in your meshconfig tokens.",
	Example: `
mesheryctl system token set [token-name] 
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenName := args[0]
		configPath := utils.GetActiveConfigPath()
		targetContext := resolveTokenTargetContext(cmd, configPath)

		if err = config.SetTokenToConfig(tokenName, configPath, targetContext); err != nil {
			assignErr := ErrAssignToken(err, tokenName, targetContext)
			utils.Log.Error(assignErr)
			return assignErr
		}
		utils.Log.Infof("Token %s set for context %s", tokenName, targetContext)
		return nil
	},
}
var listTokenCmd = &cobra.Command{
	Use:   "list",
	Short: "List tokens",
	Long:  "List all the tokens in your meshconfig",
	Example: `
mesheryctl system token list
	`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := utils.GetActiveConfigPath()
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return err
		}

		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			return utils.ErrReadConfigFile(err)
		}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		utils.Log.Info("Available tokens: ")
		for _, t := range *mctlCfg.GetTokens() {
			utils.Log.Info(t.Name)
		}
		return nil
	},
}
var viewTokenUsageMsg = "Usage: mesheryctl system token view [token-name]\nRun 'mesheryctl system token view --help' to see detailed help message"

var viewTokenCmd = &cobra.Command{
	Use:   "view [token-name]",
	Short: "View token",
	Long:  "View a specific token in meshery config",
	Example: `
// View a specific token
mesheryctl system token view [token-name]

// View token of current context
mesheryctl system token view

// View all available tokens
mesheryctl system token view --all
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", viewTokenUsageMsg))
		}
		if len(args) > 0 && viewAllTokens {
			return utils.ErrInvalidArgument(fmt.Errorf("cannot specify both a token name and --all flag\n\n%s", viewTokenUsageMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := utils.GetActiveConfigPath()
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return err
		}

		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			return utils.ErrReadConfigFile(err)
		}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		if viewAllTokens {
			utils.Log.Info("Listing all available tokens...\n")
			for _, t := range *mctlCfg.GetTokens() {
				utils.Log.Info("- token: ", t.Name)
				utils.Log.Info("  location: ", t.Location)
			}
			return nil
		}
		tokenName := ""
		if len(args) == 0 {
			targetContext := focusedSystemContext(cmd, mctlCfg.CurrentContext)
			token, err := mctlCfg.GetTokenForContext(targetContext)
			if err != nil {
				return ErrTokenContext(err)
			}
			utils.Log.Warnf("Token unspecified. Displaying token for current context \"%s\"\n", targetContext)
			utils.Log.Info("token: ", token.Name)
			utils.Log.Info("location: ", token.Location)
			return nil
		}
		tokenName = args[0]
		var tokenNames []string
		for _, t := range mctlCfg.Tokens {
			if t.Name == tokenName {
				utils.Log.Info("token: ", t.Name)
				utils.Log.Info("location: ", t.Location)
				return nil
			}
			// Collecting token names in case the provided token name is invalid and a prompt is to be shown.
			tokenNames = append(tokenNames, t.Name)
		}

		if len(tokenNames) == 0 || utils.SilentFlag {
			return fmt.Errorf("token %q not found in config", tokenName)
		}

		utils.Log.Info("Invalid token name. Select from available tokens-")
		prompt := promptui.Select{
			Label: "Select a token from the list",
			Items: tokenNames,
		}
		i, _, err := prompt.Run()
		if err != nil {
			return err
		}
		utils.Log.Info("token: ", mctlCfg.Tokens[i].Name)
		utils.Log.Info("location: ", mctlCfg.Tokens[i].Location)
		return nil
	},
}

func assignTokenWithRollback(
	tokenName, configPath, targetContext string,
	setFn func(tokenName, configPath, ctxName string) error,
	deleteFn func(tokenName, configPath string) error,
) error {
	if err := setFn(tokenName, configPath, targetContext); err != nil {
		// Avoid leaving an orphaned token in meshconfig when --set fails
		// (for example, when -c/--context names a non-existent context).
		if cleanupErr := deleteFn(tokenName, configPath); cleanupErr != nil {
			return ErrTokenRollback(tokenName, targetContext, err, cleanupErr)
		}
		return ErrAssignToken(err, tokenName, targetContext)
	}
	return nil
}

func resolveTokenTargetContext(cmd *cobra.Command, configPath string) string {
	if targetContext := focusedSystemContext(cmd, ""); targetContext != "" {
		return targetContext
	}
	// Read active configPath with an isolated viper instance so default
	// current-context lookup stays accurate without mutating global viper state.
	if _, err := os.Stat(configPath); err == nil {
		v := viper.New()
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err == nil {
			if currentContext := v.GetString("current-context"); currentContext != "" {
				return currentContext
			}
		}
	}
	return viper.GetString("current-context")
}

func init() {
	tokenCmd.AddCommand(createTokenCmd, deleteTokenCmd, setTokenCmd, listTokenCmd, viewTokenCmd)
	createTokenCmd.Flags().StringVarP(&tokenPath, "filepath", "f", "", "Add the token location")
	createTokenCmd.Flags().BoolVarP(&set, "set", "s", false, "Set as current token")
	viewTokenCmd.Flags().BoolVar(&viewAllTokens, "all", false, "set the flag to view all the tokens.")
}
