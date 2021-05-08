package context

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var context string
var allContext bool

// viewContextCmd represents the view command
var viewContextCmd = &cobra.Command{
	Use:          "view",
	Short:        "view current context",
	Long:         `Display active Meshery context`,
	Args:         cobra.MinimumNArgs(0),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return err
		}

		if allContext {
			log.Print(getYAML(configuration.Contexts))
			return nil
		}
		if len(args) != 0 {
			context = args[0]
		}
		if context == "" {
			context = viper.GetString("current-context")
			log.Printf("***This is your Current context***")
		}
		if context == "" {
			return errors.New("current context not set")
		}

		contextData, ok := configuration.Contexts[context]
		if !ok {
			log.Printf("context \"%s\" doesn't exists, run the following to create:\n\nmesheryctl system context create %s", context, context)
			return nil
		}
		log.Printf("\nContext: %s\n", context)
		log.Print(getYAML(contextData))
		return nil
	},
}

func init() {
	viewContextCmd.Flags().StringVar(&context, "context", "", "Show config for the context")
	viewContextCmd.Flags().BoolVar(&allContext, "all", false, "Show configs for all of the context")
}

// getYAML takes in a struct and converts it into yaml
func getYAML(strct interface{}) string {
	out, _ := yaml.Marshal(strct)
	return string(out)
}
