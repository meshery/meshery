package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
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
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		validTypesURL := mctlCfg.GetBaseMesheryURL() + "/api/application/types"
		client := &http.Client{}
		req, err := utils.NewRequest("GET", validTypesURL, nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
		}
		defer resp.Body.Close()

		var response []*models.ApplicationSourceTypesAPIResponse

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, utils.PerfError("failed to read response body"))
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal response body")
		}

		for _, apiResponse := range response {
			validSourceTypes = append(validSourceTypes, apiResponse.ApplicationType)
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
		var err error
		client := &http.Client{}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		appURL := mctlCfg.GetBaseMesheryURL() + "/api/application"
		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern"

		// If app file is passed via flags
		if !isValidSource(sourceType) {
			return errors.Errorf("application source type (-s) invalid or not passed.\nAllowed source types: %s", strings.Join(validSourceTypes, ", "))
		}

		// Check if the app manifest is file or URL
		if validURL := govalidator.IsURL(file); !validURL {
			content, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			text := string(content)

			// if --skip-save is not passed we save the apps first
			jsonValues, err := json.Marshal(map[string]interface{}{
				"application_data": map[string]interface{}{
					"name":             path.Base(file),
					"application_file": text,
				},
				"save": true,
			})
			if err != nil {
				return err
			}
			req, err = utils.NewRequest("POST", appURL+"/"+sourceType, bytes.NewBuffer(jsonValues))
			if err != nil {
				return err
			}

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			utils.Log.Debug("saved app file")
			var response []*models.MesheryApplication
			// failsafe (bad api call)
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

			// setup app file
			appFile = text

		} else {
			var jsonValues []byte
			url, path, err := utils.ParseURLGithub(file)
			if err != nil {
				return err
			}

			utils.Log.Debug(url)
			utils.Log.Debug(path)

			// save the app with Github URL
			if !skipSave {
				if path != "" {
					jsonValues, _ = json.Marshal(map[string]interface{}{
						"url":  url,
						"path": path,
						"save": true,
					})
				} else {
					jsonValues, _ = json.Marshal(map[string]interface{}{
						"url":  url,
						"save": true,
					})
				}
			} else { // we don't save the app
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
			}
			req, err = utils.NewRequest("POST", appURL+"/"+sourceType, bytes.NewBuffer(jsonValues))
			if err != nil {
				return err
			}

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			utils.Log.Debug("remote hosted app request success")
			var response []*models.MesheryApplication
			// failsafe (bad api call)
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

			// setup app file here
			appFile = response[0].ApplicationFile
		}

		// Start conversion of app file to pattern file
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

		return nil
	},
}

func init() {
	importCmd.Flags().StringVarP(&file, "file", "f", "", "Path/URL to app file")
	importCmd.Flags().StringVarP(&sourceType, "source-type", "s", "", "Type of source file (ex. manifest / compose / helm)")
}
