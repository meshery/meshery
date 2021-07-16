package system

import (
	"fmt"
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ctx string
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Perform CRUD operations on token",
	Long: `
	Add, Delete and Modify tokens in meshery config`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\"", args[0])))
		}
		return nil
	},
}

var addTokenCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a token to meshery config",
	Example: `
	mesheryctl system token add <token-name> -f <token-path>
	mesheryctl system token add <token-name> (default path is auth.json)
	`,
	Args: cobra.ExactArgs(1),
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
			return errors.Wrap(err, "Could not add specified token to config")
		}
		log.Printf("Token %s added successfully!", tokenName)
		return nil
	},
}
var deleteTokenCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a token from meshery config",
	Example: `
	mesheryctl system token delete <token-name>
	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenName := args[0]

		token := config.Token{
			Name:     tokenName,
			Location: tokenPath,
		}
		if err := config.DeleteTokenFromConfig(token, utils.DefaultConfigPath); err != nil {
			return errors.Wrap(err, "Could not delete specified token from config")
		}
		log.Printf("Token %s deleted successfully!", tokenName)
		return nil
	},
}
var setTokenCmd = &cobra.Command{
	Use:   "set",
	Short: "Set token for context",
	Long:  "Set Token for current context or context passed with --context flag.",
	Example: `
	mesheryctl system token set <token-name> 

	`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenName := args[0]
		if ctx == "" {
			ctx = viper.GetString("current-context")

		}
		if err := config.SetTokenToConfig(tokenName, utils.DefaultConfigPath, ctx); err != nil {
			return errors.Wrapf(err, "Could not set specified token on conteext %s", ctx)

		}
		log.Printf("Token %s set successfully! for context %s", tokenName, ctx)
		return nil
	},
}
var listTokenCmd = &cobra.Command{
	Use:   "list",
	Short: "List tokens",
	Long:  "List all the tokens in meshery config",
	Example: `
	mesheryctl system token list
	`,
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := utils.DefaultConfigPath
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Error(err)
		}

		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			return err

		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		log.Print("Available tokens: ")
		for _, t := range mctlCfg.Tokens {
			log.Info(t.Name)
		}
		return nil
	},
}
var viewTokenCmd = &cobra.Command{
	Use:   "view",
	Short: "View token",
	Long:  "View a specific token in meshery config",
	Example: `
	mesheryctl system token view <token-name>
	mesheryctl system token view (show token of current context)
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tokenName := ""
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Error(err)
		}
		if len(args) == 0 {
			token, err := mctlCfg.GetTokenForContext(viper.GetString("current-context"))
			if err != nil {
				return errors.Wrap(err, "Could not get token for the current context")
			}
			log.Info("token: ", token.Name)
			log.Info("location: ", token.Location)
			return nil
		}
		tokenName = args[0]

		for _, t := range mctlCfg.Tokens {
			if t.Name == tokenName {
				log.Info("token: ", t.Name)
				log.Info("location: ", t.Location)
				return nil
			}
		}
		return errors.Errorf("Token %s could not be found.", tokenName)
	},
}

func init() {
	tokenCmd.AddCommand(addTokenCmd, deleteTokenCmd, setTokenCmd, listTokenCmd, viewTokenCmd)
	addTokenCmd.Flags().StringVarP(&tokenPath, "filepath", "f", "", "Add the token location")
	setTokenCmd.Flags().StringVar(&ctx, "context", "", "Pass the context")
}
