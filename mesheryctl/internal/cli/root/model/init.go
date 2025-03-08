package model

import (
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		return nil
	},
}
