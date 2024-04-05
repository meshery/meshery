package connections

import (
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteConnectionCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a connection",
	Long: `Delete
a connection`,

	Example: `
// Delete a connection
mesheryctl exp connections delete [connection_id]
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			utils.Log.Error(err)
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(system.ErrGetCurrentContext(err))
			return nil
		}
		err = ctx.ValidateVersion()
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		return nil
	},

	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp connection delete \nRun 'mesheryctl exp connection delete --help' to see detailed help message"
		if len(args) != 1 {
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := fmt.Sprintf("%s/api/integrations/connections/%s", baseUrl, args[0])
		req, err := utils.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		// Check if the response status code is 200
		if resp.StatusCode == http.StatusOK {
			utils.Log.Info("Connection deleted successfully")
			return nil
		}

		return utils.ErrBadRequest(errors.New(fmt.Sprintf("failed to delete connection with id %s", args[0])))
	},
}
