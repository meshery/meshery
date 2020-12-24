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
	Long:  `Display active Meshery context`,
	Args:  cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return err
		}
		currentContext := viper.GetString("current-context")
		log.Println("Current context:", currentContext)
		return nil
	},
}
