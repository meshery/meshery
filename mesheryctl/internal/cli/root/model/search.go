package model

import (
	"fmt"
	"net/url"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var searchModelCmd = &cobra.Command{
	Use:   "search",
	Short: "Search model(s)",
	Long: `Search model(s) by search string
Documentation for models search can be found at https://docs.meshery.io/reference/mesheryctl/model/search`,
	Example: `
// Search model from current provider
mesheryctl model search [query-text]
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl model search [query-text]\nRun 'mesheryctl model search --help' to see detailed help message"
		if len(args) == 0 {
			return utils.ErrInvalidArgument(errors.New("Please provide a model name. " + errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		queryText := args[0]

		url := fmt.Sprintf("%s/%s?search=%s&pagesize=all", baseUrl, modelsApiPath, url.QueryEscape(queryText))

		modelsResponse, err := api.Fetch[models.MeshmodelsAPIResponse](url)

		if err != nil {
			return err
		}

		header := []string{"Model", "Category", "Version"}
		rows := [][]string{}

		for _, model := range modelsResponse.Models {
			if len(model.DisplayName) > 0 {
				rows = append(rows, []string{model.Name, model.Category.Name, model.Version})
			}
		}

		dataToDisplay := display.DisplayedData{
			DataType:         "model",
			Header:           header,
			Rows:             rows,
			Count:            int64(modelsResponse.Count),
			DisplayCountOnly: false,
			IsPage:           cmd.Flags().Changed("page"),
		}

		err = display.List(dataToDisplay)

		if err != nil {
			return err
		}

		return nil
	},
}
