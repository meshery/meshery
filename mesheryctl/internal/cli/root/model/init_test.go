package model

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestModelInit(t *testing.T) {
	// Shared cleanup function for test directories and artifacts
	cleanupTestArtifacts := func(dirs []string) {
		for _, dir := range dirs {
			_ = os.RemoveAll(dir)
		}
	}

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

	mesheryctlflags.InitValidators(initModelCmd)
	// Helper function to create fresh commands for init tests
	createFreshCommands := func() *cobra.Command {
		cmd := createFreshModelCmd()
		cmd.AddCommand(initModelCmd)
		return cmd
	}

	// Test constants
	const (
		initTestEC2Controller    = "test-case-aws-ec2-controller"
		initTestDynamoController = "test-case-aws-dynamodb-controller"
		initTestVersion          = "v0.1.0"
		initTestInvalidVersion   = "v1.0.0"
	)

	// Clean up any existing test directories before running tests
	cleanupDirs := []string{
		initTestEC2Controller,
		initTestDynamoController,
		"test_case_some_other_custom_dir",
		"test-case_aws-ec2-controller", // underscore variant used in invalid model name test
	}
	cleanupTestArtifacts(cleanupDirs)
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
	tests := []struct {
		Name               string
		Args               []string
		SetupHook          func()
		ExpectError        bool
		ExpectedResponse   string
		ExpectedDirs       []string
		ExpectedFiles      []string
		AfterTestRemoveDir string
		ExpectedError      error
		IsOutputGolden     bool
	}{
		// NOTE:
		// we need this test with full params on the first place,
		// to prevent side effects of using same command object in model build test.
		//
		// TODO: think about how to fix this.
		{
			Name:             "given all default parameters when model init model is initialized",
			Args:             []string{"init", initTestEC2Controller, "--version", initTestVersion, "--path", ".", "--output-format", "json"},
			ExpectError:      false,
			ExpectedResponse: "model.init.aws-ec2-controller.output.golden",
			ExpectedDirs: []string{
				initTestEC2Controller,
				filepath.Join(initTestEC2Controller, initTestVersion),
				filepath.Join(initTestEC2Controller, initTestVersion, "components"),
				filepath.Join(initTestEC2Controller, initTestVersion, "connections"),
				filepath.Join(initTestEC2Controller, initTestVersion, "credentials"),
				filepath.Join(initTestEC2Controller, initTestVersion, "relationships"),
			},
			ExpectedFiles: []string{
				filepath.Join(initTestEC2Controller, initTestVersion, "model.json"),
				filepath.Join(initTestEC2Controller, initTestVersion, "components/component.json"),
				filepath.Join(initTestEC2Controller, initTestVersion, "relationships/relationship.json"),
			},
			AfterTestRemoveDir: initTestEC2Controller,
		},
		{
			Name:             "given default parameters when model init model is initialized",
			Args:             []string{"init", initTestEC2Controller},
			ExpectError:      false,
			ExpectedResponse: "model.init.aws-ec2-controller.output.golden",
			ExpectedDirs: []string{
				initTestEC2Controller,
				filepath.Join(initTestEC2Controller, initTestVersion),
				filepath.Join(initTestEC2Controller, initTestVersion, "components"),
				filepath.Join(initTestEC2Controller, initTestVersion, "connections"),
				filepath.Join(initTestEC2Controller, initTestVersion, "credentials"),
				filepath.Join(initTestEC2Controller, initTestVersion, "relationships"),
			},
			ExpectedFiles: []string{
				filepath.Join(initTestEC2Controller, initTestVersion, "model.json"),
				filepath.Join(initTestEC2Controller, initTestVersion, "components/component.json"),
				filepath.Join(initTestEC2Controller, initTestVersion, "relationships/relationship.json"),
			},
			AfterTestRemoveDir: initTestEC2Controller,
		},
		{
			Name:             "given output format yaml when model init model is initialized in yaml format",
			Args:             []string{"init", initTestDynamoController, "--output-format", "yaml"},
			ExpectError:      false,
			ExpectedResponse: "model.init.aws-dynamodb-controller-in-yaml.output.golden",
			ExpectedDirs: []string{
				initTestDynamoController,
				filepath.Join(initTestDynamoController, initTestVersion),
				filepath.Join(initTestDynamoController, initTestVersion, "components"),
				filepath.Join(initTestDynamoController, initTestVersion, "connections"),
				filepath.Join(initTestDynamoController, initTestVersion, "credentials"),
				filepath.Join(initTestDynamoController, initTestVersion, "relationships"),
			},
			ExpectedFiles: []string{
				filepath.Join(initTestDynamoController, initTestVersion, "model.yaml"),
				filepath.Join(initTestDynamoController, initTestVersion, "components/component.yaml"),
				filepath.Join(initTestDynamoController, initTestVersion, "relationships/relationship.yaml"),
			},
			AfterTestRemoveDir: initTestDynamoController,
		},
		// Added --output-format json in this test because somehow
		// the --output-format yaml from the previous test case is propagated to this test case
		// which is only the behaviour inside the test.
		// TODO think about how to reset the flags between the test cases.
		{
			Name:             "given custom path and version when model init model is initialized with custom path and version specified",
			Args:             []string{"init", initTestEC2Controller, "--path", "test_case_some_custom_dir/subdir/one_more_subdir", "--version", "v1.2.3", "--output-format", "json"},
			ExpectError:      false,
			ExpectedResponse: "model.init.custom-dir.aws-ec2-controller.output.golden",
			ExpectedDirs: []string{
				"test_case_some_custom_dir/subdir/one_more_subdir/" + initTestEC2Controller,
				"test_case_some_custom_dir/subdir/one_more_subdir/" + initTestEC2Controller + "/v1.2.3",
				"test_case_some_custom_dir/subdir/one_more_subdir/" + initTestEC2Controller + "/v1.2.3/components",
				"test_case_some_custom_dir/subdir/one_more_subdir/" + initTestEC2Controller + "/v1.2.3/connections",
				"test_case_some_custom_dir/subdir/one_more_subdir/" + initTestEC2Controller + "/v1.2.3/credentials",
				"test_case_some_custom_dir/subdir/one_more_subdir/" + initTestEC2Controller + "/v1.2.3/relationships",
			},
			ExpectedFiles: []string{
				"test_case_some_custom_dir/subdir/one_more_subdir/" + initTestEC2Controller + "/v1.2.3/model.json",
				"test_case_some_custom_dir/subdir/one_more_subdir/" + initTestEC2Controller + "/v1.2.3/components/component.json",
				"test_case_some_custom_dir/subdir/one_more_subdir/" + initTestEC2Controller + "/v1.2.3/relationships/relationship.json",
			},
			AfterTestRemoveDir: "test_case_some_custom_dir",
		},
		{
			Name:             "given custom relative to current directory path when model init then model is initialized in custom relative path",
			Args:             []string{"init", initTestEC2Controller, "--path", "./test_case_some_custom_dir/subdir/one_more_subdir", "--version", "v1.2.3", "--output-format", "json"},
			ExpectError:      false,
			ExpectedResponse: "model.init.custom-dir.aws-ec2-controller.output.golden",
			// do not need to check all dirs and files here, as we tested it in previous test case,
			// just check the main directory is correct
			ExpectedDirs: []string{
				"./test_case_some_custom_dir/subdir/one_more_subdir/" + initTestEC2Controller,
			},
			AfterTestRemoveDir: "./test_case_some_custom_dir",
		},
		{
			Name:             "given custom relative to parent directory path when model init then model is initialized in custom relative to parent directory path",
			Args:             []string{"init", initTestEC2Controller, "--path", "../test_case_some_custom_dir/subdir/one_more_subdir", "--version", "v1.2.3", "--output-format", "json"},
			ExpectError:      false,
			ExpectedResponse: "model.init.custom-relative-parent-dir.aws-ec2-controller.output.golden",
			// do not need to check all dirs and files here, as we tested it in previous test case,
			// just check the main directory is correct
			ExpectedDirs: []string{
				"../test_case_some_custom_dir/subdir/one_more_subdir/" + initTestEC2Controller,
			},
			AfterTestRemoveDir: "../test_case_some_custom_dir",
		},
		{
			Name:             "given trailing folder separator in the path when model init then model is initialized in the path without trailing separator",
			Args:             []string{"init", initTestEC2Controller, "--path", "test_case_some_other_custom_dir/with/trailing/separator////", "--version", "v1.2.3"},
			ExpectError:      false,
			ExpectedResponse: "model.init.custom-dir-2.aws-ec2-controller.output.golden",
			// do not need to check all dirs and files here, as we tested it in previous test case,
			// just check the main directory is correct
			ExpectedDirs: []string{
				"test_case_some_other_custom_dir/with/trailing/separator/" + initTestEC2Controller,
			},
			AfterTestRemoveDir: "test_case_some_other_custom_dir",
		},
		{
			Name: "given existing model/version folder when model init then throw error",
			Args: []string{"init", initTestEC2Controller, "--path", ".", "--version", initTestInvalidVersion},
			SetupHook: func() {
				err := os.MkdirAll(filepath.Join(initTestEC2Controller, initTestInvalidVersion), initModelDirPerm)
				if err != nil {
					t.Fatal("error in SetupHook when creating folder", err)
				}
			},
			ExpectError:        true,
			ExpectedResponse:   "",
			AfterTestRemoveDir: "./" + initTestEC2Controller,
			IsOutputGolden:     false,
			ExpectedError:      ErrModelInitFromString(fmt.Sprintf(errInitFolderExists, filepath.Join(initTestEC2Controller, initTestInvalidVersion))),
		},
		{
			Name:             "given invalid version format when model init then throw error",
			Args:             []string{"init", initTestEC2Controller, "--version", "1.2"},
			ExpectError:      true,
			ExpectedResponse: "",
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --version '1.2': version must be in format vX.X.X")),
		},
		{
			Name:             "given invalid output format when model init then throw error",
			Args:             []string{"init", "test-case-aws-ec2-controller", "--output-format", "protobuf"},
			ExpectError:      true,
			ExpectedResponse: "",
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --output-format 'protobuf': valid values are json yaml")),
		},
		{
			Name:             "given no model name provided when model init then throw error",
			Args:             []string{"init", "--output-format", "json", "--version", "v0.1.0"},
			ExpectError:      true,
			ExpectedResponse: "",
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("%s", errInitOneArg)),
		},
		{
			Name:             "given too many arguments when model init then throw error",
			Args:             []string{"init", "test-case-aws-ec2-controller", "test-case-aws-dynamodb-controller", "--output-format", "json", "--version", "v0.1.0"},
			ExpectError:      true,
			ExpectedResponse: "",
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("%s", errInitOneArg)),
		},
		{
			Name:             "given invalid model name when model init then throw error",
			Args:             []string{"init", "test-case_aws-ec2-controller", "--output-format", "json", "--version", "v0.1.0"},
			ExpectError:      true,
			ExpectedResponse: "",
			IsOutputGolden:   false,
			ExpectedError:    ErrModelInit(fmt.Errorf("%s", "invalid model name: name must match pattern ^[a-z0-9-]+$")),
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			// Clean up any lingering test artifacts before this subtest starts
			cleanupTestArtifacts(cleanupDirs)

			defer func() {
				// clean up created folders on any test case outcome
				if tc.AfterTestRemoveDir != "" {
					cleanupTestArtifacts([]string{tc.AfterTestRemoveDir})
					t.Log("removed created folders")
				}
			}()
			if tc.SetupHook != nil {
				tc.SetupHook()
			}
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tc.ExpectedResponse, testdataDir)
			buff := utils.SetupMeshkitLoggerTesting(t, false)
			// Create fresh commands using helper function
			cmd := createFreshCommands()
			defer utils.ResetCommandFlags(cmd, t)
			defer utils.ResetCommandFlags(initModelCmd, t)

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

			if len(tc.ExpectedDirs) > 0 {
				for _, dir := range tc.ExpectedDirs {
					info, err := os.Stat(dir)
					assert.NoError(t, err)
					if err == nil {
						assert.True(t, info.IsDir())
					}

				}
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
