// Copyright Meshery Author
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

package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/layer5io/meshkit/models/meshmodel/registry"
	extract "github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/walker"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var importModelCmd = &cobra.Command{
	Use:   "import",
	Short: "import model",
	Long:  "import a model from the registry",
	Example: `
// Import a model
mesheryctl exp model import <URI>

// Import a model from a local directory
mesheryctl exp model import <path to model>
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return system.ErrGetCurrentContext(err)
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp model import <URI/path to Model>\nRun 'mesheryctl exp model import --help' to see detailed help message"
		if len(args) == 0 {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if validURL := govalidator.IsURL(args[0]); !validURL {
			return importFromLocalFile(args[0])
		} else if isGithubURL := strings.Contains(args[0], "github.com"); isGithubURL {
			return importFromGithubURL(args[0])
		}
		return nil
	},
}

func importFromLocalFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return utils.ErrFileRead(errors.Errorf("File path does not exist. Path provided: %s", filePath))
	}

	if strings.HasSuffix(filePath, ".zip") {
		return processZip(filePath)
	}

	if strings.HasSuffix(filePath, ".tar.gz") || strings.HasSuffix(filePath, ".tgz") {
		return processTarGz(filePath)
	}

	dirPath := filePath
	// fmt.Println("Directory path:", dirPath)
	return processDirectory(dirPath)
}

func processZip(zipFilePath string) error {
	homeDir, _ := os.UserHomeDir()
	tmpDir := homeDir + "/.meshery/temp"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return system.ErrCreateDir(err, tmpDir)
	}
	defer os.RemoveAll(tmpDir)

	if err := extract.ExtractZip(tmpDir, zipFilePath); err != nil {
		return utils.ErrExtractFile(err, "zip")
	}

	return processDirectory(tmpDir)
}

func processTarGz(tarGzFilePath string) error {
	homeDir, _ := os.UserHomeDir()
	tmpDir := homeDir + "/.meshery/temp"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return system.ErrCreateDir(err, tmpDir)
	}
	defer os.RemoveAll(tmpDir)

	if err := extract.ExtractTarGz(tmpDir, tarGzFilePath); err != nil {
		return utils.ErrExtractFile(err, "tar.gz")
	}

	return processDirectory(tmpDir)
}

func processDirectory(dirPath string) error {
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return utils.ErrDirRead(err, dirPath)
	}
	imported := true 
	for _, entry := range entries {
		if entry.IsDir() {
			subdirPath := filepath.Join(dirPath, entry.Name())
			err := processDirectory(subdirPath)
			if err != nil {
				continue
			}
		} else {
			fileContent, err := ioutil.ReadFile(filepath.Join(dirPath, entry.Name()))
			if err != nil {
				continue
			}
			if err := registerComponent(entry.Name(), fileContent); err != nil {
				// utils.Log.Info("Error registering component: ", entry.Name(), err)
				imported = false
				break
			}
		}
	}
	if imported {
		utils.Log.Info("Model imported successfully")
	} else {
		utils.Log.Info("Model not imported")
	}
	return nil
}

func importFromGithubURL(githubURL string) error {
	parsedURL, err := url.Parse(githubURL)
	if err != nil {
		return utils.ErrParsingUrl(err)
	}
	pathParts := strings.Split(parsedURL.Path, "/")

	githubInfo := walker.NewGithub()
	githubInfo.Owner(pathParts[1])
	githubInfo.Repo(pathParts[2])
	githubInfo.Branch(pathParts[4])
	githubInfo.Root(strings.Join(pathParts[5:], "/"))

	FileInterceptor := func(file walker.GithubContentAPI) error {
		content, err := base64.StdEncoding.DecodeString(file.Content)
		if err != nil {
			return err
		}
		if strings.HasSuffix(githubURL, ".zip") {
			return processZipFromGithub(content)
		}
		return registerComponent(file.Name, content)
	}

	githubInfo.RegisterFileInterceptor(FileInterceptor)

	if err := githubInfo.Walk(); err != nil {
		utils.Log.Info("Error occurred during traversal:", err)
	}
	return nil
}

func processZipFromGithub(zipContent []byte) error {
	homeDir, _ := os.UserHomeDir()
	tmpDir := homeDir + "/.meshery/tmp"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return system.ErrCreateDir(err, tmpDir)
	}
	defer os.RemoveAll(tmpDir)

	zipFilePath := filepath.Join(tmpDir, "temp.zip")
	if err := ioutil.WriteFile(zipFilePath, zipContent, 0644); err != nil {
		return err
	}

	return processZip(zipFilePath)
}

func registerComponent(fileName string, content []byte) error {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return utils.ErrLoadConfig(err)
	}
	baseURL := mctlCfg.GetBaseMesheryURL()

	registrantData := &registry.MeshModelRegistrantData{
		Host:       v1beta1.Host{Hostname: "localhost"},
		EntityType: "component",
		Entity:     content,
	}

	url := fmt.Sprintf("%s/api/meshmodel/components/register", baseURL)

	requestBody, err := json.Marshal(registrantData)
	if err != nil {
		return utils.ErrMarshal(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return utils.ErrCreatingRequest(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// For debugging, checking which files are being imported
		// utils.Log.Info("Model imported successfully: ", fileName)
	} else {
		// utils.Log.Info("Error registering component: ", fileName, resp.Status)
		return utils.ErrImportingModel(errors.New("Error registering component"))
	}

	return nil
}
