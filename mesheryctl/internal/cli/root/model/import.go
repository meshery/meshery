package model

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"mime/multipart"
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
		//if oci then validate and send
		tarData, err := compressDirectory(path)
		if err != nil {
			return err
		}
		fileName := filepath.Base(path) + ".tar.gz"
		err = sendToAPI(tarData, fileName, "dir")
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
func sendToAPI(data []byte, name string, dataType string) error {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return err
	}

	baseURL := mctlCfg.GetBaseMesheryURL()
	url := baseURL + "/api/meshmodels/register"
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	var formFile io.Writer
	if dataType == "oci" {
		formFile, _ = writer.CreateFormField("oci")
	} else {
		formFile, _ = writer.CreateFormField("dir")
	}
	_, err = formFile.Write(data)
	if err != nil {
		err = meshkitutils.ErrWriteFile(err, name)
		return err
	}

	_ = writer.Close()

	req, err := utils.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
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
	var response models.RegisterMeshmodelAPIResponse

	if err := meshkitutils.Unmarshal(string(bodyBytes), &response); err != nil {
		err = models.ErrUnmarshal(err, "response body")
		return err
	}

	utils.Log.Info("\n\033[1m", response.ErrMsg, "\033[0m\n")
	if len(response.EntityTypeSummary.SuccessfulComponents) > 0 {
		utils.Log.Infof("\033[1mTotal Component(s) Imported: %d\033[0m\n", response.EntityCount.CompCount)
		fmt.Printf("\033[1m%-30s %-30s %-10s\033[0m\n", centerText("DisplayName", 30), centerText("ModelName", 30), centerText("Version", 10))

		for _, comp := range response.EntityTypeSummary.SuccessfulComponents {
			displayName, _ := comp["DisplayName"].(string)
			modelData, _ := comp["Model"].(map[string]interface{})
			modelDisplayName, _ := modelData["displayName"].(string)
			modelVersion, _ := modelData["model"].(map[string]interface{})["version"].(string)

			displayNameLines := wrapText(displayName, 30)
			modelNameLines := wrapText(modelDisplayName, 30)
			versionLines := wrapText(modelVersion, 10)

			maxLines := max(len(displayNameLines), len(modelNameLines))

			for i := 0; i < maxLines; i++ {
				fmt.Printf(
					"%-30s %-30s %-10s\n",
					centerText(getLine(displayNameLines, i), 30),
					centerText(getLine(modelNameLines, i), 30),
					centerText(getLine(versionLines, i), 10))
			}
		}
	}

	if len(response.EntityTypeSummary.SuccessfulRelationships) > 0 {
		utils.Log.Infof("\nTotal Relationship(s) Imported: %d", (response.EntityCount.RelCount))
		//fmt.Printf("\033[1m%-20s %-20s %-20s\033[0m\n", centerText("Type", 20), centerText("SubType", 20), centerText("Kind", 20))
		fmt.Printf("\033[1m %-20s %-20s\033[0m\n", centerText("SubType", 20), centerText("Kind", 20))

		for _, rel := range response.EntityTypeSummary.SuccessfulRelationships {

			//future when we support relationship type

			// fmt.Printf("%-20s %-20s %-20s\n",
			// 	centerText(rel["RelationshipType"].(string), 20),
			// 	centerText(rel["Subtype"].(string), 20),
			// 	centerText(rel["Kind"].(string), 20))

			fmt.Printf("%-20s %-20s\n ",
				centerText(rel["Subtype"].(string), 20),
				centerText(rel["Kind"].(string), 20))
		}
	}

	modelNames := ModelNames(&response)
	if len(response.EntityTypeSummary.UnsuccessfulEntityNameWithError) > 0 {
		utils.Log.Info("\033[1m\nImport failed for these Entity(s): \033[0m")
		for _, entity := range response.EntityTypeSummary.UnsuccessfulEntityNameWithError {
			for name, errInterface := range entity {
				errMap, ok := errInterface.(map[string]interface{})
				if !ok {
					utils.Log.Infof("Error: unable to assert errInterface as map[string]interface{}")
					return nil
				}
				longDescription, ok := errMap["LongDescription"].([]interface{})
				if !ok {
					utils.Log.Infof("Error: unable to assert LongDescription as []interface{}")
					return nil
				}

				// Convert []interface{} to []string
				var longDescriptionStrings []string
				for _, desc := range longDescription {
					if strDesc, ok := desc.(string); ok {
						longDescriptionStrings = append(longDescriptionStrings, strDesc)
					}
				}

				// Join the strings with a separator, e.g., space or newline
				longDescriptionText := strings.Join(longDescriptionStrings, " ")

				utils.Log.Infof("\nEntity File Name: \033[1m%s\033[0m and error: \033[1m%s\033[0m", name, longDescriptionText)
			}
		}
	}

	utils.Log.Infof("\033[1m\n%s model(s) imported\033[0m", modelNames)
	return nil
}
func ModelNames(response *models.RegisterMeshmodelAPIResponse) string {
	msg := ""
	seen := make(map[string]bool) // map to track seen model names

	for _, model := range response.ModelName {
		if !seen[model] {
			msg += model + " "
			seen[model] = true
		}
	}

	return msg
}
func wrapText(text string, width int) []string {
	var wrapped []string
	for len(text) > width {
		spaceIndex := strings.LastIndex(text[:width], " ")
		if spaceIndex == -1 {
			spaceIndex = width
		}
		wrapped = append(wrapped, text[:spaceIndex])
		text = strings.TrimSpace(text[spaceIndex:])
	}
	wrapped = append(wrapped, text)
	return wrapped
}

func getLine(lines []string, index int) string {
	if index < len(lines) {
		return lines[index]
	}
	return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text
	}
	padding := (width - len(text)) / 2
	return fmt.Sprintf("%*s%*s", padding+len(text), text, width-len(text)-padding, "")
}
