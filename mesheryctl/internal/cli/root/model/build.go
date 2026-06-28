package model

import (
	"fmt"
	"os"
	"strings"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/modeloci"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
)

type cmdModelBuildFlags struct {
	Path string `json:"path" validate:"dirpath"`
}

var cmdModelBuildFlagsProvided cmdModelBuildFlags

var buildModelCmd = &cobra.Command{
	Use:   "build",
	Short: "Create an OCI-compliant package from the model files",
	Long: `Create an OCI-compliant package from the model files.
Model files are taken from [path]/[model-name]/[model-version] folder.
Expects input to be in the format scaffolded by the model init command.
Find more information at: https://docs.meshery.io/reference/mesheryctl/model/build`,
	Example: `
// Create an OCI-compliant package from the model files
mesheryctl model build [model-name]
mesheryctl model build [model-name]/[model-version]
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return ErrModelBuildFromStrings(errBuildUsage)
		}

		{
			err := mesheryctlflags.ValidateCmdFlags(cmd, &cmdModelBuildFlagsProvided)
			if err != nil {
				return err
			}
		}

		inputParam := args[0]
		name := ""
		version := ""

		{
			parts := strings.Split(inputParam, "/")
			// input param is supposed to be [model-name] or [model-name]/[model-version]
			// since len(args) is validated to be 1 then parts will have at least one element
			if len(parts) > 2 {
				return ErrModelBuildFromStrings(errBuildUsage)
			}
			name, version = modeloci.ParseModelInput(inputParam)
		}

		// do not validate model name, path and version,
		// as their purpose is to combine into the folder
		// only check if they are not empty and combined folder exists

		// validate name is not empty, version could be empty
		if name == "" {
			return ErrModelBuildFromStrings(errBuildUsage)
		}

		folder := modeloci.CompileFolderName(cmdModelBuildFlagsProvided.Path, name, version)
		// check if combined folder exists
		{
			// if folder does not exist return with error
			_, err := os.Stat(folder)
			if os.IsNotExist(err) {
				return ErrModelBuildFromStrings(
					errBuildUsage,
					fmt.Sprintf(
						errBuildFolderNotFound,
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
					errBuildUsage,
					errBuildMultiVersionNotSupported,
				)
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// validation (if any) is done in PreRunE (so args certainly has one element)
		name, version := modeloci.ParseModelInput(args[0])
		folder := modeloci.CompileFolderName(cmdModelBuildFlagsProvided.Path, name, version)

		utils.Log.Infof("Building meshery model from path %s", folder)
		artifactPath, err := modeloci.BuildModelOCIArtifact(cmdModelBuildFlagsProvided.Path, "", name, version)
		if err != nil {
			return ErrModelBuild(err)
		}
		utils.Log.Infof("Saving OCI artifact as %s", artifactPath)

		return nil
	},
}

func init() {
	buildModelCmd.Flags().StringVarP(&cmdModelBuildFlagsProvided.Path, "path", "p", ".", "(optional) target directory to get model from (default: current dir)")
}
