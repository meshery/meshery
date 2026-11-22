package model

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitOci "github.com/meshery/meshkit/models/oci"
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
When only [model-name] is provided, the model directory must contain exactly one version directory.
Root-level files such as README.md, .gitignore, and LICENSE do not affect this validation.
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
			name, version = buildModelParseModelInput(inputParam)
		}

		// do not validate model name, path and version,
		// as their purpose is to combine into the folder
		// only check if they are not empty and combined folder exists

		// validate name is not empty, version could be empty
		if name == "" {
			return ErrModelBuildFromStrings(errBuildUsage)
		}

		folder := buildModelCompileFolderName(cmdModelBuildFlagsProvided.Path, name, version)
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
					if entry.IsDir() && mesheryctlflags.IsSemver(entry.Name()) {
						count++
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
		name, version := buildModelParseModelInput(args[0])

		// validation done above that args contains exactly one argument
		folder := buildModelCompileFolderName(cmdModelBuildFlagsProvided.Path, name, version)

		utils.Log.Infof("Building meshery model from path %s", folder)

		buildSource, cleanup, err := prepareModelBuildSource(folder)
		if err != nil {
			err = ErrModelBuild(err)
			utils.Log.Error(err)
			return err
		}
		defer cleanup()

		img, errBuildImage := meshkitOci.BuildImage(buildSource)
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
	buildModelCmd.Flags().StringVarP(&cmdModelBuildFlagsProvided.Path, "path", "p", ".", "(optional) target directory to get model from (default: current dir)")
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

// hasGitMetadata checks whether .git exists anywhere within src (as a file or directory).
func hasGitMetadata(src string) (bool, error) {
	hasGit := false
	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Name() == ".git" {
			hasGit = true
			if d.IsDir() {
				return fs.SkipDir
			}
		}
		return nil
	})
	return hasGit, err
}

// prepareModelBuildSource prepares a build source directory that excludes any Git repository metadata (.git).
// If no .git metadata is present, it returns the original folder with a no-op cleanup to avoid copying overhead.
// If .git metadata is present, it stages a clean copy into an isolated temporary directory outside the source tree.
func prepareModelBuildSource(folder string) (string, func(), error) {
	cleanFolder := filepath.Clean(folder)
	hasGit, err := hasGitMetadata(cleanFolder)
	if err != nil {
		return "", func() {}, err
	}
	if !hasGit {
		return cleanFolder, func() {}, nil
	}

	tmpDir, err := os.MkdirTemp("", "meshery-model-build-*")
	if err != nil {
		return "", func() {}, err
	}
	cleanup := func() {
		_ = filepath.WalkDir(tmpDir, func(p string, d fs.DirEntry, err error) error {
			if err == nil && d.IsDir() {
				_ = os.Chmod(p, 0700)
			}
			return nil
		})
		_ = os.RemoveAll(tmpDir)
	}

	if err := copyModelExcludingGit(cleanFolder, tmpDir); err != nil {
		cleanup()
		return "", func() {}, err
	}

	return tmpDir, cleanup, nil
}

type dirPerm struct {
	path string
	perm os.FileMode
}

func isSubpathOrEqual(p, target string) bool {
	cleanTarget := filepath.Clean(target)
	cleanP := filepath.Clean(p)
	if cleanP == cleanTarget {
		return true
	}
	if cleanTarget == string(filepath.Separator) {
		return true
	}
	return strings.HasPrefix(cleanP, cleanTarget+string(filepath.Separator))
}

// copyModelExcludingGit recursively copies src into dst while omitting any .git file or directory.
func copyModelExcludingGit(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source path is not a directory: %s", src)
	}

	absDst, err := filepath.Abs(dst)
	if err != nil {
		return err
	}
	realDst, err := filepath.EvalSymlinks(absDst)
	if err != nil {
		realDst = absDst
	}

	var dirsToRestore []dirPerm

	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		// Prevent recursion if destination is inside source (e.g. TMPDIR set inside model source)
		if isSubpathOrEqual(absPath, absDst) || isSubpathOrEqual(absPath, realDst) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if realPath, err := filepath.EvalSymlinks(absPath); err == nil {
			if isSubpathOrEqual(realPath, absDst) || isSubpathOrEqual(realPath, realDst) {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
		}

		if d.Name() == ".git" {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		// Prevent path traversal
		if strings.HasPrefix(relPath, "..") || filepath.IsAbs(relPath) {
			return fmt.Errorf("invalid relative path: %s", relPath)
		}

		// Ensure no path segment is ".git"
		segments := strings.Split(filepath.ToSlash(relPath), "/")
		for _, seg := range segments {
			if seg == ".git" {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
		}

		targetPath := filepath.Join(dst, relPath)

		info, err := d.Info()
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Ensure owner write permission (0700) during staging so child items can be created,
			// even if the source directory has read-only permissions (e.g. 0555).
			dirsToRestore = append(dirsToRestore, dirPerm{path: targetPath, perm: info.Mode().Perm()})
			return os.MkdirAll(targetPath, info.Mode().Perm()|0700)
		}

		if d.Type()&fs.ModeSymlink != 0 {
			linkTarget, err := os.Readlink(path)
			if err != nil {
				return err
			}
			cleanTarget := filepath.Clean(linkTarget)
			targetSlash := "/" + filepath.ToSlash(cleanTarget) + "/"
			if cleanTarget == ".git" || strings.Contains(targetSlash, "/.git/") {
				return nil
			}
			return os.Symlink(linkTarget, targetPath)
		}

		if d.Type().IsRegular() {
			return copyModelFile(path, targetPath, info.Mode().Perm())
		}

		return fmt.Errorf("unsupported file type for %s: %s", path, d.Type().String())
	})
	if err != nil {
		return err
	}

	// Restore original directory permissions in reverse order (deepest directories first)
	for i := len(dirsToRestore) - 1; i >= 0; i-- {
		dp := dirsToRestore[i]
		if chmodErr := os.Chmod(dp.path, dp.perm); chmodErr != nil {
			return chmodErr
		}
	}

	return nil
}

func copyModelFile(srcFile, dstFile string, perm os.FileMode) (retErr error) {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := src.Close(); closeErr != nil {
			retErr = errors.Join(retErr, closeErr)
		}
	}()

	if err := os.MkdirAll(filepath.Dir(dstFile), 0755); err != nil {
		return err
	}

	dst, err := os.OpenFile(dstFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dst.Close(); closeErr != nil {
			retErr = errors.Join(retErr, closeErr)
		}
	}()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return retErr
}
