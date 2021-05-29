package context

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"fmt"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/manifoldco/promptui"
)

// deleteContextCmd represents the delete command
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

		if viper.GetString("current-context") == args[0] {
			var res bool
			res = utils.AskForConfirmation("Are you sure you want to delete the current context")

			if !res {
				log.Printf("Delete aborted")
				return nil
			}

			var listContexts []string
			for context := range configuration.Contexts {
				if context != args[0] {
					listContexts = append(listContexts, context)
				}
			}

			prompt := promptui.Select{
				Label: "Select context",
				Items: listContexts,
			}

			_, result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return err
			}

			fmt.Printf("The current context is now %q\n", result)
			viper.Set("current-context", result)
		}
		delete(configuration.Contexts, args[0])
		viper.Set("contexts", configuration.Contexts)
		log.Printf("deleted context %s", args[0])
		err = viper.WriteConfig()

		return err
	},
}
