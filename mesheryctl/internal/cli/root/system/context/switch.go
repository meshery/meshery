package context

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// switchContextCmd represents the update command
var switchContextCmd = &cobra.Command{
	Use:          "switch <context-name>",
	Short:        "switch context",
	Long:         `Configure mesheryctl to actively use one one context vs. the another context`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return err
		}
		_, exists := configuration.Contexts[args[0]]
		if !exists {
			return errors.New("requested context does not exist")
		}
		if viper.GetString("current-context") == args[0] {
			return errors.New("already using context '" + args[0] + "'")
		}
		configuration.CurrentContext = args[0]
		viper.Set("current-context", configuration.CurrentContext)
		log.Printf("switched to context '%s'", args[0])
		err = viper.WriteConfig()
		return err
	},
}
