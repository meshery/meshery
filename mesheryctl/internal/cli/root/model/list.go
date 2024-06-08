package model

import (
	"encoding/json"
	"fmt"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

var listModelCmd = &cobra.Command{
	Use:   "list",
	Short: "list registered models",
	Long:  "list name of all registered models",
	Example: `
// View list of models
mesheryctl model list

// View list of models with specified page number (25 models per page)
mesheryctl model list --page 2

// View number of available models in Meshery
mesheryctl model list --count
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemModelSubError("this command takes no arguments\n", "list"))
		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := fmt.Sprintf("%s/api/meshmodels/models?%s", baseUrl, getPageQueryParameter(cmd))

		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		modelsResponse := &models.MeshmodelsAPIResponse{}
		err = json.Unmarshal(data, modelsResponse)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

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

		if cmd.Flag("count").Value.String() == "true" {
			// fmt.Println("Total number of models: ", len(rows))
			whiteBoardPrinter.Println("Total number of models: ", modelsResponse.Count)
			return nil
		}

		if cmd.Flags().Changed("page") {
			utils.PrintToTable(header, rows)
		} else {
			err := utils.HandlePagination(maxRowsPerPage, "models", rows, header)
			if err != nil {
				utils.Log.Error(err)
				return err
			}
		}

		return nil
	},
}

func getPageQueryParameter(cmd *cobra.Command) string {
	if cmd.Flag("count").Value.String() == "true" {
		return "page=1"
	}
	if cmd.Flags().Changed("page") {
		return fmt.Sprintf("page=%d", pageNumberFlag)
	}
	return "pagesize=all"
}
