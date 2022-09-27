package filter

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

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
		client := &http.Client{}
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		filter := ""
		isID := false
		if len(args) > 0 {
			filter = args[0]
			// It checks if filterID is present or not
			filterID, err := utils.GetID("filter")
			if err == nil {
				for _, id := range filterID {
					if strings.HasPrefix(id, filter) {
						filter = id
					}
				}
			}
			// check if the filter argument is a valid uuid v4 string
			isID, err = regexp.MatchString("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$", filter)
			if err != nil {
				return err
			}
		}

		// Delete the filter using the id
		if isID {
			err := utils.DeleteConfiguration(filter, "filter")
			if err != nil {
				return errors.Wrap(err, utils.SystemError(fmt.Sprintf("failed to delete filter %s", args[0])))
			}
			fmt.Printf("Filter %s deleted successfully\n", args[0])
			return nil
		}

		// Read file
		fileReader, err := os.Open(file)
		if err != nil {
			return errors.New(utils.SystemError(fmt.Sprintf("failed to read file %s", file)))
		}

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
