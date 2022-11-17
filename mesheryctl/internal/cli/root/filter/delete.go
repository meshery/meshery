package filter

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete filter file",
	Long:  `delete filter file will trigger deletion of the filter file`,
	Example: `
// Delete the specified WASM filter file using name or ID
mesheryctl exp filter delete [filter-name | ID]

// Delete using the file name
mesheryctl exp filter delete test-wasm
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		// Read file
		fileReader, err := os.Open(file)
		if err != nil {
			return errors.New(utils.SystemError(fmt.Sprintf("failed to read file %s", file)))
		}

		client := &http.Client{}
		req, err := utils.NewRequest("DELETE", mctlCfg.GetBaseMesheryURL()+"/api/filter/deploy", fileReader)
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		utils.Log.Info(string(body))

		return nil
	},
}
