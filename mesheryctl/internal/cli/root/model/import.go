package model

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/models"
	meshkitutils "github.com/layer5io/meshkit/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ImportRequestBody struct {
	ImportBody struct {
		ModelFile []byte `json:"model_file"`
		URL       string `json:"url,omitempty"`
		FileName  string `json:"file_name,omitempty"`
	} `json:"importBody"`
	UploadType string `json:"uploadType"`
}

var importModelCmd = &cobra.Command{
	Use:   "import",
	Short: "import models from mesheryctl command",
	Long:  "import model by specifying the directory, file. Use 'import model [filepath]' or 'import model  [directory]'.",
	Example: `
	import model  /path/to/[file.yaml|file.json]
	import model  /path/to/models
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check prerequisites
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl model import [ file | filePath ]\nRun 'mesheryctl model import --help' to see detailed help message"
		if len(args) == 0 {
			return fmt.Errorf("[ file | filepath ] isn't specified\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("could not access the specified path: %v", err)
		}

		var tarData []byte
		var fileName string

		if info.IsDir() {
			// If the input is a directory, compress it
			tarData, err = compressDirectory(path)
			if err != nil {
				return err
			}
			fileName = filepath.Base(path) + ".tar.gz"
		} else {
			// If the input is a file, read its contents
			fileData, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("could not read the specified file: %v", err)
			}
			tarData = fileData
			fileName = filepath.Base(path)
		}

		err = registerModel(tarData, fileName, "file")
		if err != nil {
			utils.Log.Error(err)
			return err
		}
		return nil
	},
}

func compressDirectory(dirpath string) ([]byte, error) {
	tw := meshkitutils.NewTarWriter()
	defer tw.Close()

	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return meshkitutils.ErrFileWalkDir(err, path)
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return handlers.ErrOpenFile(path)
		}
		defer file.Close()

		fileData, err := io.ReadAll(file)
		if err != nil {
			return meshkitutils.ErrReadFile(err, path)
		}

		relPath, err := filepath.Rel(filepath.Dir(dirpath), path)
		if err != nil {
			return meshkitutils.ErrRelPath(err, path)
		}

		if err := tw.Compress(relPath, fileData); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err = io.Copy(gzipWriter, tw.Buffer)
	if err != nil {
		return nil, meshkitutils.ErrCopyFile(err)
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, meshkitutils.ErrCloseFile(err)
	}

	return buf.Bytes(), nil
}
func registerModel(data []byte, name string, dataType string) error {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return err
	}

	baseURL := mctlCfg.GetBaseMesheryURL()
	url := baseURL + "/api/meshmodels/register"
	importRequest := ImportRequestBody{
		UploadType: dataType,
	}
	importRequest.ImportBody.ModelFile = data
	importRequest.ImportBody.FileName = name

	// Marshal the request body to JSON
	requestBody, err := json.Marshal(importRequest)
	if err != nil {
		return err
	}

	// Create the HTTP request
	req, err := utils.NewRequest(http.MethodPost, url, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the HTTP request
	resp, err := utils.MakeRequest(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {

		err = models.ErrDoRequest(err, resp.Request.Method, url)
		return err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err = models.ErrDataRead(err, "response body")
		return err
	}
	var response models.RegistryAPIResponse

	if err := meshkitutils.Unmarshal(string(bodyBytes), &response); err != nil {
		err = models.ErrUnmarshal(err, "response body")
		return err
	}
	utils.Log.Info(response.ErrMsg)

	if len(response.EntityTypeSummary.SuccessfulComponents) > 0 {
		header := []string{"Component Name", "Model Name", "Version"}
		rows := [][]string{}
		utils.Log.Info("\n", utils.BoldString("Imported Component(s)"), "\n")

		for _, comp := range response.EntityTypeSummary.SuccessfulComponents {
			displayName, _ := comp["DisplayName"].(string)
			modelData, _ := comp["Model"].(map[string]interface{})
			modelDisplayName, _ := modelData["displayName"].(string)
			modelVersion, _ := modelData["model"].(map[string]interface{})["version"].(string)

			rows = append(rows, []string{displayName, modelDisplayName, modelVersion})
		}
		utils.PrintToTable(header, rows)

	}

	if len(response.EntityTypeSummary.SuccessfulRelationships) > 0 {
		header := []string{"From(CompoentName/ModelName)", "To(ComponentName/ModelNAme)", "Kind", "SubType"}
		rows := [][]string{}
		utils.Log.Info("\n", utils.BoldString("Imported Relationship(s)"), "\n")
		for _, rel := range response.EntityTypeSummary.SuccessfulRelationships {
			kind := rel["Kind"].(string)
			subtype := rel["Subtype"].(string)
			selectors := rel["Selectors"].([]interface{})
			for _, selector := range selectors {
				selectorMap := selector.(map[string]interface{})
				allow := selectorMap["allow"].(map[string]interface{})
				from := allow["from"].([]interface{})
				to := allow["to"].([]interface{})

				fromComponent := fmt.Sprintf("%s/%s", from[0].(map[string]interface{})["kind"], from[0].(map[string]interface{})["model"])
				toComponent := fmt.Sprintf("%s/%s", to[0].(map[string]interface{})["kind"], to[0].(map[string]interface{})["model"])
				rows = append(rows, []string{fromComponent, toComponent, kind, subtype})
			}

		}
		utils.PrintToTable(header, rows)
	}

	modelNames := ModelNames(&response)
	if len(response.EntityTypeSummary.UnsuccessfulEntityNameWithError) > 0 {
		utils.Log.Info("\n", utils.BoldString("Import failed for these Entity(s): "), "\n")
		for _, entity := range response.EntityTypeSummary.UnsuccessfulEntityNameWithError {
			entityMap, err := meshkitutils.Cast[map[string]interface{}](entity)
			if err != nil {
				utils.Log.Error(err)
				continue
			}

			names, err := meshkitutils.Cast[[]interface{}](entityMap["name"])
			if err != nil {
				utils.Log.Error(err)
				continue
			}

			entityTypes, err := meshkitutils.Cast[[]interface{}](entityMap["entityType"])
			if err != nil {
				utils.Log.Error(err)
				continue
			}

			errorDetails, err := meshkitutils.Cast[map[string]interface{}](entityMap["error"])
			if err != nil {
				utils.Log.Error(err)
				continue
			}

			longDescriptionInterface := errorDetails["LongDescription"]
			longDescriptionSlice, ok := longDescriptionInterface.([]interface{})
			if !ok {
				utils.Log.Infof("Type assertion to []interface{} failed for LongDescription: %v (type %T)", longDescriptionInterface, longDescriptionInterface)
				continue
			}

			var longDescription string
			for _, item := range longDescriptionSlice {
				str, ok := item.(string)
				if !ok {
					utils.Log.Infof("Item in LongDescription is not a string: %v (type %T)", item, item)
					continue
				}
				longDescription += str + " "
			}

			longDescription = strings.TrimSpace(longDescription)
			EntityTypeLine := ""
			for i, name := range names {
				entityType := ""
				if i < len(entityTypes) {
					entityType = entityTypes[i].(string)
				}

				if entityType == "" {
					// If entityType is not present, log normal message
					utils.Log.Infof("Entity Filename: %s and error: \n%s", name, longDescription)
				} else {
					if EntityTypeLine != "" {
						EntityTypeLine += ", "
					}
					EntityTypeLine = fmt.Sprintf("%s entity of type %s with model name %s", EntityTypeLine, utils.BoldString(entityType), utils.BoldString(name.(string)))
				}
			}
			if EntityTypeLine != "" {
				utils.Log.Infof("Import did not occur for%s and error: \n%s\n", EntityTypeLine, longDescription)
			}

		}
	}
	totalEntityCount := response.EntityCount.CompCount + response.EntityCount.RelCount
	if totalEntityCount != 0 {
		utils.Log.Infof("\n%s model(s) imported", utils.BoldString(modelNames))

	}
	return nil
}
func ModelNames(response *models.RegistryAPIResponse) string {
	var builder strings.Builder
	seen := make(map[string]bool) // map to track seen model names

	for _, model := range response.ModelName {
		if model != "" {
			if !seen[model] {
				if builder.Len() > 0 {
					builder.WriteString(", ")
				}
				builder.WriteString(model)
				seen[model] = true
			}
		}
	}
	return builder.String()
}
