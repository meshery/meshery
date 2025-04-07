package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/mod/semver"
)

var initModelCmd = &cobra.Command{
	Use:   "init",
	Short: "Generates scaffolding for convenient model creation",
	Long: `Generates a folder structure and guides user on model creation
Documentation for exp models init can be found at https://docs.meshery.io/reference/mesheryctl/exp/model/init`,
	Example: `
// generates a folder structure
mesheryctl exp model init [model-name]

// generates a folder structure and sets up model version
mesheryctl exp model init [model-name] --version [version] (default is v0.1.0)

// generates a folder structure under specified path
mesheryctl exp model init [model-name] --path [path-to-location] (default is current folder)

// generate a folder structure in json format
mesheryctl exp model init [model-name] --output-format [json|yaml|csv] (default is json)
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		{
			// validate model name
			if len(args) != 1 {
				return ErrModelInitFromString("must provide only one argument: model name")
			}
			modelName := args[0]
			input := map[string]any{"name": modelName}
			schema, err := initModelReadTemplate(initModelModelSchema)
			if err != nil {
				return ErrModelInit(err)
			}
			if err := initModelValidateDataOverSchema(schema, input); err != nil {
				return ErrModelInit(
					fmt.Errorf("invalid model name: %v", err),
				)
			}
		}

		{
			// validate format
			format, _ := cmd.Flags().GetString("output-format")
			getValidOutputFormatSlices := initModelGetValidOutputFormat()
			if !slices.Contains(getValidOutputFormatSlices, format) {
				validFormatsString := strings.Join(initModelGetValidOutputFormat(), ", ")
				return ErrModelUnsupportedOutputFormat(
					fmt.Sprintf(
						"[ %s ] are the only format supported",
						validFormatsString,
					),
				)
			}
		}

		{
			// validate version
			version, _ := cmd.Flags().GetString("version")
			if !semver.IsValid(version) {
				return ErrModelUnsupportedVersion(
					"version must follow a semver format, f.e. v1.2.3",
				)
			}
		}

		return nil

	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return ErrModelInit(err)
		}

		// validation done above that args contains exactly one argument
		modelName := args[0]
		path, _ := cmd.Flags().GetString("path")
		// immediately remove trailing folder separator
		path = strings.TrimRight(path, string(os.PathSeparator))
		version, _ := cmd.Flags().GetString("version")
		outputFormat, _ := cmd.Flags().GetString("output-format")

		utils.Log.Infof("Creating new Meshery model: %s", modelName)

		modelFolder := filepath.Join(path, modelName)
		modelVersionFolder := filepath.Join(modelFolder, version)

		{
			// if model/version folder already exists return with error
			info, err := os.Stat(modelVersionFolder)
			if !os.IsNotExist(err) && info.IsDir() {
				return ErrModelInitFromString(
					fmt.Sprintf(
						"folder %s exists, please specify different model name or version",
						modelVersionFolder,
					),
				)
			}
		}

		infoOnModelFolder, err := os.Stat(modelFolder)
		// this indicates if model folder already exists (before we created it)
		// this information will be used for clean up
		isModelFolderAlreadyExists := !os.IsNotExist(err) && infoOnModelFolder.IsDir()

		utils.Log.Infof("Creating directory structure...")
		err = os.MkdirAll(modelVersionFolder, initModelDirPerm)
		if err != nil {
			return ErrModelInit(err)
		}

		// we need to wrap the next code block in the function which returns error
		// because it has several return statements
		// and we want to capture all of them and perform action on failure.
		err = func() error {
			for _, item := range initModelData {
				item.beforeHook()
				itemFolderPath := modelVersionFolder
				if item.folderPath != "" {
					itemFolderPath = filepath.Join(modelVersionFolder, item.folderPath)
					err = os.MkdirAll(itemFolderPath, initModelDirPerm)
					if err != nil {
						return ErrModelInit(err)
					}
				}
				for name, templatePath := range item.files {
					content, err := getTemplateInOutputFormat(templatePath, outputFormat)
					if err != nil {
						return ErrModelInit(err)
					}
					filePath := filepath.Join(

						itemFolderPath,
						strings.Join(
							[]string{name, outputFormat},
							".",
						),
					)
					file, err := os.Create(filePath)
					if err != nil {
						return ErrModelInit(err)
					}
					defer file.Close() // Ensure the file is closed when the function exits

					_, err = file.Write(content)
					if err != nil {
						return ErrModelInit(err)
					}
				}
			}
			utils.Log.Infof("Created %s model at %s", modelName, modelFolder)
			utils.Log.Info("")
			utils.Log.Info(
				initModelReplacePlaceholders(
					initModelNextStepsText,
					map[string]string{
						"{modelName}":          modelName,
						"{modelVersion}":       version,
						"{modelFolder}":        modelFolder,
						"{outputFormat}":       outputFormat,
						"{modelVersionFolder}": modelVersionFolder,
					},
				),
			)
			return nil
		}()

		if err != nil {
			utils.Log.Info("Failure, cleaning up...")
			if !isModelFolderAlreadyExists {
				// if model folder didn'tv exist before -> delete it
				utils.Log.Infof("Removing %s", modelFolder)
				os.RemoveAll(modelFolder)
			} else {
				// otherwise remove only version folder
				utils.Log.Infof("Removing %s", modelVersionFolder)
				os.RemoveAll(modelVersionFolder)
			}
		}

		// TODO put a model name into generated model file
		return nil
	},
}

func init() {
	initModelCmd.Flags().StringP("path", "p", ".", "(optional) target directory (default: current dir)")
	initModelCmd.Flags().StringP("version", "", "v0.1.0", "(optional) model version (default: v0.1.0)")
	initModelCmd.Flags().StringP("output-format", "o", "json", "(optional) format to display in [json|yaml]")
}

func initModelGetValidOutputFormat() []string {
	return []string{"json", "yaml"}
	// TODO implement csv
	// return []string{"json", "yaml", "csv"}
}

const initModelDirPerm = 0755
const initModelModelSchema = "schemas/constructs/v1beta1/model/model.json"
const initModelTemplatePathModelJSON = "json_models/constructs/v1beta1/model.json"
const initModelTemplatePathComponentJSON = "json_models/constructs/v1beta1/component.json"
const initModelTemplatePathConnectionJSON = "json_models/constructs/v1beta1/connection.json"
const initModelTemplatePathRelathionshipJSON = "json_models/constructs/v1alpha3/relationship.json"

// TODO
// if csv output is not directory based
// should it have different text for csv output format?
const initModelNextStepsText = `Next steps:
1. cd {modelFolder}
2. Edit model.{outputFormat} to customize your model configuration
3. Add your components in the components/ directory
4. Define relationships in relationships/ directory
5. Add your connections in the connections/ directory
6. Define credentials in credentials/ directory
7. Use 'mesheryctl model build' to package your model

To import this model into Meshery:
$ mesheryctl model import {modelFolder}

To export this model as OCI image:
$ mesheryctl model build {modelVersionFolder} -t myregistry/{modelName}:{modelVersion}

Detailed guide: https://docs.meshery.io/guides/creating-new-model-with-mesheryctl`

// TODO
// initModelData fits well for json and yaml format
// if csv output is different (non folder based), will initModelData fit it?
var initModelData = []struct {
	folderPath string
	files      map[string]string
	beforeHook func()
}{
	{
		folderPath: "",
		// map file name to template key
		files: map[string]string{
			"model": initModelTemplatePathModelJSON,
		},
		beforeHook: func() {
			utils.Log.Info("Generating model definition...")
		},
	},
	{
		folderPath: "components",
		// map file name to template key
		files: map[string]string{
			"component": initModelTemplatePathComponentJSON,
		},
		beforeHook: func() {
			utils.Log.Info("Adding sample components...")
		},
	},
	{
		folderPath: "relationships",
		// map file name to template key
		files: map[string]string{
			"relationship": initModelTemplatePathRelathionshipJSON,
		},
		beforeHook: func() {
			utils.Log.Info("Creating sample relationships...")
		},
	},
	{
		folderPath: "connections",
		// map file name to template key
		files: map[string]string{
			"connection": initModelTemplatePathConnectionJSON,
		},
		beforeHook: func() {
			utils.Log.Info("Adding sample connections...")
		},
	},
	{
		folderPath: "credentials",
		// map file name to template key
		files: nil,
		beforeHook: func() {
			utils.Log.Info("Creating sample credentials...")
		},
	},
}

func initModelReadTemplate(templatePath string) ([]byte, error) {
	return schemas.Schemas.ReadFile(templatePath)
}

func getTemplateInOutputFormat(templatePath string, outputFormat string) ([]byte, error) {
	// outputFormat was already validated, it is one of the initModelGetValidOutputFormat()
	content, err := initModelReadTemplate(templatePath)
	if err != nil {
		return nil, err
	}

	if outputFormat == "json" {
		return content, err
	}

	if outputFormat == "yaml" {
		// TODO
		// this does not prevent the order of the fields
		// making fields to be ordered in alphabetical order;
		// need to generate yaml from schema instead
		yamlContent, err := yaml.JSONToYAML(content)
		if err != nil {
			return nil, utils.ErrJSONToYAML(err)
		}
		return yamlContent, nil
	}

	if outputFormat == "csv" {
		// impossible to reach here, as outputFormat is validated in prerun
		return nil, ErrModelUnsupportedOutputFormat("TODO implement csv")
	}

	// impossible to reach here, as outputFormat is validated in prerun
	return nil, ErrModelUnsupportedOutputFormat("unsupported output format")
}

func initModelReplacePlaceholders(input string, replacements map[string]string) string {
	for key, value := range replacements {
		input = strings.ReplaceAll(input, key, value)
	}
	return input
}

// TODO move it to more general package
func initModelValidateDataOverSchema(schema []byte, data map[string]interface{}) error {
	// version with full schema validation, like gojsonschema.Validate(schemaLoader, dataLoader),  is not working now,
	// probably because of the external references and relative paths in schema
	// it returns
	// Error: Could not read schema from HTTP, response status is 404 Not Found
	// --
	// as we need only to validate the model name
	// below we extract pattern for name property from schema json and validate manually string over pattern.
	// TODO figure out how to do a proper validation over schema.

	validationErrors := []string{}
	for property, value := range data {
		stringValue, ok := (value).(string)
		if !ok {
			continue
		}
		pattern, _ := initModelGetPatternFromSchema(schema, property)
		if pattern == "" {
			// skip if not possible to extract pattern
			continue
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			// skip on invalid regexp
			continue
		}
		if !re.MatchString(stringValue) {
			validationErrors = append(
				validationErrors,
				fmt.Sprintf("%s must match pattern %s", property, pattern),
			)
		}
	}

	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, ""))
	}

	return nil
}

func initModelGetPatternFromSchema(schema []byte, property string) (string, error) {
	// Generic structure to decode JSON
	var schemaMap map[string]interface{}

	// Unmarshal JSON schema into a map
	if err := json.Unmarshal(schema, &schemaMap); err != nil {
		return "", err
	}

	// Navigate to "properties" -> propertyName -> "pattern"
	if properties, ok := schemaMap["properties"].(map[string]interface{}); ok {
		if propSchema, ok := properties[property].(map[string]interface{}); ok {
			if pattern, ok := propSchema["pattern"].(string); ok {
				return pattern, nil
			}
		}
	}

	return "", fmt.Errorf("pattern not found for property: %s", property)
}
