package model

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	meshkitOci "github.com/layer5io/meshkit/models/oci"
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
mesheryctl exp model build [model-name]
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return ErrModelBuildFromStrings("must provide only one argument: model name")
		}
		// do not validate model name, path and version,
		// as their purpose is to combine into the folder
		// only check if combined folder exists
		modelName := args[0]
		path, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")

		{
			folder := buildModelCompileFolderName(path, modelName, version)
			// if folder does not exist return with error
			_, err := os.Stat(folder)
			if os.IsNotExist(err) {
				return ErrModelBuildFromStrings(
					fmt.Sprintf(
						"folder %s does not exist",
						folder,
					),
				)
			}
		}
		// TODO should we validate if the directory has a valid meshery model structure?

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return ErrModelBuild(err)
		}

		// validation (if any) is done in PreRunE
		modelName := args[0]
		path, _ := cmd.Flags().GetString("path")
		version, _ := cmd.Flags().GetString("version")

		// validation done above that args contains exactly one argument
		folder := buildModelCompileFolderName(path, modelName, version)
		utils.Log.Infof("Building meshery model from path %s", folder)

		img, err := meshkitOci.BuildImage(folder)
		if err != nil {
			return ErrModelBuild(err)
		}

		// modelFolder := buildModelCompileFolderName(path, modelName, "")

		// Save OCI artifact into a tar file
		tarfileName := filepath.Join(folder, "model.tar")
		err = meshkitOci.SaveOCIArtifact(img, tarfileName, modelName)
		if err != nil {
			return ErrModelBuild(err)
		}

		return nil
	},
}

func init() {
	buildModelCmd.Flags().StringP("path", "p", ".", "(optional) target directory to get model from (default: current dir)")
	buildModelCmd.Flags().StringP("version", "", "", "(optional) model version (if not specified, cmd builds all version from model folder)")
}

func buildModelCompileFolderName(path string, modelName string, version string) string {
	dirParts := make([]string, 0, 3)
	if path != "" {
		dirParts = append(dirParts, path)
	}
	dirParts = append(dirParts, modelName)
	if version != "" {
		dirParts = append(dirParts, version)
	}
	return filepath.Join(dirParts...)
}
