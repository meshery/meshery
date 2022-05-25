package pattern

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
	"github.com/layer5io/meshery/models"
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
mesheryctl pattern delete [file | URL]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
		client := &http.Client{}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/deploy"
		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern"

		// If file path not a valid URL, treat it like a local file path
		if !govalidator.IsURL(file) {
			content, err := os.ReadFile(file)
			if err != nil {
				return errors.New(utils.SystemError(fmt.Sprintf("failed to read file %s", file)))
			}

			patternFile = string(content)
		} else {
			// Else treat it like a URL
			url, path, err := utils.ParseURLGithub(file)
			if err != nil {
				return err
			}

			utils.Log.Debug(url)
			utils.Log.Debug(path)

			var jsonValues []byte

			// Send the URL and path to the server and let it fetch the patternfile and delete
			// the components
			if path != "" {
				jsonValues, _ = json.Marshal(map[string]interface{}{
					"url":  url,
					"path": path,
					"save": false,
				})
			} else {
				jsonValues, _ = json.Marshal(map[string]interface{}{
					"url":  url,
					"save": false,
				})
			}

			req, err = utils.NewRequest("POST", patternURL, bytes.NewBuffer(jsonValues))
			if err != nil {
				return err
			}

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			utils.Log.Debug("remote hosted pattern request success")
			var response []*models.MesheryPattern
			// If API returns a non 200 status, return error
			if resp.StatusCode != 200 {
				return errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrap(err, utils.PerfError("failed to read response body"))
			}

			err = json.Unmarshal(body, &response)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal response body")
			}

			patternFile = response[0].PatternFile
		}

		req, err = utils.NewRequest("DELETE", deployURL, bytes.NewBuffer([]byte(patternFile)))
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
