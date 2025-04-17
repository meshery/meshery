package model

import (
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
		{
			Name:             "model build from model name and version",
			Args:             []string{"build", "aws-ec2-controller", "--version", "v0.1.0"},
			ExpectError:      false,
			ExpectedResponse: "model.build.encouraging-message.golden",
		},
		{
			Name:             "model build no params",
			Args:             []string{"build"},
			ExpectError:      true,
			ExpectedResponse: "model.build.error.invalid-amount-of-args.golden",
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
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
		})
	}
}
