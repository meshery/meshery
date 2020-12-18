package context

import (
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	url       = ""
	tokenPath = ""
)

// createContextCmd represents the update command
var createContextCmd = &cobra.Command{
	Use:   "create",
	Short: "create context",
	Long:  `Add a new context to mesheryctl config`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Adding Context %s to config", args[0])
		log.Println("Setting as Current Context")
		_, exists := configuration.Contexts[args[0]]
		if exists {
			return errors.New("error adding context, a context with same name already exists")
		}
		configuration.Contexts[args[0]] = models.Context{BaseMesheryURL: url}
		configuration.CurrentContext = args[0]
		viper.Set("contexts", configuration.Contexts)
		viper.Set("currentcontext", configuration.CurrentContext)
		err := viper.WriteConfig()

		return err
	},
}

func init() {
	createContextCmd.Flags().StringVarP(&url, "url", "u", "", "Meshery Server URL with Port")
	_ = createContextCmd.MarkFlagRequired("url")
	createContextCmd.Flags().StringVarP(&tokenPath, "token", "t", "", "Token to be used for Authentication")
	_ = createContextCmd.MarkFlagRequired("token")
}
