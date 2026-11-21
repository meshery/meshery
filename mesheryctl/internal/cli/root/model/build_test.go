package model

import (
	archiveTar "archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-containerregistry/pkg/v1/tarball"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitOci "github.com/meshery/meshkit/models/oci"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestModelBuild(t *testing.T) {
	// Shared cleanup function for test directories and artifacts
	cleanupTestArtifacts := func(dirs []string) {
		for _, dir := range dirs {
			_ = os.RemoveAll(dir)
			_ = os.RemoveAll(dir + ".tar")
		}
	}

	// Test constants
	const (
		buildTestDynamoController    = "test-case-model-build-aws-dynamodb-controller"
		buildTestDynamoControllerGbx = "test-case-model-build-aws-dynamodb-controller-gbxter34"
		buildTestEC2Controller       = "aws-ec2-controller"
		buildTestNoVersionDir        = "test-case-model-build-no-versions"
		buildTestTwoVersions         = "test-case-model-build-two-versions"
		buildTestVersion             = "v0.1.0"
		buildTestNonExistentFolder   = "not_existing_folder"
	)

	// Clean up any existing test directories before running tests
	cleanupDirs := []string{
		buildTestDynamoController,
		buildTestDynamoControllerGbx,
		buildTestNoVersionDir,
		buildTestTwoVersions,
	}
	// Clean up all test artifacts from previous runs
	cleanupTestArtifacts(cleanupDirs)
	// Register cleanup for after test completion
	t.Cleanup(func() {
		cleanupTestArtifacts(cleanupDirs)
	})

	utils.SetupContextEnv(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// Helper function to create a fresh ModelCmd with all original properties
	createFreshModelCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:     ModelCmd.Use,
			Short:   ModelCmd.Short,
			Long:    ModelCmd.Long,
			Example: ModelCmd.Example,
			Args:    ModelCmd.Args,
			RunE:    ModelCmd.RunE,
		}
		// Copy all flags from the original ModelCmd
		ModelCmd.Flags().VisitAll(func(flag *pflag.Flag) {
			cmd.Flags().AddFlag(flag)
		})
		return cmd
	}

	// Helper function to create fresh commands for build tests
	createFreshCommands := func() *cobra.Command {
		// Create fresh build command
		freshBuildCmd := &cobra.Command{
			Use:     buildModelCmd.Use,
			Short:   buildModelCmd.Short,
			Long:    buildModelCmd.Long,
			Example: buildModelCmd.Example,
			PreRunE: buildModelCmd.PreRunE,
			RunE:    buildModelCmd.RunE,
		}
		// Copy all flags from the original build command
		buildModelCmd.Flags().VisitAll(func(flag *pflag.Flag) {
			freshBuildCmd.Flags().AddFlag(flag)
		})

		// Use the shared helper for ModelCmd
		cmd := createFreshModelCmd()
		mesheryctlflags.InitValidators(freshBuildCmd)
		cmd.AddCommand(freshBuildCmd)
		return cmd
	}

	setupHookModelInit := func(modelInitArgs ...string) func() {
		return func() {
			// Use the shared helper to create a fresh ModelCmd
			cmd := createFreshModelCmd()
			cmd.AddCommand(initModelCmd)
			mesheryctlflags.InitValidators(initModelCmd)
			cmd.SetArgs(modelInitArgs)
			buff := utils.SetupMeshkitLoggerTesting(t, false)
			cmd.SetOut(buff)
			if err := cmd.Execute(); err != nil {
				t.Fatal(err)
			}
		}
	}
	cleanUpHookRemoveDirsAndFiles := func(dirs ...string) func() {
		return func() {
			cleanupTestArtifacts(dirs)
			t.Log("removed created dirs and files")
		}
	}
	setupHookCreateRootFiles := func(modelDir string, files ...string) func() {
		return func() {
			for _, file := range files {
				filePath := filepath.Join(modelDir, file)
				if err := os.WriteFile(filePath, []byte("sample content"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
	setupHookCreateRootDirectories := func(modelDir string, dirs ...string) func() {
		return func() {
			for _, dir := range dirs {
				dirPath := filepath.Join(modelDir, dir)
				if err := os.MkdirAll(dirPath, 0o755); err != nil {
					t.Fatal(err)
				}
			}
		}
	}

	tests := []struct {
		Name             string
		Args             []string
		SetupHooks       []func()
		ExpectError      bool
		ExpectedResponse string
		ExpectedFiles    []string
		CleanupHooks     []func()
		IsOutputGolden   bool
		ExpectedError    error
	}{
		{
			Name:             "given name and version when model build then model is built",
			Args:             []string{"build", "test-case-model-build-aws-lambda-controller/v0.1.0"},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-version.golden",
			ExpectedFiles: []string{
				"test-case-model-build-aws-lambda-controller-v0-1-0.tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", "test-case-model-build-aws-lambda-controller", "--version", "v0.1.0"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					"test-case-model-build-aws-lambda-controller",
					"test-case-model-build-aws-lambda-controller-v0-1-0.tar",
				),
			},
		},
		{
			Name:             "given name only when model build then model is built",
			Args:             []string{"build", buildTestDynamoController},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				buildTestDynamoController + ".tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoController, "--version", buildTestVersion),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoController,
					buildTestDynamoController+".tar",
				),
			},
		},
		{
			Name:             "given name only when model build then model is built",
			Args:             []string{"build", buildTestDynamoController + "/"},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				buildTestDynamoController + ".tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoController, "--version", buildTestVersion),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoController,
					buildTestDynamoController+".tar",
				),
			},
		},
		{
			Name:             "given name only with README.md when model build then model is built",
			Args:             []string{"build", buildTestDynamoController},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				buildTestDynamoController + ".tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoController, "--version", buildTestVersion),
				setupHookCreateRootFiles(buildTestDynamoController, "README.md"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoController,
					buildTestDynamoController+".tar",
				),
			},
		},
		{
			Name:             "given name only with .gitignore when model build then model is built",
			Args:             []string{"build", buildTestDynamoController},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				buildTestDynamoController + ".tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoController, "--version", buildTestVersion),
				setupHookCreateRootFiles(buildTestDynamoController, ".gitignore"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoController,
					buildTestDynamoController+".tar",
				),
			},
		},
		{
			Name:             "given name only with README.md and .gitignore when model build then model is built",
			Args:             []string{"build", buildTestDynamoController},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				buildTestDynamoController + ".tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoController, "--version", buildTestVersion),
				setupHookCreateRootFiles(buildTestDynamoController, "README.md", ".gitignore"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoController,
					buildTestDynamoController+".tar",
				),
			},
		},
		{
			Name:             "given name only with LICENSE when model build then model is built",
			Args:             []string{"build", buildTestDynamoController},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				buildTestDynamoController + ".tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoController, "--version", buildTestVersion),
				setupHookCreateRootFiles(buildTestDynamoController, "LICENSE"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoController,
					buildTestDynamoController+".tar",
				),
			},
		},
		{
			Name:             "given name only with .git directory when model build then model is built",
			Args:             []string{"build", buildTestDynamoController},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				buildTestDynamoController + ".tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoController, "--version", buildTestVersion),
				setupHookCreateRootDirectories(buildTestDynamoController, ".git"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoController,
					buildTestDynamoController+".tar",
				),
			},
		},
		{
			Name:             "given name only with .github directory when model build then model is built",
			Args:             []string{"build", buildTestDynamoController},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				buildTestDynamoController + ".tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoController, "--version", buildTestVersion),
				setupHookCreateRootDirectories(buildTestDynamoController, ".github"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoController,
					buildTestDynamoController+".tar",
				),
			},
		},
		{
			Name:             "given name only with docs directory when model build then model is built",
			Args:             []string{"build", buildTestDynamoController},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				buildTestDynamoController + ".tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoController, "--version", buildTestVersion),
				setupHookCreateRootDirectories(buildTestDynamoController, "docs"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoController,
					buildTestDynamoController+".tar",
				),
			},
		},
		{
			Name:             "given name only with .git and docs directory when model build then model is built",
			Args:             []string{"build", buildTestDynamoController},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				buildTestDynamoController + ".tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoController, "--version", buildTestVersion),
				setupHookCreateRootDirectories(buildTestDynamoController, ".git", "docs"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoController,
					buildTestDynamoController+".tar",
				),
			},
		},
		{
			Name:             "given name only with README.md, docs and valid version directory when model build then model is built",
			Args:             []string{"build", buildTestDynamoController},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				buildTestDynamoController + ".tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoController, "--version", "v1.0.0"),
				setupHookCreateRootFiles(buildTestDynamoController, "README.md"),
				setupHookCreateRootDirectories(buildTestDynamoController, "docs"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoController,
					buildTestDynamoController+".tar",
				),
			},
		},
		{
			Name:             "given name only with prerelease version directory when model build then model is built",
			Args:             []string{"build", buildTestDynamoController},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				buildTestDynamoController + ".tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoController, "--version", "v1.0.0-alpha"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoController,
					buildTestDynamoController+".tar",
				),
			},
		},
		{
			Name:             "given name only with build metadata version directory when model build then model is built",
			Args:             []string{"build", buildTestDynamoController},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				buildTestDynamoController + ".tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoController, "--version", "v1.0.0+20130313144700"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoController,
					buildTestDynamoController+".tar",
				),
			},
		},
		{
			Name:             "given no params when model build then throw error",
			Args:             []string{"build"},
			ExpectError:      true,
			ExpectedResponse: "",
			IsOutputGolden:   false,
			ExpectedError:    ErrModelBuildFromStrings(errBuildUsage),
		},
		{
			Name:             "given wrong input param format when model build then throw error",
			Args:             []string{"build", buildTestEC2Controller + "/" + buildTestVersion + "/smthelse"},
			ExpectError:      true,
			ExpectedResponse: "",
			IsOutputGolden:   false,
			ExpectedError:    ErrModelBuildFromStrings(errBuildUsage),
		},
		{
			Name:             "given name and multiple versions when model build then throw error",
			Args:             []string{"build", buildTestDynamoControllerGbx},
			ExpectError:      true,
			ExpectedResponse: "",
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestDynamoControllerGbx, "--version", buildTestVersion),
				setupHookModelInit("init", buildTestDynamoControllerGbx, "--version", "v0.1.1"),
				setupHookModelInit("init", buildTestDynamoControllerGbx, "--version", "v0.1.2"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestDynamoControllerGbx,
				),
			},
			IsOutputGolden: false,
			ExpectedError:  ErrModelBuildFromStrings(errBuildUsage, errBuildMultiVersionNotSupported),
		},
		{
			Name:             "given name and two version directories when model build then throw error",
			Args:             []string{"build", buildTestTwoVersions},
			ExpectError:      true,
			ExpectedResponse: "",
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestTwoVersions, "--version", "v1.0.0"),
				setupHookModelInit("init", buildTestTwoVersions, "--version", "v2.0.0"),
				setupHookCreateRootFiles(buildTestTwoVersions, "README.md"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestTwoVersions,
				),
			},
			IsOutputGolden: false,
			ExpectedError:  ErrModelBuildFromStrings(errBuildUsage, errBuildMultiVersionNotSupported),
		},
		{
			Name:             "given name, two version directories and .git directory when model build then throw error",
			Args:             []string{"build", buildTestTwoVersions},
			ExpectError:      true,
			ExpectedResponse: "",
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestTwoVersions, "--version", "v1.0.0"),
				setupHookModelInit("init", buildTestTwoVersions, "--version", "v2.0.0"),
				setupHookCreateRootDirectories(buildTestTwoVersions, ".git"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestTwoVersions,
				),
			},
			IsOutputGolden: false,
			ExpectedError:  ErrModelBuildFromStrings(errBuildUsage, errBuildMultiVersionNotSupported),
		},
		{
			Name:             "given name, two version directories and docs directory when model build then throw error",
			Args:             []string{"build", buildTestTwoVersions},
			ExpectError:      true,
			ExpectedResponse: "",
			SetupHooks: []func(){
				setupHookModelInit("init", buildTestTwoVersions, "--version", "v1.0.0"),
				setupHookModelInit("init", buildTestTwoVersions, "--version", "v2.0.0"),
				setupHookCreateRootDirectories(buildTestTwoVersions, "docs"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestTwoVersions,
				),
			},
			IsOutputGolden: false,
			ExpectedError:  ErrModelBuildFromStrings(errBuildUsage, errBuildMultiVersionNotSupported),
		},
		{
			Name:             "given name and no version directories when model build then throw error",
			Args:             []string{"build", buildTestNoVersionDir},
			ExpectError:      true,
			ExpectedResponse: "",
			SetupHooks: []func(){
				func() {
					if err := os.MkdirAll(buildTestNoVersionDir, 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(buildTestNoVersionDir, "README.md"), []byte("sample docs"), 0o644); err != nil {
						t.Fatal(err)
					}
				},
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestNoVersionDir,
				),
			},
			IsOutputGolden: false,
			ExpectedError:  ErrModelBuildFromStrings(errBuildUsage, errBuildMultiVersionNotSupported),
		},
		{
			Name:             "given name and only hidden directory when model build then throw error",
			Args:             []string{"build", buildTestNoVersionDir},
			ExpectError:      true,
			ExpectedResponse: "",
			SetupHooks: []func(){
				func() {
					if err := os.MkdirAll(filepath.Join(buildTestNoVersionDir, ".git"), 0o755); err != nil {
						t.Fatal(err)
					}
				},
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestNoVersionDir,
				),
			},
			IsOutputGolden: false,
			ExpectedError:  ErrModelBuildFromStrings(errBuildUsage, errBuildMultiVersionNotSupported),
		},
		{
			Name:             "given name, docs directory with guide.md and README.md without valid version directory when model build then throw error",
			Args:             []string{"build", buildTestNoVersionDir},
			ExpectError:      true,
			ExpectedResponse: "",
			SetupHooks: []func(){
				func() {
					docsDir := filepath.Join(buildTestNoVersionDir, "docs")
					if err := os.MkdirAll(docsDir, 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(docsDir, "guide.md"), []byte("sample guide content"), 0o644); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(buildTestNoVersionDir, "README.md"), []byte("sample docs"), 0o644); err != nil {
						t.Fatal(err)
					}
				},
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestNoVersionDir,
				),
			},
			IsOutputGolden: false,
			ExpectedError:  ErrModelBuildFromStrings(errBuildUsage, errBuildMultiVersionNotSupported),
		},
		{
			Name:             "given name and only regular files without valid version directory when model build then throw error",
			Args:             []string{"build", buildTestNoVersionDir},
			ExpectError:      true,
			ExpectedResponse: "",
			SetupHooks: []func(){
				func() {
					if err := os.MkdirAll(buildTestNoVersionDir, 0o755); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(buildTestNoVersionDir, "README.md"), []byte("sample docs"), 0o644); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(buildTestNoVersionDir, "LICENSE"), []byte("sample license"), 0o644); err != nil {
						t.Fatal(err)
					}
				},
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					buildTestNoVersionDir,
				),
			},
			IsOutputGolden: false,
			ExpectedError:  ErrModelBuildFromStrings(errBuildUsage, errBuildMultiVersionNotSupported),
		},
		{
			Name:             "given folder does not exist when model build then throw error",
			Args:             []string{"build", buildTestEC2Controller + "/" + buildTestVersion, "--path", "./" + buildTestNonExistentFolder},
			ExpectError:      true,
			ExpectedResponse: "",
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --path './%s': directory does not exist", buildTestNonExistentFolder)),
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			// Clean up any lingering test artifacts before this subtest starts
			cleanupTestArtifacts(cleanupDirs)
			cmdModelBuildFlagsProvided.Path = "."

			if len(tc.CleanupHooks) > 0 {
				for _, cleanupHook := range tc.CleanupHooks {
					defer cleanupHook()
				}
			}
			if len(tc.SetupHooks) > 0 {
				for _, setupHook := range tc.SetupHooks {
					setupHook()
				}
			}
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tc.ExpectedResponse, testdataDir)
			buff := utils.SetupMeshkitLoggerTesting(t, false)
			// Create fresh commands using helper function
			cmd := createFreshCommands()
			cmd.SetArgs(tc.Args)
			cmd.SetOut(buff)
			err := cmd.Execute()
			if err != nil {
				// if we're supposed to get an error
				if tc.ExpectError {
					if tc.IsOutputGolden {
						expectedResponse := golden.Load()

						utils.Equals(t, expectedResponse, err.Error())
						return
					}
					utils.AssertMeshkitErrorsEqual(t, err, tc.ExpectedError)
					return

				}
				t.Fatal(err)
			}
			// response being printed in console
			actualResponse := buff.String()

			expectedResponse := golden.Load()
			assert.Equal(t, expectedResponse, actualResponse)

			if len(tc.ExpectedFiles) > 0 {
				for _, file := range tc.ExpectedFiles {
					info, err := os.Stat(file)
					assert.NoError(t, err)
					if err == nil {
						assert.False(t, info.IsDir())
					}
				}
			}
		})
	}
}

