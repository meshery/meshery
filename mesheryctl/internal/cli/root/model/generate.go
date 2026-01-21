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
	Long: `Generate models from CSV files or URLs with a template JSON file.

Use 'generate' when:
- Creating new models from CSV files (model.csv, component.csv, relationship.csv)
- Generating models from URLs using a template JSON file for model metadata
- You need to skip registration during development/testing (-r flag)

Use 'import' instead when:
- Importing pre-existing model packages (tar.gz, directories)
- Re-importing models exported from Meshery UI
- Importing from URLs without additional template configuration

Documentation: https://docs.meshery.io/reference/mesheryctl/model/generate`,
	Example: `
// Generate a model from a CSV file directory
mesheryctl model generate -f [path-to-csv-directory]

// Generate a model from a URL using a JSON template
mesheryctl model generate -f [URL] -t [path-to-template.json]

// Generate a model from a URL, skip registration
mesheryctl model generate -f [URL] -t [path-to-template.json] -r
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
	// Common flag shared with 'import' command - specifies input source
	// For 'generate': CSV directory or URL (requires -t template for URLs)
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
