package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitOci "github.com/layer5io/meshkit/models/oci"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildModelCmd = &cobra.Command{
	Use:   "build",
	Short: "Create an OCI-compliant package from the model files",
	Long: `Create an OCI-compliant package from the model files.
Model files are taken from [path]/[model-name]/[version] folder.
Expects input to be in the format scaffolded by the model init command.
Documentation for exp model init can be found at https://docs.meshery.io/reference/mesheryctl/exp/model/init
Documentation for exp model build can be found at https://docs.meshery.io/reference/mesheryctl/exp/model/build`,
	Example: `
// Create an OCI-compliant package from the model files
mesheryctl exp model build [model-name] --version [version]
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

		// validate version is not empty
		if version == "" {
			return ErrModelBuildFromStrings("--version is empty")
		}

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

		// TODO validation over schema
		// https://github.com/meshery/schemas/pull/306

		utils.Log.Infof("Building meshery model from path %s", folder)
		img, errBuildImage := meshkitOci.BuildImage(folder)
		if errBuildImage != nil {
			return ErrModelBuild(errBuildImage)
		}

		// Save OCI artifact into a tar file under current folder
		imageName := buildModelCompileImageName(modelName, version)

		utils.Log.Infof("Saving OCI artifact as %s", imageName)
		if err := meshkitOci.SaveOCIArtifact(img, imageName, modelName); err != nil {
			return ErrModelBuild(err)
		}

		return nil
	},
}

func init() {
	buildModelCmd.Flags().StringP("path", "p", ".", "(optional) target directory to get model from (default: current dir)")
	// TODO make optional (if not specified look inside the model directory and take subfolder if only one subfolder, fail if few)
	buildModelCmd.Flags().StringP("version", "", "", "(mandatory) model version")
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

func buildModelCompileImageName(modelName string, version string) string {
	return fmt.Sprintf(
		"%s-%s.%s",
		modelName,
		strings.ReplaceAll(version, ".", "-"),
		"tar",
	)
}
