package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import app manifests",
	Long:  `Import the app manifest into Meshery`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// Import app manifest
mesheryctl app import -f [file/URL] -s [source-type]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return getSourceTypes()
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		var err error

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		appURL := mctlCfg.GetBaseMesheryURL() + "/api/application"

		// If app file is passed via flags
		if !isValidSource(sourceType) {
			return errors.Errorf("application source type (-s) invalid or not passed.\nAllowed source types: %s", strings.Join(validSourceTypes, ", "))
		}

		app, err := importApp(sourceType, file, appURL, true)

		if err != nil {
			return err
		}

		fmt.Printf("App file imported successfully. \nID of the app: %s \n", utils.TruncateID(app.ID.String()))

		return nil
	},
}

func importApp(sourceType string, file string, appURL string, save bool) (*models.MesheryApplication, error) {
	var req *http.Request
	var app *models.MesheryApplication

	client := &http.Client{}

	// Check if the app manifest is file or URL
	if validURL := govalidator.IsURL(file); !validURL {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		text := string(content)

		jsonValues, err := json.Marshal(map[string]interface{}{
			"application_data": map[string]interface{}{
				"name":             path.Base(file),
				"application_file": text,
			},
			"save": save,
		})
		if err != nil {
			return nil, err
		}
		req, err = utils.NewRequest("POST", appURL+"/"+sourceType, bytes.NewBuffer(jsonValues))
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		utils.Log.Debug("app file saved")
		var response []*models.MesheryApplication
		// failsafe (bad api call)
		if resp.StatusCode != 200 {
			return nil, errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, utils.AppError("failed to read response body"))
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal response body")
		}
		// set app
		app = response[0]
	} else {
		var jsonValues []byte
		url, path, err := utils.ParseURLGithub(file)
		if err != nil {
			return nil, err
		}

		utils.Log.Debug(url)
		utils.Log.Debug(path)

		// save the app with Github URL
		if path != "" {
			jsonValues, _ = json.Marshal(map[string]interface{}{
				"url":  url,
				"path": path,
				"save": save,
			})
		} else {
			jsonValues, _ = json.Marshal(map[string]interface{}{
				"url":  url,
				"save": save,
			})
		}

		req, err = utils.NewRequest("POST", appURL+"/"+sourceType, bytes.NewBuffer(jsonValues))
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		utils.Log.Debug("remote hosted app request success")
		var response []*models.MesheryApplication
		// failsafe (bad api call)
		if resp.StatusCode != 200 {
			return nil, errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, utils.AppError("failed to read response body"))
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal response body")
		}

		// set app
		app = response[0]
	}

	return app, nil
}

func init() {
	importCmd.Flags().StringVarP(&file, "file", "f", "", "Path/URL to app file")
	importCmd.Flags().StringVarP(&sourceType, "source-type", "s", "", "Type of source file (ex. manifest / compose / helm)")
}
