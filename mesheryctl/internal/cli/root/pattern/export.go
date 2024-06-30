package pattern

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	designType string
	designID   string
	outputDir  string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a design from Meshery",
	Long:  "The 'export' command allows you to export a specific design from your Meshery server. You can specify the design by its ID and optionally define the type of design. The command also supports specifying an output directory where the exported design will be saved. By default, the exported design will be saved in the home directory.",
	Example: `
	# Export a design with a specific ID
	mesheryctl pattern export --id [design-ID]
	
	# Export a design with a specific ID and type
	mesheryctl pattern export --id [design-ID] --dtype [design-type]
	
	# Export a design and save it to a specific directory
	mesheryctl pattern export --id [design-ID] --output ./designs
	
	# Export a design with a specific type and save it to a directory
	mesheryctl pattern export --id [design-ID] --dtype [design-type] --output ./exports
	`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		return ctx.ValidateVersion()
	},
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		designID, _ := cmd.Flags().GetString("id")
		designType, _ := cmd.Flags().GetString("dtype")
		if designType == "" {
			designType = "ahpkg"
		}
		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := baseUrl + "/api/pattern/download/" + designID

		utils.Log.Info(fmt.Sprintf("Exporting Design of type %s with ID %s", designType, designID))
		switch designType {
		case "oci":
			url = url + "?oci=true"
		case "ahpkg":
			url = url + "?ahpkg=true"
		}
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
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = models.ErrDoRequest(err, resp.Request.Method, url)
			utils.Log.Error(err)
			return nil
		}
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			utils.Log.Error(ErrReadFromBody(err))
			return nil
		}
		var patternName string
		if designType != "oci" {
			var pattern models.MesheryPattern
			patternContent := buf.String()
			err = yaml.Unmarshal([]byte(patternContent), &pattern)
			if err != nil {
				utils.Log.Error(err)
				return nil
			}
			patternName = pattern.Name
		}

		if patternName == "" {
			patternName = "pattern"
		}

		designIDParts := strings.Split(designID, "-")
		lastPartOfID := designIDParts[len(designIDParts)-1]

		filename := fmt.Sprintf("%s_%s", patternName, lastPartOfID)
		if designType != "oci" {
			filename += ".yaml"
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			err = ErrRetrieveHomeDir(err)
			utils.Log.Error(err)
			return nil
		}

		outputFilePath := filepath.Join(homeDir, outputDir, filename)

		// Check if file exists and modify filename if needed
		outputFilePath = getUniqueFilename(outputFilePath)

		err = os.MkdirAll(filepath.Dir(outputFilePath), 0755)
		if err != nil {
			err = models.ErrMakeDir(err, outputFilePath)
			utils.Log.Error(err)
			return nil
		}

		file, err := os.Create(outputFilePath)
		if err != nil {
			err = ErrCreateFile(outputFilePath, err)
			utils.Log.Error(err)
			return nil
		}
		defer file.Close()

		_, err = io.Copy(file, buf)
		if err != nil {
			err = ErrCopyData(outputFilePath, err)
			utils.Log.Error(err)
			return nil
		}

		utils.Log.Info("Design exported successfully to ", outputFilePath)

		return nil
	},
}

func getUniqueFilename(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	for i := 1; ; i++ {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			break
		}
		filename = fmt.Sprintf("%s(%d)%s", base, i, ext)
	}
	return filename
}

func init() {
	exportCmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})

	exportCmd.Flags().StringVarP(&designID, "id", "i", "", "Specify the design ID to export")
	exportCmd.Flags().StringVarP(&designType, "dtype", "d", "", "Specify the design type to export")
	exportCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Specify the output directory to save the design")

	err := exportCmd.MarkFlagRequired("id")
	if err != nil {
		utils.Log.Error(ErrMarkFlagRequire("id", err))
	}
}
