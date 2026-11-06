package system

import (
	"flag"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update", false, "update golden files")

// This is an Integration test
func TestPreflightCmdIntegration(t *testing.T) {
	// skipping this integration test with --short flag
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	setupContextTestEnv(t)
	tests := []utils.CmdTestInput{
		{
			Name:             "Run preflight check",
			Args:             []string{"check", "--pre"},
			ExpectedResponse: "check.output.golden",
		},
		{
			Name:             "Run preflight check",
			Args:             []string{"check", "--preflight"},
			ExpectedResponse: "check.output.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// setting up log to grab logs
			buf := setupSystemOutCmdTest(t)
			defer func() {
				buf.Reset()
				// resetCmdFlags(SystemCmd, t)
			}()

			SystemCmd.SetArgs(tt.Args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			output := buf.String()
			actualResponse := output

			// get current directory
			_, filename, _, ok := runtime.Caller(0)
			if !ok {
				t.Fatal("Not able to get current working directory")
			}

			currDir := filepath.Dir(filename)
			testdataDir := filepath.Join(currDir, "testdata")
			// testdataDir := filepath.Join(filepath.Dir(filename), tf.dir, tf.name)
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, filepath.Join(testdataDir, "check"))

			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			assert.Equal(t, expectedResponse, actualResponse)
		})
		t.Log("PreflightCmdIntegration Test Passed")
	}
}

// TestSystemCheckFlagValidation verifies that each flag struct tag matches a registered Cobra flag on checkCmd.
func TestSystemCheckFlagValidation(t *testing.T) {
	mesheryctlflags.InitValidators(checkCmd)

	// Verify valid flags pass validation
	err := mesheryctlflags.ValidateCmdFlags(checkCmd, &systemCheckFlags)
	assert.NoError(t, err)

	// Verify each flag struct tag matches a registered Cobra flag on checkCmd
	structType := reflect.TypeOf(systemCheckFlags)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("json")
		if tag != "" && tag != "-" {
			flagName := strings.SplitN(tag, ",", 2)[0]
			assert.NotNil(t, checkCmd.Flags().Lookup(flagName), "expected flag --%s to be registered on checkCmd", flagName)
		}
	}
}
