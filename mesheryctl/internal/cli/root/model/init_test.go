package model

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
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
		Name             string
		Args             []string
		ExpectError      bool
		ExpectedResponse string
		ExpectedDirs     []string
		ExpectedFiles    []string
	}{
		{
			Name:             "model init with default params",
			Args:             []string{"init", "aws-ec2-controller"},
			ExpectError:      false,
			ExpectedResponse: "model.init.aws-ec2-controller.output.golden",
			ExpectedDirs: []string{
				"aws-ec2-controller",
				"aws-ec2-controller/v0.1.0",
				"aws-ec2-controller/v0.1.0/components",
				"aws-ec2-controller/v0.1.0/connections",
				"aws-ec2-controller/v0.1.0/credentials",
				"aws-ec2-controller/v0.1.0/relationships",
			},
			ExpectedFiles: []string{
				"aws-ec2-controller/v0.1.0/model.json",
				"aws-ec2-controller/v0.1.0/components/component.json",
				"aws-ec2-controller/v0.1.0/connections/connection.json",
				"aws-ec2-controller/v0.1.0/relationships/relationship.json",
			},
		},
		{
			Name:             "model init with custom path and version",
			Args:             []string{"init", "aws-ec2-controller", "--path", "some_custom_dir/subdir/one_more_subdir", "--version", "v1.2.3"},
			ExpectError:      false,
			ExpectedResponse: "model.init.custom-dir.aws-ec2-controller.output.golden",
			ExpectedDirs: []string{
				"some_custom_dir/subdir/one_more_subdir/aws-ec2-controller",
				"some_custom_dir/subdir/one_more_subdir/aws-ec2-controller/v1.2.3",
				"some_custom_dir/subdir/one_more_subdir/aws-ec2-controller/v1.2.3/components",
				"some_custom_dir/subdir/one_more_subdir/aws-ec2-controller/v1.2.3/connections",
				"some_custom_dir/subdir/one_more_subdir/aws-ec2-controller/v1.2.3/credentials",
				"some_custom_dir/subdir/one_more_subdir/aws-ec2-controller/v1.2.3/relationships",
			},
			ExpectedFiles: []string{
				"some_custom_dir/subdir/one_more_subdir/aws-ec2-controller/v1.2.3/model.json",
				"some_custom_dir/subdir/one_more_subdir/aws-ec2-controller/v1.2.3/components/component.json",
				"some_custom_dir/subdir/one_more_subdir/aws-ec2-controller/v1.2.3/connections/connection.json",
				"some_custom_dir/subdir/one_more_subdir/aws-ec2-controller/v1.2.3/relationships/relationship.json",
			},
		},
		{
			Name:             "model init with invalid version format",
			Args:             []string{"init", "aws-ec2-controller", "--version", "1.2"},
			ExpectError:      true,
			ExpectedResponse: "model.init.invalid-version-format.output.golden",
		},
		// TODO test not covered branches and corner cases
		// TODO remove created folders during tests
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tc.ExpectedResponse, testdataDir)
			buff := utils.SetupMeshkitLoggerTesting(t, false)
			cmd := ModelCmd
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
		t.Log("model init test finished")
	}
}
