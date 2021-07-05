package auth

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
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
			fmt.Println("authentication failed")
			return nil
		}

		fmt.Println("successfully authenticated")
		fmt.Println(string(tokenData))

		return nil
	},
}
