package model

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
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
		buildTestVersion             = "v0.1.0"
		buildTestNonExistentFolder   = "not_existing_folder"
	)

	// Clean up any existing test directories before running tests
	cleanupDirs := []string{
		buildTestDynamoController,
		buildTestDynamoControllerGbx,
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
		ExpectedContains []string
		ExpectedFiles    []string
		CleanupHooks     []func()
		ExpectedError    error `default:"nil"`
	}{
		{
			Name:             "model build from model name and version",
			Args:             []string{"build", "test-case-model-build-aws-lambda-controller/v0.1.0"},
			ExpectError:      false,
			ExpectedContains: []string{
				"Building meshery model from path test-case-model-build-aws-lambda-controller/v0.1.0",
				"Saving OCI artifact as test-case-model-build-aws-lambda-controller-v0-1-0.tar",
			},
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
			Args:             []string{"build", buildTestDynamoController},
			ExpectError:      false,
			ExpectedContains: []string{
				"Building meshery model from path test-case-model-build-aws-dynamodb-controller",
				"Saving OCI artifact as test-case-model-build-aws-dynamodb-controller.tar",
			},
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
			Name:             "model build from model name only (no version) with slash in the end",
			Args:             []string{"build", buildTestDynamoController + "/"},
			ExpectError:      false,
			ExpectedContains: []string{
				"Building meshery model from path test-case-model-build-aws-dynamodb-controller",
				"Saving OCI artifact as test-case-model-build-aws-dynamodb-controller.tar",
			},
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
			Name:             "model build no params",
			Args:             []string{"build"},
			ExpectError:      true,
			ExpectedError:    ErrModelBuildFromStrings(errBuildUsage),
		},
		{
			Name:             "model build wrong input param format",
			Args:             []string{"build", buildTestEC2Controller + "/" + buildTestVersion + "/smthelse"},
			ExpectError:      true,
			ExpectedError:    ErrModelBuildFromStrings(errBuildUsage),
		},
		{
			Name:             "model build from model name only (no version) not supporting multiple versions",
			Args:             []string{"build", buildTestDynamoControllerGbx},
			ExpectError:      true,
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
			ExpectedError:  ErrModelBuildFromStrings(errBuildUsage, errBuildMultiVersionNotSupported),
		},
		{
			Name:             "model build folder does not exist",
			Args:             []string{"build", buildTestEC2Controller + "/" + buildTestVersion, "--path", "./" + buildTestNonExistentFolder},
			ExpectError:      true,
			ExpectedError:    ErrModelBuildFromStrings(errBuildUsage, fmt.Sprintf(errBuildFolderNotFound, filepath.Join(buildTestNonExistentFolder, buildTestEC2Controller, buildTestVersion))),
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
			buff := utils.SetupMeshkitLoggerTesting(t, false)
			// Create fresh commands using helper function
			cmd := createFreshCommands()
			cmd.SetArgs(tc.Args)
			cmd.SetOut(buff)
			err := cmd.Execute()
			if err != nil {
				// if we're supposed to get an error
				if tc.ExpectError {
					assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(tc.ExpectedError))
					assert.Equal(t, errors.GetCode(err), errors.GetCode(tc.ExpectedError))
					assert.Equal(t, errors.GetLDescription(err), errors.GetLDescription(tc.ExpectedError))
					return

				}
				t.Fatal(err)
			}
			// response being printed in console
			actualResponse := buff.String()
			for _, expected := range tc.ExpectedContains {
				assert.Contains(t, actualResponse, expected)
			}

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
