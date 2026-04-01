package model

import (
	"fmt"
	"os"
	"strings"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type cmdModelGenerateFlags struct {
	File     string `json:"file" validate:"omitempty,dirpath|filepath|url"`
	Template string `json:"template" validate:"omitempty,filepath"`
	Register bool   `json:"register" validate:"boolean"`
}

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

var modelGenerateFlags cmdModelGenerateFlags

var generateModelCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate models from a file",
	Long: `Generate models by specifying the directory, file, or URL. You can also provide a template JSON file and registrant name
Find more information at: https://docs.meshery.io/reference/mesheryctl/model/generate`,
	Example: ` 
// Generate a model from a CSV directory
mesheryctl model generate --f [path-to-csv-directory]

// Generate a model from a URL based on a JSON template
mesheryctl model generate --f [URL] -t [path-to-template.json]

// Generate a model from a URL based on a JSON template skipping registration
mesheryctl model generate --f [URL] -t [path-to-template.json] -r
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &modelGenerateFlags)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if modelGenerateFlags.File == "" && len(args) == 0 {
			return utils.ErrInvalidArgument(fmt.Errorf(errGenerateMissingArgsMsg, errGenerateUsageMsg))
		}

		if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments\n\n%s", errGenerateUsageMsg))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		path := modelGenerateFlags.File
		// If file flag is not provided, use the argument as the path
		if path == "" {
			path = args[0]

		}

		isUrl := utils.IsValidUrl(path)
		if isUrl {
			if modelGenerateFlags.Template == "" {
				return ErrTemplateFileNotPresent()
			}

			urlModelGenerator := &UrlModelGenerator{
				TemplateFile: modelGenerateFlags.Template,
				Url:          path,
				SkipRegister: modelGenerateFlags.Register,
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
			SkipRegister:     modelGenerateFlags.Register,
		}

		return csvModelGenerator.Generate()
	},
}

func init() {
	generateModelCmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})

	generateModelCmd.Flags().StringVarP(&modelGenerateFlags.File, "file", "f", "", "Specify path to the file or directory")
	generateModelCmd.Flags().StringVarP(&modelGenerateFlags.Template, "template", "t", "", "Specify path to the template JSON file")
	generateModelCmd.Flags().BoolVarP(&modelGenerateFlags.Register, "register", "r", false, "Skip registration of the model")

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
