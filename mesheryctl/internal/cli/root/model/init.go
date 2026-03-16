package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gopkg.in/yaml.v3"
)

type cmdModelInitFlags struct {
	Path         string `json:"path" validate:"omitempty,relabspath"`
	Version      string `json:"version" validate:"semver"`
	OutputFormat string `json:"output-format" validate:"oneof=json yaml"`
}

var modelInitFlags cmdModelInitFlags

var initModelCmd = &cobra.Command{
	Use:   "init",
	Short: "Generates scaffolding for convenient model creation",
	Long: `Generates a folder structure and guides user on model creation
Find more information at: https://docs.meshery.io/reference/mesheryctl/model/init`,
	Example: `
// generates a folder structure
mesheryctl model init [model-name]

// generates a folder structure and sets up model version
mesheryctl model init [model-name] --version [version] (default is v0.1.0)

// generates a folder structure under specified path
mesheryctl model init [model-name] --path [path-to-location] (default is current folder)

// generate a folder structure in json format
mesheryctl model init [model-name] --output-format [json|yaml|csv] (default is json)
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &modelInitFlags)
	},
	Args: func(cmd *cobra.Command, args []string) error {

		// validate model name
		if len(args) != 1 {
			return ErrModelInitFromString(errInitOneArg)
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

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// validation done above that args contains exactly one argument
		modelName := args[0]
		path := modelInitFlags.Path
		// immediately remove trailing folder separator
		path = strings.TrimRight(path, string(os.PathSeparator))
		// this will make it in one format
		path = filepath.Join(path)

		utils.Log.Infof("Creating new Meshery model: %s", modelName)

		modelFolder := filepath.Join(path, modelName)
		modelVersionFolder := filepath.Join(modelFolder, modelInitFlags.Version)

		{
			// if model/version folder already exists return with error
			info, err := os.Stat(modelVersionFolder)
			if !os.IsNotExist(err) && info.IsDir() {
				return ErrModelInitFromString(
					fmt.Sprintf(
						errInitFolderExists,
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
					content, err := getTemplateInOutputFormat(templatePath, modelInitFlags.OutputFormat)
					if err != nil {
						return ErrModelInit(err)
					}
					filePath := filepath.Join(

						itemFolderPath,
						strings.Join(
							[]string{name, modelInitFlags.OutputFormat},
							".",
						),
					)
					file, err := os.Create(filePath)
					if err != nil {
						return ErrModelInit(err)
					}
					defer func() { _ = file.Close() }() // Ensure the file is closed when the function exits

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
						"{path}":               path,
						"{modelName}":          modelName,
						"{modelVersion}":       modelInitFlags.Version,
						"{modelFolder}":        modelFolder,
						"{outputFormat}":       modelInitFlags.OutputFormat,
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
				_ = os.RemoveAll(modelFolder)
			} else {
				// otherwise remove only version folder
				utils.Log.Infof("Removing %s", modelVersionFolder)
				_ = os.RemoveAll(modelVersionFolder)
			}
			return err
		}

		// TODO put a model name into generated model file
		return nil
	},
}

func init() {
	initModelCmd.Flags().StringVarP(&modelInitFlags.Path, "path", "p", ".", "(optional) target directory (default: current dir)")
	initModelCmd.Flags().StringVarP(&modelInitFlags.Version, "version", "", "v0.1.0", "(optional) model version (default: v0.1.0)")
	initModelCmd.Flags().StringVarP(&modelInitFlags.OutputFormat, "output-format", "o", "json", "(optional) format to display in [json|yaml]")
}

const (
	initModelDirPerm                  = 0o755
	initModelModelSchema              = "schemas/constructs/v1beta1/model/model.yaml"
	initModelTemplatePathModel        = "schemas/constructs/v1beta1/model/templates/model_template"
	initModelTemplatePathComponent    = "schemas/constructs/v1beta1/component/templates/component_template"
	initModelTemplatePathRelationship = "schemas/constructs/v1alpha3/relationship/templates/relationship_template"
)

// TODO: Connection templates are temporarily disabled.
// This constant is not currently in use.
// const initModelTemplatePathConnection = "schemas/constructs/v1beta1/connection/connection_template"

// TODO
// if csv output is not directory based
// should it have different text for csv output format?
const initModelNextStepsText = `Next steps:
1. cd {modelVersionFolder}
2. Edit model.{outputFormat} to customize your model configuration
3. Add your components in the components/ directory
4. Define relationships in relationships/ directory
5. Add your connections in the connections/ directory
6. Define credentials in credentials/ directory
7. Use 'mesheryctl model build' to package your model

To import this model into Meshery:
$ mesheryctl model import {modelFolder}

To export this model as OCI image:
$ mesheryctl model build {modelName}/{modelVersion} --path {path}

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
			"model": initModelTemplatePathModel,
		},
		beforeHook: func() {
			utils.Log.Info("Generating model definition...")
		},
	},
	{
		folderPath: "components",
		// map file name to template key
		files: map[string]string{
			"component": initModelTemplatePathComponent,
		},
		beforeHook: func() {
			utils.Log.Info("Adding sample components...")
		},
	},
	{
		folderPath: "relationships",
		// map file name to template key
		files: map[string]string{
			"relationship": initModelTemplatePathRelationship,
		},
		beforeHook: func() {
			utils.Log.Info("Creating sample relationships...")
		},
	},
	{
		folderPath: "connections",
		// map file name to template key
		files: nil,
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
	if outputFormat == "json" || outputFormat == "yaml" {
		return initModelReadTemplate(
			strings.Join(
				[]string{templatePath, outputFormat},
				".",
			),
		)
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

func initModelValidateDataOverSchema(schema []byte, data map[string]interface{}) error {
	// Parse the YAML schema to extract per-property sub-schemas
	var schemaMap map[string]interface{}
	if err := yaml.Unmarshal(schema, &schemaMap); err != nil {
		return fmt.Errorf("failed to parse schema YAML: %s", err)
	}

	properties, _ := schemaMap["properties"].(map[string]interface{})
	if properties == nil {
		return nil
	}

	c := jsonschema.NewCompiler()

	var validationErrors []string
	for property, value := range data {
		propSchema, ok := properties[property].(map[string]interface{})
		if !ok {
			continue
		}

		// Build a minimal JSON Schema for this single property, excluding $ref
		// to avoid unresolvable external references.
		minimalSchema := make(map[string]interface{})
		for k, v := range propSchema {
			if k == "$ref" {
				continue
			}
			minimalSchema[k] = v
		}

		schemaJSON, err := json.Marshal(minimalSchema)
		if err != nil {
			continue
		}

		resourceName := property + ".json"
		schemaDoc, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaJSON))
		if err != nil {
			continue
		}
		if err := c.AddResource(resourceName, schemaDoc); err != nil {
			continue
		}
		sch, err := c.Compile(resourceName)
		if err != nil {
			continue
		}

		if err := sch.Validate(value); err != nil {
			validationErr, ok := err.(*jsonschema.ValidationError)
			if ok {
				msgs := collectValidationErrors(validationErr)
				for _, msg := range msgs {
					validationErrors = append(validationErrors, fmt.Sprintf("%s %s", property, msg))
				}
			} else {
				validationErrors = append(validationErrors, fmt.Sprintf("%s %s", property, err))
			}
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(validationErrors, "; "))
	}

	return nil
}

var englishPrinter = message.NewPrinter(language.English)

// collectValidationErrors extracts leaf error messages from a ValidationError tree.
func collectValidationErrors(ve *jsonschema.ValidationError) []string {
	if len(ve.Causes) == 0 && ve.ErrorKind != nil {
		return []string{ve.ErrorKind.LocalizedString(englishPrinter)}
	}
	var msgs []string
	for _, cause := range ve.Causes {
		msgs = append(msgs, collectValidationErrors(cause)...)
	}
	return msgs
}

