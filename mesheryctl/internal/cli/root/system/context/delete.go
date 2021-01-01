package context

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteContextCmd represents the update command
var deleteContextCmd = &cobra.Command{
	Use:   "delete <context-name>",
	Short: "delete context",
	Long:  `Delete an existing context (a named Meshery deployment) from Meshery config file`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return err
		}
		_, exists := configuration.Contexts[args[0]]
		if !exists {
			return errors.New("no context to delete")
		}
		delete(configuration.Contexts, args[0])
		if viper.GetString("current-context") == args[0] {
			viper.Set("current-context", "")
		}
		viper.Set("contexts", configuration.Contexts)
		log.Printf("deleted context %s", args[0])
		err = viper.WriteConfig()
		return err
	},
}
