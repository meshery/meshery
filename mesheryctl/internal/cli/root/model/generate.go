package model

import (
	"fmt"
	"os"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var generateModelCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate models from mesheryctl command",
	Long:  "Generate models by specifying the directory, file, or URL. You can also provide a template JSON file and registrant name.",
	Example: `
	mesheryctl model generate -f [ URI ]
	mesheryctl model generate -f [ URI ] -t [ path to template file ] ( only required in case of URL )
	mesheryctl model generate -f [ URI ] -t [ path to template file ] -r ( to skip registration by default registration is true)
 
	
	mesheryctl model generate --f /path/to/csv-drectory
    mesheryctl model generate --f http://example.com/model -t /path/to/template.json 
	mesheryctl model generate --f http://example.com/model -t /path/to/template.json -r
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl model generate [ file | filePath | URL ]\nRun 'mesheryctl model generate --help' to see detailed help message"
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
		isUrl := utils.IsValidUrl(path)

		if isUrl {
			if templateFile == "" {
				return ErrTemplateFileNotPresent()
			}

			fileData, err := os.ReadFile(templateFile)
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
			modelcsvpath, componentcsvpath, relationshipcsvpath, err := utils.GetCsv(path)
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

	generateModelCmd.Flags().StringVarP(&location, "file", "f", "", "Specify path to the file or directory")
	generateModelCmd.Flags().StringVarP(&templateFile, "template", "t", "", "Specify path to the template JSON file")
	generateModelCmd.Flags().BoolVarP(&register, "register", "r", false, "Skip registration of the model")

}
