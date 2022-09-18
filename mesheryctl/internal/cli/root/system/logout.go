package system

import (
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove authentication for Meshery Server",
	Long: `
Remove authentication for Meshery Server

This command removes the authentication token from the user's filesystem`,
	Args: cobra.MinimumNArgs(0),
	Example: `
// Logout current session with your Meshery Provider.
mesheryctl system logout
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		token, err := mctlCfg.GetTokenForContext(mctlCfg.GetCurrentContextName())
		if err != nil {
			log.Error("failed to find token path for the current context")
			return nil
		}

		// Replace the content of the token file with empty content
		if err := os.WriteFile(token.GetLocation(), []byte{}, 0666); err != nil {
			log.Error("logout failed: ", err)
			return nil
		}

		log.Println("successfully logged out")
		return nil
	},
}
