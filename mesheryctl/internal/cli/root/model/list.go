package model

import (
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type cmdModelListFlags struct {
	Page     int  `json:"page" validate:"gte=1"`
	PageSize int  `json:"pagesize" validate:"gte=1"`
	Count    bool `json:"count" validate:"boolean"`
}

var modelListFlags cmdModelListFlags

var listModelCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered models",
	Long: `List all registered models by pagingation (10 models per page)
Find more information at: https://docs.meshery.io/reference/mesheryctl/model/list`,
	Example: `
// List of models
mesheryctl model list

// List of models for a specified page
mesheryctl model list --page [page-number] --pagesize [pagesize]

// Display number of available models in Meshery
mesheryctl model list --count
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flagValidator, ok := cmd.Context().Value(mesheryctlflags.FlagValidatorKey).(*mesheryctlflags.FlagValidator)
		if !ok || flagValidator == nil {
			return utils.ErrCommandContextMissing("flags-validator")
		}
		err := flagValidator.Validate(modelListFlags)
		if err != nil {
			return utils.ErrFlagsInvalid(err)
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return utils.ErrInvalidArgument(errors.New(utils.SystemModelSubError("this command takes no arguments\n", "list")))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		modelData := display.DisplayDataAsync{
			UrlPath:          modelsApiPath,
			DataType:         "model",
			Header:           []string{"ID", "Model", "Category", "Version"},
			Page:             modelListFlags.Page,
			PageSize:         modelListFlags.PageSize,
			IsPage:           cmd.Flags().Changed("page"),
			DisplayCountOnly: modelListFlags.Count,
		}

		return display.ListAsyncPagination(modelData, generateModelDataToDisplay)
	},
}

func init() {
	listModelCmd.Flags().IntVarP(&modelListFlags.Page, "page", "p", 1, "(optional) List next set of models with --page (default = 1)")
	listModelCmd.Flags().IntVarP(&modelListFlags.PageSize, "pagesize", "s", 10, "(optional) List next set of models with --pagesize (default = 10)")
	listModelCmd.Flags().BoolVarP(&modelListFlags.Count, "count", "c", false, "(optional) Get the number of models in total")
}
