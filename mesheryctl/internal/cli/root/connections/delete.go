package connections

import (
	"fmt"
	"strings"

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
		const errMsg = "[ connection-id ] is required\n\nUsage: mesheryctl connection delete --help' to see detailed help message"
		if len(args) != 1 {
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}

		if !utils.IsUUID(args[0]) {
			return utils.ErrInvalidUUID(fmt.Errorf("invalid connection ID: %q", args[0]))
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := api.Delete(fmt.Sprintf("%s/%s", connectionApiPath, args[0]))
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				return errConnectionNotFound(fmt.Errorf("no connection with id %q found", args[0]))
			}

			return err
		}

		utils.Log.Info("Connection with ID: %q is deleted.", args[0])
		return nil
	},
}
