package model

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// TODO must use schemas repo instead of local temp package
	// "github.com/meshery/schemas"
	schemas "github.com/layer5io/meshery/mesheryctl/internal/cli/root/model/temp_schemas"
)

var initModelCmd = &cobra.Command{
	Use:   "init [model-name]",
	Short: "generates scaffolding for convenient model creation",
	Long:  "generates a folder structure and guides user on model creation",
	Example: `
// generates a folder structure
mesheryctl model init [model-name]

// generates a folder structure and sets up model version
mesheryctl model init [model-name] --version [version] (default is 0.1.0)

// generates a folder structure under specified path
mesheryctl model init [model-name] --path [path-to-location] (default is current folder)

// generate a folder structure in json format
mesheryctl model init [model-name] --output-format [json|yaml|csv] (default is json)
    `,
	Args: cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("output-format")
		getValidOutputFormatSlices := getValidOutputFormat()
		if !slices.Contains(getValidOutputFormatSlices, format) {
			validFormatsString := strings.Join(getValidOutputFormat(), ", ")
			return ErrModelUnsupportedOutputFormat(
				fmt.Sprintf(
					"[ %s ] are the only format supported",
					validFormatsString,
				),
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			// TODO use meshkit error format instead during implementation phase
			return err
		}

		modelName := args[0]
		path, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")
		// TODO take into account outputFormat
		outputFormat, _ := cmd.Flags().GetString("output-format")

		utils.Log.Info("init command will be here soon")
		utils.Log.Infof("model name = %s", modelName)
		utils.Log.Infof("path = %s", path)
		utils.Log.Infof("version = %s", version)
		utils.Log.Infof("output format = %s", outputFormat)

		for _, templatePath := range []string{
			initModelTemplatePathModelJSON,
			initModelTemplatePathDesignJSON,
			initModelTemplatePathComponentJSON,
			initModelTemplatePathConnectionJSON,
			initModelTemplatePathRelathionshipJSON,
		} {
			// modelJSONContent, err := readTemplate(templatePath)
			_, err := initModelReadTemplate(templatePath)
			if err != nil {
				utils.Log.Error(err)
				// TODO use meshkit error format instead during implementation phase
				return err
			}

			//utils.Log.Debug(string(modelJSONContent))
		}

		const DirPerm = 0755
		mainFolderPath := strings.Join([]string{path, modelName, version}, string(os.PathSeparator))
		err = os.MkdirAll(mainFolderPath, DirPerm)
		if err != nil {
			// TODO use meshkit error format
			utils.Log.Error(err)
			return err
		}

		for _, item := range initModelData {
			if item.folderPath != "" {
				itemFolderPath := strings.Join([]string{mainFolderPath, item.folderPath}, string(os.PathSeparator))
				err = os.MkdirAll(itemFolderPath, DirPerm)
				if err != nil {
					// TODO use meshkit error format
					utils.Log.Error(err)
					return err
				}
			}
			for name, templatePath := range item.files {
				content, err := getTemplateInOutputFormat(templatePath, outputFormat)
				if err != nil {
					utils.Log.Error(err)
					// TODO use meshkit error format instead during implementation phase
					return err
				}
				filePath := strings.Join(
					[]string{
						mainFolderPath,
						strings.Join(
							[]string{name, outputFormat},
							".",
						),
					},
					string(os.PathSeparator),
				)
				file, err := os.Create(filePath)
				if err != nil {
					utils.Log.Error(err)
					// TODO use meshkit error format instead during implementation phase
					return err
				}
				defer file.Close() // Ensure the file is closed when the function exits

				_, err = file.Write(content)
				if err != nil {
					utils.Log.Error(err)
					// TODO use meshkit error format instead during implementation phase
					return err
				}

			}

		}

		// TODO maybe clean partial data (if error occurs in the middle of execution)
		return nil
	},
}

func init() {
	initModelCmd.Flags().StringP("path", "p", ".", "(optional) target directory (default: current dir)")
	initModelCmd.Flags().StringP("version", "", "0.1.0", "(optional) model version (default: 0.1.0)")
	initModelCmd.Flags().StringP("output-format", "o", "json", "(optional) format to display in [json|yaml]")
}

func initModelGetValidOutputFormat() []string {
	return []string{"json"}
	// TODO implement yamls, csv
	// return []string{"json", "yaml", "csv"}
}

const initModelTemplatePathModelJSON = "json_models/constructs/v1beta1/model.json"
const initModelTemplatePathDesignJSON = "json_models/constructs/v1beta1/design.json"
const initModelTemplatePathComponentJSON = "json_models/constructs/v1beta1/component.json"
const initModelTemplatePathConnectionJSON = "json_models/constructs/v1beta1/connection.json"
const initModelTemplatePathRelathionshipJSON = "json_models/constructs/v1alpha3/relationship.json"

var initModelData = []struct {
	folderPath string
	files      map[string]string
}{
	{
		folderPath: "",
		// map file name to template key
		files: map[string]string{
			"model": initModelTemplatePathModelJSON,
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
		// TODO process yaml
		return nil, errors.New("implement yaml")
	}

	if outputFormat == "csv" {
		// TODO process csv
		return nil, errors.New("implement csv")
	}

	// should impossible to get here
	// TODO use meshkit error format instead during implementation phase
	return nil, errors.New("invalid format")
}
