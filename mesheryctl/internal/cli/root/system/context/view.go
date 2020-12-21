package context

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// viewContextCmd represents the update command
var viewContextCmd = &cobra.Command{
	Use:   "view",
	Short: "view current context",
	Long:  `view current active meshery context`,
	Args:  cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		currentContext := viper.GetString("current-context")
		log.Println("file used", viper.ConfigFileUsed())
		log.Println("Current context:", currentContext)
		return nil
	},
}
