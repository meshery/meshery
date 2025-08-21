package model

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestModelBuild(t *testing.T) {
	// Shared cleanup function for test directories and artifacts
	cleanupTestArtifacts := func(dirs []string) {
		for _, dir := range dirs {
			os.RemoveAll(dir)
			os.RemoveAll(dir + ".tar")
		}
	}

	// Clean up any existing test directories before running tests
	cleanupDirs := []string{
		"test-case-model-build-aws-dynamodb-controller",
		"test-case-model-build-aws-dynamodb-controller-gbxter34",
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
		cmd.AddCommand(freshBuildCmd)
		return cmd
	}

	setupHookModelInit := func(modelInitArgs ...string) func() {
		return func() {
			// Use the shared helper to create a fresh ModelCmd
			cmd := createFreshModelCmd()
			cmd.AddCommand(initModelCmd)
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

	tests := []struct {
		Name             string
		Args             []string
		SetupHooks       []func()
		ExpectError      bool
		ExpectedResponse string
		ExpectedFiles    []string
		CleanupHooks     []func()
	}{
		{
			Name:             "model build from model name and version",
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
			Name:             "model build from model name only (no version)",
			Args:             []string{"build", "test-case-model-build-aws-dynamodb-controller"},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				"test-case-model-build-aws-dynamodb-controller.tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", "test-case-model-build-aws-dynamodb-controller", "--version", "v0.1.0"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					"test-case-model-build-aws-dynamodb-controller",
					"test-case-model-build-aws-dynamodb-controller.tar",
				),
			},
		},
		{
			Name:             "model build from model name only (no version) with slash in the end",
			Args:             []string{"build", "test-case-model-build-aws-dynamodb-controller/"},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-only.golden",
			ExpectedFiles: []string{
				"test-case-model-build-aws-dynamodb-controller.tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", "test-case-model-build-aws-dynamodb-controller", "--version", "v0.1.0"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					"test-case-model-build-aws-dynamodb-controller",
					"test-case-model-build-aws-dynamodb-controller.tar",
				),
			},
		},
		{
			Name:             "model build no params",
			Args:             []string{"build"},
			ExpectError:      true,
			ExpectedResponse: "model.build.error.usage.golden",
		},
		{
			Name:             "model build wrong input param format",
			Args:             []string{"build", "aws-ec2-controller/v0.1.0/smthelse"},
			ExpectError:      true,
			ExpectedResponse: "model.build.error.usage.golden",
		},
		{
			Name:             "model build from model name only (no version) not supporting multiple versions",
			Args:             []string{"build", "test-case-model-build-aws-dynamodb-controller-gbxter34"},
			ExpectError:      true,
			ExpectedResponse: "model.build.error.not-supporting-multiple-versions-build.golden",
			SetupHooks: []func(){
				setupHookModelInit("init", "test-case-model-build-aws-dynamodb-controller-gbxter34", "--version", "v0.1.0"),
				setupHookModelInit("init", "test-case-model-build-aws-dynamodb-controller-gbxter34", "--version", "v0.1.1"),
				setupHookModelInit("init", "test-case-model-build-aws-dynamodb-controller-gbxter34", "--version", "v0.1.2"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					"test-case-model-build-aws-dynamodb-controller-gbxter34",
				),
			},
		},
		{
			Name:             "model build folder does not exist",
			Args:             []string{"build", "aws-ec2-controller/v0.1.0", "--path", "./not_existing_folder"},
			ExpectError:      true,
			ExpectedResponse: "model.build.error.folder-does-not-exist.golden",
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			// Clean up any lingering test artifacts before this subtest starts
			cleanupTestArtifacts(cleanupDirs)

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
					expectedResponse := golden.Load()

					utils.Equals(t, expectedResponse, err.Error())
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
