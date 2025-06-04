package connections

import (
	"fmt"
	"net/http"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var deleteConnectionCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a connection",
	Long: `Delete a connection providing the connection ID.
Documentation for connection can be found at https://docs.meshery.io/reference/mesheryctl/exp/connection/delete`,

	Example: `
// Delete a specific connection
mesheryctl exp connection delete [connection_id]
`,

	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "[ Connection ID ] isn't specified\n\nUsage: mesheryctl exp connection delete \nRun 'mesheryctl exp connection delete --help' to see detailed help message"
		if len(args) != 1 {
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		connectionDeletePath := fmt.Sprintf("%s/%s", connectionApiPath, args[0])

		resp, err := api.Delete(connectionDeletePath)
		if err != nil {
			return err
		}

		// Check if the response status code is 200
		if resp.StatusCode == http.StatusOK {
			utils.Log.Info("Connection deleted")
			return nil
		}

		return utils.ErrBadRequest(errors.New(fmt.Sprintf("failed to delete connection with id %s", args[0])))
	},
}
