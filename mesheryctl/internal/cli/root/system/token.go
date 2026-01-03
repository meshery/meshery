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
	ctx           string
	viewAllTokens bool
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
	Long:  "Create the token with provided token name (optionally token path) to your meshconfig tokens.",
	Example: `
mesheryctl system token create [token-name] -f [token-path]
mesheryctl system token create [token-name] (default path is auth.json)
mesheryctl system token create [token-name] -f [token-path] --set
	`,
	Args: checkTokenName(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenName := args[0]
		if tokenPath == "" {
			tokenPath = "auth.json"
		}

		token := config.Token{
			Name:     tokenName,
			Location: tokenPath,
		}
		if err := config.AddTokenToConfig(token, utils.DefaultConfigPath); err != nil {
			return errors.Wrap(err, "Could not create specified token to config")
		}
		utils.Log.Info(fmt.Sprintf("Token %s created.", tokenName))
		if set {
			if ctx == "" {
				ctx = viper.GetString("current-context")
			}
			if err = config.SetTokenToConfig(tokenName, utils.DefaultConfigPath, ctx); err != nil {
				return errors.Wrapf(err, "Could not set token \"%s\" on context %s", tokenName, ctx)
			}
			utils.Log.Infof(fmt.Sprintf("Token: %s set on context %s.", tokenName, ctx))
		}
		return nil
	},
}
var deleteTokenCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a token from your meshconfig",
	Long:  "Delete the token with provided token name from your meshconfig tokens.",
	Example: `
mesheryctl system token delete [token-name]
	`,
	Args: checkTokenName(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenName := args[0]

		if err = config.DeleteTokenFromConfig(tokenName, utils.DefaultConfigPath); err != nil {
			return errors.Wrapf(err, "Could not delete token \"%s\" from config", tokenName)
		}
		cmd.Printf("Token %s deleted.", tokenName)
		return nil
	},
}
var setTokenCmd = &cobra.Command{
	Use:   "set",
	Short: "Set token for context",
	Long:  "Set token for current context or context specified with --context flag.",
	Example: `
mesheryctl system token set [token-name] 
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenName := args[0]
		if ctx == "" {
			ctx = viper.GetString("current-context")
		}

		if err = config.SetTokenToConfig(tokenName, utils.DefaultConfigPath, ctx); err != nil {
			return errors.Wrapf(err, "Could not set token \"%s\" on context %s", tokenName, ctx)

		}
		cmd.Printf("Token %s set for context %s", tokenName, ctx)
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
		if _, err := os.Stat(utils.DefaultConfigPath); os.IsNotExist(err) {
			return err
		}

		viper.SetConfigFile(utils.DefaultConfigPath)
		err := viper.ReadInConfig()
		if err != nil {
			utils.LogError.Error(utils.ErrReadConfigFile(err))
			return nil
		}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.LogError.Error(err)
			return nil
		}
		cmd.Print("Available tokens: ")
		for _, t := range *mctlCfg.GetTokens() {
			utils.Log.Info(t.Name)
		}
		return nil
	},
}
var viewTokenCmd = &cobra.Command{
	Use:   "view",
	Short: "View token",
	Long:  "View a specific token in meshery config",
	Example: `
mesheryctl system token view [token-name]
mesheryctl system token view (show token of current context)
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(utils.DefaultConfigPath); os.IsNotExist(err) {
			return err
		}

		viper.SetConfigFile(utils.DefaultConfigPath)
		err := viper.ReadInConfig()
		if err != nil {
			utils.LogError.Error(utils.ErrReadConfigFile(err))
			return nil
		}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.LogError.Error(err)
			return nil
		}
		if viewAllTokens {
			utils.Log.Info("Listing all available tokens...\n")
			for _, t := range *mctlCfg.GetTokens() {
				utils.Log.Infof("- token: %s", t.Name)
				utils.Log.Infof("  location: %s", t.Location)
			}
			return nil
		}
		tokenName := ""
		if len(args) == 0 {
			token, err := mctlCfg.GetTokenForContext(viper.GetString("current-context"))
			if err != nil {
				utils.LogError.Error(ErrTokenContext(err))
				return nil
			}
			utils.Log.Warnf("Token unspecified. Displaying token for current context \"%s\"\n", viper.GetString("current-context"))
			utils.Log.Infof("token: %s", token.Name)
			utils.Log.Infof("location: %s", token.Location)
			return nil
		}
		tokenName = args[0]
		var tokenNames []string
		for _, t := range mctlCfg.Tokens {
			if t.Name == tokenName {
				utils.Log.Infof("token: %s", t.Name)
				utils.Log.Infof("location: %s", t.Location)
				return nil
			}
			// Collecting token names in case the provided token name is invalid and a prompt is to be shown.
			tokenNames = append(tokenNames, t.Name)
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
		utils.Log.Infof("token: %s", mctlCfg.Tokens[i].Name)
		utils.Log.Infof("location: %s", mctlCfg.Tokens[i].Location)
		return nil
	},
}

func init() {
	tokenCmd.AddCommand(createTokenCmd, deleteTokenCmd, setTokenCmd, listTokenCmd, viewTokenCmd)
	createTokenCmd.Flags().StringVarP(&tokenPath, "filepath", "f", "", "Add the token location")
	createTokenCmd.Flags().BoolVarP(&set, "set", "s", false, "Set as current token")
	setTokenCmd.Flags().StringVar(&ctx, "context", "", "Pass the context")
	viewTokenCmd.Flags().BoolVar(&viewAllTokens, "all", false, "set the flag to view all the tokens.")
}
