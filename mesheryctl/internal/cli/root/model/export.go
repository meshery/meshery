package model

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// mesheryctl model export <designname>
var exportModal = &cobra.Command{
	Use:   "export",
	Short: "export registered models",
	Long:  "export the registered model to the specified output type",
	Example: `
// Export a model by name
mesheryctl model export [model-name] -o [oci|tar]  (default is oci)
mesheryctl model export [model-name] -t json (default is yaml)
mesheryctl model export [model-name] -l /home/meshery/
mesheryctl model export [model-name] --discard-components --discard-relationships
mesheryctl model export [model-name] --version v0.7.3

    `,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl model export [model-name]\nRun 'mesheryctl model export --help' to see detailed help message"
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
		outFormat := "yaml"
		if outFormatFlag == "json" {
			outFormat = "json"
		}
		outfileType := "oci"
		if outTypeFlag == "tar" {
			outfileType = "tar"
		}
		url := fmt.Sprintf("%s/api/meshmodels/export?name=%s&output_format=%s&file_type=%s&components=%t&relationships=%t&%s", baseUrl, modelName, outFormat, outfileType, !discardComponentsFlag, !discardRelationshipsFlag, utils.GetPageQueryParameter(cmd, pageNumberFlag))
		if versionFlag != "" {
			url += fmt.Sprintf("&version=%s", versionFlag)
		}
		return exportModel(args[0], cmd, url, false)
	},
}
