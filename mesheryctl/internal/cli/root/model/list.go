package model

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemModelSubError("this command takes no arguments\n", "list"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt("page")

		modelsResponse, err := api.Fetch[models.MeshmodelsAPIResponse](fmt.Sprintf("%s?%s", modelsApiPath, utils.GetPageQueryParameter(cmd, page)))

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
