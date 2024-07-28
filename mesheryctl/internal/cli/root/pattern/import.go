// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pattern

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

var (
	name       string
	sourceType string // source type of the design file
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a Meshery design",
	Long: `
		Import Helm Charts, Kubernetes Manifest, Docker Compose or Meshery designs by passing
		remote URL or local file system path to the file. Source type must be provided.

		YAML and TGZ (with helm only) format of file is accepted, if you are importing Meshery Design OCI file format is also supported

		If you are providing remote URL, it should be a direct URL to a downloadable file.
		For example, if the file is stored on GitHub, the URL should be 'https://raw.githubusercontent.com/path-to-file'.
	`,
	Example: `
// Import pattern manifest
mesheryctl pattern import -f [file/URL] -s [source-type] -n [name]
	`,
	Args: func(_ *cobra.Command, args []string) error {

		if file == "" {
			utils.Log.Debug("manifest path not provided")
			return ErrPatternManifest()
		}

		if sourceType == "" {
			utils.Log.Debug("source-type not provided")
			return ErrPatternSourceType()
		}

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return getSourceTypes()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern"

		// If pattern file is passed via flags
		if sourceType, err = getFullSourceType(sourceType); err != nil {
			return ErrInValidSource(sourceType, validSourceTypes)
		}

		pattern, err := importPattern(sourceType, file, patternURL, true)

		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		fmt.Printf("The design file '%s' has been imported. Design ID: %s, Source Type: %s ", pattern.Name, utils.TruncateID(pattern.ID.String()), sourceType)

		return nil
	},
}

func importPattern(sourceType string, file string, patternURL string, save bool) (*models.MesheryPattern, error) {
	var req *http.Request
	var pattern *models.MesheryPattern

	// If design name is not provided
	// use file name as default
	patternName := path.Base(file)
	if name != "" {
		patternName = name
	}

	// Check if the pattern manifest is file or URL
	if validURL := govalidator.IsURL(file); !validURL {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, utils.ErrFileRead(err)
		}

		jsonValues, err := json.Marshal(map[string]interface{}{
			"pattern_data": map[string]interface{}{
				"name":         patternName,
				"pattern_file": content,
			},
			"save": save,
		})
		if err != nil {
			return nil, utils.ErrMarshal(err)
		}
		req, err = utils.NewRequest("POST", patternURL+"/"+sourceType, bytes.NewBuffer(jsonValues))
		if err != nil {
			return nil, err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			return nil, err
		}
		utils.Log.Debug("pattern file saved")
		var response []*models.MesheryPattern
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Debug("failed to read response body")
			return nil, utils.ErrReadResponseBody(err)
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			utils.Log.Debug("failed to unmarshal JSON response")
			return nil, utils.ErrUnmarshal(err)
		}
		// set pattern
		pattern = response[0]
	} else {
		var jsonValues []byte

		jsonValues, _ = json.Marshal(map[string]interface{}{
			"url":  file,
			"save": save,
			"name": patternName,
		})

		req, err := utils.NewRequest("POST", patternURL+"/"+sourceType, bytes.NewBuffer(jsonValues))
		if err != nil {
			return nil, utils.ErrCreatingRequest(err)
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			return nil, utils.ErrRequestResponse(err)
		}
		utils.Log.Debug("Fetched the design from the remote host")
		var response []*models.MesheryPattern
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Debug("failed to read response body")
			return nil, utils.ErrReadResponseBody(err)
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			utils.Log.Debug("failed to unmarshal JSON response")
			return nil, utils.ErrUnmarshal(err)
		}

		// set pattern
		pattern = response[0]
	}

	return pattern, nil
}

func init() {
	importCmd.Flags().StringVarP(&file, "file", "f", "", "Path/URL to pattern file")
	importCmd.Flags().StringVarP(&sourceType, "source-type", "s", "", "Type of source file (ex. manifest / compose / helm / design)")
	importCmd.Flags().StringVarP(&name, "name", "n", "", "Name for the pattern file")
}

// getSourceTypes retrieves the list of valid source types from the Meshery server.
// It makes a GET request to the "/api/pattern/types" endpoint of the Meshery server
// and populates the `validSourceTypes` slice with the design types received in the response.
// If any error occurs during the process, it logs the error and returns nil.
func getSourceTypes() error {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		utils.Log.Error(err)
		return nil
	}
	validTypesURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/types"
	req, err := utils.NewRequest("GET", validTypesURL, nil)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	defer resp.Body.Close()

	var response []*models.PatternSourceTypesAPIResponse

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Log.Error(utils.ErrReadResponseBody(errors.Wrap(err, "couldn't read response from server. Please try again after some time")))
		return nil
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		utils.Log.Error(utils.ErrUnmarshal(errors.Wrap(err, "couldn't process response received from server")))
		return nil
	}

	for _, apiResponse := range response {
		validSourceTypes = append(validSourceTypes, apiResponse.DesignType)
	}

	return nil
}

// returns full source name e.g. helm -> `Helm Chart`
// user passes only helm, compose or manifest but server accepts full source type
// e.g `Heml Chart`, `Docker Compose`, `Kubernetes Manifest`
func getFullSourceType(sType string) (string, error) {
	for _, validType := range validSourceTypes {
		lowerType := strings.ToLower(validType)
		// user may pass Pascal Case source e.g Helm
		sType = strings.ToLower(sType)
		if strings.Contains(lowerType, sType) {
			return validType, nil
		}
	}

	return sType, fmt.Errorf("no matching source type found")
}
