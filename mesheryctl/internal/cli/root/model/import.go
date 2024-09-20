package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/encoding"
	meshkitutils "github.com/layer5io/meshkit/utils"
	schemav1beta1 "github.com/meshery/schemas/models/v1beta1"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	location     string
	templateFile string
	register     bool
)

var importModelCmd = &cobra.Command{
	Use:   "import",
	Short: "Import models from mesheryctl command",
	Long:  "Import models by specifying the directory, file, or URL. You can also provide a template JSON file and registrant name.",
	Example: `
	mesehryctl model import -f [ URI ]
	mesehryctl model import -f [ URI ] -t [ path to template file ] ( only required in case of URL )
	mesehryctl model import -f [ URI ] -t [ path to template file ] -r ( to skip registration by default registration is true)
 
	mesehryctl model import -f OCI 
	mesehryctl model import -f model.tar.gz 
	mesehryctl model import --f /path/to/models
    mesehryctl model import --f http://example.com/model -t /path/to/template.json 
	mesehryctl model import --f http://example.com/model -t /path/to/template.json -r
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl model import [ file | filePath | URL ]\nRun 'mesheryctl model import --help' to see detailed help message"
		if location == "" && len(args) == 0 {
			return fmt.Errorf("[ file | filepath | URL ] isn't specified\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		if location != "" {
			path = location
		} else {
			path = args[0]
		}

		validURL := strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "git://")
		if validURL {
			if templateFile == "" {
				return fmt.Errorf("a template file must be specified when importing from a URL")
			}

			fileData, err := os.ReadFile(templateFile)
			if err != nil {
				return fmt.Errorf("could not read the specified template file: %v", err)
			}
			err = registerModel(fileData, "", "url", path, !register)
			if err != nil {
				utils.Log.Error(err)
				return nil
			}
			locationForModel := utils.MesheryFolder + "/models"
			utils.Log.Info("Model can be accessed from ", locationForModel)
			return nil
		}

		// Handle the case when the path is a file or directory
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("could not access the specified path: %v", err)
		}

		var tarData []byte
		var fileName string

		if info.IsDir() {
			var buf bytes.Buffer
			err = meshkitutils.Compress(path, &buf)
			if err != nil {
				return err
			}
			tarData = buf.Bytes()
			fileName = filepath.Base(path) + ".tar.gz"
		} else {
			fileData, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("could not read the specified file: %v", err)
			}
			tarData = fileData
			fileName = filepath.Base(path)
		}

		err = registerModel(tarData, fileName, "file", "", !register)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		return nil
	},
}

func registerModel(data []byte, filename string, dataType string, sourceURI string, register bool) error {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return err
	}

	baseURL := mctlCfg.GetBaseMesheryURL()
	url := baseURL + "/api/meshmodels/register"
	var importRequest schemav1beta1.ImportRequest
	importRequest.UploadType = dataType
	if dataType == "file" {
		importRequest.ImportBody.ModelFile = data
	} else {
		err = encoding.Unmarshal(data, &importRequest.ImportBody.Model)
		if err != nil {
			return err
		}
	}
	importRequest.ImportBody.Url = sourceURI
	importRequest.ImportBody.FileName = filename
	importRequest.Register = register
	requestBody, err := json.Marshal(importRequest)
	if err != nil {
		return err
	}

	req, err := utils.NewRequest(http.MethodPost, url, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

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

	if err := encoding.Unmarshal((bodyBytes), &response); err != nil {
		err = models.ErrUnmarshal(err, "response body")
		return err
	}
	displayEntities(&response)
	return nil
}

func displayEntities(response *models.RegistryAPIResponse) {
	displaySummary(response)
	ok := displayEmtpyModel(response)
	if !ok {
		return
	}
	displayEntitisIfModel(response)
}
func displayEmtpyModel(response *models.RegistryAPIResponse) bool {
	if len(response.ModelName) != 0 && response.EntityCount.CompCount == 0 && response.EntityCount.RelCount == 0 {
		if response.EntityCount.TotalErrCount == 0 {
			return false
		}
	}
	return true
}

// TO check the case if we were never able to read the file at first palce
func hasExtension(name string) bool {
	extension := filepath.Ext(name)
	return extension == ".json" || extension == ".yaml" || extension == ".yml" || extension == ".tar.gz" || extension == ".tar" || extension == ".zip" || extension == ".tgz"
}

func displayEntitisIfModel(response *models.RegistryAPIResponse) {
	var modelsWithoutExtension []string
	var modelsWithExtension []string

	// Separate models into those with and without extensions
	for _, model := range response.ModelName {
		if model != "" {
			if hasExtension(model) {
				modelsWithExtension = append(modelsWithExtension, model)
			} else {
				modelsWithoutExtension = append(modelsWithoutExtension, model)
			}
		}
	}

	// Function to display models and their components, relationships, and entities
	displayModelInfo := func(model string, hasExtension bool) {
		if !hasExtension {
			boldModel := utils.BoldString("MODEL")
			utils.Log.Infof("\n%s: %s", boldModel, model)
		}
		displaySuccessfulComponents(response, model)
		displaySuccessfulRelationships(response, model)
		displayUnsuccessfulEntities(response, model)
	}

	for _, model := range modelsWithoutExtension {
		displayModelInfo(model, false)
	}
	for _, model := range modelsWithExtension {
		displayModelInfo(model, true)
	}
}

func displaySuccessfulComponents(response *models.RegistryAPIResponse, modelName string) {
	if len(response.EntityTypeSummary.SuccessfulComponents) > 0 {
		header := []string{"Component", "Category", "Version"}
		rows := [][]string{}

		for _, comp := range response.EntityTypeSummary.SuccessfulComponents {
			displayName, _ := comp["DisplayName"].(string)
			modelDisplayName, _ := comp["Model"].(string)
			category, _ := comp["Category"].(string)
			modelVersion, _ := comp["Version"].(string)
			if modelDisplayName == modelName {
				rows = append(rows, []string{displayName, category, modelVersion})

			}
		}
		if len(rows) > 0 {
			fmt.Println("")
			utils.PrintToTable(header, rows)
		}
	}
}
func displaySuccessfulRelationships(response *models.RegistryAPIResponse, model string) {
	if len(response.EntityTypeSummary.SuccessfulRelationships) > 0 {
		header := []string{"From", "To"}
		seen := make(map[string]bool)
		relationshipMap := make(map[string][][]string)

		for _, rel := range response.EntityTypeSummary.SuccessfulRelationships {
			kind := rel["Kind"].(string)
			subtype := rel["Subtype"].(string)
			relationshipType := rel["RelationshipType"].(string)
			modelName := rel["Model"].(string)
			if modelName != model {
				continue
			}
			selectors := rel["Selectors"].([]interface{})
			for _, selector := range selectors {
				selectorMap := selector.(map[string]interface{})
				allow := selectorMap["allow"].(map[string]interface{})
				from := allow["from"].([]interface{})
				to := allow["to"].([]interface{})
				fromComponent := fmt.Sprintf("%s", from[0].(map[string]interface{})["kind"])
				toComponent := fmt.Sprintf("%s", to[0].(map[string]interface{})["kind"])
				key := fmt.Sprintf("%s/%s/%s", kind, subtype, relationshipType)
				if seen[key+fromComponent+toComponent] {
					continue
				}
				seen[key+fromComponent+toComponent] = true
				relationshipMap[key] = append(relationshipMap[key], []string{fromComponent, toComponent})
			}
		}
		for key, rows := range relationshipMap {
			if len(rows) > 0 {
				fmt.Println("")
				boldRelationships := utils.BoldString("RELATIONSHIP:")
				if len(rows) > 1 {
					boldRelationships = utils.BoldString("RELATIONSHIPS:")
				}
				parts := strings.Split(key, "/")
				utils.Log.Infof("  %s Kind of %s, sub type %s and type %s", boldRelationships, parts[0], parts[1], parts[2])
				utils.PrintToTable(header, rows)
			}
		}
	}
}

func displayUnsuccessfulEntities(response *models.RegistryAPIResponse, modelName string) {
	if len(response.EntityTypeSummary.UnsuccessfulEntityNameWithError) > 0 {
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

			longDescription := buildLongDescription(errorDetails["LongDescription"])

			EntityTypeLine := buildEntityTypeLine(names, entityTypes, longDescription, modelName)
			if EntityTypeLine != "" {
				fmt.Println("")
				utils.Log.Infof("  %s: Import did not occur for%s error: \n  %s", utils.BoldString("ERROR"), EntityTypeLine, longDescription)
			}

		}
	}
}

func buildLongDescription(longDescriptionInterface interface{}) string {
	longDescriptionSlice, ok := longDescriptionInterface.([]interface{})
	if !ok {
		utils.Log.Infof("Type assertion to []interface{} failed for LongDescription: %v (type %T)", longDescriptionInterface, longDescriptionInterface)
		return ""
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

	return strings.TrimSpace(longDescription)
}

func buildEntityTypeLine(names, entityTypes []interface{}, longDescription, modelName string) string {
	compCount, relCount := 0, 0
	EntityTypeLine := ""
	for i, name := range names {
		entityType := ""
		if i < len(entityTypes) {
			entityType = entityTypes[i].(string)
		}
		if modelName != "" {
			if modelName != name.(string) {
				continue
			}
		} else if modelName == "" {
			if entityType != "Unknown" {
				continue
			}
		}
		if entityType == "unknown" {
			utils.Log.Infof("\n%s: Import process for file %s encountered error: \n    %s", utils.BoldString("ERROR"), name.(string), longDescription)
		} else if entityType == "component" {
			compCount++
		} else if entityType == "relationship" {
			relCount++
		}

	}
	if compCount > 0 {
		word := "entity"
		if compCount > 1 {
			word = "entities"
		}
		msg := fmt.Sprintf(" %d %s of type component", compCount, word)
		EntityTypeLine = msg
	}
	if compCount > 0 && relCount > 0 {
		EntityTypeLine = fmt.Sprintf("%s and", EntityTypeLine)
	}
	if relCount > 0 {
		word := "entity"
		if relCount > 1 {
			word = "entities"
		}
		msg := fmt.Sprintf(" %d %s of type relationship", relCount, word)
		EntityTypeLine = fmt.Sprintf("%s%s", EntityTypeLine, msg)
	}
	return EntityTypeLine
}

func displaySummary(response *models.RegistryAPIResponse) {
	if response.EntityCount.ModelCount == 0 {
		return
	}
	boldSummary := utils.BoldString("SUMMARY")
	utils.Log.Infof("%s: %s", boldSummary, response.ErrMsg)
}
func ModelNames(response *models.RegistryAPIResponse) string {
	var builder strings.Builder
	seen := make(map[string]bool)

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
func init() {
	importModelCmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})

	importModelCmd.Flags().StringVarP(&location, "file", "f", "", "Specify path to the file or directory")
	importModelCmd.Flags().StringVarP(&templateFile, "template", "t", "", "Specify path to the template JSON file")
	importModelCmd.Flags().BoolVarP(&register, "register", "r", false, "Skip registration of the model")

}
