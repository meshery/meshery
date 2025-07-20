package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitOci "github.com/meshery/meshkit/models/oci"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var buildModelCmd = &cobra.Command{
	Use:   "build",
	Short: "Create an OCI-compliant package from the model files",
	Long: `Create an OCI-compliant package from the model files.
Model files are taken from [path]/[model-name]/[model-version] folder.
Expects input to be in the format scaffolded by the model init command.
Documentation for exp model and subcommands can be found at https://docs.meshery.io/reference/mesheryctl/exp/model`,
	Example: `
// Create an OCI-compliant package from the model files
mesheryctl exp model build [model-name]/[model-version]
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp model build [model-name]/[model-version]\nRun 'mesheryctl exp model build --help' to see detailed help message"
		if len(args) != 1 {
			return ErrModelBuildFromStrings(errMsg)
		}

		modelNameVersion := args[0]
		path, _ := cmd.Flags().GetString("path")
		name := ""
		version := ""

		{
			parts := strings.Split(modelNameVersion, "/")
			// input param is supposed to be [model-name]/[version]
			if len(parts) != 2 {
				return ErrModelBuildFromStrings(errMsg)
			}
			name = parts[0]
			version = parts[1]
		}

		// do not validate model name, path and version,
		// as their purpose is to combine into the folder
		// only check if they y are not empty and combined folder exists

		// validate name and version are not empty
		if name == "" || version == "" {
			return ErrModelBuildFromStrings(errMsg)
		}

		// check if combined folder exists
		{
			folder := buildModelCompileFolderName(path, name, version)
			// if folder does not exist return with error
			_, err := os.Stat(folder)
			if os.IsNotExist(err) {
				return ErrModelBuildFromStrings(
					errMsg,
					fmt.Sprintf(
						"\nfolder %s does not exist",
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
		modelNameVersion := args[0]
		parts := strings.Split(modelNameVersion, "/")
		// it was validated that parts is exactly two elements slice in preRunE
		name := parts[0]
		version := parts[1]
		path, _ := cmd.Flags().GetString("path")

		// validation done above that args contains exactly one argument
		folder := buildModelCompileFolderName(path, name, version)

		// TODO validation over schema
		// https://github.com/meshery/schemas/pull/306

		utils.Log.Infof("Building meshery model from path %s", folder)
		img, errBuildImage := meshkitOci.BuildImage(folder)
		if errBuildImage != nil {
			return ErrModelBuild(errBuildImage)
		}

		// Save OCI artifact into a tar file under current folder
		imageName := buildModelCompileImageName(name, version)

		utils.Log.Infof("Saving OCI artifact as %s", imageName)
		if err := meshkitOci.SaveOCIArtifact(img, imageName, name); err != nil {
			return ErrModelBuild(err)
		}

		return nil
	},
}

func init() {
	buildModelCmd.Flags().StringP("path", "p", ".", "(optional) target directory to get model from (default: current dir)")
}

func buildModelCompileFolderName(path string, name string, version string) string {
	dirParts := make([]string, 0, 3)
	if path != "" {
		dirParts = append(dirParts, path)
	}
	dirParts = append(dirParts, name)
	if version != "" {
		dirParts = append(dirParts, version)
	}
	return filepath.Join(dirParts...)
}

func buildModelCompileImageName(name string, version string) string {
	return fmt.Sprintf(
		"%s-%s.%s",
		name,
		strings.ReplaceAll(version, ".", "-"),
		"tar",
	)
}
