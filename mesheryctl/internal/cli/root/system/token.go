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
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Manage Meshery user tokens",
	Long: `Manipulate user tokens and their context assignments in your meshconfig.
Find more information at: https://docs.meshery.io/reference/mesheryctl/system/token`,
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

type cmdTokenCreateFlags struct {
	Filepath string `json:"filepath" validate:"omitempty"`
	Set      bool   `json:"set" validate:"boolean"`
}

var tokenCreateFlags cmdTokenCreateFlags

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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &tokenCreateFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenName := args[0]
		configPath := utils.GetActiveConfigPath()
		if tokenCreateFlags.Filepath == "" {
			tokenCreateFlags.Filepath = "auth.json"
		}

		token := config.Token{
			Name:     tokenName,
			Location: tokenCreateFlags.Filepath,
		}
		if err := config.AddTokenToConfig(token, configPath); err != nil {
			return errors.Wrap(err, "Could not create specified token to config")
		}
		utils.Log.Info(fmt.Sprintf("Token %s created.", tokenName))
		if tokenCreateFlags.Set {
			ctx := tokenSetFlags.Context
			if ctx == "" {
				ctx = viper.GetString("current-context")
			}
			if err := config.SetTokenToConfig(tokenName, configPath, ctx); err != nil {
				return ErrSetTokenToConfig(err)
			}
			utils.Log.Info(fmt.Sprintf("Token: %s set on context %s.", tokenName, ctx))
		}
		return nil
	},
}

var deleteTokenCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a token from your meshconfig",
	Long: `Delete the token with provided token name from your meshconfig tokens.
Find more information at: https://docs.meshery.io/reference/mesheryctl/system/token/delete`,
	Example: `
mesheryctl system token delete [token-name]
	`,
	Args: checkTokenName(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenName := args[0]
		configPath := utils.GetActiveConfigPath()

		if err := config.DeleteTokenFromConfig(tokenName, configPath); err != nil {
			return ErrDeleteTokenFromConfig(err)
		}
		utils.Log.Infof("Token %s deleted.", tokenName)
		return nil
	},
}

type cmdTokenSetFlags struct {
	Context string `json:"context" validate:"omitempty"`
}

var tokenSetFlags cmdTokenSetFlags

var setTokenCmd = &cobra.Command{
	Use:   "set",
	Short: "Set token for context",
	Long: `Set token for current context or context specified with --context flag.
Find more information at: https://docs.meshery.io/reference/mesheryctl/system/token/set`,
	Example: `
mesheryctl system token set [token-name]
	`,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &tokenSetFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenName := args[0]
		configPath := utils.GetActiveConfigPath()
		if tokenSetFlags.Context == "" {
			tokenSetFlags.Context = viper.GetString("current-context")
		}

		if err := config.SetTokenToConfig(tokenName, configPath, tokenSetFlags.Context); err != nil {
			return ErrSetTokenToConfig(err)
		}
		utils.Log.Infof("Token %s set for context %s", tokenName, tokenSetFlags.Context)
		return nil
	},
}

var listTokenCmd = &cobra.Command{
	Use:   "list",
	Short: "List tokens",
	Long: `List all the tokens in your meshconfig.
Find more information at: https://docs.meshery.io/reference/mesheryctl/system/token/list`,
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
		if err := viper.ReadInConfig(); err != nil {
			return utils.ErrReadConfigFile(err)
		}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		utils.Log.Info("Available tokens: ")
		for _, t := range *mctlCfg.GetTokens() {
			utils.Log.Info(t.Name)
		}
		return nil
	},
}

type cmdTokenViewFlags struct {
	All bool `json:"all" validate:"boolean"`
}

var tokenViewFlags cmdTokenViewFlags

var viewTokenCmd = &cobra.Command{
	Use:   "view",
	Short: "View token",
	Long: `View a specific token in meshery config.
Find more information at: https://docs.meshery.io/reference/mesheryctl/system/token/view`,
	Example: `
mesheryctl system token view [token-name]
mesheryctl system token view (show token of current context)
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &tokenViewFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := utils.GetActiveConfigPath()
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return err
		}

		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			return utils.ErrReadConfigFile(err)
		}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		if tokenViewFlags.All {
			utils.Log.Info("Listing all available tokens...\n")
			for _, t := range *mctlCfg.GetTokens() {
				utils.Log.Info("- token: ", t.Name)
				utils.Log.Info("  location: ", t.Location)
			}
			return nil
		}
		if len(args) == 0 {
			token, err := mctlCfg.GetTokenForContext(viper.GetString("current-context"))
			if err != nil {
				utils.Log.Error(ErrTokenContext(err))
				return nil
			}
			utils.Log.Warnf("Token unspecified. Displaying token for current context \"%s\"\n", viper.GetString("current-context"))
			utils.Log.Info("token: ", token.Name)
			utils.Log.Info("location: ", token.Location)
			return nil
		}
		tokenName := args[0]
		var tokenNames []string
		for _, t := range mctlCfg.Tokens {
			if t.Name == tokenName {
				utils.Log.Info("token: ", t.Name)
				utils.Log.Info("location: ", t.Location)
				return nil
			}
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
		utils.Log.Info("token: ", mctlCfg.Tokens[i].Name)
		utils.Log.Info("location: ", mctlCfg.Tokens[i].Location)
		return nil
	},
}

func checkTokenName(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != n || (len(args) > 0 && args[0] == "") {
			return utils.ErrInvalidArgument(fmt.Errorf("token name is required in command, accepts %d arg(s), received %d or empty string", n, len(args)))
		}
		return nil
	}
}

func init() {
	tokenCmd.AddCommand(createTokenCmd)
	createTokenCmd.Flags().StringVarP(&tokenCreateFlags.Filepath, "filepath", "f", "", "Add the token location")
	createTokenCmd.Flags().BoolVarP(&tokenCreateFlags.Set, "set", "s", false, "Set as current token")

	tokenCmd.AddCommand(deleteTokenCmd)

	tokenCmd.AddCommand(setTokenCmd)
	setTokenCmd.Flags().StringVar(&tokenSetFlags.Context, "context", "", "Pass the context")

	tokenCmd.AddCommand(listTokenCmd)

	tokenCmd.AddCommand(viewTokenCmd)
	viewTokenCmd.Flags().BoolVar(&tokenViewFlags.All, "all", false, "set the flag to view all the tokens.")
}
