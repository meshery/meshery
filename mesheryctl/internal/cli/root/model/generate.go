package model

import (
	"fmt"
	"os"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/layer5io/meshkit/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var generateModelCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate models from a file",
	Long:  "Generate models by specifying the directory, file, or URL. You can also provide a template JSON file and registrant name.",
	Example: ` 
// Generate a model from a CSV file(s)
mesheryctl model generate --f [path-to-csv-drectory]

// Generate a model from a Uri baesd on a JSON template
mesheryctl model generate --f [URL] -t [path-to-template.json]

// Generate a model from a Uri baesd on a JSON template skipping registration
mesheryctl model generate --f [URL] -t [path-to-template.json] -r
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl model generate [ file | filePath | URL ]\nRun 'mesheryctl model generate --help' to see detailed help message"
		file, _ := cmd.Flags().GetString("file")
		if file == "" && len(args) == 0 {
			return fmt.Errorf("[ file | filepath | URL ] isn't specified\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		file, _ := cmd.Flags().GetString("file")
		if file != "" {
			path = file
		} else {
			path = args[0]
		}
		isUrl := utils.IsValidUrl(path)

		register, _ := cmd.Flags().GetBool("register")

		if isUrl {
			template, _ := cmd.Flags().GetString("template")
			if template == "" {
				return ErrTemplateFileNotPresent()
			}

			fileData, err := os.ReadFile(template)
			if err != nil {
				return utils.ErrFileRead(err)
			}

			err = registerModel(fileData, nil, nil, "", "url", path, !register)
			if err != nil {
				utils.Log.Error(err)
				return nil
			}
			locationForModel := utils.MesheryFolder + "/models"
			utils.Log.Info("Model can be accessed from ", locationForModel)
			return nil
		} else {

			err := meshkitRegistryUtils.SetLogger(true)
			if err != nil {
				utils.Log.Info("Error setting logger: ", err)
			}

			modelcsvpath, componentcsvpath, relationshipcsvpath, err := meshkitRegistryUtils.GetCsv(path)
			if err == nil {
				modelData, err := os.ReadFile(modelcsvpath)
				if err != nil {
					return utils.ErrFileRead(err)
				}
				componentData, err := os.ReadFile(componentcsvpath)
				if err != nil {
					return utils.ErrFileRead(err)
				}
				relationshipData, err := os.ReadFile(relationshipcsvpath)
				if err != nil {
					return utils.ErrFileRead(err)
				}
				err = registerModel(modelData, componentData, relationshipData, "model.csv", "csv", "", !register)
				if err != nil {
					utils.Log.Error(err)
					return nil
				}
				locationForModel := utils.MesheryFolder + "/models"
				utils.Log.Info("Model can be accessed from ", locationForModel)
				locationForLogs := utils.MesheryFolder + "/logs/registry"
				utils.Log.Info("Logs for the csv generation can be accessed ", locationForLogs)
			}
		}
		return nil
	},
}

func init() {
	generateModelCmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})

	generateModelCmd.Flags().StringP("file", "f", "", "Specify path to the file or directory")
	generateModelCmd.Flags().StringP("template", "t", "", "Specify path to the template JSON file")
	generateModelCmd.Flags().BoolP("register", "r", false, "Skip registration of the model")

}
