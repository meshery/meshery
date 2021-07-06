package system

import (
	"io/ioutil"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
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

		token, err := mctlCfg.GetTokenForContext(mctlCfg.CurrentContext)
		if err != nil {
			log.Error("failed to find token path for the current context")
			return nil
		}

		// Replace the content of the token file with empty content
		if err := ioutil.WriteFile(token.GetLocation(), []byte{}, 0666); err != nil {
			log.Error("logout failed: ", err)
			return nil
		}

		log.Println("successfully logged out")
		return nil
	},
}
