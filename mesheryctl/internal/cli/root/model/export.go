package model

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type outputDetail struct {
	Format string
	Type   string
	Path   string
}

// HTTPClient interface for dependency injection
type HTTPClient interface {
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
	MakeRequest(req *http.Request) (*http.Response, error)
}

// DefaultHTTPClient implements HTTPClient
type DefaultHTTPClient struct{}

func (d *DefaultHTTPClient) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return utils.NewRequest(method, url, body)
}

func (d *DefaultHTTPClient) MakeRequest(req *http.Request) (*http.Response, error) {
	return utils.MakeRequest(req)
}

var defaultHTTPClient HTTPClient = &DefaultHTTPClient{}

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
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl model export [model-name ]\nRun 'mesheryctl model export --help' to see detailed help message"
		if len(args) == 0 {
			return utils.ErrInvalidArgument(errors.New("Please provide a model name. " + errMsg))
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		outputFormat, _ := cmd.Flags().GetString("output-format")
		validFormat := map[string]bool{"yaml": true, "json": true}
		if !validFormat[outputFormat] {
			return ErrModelUnsupportedOutputFormat(fmt.Sprintf("invalid value %q  for flag -t, allowed: yaml, json", outputFormat))
		}
		outputType, _ := cmd.Flags().GetString("output-type")
		validType := map[string]bool{"oci": true, "tar": true}
		if !validType[outputType] {
			return ErrModelUnsupportedOutputType(fmt.Sprintf("invalid value %q for flag -o, allowed: oci, tar", outputType))
		}

		page, _ := cmd.Flags().GetInt("page")
		if page < 1 {
			return ErrModelInvalidPageNumber(fmt.Sprintf("invalid page number %d, only page number > 0 allowed", page))
		}
		version, _ := cmd.Flags().GetString("version")
		if version != "" {
			matched, _ := regexp.MatchString(`^v[0-9]+\.[0-9]+\.[0-9]+$`, version)
			if !matched {
				return ErrModelUnsupportedVersion(fmt.Sprintf("invalid format %q, expected format: vMAJOR:MINOR:PATCH (ex:v0.7.3)", version))
			}
		}
		outputPath, _ := cmd.Flags().GetString(("output-location"))
		info, err := os.Stat(outputPath)
		if os.IsNotExist(err) {
			return ErrFolderNotPresent(fmt.Sprintf("output location does not exist %q", outputPath))
		} else if err != nil {
			return ErrFolderNotPresent(fmt.Sprintf("error accessing output location %q: %v", outputPath, err))
		}
		if !info.IsDir() {
			return ErrFolderNotPresent(fmt.Sprintf("output location %q is not a directory", outputPath))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return ErrProcessingConfig(fmt.Sprintf("error processing config %s", err))
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

		return exportWithClient(modelName, url, output, defaultHTTPClient)
	},
}

// function to accept HTTPClient for DI
func exportWithClient(modelName string, url string, output *outputDetail, client HTTPClient) error {
	req, err := client.NewRequest("GET", url, nil)
	if err != nil {
		return ErrExportModel(err, modelName)
	}
	if req == nil {
		return ErrExportModel(fmt.Errorf("request is nil"), modelName)
	}

	resp, err := client.MakeRequest(req)
	if err != nil {
		return ErrorExportModel(fmt.Errorf("request failed: %w", err))
	}
	if resp == nil {
		return ErrExportModel(fmt.Errorf("response is nil"), modelName)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrExportModel(fmt.Errorf("failed to export model: status %d %s", resp.StatusCode, resp.Status))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ErrExportModel(err, modelName)
	}

	var exportedModelPath string
	if output.Type == "oci" {
		exportedModelPath = filepath.Join(output.Path, modelName+".tar")
	} else if output.Type == "tar" {
		exportedModelPath = filepath.Join(output.Path, modelName+".tar.gz")
	} else {
		return ErrModelUnsupportedOutputType(fmt.Sprintf("unsupported output type: %s", output.Type))
	}

	err = os.WriteFile(exportedModelPath, data, 0644)
	if err != nil {
		return ErrExportModel(err, modelName)
	}

	return nil
}

func export(modelName string, url string, output *outputDetail) error {
	return exportWithClient(modelName, url, output, defaultHTTPClient)
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
