package connections

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
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
		const errMsg = "mesheryctl connection delete needs connection-id \nRun 'mesheryctl connection delete --help' to see detailed help message"
		if len(args) != 1 {
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		connectionID := args[0]

		if _, err := uuid.Parse(connectionID); err != nil {
			return utils.ErrInvalidArgument(
				errors.Errorf("invalid connection ID : %q", connectionID),
			)
		}

		_, err := api.Delete(fmt.Sprintf("%s/%s", connectionApiPath, args[0]))
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				return utils.ErrInvalidArgument(
					errors.Errorf("invalid connection ID : %q not found", args[0]),
				)
			}

			return err
		}

		utils.Log.Info("Connection with ID : %q is deleted.", connectionID)
		return nil
	},
}
