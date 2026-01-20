package model

import (
	"fmt"
	"os"
	"strings"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type ModelGenerator interface {
	Generate() error
}

type UrlModelGenerator struct {
	TemplateFile string
	Url          string
	SkipRegister bool
}

type CsvModelGenerator struct {
	ModelFile        string
	ComponentFile    string
	RelationshipFile string
	SkipRegister     bool
}

var generateModelCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate models from a file",
	Long: `Generate models by specifying the directory, file, or URL. You can also provide a template JSON file and registrant name
Documentation for models generate can be found at https://docs.meshery.io/reference/mesheryctl/model/generate`,
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

		// Path is a url
		if isUrl {
			template, _ := cmd.Flags().GetString("template")
			if template == "" {
				return ErrTemplateFileNotPresent()
			}

			urlModelGenerator := &UrlModelGenerator{
				TemplateFile: template,
				Url:          path,
				SkipRegister: register,
			}
			return urlModelGenerator.Generate()
		}

		// Path is a file or directory
		err := meshkitRegistryUtils.SetLogger(true)
		if err != nil {
			utils.Log.Info("Error setting logger: ", err)
		}

		modelcsvpath, componentcsvpath, relationshipcsvpath, err := meshkitRegistryUtils.GetCsv(path)

		if err != nil {
			return err
		}

		csvModelGenerator := &CsvModelGenerator{
			ModelFile:        modelcsvpath,
			ComponentFile:    componentcsvpath,
			RelationshipFile: relationshipcsvpath,
			SkipRegister:     register,
		}

		return csvModelGenerator.Generate()
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

func (u *UrlModelGenerator) Generate() error {
	utils.Log.Info("Generating model from URL: ", u.Url)

	fileData, err := os.ReadFile(u.TemplateFile)
	if err != nil {
		return utils.ErrFileRead(err)
	}

	err = registerModel(fileData, nil, nil, "", "url", u.Url, !u.SkipRegister)
	if err != nil {
		return err
	}

	locationForModel := utils.MesheryFolder + "/models"
	utils.Log.Info("Model can be accessed from ", locationForModel)

	return nil
}

func (c *CsvModelGenerator) Generate() error {
	utils.Log.Info("Generating model from CSV files")

	var modelData, componentData, relationshipData []byte
	var err error

	filePaths := []struct {
		path string
		data *[]byte
	}{
		{c.ModelFile, &modelData},
		{c.ComponentFile, &componentData},
		{c.RelationshipFile, &relationshipData},
	}

	for _, f := range filePaths {
		*f.data, err = os.ReadFile(f.path)
		if err != nil {
			return utils.ErrFileRead(err)
		}
	}

	err = registerModel(modelData, componentData, relationshipData, "model.csv", "csv", "", !c.SkipRegister)
	if err != nil {
		return err
	}

	locationForModel := utils.MesheryFolder + "/models"
	utils.Log.Info("Model can be accessed from ", locationForModel)

	locationForLogs := utils.MesheryFolder + "/logs/registry"
	utils.Log.Info("Logs for the csv generation can be accessed ", locationForLogs)

	return nil
}
