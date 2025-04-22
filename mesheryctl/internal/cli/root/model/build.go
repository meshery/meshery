package model

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitOci "github.com/layer5io/meshkit/models/oci"
	"github.com/meshery/schemas"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xeipuuv/gojsonschema"
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

		if err := buildModelValidateModelOverSchema(folder); err != nil {
			return err
		}

		utils.Log.Infof("Building meshery model from path %s", folder)
		img, errBuildImage := meshkitOci.BuildImage(folder)
		if errBuildImage != nil {
			return ErrModelBuild(errBuildImage)
		}

		// modelFolder := buildModelCompileFolderName(path, modelName, "")

		// Save OCI artifact into a tar file
		tarfileName := filepath.Join(folder, "model.tar")
		if err := meshkitOci.SaveOCIArtifact(img, tarfileName, modelName); err != nil {
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

func buildModelValidateModelOverSchema(folder string) error {
	utils.Log.Infof("Validating meshery model over schema from path %s", folder)
	// TODO determine format (json, yaml, csv)
	modelFile := filepath.Join(folder, "model.json")

	utils.Log.Infof("modelFile is %s", modelFile)

	utils.Log.Debug("Creating temp folder")
	tempFolder, err := createTempFolder()
	if err != nil {
		return ErrModelBuild(
			errors.Join(
				fmt.Errorf("failed to create temp folder"),
				err,
			),
		)
	}
	utils.Log.Debugf("Created temp folder %s", tempFolder)
	defer func() {
		utils.Log.Debug("Removing temp folder")
		if err := os.RemoveAll(tempFolder); err != nil {
			utils.Log.Warnf("failed to remove temp folder %s", tempFolder)
		} else {
			utils.Log.Debugf("Removed temp folder %s", tempFolder)
		}
	}()

	if err := copyEmbeddedDirToLocalFolder("schemas", schemas.Schemas, tempFolder); err != nil {
		return ErrModelBuild(
			errors.Join(
				fmt.Errorf("error copying files from embedded schemas to temp folder"),
				err,
			),
		)
	}

	modelConstruct := "constructs/v1beta1/model/model.json"
	modelSchema := filepath.Join(tempFolder, modelConstruct)
	// abs path is necessary to be able to resolve relative refs in schema
	modelSchemaAbsPath, errAbs := filepath.Abs(modelSchema)
	if errAbs != nil {
		return ErrModelBuild(
			errors.Join(
				fmt.Errorf("error determining abs path to schema"),
				errAbs,
			),
		)
	}

	schemaLoader := gojsonschema.NewReferenceLoader(
		fmt.Sprintf("file://%s", modelSchemaAbsPath),
	)

	documentLoader := gojsonschema.NewReferenceLoader(modelFile)

	// Validate

	utils.Log.Debugf(
		"Validating file %s over schema %s",
		modelFile,
		modelConstruct,
	)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return ErrModelBuild(
			errors.Join(
				fmt.Errorf("error validating %s over schema %s", modelFile, modelConstruct),
				err,
			),
		)
	}

	if result.Valid() {
		utils.Log.Debug("✅ JSON is valid!")
	} else {
		utils.Log.Info("❌ JSON is invalid:")
		for _, err := range result.Errors() {
			utils.Log.Infof("- %s\n", err)
		}
	}

	return nil
}

// returns temp dir path
func createTempFolder() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Define the base directory
	baseDir := filepath.Join(homeDir, ".meshery")

	// Create the base directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}

	// Create a temp directory inside .meshery
	tmpDir, err := os.MkdirTemp(baseDir, "tmp-*")
	if err != nil {
		return "", err
	}

	return tmpDir, nil
}

func copyEmbeddedDirToLocalFolder(embeddedPath string, efs embed.FS, targetPath string) error {
	return fs.WalkDir(efs, embeddedPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(embeddedPath, path)
		if err != nil {
			return err
		}
		target := filepath.Join(targetPath, relPath)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		data, err := efs.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(target, data, 0644)
	})
}
