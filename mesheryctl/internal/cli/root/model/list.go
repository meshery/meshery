package model

import (
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
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
		pageSize, _ := cmd.Flags().GetInt("pagesize")
		modelData := display.DisplayDataAsync{
			UrlPath:          modelsApiPath,
			DataType:         "model",
			Header:           []string{"Model", "Category", "Version"},
			Page:             page,
			PageSize:         pageSize,
			IsPage:           cmd.Flags().Changed("page"),
			DisplayCountOnly: cmd.Flags().Changed("count"),
		}

		return display.ListAsyncPagination(modelData, generateModelDataToDisplay)
	},
}

func init() {
	listModelCmd.Flags().IntP("page", "p", 1, "(optional) List next set of models with --page (default = 0)")
	listModelCmd.Flags().IntP("pagesize", "s", 0, "(optional) List next set of models with --pagesize (default = 0)")
	listModelCmd.Flags().BoolP("count", "c", false, "(optional) Get the number of models in total")
}