func TestModelBuild_Security_ExcludesGitMetadata(t *testing.T) {
	utils.SetupContextEnv(t)
	cmdModelBuildFlagsProvided.Path = "."

	createFreshModelCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:     ModelCmd.Use,
			Short:   ModelCmd.Short,
			Long:    ModelCmd.Long,
			Example: ModelCmd.Example,
			Args:    ModelCmd.Args,
			RunE:    ModelCmd.RunE,
		}
		ModelCmd.Flags().VisitAll(func(flag *pflag.Flag) {
			cmd.Flags().AddFlag(flag)
		})
		return cmd
	}

	createFreshBuildCmd := func() *cobra.Command {
		freshBuildCmd := &cobra.Command{
			Use:     buildModelCmd.Use,
			Short:   buildModelCmd.Short,
			Long:    buildModelCmd.Long,
			Example: buildModelCmd.Example,
			PreRunE: buildModelCmd.PreRunE,
			RunE:    buildModelCmd.RunE,
		}
		buildModelCmd.Flags().VisitAll(func(flag *pflag.Flag) {
			freshBuildCmd.Flags().AddFlag(flag)
		})
		cmd := createFreshModelCmd()
		mesheryctlflags.InitValidators(freshBuildCmd)
		cmd.AddCommand(freshBuildCmd)
		return cmd
	}

	const modelName = "test-case-model-build-security-git-model"
	tarName := modelName + ".tar"

	cleanup := func() {
		_ = os.RemoveAll(modelName)
		_ = os.RemoveAll(tarName)
	}
	cleanup()
	defer cleanup()

	// Initialize model with v1.0.0
	setupCmd := createFreshModelCmd()
	setupCmd.AddCommand(initModelCmd)
	mesheryctlflags.InitValidators(initModelCmd)
	setupCmd.SetArgs([]string{"init", modelName, "--version", "v1.0.0"})
	buffInit := utils.SetupMeshkitLoggerTesting(t, false)
	setupCmd.SetOut(buffInit)
	if err := setupCmd.Execute(); err != nil {
		t.Fatalf("failed to init test model: %v", err)
	}

	// Create root files to verify preservation
	readmeExpected := "# Security Test Model\nThis is a sample model."
	gitignoreExpected := "*.tar\n"
	licenseExpected := "Apache License 2.0"
	ciExpected := "name: CI\non: [push]"

	if err := os.WriteFile(filepath.Join(modelName, "README.md"), []byte(readmeExpected), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(modelName, ".gitignore"), []byte(gitignoreExpected), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(modelName, "LICENSE"), []byte(licenseExpected), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create .github/workflows/ci.yml
	githubDir := filepath.Join(modelName, ".github", "workflows")
	if err := os.MkdirAll(githubDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(githubDir, "ci.yml"), []byte(ciExpected), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create .git directory and config with dummy credentials
	gitDir := filepath.Join(modelName, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	dummyGitConfig := `[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
[remote "origin"]
	url = https://test-user:dummy-token@example.invalid/private/repository.git
	fetch = +refs/heads/*:refs/remotes/origin/*
`
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte(dummyGitConfig), 0o600); err != nil {
		t.Fatal(err)
	}

	// Create nested directory with a nested .git directory and config
	nestedGitDir := filepath.Join(modelName, "docs", ".git")
	if err := os.MkdirAll(nestedGitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nestedGitDir, "config"), []byte(dummyGitConfig), 0o600); err != nil {
		t.Fatal(err)
	}

	// Run mesheryctl model build
	cmd := createFreshBuildCmd()
	cmd.SetArgs([]string{"build", modelName})
	buff := utils.SetupMeshkitLoggerTesting(t, false)
	cmd.SetOut(buff)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("model build failed: %v", err)
	}

	// Verify tar exists
	info, err := os.Stat(tarName)
	assert.NoError(t, err)
	assert.False(t, info.IsDir())

	// Uncompress OCI artifact using MeshKit helper into an isolated temporary directory
	extractDir := t.TempDir()
	err = meshkitOci.UnCompressOCIArtifact(tarName, extractDir)
	assert.NoError(t, err)

	// 1. Assert .git does not exist
	_, err = os.Stat(filepath.Join(extractDir, ".git"))
	assert.True(t, os.IsNotExist(err), ".git directory must not exist in extracted artifact")

	_, err = os.Stat(filepath.Join(extractDir, ".git", "config"))
	assert.True(t, os.IsNotExist(err), ".git/config must not exist in extracted artifact")

	// 2. Assert nested .git does not exist
	_, err = os.Stat(filepath.Join(extractDir, "docs", ".git"))
	assert.True(t, os.IsNotExist(err), "nested .git directory must not exist in extracted artifact")

	// 3. Walk all extracted files to ensure no entry begins with or contains .git, and no dummy credentials leaked
	err = filepath.Walk(extractDir, func(path string, d os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(extractDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		// Ensure no path segment is ".git"
		for _, seg := range strings.Split(filepath.ToSlash(rel), "/") {
			assert.NotEqual(t, ".git", seg, "artifact entry contains .git: %s", rel)
		}

		if !d.IsDir() {
			content, err := os.ReadFile(path)
			if assert.NoError(t, err) {
				assert.NotContains(t, string(content), "dummy-token", "found credential token leaked in %s", rel)
				assert.NotContains(t, string(content), "example.invalid", "found remote URL leaked in %s", rel)
			}
		}
		return nil
	})
	assert.NoError(t, err)

	// 4. Assert model and documentation files are preserved intact
	readmeContent, err := os.ReadFile(filepath.Join(extractDir, "README.md"))
	assert.NoError(t, err)
	assert.Equal(t, readmeExpected, string(readmeContent))

	gitignoreContent, err := os.ReadFile(filepath.Join(extractDir, ".gitignore"))
	assert.NoError(t, err)
	assert.Equal(t, gitignoreExpected, string(gitignoreContent))

	licenseContent, err := os.ReadFile(filepath.Join(extractDir, "LICENSE"))
	assert.NoError(t, err)
	assert.Equal(t, licenseExpected, string(licenseContent))

	ciContent, err := os.ReadFile(filepath.Join(extractDir, ".github", "workflows", "ci.yml"))
	assert.NoError(t, err)
	assert.Equal(t, ciExpected, string(ciContent))

	// 5. Assert version directory and model.json are preserved intact
	_, err = os.Stat(filepath.Join(extractDir, "v1.0.0", "model.json"))
	assert.NoError(t, err)
}

