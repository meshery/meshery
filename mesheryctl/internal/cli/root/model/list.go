package model

import (
	"fmt"

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
	Long:  "List all registered models by pagingation (25 models per page)",
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
		url := fmt.Sprintf("%s/api/meshmodels/models?%s", baseUrl, utils.GetPageQueryParameter(cmd, page))

		models, err := fetchModels(url)

		if err != nil {
			return err
		}

		count, _ := cmd.Flags().GetBool("count")
		return listModel(cmd, models, count)
	},
}

func init() {
	listModelCmd.Flags().IntP("page", "p", 1, "(optional) List next set of models with --page (default = 1)")
	listModelCmd.Flags().BoolP("count", "c", false, "(optional) Get the number of models in total")
}

func listModel(cmd *cobra.Command, modelsResponse *models.MeshmodelsAPIResponse, displayCountOnly bool) error {

	header := []string{"Model", "Category", "Version"}
	rows := [][]string{}

	for _, model := range modelsResponse.Models {
		if len(model.DisplayName) > 0 {
			rows = append(rows, []string{model.Name, model.Category.Name, model.Version})
		}
	}

	if len(rows) == 0 {
		// if no model is found
		// fmt.Println("No model(s) found")
		whiteBoardPrinter.Println("No model(s) found")
		return nil
	}

	utils.DisplayCount("models", modelsResponse.Count)

	if displayCountOnly {
		return nil
	}

	if cmd.Flags().Changed("page") {
		utils.PrintToTable(header, rows)
	} else {
		err := utils.HandlePagination(maxRowsPerPage, "models", rows, header)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
	}

	return nil
}
