package system

import (
	"io/ioutil"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with meshery server",
	Long: `
Authenticate with meshery server

The authetication mode is web-based browser flow`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		tokenData, err := utils.InitiateLogin(mctlCfg)
		if err != nil {
			log.Println("authentication failed")
			return nil
		}

		log.Println("successfully authenticated")

		token, err := mctlCfg.GetTokenForContext(mctlCfg.CurrentContext)
		if err != nil {
			log.Error("failed to find token path for the current context")
			return nil
		}

		if err := ioutil.WriteFile(token.GetLocation(), tokenData, 0666); err != nil {
			log.Error("failed to write the token to the filesystem: ", err)
		}

		return nil
	},
}
