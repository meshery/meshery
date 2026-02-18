package model

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type exportModelFlags struct {
	DiscardComponents    bool   `json:"discard-components" validate:"boolean"`
	DiscardRelationships bool   `json:"discard-relationships" validate:"boolean"`
	OutputFormat         string `json:"output-format" validate:"required,oneof=json yaml"`
	OutputLocation       string `json:"output-location" validate:"required,dirpath"`
	OutputType           string `json:"output-type" validate:"required,oneof=oci tar"`
	Page                 int    `json:"page" validate:"omitempty,min=1"`
	Version              string `json:"version" validate:"omitempty,semver"`
}

type outputDetail struct {
	Format string
	Type   string
	Path   string
}

var (
	exportModelFlagsProvided exportModelFlags
)

// mesheryctl model export <designname>
var exportModelCmd = &cobra.Command{
	Use:   "export",
	Short: "Export registered models",
	Long: `Export the registered model to the specified output type
Documentation for models export can be found at https://docs.meshery.io/reference/mesheryctl/model/export`,
	Example: `
// Export a model by name 
mesheryctl model export [model-name] -o [oci|tar]  (default is oci)

// Export a model by name in JSON type
mesheryctl model export [model-name] -t [yaml|json] (default is YAML)

// Export a model by name in YAML type in a specific location
mesheryctl model export [model-name] -l [path-to-location]

// Export a model by name in YAML type discarding components and relationships
mesheryctl model export [model-name] --discard-components --discard-relationships

// Export a model version by name in YAML type
mesheryctl model export [model-name] --version [version (ex: v0.7.3)]
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flagValidator, ok := cmd.Context().Value("flags-validator").(*mesheryctlflags.FlagValidator)
		if !ok || flagValidator == nil {
			return utils.ErrCommandContextMissing("flags-validator")
		}
		err := flagValidator.Validate(exportModelFlagsProvided)
		if err != nil {
			return utils.ErrFlagsInvalid(err)
		}
		return nil
	},
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl model export [model-name]\nRun 'mesheryctl model export --help' to see detailed help message"
		if len(args) == 0 {
			return utils.ErrInvalidArgument(errors.New("Please provide a model name. " + errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		queryParams := url.Values{}
		queryParams.Set("name", modelName)
		queryParams.Set("output_format", exportModelFlagsProvided.OutputFormat)
		queryParams.Set("file_type", exportModelFlagsProvided.OutputType)
		queryParams.Set("components", fmt.Sprintf("%t", !exportModelFlagsProvided.DiscardComponents))
		queryParams.Set("relationships", fmt.Sprintf("%t", !exportModelFlagsProvided.DiscardRelationships))
		queryParams.Set("page", fmt.Sprintf("%d", exportModelFlagsProvided.Page))
		if exportModelFlagsProvided.Version != "" {
			queryParams.Set("version", exportModelFlagsProvided.Version)
		}

		urlPath := fmt.Sprintf("api/meshmodels/export?%s", queryParams.Encode())

		output := &outputDetail{
			Format: exportModelFlagsProvided.OutputFormat,
			Type:   exportModelFlagsProvided.OutputType,
			Path:   exportModelFlagsProvided.OutputLocation,
		}

		exportedModelData, err := api.FetchData(urlPath)
		if err != nil {
			return err
		}

		var exportedModelPath string

		extension := "tar"
		if output.Type != "oci" {
			extension = "tar.gz"
		}

		exportedModelPath = filepath.Join(output.Path, fmt.Sprintf("%s.%s", modelName, extension))

		err = os.WriteFile(exportedModelPath, exportedModelData, 0o644)
		if err != nil {
			return utils.ErrCreateFile(exportedModelPath, err)
		}

		utils.Log.Infof("Exported model to %s", exportedModelPath)

		return nil
	},
}

func init() {
	exportModelCmd.Flags().BoolVarP(&exportModelFlagsProvided.DiscardComponents, "discard-components", "c", false, "(optional) whether to discard components in the exported model definition (default = false)")
	exportModelCmd.Flags().BoolVarP(&exportModelFlagsProvided.DiscardRelationships, "discard-relationships", "r", false, "(optional) whether to discard relationships in the exported model definition (default = false)")
	exportModelCmd.Flags().StringVarP(&exportModelFlagsProvided.OutputFormat, "output-format", "t", "yaml", "(optional) format to display in [json|yaml] (default = yaml)")
	exportModelCmd.Flags().StringVarP(&exportModelFlagsProvided.OutputLocation, "output-location", "l", "./", "(optional) output location (default = current directory)")
	exportModelCmd.Flags().StringVarP(&exportModelFlagsProvided.OutputType, "output-type", "o", "oci", "(optional) format to display in [oci|tar] (default = oci)")
	exportModelCmd.Flags().IntVarP(&exportModelFlagsProvided.Page, "page", "p", 1, "(optional) page number for paginated results (default = 1)")
	exportModelCmd.Flags().StringVarP(&exportModelFlagsProvided.Version, "version", "", "", "(optional) model version to export (default = \"\", format: vX.X.X)")
}
