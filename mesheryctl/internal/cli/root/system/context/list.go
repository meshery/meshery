package context

import (
	"sort"

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
			// return errors.New("no available contexts")
			log.Print("No contexts available. Use `mesheryctl system context create <name>` to create a new Meshery deployment context.\n")
			return nil
		}

		if currContext == "" {
			currContext = viper.GetString("current-context")
		}
		if currContext == "" {
			log.Print("Current context not set\n")
		} else {
			log.Printf("Current context: %s\n", currContext)
		}
		log.Print("Available contexts:\n")

		//sorting the contexts to get a consistent order on each subsequent run
		keys := make([]string, 0, len(contexts))
		for k := range contexts {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			log.Printf("- %s", k)
		}

		if currContext == "" {
			log.Print("\nRun `mesheryctl system context switch <context name>` to set the current context.")
		}
		return nil
	},
}
