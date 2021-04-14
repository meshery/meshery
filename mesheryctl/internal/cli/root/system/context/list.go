package context

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var currContext string

// listContextCmd represents the list command
var listContextCmd = &cobra.Command{
	Use:          "list",
	Short:        "list contexts",
	Long:         `List current context and available contexts`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return err
		}
		var contexts = configuration.Contexts
		if contexts == nil {
			return errors.New("no contexts set")
		}

		if currContext == "" {
			currContext = viper.GetString("current-context")
		}
		if currContext == "" {
			return errors.New("current context not set")
		}
		log.Printf("Current context: %s\n", currContext)
		log.Print("Available contexts:\n")
		for context, _ := range contexts {
			log.Printf("%s", context)
		}
		return nil
	},
}

func init() {
	listContextCmd.Flags().StringVar(&currContext, "context", "", "Show config for the context")
}
