package model

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelBuild(t *testing.T) {
	utils.SetupContextEnv(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	setupHookModelInit := func(modelInitArgs ...string) func() {
		return func() {
			// TODO replace ModelExpCmd with  ModelCmd
			cmd := ModelExpCmd
			// cmd := ModelCmd
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
			errs := make([]error, 0, 1)
			for _, dir := range dirs {
				if err := os.RemoveAll(dir); err != nil {
					errs = append(errs, err)
				}
			}
			if len(errs) > 0 {
				t.Fatal(errors.Join(errs...))
			}

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
			Args:             []string{"build", "test-case-aws-lambda-controller", "--version", "v0.1.0"},
			ExpectError:      false,
			ExpectedResponse: "model.build.from-model-name-version.golden",
			ExpectedFiles: []string{
				"test-case-aws-lambda-controller-v0-1-0.tar",
			},
			SetupHooks: []func(){
				setupHookModelInit("init", "test-case-aws-lambda-controller", "--version", "v0.1.0"),
			},
			CleanupHooks: []func(){
				cleanUpHookRemoveDirsAndFiles(
					"test-case-aws-lambda-controller",
					"test-case-aws-lambda-controller-v0-1-0.tar",
				),
			},
		},
		{
			Name:             "model build no params",
			Args:             []string{"build"},
			ExpectError:      true,
			ExpectedResponse: "model.build.error.invalid-amount-of-args.golden",
		},
		{
			Name:             "model build folder does not exist",
			Args:             []string{"build", "aws-ec2-controller", "--version", "v0.1.0", "--path", "./not_existing_folder"},
			ExpectError:      true,
			ExpectedResponse: "model.build.error.folder-does-not-exist.golden",
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
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
