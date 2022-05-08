package system

import (
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var (
	noneProviderFlag bool
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate to a Meshery Server",
	Long: `
Authenticate to the Local or a Remote Provider of a Meshery Server

The authentication mode is web-based browser flow`,
	Args: cobra.MinimumNArgs(0),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
		hcOptions := &HealthCheckOptions{
			IsPreRunE:  true,
			PrintLogs:  false,
			Subcommand: cmd.Use,
		}
		hc, err := NewHealthChecker(hcOptions)
		if err != nil {
			return errors.Wrapf(err, "failed to initialize healthchecker")
		}
		// execute healthchecks
		err = hc.RunPreflightHealthChecks()
		if err != nil {
			cmd.SilenceUsage = true
		}

		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}

		isRunning, err := utils.IsMesheryRunning(currCtx.GetPlatform())
		if err != nil {
			log.Error("failed to check Meshery Server status: ", err)
			return nil
		}

		if !isRunning {
			log.Error(`Meshery Server is not running. Run "mesheryctl system start" to start Meshery.`)
			return nil
		}

		var tokenData []byte
		if noneProviderFlag {
			tokenData, err = utils.InitiateLoginNone(mctlCfg)
		} else {
			tokenData, err = utils.InitiateLogin(mctlCfg)
		}

		if err != nil {
			log.Println("authentication failed:", err)
			return nil
		}

		log.Println("successfully authenticated")

		token, err := mctlCfg.GetTokenForContext(mctlCfg.GetCurrentContextName())
		if err != nil {
			// Attempt to create token if it doesn't already exists
			token.Location = utils.AuthConfigFile

			// Write new entry in the config
			if err := config.AddTokenToConfig(token, utils.DefaultConfigPath); err != nil {
				log.Error("failed to find token path for the current context")
				return nil
			}
		}

		if err := os.WriteFile(token.GetLocation(), tokenData, 0666); err != nil {
			log.Error("failed to write the token to the filesystem: ", err)
		}

		return nil
	},
}

func init() {
	loginCmd.Flags().BoolVarP(&noneProviderFlag, "none", "", false, "login Meshery with 'None' provider")
}
