package model

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listModelCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered models",
	Long: `List all registered models by pagingation (25 models per page)
Documentation for models list can be found at https://docs.meshery.io/reference/mesheryctl/model/list`,
	Example: `
// List of models
mesheryctl model list

// List of models for a specified page
mesheryctl model list --page [page-number]

// Display number of available models in Meshery
mesheryctl model list --count
    `,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemModelSubError("this command takes no arguments\n", "list"))
		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		page, _ := cmd.Flags().GetInt("page")
		url := fmt.Sprintf("%s/%s?%s", baseUrl, modelsApiPath, utils.GetPageQueryParameter(cmd, page))

		modelsResponse, err := api.Fetch[models.MeshmodelsAPIResponse](url)

		if err != nil {
			return err
		}

		return displayModels(modelsResponse, cmd)
	},
}

func init() {
	listModelCmd.Flags().IntP("page", "p", 1, "(optional) List next set of models with --page (default = 1)")
	listModelCmd.Flags().BoolP("count", "c", false, "(optional) Get the number of models in total")
}
