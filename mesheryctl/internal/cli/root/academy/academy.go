package academy

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands = []*cobra.Command{createCmd}
)

var AcademyCmd = &cobra.Command{
	Use:   "academy",
	Short: "Layer5 Academy related commands",
	Long:  `Manage scaffolding and creation of Layer5 Academy content.`,
	Example: `
// Scaffold a full learning path tree
mesheryctl academy create --type learning-path --title "My Path" --description "Desc" --level beginner --org 123e4567-e89b-12d3-a456-426614174000

// Scaffold a single course into an existing tree
mesheryctl academy create course "New Course" --description "Desc" --into ./my-path
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.RootError(fmt.Sprintf("'%s' is an invalid command. Use 'mesheryctl academy --help' to display usage guide.\n", args[0])))
		}
		return nil
	},
}

func init() {
	AcademyCmd.AddCommand(availableSubcommands...)
}
