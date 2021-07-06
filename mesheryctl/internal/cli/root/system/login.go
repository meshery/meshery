package system

import (
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate mesheryctl",
	Long:  `Authenticate meheryctl`,
	Args:  cobra.MinimumNArgs(0),
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
		log.Println(string(tokenData))

		return nil
	},
}
