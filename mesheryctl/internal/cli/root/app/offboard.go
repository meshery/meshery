package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var offboardCmd = &cobra.Command{
	Use:   "offboard",
	Short: "Offboard application",
	Long:  `Offboard application will trigger undeploy of application`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// Offboard application by providing file path
mesheryctl app offboard -f [filepath]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if file == "" {
			const errMsg = `Usage: mesheryctl app offboard -f [filepath]`
			return fmt.Errorf("no file path provided \n\n%v", errMsg)
		}
		var req *http.Request
		var err error
		client := &http.Client{}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/application/deploy"
		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern"

		// Read file
		if !govalidator.IsURL(file) {
			content, err := os.ReadFile(file)
			if err != nil {
				return errors.New(utils.AppError(fmt.Sprintf("failed to read file %s\n", file)))
			}

			appFile = string(content)
		} else {
			utils.Log.Info("URLs are not currently supported")
		}

		// Convert App File into Pattern File
		jsonValues, _ := json.Marshal(map[string]interface{}{
			"K8sManifest": appFile,
		})

		req, err = utils.NewRequest("POST", patternURL, bytes.NewBuffer(jsonValues))
		if err != nil {
			return err
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var response []*models.MesheryPattern
		// bad api call
		if resp.StatusCode != 200 {
			return errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, utils.PerfError("failed to read response body"))
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal response body")
		}

		utils.Log.Debug("application file converted to pattern file")

		patternFile := response[0].PatternFile

		req, err = utils.NewRequest("DELETE", deployURL, bytes.NewBuffer([]byte(patternFile)))
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if res.StatusCode == 200 {
			utils.Log.Info("app successfully offboarded")
		}
		utils.Log.Info(string(body))

		return nil
	},
}

func init() {
	offboardCmd.Flags().StringVarP(&file, "file", "f", "", "Path to app file")
}
