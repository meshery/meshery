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

package design

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/meshery/schemas/models/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var name string

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a Meshery design",
	Long: `
		Import Helm Charts, Kubernetes Manifest, Docker Compose or Meshery designs by passing
		remote URL or local file system path to the file. Providing source type is optional.

		YAML and TGZ (with helm only) format of file is accepted, if you are importing Meshery Design OCI file format is also supported

		If you are providing remote URL, it should be a direct URL to a downloadable file.
		For example, if the file is stored on GitHub, the URL should be 'https://raw.githubusercontent.com/path-to-file'.
	`,
	Example: `
// Import design manifest
mesheryctl design import -f [file/URL] -s [source-type] -n [name]

mesheryctl design import -f design.tar
mesheryctl design import -f design.yml -n design-name
mesheryctl design import -f design.yml -s "Kubernetes Manifest" -n design-name
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl design import -f [file/URL] -s [source-type] -n [name]\n"
		if file == "" {
			utils.Log.Debug("File path not provided\n" + errMsg)
			return ErrDesignFileNotProvided()
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

		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/import"

		// If pattern file is passed via flags
		if sourceType != "" {
			if sourceType, err = getFullSourceType(sourceType); err != nil {
				utils.Log.Debugf("%s is not a valid source type. Valid types are Helm Chart, Kubernetes Manifest, Docker Compose, Meshery Design", sourceType)
				validSourceTypes := []string{"Helm Chart", "Kubernetes Manifest", "Docker Compose", "Meshery Design"}
				return ErrInValidSource(sourceType, validSourceTypes)
			}
		}

		switch sourceType {
		case "Helm Chart":
			sourceType = core.IacFileTypes.HELM_CHART
		case "Kubernetes Manifest":
			sourceType = core.IacFileTypes.KUBERNETES_MANIFEST
		case "Meshery Design":
			sourceType = core.IacFileTypes.MESHERY_DESIGN
		case "Docker Compose":
			sourceType = core.IacFileTypes.DOCKER_COMPOSE
		}

		pattern, err := importPattern(sourceType, file, patternURL, true)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		utils.Log.Infof("The design file '%s' has been imported. Design ID: %s", pattern.Name, utils.TruncateID(pattern.ID.String()))

		return nil
	},
}

func importPattern(sourceType string, file string, patternURL string, save bool) (*models.MesheryPattern, error) {
	var req *http.Request
	var pattern *models.MesheryPattern

	// If design name is not provided
	// use file name as default
	var patternName string
	fileName := path.Base(file)
	if name == "" {
		patternName = fileName
	} else {
		patternName = name
	}

	validURL := utils.IsValidUrl(file)
	// Check if the pattern is file or URL
	if !validURL {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, utils.ErrFileRead(err)
		}

		jsonValues, err := json.Marshal(map[string]interface{}{
			"name":      patternName,
			"file":      content,
			"file_name": fileName,
		})
		if err != nil {
			return nil, utils.ErrMarshal(err)
		}
		req, err = utils.NewRequest("POST", patternURL, bytes.NewBuffer(jsonValues))
		if err != nil {
			return nil, err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			return nil, err
		}
		utils.Log.Debug("design file saved")
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
			"name": patternName,
		})

		req, err := utils.NewRequest("POST", patternURL, bytes.NewBuffer(jsonValues))
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
	importCmd.Flags().StringVarP(&file, "file", "f", "", "Path/URL to design file")
	importCmd.Flags().StringVarP(&sourceType, "source-type", "s", "", "Type of source file (ex. manifest / compose / helm / design)")
	importCmd.Flags().StringVarP(&name, "name", "n", "", "Name for the design file")
}