func TestModelBuild_Staging_TmpDirInsideModelSource(t *testing.T) {
	cmdModelBuildFlagsProvided.Path = "."

	createFreshModelCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "model",
			Short: "Model management",
		}
		ModelCmd.Flags().VisitAll(func(flag *pflag.Flag) {
			cmd.Flags().AddFlag(flag)
		})
		return cmd
	}

	createFreshBuildCmd := func() *cobra.Command {
		freshBuildCmd := &cobra.Command{
			Use:     buildModelCmd.Use,
			Short:   buildModelCmd.Short,
			Long:    buildModelCmd.Long,
			Example: buildModelCmd.Example,
			PreRunE: buildModelCmd.PreRunE,
			RunE:    buildModelCmd.RunE,
		}
		buildModelCmd.Flags().VisitAll(func(flag *pflag.Flag) {
			freshBuildCmd.Flags().AddFlag(flag)
		})
		cmd := createFreshModelCmd()
		mesheryctlflags.InitValidators(freshBuildCmd)
		cmd.AddCommand(freshBuildCmd)
		return cmd
	}

	const modelName = "test-case-model-build-tmpdir-recursion"
	tarName := modelName + ".tar"

	cleanup := func() {
		_ = os.RemoveAll(modelName)
		_ = os.RemoveAll(tarName)
	}
	cleanup()
	defer cleanup()

	// Initialize model with v1.0.0
	setupCmd := createFreshModelCmd()
	setupCmd.AddCommand(initModelCmd)
	mesheryctlflags.InitValidators(initModelCmd)
	setupCmd.SetArgs([]string{"init", modelName, "--version", "v1.0.0"})
	buffInit := utils.SetupMeshkitLoggerTesting(t, false)
	setupCmd.SetOut(buffInit)
	if err := setupCmd.Execute(); err != nil {
		t.Fatalf("failed to init test model: %v", err)
	}

	// Create .git directory to trigger staging
	gitDir := filepath.Join(modelName, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte("sample git config"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Create TMPDIR inside model directory
	internalTmp := filepath.Join(modelName, "internal_tmp")
	if err := os.MkdirAll(internalTmp, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("TMPDIR", internalTmp)

	// Run mesheryctl model build
	cmd := createFreshBuildCmd()
	cmd.SetArgs([]string{"build", modelName})
	buff := utils.SetupMeshkitLoggerTesting(t, false)
	cmd.SetOut(buff)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("model build failed with TMPDIR inside model: %v", err)
	}

	// Verify tar exists
	info, err := os.Stat(tarName)
	assert.NoError(t, err)
	assert.False(t, info.IsDir())

	// Reset TMPDIR before creating extractDir so t.TempDir() doesn't write to internalTmp
	t.Setenv("TMPDIR", "")
	extractDir := t.TempDir()
	err = meshkitOci.UnCompressOCIArtifact(tarName, extractDir)
	assert.NoError(t, err)

	// Assert that internal_tmp inside extracted artifact does not contain meshery-model-build-*
	entries, err := os.ReadDir(filepath.Join(extractDir, "internal_tmp"))
	if err == nil {
		for _, e := range entries {
			assert.False(t, strings.HasPrefix(e.Name(), "meshery-model-build-"), "found staging directory leak in artifact: %s", e.Name())
		}
	}

	// Assert staging dir was cleaned up from model source
	sourceEntries, err := os.ReadDir(internalTmp)
	assert.NoError(t, err)
	assert.Empty(t, sourceEntries, "staging directory was not cleaned up from internal TMPDIR")
}

func TestModelBuild_Staging_ReadOnlySourceDirectory(t *testing.T) {
	cmdModelBuildFlagsProvided.Path = "."

	// Isolate TMPDIR so we can verify staging cleanup of 0555 and 0755 directories
	stagingTmpParent := t.TempDir()
	t.Setenv("TMPDIR", stagingTmpParent)

	createFreshModelCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "model",
			Short: "Model management",
		}
		ModelCmd.Flags().VisitAll(func(flag *pflag.Flag) {
			cmd.Flags().AddFlag(flag)
		})
		return cmd
	}

	createFreshBuildCmd := func() *cobra.Command {
		freshBuildCmd := &cobra.Command{
			Use:     buildModelCmd.Use,
			Short:   buildModelCmd.Short,
			Long:    buildModelCmd.Long,
			Example: buildModelCmd.Example,
			PreRunE: buildModelCmd.PreRunE,
			RunE:    buildModelCmd.RunE,
		}
		buildModelCmd.Flags().VisitAll(func(flag *pflag.Flag) {
			freshBuildCmd.Flags().AddFlag(flag)
		})
		cmd := createFreshModelCmd()
		mesheryctlflags.InitValidators(freshBuildCmd)
		cmd.AddCommand(freshBuildCmd)
		return cmd
	}

	const modelName = "test-case-model-build-readonly-dir"
	tarName := modelName + ".tar"

	cleanup := func() {
		// Restore write permissions so cleanup succeeds
		_ = filepath.WalkDir(modelName, func(p string, d fs.DirEntry, err error) error {
			if err == nil && d.IsDir() {
				_ = os.Chmod(p, 0o755)
			}
			return nil
		})
		_ = os.RemoveAll(modelName)
		_ = os.RemoveAll(tarName)
	}
	cleanup()
	defer cleanup()

	// Initialize model with v1.0.0
	setupCmd := createFreshModelCmd()
	setupCmd.AddCommand(initModelCmd)
	mesheryctlflags.InitValidators(initModelCmd)
	setupCmd.SetArgs([]string{"init", modelName, "--version", "v1.0.0"})
	buffInit := utils.SetupMeshkitLoggerTesting(t, false)
	setupCmd.SetOut(buffInit)
	if err := setupCmd.Execute(); err != nil {
		t.Fatalf("failed to init test model: %v", err)
	}

	// Create .git directory to trigger staging
	gitDir := filepath.Join(modelName, ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte("sample git config"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Create a subdirectory with read-only permissions (0555) containing a file
	readOnlyDir := filepath.Join(modelName, "readonly-components")
	if err := os.MkdirAll(readOnlyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	compFile := filepath.Join(readOnlyDir, "component.json")
	if err := os.WriteFile(compFile, []byte(`{"name":"test-readonly-component"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	// Change directory mode to 0555 (read & execute only, no write bit)
	if err := os.Chmod(readOnlyDir, 0o555); err != nil {
		t.Fatal(err)
	}

	// Create a standard subdirectory with normal permissions (0755) containing a file
	standardDir := filepath.Join(modelName, "standard-components")
	if err := os.MkdirAll(standardDir, 0o755); err != nil {
		t.Fatal(err)
	}
	stdCompFile := filepath.Join(standardDir, "component.json")
	if err := os.WriteFile(stdCompFile, []byte(`{"name":"test-standard-component"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	// Run mesheryctl model build
	cmd := createFreshBuildCmd()
	cmd.SetArgs([]string{"build", modelName})
	buff := utils.SetupMeshkitLoggerTesting(t, false)
	cmd.SetOut(buff)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("model build failed with read-only directory: %v", err)
	}

	// Verify tar exists
	info, err := os.Stat(tarName)
	assert.NoError(t, err)
	assert.False(t, info.IsDir())

	// Assert staging tree was completely removed by cleanup() despite containing 0555 directory
	stagingEntries, err := os.ReadDir(stagingTmpParent)
	assert.NoError(t, err)
	assert.Empty(t, stagingEntries, "staging directory was not cleaned up from TMPDIR after build")

	// Reset TMPDIR before creating extractDir
	t.Setenv("TMPDIR", "")
	extractDir := t.TempDir()
	err = meshkitOci.UnCompressOCIArtifact(tarName, extractDir)
	assert.NoError(t, err)

	extractedReadOnlyFile, err := os.ReadFile(filepath.Join(extractDir, "readonly-components", "component.json"))
	assert.NoError(t, err)
	assert.Equal(t, `{"name":"test-readonly-component"}`, string(extractedReadOnlyFile))

	extractedStandardFile, err := os.ReadFile(filepath.Join(extractDir, "standard-components", "component.json"))
	assert.NoError(t, err)
	assert.Equal(t, `{"name":"test-standard-component"}`, string(extractedStandardFile))

	// Verify tar layer headers preserve 0555 mode for readonly and 0755 for standard
	img, err := tarball.ImageFromPath(tarName, nil)
	assert.NoError(t, err)
	layers, err := img.Layers()
	assert.NoError(t, err)
	assert.NotEmpty(t, layers)
	layerReader, err := layers[0].Uncompressed()
	assert.NoError(t, err)
	defer func() { _ = layerReader.Close() }()
	tr := archiveTar.NewReader(layerReader)
	var foundReadonlyDir, foundStandardDir bool
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		if hdr.Name == "readonly-components" || hdr.Name == "readonly-components/" {
			foundReadonlyDir = true
			assert.Equal(t, int64(0o555), hdr.Mode&0o777)
		}
		if hdr.Name == "standard-components" || hdr.Name == "standard-components/" {
			foundStandardDir = true
			assert.Equal(t, int64(0o755), hdr.Mode&0o777)
		}
	}
	assert.True(t, foundReadonlyDir, "readonly-components directory not found in tar layer")
	assert.True(t, foundStandardDir, "standard-components directory not found in tar layer")
}
