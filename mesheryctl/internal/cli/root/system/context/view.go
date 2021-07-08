package context

import (
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var context string
var allContext bool
var tokenNameLocation map[string]string = map[string]string{} //maps each token name to its specified location

// viewContextCmd represents the view command
var viewContextCmd = &cobra.Command{
	Use:          "view",
	Short:        "view current context",
	Long:         `Display active Meshery context`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.Unmarshal(&configuration)
		if err != nil {
			return err
		}
		//Storing all the tokens separately in a map, to get tokenlocation by token name.
		for _, tok := range configuration.Tokens {
			tokenNameLocation[tok.Name] = tok.Location
		}

		if allContext {
			for k, v := range configuration.Contexts {
				if v.Token == "" {
					log.Warnf("[Warning]: Token not specified/empty for context \"%s\"", k)
					v.TokenLocation = ""
				} else {
					addTokenLocation(&v, k)
				}
				configuration.Contexts[k] = v
			}

			log.Print(getYAML(configuration.Contexts))
			return nil
		}
		if len(args) != 0 {
			context = args[0]
		}
		if context == "" {
			context = viper.GetString("current-context")

		}
		if context == "" {
			return errors.New("current context not set")
		}

		contextData, ok := configuration.Contexts[context]
		if !ok {
			log.Printf("context \"%s\" doesn't exists, run the following to create:\n\nmesheryctl system context create %s", context, context)
			return nil
		}
		if contextData.Token == "" {
			log.Warnf("[Warning]: Token not specified/empty for context \"%s\"", context)
			contextData.TokenLocation = ""
		} else {
			addTokenLocation(&contextData, context)
		}
		log.Printf("\nCurrent Context: %s\n", context)
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

func addTokenLocation(c *config.Context, name string) {
	tokenName := c.Token
	tokenlocation, ok := tokenNameLocation[tokenName]
	if !ok {
		log.Warnf("[Warning]: Token \"%s\" could not be found! for context \"%s\"", tokenName, name)
	}
	c.TokenLocation = tokenlocation
}
