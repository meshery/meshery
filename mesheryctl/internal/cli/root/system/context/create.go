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
	Short: "Create a new context (a named Meshery deployment)",
	Long:  `Add a new context to Meshery config.yaml file`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return err
		}
		_, exists := configuration.Contexts[args[0]]
		if exists {
			return errors.New("error adding context: a context with same name already exists")
		}
		listOfAdapters := []string{"meshery-istio", "meshery-linkerd", "meshery-consul", "meshery-octarine", "meshery-nsm", "meshery-kuma", "meshery-cpx", "meshery-osm", "meshery-nginx-sm"}
		configuration.Contexts[args[0]] = models.Context{Endpoint: url, Adapters: listOfAdapters}
		configuration.CurrentContext = args[0]
		viper.Set("contexts", configuration.Contexts)
		viper.Set("current-context", configuration.CurrentContext)
		err = viper.WriteConfig()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Added %s context and set as current context", args[0])
		return nil
	},
}

func init() {
	createContextCmd.Flags().StringVarP(&url, "url", "u", "localhost:9081", "Meshery Server URL with Port")
}
