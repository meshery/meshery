package context

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// viewContextCmd represents the update command
var viewContextCmd = &cobra.Command{
	Use:          "view",
	Short:        "view current context",
	Long:         `Display active Meshery context`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return err
		}
		currentContext := viper.GetString("current-context")
		if currentContext == "" {
			return errors.New("context not set")
		}
		log.Println("Current context:", currentContext)
		return nil
	},
}
