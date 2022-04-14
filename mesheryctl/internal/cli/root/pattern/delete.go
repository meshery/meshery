package pattern

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
	Short: "Delete pattern file",
	Long:  `delete pattern file will trigger deletion of the pattern file`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
	// delete a pattern file
	mesheryctl pattern delete [pattern-name]
	`,
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
		req, err := utils.NewRequest("DELETE", mctlCfg.GetBaseMesheryURL()+"/api/pattern/deploy", fileReader)
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

func init() {
	deleteCmd.Flags().StringVarP(&file, "file", "f", "", "Path to pattern file")
}
