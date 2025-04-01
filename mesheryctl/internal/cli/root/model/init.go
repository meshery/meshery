package model

import (
	"fmt"
	"os"
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
	Use:   "init [model-name]",
	Short: "generates scaffolding for convenient model creation",
	Long:  "generates a folder structure and guides user on model creation",
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
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
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

		version, _ := cmd.Flags().GetString("version")
		if !semver.IsValid(version) {
			return ErrModelUnsupportedVersion(
				"version must follow a semver format, f.e. v1.2.3",
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return ErrModelInit(err)
		}

		modelName := args[0]
		path, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		outputFormat, _ := cmd.Flags().GetString("output-format")

		utils.Log.Infof("Creating new Meshery model: %s", modelName)

		const DirPerm = 0755
		modelFolder := strings.Join([]string{path, modelName}, string(os.PathSeparator))
		modelVersionFolder := strings.Join([]string{modelFolder, version}, string(os.PathSeparator))
		utils.Log.Infof("Creating directory structure...")
		err = os.MkdirAll(modelVersionFolder, DirPerm)
		if err != nil {
			return ErrModelInit(err)
		}

		for _, item := range initModelData {
			item.beforeHook()
			itemFolderPath := modelVersionFolder
			if item.folderPath != "" {
				itemFolderPath = strings.Join([]string{modelVersionFolder, item.folderPath}, string(os.PathSeparator))
				err = os.MkdirAll(itemFolderPath, DirPerm)
				if err != nil {
					return ErrModelInit(err)
				}
			}
			for name, templatePath := range item.files {
				content, err := getTemplateInOutputFormat(templatePath, outputFormat)
				if err != nil {
					return ErrModelInit(err)
				}
				filePath := strings.Join(
					[]string{
						itemFolderPath,
						strings.Join(
							[]string{name, outputFormat},
							".",
						),
					},
					string(os.PathSeparator),
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

		// TODO think about cleaning up partial data (if error occurs in the middle of execution).
		// if delete a folder, only delete if the folder was created
		// if the user specifies an existing folder it should not be deleted.
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
