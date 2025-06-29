package model

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelInit(t *testing.T) {
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
	}{
		// NOTE:
		// we need this test with full params on the first place,
		// to prevent side effects of using same command object in model build test.
		//
		// TODO: think about how to fix this.
		{
			Name:             "model init with all default params",
			Args:             []string{"init", "test-case-aws-ec2-controller", "--version", "v0.1.0", "--path", ".", "--output-format", "json"},
			ExpectError:      false,
			ExpectedResponse: "model.init.aws-ec2-controller.output.golden",
			ExpectedDirs: []string{
				"test-case-aws-ec2-controller",
				"test-case-aws-ec2-controller/v0.1.0",
				"test-case-aws-ec2-controller/v0.1.0/components",
				"test-case-aws-ec2-controller/v0.1.0/connections",
				"test-case-aws-ec2-controller/v0.1.0/credentials",
				"test-case-aws-ec2-controller/v0.1.0/relationships",
			},
			ExpectedFiles: []string{
				"test-case-aws-ec2-controller/v0.1.0/model.json",
				"test-case-aws-ec2-controller/v0.1.0/components/component.json",
				"test-case-aws-ec2-controller/v0.1.0/connections/connection.json",
				"test-case-aws-ec2-controller/v0.1.0/relationships/relationship.json",
			},
			AfterTestRemoveDir: "test-case-aws-ec2-controller",
		},
		{
			Name:             "model init with default params",
			Args:             []string{"init", "test-case-aws-ec2-controller"},
			ExpectError:      false,
			ExpectedResponse: "model.init.aws-ec2-controller.output.golden",
			ExpectedDirs: []string{
				"test-case-aws-ec2-controller",
				"test-case-aws-ec2-controller/v0.1.0",
				"test-case-aws-ec2-controller/v0.1.0/components",
				"test-case-aws-ec2-controller/v0.1.0/connections",
				"test-case-aws-ec2-controller/v0.1.0/credentials",
				"test-case-aws-ec2-controller/v0.1.0/relationships",
			},
			ExpectedFiles: []string{
				"test-case-aws-ec2-controller/v0.1.0/model.json",
				"test-case-aws-ec2-controller/v0.1.0/components/component.json",
				"test-case-aws-ec2-controller/v0.1.0/connections/connection.json",
				"test-case-aws-ec2-controller/v0.1.0/relationships/relationship.json",
			},
			AfterTestRemoveDir: "test-case-aws-ec2-controller",
		},
		{
			Name:             "model init with yaml output format",
			Args:             []string{"init", "test-case-aws-dynamodb-controller", "--output-format", "yaml"},
			ExpectError:      false,
			ExpectedResponse: "model.init.aws-dynamodb-controller-in-yaml.output.golden",
			ExpectedDirs: []string{
				"test-case-aws-dynamodb-controller",
				"test-case-aws-dynamodb-controller/v0.1.0",
				"test-case-aws-dynamodb-controller/v0.1.0/components",
				"test-case-aws-dynamodb-controller/v0.1.0/connections",
				"test-case-aws-dynamodb-controller/v0.1.0/credentials",
				"test-case-aws-dynamodb-controller/v0.1.0/relationships",
			},
			ExpectedFiles: []string{
				"test-case-aws-dynamodb-controller/v0.1.0/model.yaml",
				"test-case-aws-dynamodb-controller/v0.1.0/components/component.yaml",
				"test-case-aws-dynamodb-controller/v0.1.0/connections/connection.yaml",
				"test-case-aws-dynamodb-controller/v0.1.0/relationships/relationship.yaml",
			},
			AfterTestRemoveDir: "test-case-aws-dynamodb-controller",
		},
		// Added --output-format json in this test because somehow
		// the --output-format yaml from the previous test case is propagated to this test case
		// which is only the behaviour inside the test.
		// TODO think about how to reset the flags between the test cases.
		{
			Name:             "model init with custom path and version",
			Args:             []string{"init", "test-case-aws-ec2-controller", "--path", "test_case_some_custom_dir/subdir/one_more_subdir", "--version", "v1.2.3", "--output-format", "json"},
			ExpectError:      false,
			ExpectedResponse: "model.init.custom-dir.aws-ec2-controller.output.golden",
			ExpectedDirs: []string{
				"test_case_some_custom_dir/subdir/one_more_subdir/test-case-aws-ec2-controller",
				"test_case_some_custom_dir/subdir/one_more_subdir/test-case-aws-ec2-controller/v1.2.3",
				"test_case_some_custom_dir/subdir/one_more_subdir/test-case-aws-ec2-controller/v1.2.3/components",
				"test_case_some_custom_dir/subdir/one_more_subdir/test-case-aws-ec2-controller/v1.2.3/connections",
				"test_case_some_custom_dir/subdir/one_more_subdir/test-case-aws-ec2-controller/v1.2.3/credentials",
				"test_case_some_custom_dir/subdir/one_more_subdir/test-case-aws-ec2-controller/v1.2.3/relationships",
			},
			ExpectedFiles: []string{
				"test_case_some_custom_dir/subdir/one_more_subdir/test-case-aws-ec2-controller/v1.2.3/model.json",
				"test_case_some_custom_dir/subdir/one_more_subdir/test-case-aws-ec2-controller/v1.2.3/components/component.json",
				"test_case_some_custom_dir/subdir/one_more_subdir/test-case-aws-ec2-controller/v1.2.3/connections/connection.json",
				"test_case_some_custom_dir/subdir/one_more_subdir/test-case-aws-ec2-controller/v1.2.3/relationships/relationship.json",
			},
			AfterTestRemoveDir: "test_case_some_custom_dir",
		},
		{
			Name:             "model init with custom relative to current directory path",
			Args:             []string{"init", "test-case-aws-ec2-controller", "--path", "./test_case_some_custom_dir/subdir/one_more_subdir", "--version", "v1.2.3", "--output-format", "json"},
			ExpectError:      false,
			ExpectedResponse: "model.init.custom-dir.aws-ec2-controller.output.golden",
			// do not need to check all dirs and files here, as we tested it in previous test case,
			// just check the main directory is correct
			ExpectedDirs: []string{
				"./test_case_some_custom_dir/subdir/one_more_subdir/test-case-aws-ec2-controller",
			},
			AfterTestRemoveDir: "./test_case_some_custom_dir",
		},
		{
			Name:             "model init with custom relative to parent directory path",
			Args:             []string{"init", "test-case-aws-ec2-controller", "--path", "../test_case_some_custom_dir/subdir/one_more_subdir", "--version", "v1.2.3", "--output-format", "json"},
			ExpectError:      false,
			ExpectedResponse: "model.init.custom-relative-parent-dir.aws-ec2-controller.output.golden",
			// do not need to check all dirs and files here, as we tested it in previous test case,
			// just check the main directory is correct
			ExpectedDirs: []string{
				"../test_case_some_custom_dir/subdir/one_more_subdir/test-case-aws-ec2-controller",
			},
			AfterTestRemoveDir: "../test_case_some_custom_dir",
		},
		{
			Name:             "model init with trailing folder separator in the path",
			Args:             []string{"init", "test-case-aws-ec2-controller", "--path", "test_case_some_other_custom_dir/with/trailing/separator////"},
			ExpectError:      false,
			ExpectedResponse: "model.init.custom-dir-2.aws-ec2-controller.output.golden",
			// do not need to check all dirs and files here, as we tested it in previous test case,
			// just check the main directory is correct
			ExpectedDirs: []string{
				"test_case_some_other_custom_dir/with/trailing/separator/test-case-aws-ec2-controller",
			},
			AfterTestRemoveDir: "test_case_some_other_custom_dir",
		},
		{
			Name: "model init fail if model/version folder esists",
			Args: []string{"init", "test-case-aws-ec2-controller", "--path", ".", "--version", "v1.0.0"},
			SetupHook: func() {
				err := os.MkdirAll("./test-case-aws-ec2-controller/v1.0.0", initModelDirPerm)
				if err != nil {
					t.Fatal("error in SetupHook when creating folder", err)
				}
			},
			ExpectError:        true,
			ExpectedResponse:   "model.init.folder-exists.output.golden",
			AfterTestRemoveDir: "./test-case-aws-ec2-controller",
		},
		{
			Name:             "model init with invalid version format",
			Args:             []string{"init", "test-case-aws-ec2-controller", "--version", "1.2"},
			ExpectError:      true,
			ExpectedResponse: "model.init.invalid-version-format.output.golden",
		},
		{
			Name:             "model init with invalid output format",
			Args:             []string{"init", "test-case-aws-ec2-controller", "--output-format", "protobuf"},
			ExpectError:      true,
			ExpectedResponse: "model.init.invalid-output-format.output.golden",
		},
		{
			Name:             "model init no model name",
			Args:             []string{"init", "--output-format", "json", "--version", "v0.1.0"},
			ExpectError:      true,
			ExpectedResponse: "model.init.wrong-number-of-arguments.output.golden",
		},
		{
			Name:             "model init too many arguments",
			Args:             []string{"init", "test-case-aws-ec2-controller", "test-case-aws-dynamodb-controller", "--output-format", "json", "--version", "v0.1.0"},
			ExpectError:      true,
			ExpectedResponse: "model.init.wrong-number-of-arguments.output.golden",
		},
		{
			Name:             "model init invalid model name (underscore)",
			Args:             []string{"init", "test-case_aws-ec2-controller", "--output-format", "json", "--version", "v0.1.0"},
			ExpectError:      true,
			ExpectedResponse: "model.init.invalid-model-name.output.golden",
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			defer func() {
				// clean up created folders on any test case outcome
				if tc.AfterTestRemoveDir != "" {
					if err := os.RemoveAll(tc.AfterTestRemoveDir); err != nil {
						t.Fatal(err)
					}
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
			// TODO replace ModelExpCmd with  ModelCmd
			cmd := ModelExpCmd
			// cmd := ModelCmd
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
