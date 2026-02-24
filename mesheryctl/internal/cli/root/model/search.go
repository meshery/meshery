package model

import (
	"fmt"
	"net/url"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var searchModelCmd = &cobra.Command{
	Use:   "search",
	Short: "Search model(s)",
	Long: `Search model(s) by search string
Find more information at: https://docs.meshery.io/reference/mesheryctl/model/search`,
	Example: `
// Search model from current provider
mesheryctl model search [query-text]

// Search list of models for a specified page
mesheryctl model search [query-text] --page [page-number]
	`,
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return utils.ErrInvalidArgument(errors.New(errSearchModelName))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("pagesize")
		modelData := display.DisplayDataAsync{
			UrlPath:  fmt.Sprintf("%s?search=%s", modelsApiPath, url.QueryEscape(args[0])),
			DataType: "model",
			Header:   []string{"ID", "Model", "Category", "Version"},
			Page:     page,
			PageSize: pageSize,
			IsPage:   cmd.Flags().Changed("page"),
		}

		return display.ListAsyncPagination(modelData, generateModelDataToDisplay)
	},
}

func init() {
	searchModelCmd.Flags().IntP("page", "p", 1, "(optional) List next set of models with --page (default = 1)")
	searchModelCmd.Flags().IntP("pagesize", "s", 10, "(optional) List next set of models with --pagesize (default = 10)")
}
