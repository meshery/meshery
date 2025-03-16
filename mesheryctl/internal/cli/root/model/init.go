package model

import (
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// TODO must use schemas repo instead of local temp package
	// "github.com/meshery/schemas"
	schemas "github.com/layer5io/meshery/mesheryctl/internal/cli/root/model/temp_schemas"
)

var initModelCmd = &cobra.Command{
	Use:   "init",
	Short: "generates scaffolding for convenient model creation",
	Long:  "generates a folder structure and guides user on model creation",
	Example: `
// generates a folder structure
mesheryctl model init

// generates a folder structure and sets up model version
mesheryctl model init --version [ version ] (default is 0.1.0) 

// generates a folder structure under specified path
mesheryctl model init --path path/to/some/particular_folder (default is current folder) 

// generate a folder structure in json format
mesheryctl model init --output-format yaml (default is json) 
    `,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			// TODO use meshkit error format instead during implementation phase
			return err
		}

		utils.Log.Info("init command will be here soon")

		for _, templatePath := range []string{
			templatePathModelJSON,
			templatePathDesignJSON,
			templatePathComponentJSON,
			templatePathConnectionJSON,
			templatePathRelathionshipJSON,
		} {
			modelJSONContent, err := readTemplate(templatePath)
			if err != nil {
				utils.Log.Error(err)
				// TODO use meshkit error format instead during implementation phase
				return err
			}

			utils.Log.Debug(string(modelJSONContent))
		}

		return nil
	},
}

func init() {
	initModelCmd.Flags().StringVarP(&targetDirectory, "path", "p", ".", "(optional) target directory (default: current dir)")
	initModelCmd.Flags().StringVarP(&versionFlag, "version", "", "0.1.0", "(optional) model version (default: 0.1.0)")
	initModelCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "json", "(optional) format to display in [json|yaml]")
}

const templatePathModelJSON = "json_models/constructs/v1beta1/model.json"
const templatePathDesignJSON = "json_models/constructs/v1beta1/design.json"
const templatePathComponentJSON = "json_models/constructs/v1beta1/component.json"
const templatePathConnectionJSON = "json_models/constructs/v1beta1/connection.json"
const templatePathRelathionshipJSON = "json_models/constructs/v1alpha3/relationship.json"

func readTemplate(templatePath string) ([]byte, error) {
	return schemas.Schemas.ReadFile(templatePath)
}
