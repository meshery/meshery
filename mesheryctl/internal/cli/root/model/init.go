package model

import (
	"encoding/json"
	"os"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
)

var initModelCmd = &cobra.Command{
	Use:   "init",
	Short: "generates scaffolding for convenient model creation",
	Long:  "generates a folder structure and guides user on model creation",
	Example: `
// generates a folder structure
mesheryctl model init

// generates a folder structure and sets up model version
mesheryctl model init --version 2.0.8 (default is 0.1.0) 

// generates a folder structure under specified path
mesheryctl model init --path path/to/some/particular_folder (default is current folder) 

// generate a folder structure in json format
mesheryctl model init --output-format yaml (default is json) 
    `,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		utils.Log.Info("init command will be here soon")

		// ---
		// Code below is work in progress.
		//

		model := model.ModelDefinition{
			Id:          uuid.Nil,
			Name:        "model name", // TODO
			DisplayName: "Human-readable name for the model.",
			Model:       model.Model{Version: versionFlag},
			Description: "Description of the model.",
			Status:      model.ModelDefinitionStatus("duplicate|maintenance|enabled|ignored"),
			CategoryId:  uuid.Nil,
			Category: category.CategoryDefinition{
				Id:       uuid.Nil,
				Name:     "Category of the model.",
				Metadata: map[string]any{"key": "value"},
			},
			SchemaVersion: "http://json-schema.org/draft-07/schema#",
			SubCategory:   "Sub-category of the model.", //
			Metadata:      nil,                          // TODO
			Registrant:    connection.Connection{},      // TODO
			Version:       versionFlag,
		}
		jsonData, err := json.MarshalIndent(model, "", "  ")
		if err != nil {
			// TODO meshkit error format
			utils.Log.Error(err)
			return nil
		}

		// Write the JSON to a file
		filename := "model.json"
		if err := os.WriteFile(filename, jsonData, 0644); err != nil {
			// TODO meshkit error format
			utils.Log.Error(err)
			return nil
		}

		utils.Log.Debug("Done!")

		//
		// End of work in progress block.
		// ---

		return nil
	},
}
