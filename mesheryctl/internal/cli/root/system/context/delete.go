package context

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteContextCmd represents the update command
var deleteContextCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete context",
	Long:  `delete specified context from config`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Delete Context %s from config", args[0])
		_, exists := configuration.Contexts[args[0]]
		if !exists {
			return errors.New("context requested to delete does not exist")
		}
		delete(configuration.Contexts, args[0])
		viper.Set("contexts", configuration.Contexts)
		err := viper.WriteConfig()
		return err
	},
}

func init() {

}
