package connections

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
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
		if len(args) != 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errNoArgMsg, deleteUsageMsg))
		}

		if !utils.IsUUID(args[0]) {
			return utils.ErrInvalidUUID(fmt.Errorf("invalid connection ID: %q", args[0]))
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := api.Delete(fmt.Sprintf("%s/%s", connectionApiPath, args[0]))
		if err != nil {
			if errors.GetCode(err) == utils.ErrNotFoundCode {
				return utils.ErrNotFound(fmt.Errorf("no connection with ID %q found", args[0]))
			}

			return err
		}

		utils.Log.Infof("Connection with ID: %q is deleted", args[0])
		return nil
	},
}
