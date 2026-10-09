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
	"os"
	"path"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/patterns"
	pattern "github.com/meshery/schemas/models/v1beta3/design"

	"github.com/meshery/schemas/models/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var name string
var sourceType string

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
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/import"

		// If pattern file is passed via flags
		if sourceType != "" {
			validSourceTypes, err := getDesignSourceTypes()
			if err != nil {
				return err
			}
			if sourceType, err = retrieveProvidedSourceType(sourceType, validSourceTypes); err != nil {
				return err
			}
		}

		switch sourceType {
		case "Helm Chart":
			sourceType = string(core.HelmChart)
		case "Kubernetes Manifest":
			sourceType = string(core.K8sManifest)
		case "Meshery Design":
			sourceType = string(core.MesheryDesign)
		case "Docker Compose":
			sourceType = string(core.DockerCompose)
		}

		pattern, err := importPattern(sourceType, file, patternURL)
		if err != nil {
			return err
		}

		utils.Log.Infof("The design file '%s' has been imported. Design ID: %s", pattern.Name, utils.TruncateID(pattern.ID.String()))

		return nil
	},
}

// marshalDesignImportBody builds the `POST /api/pattern/import` request body
// through the schemas-generated union type rather than a hand-written
// map[string]interface{}. The wire field names - notably camelCase `fileName` -
// then come from the generated struct tags, so they cannot drift from the
// contract: sending the legacy snake_case `file_name` leaves the server's oneOf
// unmatched and the request is rejected with "Invalid design import request"
// (meshery-server-1422), the bug this endpoint regressed with before.
//
// `variant` selects the oneOf arm via the generated
// FromMesheryPatternImport{File,URL}Payload builders.
func marshalDesignImportBody(variant func(*pattern.MesheryPatternImportRequestBody) error) ([]byte, error) {
	var body pattern.MesheryPatternImportRequestBody
	if err := variant(&body); err != nil {
		return nil, utils.ErrMarshal(err)
	}

	jsonValues, err := json.Marshal(body)
	if err != nil {
		return nil, utils.ErrMarshal(err)
	}
	return jsonValues, nil
}

func importPattern(sourceType string, file string, patternURL string) (*models.MesheryPattern, error) {
	// If design name is not provided
	// use file name as default
	var patternName string
	fileName := path.Base(file)
	if name == "" {
		patternName = fileName
	} else {
		patternName = name
	}

	// Check if the design is a file or a URL and build the matching oneOf arm.
	if !utils.IsValidUrl(file) {
		content, err := os.ReadFile(file)
		if err != nil {
			return nil, utils.ErrFileRead(err)
		}

		jsonValues, err := marshalDesignImportBody(func(body *pattern.MesheryPatternImportRequestBody) error {
			return body.FromMesheryPatternImportFilePayload(pattern.MesheryPatternImportFilePayload{
				Name:     &patternName,
				File:     content,
				FileName: fileName,
			})
		})
		if err != nil {
			return nil, err
		}

		imported, err := postDesignImport(patternURL, jsonValues)
		if err != nil {
			return nil, err
		}
		utils.Log.Debug("design file saved")
		return imported, nil
	}

	jsonValues, err := marshalDesignImportBody(func(body *pattern.MesheryPatternImportRequestBody) error {
		return body.FromMesheryPatternImportURLPayload(pattern.MesheryPatternImportURLPayload{
			Name: &patternName,
			Url:  file,
		})
	})
	if err != nil {
		return nil, err
	}

	imported, err := postDesignImport(patternURL, jsonValues)
	if err != nil {
		return nil, err
	}
	utils.Log.Debug("Fetched the design from the remote host")
	return imported, nil
}

// postDesignImport POSTs an already-marshalled import body and returns the
// first design from the response collection. Both oneOf arms share it so the
// response handling cannot drift between them - the divergence that let the
// snake_case `file_name` bug survive in only one arm.
func postDesignImport(patternURL string, jsonValues []byte) (*models.MesheryPattern, error) {
	req, err := utils.NewRequest("POST", patternURL, bytes.NewBuffer(jsonValues))
	if err != nil {
		return nil, utils.ErrCreatingRequest(err)
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return nil, utils.ErrRequestResponse(err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Log.Debug("failed to read response body")
		return nil, utils.ErrReadResponseBody(err)
	}

	var response []*models.MesheryPattern
	if err := json.Unmarshal(body, &response); err != nil {
		utils.Log.Debug("failed to unmarshal JSON response")
		return nil, utils.ErrUnmarshal(err)
	}

	// The server returns the saved design(s) as a collection. Guard the index
	// so an empty (or null) collection surfaces as a structured error instead
	// of panicking with an out-of-range index.
	if len(response) == 0 || response[0] == nil {
		return nil, ErrDesignInvalidApiResponse("design import returned no design")
	}

	return response[0], nil
}

func init() {
	importCmd.Flags().StringVarP(&file, "file", "f", "", "Path/URL to design file")
	importCmd.Flags().StringVarP(&sourceType, "source-type", "s", "", "Type of source file (ex. manifest / compose / helm / design)")
	importCmd.Flags().StringVarP(&name, "name", "n", "", "Name for the design file")
}

func readPatternFromFile(filePath string) (*pattern.PatternFile, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, utils.ErrFileRead(err)
	}

	patternFile, err := patterns.GetPatternFormat(string(fileContent))
	if err != nil {
		return nil, err
	}

	return patternFile, nil
}
