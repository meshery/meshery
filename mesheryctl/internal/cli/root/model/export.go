package model

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type outputDetail struct {
	Format string
	Type   string
	Path   string
}

// mesheryctl model export <designname>
var exportModelCmd = &cobra.Command{
	Use:   "export",
	Short: "Export registered models",
	Long:  "Export the registered model to the specified output type",
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
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl model export [model-name ]\nRun 'mesheryctl model export --help' to see detailed help message"
		if len(args) == 0 {
			return utils.ErrInvalidArgument(errors.New("Please provide a model name. " + errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}
		baseUrl := mctlCfg.GetBaseMesheryURL()
		modelName := args[0]
		outputFormat, _ := cmd.Flags().GetString("output-format")
		outputType, _ := cmd.Flags().GetString("output-type")
		discardComponents, _ := cmd.Flags().GetBool("discard-components")
		discardRelationships, _ := cmd.Flags().GetBool("discard-relationships")
		version, _ := cmd.Flags().GetString("version")
		page, _ := cmd.Flags().GetInt("page")
		url := fmt.Sprintf("%s/api/meshmodels/export?name=%s&output_format=%s&file_type=%s&components=%t&relationships=%t&%s", baseUrl, modelName, outputFormat, outputType, !discardComponents, !discardRelationships, utils.GetPageQueryParameter(cmd, page))
		if version != "" {
			url += fmt.Sprintf("&version=%s", version)
		}

		outputPath, _ := cmd.Flags().GetString(("output-location"))

		output := &outputDetail{
			Format: outputFormat,
			Type:   outputType,
			Path:   outputPath,
		}

		return export(args[0], url, output)
	},
}

func export(modelName string, url string, output *outputDetail) error {
	// Find the entity with the model name
	req, err := utils.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	// ensure proper cleaning of resources
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	var exportedModelPath string
	if output.Type != "oci" {
		exportedModelPath = filepath.Join(output.Path, modelName+"."+"tar.gz")
	} else {
		exportedModelPath = output.Path + modelName + ".tar"
	}
	err = os.WriteFile(exportedModelPath, data, 0644)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	utils.Log.Infof("Exported model to %s", exportedModelPath)
	return nil
}

func init() {
	exportModelCmd.Flags().StringP("output-format", "t", "yaml", "(optional) format to display in [json|yaml] (default = yaml)")
	exportModelCmd.Flags().StringP("output-location", "l", "./", "(optional) output location (default = current directory)")
	exportModelCmd.Flags().StringP("output-type", "o", "oci", "(optional) format to display in [oci|tar] (default = oci)")
	exportModelCmd.Flags().BoolP("discard-components", "c", false, "(optional) whether to discard components in the exported model definition (default = false)")
	exportModelCmd.Flags().BoolP("discard-relationships", "r", false, "(optional) whether to discard relationships in the exported model definition (default = false)")
	exportModelCmd.Flags().StringP("version", "", "", "(optional) model version to export (default = \"\")")
	exportModelCmd.Flags().IntP("page", "p", 1, "(optional) List next set of models with --page (default = 1)")
}
