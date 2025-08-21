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
Documentation for model build can be found at https://docs.meshery.io/reference/mesheryctl/model/build`,
	Example: `
// Create an OCI-compliant package from the model files
mesheryctl model build [model-name]
mesheryctl model build [model-name]/[model-version]
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Usage:\nmesheryctl model build [model-name]\nor\nmesheryctl model build [model-name]/[model-version]\n\nRun 'mesheryctl model build --help' to see detailed help message"
		if len(args) != 1 {
			return ErrModelBuildFromStrings(errMsg)
		}

		inputParam := args[0]
		path, _ := cmd.Flags().GetString("path")
		name := ""
		version := ""

		{
			parts := strings.Split(inputParam, "/")
			// input param is supposed to be [model-name] or [model-name]/[model-version]
			// since len(args) is validated to be 1 then parts will have at least one element
			if len(parts) > 2 {
				return ErrModelBuildFromStrings(errMsg)
			}
			name, version = buildModelParseModelInput(inputParam)
		}

		// do not validate model name, path and version,
		// as their purpose is to combine into the folder
		// only check if they are not empty and combined folder exists

		// validate name is not empty, version could be empty
		if name == "" {
			return ErrModelBuildFromStrings(errMsg)
		}

		folder := buildModelCompileFolderName(path, name, version)
		// check if combined folder exists
		{
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

		// check that if only model name specified then model folder contains only one subfolder
		if version == "" {
			buildModelHasExactlyOneSubfolder := func(dir string) bool {
				entries, err := os.ReadDir(dir)
				if err != nil {
					return false
				}

				count := 0
				for _, entry := range entries {
					if entry.IsDir() {
						count++
					} else {
						// If there is a file, it's not exactly one subfolder
						return false
					}
				}

				return count == 1
			}

			if !buildModelHasExactlyOneSubfolder(folder) {
				return ErrModelBuildFromStrings(
					errMsg,
					"\nCommand does not support multiple versions build under one image",
				)
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return ErrModelBuild(err)
		}

		// validation (if any) is done in PreRunE (so args certainly has one element)
		name, version := buildModelParseModelInput(args[0])
		path, _ := cmd.Flags().GetString("path")

		// validation done above that args contains exactly one argument
		folder := buildModelCompileFolderName(path, name, version)

		utils.Log.Infof("Building meshery model from path %s", folder)
		img, errBuildImage := meshkitOci.BuildImage(folder)
		if errBuildImage != nil {
			return ErrModelBuild(errBuildImage)
		}

		// Save OCI artifact into a tar file under current folder
		imageName := buildModelCompileImageName(name, version, "tar")

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

// parseModelInput parses the user input into model name and version.
// The input can be `model-name` or `model-name/version`.
func buildModelParseModelInput(input string) (name, version string) {
	parts := strings.Split(input, "/")
	name = parts[0]
	if len(parts) > 1 {
		version = parts[1]
	}
	return
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

func buildModelCompileImageName(name string, version string, extension string) string {
	if version != "" {
		return fmt.Sprintf(
			"%s-%s.%s",
			name,
			strings.ReplaceAll(version, ".", "-"),
			extension,
		)
	}
	return fmt.Sprintf(
		"%s.%s",
		name,
		extension,
	)
}
