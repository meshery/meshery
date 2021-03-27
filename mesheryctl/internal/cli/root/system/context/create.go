package context

import (
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	url       = ""
	tokenPath = ""
	set       = false
	adapters  = []string{}
	platform  = ""
)

// createContextCmd represents the update command
var createContextCmd = &cobra.Command{
	Use:   "create <context-name>",
	Short: "Create a new context (a named Meshery deployment)",
	Long:  `Add a new context to Meshery config.yaml file`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tempContext := utils.TemplateContext

		err := utils.ValidateURL(url)
		if err != nil {
			return err
		}
		tempContext.Endpoint = url
		tempContext.Platform = platform
		if len(adapters) >= 1 {
			tempContext.Adapters = adapters
		}

		err = utils.AddContextToConfig(args[0], tempContext, viper.ConfigFileUsed(), set)
		if err != nil {
			return err
		}

		log.Printf("Added `%s` context", args[0])
		return nil
	},
}

func init() {
	createContextCmd.Flags().StringVarP(&url, "url", "u", utils.MesheryEndpoint, "Meshery Server URL with Port")
	createContextCmd.Flags().BoolVarP(&set, "set", "s", false, "Set as current context")
	createContextCmd.Flags().StringArrayVarP(&adapters, "adapters", "a", []string{}, "List of adapters")
	createContextCmd.Flags().StringVarP(&platform, "platform", "p", "docker", "Platform to deploy Meshery")
}
