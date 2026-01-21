package connections

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var deleteConnectionCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a connection",
	Long:  `Delete a connection`,

	Example: `
// Delete a connection
mesheryctl connection delete [connection_id]
`,

	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl connection delete \nRun 'mesheryctl connection delete --help' to see detailed help message"
		if len(args) != 1 {
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := api.Delete(fmt.Sprintf("%s/%s", connectionApiPath, args[0]))
		if err != nil {
			return err
		}

		utils.Log.Info("Connection deleted.")
		return nil
	},
}
