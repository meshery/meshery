package model

import (
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildModelCmd = &cobra.Command{
	Use:   "build",
	Short: "Create an OCI-compliant package from the model files",
	Long: `Create an OCI-compliant package from the model files.
Expects input to be in the format scaffolded by the model init command.
Documentation for exp model init can be found at https://docs.meshery.io/reference/mesheryctl/exp/model/init
Documentation for exp model build can be found at https://docs.meshery.io/reference/mesheryctl/exp/model/build`,
	Example: `
// Create an OCI-compliant package from the model files
mesheryctl exp model build [path/to/model/version/folder]
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return ErrModelBuildFromStrings(
				"must provide exactly one argument: path to model definition",
				" if you have more then one version in the model folder, you must provide path to particular version",
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return ErrModelBuild(err)
		}

		// validation done above that args contains exactly one argument
		path := args[0]
		utils.Log.Infof("Building meshery model from path %s", path)

		return nil
	},
}
